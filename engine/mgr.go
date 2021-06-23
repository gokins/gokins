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
	Mgr.shellRun = runners.NewEngine(runners.Config{
		Workspace: comm.WorkPath,
		Plugin:    []string{"shell@sh"},
	}, &shellRunner{})
	err := Mgr.shellRun.Start(comm.Ctx)
	if err != nil {
		return err
	}
	go func() {
		for !hbtp.EndContext(comm.Ctx) {
			Mgr.run()
			time.Sleep(time.Second)
		}
		Mgr.buildEgn.Stop()
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
