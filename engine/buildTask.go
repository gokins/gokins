package engine

import (
	"context"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/gokins/comm"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

type BuildTask struct {
	bd   *runtime.Build
	ctx  context.Context
	cncl context.CancelFunc
}

func NewBuildTask(bd *runtime.Build) *BuildTask {
	c := &BuildTask{bd: bd}
	c.ctx, c.cncl = context.WithTimeout(comm.Ctx, time.Hour*3)
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

}
