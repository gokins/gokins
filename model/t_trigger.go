package model

import (
	"time"
)

type TTrigger struct {
	Id         string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid        int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid        string    `xorm:"index VARCHAR(64)" json:"uid"`
	PipelineId string    `xorm:"VARCHAR(64)" json:"pipelineId"`
	Types      string    `xorm:"VARCHAR(50)" json:"types"`
	Name       string    `xorm:"VARCHAR(100)" json:"name"`
	Desc       string    `xorm:"VARCHAR(255)" json:"desc"`
	Params     string    `xorm:"JSON" json:"params"`
	Enabled    int       `xorm:"INT(1)" json:"enabled"`
	Created    time.Time `xorm:"DATETIME" json:"created"`
	Updated    time.Time `xorm:"DATETIME" json:"updated"`
}
