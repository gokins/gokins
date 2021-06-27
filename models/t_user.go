package models

import (
	"time"
)

type TUser struct {
	Id        string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid       int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Name      string    `xorm:"VARCHAR(100)" json:"name"`
	Nick      string    `xorm:"VARCHAR(100)" json:"nick"`
	Avatar    string    `xorm:"VARCHAR(500)" json:"avatar"`
	Created   time.Time `xorm:"DATETIME" json:"created"`
	LoginTime time.Time `xorm:"DATETIME" json:"loginTime"`
}
