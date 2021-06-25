package model

import (
	"time"
)

type TUserMsg struct {
	Id          int64     `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Mid         string    `xorm:"VARCHAR(64)" json:"mid"`
	Uid         string    `xorm:"comment('收件人') VARCHAR(64)" json:"uid"`
	Created     time.Time `xorm:"DATETIME" json:"created"`
	Readtm      time.Time `xorm:"DATETIME" json:"readtm"`
	Status      int       `xorm:"default 0 INT(11)" json:"status"`
	Deleted     int       `xorm:"default 0 INT(1)" json:"deleted"`
	DeletedTime time.Time `xorm:"DATETIME" json:"deletedTime"`
}
