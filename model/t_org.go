package model

import (
	"time"
)

type TOrg struct {
	Id          string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid         int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Name        string    `xorm:"VARCHAR(100)" json:"name"`
	Title       string    `xorm:"VARCHAR(200)" json:"title"`
	Desc        string    `xorm:"TEXT" json:"desc"`
	Public      int       `xorm:"default 0 comment('公开') INT(1)" json:"public"`
	Created     time.Time `xorm:"DATETIME" json:"created"`
	Deleted     int       `xorm:"default 0 INT(1)" json:"deleted"`
	DeletedTime time.Time `xorm:"DATETIME" json:"deletedTime"`
}
