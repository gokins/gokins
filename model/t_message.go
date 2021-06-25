package model

import (
	"time"
)

type TMessage struct {
	Id      string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid     int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid     string    `xorm:"comment('发送者（可空）') VARCHAR(64)" json:"uid"`
	Title   string    `xorm:"VARCHAR(255)" json:"title"`
	Content string    `xorm:"LONGTEXT" json:"content"`
	Types   string    `xorm:"VARCHAR(50)" json:"types"`
	Created time.Time `xorm:"DATETIME" json:"created"`
	Infos   string    `xorm:"TEXT" json:"infos"`
	Url     string    `xorm:"VARCHAR(500)" json:"url"`
}
