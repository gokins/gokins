package model

import (
	"time"
)

type TYmlPlugin struct {
	Aid         int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Name        string    `xorm:"VARCHAR(64)" json:"name"`
	YmlContent  string    `xorm:"LONGTEXT" json:"ymlContent"`
	Deleted     int       `xorm:"default 0 INT(1)" json:"deleted"`
	DeletedTime time.Time `xorm:"DATETIME" json:"deletedTime"`
}
