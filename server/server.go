package server

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gokins-main/core"
	utils2 "github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/engine"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"xorm.io/xorm"
)

func Run() error {
	if comm.WorkPath == "" {
		pth := filepath.Join(utils2.HomePath(), ".gokins")
		comm.WorkPath = utils2.EnvDefault("GOKINS_WORKPATH", pth)
	}
	os.MkdirAll(comm.WorkPath, 0750)
	core.InitLog(comm.WorkPath)
	err := parseConfig()
	if err != nil {
		//return err
		comm.Cfg.Server.Secret = "123456"
	}
	/*err = initDb()
	if err != nil {
		return err
	}*/
	err = engine.Start()
	if err != nil {
		return err
	}
	//go runWeb()
	//runHbtp()
	runWeb()
	hbtp.Infof("gokins running in %s", comm.WorkPath)
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
func initConfig() error {
	bts, err := yaml.Marshal(&comm.Cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(comm.WorkPath, "app.yml"), bts, 0644)
}

func initDb() error {
	dvs := "mysql"
	if comm.Cfg.Database.Driver != "" {
		dvs = comm.Cfg.Database.Driver
	}
	db, err := xorm.NewEngine(dvs, comm.Cfg.Database.Url)
	if err != nil {
		return err
	}
	comm.Db = db
	return nil
}
