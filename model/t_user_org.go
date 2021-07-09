package model

import (
	"time"
)

type TUserOrg struct {
	Aid      int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid      string    `xorm:"index index(uoid) VARCHAR(64)" json:"uid"`
	OrgId    string    `xorm:"index index(uoid) VARCHAR(64)" json:"orgId"`
	Created  time.Time `xorm:"DATETIME" json:"created"`
	PermAdm  int       `xorm:"default 0 comment('管理员') INT(1)" json:"permAdm"`
	PermRw   int       `xorm:"default 0 comment('编辑权限') INT(1)" json:"permRw"`
	PermExec int       `xorm:"default 0 comment('执行权限') INT(1)" json:"permExec"`
	PermDown int       `xorm:"comment('下载制品权限') INT(1)" json:"permDown"`
}
