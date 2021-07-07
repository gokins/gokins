package models

import (
	"time"
)

type TArtifactory struct {
	Id          string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid         int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid         string    `xorm:"VARCHAR(64)" json:"uid"`
	OrgId       string    `xorm:"VARCHAR(64)" json:"orgId"`
	Identifier  string    `xorm:"VARCHAR(50)" json:"identifier"`
	Name        string    `xorm:"VARCHAR(200)" json:"name"`
	Disabled    int       `xorm:"default 0 comment('是否归档(1归档|0正常)') INT(1)" json:"disabled"`
	Source      string    `xorm:"VARCHAR(50)" json:"source"`
	Desc        string    `xorm:"VARCHAR(500)" json:"desc"`
	Logo        string    `xorm:"VARCHAR(255)" json:"logo"`
	Created     time.Time `xorm:"DATETIME" json:"created"`
	Updated     time.Time `xorm:"DATETIME" json:"updated"`
	Deleted     int       `xorm:"INT(1)" json:"deleted"`
	DeletedTime time.Time `xorm:"DATETIME" json:"deletedTime"`

	Nick  string `xorm:"-" json:"nick"`
	Avat  string `xorm:"-" json:"avat"`
	Artln int64  `xorm:"-" json:"artln"`
}
