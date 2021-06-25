package model

type TUserRole struct {
	UserId    int64  `xorm:"not null pk BIGINT(20)" json:"userId"`
	RoleCodes string `xorm:"TEXT" json:"roleCodes"`
	Limits    string `xorm:"TEXT" json:"limits"`
}
