package engine

import (
	"github.com/gokins-main/core/runtime"
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
