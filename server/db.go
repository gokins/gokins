package server

import (
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gokins/core"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/migrates"
	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

func initDb() error {
	var err error
	dvs := comm.DATASOURCE_DRIVER_MYSQL
	ul := comm.Cfg.Datasource.Url
	if comm.Cfg.Datasource.Driver != "" {
		dvs = comm.Cfg.Datasource.Driver
	}
	comm.IsMySQL = dvs == comm.DATASOURCE_DRIVER_MYSQL
	if !comm.Installed {
		switch dvs {
		case comm.DATASOURCE_DRIVER_MYSQL:
			err = migrates.UpMysqlMigrate(ul)
			break
		case comm.DATASOURCE_DRIVER_POSTGRES:
			err = migrates.UpPostgresMigrate(ul)
			break
		default:
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
	db.ShowSQL(core.Debug)
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
