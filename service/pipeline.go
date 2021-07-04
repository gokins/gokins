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
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

func Run(pipeId string, sha string) (*model.TPipelineVersion, error) {
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
	return preBuild(pipe, tpipe, sha)
}

func ReBuild(tvp *model.TPipelineVersion) (*model.TPipelineVersion, error) {
	tpipe := &model.TPipeline{}
	ok, _ := comm.Db.Where("id=? and deleted != 1", tvp.PipelineId).Get(tpipe)
	if !ok {
		return nil, errors.New("流水线不存在")
	}
	if tvp.Content == "" {
		return nil, errors.New("流水线Yaml为空")
	}
	pipe := &bean.Pipeline{}
	err := json.Unmarshal([]byte(tvp.Content), pipe)
	if err != nil {
		return nil, err
	}
	return preBuild(pipe, tpipe, tvp.Sha)
}

func preBuild(pipe *bean.Pipeline, tpipe *model.TPipeline, sha string) (*model.TPipelineVersion, error) {
	err := pipe.Check()
	if err != nil {
		return nil, err
	}
	pipe.ConvertCmd()

	m := convertVar(tpipe.Id, pipe.Vars)
	var tVars = []*model.TPipelineVar{}
	err = comm.Db.Where("pipeline_id = ? ", tpipe.Id).Find(&tVars)
	if err != nil {
		return nil, err
	}
	for _, v := range tVars {
		m[v.Name] = &runtime.Variables{
			Name:   v.Name,
			Value:  v.Value,
			Secret: v.Public == 1,
		}
	}
	replaceStages(pipe.Stages, m)

	number := int64(0)
	_, err = comm.Db.
		SQL("SELECT max(number) FROM t_pipeline_version WHERE pipeline_id = ?", tpipe.Id).
		Get(&number)
	if err != nil {
		return nil, err
	}
	tpv := &model.TPipelineVersion{
		Id:                  utils.NewXid(),
		Number:              number + 1,
		Events:              "run",
		Sha:                 sha,
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
		PipelineId:        tpipe.Id,
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
			Sha:      sha,
			CloneURL: tpipe.Url,
		},
		Vars: m,
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
			djs, err := json.Marshal(step.Waits)
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
				Waits:             string(djs),
				Sort:              j,
			}
			rtp := &runtime.Step{
				Id:          tsp.Id,
				BuildId:     tb.Id,
				StageId:     ts.Id,
				DisplayName: step.DisplayName,
				Step:        step.Step,
				Status:      common.BuildStatusPending,
				Name:        step.Name,
				Commands:    step.Commands,
				Waits:       step.Waits,
				Env:         step.Env,
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

func convertVar(pipelineId string, vm map[string]string) map[string]*runtime.Variables {
	vms := make(map[string]*runtime.Variables, 0)
	for k, v := range vm {
		k1, kok := replaceVariable(pipelineId, common.RegVar, k)
		v1, vok := replaceVariable(pipelineId, common.RegVar, v)
		vms[k1] = &runtime.Variables{
			Name:   k1,
			Value:  v1,
			Secret: kok || vok,
		}
	}
	return vms
}

func replaceVariable(pipelineId string, reg *regexp.Regexp, s string) (string, bool) {
	isSecret := false
	if reg.MatchString(s) {
		all := reg.FindAllStringSubmatch(s, -1)
		for _, v := range all {
			tVars := &model.TPipelineVar{}
			ok, err := comm.Db.Where("pipeline_id = ? and name = ?", pipelineId, v[1]).Get(tVars)
			if err != nil {
				logrus.Debugf("replaceVariable db err %v", err)
			}
			if !ok {
				continue
			}
			if tVars.Public != 1 {
				isSecret = true
			}
			s = strings.ReplaceAll(s, v[0], tVars.Value)
		}
	}
	return s, isSecret
}

func replaceStages(stages []*bean.Stage, mVars map[string]*runtime.Variables) {
	for _, stage := range stages {
		replaceStage(stage, mVars)
	}
}
func replaceStage(stage *bean.Stage, mVars map[string]*runtime.Variables) {
	stage.Stage = replace(stage.Stage, mVars)
	stage.Name = replace(stage.Name, mVars)
	stage.DisplayName = replace(stage.DisplayName, mVars)
	if stage.Steps != nil && len(stage.Steps) > 0 {
		replaceSteps(stage.Steps, mVars)
	}
}
func replaceSteps(steps []*bean.Step, mVars map[string]*runtime.Variables) {
	for _, step := range steps {
		replaceStep(step, mVars)
	}
}
func replaceStep(step *bean.Step, mVars map[string]*runtime.Variables) {
	step.Step = replace(step.Step, mVars)
	step.Name = replace(step.Name, mVars)
	step.DisplayName = replace(step.DisplayName, mVars)
	step.Image = replace(step.Image, mVars)
	if step.Env != nil && len(step.Env) > 0 {
		step.Env = replaceEnvs(step.Env, mVars)
	}
}

func replaceEnvs(envs map[string]string, mVars map[string]*runtime.Variables) map[string]string {
	m := map[string]string{}
	for k, v := range envs {
		m[replace(k, mVars)] = replace(v, mVars)
	}
	return m
}

func replace(s string, mVars map[string]*runtime.Variables) string {
	if common.RegVar.MatchString(s) {
		all := common.RegVar.FindAllStringSubmatch(s, -1)
		for _, v2 := range all {
			rVar, ok := mVars[v2[1]]
			if !ok {
				continue
			}
			if rVar.Secret {
				s = strings.ReplaceAll(s, v2[0], "***")
			}
			s = strings.ReplaceAll(s, v2[0], rVar.Value)
		}
	}
	return s
}
