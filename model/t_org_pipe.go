package model

import (
	"time"
)

type TOrgPipe struct {
	Aid     int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	OrgId   string    `xorm:"index VARCHAR(64)" json:"orgId"`
	PipeId  string    `xorm:"comment('收件人') VARCHAR(64)" json:"pipeId"`
	Created time.Time `xorm:"DATETIME" json:"created"`
	Public  int       `xorm:"default 0 comment('公开') INT(1)" json:"public"`
}
