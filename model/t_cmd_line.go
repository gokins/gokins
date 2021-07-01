package model

import (
	"time"
)

type TCmdLine struct {
	Id       string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	GroupId  string    `xorm:"VARCHAR(64)" json:"groupId"`
	BuildId  string    `xorm:"VARCHAR(64)" json:"buildId"`
	StepId   string    `xorm:"VARCHAR(64)" json:"stepId"`
	Status   string    `xorm:"VARCHAR(50)" json:"status"`
	Num      int       `xorm:"INT(11)" json:"num"`
	Code     int       `xorm:"INT(11)" json:"code"`
	Content  string    `xorm:"TEXT" json:"content"`
	Created  time.Time `xorm:"DATETIME" json:"created"`
	Started  time.Time `xorm:"DATETIME" json:"started"`
	Finished time.Time `xorm:"DATETIME" json:"finished"`
}
