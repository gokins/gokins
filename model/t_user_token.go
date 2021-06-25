package model

import (
	"time"
)

type TUserToken struct {
	Id           int64     `xorm:"pk autoincr BIGINT(20)" json:"id"`
	Uid          int64     `xorm:"BIGINT(20)" json:"uid"`
	Type         string    `xorm:"VARCHAR(50)" json:"type"`
	Openid       string    `xorm:"VARCHAR(100)" json:"openid"`
	Name         string    `xorm:"VARCHAR(255)" json:"name"`
	Nick         string    `xorm:"VARCHAR(255)" json:"nick"`
	Avatar       string    `xorm:"VARCHAR(500)" json:"avatar"`
	AccessToken  string    `xorm:"TEXT" json:"accessToken"`
	RefreshToken string    `xorm:"TEXT" json:"refreshToken"`
	ExpiresIn    int64     `xorm:"default 0 BIGINT(20)" json:"expiresIn"`
	ExpiresTime  time.Time `xorm:"DATETIME" json:"expiresTime"`
	RefreshTime  time.Time `xorm:"DATETIME" json:"refreshTime"`
	CreateTime   time.Time `xorm:"DATETIME" json:"createTime"`
	Tokens       string    `xorm:"TEXT" json:"tokens"`
	Uinfos       string    `xorm:"TEXT" json:"uinfos"`
}
