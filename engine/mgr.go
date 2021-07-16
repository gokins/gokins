package engine

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/gokins-main/core/common"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/runner/runners"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
)

var Mgr = &Manager{}

type Manager struct {
	buildEgn *BuildEngine
	jobEgn   *JobEngine
	shellRun *runners.Engine
	brun     *baseRunner
	hrun     *HbtpRunner
	timerEgn *TimerEngine
}

func Start() error {
	Mgr.buildEgn = StartBuildEngine()
	Mgr.jobEgn = StartJobEngine()
	Mgr.timerEgn = StartTimerEngine()

	Mgr.brun = &baseRunner{}
	Mgr.hrun = &HbtpRunner{}
	//runners
	comm.Cfg.Server.Shells = append(comm.Cfg.Server.Shells, "shell@ssh")
	if len(comm.Cfg.Server.Shells) > 0 {
		Mgr.shellRun = runners.NewEngine(runners.Config{
			Name:      "mainRunner",
			Workspace: filepath.Join(comm.WorkPath, common.PathRunner),
			Plugin:    comm.Cfg.Server.Shells,
		}, Mgr.brun)
		go func() {
			err := Mgr.shellRun.Run(comm.Ctx)
			if err != nil {
				logrus.Errorf("runner err:%v", err)
			}
		}()
	}

	go func() {
		os.RemoveAll(filepath.Join(comm.WorkPath, common.PathTmp))
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
func (c *Manager) HRun() *HbtpRunner {
	return c.hrun
}

func (c *Manager) TimerEng() *TimerEngine {
	return c.timerEgn
}
