package model

type TPipeline struct {
	Id           string `xorm:"not null pk VARCHAR(64)" json:"id"`
	Name         string `xorm:"VARCHAR(255)" json:"name"`
	DisplayName  string `xorm:"VARCHAR(255)" json:"displayName"`
	PipelineType string `xorm:"VARCHAR(255)" json:"pipelineType"`
	JsonContent  string `xorm:"LONGTEXT" json:"jsonContent"`
	CreateUserId string `xorm:"VARCHAR(255)" json:"createUserId"`
}
