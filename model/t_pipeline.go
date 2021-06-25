package model

type TPipeline struct {
	Id           string `xorm:"not null pk VARCHAR(64)" json:"id"`
	Name         string `xorm:"VARCHAR(255)" json:"name"`
	RepoId       string `xorm:"VARCHAR(64)" json:"repoId"`
	DisplayName  string `xorm:"VARCHAR(255)" json:"displayName"`
	PipelineType string `xorm:"VARCHAR(255)" json:"pipelineType"`
	JsonContent  string `xorm:"LONGTEXT" json:"jsonContent"`
}
