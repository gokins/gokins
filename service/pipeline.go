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

func Run(pipeId string) (*model.TPipelineVersion, error) {
	tpipe := &model.TPipeline{}
	ok, _ := comm.Db.Where("id=? and deleted != 1", pipeId).Get(tpipe)
	if !ok {
		return nil, errors.New("流水线不存在")
	}
	if tpipe.JsonContent == "" {
		return nil, errors.New("流水线Yaml为空")
	}
	pipe := &bean.Pipeline{}
	err := json.Unmarshal([]byte(tpipe.JsonContent), pipe)
	if err != nil {
		return nil, err
	}

	err = pipe.Check()
	if err != nil {
		return nil, err
	}

	pipe.ConvertCmd()
	number := int64(0)
	_, err = comm.Db.
		SQL("SELECT max(number) FROM t_pipeline_version WHERE pipeline_id = ?", pipeId).
		Get(&number)
	if err != nil {
		return nil, err
	}

	tpv := &model.TPipelineVersion{
		Id:                  utils.NewXid(),
		Number:              number + 1,
		Branch:              "",
		Events:              "run",
		Sha:                 "",
		PipelineName:        tpipe.Name,
		PipelineDisplayName: tpipe.DisplayName,
		PipelineId:          tpipe.Id,
		Version:             "",
		Content:             tpipe.JsonContent,
		Created:             time.Now(),
		Deleted:             0,
		RepoCloneUrl:        tpipe.Url,
	}
	_, err = comm.Db.InsertOne(tpv)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	rb := &runtime.Build{
		Id:         tb.Id,
		PipelineId: tb.PipelineId,
		Status:     common.BuildStatusPending,
		Created:    time.Now(),
		Repo: &runtime.Repository{
			Name:     tpipe.Username,
			Token:    tpipe.AccessToken,
			Sha:      "",
			CloneURL: tpipe.Url,
		},
		Variables: pipe.Variables,
	}

	for i, stage := range pipe.Stages {
		ts := &model.TStage{
			Id:                utils.NewXid(),
			PipelineVersionId: tpv.Id,
			BuildId:           tb.Id,
			Status:            common.BuildStatusPending,
			Name:              stage.Name,
			DisplayName:       stage.DisplayName,
			Created:           time.Now(),
			Stage:             stage.Stage,
			Sort:              i,
		}
		rt := &runtime.Stage{
			Id:          ts.Id,
			BuildId:     tb.Id,
			Status:      common.BuildStatusPending,
			Name:        stage.Name,
			DisplayName: stage.DisplayName,
			Created:     time.Now(),
			Stage:       stage.Stage,
		}
		_, err = comm.Db.InsertOne(ts)
		if err != nil {
			return nil, err
		}
		for j, step := range stage.Steps {
			cmds, err := json.Marshal(step.Commands)
			if err != nil {
				return nil, err
			}
			djs, err := json.Marshal(step.DependsOn)
			if err != nil {
				return nil, err
			}
			de, err := json.Marshal(step.Environments)
			if err != nil {
				return nil, err
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
				Sort:              j,
			}
			rtp := &runtime.Step{
				Id:           tsp.Id,
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
				return nil, err
			}
			rt.Steps = append(rt.Steps, rtp)
		}
		rb.Stages = append(rb.Stages, rt)
	}
	engine.Mgr.BuildEgn().Put(rb)
	return tpv, nil
}
