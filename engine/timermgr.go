package engine

import (
	"encoding/json"
	"errors"
	"github.com/gokins/core/common"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/model"
	"github.com/gokins/gokins/service"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"sync"
	"time"
)

type TimerEngine struct {
	tasklk sync.RWMutex
	tasks  map[string]*timerExec
}
type timerExec struct {
	tt   *model.TTrigger
	typ  int64
	tms  time.Time
	tick time.Time
}

func StartTimerEngine() *TimerEngine {
	c := &TimerEngine{
		tasks: make(map[string]*timerExec),
	}
	go func() {
		c.refresh()
		for !hbtp.EndContext(comm.Ctx) {
			c.run()
			time.Sleep(time.Millisecond * 10)
		}
	}()
	return c
}
func (c *TimerEngine) run() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("TimerEngine run recover:%v", err)
			logrus.Warnf("TimerEngine stack:%s", string(debug.Stack()))
		}
	}()

	c.tasklk.RLock()
	defer c.tasklk.RUnlock()
	for _, v := range c.tasks {
		c.execItem(v)
	}
}
func (c *TimerEngine) execItem(v *timerExec) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("TimerEngine execItem recover:%v", err)
			logrus.Warnf("TimerEngine stack:%s", string(debug.Stack()))
		}
	}()
	if time.Since(v.tick) > 0 {
		now := time.Now()
		logrus.Debugf("Timer(%s[%d]:%s) tick on:%s", v.tt.Name, v.typ, now.Format(common.TimeFmt), v.tick.Format(common.TimeFmt))
		switch v.typ {
		case 0:
			go c.Delete(v.tt.Id)
			time.Sleep(time.Millisecond * 10)
		case 1:
			v.tick = now.Add(time.Minute)
		case 2:
			v.tick = now.Add(time.Hour)
		case 3:
			v.tick = now.Add(time.Hour * 24)
		case 4:
			v.tick = now.Add(time.Hour * 24 * 7)
		case 5:
			v.tick = now.Add(time.Hour * 24 * 30)
		}

		rb, err := service.TriggerTimer(v.tt)
		if err != nil {
			logrus.Errorf("TriggerTimer err:%v", err)
		} else {
			Mgr.BuildEgn().Put(rb)
		}
	}
}

func (c *TimerEngine) refresh() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("TimerEngine refresh recover:%v", err)
			logrus.Warnf("TimerEngine stack:%s", string(debug.Stack()))
		}
	}()
	var ls []*model.TTrigger
	comm.Db.Where("enabled = 1 AND types = 'timer'").Find(&ls)

	c.tasklk.Lock()
	defer c.tasklk.Unlock()
	for _, v := range ls {
		err := c.resetOne(v)
		if err != nil {
			logrus.Errorf("TimerEngine resetOne err:%v", err)
		}
	}
}
func (c *TimerEngine) resetOne(tmr *model.TTrigger) error {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("TimerEngine refresh recover:%v", err)
			logrus.Warnf("TimerEngine stack:%s", string(debug.Stack()))
		}
	}()
	if tmr.Types != "timer" {
		return errors.New("type is err:" + tmr.Types)
	}
	mp := hbtp.Map{}
	err := json.Unmarshal([]byte(tmr.Params), &mp)
	if err != nil {
		return err
	}
	typ, err := mp.GetInt("timerType")
	if err != nil {
		return err
	}
	dates := mp.GetString("dates")
	tms, err := time.ParseInLocation(time.RFC3339Nano, dates, time.Local)
	if err != nil {
		return err
	}
	switch typ {
	case 0:
		if time.Since(tms) < 0 {
			t, ok := c.tasks[tmr.Id]
			if !ok {
				t = &timerExec{
					tt: tmr,
				}
				c.tasks[tmr.Id] = t
			}
			t.typ = typ
			t.tms = tms
			t.tick = tms
			logrus.Debugf("Timer add(%s[%d]:%s) tick on:%s", tmr.Name, typ, tms.Format(common.TimeFmt), t.tick.Format(common.TimeFmt))
		}
	case 1, 2, 3, 4, 5:
		now := time.Now()
		t, ok := c.tasks[tmr.Id]
		if !ok {
			t = &timerExec{
				tt: tmr,
			}
			c.tasks[tmr.Id] = t
		}
		t.typ = typ
		t.tms = tms
		switch typ {
		case 1:
			t.tick = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), tms.Second(), 0, time.Local)
			if time.Since(t.tick) > 0 {
				t.tick = t.tick.Add(time.Minute)
			}
		case 2:
			t.tick = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), tms.Minute(), tms.Second(), 0, time.Local)
			if time.Since(t.tick) > 0 {
				t.tick = t.tick.Add(time.Hour)
			}
		case 3:
			t.tick = time.Date(now.Year(), now.Month(), now.Day(), tms.Hour(), tms.Minute(), tms.Second(), 0, time.Local)
			if time.Since(t.tick) > 0 {
				t.tick = t.tick.Add(time.Hour * 24)
			}
		case 4:
			t.tick = time.Date(now.Year(), now.Month(), tms.Day(), tms.Hour(), tms.Minute(), tms.Second(), 0, time.Local)
			if time.Since(t.tick) > 0 {
				t.tick = t.tick.Add(time.Hour * 24 * 7)
			}
		case 5:
			t.tick = time.Date(now.Year(), now.Month(), tms.Day(), tms.Hour(), tms.Minute(), tms.Second(), 0, time.Local)
			if time.Since(t.tick) > 0 {
				t.tick = t.tick.Add(time.Hour * 24 * 30)
			}
		}
		logrus.Debugf("Timer add(%s[%d]:%s) tick on:%s", tmr.Name, typ, tms.Format(common.TimeFmt), t.tick.Format(common.TimeFmt))
	}
	return nil
}
func (c *TimerEngine) Refresh(tmrid string) error {
	if tmrid == "" {
		return errors.New("param err")
	}
	tmr := &model.TTrigger{}
	ok, _ := comm.Db.Where("id=?", tmrid).Get(tmr)
	if !ok || tmr.Enabled != 1 {
		c.Delete(tmrid)
		return errors.New("not found")
	}
	return c.resetOne(tmr)
}
func (c *TimerEngine) Delete(tmrid string) {
	c.tasklk.Lock()
	delete(c.tasks, tmrid)
	c.tasklk.Unlock()
}
