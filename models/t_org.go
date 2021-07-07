package models

import (
	"time"
)

type TOrg struct {
	Id          string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid         int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid         string    `xorm:"VARCHAR(64)" json:"uid"`
	Name        string    `xorm:"VARCHAR(100)" json:"name"`
	Desc        string    `xorm:"TEXT" json:"desc"`
	Public      int       `xorm:"default 0 comment('公开') INT(1)" json:"public"`
	Created     time.Time `xorm:"comment('创建时间') DATETIME" json:"created"`
	Updated     time.Time `xorm:"comment('更新时间') DATETIME" json:"updated"`
	Deleted     int       `xorm:"default 0 INT(1)" json:"-"`
	DeletedTime time.Time `xorm:"DATETIME" json:"-"`
}

type TOrgInfo struct {
	Id          string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid         int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid         string    `xorm:"VARCHAR(64)" json:"uid"`
	Name        string    `xorm:"VARCHAR(100)" json:"name"`
	Desc        string    `xorm:"TEXT" json:"desc"`
	Public      int       `xorm:"default 0 comment('公开') INT(1)" json:"public"`
	Created     time.Time `xorm:"comment('创建时间') DATETIME" json:"created"`
	Updated     time.Time `xorm:"comment('更新时间') DATETIME" json:"updated"`
	Deleted     int       `xorm:"default 0 INT(1)" json:"-"`
	DeletedTime time.Time `xorm:"DATETIME" json:"-"`

	Nick   string `xorm:"-" json:"nick"`
	Avat   string `xorm:"-" json:"avat"`
	Pipeln int64  `xorm:"-" json:"pipeln"`
	Userln int64  `xorm:"-" json:"userln"`
}

func (TOrgInfo) TableName() string {
	return "t_org"
}
