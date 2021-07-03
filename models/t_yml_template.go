package models

type TYmlTemplate struct {
	Aid        int64  `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Name       string `xorm:"VARCHAR(64)" json:"name"`
	YmlContent string `xorm:"LONGTEXT" json:"ymlContent"`
}
