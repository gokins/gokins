package model

import (
	"time"
)

type TJob struct {
	Id                string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	BuildId           string    `xorm:"VARCHAR(64)" json:"buildId"`
	StageId           string    `xorm:"comment('流水线id') VARCHAR(100)" json:"stageId"`
	DisplayName       string    `xorm:"VARCHAR(255)" json:"displayName"`
	PipelineVersionId string    `xorm:"comment('流水线id') VARCHAR(64)" json:"pipelineVersionId"`
	Job               string    `xorm:"VARCHAR(255)" json:"job"`
	Status            string    `xorm:"comment('构建状态') VARCHAR(100)" json:"status"`
	ExitCode          int64     `xorm:"comment('退出码') BIGINT(20)" json:"exitCode"`
	Error             string    `xorm:"comment('错误信息') VARCHAR(500)" json:"error"`
	Name              string    `xorm:"comment('名字') VARCHAR(100)" json:"name"`
	Started           time.Time `xorm:"comment('开始时间') DATETIME" json:"started"`
	Finished          time.Time `xorm:"comment('结束时间') DATETIME" json:"finished"`
	Created           time.Time `xorm:"comment('创建时间') DATETIME" json:"created"`
	Updated           time.Time `xorm:"comment('更新时间') DATETIME" json:"updated"`
	Version           string    `xorm:"comment('版本') VARCHAR(255)" json:"version"`
	Errignore         string    `xorm:"VARCHAR(5)" json:"errignore"`
	Number            int64     `xorm:"BIGINT(20)" json:"number"`
	Commands          string    `xorm:"TEXT" json:"commands"`
	DependsOn         string    `xorm:"JSON" json:"dependsOn"`
	Image             string    `xorm:"VARCHAR(255)" json:"image"`
	Environments      string    `xorm:"JSON" json:"environments"`
	Sort              int64     `xorm:"BIGINT(10)" json:"sort"`
}
