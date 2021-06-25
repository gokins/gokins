package engine

import (
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/runner/runners"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

var Mgr = &Manager{}

type Manager struct {
	buildEgn *BuildEngine
	jobEgn   *JobEngine
	shellRun *runners.Engine
}

func Start() error {
	Mgr.buildEgn = StartBuildEngine()
	Mgr.jobEgn = StartJobEngine()
	if len(comm.Cfg.Server.Shells) > 0 {
		runr := runners.NewEngine(runners.Config{
			Workspace: comm.WorkPath,
			Plugin:    comm.Cfg.Server.Shells,
		}, &baseRunner{})
		err := runr.Start(comm.Ctx)
		if err != nil {
			return err
		}
		Mgr.shellRun = runr
	}
	go func() {
		for !hbtp.EndContext(comm.Ctx) {
			//Mgr.run()
			time.Sleep(time.Millisecond * 100)
		}
		Mgr.buildEgn.Stop()
		if Mgr.shellRun != nil {
			Mgr.shellRun.Stop()
		}
	}()
	return nil
}
func (c *Manager) run() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("Manager run recover:%v", err)
			logrus.Warnf("Manager stack:%s", string(debug.Stack()))
		}
	}()

}

func (c *Manager) BuildEgn() *BuildEngine {
	return c.buildEgn
}
