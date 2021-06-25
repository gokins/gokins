package model

type SchemaMigrations struct {
	Version int64 `xorm:"not null pk BIGINT(20)" json:"version"`
	Dirty   int   `xorm:"not null TINYINT(1)" json:"dirty"`
}
