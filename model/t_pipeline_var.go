package model

type TPipelineVar struct {
	Aid        int64  `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	Uid        string `xorm:"VARCHAR(64)" json:"uid"`
	PipelineId string `xorm:"VARCHAR(64)" json:"pipelineId"`
	Name       string `xorm:"VARCHAR(255)" json:"name"`
	Value      string `xorm:"VARCHAR(255)" json:"value"`
	Remarks    string `xorm:"VARCHAR(255)" json:"remarks"`
	Public     int    `xorm:"default 0 comment('公开') INT(1)" json:"public"`
}
