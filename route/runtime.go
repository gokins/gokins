package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/engine"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/service"
	"github.com/gokins-main/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"os"
	"path/filepath"
)

type RuntimeController struct{}

func (RuntimeController) GetPath() string {
	return "/api/runtime"
}
func (c *RuntimeController) Routes(g gin.IRoutes) {
	g.Use(service.MidUserCheck)
	g.POST("/build", util.GinReqParseJson(c.build))
	g.POST("/cancel", util.GinReqParseJson(c.cancel))
	g.POST("/logs", util.GinReqParseJson(c.logs))
}
func (RuntimeController) build(c *gin.Context, m *hbtp.Map) {
	bdid := m.GetString("buildId")
	if bdid == "" {
		c.String(500, "param err")
		return
	}
	v, ok := engine.Mgr.BuildEgn().Get(bdid)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	show, ok := v.Show()
	if !ok {
		c.String(404, "Not Found")
		return
	}
	c.JSON(200, show)
}
func (RuntimeController) cancel(c *gin.Context, m *hbtp.Map) {
	bdid := m.GetString("buildId")
	if bdid == "" {
		c.String(500, "param err")
		return
	}
	v, ok := engine.Mgr.BuildEgn().Get(bdid)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	v.Cancel()
	c.String(200, "ok")
}
func (RuntimeController) logs(c *gin.Context, m *hbtp.Map) {
	jobId := m.GetString("jobId")
	offset, _ := m.GetInt("offset")
	limit, _ := m.GetInt("limit")
	if jobId == "" {
		c.String(500, "param err")
		return
	}
	job := &model.TJob{}
	ok, _ := comm.Db.Where("id=?", jobId).Get(job)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	dir := filepath.Join(comm.WorkPath, common.PathBuild, job.BuildId, common.PathJobs)
	logpth := filepath.Join(dir, fmt.Sprintf("%v.log", job.Id))
	fl, err := os.Open(logpth)
	if err != nil {
		c.String(404, "Not Found File")
		return
	}
	defer fl.Close()
	off := offset
	if offset > 0 {
		off, err = fl.Seek(offset, 0)
		if err != nil {
			c.String(510, "err:%v", err)
			return
		}
	}
	ls := make([]*bean.LogOutJsonRes, 0)
	bts := make([]byte, 1024*5)
	linebuf := &bytes.Buffer{}
	for !hbtp.EndContext(c) {
		rn, err := fl.Read(bts)
		if rn > 0 {
			for i := 0; i < rn; i++ {
				off++
				b := bts[i]
				if linebuf == nil && b == '{' {
					linebuf.Reset()
				}
				if linebuf != nil {
					if b == '\n' {
						e := &bean.LogOutJsonRes{}
						err := json.Unmarshal(linebuf.Bytes(), e)
						linebuf.Reset()
						if err == nil {
							/*if e.Type == hbtpBean.TypeCmdLogLineSys {
								continue
							}*/
							e.Offset = off - 1
							ls = append(ls, e)
						}
						if limit > 0 && limit >= int64(len(ls)) {
							break
						}
					} else {
						linebuf.WriteByte(b)
					}
				}
			}
		}
		if err != nil {
			break
		}
	}
}
