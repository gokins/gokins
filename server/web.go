package server

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gokins/gokins/util/httpex"

	"github.com/gin-gonic/gin"
	"github.com/gokins/core"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/route"
	"github.com/gokins/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
)

func runWeb() {
	defer func() {
		if err := recover(); err != nil {
			hbtp.Errorf("Web recover:%v", err)
		}
	}()
	comm.WebEgn = gin.Default()
	comm.WebEgn.Use(midUiHandle)
	err := comm.WebEgn.Run(comm.WebHost)
	if err != nil {
		logrus.Errorf("Web err:%v", err)
		//comm.HbtpEgn.Stop()
	}
	comm.Cancel()
	time.Sleep(time.Millisecond * 100)
}

func regApi() {
	if core.Debug {
		comm.WebEgn.Use(util.MidAccessAllowFun)
	}
	util.GinRegController(comm.WebEgn, &route.ApiController{})
	util.GinRegController(comm.WebEgn, &route.ArtifactController{})
	util.GinRegController(comm.WebEgn, &route.ArtPublicController{})
	util.GinRegController(comm.WebEgn, &route.LoginController{})
	util.GinRegController(comm.WebEgn, &route.UserController{})
	util.GinRegController(comm.WebEgn, &route.OrgController{})
	util.GinRegController(comm.WebEgn, &route.PipelineController{})
	util.GinRegController(comm.WebEgn, &route.PipelineVersionController{})
	util.GinRegController(comm.WebEgn, &route.RuntimeController{})
	util.GinRegController(comm.WebEgn, &route.YmlController{})
	util.GinRegController(comm.WebEgn, &route.TriggerController{})
	util.GinRegController(comm.WebEgn, &route.HookController{})
}
func midUiHandle(c *gin.Context) {
	c.Next()
	if c.Writer.Status() != http.StatusNotFound || c.Writer.Size() > 0 {
		return
	}
	pth := c.Request.URL.Path
	if !comm.Installed && !strings.HasPrefix(pth, "/gokinsui/") && pth != "/install" {
		httpex.ResMsgUrl(c, "未安装,跳转中...", "/install")
		return
	}
	r, err := getFile(pth[1:])
	if err != nil {
		r, err = getFile("index.html")
	}
	if err != nil {
		//c.String(404, "rdr err:"+err.Error())
		httpex.ResMsgUrl(c, "未找到内容,跳转中...", "/")
		return
	}
	rd, err := r.Open()
	if err != nil {
		//c.String(500, "open err:"+err.Error())
		httpex.ResMsgUrl(c, "内容有误,跳转中...", "/")
		return
	}
	defer rd.Close()
	c.Writer.Header().Set("Cache-Control", "max-age=360000000")

	ext := filepath.Ext(r.Name)
	if ext == ".html" {
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Writer.Header().Set("Content-Type", "text/html")
	} else if ext == ".css" {
		c.Writer.Header().Set("Content-Type", "text/css")
	} else if ext == ".js" {
		c.Writer.Header().Set("Content-Type", "application/javascript")
	} else if ext == ".svg" {
		c.Writer.Header().Set("Content-Type", "image/svg+xml")
	} else if ext == ".woff2" {
		//c.Writer.Header().Set("Content-Type", "image/svg+xml")
	} else if ext == ".ttf" || ext == ".ttc" {
		c.Writer.Header().Set("Content-Type", "application/x-font-ttf")
	}
	c.Status(200)
	bts := make([]byte, 1024)
	for !hbtp.EndContext(c) {
		n, err := rd.Read(bts)
		if n > 0 {
			c.Writer.Write(bts[:n])
		}
		if err != nil {
			break
		}
	}
}

var rder *zip.Reader

func getRdr() (*zip.Reader, error) {
	if rder != nil {
		return rder, nil
	}
	bts, err := base64.StdEncoding.DecodeString(comm.StaticPkg)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(bts)
	r, err := zip.NewReader(buf, buf.Size())
	if err != nil {
		return nil, err
	}
	rder = r
	return rder, nil
}
func getFile(pth string) (*zip.File, error) {
	if pth == "" {
		return nil, errors.New("param err")
	}
	//println("getFile:" + pth)
	r, err := getRdr()
	if err != nil {
		return nil, err
	}
	for _, f := range r.File {
		//println("f.name:", f.Name)
		if pth == f.Name {
			return f, nil
		}
	}
	return nil, errors.New("file not found")
}
