package model

import (
	"time"
)

type TPipelineVersion struct {
	Id                  string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Number              int64     `xorm:"comment('构建次数') BIGINT(20)" json:"number"`
	Trigger             string    `xorm:"comment('触发方式') VARCHAR(100)" json:"trigger"`
	Events              string    `xorm:"comment('事件push、pr、note') VARCHAR(100)" json:"events"`
	Ref                 string    `xorm:"VARCHAR(255)" json:"ref"`
	Branch              string    `xorm:"VARCHAR(255)" json:"branch"`
	RepoId              string    `xorm:"VARCHAR(64)" json:"repoId"`
	RepoName            string    `xorm:"VARCHAR(255)" json:"repoName"`
	CommitSha           string    `xorm:"VARCHAR(255)" json:"commitSha"`
	CommitMessage       string    `xorm:"comment('提交信息') TEXT" json:"commitMessage"`
	PipelineName        string    `xorm:"VARCHAR(255)" json:"pipelineName"`
	PipelineDisplayName string    `xorm:"VARCHAR(255)" json:"pipelineDisplayName"`
	PipelineId          string    `xorm:"VARCHAR(64)" json:"pipelineId"`
	Version             string    `xorm:"VARCHAR(255)" json:"version"`
	YmlContent          string    `xorm:"LONGTEXT" json:"ymlContent"`
	Created             time.Time `xorm:"DATETIME" json:"created"`
	CreateUser          string    `xorm:"VARCHAR(255)" json:"createUser"`
	CreateUserId        string    `xorm:"VARCHAR(64)" json:"createUserId"`
	Deleted             int       `xorm:"default 0 TINYINT(1)" json:"deleted"`
	TargetRepoName      string    `xorm:"VARCHAR(255)" json:"targetRepoName"`
	TargetRepoSha       string    `xorm:"VARCHAR(255)" json:"targetRepoSha"`
	TargetRepoRef       string    `xorm:"VARCHAR(255)" json:"targetRepoRef"`
	TargetRepoCloneUrl  string    `xorm:"VARCHAR(255)" json:"targetRepoCloneUrl"`
	Status              string    `xorm:"comment('构建状态') VARCHAR(100)" json:"status"`
	Error               string    `xorm:"comment('错误信息') VARCHAR(500)" json:"error"`
	Note                string    `xorm:"VARCHAR(255)" json:"note"`
	Title               string    `xorm:"VARCHAR(255)" json:"title"`
	PrNumber            int64     `xorm:"BIGINT(20)" json:"prNumber"`
	RepoCloneUrl        string    `xorm:"VARCHAR(255)" json:"repoCloneUrl"`
}
