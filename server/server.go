package server

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gokins-main/core"
	utils2 "github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/engine"
	"github.com/gokins-main/gokins/migrates"
	"github.com/gokins-main/gokins/route"
	"github.com/gokins-main/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"xorm.io/xorm"
)

func Run() error {
	if comm.WorkPath == "" {
		pth := filepath.Join(utils2.HomePath(), ".gokins")
		comm.WorkPath = utils2.EnvDefault("GOKINS_WORKPATH", pth)
	}

	os.MkdirAll(comm.WorkPath, 0750)
	core.InitLog(comm.WorkPath)
	go runWeb()

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

	comm.Installed = true
	regApi()

	err = engine.Start()
	if err != nil {
		return err
	}

	go runHbtp()
	hbtp.Infof("gokins running in %s", comm.WorkPath)
	for !hbtp.EndContext(comm.Ctx) {
		time.Sleep(time.Millisecond * 100)
	}
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

func initDb() error {
	var err error
	dvs := "mysql"
	ul := comm.Cfg.Datasource.Url
	if comm.Cfg.Datasource.Driver != "" {
		dvs = comm.Cfg.Datasource.Driver
	}
	if dvs == "mysql" {
		err = migrates.UpMysqlMigrate(ul)
	}
	if err != nil {
		return err
	}
	db, err := xorm.NewEngine(dvs, comm.Cfg.Datasource.Url)
	if err != nil {
		return err
	}
	comm.Db = db
	return nil
}
