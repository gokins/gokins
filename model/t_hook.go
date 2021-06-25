package model

type THook struct {
	Id       string `xorm:"not null pk VARCHAR(64)" json:"id"`
	Type     string `xorm:"VARCHAR(255)" json:"type"`
	Snapshot string `xorm:"LONGTEXT" json:"snapshot"`
	Status   string `xorm:"VARCHAR(255)" json:"status"`
	Msg      string `xorm:"VARCHAR(255)" json:"msg"`
	HookType string `xorm:"VARCHAR(255)" json:"hookType"`
}
