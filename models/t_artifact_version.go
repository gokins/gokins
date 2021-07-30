package models

import (
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/gokins/core/common"
	"github.com/gokins/gokins/comm"
)

type TArtifactVersion struct {
	Id        string    `xorm:"not null pk VARCHAR(64)" json:"id"`
	Aid       int64     `xorm:"not null pk autoincr BIGINT(20)" json:"aid"`
	RepoId    string    `xorm:"VARCHAR(64)" json:"repoId"`
	PackageId string    `xorm:"VARCHAR(64)" json:"packageId"`
	Name      string    `xorm:"VARCHAR(100)" json:"name"`
	Version   string    `xorm:"VARCHAR(100)" json:"version"`
	Sha       string    `xorm:"VARCHAR(100)" json:"sha"`
	Desc      string    `xorm:"VARCHAR(500)" json:"desc"`
	Preview   int       `xorm:"INT(1)" json:"preview"`
	Created   time.Time `xorm:"DATETIME" json:"created"`
	Updated   time.Time `xorm:"DATETIME" json:"updated"`

	Files []hbtp.Map `xorm:"-" json:"files"`
}

func (c *TArtifactVersion) ReadFiles() error {
	dir := filepath.Join(comm.WorkPath, common.PathArtifacts, c.Id)
	fls, err := c.readDir(dir)
	c.Files = fls
	return err
}
func (c *TArtifactVersion) readDir(pth string) ([]hbtp.Map, error) {
	var rts []hbtp.Map
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
			var chd []hbtp.Map
			chd = append(chd, fls...)
			rts = append(rts, hbtp.Map{
				"name":  v.Name(),
				"dir":   true,
				"size":  0,
				"child": chd,
			})
		} else {
			rts = append(rts, hbtp.Map{
				"name": v.Name(),
				"dir":  false,
				"size": v.Size(),
			})
		}
	}
	return rts, nil
}
