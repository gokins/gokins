package model

import (
	"time"
)

type TPipeline struct {
	Id           string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Uid          string    `xorm:"VARCHAR(64)" json:"uid"`
	Name         string    `xorm:"VARCHAR(255)" json:"name"`
	DisplayName  string    `xorm:"VARCHAR(255)" json:"displayName"`
	PipelineType string    `xorm:"VARCHAR(255)" json:"pipelineType"`
	YmlContent   string    `xorm:"LONGTEXT" json:"ymlContent"`
	AccessToken  string    `xorm:"VARCHAR(255)" json:"accessToken"`
	Url          string    `xorm:"VARCHAR(255)" json:"url"`
	Username     string    `xorm:"VARCHAR(255)" json:"username"`
	Deleted      int       `xorm:"default 0 INT(1)" json:"deleted"`
	DeletedTime  time.Time `xorm:"DATETIME" json:"deletedTime"`
}
