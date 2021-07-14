package models

import (
	"time"
)

type TArtifactPackage struct {
	Id          string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid         int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	RepoId      string    `xorm:"index VARCHAR(64)" json:"repoId"`
	Name        string    `xorm:"VARCHAR(100)" json:"name"`
	DisplayName string    `xorm:"VARCHAR(255)" json:"displayName"`
	Desc        string    `xorm:"VARCHAR(500)" json:"desc"`
	Created     time.Time `xorm:"DATETIME" json:"created"`
	Updated     time.Time `xorm:"DATETIME" json:"updated"`
	Deleted     int       `xorm:"INT(1)" json:"deleted"`
	DeletedTime time.Time `xorm:"DATETIME" json:"deletedTime"`

	Verln int64 `xorm:"-" json:"verln"`
}
