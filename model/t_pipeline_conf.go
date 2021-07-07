package model

type TPipelineConf struct {
	Aid         int    `xorm:"not null pk autoincr INT(20)" json:"aid"`
	PipelineId  string `xorm:"not null VARCHAR(64)" json:"pipelineId"`
	Url         string `xorm:"VARCHAR(255)" json:"url"`
	AccessToken string `xorm:"VARCHAR(255)" json:"accessToken"`
	YmlContent  string `xorm:"LONGTEXT" json:"ymlContent"`
	Username    string `xorm:"VARCHAR(255)" json:"username"`
}
