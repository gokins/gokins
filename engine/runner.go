package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/runner/runners"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type baseRunner struct{}

func (c *baseRunner) PullJob(plugs []string) (*runners.RunJob, error) {
	tms := time.Now()
	for time.Since(tms).Seconds() < 5 {
		v := Mgr.jobEgn.Pull(plugs)
		if v != nil {
			return v, nil
		}
	}
	return nil, errors.New("not found")
}
func (c *baseRunner) CheckCancel(buildId string) bool {
	v, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return true
	}
	return v.stopd()
}
func (c *baseRunner) Update(m *runners.UpdateJobInfo) error {
	job, ok := Mgr.jobEgn.GetJob(m.Id)
	if !ok {
		return errors.New("not found job")
	}
	tsk, ok := Mgr.buildEgn.Get(job.step.BuildId)
	if !ok {
		return errors.New("not found task")
	}
	tsk.UpJob(job, m.Status, m.Error, m.ExitCode)
	return nil
}

func (c *baseRunner) UpdateCmd(jobid, cmdid string, fs, code int) error {
	job, ok := Mgr.jobEgn.GetJob(jobid)
	if !ok {
		return errors.New("not found job")
	}
	tsk, ok := Mgr.buildEgn.Get(job.step.BuildId)
	if !ok {
		return errors.New("not found task")
	}
	job.RLock()
	cmd, ok := job.cmdmp[cmdid]
	job.RUnlock()
	if !ok {
		return errors.New("not found cmd")
	}
	tsk.UpJobCmd(cmd, fs, code)
	return nil
}
func (c *baseRunner) PushOutLine(jobid, cmdid, bs string, iserr bool) error {
	job, ok := Mgr.jobEgn.GetJob(jobid)
	if !ok {
		return errors.New("not found")
	}

	bts, err := json.Marshal(&bean.LogOutJson{
		Id:      cmdid,
		Content: bs,
		Times:   time.Now(),
		Errs:    iserr,
	})
	if err != nil {
		return err
	}

	dir := filepath.Join(comm.WorkPath, common.PathBuild, job.step.BuildId, common.PathJobs, job.step.Id)
	logpth := filepath.Join(dir, "build.log")
	os.MkdirAll(dir, 0755)
	logfl, err := os.OpenFile(logpth, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer logfl.Close()
	logfl.Write(bts)
	logfl.WriteString("\n")
	return nil
}
func (c *baseRunner) FindJobId(buildId, stgNm, stpNm string) (string, bool) {
	if buildId == "" || stgNm == "" || stpNm == "" {
		return "", false
	}
	build, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return "", false
	}
	build.staglk.RLock()
	defer build.staglk.RUnlock()
	stg, ok := build.stages[stgNm]
	if !ok {
		return "", false
	}
	for _, v := range stg.jobs {
		if v.step.Name == stpNm {
			return v.step.Id, true
		}
	}
	return "", false
}

func (c *baseRunner) ReadDir(fs int, buildId string, pth string) ([]*runners.DirEntry, error) {
	if buildId == "" || pth == "" {
		return nil, errors.New("param err")
	}
	build, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return nil, errors.New("not found build")
	}
	pths := ""
	if fs == 1 {
		pths = filepath.Join(build.repoPaths, pth)
	} else if fs == 2 {
		pths = filepath.Join(build.buildPath, common.PathJobs, pth)
	}
	fls, err := os.ReadDir(pths)
	if err != nil {
		return nil, err
	}
	var ls []*runners.DirEntry
	for _, v := range fls {
		e := &runners.DirEntry{
			Name:  v.Name(),
			IsDir: v.IsDir(),
		}
		ifo, err := v.Info()
		if err == nil {
			e.Size = ifo.Size()
		}
		ls = append(ls, e)
	}
	return ls, nil
}
func (c *baseRunner) ReadFile(fs int, buildId string, pth string) (int64, io.ReadCloser, error) {
	if buildId == "" || pth == "" {
		return 0, nil, errors.New("param err")
	}
	build, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return 0, nil, errors.New("not found build")
	}
	pths := ""
	if fs == 1 {
		pths = filepath.Join(build.repoPaths, pth)
	} else if fs == 2 {
		pths = filepath.Join(build.buildPath, common.PathJobs, pth)
	}
	if pths == "" {
		return 0, nil, errors.New("path param err")
	}
	stat, err := os.Stat(pths)
	if err != nil {
		return 0, nil, err
	}
	fl, err := os.Open(pths)
	if err != nil {
		return 0, nil, err
	}
	return stat.Size(), fl, nil
}

func (c *baseRunner) GetEnv(jobid, key string) (string, bool) {
	if jobid == "" || key == "" {
		return "", false
	}
	job, ok := Mgr.jobEgn.GetJob(jobid)
	if !ok {
		return "", false
	}
	dir := filepath.Join(comm.WorkPath, common.PathBuild, job.step.BuildId, common.PathJobs, job.step.Id)
	bts, err := ioutil.ReadFile(filepath.Join(dir, "build.env"))
	if err != nil {
		return "", false
	}
	mp := hbtp.NewMaps(bts)
	v, ok := mp.Get(key)
	if !ok {
		return "", false
	}
	switch v.(type) {
	case string:
		return v.(string), true
	}
	return fmt.Sprintf("%v", v), true
}
func (c *baseRunner) GenEnv(jobid string, env utils.EnvVal) error {
	if jobid == "" || env == nil {
		return errors.New("param err")
	}
	job, ok := Mgr.jobEgn.GetJob(jobid)
	if !ok {
		return errors.New("not found job")
	}
	bts, err := json.Marshal(env)
	if err != nil {
		return err
	}
	dir := filepath.Join(comm.WorkPath, common.PathBuild, job.step.BuildId, common.PathJobs, job.step.Id)
	err = ioutil.WriteFile(filepath.Join(dir, "build.env"), bts, 0640)
	return err
}

func (c *baseRunner) UploadFile(jobid string, name, pth string) (io.WriteCloser, error) {
	if jobid == "" || pth == "" {
		return nil, errors.New("param err")
	}
	job, ok := Mgr.jobEgn.GetJob(jobid)
	if !ok {
		return nil, errors.New("not found job")
	}
	pths := filepath.Join(job.task.buildPath, common.PathJobs, job.step.Id, common.PathArts, name, pth)
	dir := filepath.Dir(pths)
	os.MkdirAll(dir, 0750)
	fl, err := os.OpenFile(pths, os.O_CREATE|os.O_RDWR, 0640)
	/*if err!=nil{
		return nil,err
	}*/
	return fl, err
}
