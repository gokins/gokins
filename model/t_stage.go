package model

import (
	"time"
)

type TStage struct {
	Id                string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	PipelineVersionId string    `xorm:"comment('流水线id') VARCHAR(64)" json:"pipelineVersionId"`
	BuildId           string    `xorm:"VARCHAR(64)" json:"buildId"`
	Status            string    `xorm:"comment('构建状态') VARCHAR(100)" json:"status"`
	ExitCode          int64     `xorm:"comment('退出码') BIGINT(20)" json:"exitCode"`
	Error             string    `xorm:"comment('错误信息') VARCHAR(500)" json:"error"`
	Name              string    `xorm:"comment('名字') VARCHAR(255)" json:"name"`
	DisplayName       string    `xorm:"VARCHAR(255)" json:"displayName"`
	Started           time.Time `xorm:"comment('开始时间') DATETIME" json:"started"`
	Finished          time.Time `xorm:"comment('结束时间') DATETIME" json:"finished"`
	Created           time.Time `xorm:"comment('创建时间') DATETIME" json:"created"`
	Updated           time.Time `xorm:"comment('更新时间') DATETIME" json:"updated"`
	Version           string    `xorm:"comment('版本') VARCHAR(255)" json:"version"`
	OnSuccess         string    `xorm:"VARCHAR(5)" json:"onSuccess"`
	OnFailure         string    `xorm:"VARCHAR(5)" json:"onFailure"`
	Sort              int64     `xorm:"BIGINT(10)" json:"sort"`
	Stage             string    `xorm:"VARCHAR(255)" json:"stage"`
}
