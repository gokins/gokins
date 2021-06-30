package engine

import (
	"context"
	"errors"
	"fmt"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/util"
	"github.com/gokins-main/runner/runners"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	ghttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/gokins/comm"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
)

type taskStage struct {
	sync.RWMutex
	wg    sync.WaitGroup
	stage *runtime.Stage
	jobs  map[string]*jobSync
}

func (c *taskStage) status(stat, errs string, event ...string) {
	c.Lock()
	defer c.Unlock()
	c.stage.Status = stat
	c.stage.Error = errs
	if len(event) > 0 {
		c.stage.Event = event[0]
	}
}

type BuildTask struct {
	egn   *BuildEngine
	ctx   context.Context
	cncl  context.CancelFunc
	bdlk  sync.RWMutex
	build *runtime.Build

	bngtm     time.Time
	endtm     time.Time
	ctrlendtm time.Time

	staglk sync.RWMutex
	stages map[string]*taskStage

	buildPath string
	isClone   bool
	repoPath  string
}

func (c *BuildTask) status(stat, errs string, event ...string) {
	c.bdlk.Lock()
	defer c.bdlk.Unlock()
	c.build.Status = stat
	c.build.Error = errs
	if len(event) > 0 {
		c.build.Event = event[0]
	}
}

func NewBuildTask(egn *BuildEngine, bd *runtime.Build) *BuildTask {
	c := &BuildTask{egn: egn, build: bd}
	return c
}

func (c *BuildTask) stopd() bool {
	if c.ctx == nil {
		return true
	}
	return hbtp.EndContext(c.ctx)
}
func (c *BuildTask) stop() {
	c.ctrlendtm = time.Time{}
	if c.cncl != nil {
		c.cncl()
	}
}
func (c *BuildTask) Cancel() {
	c.ctrlendtm = time.Now()
	if c.cncl != nil {
		c.cncl()
	}
}
func (c *BuildTask) run() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask run recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

	defer func() {
		c.endtm = time.Now()
		c.build.Finished = time.Now()
		c.updateBuild(c.build)
		if c.isClone {
			os.RemoveAll(c.repoPath)
		}
	}()

	c.buildPath = filepath.Join(comm.WorkPath, common.PathBuild, c.build.Id)
	err := os.MkdirAll(c.buildPath, 0750)
	if err != nil {
		c.status(common.BuildStatusError, "build path err:"+err.Error(), common.BuildEventPath)
		return
	}

	c.bngtm = time.Now()
	c.stages = make(map[string]*taskStage)

	c.build.Started = time.Now()
	c.build.Status = common.BuildStatusPending
	if !c.check() {
		c.build.Status = common.BuildStatusError
		return
	}
	c.ctx, c.cncl = context.WithTimeout(comm.Ctx, time.Hour*2+time.Minute*5)
	c.build.Status = common.BuildStatusPreparation
	err = c.getRepo()
	if err != nil {
		logrus.Errorf("clone repo err:%v", err)
		c.status(common.BuildStatusError, "repo err", common.BuildEventGetRepo)
		return
	}
	c.build.Status = common.BuildStatusRunning
	for _, v := range c.build.Stages {
		v.Status = common.BuildStatusPending
		for _, e := range v.Steps {
			e.Status = common.BuildStatusPending
		}
	}
	c.updateBuild(c.build)
	for _, v := range c.build.Stages {
		c.runStage(v)
		if v.Status != common.BuildStatusOk {
			c.build.Status = v.Status
			return
		}
	}
	c.build.Status = common.BuildStatusOk
}
func (c *BuildTask) check() bool {
	if c.build.Repo == nil {
		c.status(common.BuildEventCheckParam, "repo param err")
		return false
	}
	if c.build.Repo.CloneURL == "" {
		c.status(common.BuildEventCheckParam, "repo param err:clone url")
		return false
	}
	s, err := os.Stat(c.build.Repo.CloneURL)
	if err == nil && s.IsDir() {
		c.isClone = false
		c.repoPath = c.build.Repo.CloneURL
	} else {
		if !common.RegUrl.MatchString(c.build.Repo.CloneURL) {
			c.status(common.BuildEventCheckParam, "repo param err:clone url")
			return false
		}
		c.isClone = true
	}
	if c.build.Stages == nil || len(c.build.Stages) <= 0 {
		c.build.Event = common.BuildEventCheckParam
		c.build.Error = "build Stages is empty"
		return false
	}
	stages := make(map[string]*taskStage)
	for _, v := range c.build.Stages {
		if v.BuildId != c.build.Id {
			c.build.Event = common.BuildEventCheckParam
			c.build.Error = fmt.Sprintf("Stage Build id err:%s/%s", v.BuildId, c.build.Id)
			return false
		}
		if v.Name == "" {
			c.build.Event = common.BuildEventCheckParam
			c.build.Error = "build Stage name is empty"
			return false
		}
		if v.Steps == nil || len(v.Steps) <= 0 {
			c.build.Event = common.BuildEventCheckParam
			c.build.Error = "build Stages is empty"
			return false
		}
		if _, ok := stages[v.Name]; ok {
			c.build.Event = common.BuildEventCheckParam
			c.build.Error = fmt.Sprintf("build Stages.%s is repeat", v.Name)
			return false
		}
		vs := &taskStage{
			stage: v,
			jobs:  make(map[string]*jobSync),
		}
		stages[v.Name] = vs
		for _, e := range v.Steps {
			if e.BuildId != c.build.Id {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = fmt.Sprintf("Job Build id err:%s/%s", v.BuildId, c.build.Id)
				return false
			}
			if e.StageId != v.Id {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = fmt.Sprintf("Job Stage id err:%s/%s", v.BuildId, c.build.Id)
				return false
			}
			e.Step = strings.TrimSpace(e.Step)
			if e.Step == "" {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = "build Step Plugin is empty"
				return false
			}
			if e.Name == "" {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = "build Step name is empty"
				return false
			}
			if _, ok := vs.jobs[e.Name]; ok {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = fmt.Sprintf("build Job.%s is repeat", e.Name)
				return false
			}
			job := &jobSync{
				step:  e,
				cmdmp: make(map[string]*cmdSync),
			}
			if err := c.genCmds(job); err != nil {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = fmt.Sprintf("build Job.%s Commands err", e.Name)
				return false
			}
			vs.jobs[e.Name] = job
		}
	}
	/*for _,v:=range stages{
		for _,e:=range v.jobs{
			err:=Mgr.jobEgn.Put(e)
			if err!=nil{
				c.build.Event = common.BuildEventPutJob
				c.build.Error=err.Error()
				return false
			}
		}
	}*/

	for k, v := range stages {
		c.stages[k] = v
	}
	return true
}

func (c *BuildTask) genCmds(job *jobSync) error {
	runjb := &runners.RunJob{
		Id:              job.step.Id,
		StageId:         job.step.StageId,
		BuildId:         job.step.BuildId,
		Step:            job.step.Step,
		Name:            job.step.Name,
		Environments:    job.step.Environments,
		Artifacts:       job.step.Artifacts,
		DependArtifacts: job.step.DependArtifacts,
		IsClone:         c.isClone,
		RepoPath:        c.repoPath,
	}
	var err error
	switch job.step.Commands.(type) {
	case string:
		c.appendcmds(runjb, utils.NewXid(), job.step.Commands.(string))
	case []interface{}:
		err = c.gencmds(runjb, job.step.Commands.([]interface{}))
	case []string:
		var ls []interface{}
		ts := job.step.Commands.([]string)
		for _, v := range ts {
			ls = append(ls, v)
		}
		err = c.gencmds(runjb, ls)
	default:
		err = errors.New("commands format err")
	}
	if err != nil {
		return err
	}
	if len(runjb.Commands) <= 0 {
		return errors.New("command format empty")
	}
	job.runjb = runjb
	for i, v := range runjb.Commands {
		job.cmdmp[v.Id] = &cmdSync{
			cmd:    v,
			status: common.BuildStatusPending,
		}
		cmd := &model.TCmdLine{
			Id:      v.Id,
			GroupId: v.Gid,
			BuildId: job.step.BuildId,
			StepId:  job.step.Id,
			Status:  common.BuildStatusPending,
			Num:     i + 1,
			Content: v.Conts,
			Created: time.Now(),
		}
		_, err = comm.Db.InsertOne(cmd)
		if err != nil {
			comm.Db.Where("build_id=? and job_id=?", cmd.BuildId, cmd.StepId).Delete(cmd)
			return err
		}
	}
	return nil
}
func (c *BuildTask) runStage(stage *runtime.Stage) {
	defer func() {
		stage.Finished = time.Now()
		c.updateStage(stage)
		logrus.Debugf("stage %s end!!!", stage.Name)
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask runStage recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()
	stage.Started = time.Now()
	stage.Status = common.BuildStatusRunning
	//c.logfile.WriteString(fmt.Sprintf("\n****************Stage+ %s\n", stage.Name))
	c.updateStage(stage)
	c.staglk.RLock()
	stg, ok := c.stages[stage.Name]
	c.staglk.RUnlock()
	if !ok {
		stg.status(common.BuildStatusError, fmt.Sprintf("not found stage?:%s", stage.Name))
		return
	}

	c.staglk.RLock()
	for _, v := range stage.Steps {
		jb, ok := stg.jobs[v.Name]
		if !ok {
			jb.status(common.BuildStatusError, "")
			break
		}
		stg.wg.Add(1)
		go c.runStep(stg, jb)
	}
	c.staglk.RUnlock()
	stg.wg.Wait()
	for _, v := range stg.jobs {
		v.RLock()
		ign := v.step.ErrIgnore
		status := v.step.Status
		errs := v.step.Error
		v.RUnlock()
		if !ign && status == common.BuildStatusError {
			stg.status(status, errs)
			return
		} else if status == common.BuildStatusCancel {
			stg.status(status, errs)
			return
		}
	}

	stage.Status = common.BuildStatusOk
}
func (c *BuildTask) runStep(stage *taskStage, job *jobSync) {
	defer stage.wg.Done()
	defer func() {
		go c.updateStep(job)
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask runStep recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

	if len(job.runjb.Commands) <= 0 {
		job.status(common.BuildStatusError, "command format empty", common.BuildEventJobCmds)
		return
	}

	job.RLock()
	dendons := job.step.DependsOn
	job.RUnlock()
	if len(dendons) > 0 {
		ls := make([]*jobSync, 0)
		for _, v := range dendons {
			if v == "" {
				continue
			}
			e, ok := stage.jobs[v]
			//core.Log.Debugf("job(%s) depend %s(ok:%t)",job.step.Name,v,ok)
			if !ok {
				job.status(common.BuildStatusError, fmt.Sprintf("depend on %s not found", v))
				return
			}
			if e.step.Name == job.step.Name {
				job.status(common.BuildStatusError, fmt.Sprintf("depend on %s is your self", job.step.Name))
				return
			}
			ls = append(ls, e)
		}
		for !hbtp.EndContext(comm.Ctx) {
			time.Sleep(time.Millisecond * 100)
			if c.stopd() {
				job.status(common.BuildStatusCancel, "")
				return
			}
			waitln := len(ls)
			for _, v := range ls {
				v.Lock()
				vStats := v.step.Status
				v.Unlock()
				if vStats == common.BuildStatusOk {
					waitln--
				} else if vStats == common.BuildStatusCancel {
					job.status(common.BuildStatusCancel, "")
					return
				} else if vStats == common.BuildStatusError {
					if v.step.ErrIgnore {
						waitln--
					} else {
						job.status(common.BuildStatusError, fmt.Sprintf("depend on %s is err", v.step.Name))
						return
					}
				}
			}
			if waitln <= 0 {
				break
			}
		}
	}

	job.Lock()
	job.step.Status = common.BuildStatusPreparation
	job.step.Started = time.Now()
	job.Unlock()
	go c.updateStep(job)
	err := Mgr.jobEgn.Put(job)
	if err != nil {
		job.status(common.BuildStatusError, fmt.Sprintf("command run err:%v", err))
		return
	}
	for !hbtp.EndContext(comm.Ctx) {
		job.Lock()
		stats := job.step.Status
		job.Unlock()
		if common.BuildStatusEnded(stats) {
			break
		}
		if c.stopd() && time.Since(c.ctrlendtm).Seconds() > 3 {
			job.status(common.BuildStatusCancel, "cancel")
			break
		}
		/*if c.ctrlend && time.Since(c.ctrlendtm).Seconds() > 3 {
			job.status(common.BuildStatusError, "cancel")
			break
		}*/
		time.Sleep(time.Millisecond * 10)
	}
	/*job.Lock()
	defer job.Unlock()
	if c.ctrlend && job.step.Status == common.BuildStatusError {
		job.step.Status = common.BuildStatusCancel
	}*/
}

func (c *BuildTask) getRepo() error {
	if !c.isClone {
		return nil
	}
	clonePath := filepath.Join(c.buildPath, common.PathRepo)
	err := gitClone(c.ctx, clonePath, c.build.Repo)
	if err != nil {
		return err
	}
	return nil
}

func gitClone(ctx context.Context, dir string, repo *runtime.Repository) error {
	clonePath := filepath.Join(dir)
	bauth := &ghttp.BasicAuth{
		Username: repo.Name,
		Password: repo.Token,
	}
	gc := &git.CloneOptions{
		URL:  repo.CloneURL,
		Auth: bauth,
	}
	logrus.Debugf("gitClone : clone url: %s sha: %s", repo.CloneURL, repo.Sha)
	repository, err := util.CloneRepo(clonePath, gc, ctx)
	if err != nil {
		return err
	}
	if repo.Sha != "" {
		err = util.CheckOutHash(repository, repo.Sha)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *BuildTask) appendcmds(runjb *runners.RunJob, gid string, conts string) {
	if gid == "" {
		return
	}
	m := &runners.CmdContent{
		Id:    utils.NewXid(),
		Gid:   gid,
		Conts: conts,
		Times: time.Now(),
	}
	logrus.Debugf("append cmd(%d)-%s:%s", len(runjb.Commands), gid, m.Conts)
	//job.Commands[m.Id] = m
	runjb.Commands = append(runjb.Commands, m)
}
func (c *BuildTask) gencmds(runjb *runners.RunJob, cmds []interface{}) (rterr error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask gencmds recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
			rterr = fmt.Errorf("%v", err)
		}
	}()
	for _, v := range cmds {
		switch v.(type) {
		case string:
			gid := utils.NewXid()
			//grp:=&hbtpBean.CmdGroupJson{Id: utils.NewXid()}
			c.appendcmds(runjb, gid, v.(string))
		case []interface{}:
			gid := utils.NewXid()
			for _, v1 := range v.([]interface{}) {
				c.appendcmds(runjb, gid, fmt.Sprintf("%v", v1))
			}
		case map[interface{}]interface{}:
			for _, v1 := range v.(map[interface{}]interface{}) {
				gid := utils.NewXid()
				switch v1.(type) {
				case string:
					c.appendcmds(runjb, gid, fmt.Sprintf("%v", v1))
				case []interface{}:
					for _, v2 := range v1.([]interface{}) {
						c.appendcmds(runjb, gid, fmt.Sprintf("%v", v2))
					}
				}
			}
		case map[string]interface{}:
			for _, v1 := range v.(map[string]interface{}) {
				gid := utils.NewXid()
				switch v1.(type) {
				case string:
					c.appendcmds(runjb, gid, fmt.Sprintf("%v", v1))
				case []interface{}:
					for _, v2 := range v1.([]interface{}) {
						c.appendcmds(runjb, gid, fmt.Sprintf("%v", v2))
					}
				}
			}
		}
	}
	return nil
}

func (c *BuildTask) UpJob(job *jobSync, stat, errs string, code int) {
	if job == nil || stat == "" {
		return
	}
	job.Lock()
	job.step.Status = stat
	job.step.Error = errs
	job.step.ExitCode = code
	job.Unlock()
	go c.updateStep(job)
}
func (c *BuildTask) UpJobCmd(job *jobSync, cmdid string, fs int) {
	if job == nil || cmdid == "" {
		return
	}
	job.RLock()
	cmd, ok := job.cmdmp[cmdid]
	job.RUnlock()
	if !ok {
		return
	}
	cmd.Lock()
	defer cmd.Unlock()
	switch fs {
	case 1:
		cmd.status = common.BuildStatusRunning
		cmd.started = time.Now()
	case 2:
		cmd.status = common.BuildStatusOk
		cmd.finished = time.Now()
	default:
		return
	}
	go c.updateStepCmd(cmd)
}
func (c *BuildTask) Show() (*runtime.BuildShow, bool) {
	if c.stopd() {
		return nil, false
	}
	rtbd := &runtime.BuildShow{
		Id:         c.build.Id,
		PipelineId: c.build.PipelineId,
		Status:     c.build.Status,
		Error:      c.build.Error,
		Event:      c.build.Event,
		Started:    c.build.Started,
		Finished:   c.build.Finished,
		Created:    c.build.Created,
		Updated:    c.build.Updated,
	}
	for _, v := range c.build.Stages {
		c.staglk.RLock()
		stg, ok := c.stages[v.Name]
		c.staglk.RUnlock()
		if !ok {
			continue
		}
		stg.RLock()
		rtstg := &runtime.StageShow{
			Id:       stg.stage.Id,
			BuildId:  stg.stage.BuildId,
			Status:   stg.stage.Status,
			Event:    stg.stage.Event,
			Error:    stg.stage.Error,
			Started:  stg.stage.Started,
			Stopped:  stg.stage.Stopped,
			Finished: stg.stage.Finished,
			Created:  stg.stage.Created,
			Updated:  stg.stage.Updated,
		}
		stg.RUnlock()
		rtbd.Stages = append(rtbd.Stages, rtstg)
		for _, st := range v.Steps {
			c.staglk.RLock()
			job, ok := stg.jobs[st.Name]
			c.staglk.RUnlock()
			if !ok {
				continue
			}
			job.RLock()
			rtstp := &runtime.StepShow{
				Id:       job.step.Id,
				StageId:  job.step.StageId,
				BuildId:  job.step.BuildId,
				Status:   job.step.Status,
				Event:    job.step.Event,
				Error:    job.step.Error,
				ExitCode: job.step.ExitCode,
				Started:  job.step.Started,
				Stopped:  job.step.Stopped,
				Finished: job.step.Finished,
			}
			rtstg.Steps = append(rtstg.Steps, rtstp)
			for _, cmd := range job.cmdmp {
				rtstp.Cmds = append(rtstp.Cmds, &runtime.CmdShow{
					Id:       cmd.cmd.Id,
					Status:   cmd.status,
					Started:  cmd.started,
					Finished: cmd.finished,
				})
			}
			job.RUnlock()
		}
	}
	return rtbd, true
}
