package engine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	ghttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/util"
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
	stages map[string]*taskStage // key:name
	joblk  sync.RWMutex
	jobs   map[string]*jobSync //key:id

	buildPath string
	repoPaths string //fs
	workpgss  int

	isClone  bool
	repoPath string
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
func (c *BuildTask) clears() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask clears recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

	if c.isClone {
		os.RemoveAll(c.repoPaths)
	}
	for _, v := range c.jobs {
		pth := filepath.Join(c.buildPath, common.PathJobs, v.step.Id, common.PathArts)
		os.RemoveAll(pth)
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
		c.clears()
	}()

	c.buildPath = filepath.Join(comm.WorkPath, common.PathBuild, c.build.Id)
	c.repoPaths = filepath.Join(c.buildPath, common.PathRepo)
	err := os.MkdirAll(c.buildPath, 0750)
	if err != nil {
		c.status(common.BuildStatusError, "build path err:"+err.Error(), common.BuildEventPath)
		return
	}

	c.bngtm = time.Now()
	c.stages = make(map[string]*taskStage)
	c.jobs = make(map[string]*jobSync)

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
		c.status(common.BuildStatusError, "repo err:"+err.Error(), common.BuildEventGetRepo)
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
	logrus.Debugf("BuildTask run build:%s,pgss:%d", c.build.Id, c.workpgss)
	c.workpgss = 100
	for _, v := range c.build.Stages {
		c.runStage(v)
		if v.Status != common.BuildStatusOk {
			c.build.Status = v.Status
			return
		}
	}
	c.build.Status = common.BuildStatusOk
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
		stg.RLock()
		jb, ok := stg.jobs[v.Name]
		stg.RUnlock()
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
		job.ended = true
		job.step.Finished = time.Now()
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
	dendons := job.step.Waits
	job.RUnlock()
	if len(dendons) > 0 {
		ls := make([]*jobSync, 0)
		for _, v := range dendons {
			if v == "" {
				continue
			}
			stage.RLock()
			e, ok := stage.jobs[v]
			stage.RUnlock()
			//core.Log.Debugf("job(%s) depend %s(ok:%t)",job.step.Name,v,ok)
			if !ok {
				job.status(common.BuildStatusError, fmt.Sprintf("wait on %s not found", v))
				return
			}
			if e.step.Name == job.step.Name {
				job.status(common.BuildStatusError, fmt.Sprintf("wait on %s is your self", job.step.Name))
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
						job.status(common.BuildStatusError, fmt.Sprintf("wait on %s is err", v.step.Name))
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
	job.ended = false
	job.step.Status = common.BuildStatusPreparation
	job.step.Started = time.Now()
	job.Unlock()
	go c.updateStep(job)
	err := Mgr.jobEgn.Put(job)
	if err != nil {
		job.status(common.BuildStatusError, fmt.Sprintf("command run err:%v", err))
		return
	}
	logrus.Debugf("BuildTask put step:%s", job.step.Name)
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
	os.MkdirAll(c.repoPaths, 0750)
	err := c.gitClone(c.ctx, c.repoPaths, c.build.Repo)
	if err != nil {
		return err
	}
	return nil
}

var regBfb = regexp.MustCompile(`:\s+(\d+)% \(\d+\/\d+\)`)

func (c *BuildTask) Write(bts []byte) (n int, err error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask gitWrite recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()
	ln := len(bts)
	line := string(bts)
	if ln > 0 && regBfb.MatchString(line) {
		subs := regBfb.FindAllStringSubmatch(line, -1)[0]
		if len(subs) > 1 {
			p, err := strconv.Atoi(subs[1])
			if err == nil {
				c.workpgss = int(float64(p) * 0.8)
			}
		}
	}
	println("BuildTask git log:", line)
	return ln, nil
}
func (c *BuildTask) gitClone(ctx context.Context, dir string, repo *runtime.Repository) error {
	clonePath := filepath.Join(dir)
	gc := &git.CloneOptions{
		URL:      repo.CloneURL,
		Progress: c,
	}
	if repo.Name != "" {
		gc.Auth = &ghttp.BasicAuth{
			Username: repo.Name,
			Password: repo.Token,
		}
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
func (c *BuildTask) UpJobCmd(cmd *cmdSync, fs, code int) {
	if cmd == nil {
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
		if code != 0 {
			cmd.code = code
			cmd.status = common.BuildStatusError
		}
		cmd.finished = time.Now()
	case 3:
		cmd.code = code
		cmd.status = common.BuildStatusCancel
		cmd.finished = time.Now()
	case -1:
		cmd.code = code
		cmd.status = common.BuildStatusError
		cmd.finished = time.Now()
	default:
		return
	}
	go c.updateStepCmd(cmd)
}
func (c *BuildTask) WorkProgress() int {
	return c.workpgss
}
func (c *BuildTask) Show() (*runtime.BuildShow, bool) {
	if c.stopd() {
		return nil, false
	}
	c.bdlk.RLock()
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
	c.bdlk.RUnlock()
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

func (c *BuildTask) GetJob(id string) (*jobSync, bool) {
	if id == "" {
		return nil, false
	}
	c.joblk.RLock()
	defer c.joblk.RUnlock()
	job, ok := c.jobs[id]
	return job, ok
}
