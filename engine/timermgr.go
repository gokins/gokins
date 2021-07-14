package engine

import (
	"container/list"
	"context"
	"encoding/json"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/models"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"sync"
	"time"
)

type (
	TimerEngine struct {
		lk        sync.Mutex
		tasks     map[string]*TimerTask
		ls        list.List
		ctx       context.Context
		ctxCancel context.CancelFunc
	}
)

func StartTimerEngine() *TimerEngine {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("TimerEngine run Start:%v", err)
			logrus.Warnf("TimerEngine stack:%s", string(debug.Stack()))
		}
	}()
	c := &TimerEngine{}
	ctx, cancelFunc := context.WithCancel(comm.Ctx)
	c.ctx = ctx
	c.ctxCancel = cancelFunc
	c.initTasks()
	go func() {
		for !hbtp.EndContext(c.ctx) {
			c.refreshTasks()
			time.Sleep(time.Millisecond * 200)
		}
	}()
	go func() {
		for !hbtp.EndContext(c.ctx) {
			c.exec()
			time.Sleep(time.Second * 1)
		}
	}()
	return c
}

func (c *TimerEngine) initTasks() {
	c.tasks = make(map[string]*TimerTask)
	var ls []*models.TimerTriggerRun
	sql := "SELECT `tt`.*, " +
		" ( SELECT created FROM `t_trigger_run` WHERE tt.id = t_trigger_run.tid ORDER BY created DESC LIMIT 1 ) AS r_created " +
		" FROM t_trigger AS tt " +
		" WHERE ( enabled != 0 AND types = 'timer' );"
	_ = comm.Db.SQL(sql).Find(&ls)
	for _, v := range ls {
		date := v.RunCreated
		tys := 0
		if date.IsZero() {
			param := &timerParam{}
			err := json.Unmarshal([]byte(v.Params), param)
			if err != nil {
				logrus.Debugf("initTasks err:%v", err)
				continue
			}
			if param.Date.IsZero() {
				logrus.Debugf("initTasks err:%v time config is empty", v.Name)
				continue
			}
			tys = param.TimerType
			date = param.Date
		}
		ctx, cfn := context.WithCancel(c.ctx)
		c.tasks[v.Id] = &TimerTask{
			ctx:       ctx,
			ctxCancel: cfn,
			date:      date,
			timerType: tys,
			end:       false,
		}
	}
}

func (c *TimerEngine) AddTask(tt *model.TTrigger) {
	c.lk.Lock()
	defer c.lk.Unlock()
	if task, ok := c.tasks[tt.Id]; ok {
		task.stop()
		delete(c.tasks, tt.Id)
	}
	param := &timerParam{}
	err := json.Unmarshal([]byte(tt.Params), param)
	if err != nil {
		logrus.Debugf("addTask err:%v", err)
		return
	}
	if param.Date.IsZero() {
		logrus.Debugf("addTask err:%v time config is empty", tt.Name)
		return
	}
	ctx, can := context.WithCancel(c.ctx)
	c.tasks[tt.Id] = &TimerTask{
		ctx:       ctx,
		ctxCancel: can,
		date:      param.Date,
		timerType: param.TimerType,
		end:       false,
	}
}

func (c *TimerEngine) refreshTasks() {
	c.lk.Lock()
	defer c.lk.Unlock()
	for k, task := range c.tasks {
		if task.end || hbtp.EndContext(task.ctx) {
			delete(c.tasks, k)
			continue
		}
		if isMatch(task) {
			task.date = time.Now()
			c.ls.PushBack(task)
		}
	}
}

func (c *TimerEngine) RemoveTasks(id string) {
	c.lk.Lock()
	defer c.lk.Unlock()
	if task, ok := c.tasks[id]; ok {
		task.stop()
		delete(c.tasks, id)
	}
}

func (c *TimerEngine) exec() {
	c.lk.Lock()
	defer c.lk.Unlock()
	for e := c.ls.Front(); e != nil; {
		task := e.Value.(*TimerTask)
		next := e.Next()
		c.ls.Remove(e)
		e = next
		if !task.end || !hbtp.EndContext(task.ctx) {
			task.run()
		}
	}

}

func (c *TimerEngine) Cancel() {
	if c.ctxCancel != nil {
		c.ctxCancel()
	}
}
