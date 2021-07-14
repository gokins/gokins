package models

import (
	"time"
)

type TTriggerRun struct {
	Id            string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid           int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Tid           string    `xorm:"comment('触发器ID') index VARCHAR(64)" json:"tid"`
	PipeVersionId string    `xorm:"VARCHAR(64)" json:"pipeVersionId"`
	Infos         string    `xorm:"JSON" json:"infos"`
	Error         string    `xorm:"VARCHAR(255)" json:"error"`
	Created       time.Time `xorm:"DATETIME" json:"created"`
	//version
	Number              int64  `xorm:"-" json:"number"`
	PipelineName        string `xorm:"-" json:"pipelineName"`
	PipelineDisplayName string `xorm:"-" json:"pipelineDisplayName"`
	//Build
	BStatus string `xorm:"-" json:"bStatus"`
}

type TimerTriggerRun struct {
	Id         string `xorm:"not null pk VARCHAR(64)" json:"id"`
	Uid        string `xorm:"index VARCHAR(64)" json:"uid"`
	PipelineId string `xorm:"VARCHAR(64)" json:"pipelineId"`
	Types      string `xorm:"VARCHAR(50)" json:"types"`
	Name       string `xorm:"VARCHAR(100)" json:"name"`
	Params     string `xorm:"JSON" json:"params"`
	Enabled    int    `xorm:"INT(1)" json:"enabled"`
	//run
	RunCreated time.Time `xorm:"DATETIME" json:"created"`
}
