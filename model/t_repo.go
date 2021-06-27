package model

type TRepo struct {
	Id   string `xorm:"not null pk VARCHAR(64)" json:"id"`
	Name string `xorm:"VARCHAR(255)" json:"name"`
	Url  string `xorm:"TEXT" json:"url"`
}
