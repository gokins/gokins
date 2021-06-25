package model

import (
	"time"
)

type TUserMsg struct {
	Aid         int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid         string    `xorm:"comment('收件人') index VARCHAR(64)" json:"uid"`
	MsgId       string    `xorm:"VARCHAR(64)" json:"msgId"`
	Created     time.Time `xorm:"DATETIME" json:"created"`
	Readtm      time.Time `xorm:"DATETIME" json:"readtm"`
	Status      int       `xorm:"default 0 INT(11)" json:"status"`
	Deleted     int       `xorm:"default 0 INT(1)" json:"deleted"`
	DeletedTime time.Time `xorm:"DATETIME" json:"deletedTime"`
}
