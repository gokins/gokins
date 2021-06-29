package service

import (
	"encoding/json"
	"errors"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/engine"
	"github.com/gokins-main/gokins/model"
	"time"
)

func Run(pipeId string, repoId string) error {
	tpipe := &model.TPipeline{}
	ok, _ := comm.Db.Where("id=?", pipeId).Get(tpipe)
	if !ok {
		return errors.New("流水线不存在")
	}
	if tpipe.JsonContent == "" {
		return errors.New("流水线Yaml为空")
	}
	trepo := &model.TRepo{}
	ok, _ = comm.Db.Where("id=?", repoId).Get(trepo)
	if !ok {
		return errors.New("仓库不存在")
	}
	pipe := &bean.Pipeline{}
	err := json.Unmarshal([]byte(tpipe.JsonContent), pipe)
	if err != nil {
		return err
	}

	err = pipe.Check()
	if err != nil {
		return err
	}

	pipe.ConvertCmd()
	number := int64(0)
	_, err = comm.Db.
		SQL("SELECT max(number) FROM t_pipeline_version WHERE pipeline_id = ?", pipeId).
		Get(&number)
	if err != nil {
		return err
	}

	tpv := &model.TPipelineVersion{
		Id:                  utils.NewXid(),
		Number:              number + 1,
		Branch:              "0",
		Events:              "run",
		RepoId:              trepo.Id,
		RepoName:            trepo.Name,
		Sha:                 "",
		PipelineName:        pipe.Name,
		PipelineDisplayName: pipe.DisplayName,
		PipelineId:          tpipe.Id,
		Version:             "",
		Content:             tpipe.JsonContent,
		Created:             time.Now(),
		Deleted:             0,
		RepoCloneUrl:        trepo.Url,
	}
	_, err = comm.Db.InsertOne(tpv)
	if err != nil {
		return err
	}

	tb := &model.TBuild{
		Id:                utils.NewXid(),
		PipelineId:        pipeId,
		PipelineVersionId: tpv.Id,
		Status:            common.BuildStatusPending,
		Created:           time.Now(),
		Version:           "",
	}
	_, err = comm.Db.InsertOne(tb)
	if err != nil {
		return err
	}

	rb := &runtime.Build{
		Id:         tb.Id,
		PipelineId: tb.PipelineId,
		Status:     common.BuildStatusPending,
		Created:    time.Now(),
		Repo: &runtime.Repository{
			Name:     "",
			Token:    "",
			Sha:      "",
			CloneURL: trepo.Url,
		},
		Variables: pipe.Variables,
		Stages:    nil,
	}

	rstages := make([]*runtime.Stage, 0)
	for _, stage := range pipe.Stages {
		ts := &model.TStage{
			Id:                utils.NewXid(),
			PipelineVersionId: tpv.Id,
			BuildId:           tb.Id,
			Status:            common.BuildStatusPending,
			Name:              stage.Name,
			DisplayName:       stage.DisplayName,
			Created:           time.Now(),
			Stage:             stage.Stage,
		}
		rt := &runtime.Stage{
			Id:          utils.NewXid(),
			BuildId:     tb.Id,
			Status:      common.BuildStatusPending,
			Name:        stage.Name,
			DisplayName: stage.DisplayName,
			Created:     time.Now(),
			Stage:       stage.Stage,
		}
		_, err = comm.Db.InsertOne(ts)
		if err != nil {
			return err
		}
		rsteps := make([]*runtime.Step, 0)
		for _, step := range stage.Steps {
			cmds, err := json.Marshal(step.Commands)
			if err != nil {
				return err
			}
			djs, err := json.Marshal(step.DependsOn)
			if err != nil {
				return err
			}
			de, err := json.Marshal(step.Environments)
			if err != nil {
				return err
			}
			tsp := &model.TStep{
				Id:                utils.NewXid(),
				BuildId:           tb.Id,
				StageId:           ts.Id,
				DisplayName:       step.DisplayName,
				PipelineVersionId: tpv.Id,
				Step:              step.Step,
				Status:            common.BuildStatusPending,
				Name:              step.Name,
				Created:           time.Now(),
				Commands:          string(cmds),
				DependsOn:         string(djs),
				Environments:      string(de),
			}
			rtp := &runtime.Step{
				Id:           utils.NewXid(),
				BuildId:      tb.Id,
				StageId:      ts.Id,
				DisplayName:  step.DisplayName,
				Step:         step.Step,
				Status:       common.BuildStatusPending,
				Name:         step.Name,
				Commands:     step.Commands,
				DependsOn:    step.DependsOn,
				Environments: step.Environments,
			}
			_, err = comm.Db.InsertOne(tsp)
			if err != nil {
				return err
			}
			rsteps = append(rsteps, rtp)
		}
		rstages = append(rstages, rt)
	}
	rb.Stages = rstages
	engine.Mgr.BuildEgn().Put(rb)
	return nil
}
