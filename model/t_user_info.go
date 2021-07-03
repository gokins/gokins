package model

import (
	"time"
)

type TUserInfo struct {
	Id       string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Phone    string    `xorm:"VARCHAR(100)" json:"phone"`
	Email    string    `xorm:"VARCHAR(200)" json:"email"`
	Remark   string    `xorm:"TEXT" json:"remark"`
	Birthday time.Time `xorm:"DATETIME" json:"birthday"`
	PermUser int       `xorm:"INT(1)" json:"permUser"`
	PermOrg  int       `xorm:"INT(1)" json:"permOrg"`
	PermPipe int       `xorm:"INT(1)" json:"permPipe"`
}
