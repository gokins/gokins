package model

type TUserRepo struct {
	Id     string `xorm:"not null pk VARCHAR(64)" json:"id"`
	RepoId string `xorm:"not null VARCHAR(64)" json:"repoId"`
	UserId int64  `xorm:"not null BIGINT(20)" json:"userId"`
}
