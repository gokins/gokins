package models

import (
	"time"
)

type TPipeline struct {
	Id           string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Uid          string    `xorm:"VARCHAR(64)" json:"uid"`
	Name         string    `xorm:"VARCHAR(255)" json:"name"`
	DisplayName  string    `xorm:"VARCHAR(255)" json:"displayName"`
	PipelineType string    `xorm:"VARCHAR(255)" json:"pipelineType"`
	AccessToken  string    `xorm:"-" json:"accessToken"`
	Url          string    `xorm:"-" json:"url"`
	Username     string    `xorm:"-" json:"username"`
	Deleted      int       `xorm:"default 0 INT(1)" json:"-"`
	DeletedTime  time.Time `xorm:"DATETIME" json:"-"`
	CreateTime   time.Time `xorm:"DATETIME" json:"-"`

	Nick    string    `xorm:"-" json:"nick"`
	Avat    string    `xorm:"-" json:"avat"`
	Buildln int64     `xorm:"-" json:"buildln"`
	Build   *RunBuild `xorm:"-" json:"build"`
}

type TPipelineInfo struct {
	Id           string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Uid          string    `xorm:"VARCHAR(64)" json:"uid"`
	Name         string    `xorm:"VARCHAR(255)" json:"name"`
	DisplayName  string    `xorm:"VARCHAR(255)" json:"displayName"`
	PipelineType string    `xorm:"VARCHAR(255)" json:"pipelineType"`
	YmlContent   string    `xorm:"-" json:"ymlContent"`
	AccessToken  string    `xorm:"-" json:"accessToken"`
	Url          string    `xorm:"-" json:"url"`
	Username     string    `xorm:"-" json:"username"`
	Created      time.Time `xorm:"DATETIME" json:"created"`
}

func (TPipelineInfo) TableName() string {
	return "t_pipeline"
}
