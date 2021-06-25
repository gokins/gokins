package model

import (
	"time"
)

type TParam struct {
	Aid   int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Name  string    `xorm:"VARCHAR(100)" json:"name"`
	Title string    `xorm:"VARCHAR(255)" json:"title"`
	Data  string    `xorm:"TEXT" json:"data"`
	Times time.Time `xorm:"DATETIME" json:"times"`
}
