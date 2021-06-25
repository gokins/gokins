package model

import (
	"time"
)

type TParam struct {
	Id    int64     `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Name  string    `xorm:"VARCHAR(100)" json:"name"`
	Title string    `xorm:"VARCHAR(255)" json:"title"`
	Data  string    `xorm:"TEXT" json:"data"`
	Times time.Time `xorm:"DATETIME" json:"times"`
}
