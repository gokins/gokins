package model

type TVariable struct {
	Id    string `xorm:"not null pk VARCHAR(64)" json:"id"`
	Name  string `xorm:"VARCHAR(255)" json:"name"`
	Value string `xorm:"VARCHAR(255)" json:"value"`
}
