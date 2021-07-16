package engine

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/runner/runners"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"sync"
	"time"
)

type JobEngine struct {
	tmr   *utils.Timer
	exelk sync.RWMutex
	execs map[string]*executer
	joblk sync.RWMutex
	jobs  map[string]*jobSync
}
type executer struct {
	sync.RWMutex
	plug  string
	tms   time.Time
	jobwt *list.List
}
type cmdSync struct {
	sync.RWMutex
	cmd      *runners.CmdContent
	code     int
	status   string
	started  time.Time
	finished time.Time
}
type jobSync struct {
	sync.RWMutex
	task  *BuildTask
	step  *runtime.Step
	runjb *runners.RunJob
	cmdmp map[string]*cmdSync
	ended bool
}

func (c *jobSync) status(stat, errs string, event ...string) {
	c.Lock()
	defer c.Unlock()
	c.step.Status = stat
	c.step.Error = errs
	if len(event) > 0 {
		c.step.Event = event[0]
	}
}

func StartJobEngine() *JobEngine {
	c := &JobEngine{
		tmr:   utils.NewTimer(time.Second * 30),
		execs: make(map[string]*executer),
		jobs:  make(map[string]*jobSync),
	}
	go func() {
		for !hbtp.EndContext(comm.Ctx) {
			c.run()
			time.Sleep(time.Second)
		}
	}()
	return c
}
func (c *JobEngine) run() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("JobEngine run recover:%v", err)
			logrus.Warnf("JobEngine stack:%s", string(debug.Stack()))
		}
	}()

	if !c.tmr.Tick() {
		return
	}
	func() {
		c.exelk.RLock()
		defer c.exelk.RUnlock()
		for k, v := range c.execs {
			v.RLock()
			if time.Since(v.tms).Minutes() > 5 {
				go c.rmExec(k, v)
			}
			v.RUnlock()
		}
	}()
	func() {
		c.joblk.Lock()
		defer c.joblk.Unlock()
		for k, v := range c.jobs {
			if v.ended {
				delete(c.jobs, k)
			}
		}
	}()
}
func (c *JobEngine) rmExec(k string, ex *executer) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("JobEngine stopsJob recover:%v", err)
			logrus.Warnf("JobEngine stack:%s", string(debug.Stack()))
		}
	}()

	c.exelk.Lock()
	defer c.exelk.Unlock()
	ex.Lock()
	defer ex.Unlock()
	for e := ex.jobwt.Front(); e != nil; e = e.Next() {
		job := e.Value.(*jobSync)
		job.ended = true
	}
	delete(c.execs, k)
}
func (c *JobEngine) Put(job *jobSync) error {
	if job == nil || job.step.Step == "" {
		return errors.New("step plugin empty")
	}
	c.exelk.RLock()
	e, ok := c.execs[job.step.Step]
	c.exelk.RUnlock()
	if !ok {
		return fmt.Errorf("Not Found Plugin:%s", job.step.Step)
	}
	e.Lock()
	defer e.Unlock()
	e.jobwt.PushBack(job)
	return nil
}
func (c *JobEngine) Pull(name string, plugs []string) *runners.RunJob {
	for _, v := range plugs {
		if v == "" {
			continue
		}
		c.exelk.RLock()
		ex, ok := c.execs[v]
		c.exelk.RUnlock()
		if !ok {
			ex = &executer{
				plug:  v,
				tms:   time.Now(),
				jobwt: list.New(),
			}
			c.exelk.Lock()
			c.execs[v] = ex
			c.exelk.Unlock()
		}
		var job *jobSync
		ex.Lock()
		ex.tms = time.Now()
		e := ex.jobwt.Front()
		if e != nil {
			job = e.Value.(*jobSync)
			ex.jobwt.Remove(e)
			c.joblk.Lock()
			c.jobs[job.step.Id] = job
			c.joblk.Unlock()
		}
		ex.Unlock()
		if job != nil {
			return job.runjb
		}
	}
	return nil
}

/*func (c *JobEngine) GetJob(id string) (*jobSync, bool) {
	if id == "" {
		return nil, false
	}
	c.joblk.RLock()
	defer c.joblk.RUnlock()
	job, ok := c.jobs[id]
	return job, ok
}*/
