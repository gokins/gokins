package model

import (
	"time"
)

type TBuild struct {
	Id                string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	PipelineId        string    `xorm:"VARCHAR(64)" json:"pipelineId"`
	PipelineVersionId string    `xorm:"VARCHAR(64)" json:"pipelineVersionId"`
	Status            string    `xorm:"comment('构建状态') VARCHAR(100)" json:"status"`
	Error             string    `xorm:"comment('错误信息') VARCHAR(500)" json:"error"`
	Event             string    `xorm:"comment('事件') VARCHAR(100)" json:"event"`
	Started           time.Time `xorm:"comment('开始时间') DATETIME" json:"started"`
	Finished          time.Time `xorm:"comment('结束时间') DATETIME" json:"finished"`
	Created           time.Time `xorm:"comment('创建时间') DATETIME" json:"created"`
	Updated           time.Time `xorm:"comment('更新时间') DATETIME" json:"updated"`
	Version           string    `xorm:"comment('版本') VARCHAR(255)" json:"version"`
}
