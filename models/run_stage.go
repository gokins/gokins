package models

import (
	"time"
)

type RunStage struct {
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
	Stage             string    `xorm:"VARCHAR(255)" json:"stage"`
	Stepids           []string  `xorm:"-" json:"stepids"`
}

func (RunStage) TableName() string {
	return "t_stage"
}

type RunStep struct {
	Id                string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	BuildId           string    `xorm:"VARCHAR(64)" json:"buildId"`
	StageId           string    `xorm:"comment('流水线id') VARCHAR(100)" json:"stageId"`
	DisplayName       string    `xorm:"VARCHAR(255)" json:"displayName"`
	PipelineVersionId string    `xorm:"comment('流水线id') VARCHAR(64)" json:"pipelineVersionId"`
	Step              string    `xorm:"VARCHAR(255)" json:"step"`
	Status            string    `xorm:"comment('构建状态') VARCHAR(100)" json:"status"`
	ExitCode          int64     `xorm:"comment('退出码') BIGINT(20)" json:"exitCode"`
	Error             string    `xorm:"comment('错误信息') VARCHAR(500)" json:"error"`
	Name              string    `xorm:"comment('名字') VARCHAR(100)" json:"name"`
	Started           time.Time `xorm:"comment('开始时间') DATETIME" json:"started"`
	Finished          time.Time `xorm:"comment('结束时间') DATETIME" json:"finished"`
	Created           time.Time `xorm:"comment('创建时间') DATETIME" json:"created"`
	Updated           time.Time `xorm:"comment('更新时间') DATETIME" json:"updated"`
	Errignore         int       `xorm:"INT(11)" json:"errignore"`
	Waits             string    `xorm:"JSON" json:"-"`
	Waitings          []string  `xorm:"-" json:"waits"`
	Sort              int64     `xorm:"BIGINT(10)" json:"sort"`
}

func (RunStep) TableName() string {
	return "t_step"
}
