package model

import (
	"time"
)

type TOrg struct {
	Id          string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid         int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid         string    `xorm:"index VARCHAR(64)" json:"uid"`
	Name        string    `xorm:"VARCHAR(200)" json:"name"`
	Desc        string    `xorm:"TEXT" json:"desc"`
	Public      int       `xorm:"default 0 comment('公开') INT(1)" json:"public"`
	Created     time.Time `xorm:"comment('创建时间') DATETIME" json:"created"`
	Updated     time.Time `xorm:"comment('更新时间') DATETIME" json:"updated"`
	Deleted     int       `xorm:"default 0 INT(1)" json:"deleted"`
	DeletedTime time.Time `xorm:"DATETIME" json:"deletedTime"`
}
