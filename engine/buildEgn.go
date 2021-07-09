package engine

import (
	"container/list"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/gokins/comm"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"sync"
	"time"
)

type BuildEngine struct {
	tskwlk sync.RWMutex
	taskw  *list.List

	tskslk sync.RWMutex
	tasks  map[string]*BuildTask
}

func StartBuildEngine() *BuildEngine {
	if comm.Cfg.Server.RunLimit < 2 {
		comm.Cfg.Server.RunLimit = 5
	}
	c := &BuildEngine{
		taskw: list.New(),
		tasks: make(map[string]*BuildTask),
	}
	go func() {
		c.init()
		for !hbtp.EndContext(comm.Ctx) {
			c.run()
			time.Sleep(time.Second)
		}
	}()
	return c
}
func (c *BuildEngine) Stop() {
	c.tskslk.RLock()
	defer c.tskslk.RUnlock()
	for _, v := range c.tasks {
		v.stop()
	}
}
func (c *BuildEngine) init() {

	/*// TODO: 调试不执行
	if comm.Debugs {
		return
	}*/

	cont := "server restart"
	comm.Db.Exec(
		"update `t_build` set `status`=?,`error`=? where `status`!=? and `status`!=? and `status`!=?",
		common.BuildStatusCancel, cont, common.BuildStatusOk, common.BuildStatusError, common.BuildStatusCancel,
	)
	comm.Db.Exec(
		"update `t_stage` set `status`=?,`error`=? where `status`!=? and `status`!=? and `status`!=?",
		common.BuildStatusCancel, cont, common.BuildStatusOk, common.BuildStatusError, common.BuildStatusCancel,
	)
	comm.Db.Exec(
		"update `t_step` set `status`=?,`error`=? where `status`!=? and `status`!=? and `status`!=?",
		common.BuildStatusCancel, cont, common.BuildStatusOk, common.BuildStatusError, common.BuildStatusCancel,
	)
	comm.Db.Exec(
		"update `t_cmd_line` set `status`=? where `status`!=? and `status`!=? and `status`!=?",
		common.BuildStatusCancel, common.BuildStatusOk, common.BuildStatusError, common.BuildStatusCancel,
	)
}

func (c *BuildEngine) run() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildEngine run recover:%v", err)
			logrus.Warnf("BuildEngine stack:%s", string(debug.Stack()))
		}
	}()

	c.tskwlk.RLock()
	ln1 := c.taskw.Len()
	c.tskwlk.RUnlock()
	c.tskslk.RLock()
	ln2 := len(c.tasks)
	c.tskslk.RUnlock()
	if ln1 > 0 && ln2 < comm.Cfg.Server.RunLimit {
		c.tskwlk.RLock()
		e := c.taskw.Front()
		c.tskwlk.RUnlock()
		if e == nil {
			return
		}
		c.tskwlk.Lock()
		c.taskw.Remove(e)
		c.tskwlk.Unlock()
		v := NewBuildTask(c, e.Value.(*runtime.Build))
		c.tskslk.Lock()
		c.tasks[v.build.Id] = v
		c.tskslk.Unlock()
		go c.startBuild(v)
	}
}
func (c *BuildEngine) startBuild(v *BuildTask) {
	v.run()
	c.tskslk.Lock()
	defer c.tskslk.Unlock()
	delete(c.tasks, v.build.Id)
}
func (c *BuildEngine) Put(bd *runtime.Build) {
	c.tskwlk.Lock()
	defer c.tskwlk.Unlock()
	c.taskw.PushBack(bd)
}
func (c *BuildEngine) Get(buildid string) (*BuildTask, bool) {
	if buildid == "" {
		return nil, false
	}
	c.tskslk.Lock()
	defer c.tskslk.Unlock()
	v, ok := c.tasks[buildid]
	return v, ok
}
