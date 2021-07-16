package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gokins-main/gokins/service"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/runner/runners"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
)

type baseRunner struct{}

func (c *baseRunner) ServerInfo() (*runners.ServerInfo, error) {
	return &runners.ServerInfo{
		WebHost:   comm.Cfg.Server.Host,
		DownToken: comm.Cfg.Server.DownToken,
	}, nil
}

func (c *baseRunner) PullJob(name string, plugs []string) (*runners.RunJob, error) {
	tms := time.Now()
	for time.Since(tms).Seconds() < 5 {
		v := Mgr.jobEgn.Pull(name, plugs)
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
	tsk, ok := Mgr.buildEgn.Get(m.BuildId)
	if !ok {
		return errors.New("not found build")
	}
	job, ok := tsk.GetJob(m.JobId)
	if !ok {
		return errors.New("not found job")
	}
	tsk.UpJob(job, m.Status, m.Error, m.ExitCode)
	return nil
}

func (c *baseRunner) UpdateCmd(buildId, jobId, cmdId string, fs, code int) error {
	tsk, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return errors.New("not found build")
	}
	job, ok := tsk.GetJob(jobId)
	if !ok {
		return errors.New("not found job")
	}
	job.RLock()
	cmd, ok := job.cmdmp[cmdId]
	job.RUnlock()
	if !ok {
		return errors.New("not found cmd")
	}
	tsk.UpJobCmd(cmd, fs, code)
	return nil
}
func (c *baseRunner) PushOutLine(buildId, jobId, cmdId, bs string, iserr bool) error {
	tsk, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return errors.New("not found build")
	}
	job, ok := tsk.GetJob(jobId)
	if !ok {
		return errors.New("not found job")
	}

	bts, err := json.Marshal(&bean.LogOutJson{
		Id:      cmdId,
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
	stg.RLock()
	defer stg.RUnlock()
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
		pths = filepath.Join(comm.WorkPath, common.PathArtifacts, pth)
	} else if fs == 3 {
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
		pths = filepath.Join(comm.WorkPath, common.PathArtifacts, pth)
	} else if fs == 3 {
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

func (c *baseRunner) GetEnv(buildId, jobId, key string) (string, bool) {
	if jobId == "" || key == "" {
		return "", false
	}
	tsk, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return "", false
	}
	job, ok := tsk.GetJob(jobId)
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
func (c *baseRunner) GenEnv(buildId, jobId string, env utils.EnvVal) error {
	if jobId == "" || env == nil {
		return errors.New("param err")
	}
	tsk, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return errors.New("not found build")
	}
	job, ok := tsk.GetJob(jobId)
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

func (c *baseRunner) UploadFile(fs int, buildId, jobId string, dir, pth string) (io.WriteCloser, error) {
	if jobId == "" || pth == "" {
		return nil, errors.New("param err")
	}
	tsk, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return nil, errors.New("not found build")
	}
	job, ok := tsk.GetJob(jobId)
	if !ok {
		return nil, errors.New("not found job")
	}
	pths := ""
	if fs == 1 {
		pths = filepath.Join(comm.WorkPath, common.PathArtifacts, dir, pth)
	} else if fs == 2 {
		pths = filepath.Join(job.task.buildPath, common.PathJobs, job.step.Id, common.PathArts, dir, pth)
	}
	if pths == "" {
		return nil, errors.New("path param err")
	}
	dirs := filepath.Dir(pths)
	os.MkdirAll(dirs, 0750)
	fl, err := os.OpenFile(pths, os.O_CREATE|os.O_RDWR, 0640)
	/*if err!=nil{
		return nil,err
	}*/
	return fl, err
}
func (c *baseRunner) FindArtVersionId(buildId, idnt string, names string) (string, error) {
	tnms := strings.Split(strings.TrimSpace(names), "@")
	name := tnms[0]
	vers := ""
	if len(tnms) > 1 {
		vers = tnms[1]
	}
	if buildId == "" || idnt == "" || name == "" {
		return "", errors.New("param err")
	}
	build, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return "", errors.New("not found build")
	}

	arty := &model.TArtifactory{}
	ok, _ = comm.Db.Where("deleted!=1 and identifier=? and org_id in (select org_id from t_org_pipe where pipe_id=?)",
		idnt, build.build.PipelineId).Get(arty)
	if !ok {
		return "", errors.New("not found artifactory")
	}

	pv := &model.TPipelineVersion{}
	ok = service.GetIdOrAid(build.build.PipelineVersionId, pv)
	if !ok {
		return "", errors.New("not found pv")
	}
	usr := &model.TUser{}
	ok = service.GetIdOrAid(pv.Uid, usr)
	if !ok {
		return "", errors.New("not found user")
	}
	perm := service.NewOrgPerm(usr, arty.OrgId)
	if !perm.CanExec() {
		return "", fmt.Errorf("user put '%s' no permission", idnt)
	}

	artp := &model.TArtifactPackage{}
	ok, _ = comm.Db.Where("deleted!=1 and repo_id=? and name=?", arty.Id, name).Get(artp)
	if !ok {
		return "", fmt.Errorf("not found artifact '%s'", names)
	}
	artv := &model.TArtifactVersion{}
	ses := comm.Db.Where("package_id=?", artp.Id)
	if vers != "" {
		ses.And("version=? or sha=?", vers)
	}
	ok, _ = ses.OrderBy("aid DESC").Get(artv)
	if !ok {
		return "", fmt.Errorf("not found artifacts '%s'", names)
	}
	return artv.Id, nil
}
func (c *baseRunner) NewArtVersionId(buildId, idnt string, name string) (string, error) {
	name = strings.Split(strings.TrimSpace(name), "@")[0]
	if buildId == "" || idnt == "" || name == "" {
		return "", errors.New("param err")
	}
	build, ok := Mgr.buildEgn.Get(buildId)
	if !ok {
		return "", errors.New("not found build")
	}

	arty := &model.TArtifactory{}
	ok, _ = comm.Db.Where("deleted!=1 and identifier=? and org_id in (select org_id from t_org_pipe where pipe_id=?)",
		idnt, build.build.PipelineId).Get(arty)
	if !ok {
		return "", errors.New("not found artifactory")
	}
	if arty.Disabled == 1 {
		return "", errors.New("artifactory already disabled")
	}

	artp := &model.TArtifactPackage{}
	ok, _ = comm.Db.Where("deleted!=1 and repo_id=? and name=?", arty.Id, name).Get(artp)
	if !ok {
		artp.Id = utils.NewXid()
		artp.RepoId = arty.Id
		artp.Name = name
		artp.Created = time.Now()
		artp.Updated = time.Now()
		_, err := comm.Db.InsertOne(artp)
		if err != nil {
			return "", err
		}
	}
	artv := &model.TArtifactVersion{
		Id:        utils.NewXid(),
		RepoId:    arty.Id,
		PackageId: artp.Id,
		Name:      artp.Name,
		Preview:   1,
		Created:   time.Now(),
		Updated:   time.Now(),
	}
	artv.Sha = artv.Id
	_, err := comm.Db.InsertOne(artv)
	if err != nil {
		return "", err
	}
	return artv.Id, nil
}
