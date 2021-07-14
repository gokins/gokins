package model

import (
	"time"
)

type TArtifactVersion struct {
	Id        string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid       int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	RepoId    string    `xorm:"index(rpnm) VARCHAR(64)" json:"repoId"`
	PackageId string    `xorm:"VARCHAR(64)" json:"packageId"`
	Name      string    `xorm:"index(rpnm) VARCHAR(100)" json:"name"`
	Version   string    `xorm:"VARCHAR(100)" json:"version"`
	Sha       string    `xorm:"VARCHAR(100)" json:"sha"`
	Desc      string    `xorm:"VARCHAR(500)" json:"desc"`
	Preview   int       `xorm:"INT(1)" json:"preview"`
	Created   time.Time `xorm:"DATETIME" json:"created"`
	Updated   time.Time `xorm:"DATETIME" json:"updated"`
}
