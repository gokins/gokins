package model

import (
	"time"
)

type TUserToken struct {
	Aid          int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid          int64     `xorm:"index BIGINT(20)" json:"uid"`
	Type         string    `xorm:"VARCHAR(50)" json:"type"`
	Openid       string    `xorm:"index VARCHAR(100)" json:"openid"`
	Name         string    `xorm:"VARCHAR(255)" json:"name"`
	Nick         string    `xorm:"VARCHAR(255)" json:"nick"`
	Avatar       string    `xorm:"VARCHAR(500)" json:"avatar"`
	AccessToken  string    `xorm:"TEXT" json:"accessToken"`
	RefreshToken string    `xorm:"TEXT" json:"refreshToken"`
	ExpiresIn    int64     `xorm:"default 0 BIGINT(20)" json:"expiresIn"`
	ExpiresTime  time.Time `xorm:"DATETIME" json:"expiresTime"`
	RefreshTime  time.Time `xorm:"DATETIME" json:"refreshTime"`
	Created      time.Time `xorm:"DATETIME" json:"created"`
	Tokens       string    `xorm:"TEXT" json:"tokens"`
	Uinfos       string    `xorm:"TEXT" json:"uinfos"`
}
