package engine

import (
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"time"
)

func (c *BuildTask) updateBuild(build *runtime.Build) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask updateBuild recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

	e := &model.TBuild{
		Status:   build.Status,
		Error:    build.Error,
		Event:    build.Event,
		Started:  build.Started,
		Finished: build.Finished,
		Updated:  time.Now(),
	}
	_, err := comm.Db.Cols("status", "event", "error", "started", "finished", "updated").
		Where("id=?", build.Id).Update(e)
	if err != nil {
		logrus.Errorf("BuildTask.updateBuild db err:%v", err)
	}

	if !common.BuildStatusEnded(e.Status) {
		return
	}
	stge := &model.TStage{
		Status:   common.BuildStatusCancel,
		Finished: time.Now(),
		Updated:  time.Now(),
	}
	_, err = comm.Db.Cols("status", "finished", "updated").
		Where("build_id=? and `status`!=? and `status`!=? and `status`!=?",
			build.Id, common.BuildStatusOk, common.BuildStatusError, common.BuildStatusCancel).Update(stge)
	if err != nil {
		logrus.Errorf("BuildTask.updateBuild stage err:%v", err)
	}
	stpe := &model.TStep{
		Status:   common.BuildStatusCancel,
		Finished: time.Now(),
		Updated:  time.Now(),
	}
	_, err = comm.Db.Cols("status", "finished", "updated").
		Where("build_id=? and `status`!=? and `status`!=? and `status`!=?",
			build.Id, common.BuildStatusOk, common.BuildStatusError, common.BuildStatusCancel).Update(stpe)
	if err != nil {
		logrus.Errorf("BuildTask.updateBuild step err:%v", err)
	}
	cmde := &model.TCmdLine{
		Status:   common.BuildStatusCancel,
		Finished: time.Now(),
	}
	_, err = comm.Db.Cols("status", "finished").
		Where("build_id=? and `status`!=? and `status`!=? and `status`!=?",
			build.Id, common.BuildStatusOk, common.BuildStatusError, common.BuildStatusCancel).Update(cmde)
	if err != nil {
		logrus.Errorf("BuildTask.updateStage step err:%v", err)
	}
}
func (c *BuildTask) updateStage(stage *runtime.Stage) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask updateBuild recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

	e := &model.TStage{
		Status:   stage.Status,
		Error:    stage.Error,
		Started:  stage.Started,
		Finished: stage.Finished,
		Updated:  time.Now(),
	}
	_, err := comm.Db.Cols("status", "error", "started", "finished", "updated").
		Where("id=?", stage.Id).Update(e)
	if err != nil {
		logrus.Errorf("BuildTask.updateStage db err:%v", err)
	}

	if !common.BuildStatusEnded(e.Status) {
		return
	}
	stpe := &model.TStep{
		Status:   common.BuildStatusCancel,
		Finished: time.Now(),
		Updated:  time.Now(),
	}
	_, err = comm.Db.Cols("status", "finished", "updated").
		Where("stage_id=? and `status`!=? and `status`!=? and `status`!=?",
			stage.Id, common.BuildStatusOk, common.BuildStatusError, common.BuildStatusCancel).Update(stpe)
	if err != nil {
		logrus.Errorf("BuildTask.updateStage step err:%v", err)
	}
}
func (c *BuildTask) updateStep(job *jobSync) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask updateBuild recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()

	job.RLock()
	defer job.RUnlock()
	e := &model.TStep{
		Status:   job.step.Status,
		Event:    job.step.Event,
		Error:    job.step.Error,
		ExitCode: job.step.ExitCode,
		Started:  job.step.Started,
		Finished: job.step.Finished,
		Updated:  time.Now(),
	}
	_, err := comm.Db.Cols("status", "event", "error", "exit_code", "started", "finished", "updated").
		Where("id=?", job.step.Id).Update(e)
	if err != nil {
		logrus.Errorf("BuildTask.updateStep db err:%v", err)
	}

	if !common.BuildStatusEnded(e.Status) {
		return
	}
	cmde := &model.TCmdLine{
		Status:   common.BuildStatusCancel,
		Finished: time.Now(),
	}
	_, err = comm.Db.Cols("status", "finished").
		Where("step_id=? and `status`!=? and `status`!=? and `status`!=?",
			job.step.Id, common.BuildStatusOk, common.BuildStatusError, common.BuildStatusCancel).Update(cmde)
	if err != nil {
		logrus.Errorf("BuildTask.updateStage step err:%v", err)
	}
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
		Code:   cmd.code,
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
	_, err := comm.Db.Cols(cols...).Where("id=?", cmd.cmd.Id).Update(cmde)
	if err != nil {
		logrus.Errorf("BuildTask.updateStep db err:%v", err)
	}
}
