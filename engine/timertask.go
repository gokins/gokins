package engine

import (
	"context"
	"fmt"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"time"
)

type (
	TimerTask struct {
		ctx       context.Context
		ctxCancel context.CancelFunc
		date      time.Time
		timerType int
		end       bool
	}
	timerParam struct {
		TimerType int       `json:"timerType"`
		Date      time.Time `json:"dates"`
	}
)

func (c *TimerTask) run() {
	if hbtp.EndContext(c.ctx) {
		return
	}
	if c.timerType == 0 {
		defer c.stop()
		c.end = true
	}
	//TODO do job
	fmt.Println("job type :", c.timerType, " job time:", time.Now())
}

func (c *TimerTask) stop() {
	if c.ctxCancel != nil {
		c.ctxCancel()
	}
}

func isMatch(t *TimerTask) bool {
	switch t.timerType {
	case 0:
		return !t.end
	case 1:
		return matchMinute(t.date)
	case 2:
		return matchHour(t.date)
	case 3:
		return matchDay(t.date)
	default:
		return false
	}
}

func matchSecond(tm time.Time) bool {
	return time.Since(tm).Seconds() >= 1
}

func matchMinute(tm time.Time) bool {
	return time.Since(tm).Minutes() >= 1
}

func matchHour(tm time.Time) bool {
	return time.Since(tm).Hours() >= 1
}

func matchDay(tm time.Time) bool {
	return time.Since(tm).Hours() >= 24
}
