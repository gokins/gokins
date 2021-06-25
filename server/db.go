package server

import (
	"github.com/boltdb/bolt"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/migrates"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"xorm.io/xorm"
)

func initDb() error {
	var err error
	dvs := "mysql"
	ul := comm.Cfg.Datasource.Url
	if comm.Cfg.Datasource.Driver != "" {
		dvs = comm.Cfg.Datasource.Driver
	}
	if !comm.Installed {
		if dvs == "mysql" {
			err = migrates.UpMysqlMigrate(ul)
		} else {
			err = migrates.UpSqliteMigrate(ul)
		}
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

func initCache() error {
	pth := filepath.Join(comm.WorkPath, "cache.dat")
	os.Remove(pth)
	db, err := bolt.Open(pth, 0640, nil)
	if err != nil {
		logrus.Errorf("InitCache err:%v", err)
		return err
	}
	comm.BCache = db
	return nil
}
