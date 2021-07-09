package models

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/gokins-main/core/common"
	"github.com/gokins-main/gokins/comm"
)

type FlInfo struct {
	Name  string    `json:"name"`
	Dir   bool      `json:"dir"`
	Size  int64     `json:"size"`
	Child []*FlInfo `json:"child"`
}
type TArtifactVersion struct {
	Id          string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid         int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	RepoId      string    `xorm:"VARCHAR(64)" json:"repoId"`
	PackageId   string    `xorm:"VARCHAR(64)" json:"packageId"`
	Name        string    `xorm:"VARCHAR(100)" json:"name"`
	Version     string    `xorm:"VARCHAR(100)" json:"version"`
	Sha         string    `xorm:"VARCHAR(100)" json:"sha"`
	Desc        string    `xorm:"VARCHAR(500)" json:"desc"`
	Preview     int       `xorm:"INT(1)" json:"preview"`
	Created     time.Time `xorm:"DATETIME" json:"created"`
	Updated     time.Time `xorm:"DATETIME" json:"updated"`
	Deleted     int       `xorm:"INT(1)" json:"deleted"`
	DeletedTime time.Time `xorm:"DATETIME" json:"deletedTime"`

	Files []*FlInfo `xorm:"-" json:"files"`
}

func (c *TArtifactVersion) ReadFiles() error {
	dir := filepath.Join(comm.WorkPath, common.PathArtifacts, c.Id)
	fls, err := c.readDir(dir)
	c.Files = fls
	return err
}
func (c *TArtifactVersion) readDir(pth string) ([]*FlInfo, error) {
	var rts []*FlInfo
	fls, err := ioutil.ReadDir(pth)
	if err != nil {
		return nil, err
	}
	for _, v := range fls {
		if v.IsDir() {
			fls, err := c.readDir(filepath.Join(pth, v.Name()))
			if err != nil {
				return nil, err
			}
			rts = append(rts, &FlInfo{
				Name:  v.Name(),
				Dir:   true,
				Size:  0,
				Child: rts,
			})
			rts = append(rts, fls...)
		} else {
			rts = append(rts, &FlInfo{
				Name:  v.Name(),
				Dir:   false,
				Size:  v.Size(),
				Child: nil,
			})
		}
	}
	return rts, nil
}
