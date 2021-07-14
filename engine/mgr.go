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
}

func Start() error {
	Mgr.buildEgn = StartBuildEngine()
	Mgr.jobEgn = StartJobEngine()

	Mgr.brun = &baseRunner{}
	Mgr.hrun = &HbtpRunner{}
	//runners
	comm.Cfg.Server.Shells = append(comm.Cfg.Server.Shells, "shell@ssh")
	if len(comm.Cfg.Server.Shells) > 0 {
		runr := runners.NewEngine(runners.Config{
			Workspace: filepath.Join(comm.WorkPath, common.PathRunner),
			Plugin:    comm.Cfg.Server.Shells,
		}, Mgr.brun)
		err := runr.Start(comm.Ctx)
		if err != nil {
			return err
		}
		Mgr.shellRun = runr
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
