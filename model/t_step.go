package model

import (
	"time"
)

type TStep struct {
	Id                string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	BuildId           string    `xorm:"VARCHAR(64)" json:"buildId"`
	StageId           string    `xorm:"comment('流水线id') VARCHAR(100)" json:"stageId"`
	DisplayName       string    `xorm:"VARCHAR(255)" json:"displayName"`
	PipelineVersionId string    `xorm:"comment('流水线id') VARCHAR(64)" json:"pipelineVersionId"`
	Step              string    `xorm:"VARCHAR(255)" json:"step"`
	Status            string    `xorm:"comment('构建状态') VARCHAR(100)" json:"status"`
	Event             string    `xorm:"comment('事件') VARCHAR(100)" json:"event"`
	ExitCode          int       `xorm:"comment('退出码') INT(11)" json:"exitCode"`
	Error             string    `xorm:"comment('错误信息') VARCHAR(500)" json:"error"`
	Name              string    `xorm:"comment('名字') VARCHAR(100)" json:"name"`
	Started           time.Time `xorm:"comment('开始时间') DATETIME" json:"started"`
	Finished          time.Time `xorm:"comment('结束时间') DATETIME" json:"finished"`
	Created           time.Time `xorm:"comment('创建时间') DATETIME" json:"created"`
	Updated           time.Time `xorm:"comment('更新时间') DATETIME" json:"updated"`
	Version           string    `xorm:"comment('版本') VARCHAR(255)" json:"version"`
	Errignore         int       `xorm:"INT(11)" json:"errignore"`
	Commands          string    `xorm:"TEXT" json:"commands"`
	Waits             string    `xorm:"JSON" json:"waits"`
	Sort              int       `xorm:"INT(11)" json:"sort"`
}
