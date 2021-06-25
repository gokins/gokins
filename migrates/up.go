package migrates

import (
	"database/sql"
	"errors"
	"github.com/gokins-main/gokins/comm"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"strings"
)

func UpMysqlMigrate(ul string) error {
	if ul == "" {
		return errors.New("database config not found")
	}
	db, err := sql.Open("mysql", ul)
	if err != nil {
		//core.Log.Errorf("could not connect to postgresql database... %v", err)
		println("open db err:" + err.Error())
		return err
	}
	err = db.Ping()
	defer db.Close()
	if err != nil {
		return err
	}

	// Run migrations
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		println("could not start sql migration... ", err.Error())
		return err
	}
	defer driver.Close()
	var nms []string
	tms := comm.AssetNames()
	for _, v := range tms {
		if strings.HasPrefix(v, "mysql") {
			nms = append(nms, strings.Replace(v, "mysql/", "", 1))
		}
	}
	s := bindata.Resource(nms, func(name string) ([]byte, error) {
		return comm.Asset("mysql/" + name)
	})
	sc, err := bindata.WithInstance(s)
	if err != nil {
		return err
	}
	defer sc.Close()
	mgt, err := migrate.NewWithInstance(
		"bindata", sc,
		"mysql", driver)
	if err != nil {
		return err
	}
	defer mgt.Close()
	err = mgt.Up()
	if err != nil && err != migrate.ErrNoChange {
		mgt.Down()
		return err
	}

	return nil
}
func UpSqliteMigrate(ul string) error {
	if ul == "" {
		return errors.New("database config not found")
	}
	db, err := sql.Open("sqlite3", ul)
	if err != nil {
		//core.Log.Errorf("could not connect to postgresql database... %v", err)
		println("open db err:" + err.Error())
		return err
	}
	err = db.Ping()
	defer db.Close()
	if err != nil {
		return err
	}

	// Run migrations
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		println("could not start sql migration... ", err.Error())
		return err
	}
	defer driver.Close()
	var nms []string
	tms := comm.AssetNames()
	for _, v := range tms {
		if strings.HasPrefix(v, "sqlite") {
			nms = append(nms, strings.Replace(v, "sqlite/", "", 1))
		}
	}
	s := bindata.Resource(nms, func(name string) ([]byte, error) {
		return comm.Asset("sqlite/" + name)
	})
	sc, err := bindata.WithInstance(s)
	if err != nil {
		return err
	}
	defer sc.Close()
	mgt, err := migrate.NewWithInstance(
		"bindata", sc,
		"sqlite3", driver)
	if err != nil {
		return err
	}
	defer mgt.Close()
	err = mgt.Up()
	if err != nil && err != migrate.ErrNoChange {
		mgt.Down()
		return err
	}

	return nil
}
