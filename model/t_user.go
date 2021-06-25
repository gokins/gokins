package model

import (
	"time"
)

type TUser struct {
	Id         int64     `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Name       string    `xorm:"VARCHAR(100)" json:"name"`
	Pass       string    `xorm:"VARCHAR(255)" json:"pass"`
	Nick       string    `xorm:"VARCHAR(100)" json:"nick"`
	Avatar     string    `xorm:"VARCHAR(500)" json:"avatar"`
	CreateTime time.Time `xorm:"DATETIME" json:"createTime"`
	LoginTime  time.Time `xorm:"DATETIME" json:"loginTime"`
}
