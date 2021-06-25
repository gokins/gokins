package model

import (
	"time"
)

type TPermssion struct {
	Id     int64     `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Xid    string    `xorm:"not null pk VARCHAR(64)" json:"xid"`
	Parent string    `xorm:"index VARCHAR(64)" json:"parent"`
	Title  string    `xorm:"VARCHAR(100)" json:"title"`
	Value  string    `xorm:"VARCHAR(100)" json:"value"`
	Times  time.Time `xorm:"default CURRENT_TIMESTAMP DATETIME" json:"times"`
	Sort   int       `xorm:"default 10 INT(11)" json:"sort"`
}
