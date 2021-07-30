package migrates

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gokins/gokins/comm"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"path/filepath"
	"strings"
)

func InitMysqlMigrate(host, dbs, user, pass string) (wait bool, rtul string, errs error) {
	wait = false
	if host == "" || dbs == "" || user == "" {
		errs = errors.New("database config not found")
		return
	}
	wait = true
	ul := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&multiStatements=true",
		user,
		pass,
		host,
		dbs)
	db, err := sql.Open("mysql", ul)
	if err != nil {
		errs = err
		return
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		uls := fmt.Sprintf("%s:%s@tcp(%s)/?parseTime=true&multiStatements=true",
			user,
			pass,
			host)
		db, err = sql.Open("mysql", uls)
		if err != nil {
			println("open dbs err:" + err.Error())
			errs = err
			return
		}
		defer db.Close()
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8mb4;", dbs))
		if err != nil {
			println("create dbs err:" + err.Error())
			errs = err
			return
		}
		db.Exec(fmt.Sprintf("USE `%s`;", dbs))
		err = db.Ping()
	}
	defer db.Close()
	wait = false
	if err != nil {
		errs = err
		return
	}

	// Run migrations
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		println("could not start sql migration... ", err.Error())
		errs = err
		return
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
		errs = err
		return
	}
	defer sc.Close()
	mgt, err := migrate.NewWithInstance(
		"bindata", sc,
		"mysql", driver)
	if err != nil {
		errs = err
		return
	}
	defer mgt.Close()
	err = mgt.Up()
	if err != nil && err != migrate.ErrNoChange {
		mgt.Down()
		errs = err
		return
	}

	return false, ul, nil
}

func InitSqliteMigrate() (rtul string, errs error) {
	ul := filepath.Join(comm.WorkPath, "db.dat")
	db, err := sql.Open("sqlite3", ul)
	if err != nil {
		errs = err
		return
	}
	defer db.Close()

	// Run migrations
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		println("could not start sql migration... ", err.Error())
		errs = err
		return
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
		errs = err
		return
	}
	defer sc.Close()
	mgt, err := migrate.NewWithInstance(
		"bindata", sc,
		"sqlite3", driver)
	if err != nil {
		errs = err
		return
	}
	defer mgt.Close()
	err = mgt.Up()
	if err != nil && err != migrate.ErrNoChange {
		mgt.Down()
		errs = err
		return
	}

	return ul, nil
}
