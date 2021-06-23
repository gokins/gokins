package engine

import (
	"container/list"
	"context"
	"errors"
	"fmt"
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
	"github.com/gokins-main/gokins/util/gitex"
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
	build *runtime.Build
	ctx   context.Context
	cncl  context.CancelFunc

	bngtm     time.Time
	endtm     time.Time
	ctrlend   bool //手动停止
	ctrlendtm time.Time

	staglk sync.Mutex
	stages map[string]*taskStage

	joblk sync.Mutex
	jobls *list.List

	buildPath string
	isClone   bool
	repoPath  string
}

func (c *BuildTask) status(stat, errs string, event ...string) {
	//c.Lock()
	//defer c.Unlock()
	c.build.Status = stat
	c.build.Error = errs
	if len(event) > 0 {
		c.build.Event = event[0]
	}
}

func NewBuildTask(egn *BuildEngine, bd *runtime.Build) *BuildTask {
	c := &BuildTask{egn: egn, build: bd}
	c.ctx, c.cncl = context.WithTimeout(comm.Ctx, time.Hour*2+time.Minute*5)
	return c
}

func (c *BuildTask) stopd() bool {
	return hbtp.EndContext(c.ctx)
}
func (c *BuildTask) stop() {
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
		c.updateBuild()
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
	c.jobls = list.New()

	c.build.Started = time.Now()
	c.build.Status = common.BuildStatusPending
	if !c.check() {
		c.build.Status = common.BuildStatusError
		return
	}
	c.build.Status = common.BuildStatusPreparation
	err = c.getRepo()
	if err != nil {
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
	c.updateBuild()
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
		if _, ok := c.stages[v.Name]; ok {
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
			vs.jobs[e.Name] = &jobSync{
				step: e,
			}
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
	c.staglk.Lock()
	stg, ok := c.stages[stage.Name]
	c.staglk.Unlock()
	if !ok {
		stg.status(common.BuildStatusError, fmt.Sprintf("not found stage?:%s", stage.Name))
		return
	}

	c.staglk.Lock()
	for _, v := range stage.Steps {
		jb, ok := stg.jobs[v.Name]
		if !ok {
			jb.status(common.BuildStatusError, "")
			break
		}
		stg.wg.Add(1)
		go c.runStep(stg, jb)
	}
	c.staglk.Unlock()
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
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask runStep recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

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
			if c.ctrlend {
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
	job.step.IsClone = c.isClone
	job.step.RepoPath = c.repoPath
	job.Unlock()
	c.updateStep(job.step)
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
		if c.ctrlend && time.Since(c.ctrlendtm).Seconds() > 3 {
			job.status(common.BuildStatusError, "cancel")
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	job.Lock()
	defer job.Unlock()
	if c.ctrlend && job.step.Status == common.BuildStatusError {
		job.step.Status = common.BuildStatusCancel
	}
}

func (c *BuildTask) getRepo() error {
	repo := c.build.Repo
	if repo == nil {
		return errors.New("getRepo err:  repo is empty ")
	}
	if repo.CloneURL == "" {
		return errors.New("getRepo err:  cloneUrl is empty ")
	}
	s, err := os.Stat(repo.CloneURL)
	if err == nil && s.IsDir() {
		c.isClone = false
		c.repoPath = repo.CloneURL
		return nil
	}

	clonePath := filepath.Join(c.buildPath, common.PathRepo)
	err = gitClone(c.ctx, clonePath, c.build.Repo)
	if err != nil {
		return err
	}
	c.isClone = true
	c.repoPath = clonePath
	return nil
}

func gitClone(ctx context.Context, dir string, repo *runtime.Repository) error {
	clonePath := filepath.Join(dir, common.PathRepo)
	bauth := &ghttp.BasicAuth{
		Username: repo.Name,
		Password: repo.Token,
	}
	gc := &git.CloneOptions{
		URL:  repo.CloneURL,
		Auth: bauth,
	}
	logrus.Debugf("gitClone : clone url: %s sha: %s", repo.CloneURL, repo.Sha)
	repository, err := gitex.CloneRepo(clonePath, gc, ctx)
	if err != nil {
		return err
	}
	if repo.Sha != "" {
		err = gitex.CheckOutHash(repository, repo.Sha)
		if err != nil {
			return err
		}
	}
	return nil
}
