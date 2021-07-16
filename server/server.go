package server

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gokins-main/core"
	utils2 "github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/engine"
	"github.com/gokins-main/gokins/route"
	"github.com/gokins-main/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func Run() error {
	if comm.WorkPath == "" {
		pth := filepath.Join(utils2.HomePath(), ".gokins")
		comm.WorkPath = utils2.EnvDefault("GOKINS_WORKPATH", pth)
	}
	if !comm.NotUpPass {
		comm.NotUpPass = utils2.EnvDefault("GOKINS_NOTUPDATEPASS") == "true"
	}

	os.MkdirAll(comm.WorkPath, 0750)
	core.InitLog(comm.WorkPath)
	go runWeb()
	time.Sleep(time.Millisecond * 10)
	err := parseConfig()
	if err != nil {
		logrus.Debugf("parseConfig err:%v", err)
		comm.WebEgn.GET("/install", route.Install)
		util.GinRegController(comm.WebEgn, &route.InstallController{})
		for !comm.Installed {
			time.Sleep(time.Millisecond * 100)
			if hbtp.EndContext(comm.Ctx) {
				return errors.New("ctx dead")
			}
		}
	}

	err = initDb()
	if err != nil {
		return err
	}
	err = initCache()
	if err != nil {
		return err
	}
	defer comm.BCache.Close()

	regApi()
	comm.Installed = true
	err = engine.Start()
	if err != nil {
		return err
	}

	go runHbtp()
	hbtp.Infof("gokins running in %s", comm.WorkPath)
	for !hbtp.EndContext(comm.Ctx) {
		time.Sleep(time.Millisecond * 100)
	}
	time.Sleep(time.Second)
	return nil
}
func parseConfig() error {
	bts, err := ioutil.ReadFile(filepath.Join(comm.WorkPath, "app.yml"))
	if err != nil {
		bts, err = ioutil.ReadFile(filepath.Join(comm.WorkPath, "app.yaml"))
	}
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bts, &comm.Cfg)
}
