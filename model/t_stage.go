package model

import (
	"time"
)

type TStage struct {
	Id                string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	PipelineVersionId string    `xorm:"comment('流水线id') VARCHAR(64)" json:"pipelineVersionId"`
	BuildId           string    `xorm:"VARCHAR(64)" json:"buildId"`
	Status            string    `xorm:"comment('构建状态') VARCHAR(100)" json:"status"`
	Error             string    `xorm:"comment('错误信息') VARCHAR(500)" json:"error"`
	Name              string    `xorm:"comment('名字') VARCHAR(255)" json:"name"`
	DisplayName       string    `xorm:"VARCHAR(255)" json:"displayName"`
	Started           time.Time `xorm:"comment('开始时间') DATETIME" json:"started"`
	Finished          time.Time `xorm:"comment('结束时间') DATETIME" json:"finished"`
	Created           time.Time `xorm:"comment('创建时间') DATETIME" json:"created"`
	Updated           time.Time `xorm:"comment('更新时间') DATETIME" json:"updated"`
	Sort              int       `xorm:"INT(11)" json:"sort"`
	Stage             string    `xorm:"VARCHAR(255)" json:"stage"`
}
