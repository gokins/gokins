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
	Active    int       `xorm:"default 0 INT(1)" json:"active"`
}

type TUserOrgInfo struct {
	Id        string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid       int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Name      string    `xorm:"VARCHAR(100)" json:"name"`
	Nick      string    `xorm:"VARCHAR(100)" json:"nick"`
	Avatar    string    `xorm:"VARCHAR(500)" json:"avatar"`
	Created   time.Time `xorm:"DATETIME" json:"created"`
	LoginTime time.Time `xorm:"DATETIME" json:"loginTime"`

	PermAdm  int       `xorm:"default 0 comment('管理员') INT(1)" json:"permAdm"`
	PermRw   int       `xorm:"default 0 comment('1只读,2读写') INT(1)" json:"permRw"`
	PermExec int       `xorm:"default 0 comment('执行权限') INT(1)" json:"permExec"`
	PermDown int       `xorm:"comment('下载制品权限') INT(1)" json:"permDown"`
	JoinTime time.Time `xorm:"join_time" json:"joinTime"`
}

func (TUserOrgInfo) TableName() string {
	return "t_user"
}
