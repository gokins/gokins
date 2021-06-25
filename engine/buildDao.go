package engine

import (
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (c *BuildTask) updateBuild() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask updateBuild recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

}
func (c *BuildTask) updateStage(stage *runtime.Stage) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask updateBuild recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

}
func (c *BuildTask) updateStep(step *runtime.Step) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask updateBuild recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

}
func (c *BuildTask) updateStepCmd(cmd *cmdSync) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask updateBuild recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

	cmd.RLock()
	defer cmd.RUnlock()
	cmde := &model.TCmdLine{
		Status: cmd.status,
	}
	cols := []string{"status"}
	switch cmd.status {
	case common.BuildStatusRunning:
		cmde.Started = cmd.started
		cols = append(cols, "started")
	default:
		cmde.Finished = cmd.finished
		cols = append(cols, "finished")
	}
	comm.Db.Cols(cols...).Where("id=?", cmd.cmd.Id).Update(cmde)
}
