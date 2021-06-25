package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/runner/runners"
	"github.com/sirupsen/logrus"
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
func (c *baseRunner) Update(m *runners.UpdateJobInfo) error {
	Mgr.jobEgn.joblk.RLock()
	job, ok := Mgr.jobEgn.jobs[m.Id]
	Mgr.jobEgn.joblk.RUnlock()
	if !ok {
		return errors.New("not found")
	}
	job.Lock()
	job.step.Status = m.Status
	job.step.Error = m.Error
	job.step.ExitCode = m.ExitCode
	job.Unlock()
	return nil
}
func (c *baseRunner) CheckCancel(buildId string) bool {
	Mgr.buildEgn.tskslk.RLock()
	defer Mgr.buildEgn.tskslk.RUnlock()
	v, ok := Mgr.buildEgn.tasks[buildId]
	if !ok {
		return true
	}
	return v.stopd()
}

func (c *baseRunner) UpdateCmd(jobid, cmdid string, fs int) error {
	Mgr.jobEgn.joblk.RLock()
	job, ok := Mgr.jobEgn.jobs[jobid]
	Mgr.jobEgn.joblk.RUnlock()
	if !ok {
		return errors.New("not found")
	}
	job.Lock()
	cmd, ok := job.cmdmp[cmdid]
	job.Unlock()
	if !ok {
		return errors.New("not found")
	}
	logrus.Debugf("UpdateCmd cmdid:%s", cmd.Id)
	return nil
}
func (c *baseRunner) PushOutLine(jobid, cmdid, bs string, iserr bool) error {
	Mgr.jobEgn.joblk.RLock()
	job, ok := Mgr.jobEgn.jobs[jobid]
	Mgr.jobEgn.joblk.RUnlock()
	if !ok {
		return errors.New("not found")
	}
	job.Lock()
	buildid := job.step.BuildId
	job.Unlock()
	/*Mgr.buildEgn.tskslk.Lock()
	task, ok := Mgr.buildEgn.tasks[buildid]
	Mgr.buildEgn.tskslk.Unlock()
	if !ok {
		return errors.New("not found")
	}*/

	bts, err := json.Marshal(&comm.LogOutJson{
		Id:      cmdid,
		Content: bs,
		Times:   time.Now(),
		Errs:    iserr,
	})
	if err != nil {
		return err
	}

	dir := filepath.Join(comm.WorkPath, common.PathBuild, buildid, common.PathJobs)
	logpth := filepath.Join(dir, fmt.Sprintf("%v.log", jobid))
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
