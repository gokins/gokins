package model

import (
	"time"
)

type TPipelineVersion struct {
	Id                  string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Uid                 string    `xorm:"VARCHAR(64)" json:"uid"`
	Number              int64     `xorm:"comment('构建次数') BIGINT(20)" json:"number"`
	Events              string    `xorm:"comment('事件push、pr、note') VARCHAR(100)" json:"events"`
	Sha                 string    `xorm:"VARCHAR(255)" json:"sha"`
	PipelineName        string    `xorm:"VARCHAR(255)" json:"pipelineName"`
	PipelineDisplayName string    `xorm:"VARCHAR(255)" json:"pipelineDisplayName"`
	PipelineId          string    `xorm:"VARCHAR(64)" json:"pipelineId"`
	Version             string    `xorm:"VARCHAR(255)" json:"version"`
	Content             string    `xorm:"LONGTEXT" json:"content"`
	Created             time.Time `xorm:"DATETIME" json:"created"`
	Deleted             int       `xorm:"default 0 TINYINT(1)" json:"deleted"`
	PrNumber            int64     `xorm:"BIGINT(20)" json:"prNumber"`
	RepoCloneUrl        string    `xorm:"VARCHAR(255)" json:"repoCloneUrl"`
}
