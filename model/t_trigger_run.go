package model

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
}
