package model

import (
	"time"
)

type TRole struct {
	Id    int64     `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Xid   string    `xorm:"not null pk VARCHAR(64)" json:"xid"`
	Title string    `xorm:"VARCHAR(100)" json:"title"`
	Perms string    `xorm:"TEXT" json:"perms"`
	Times time.Time `xorm:"default CURRENT_TIMESTAMP DATETIME" json:"times"`
}
