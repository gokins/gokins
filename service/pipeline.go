package service

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"gopkg.in/yaml.v3"
)

func Run(uid, pipeId, sha, event string) (*model.TPipelineVersion, *runtime.Build, error) {
	tpipe := &model.TPipelineConf{}
	ok, _ := comm.Db.Where("pipeline_id=?", pipeId).Get(tpipe)
	if !ok {
		return nil, nil, errors.New("流水线不存在")
	}
	if tpipe.YmlContent == "" {
		return nil, nil, errors.New("流水线Yaml为空")
	}
	pipe := &bean.Pipeline{}
	err := yaml.Unmarshal([]byte(tpipe.YmlContent), pipe)
	if err != nil {
		return nil, nil, err
	}
	return preBuild(uid, pipe, tpipe, sha, event)
}

func ReBuild(uid string, tvp *model.TPipelineVersion) (*model.TPipelineVersion, *runtime.Build, error) {
	tpipe := &model.TPipelineConf{}
	ok, _ := comm.Db.Where("pipeline_id=?", tvp.PipelineId).Get(tpipe)
	if !ok {
		return nil, nil, errors.New("流水线不存在")
	}
	if tvp.Content == "" {
		return nil, nil, errors.New("流水线Yaml为空")
	}
	pipe := &bean.Pipeline{}
	err := yaml.Unmarshal([]byte(tvp.Content), pipe)
	if err != nil {
		return nil, nil, err
	}
	return preBuild(uid, pipe, tpipe, "", "rebuild", tvp)
}

func preBuild(uid string, pipe *bean.Pipeline, tpipe *model.TPipelineConf, sha, event string,
	tvp ...*model.TPipelineVersion) (*model.TPipelineVersion, *runtime.Build, error) {
	tp := &model.TPipeline{}
	ok, _ := comm.Db.Where("id=? and deleted != 1", tpipe.PipelineId).Get(tp)
	if !ok {
		return nil, nil, errors.New("流水线不存在")
	}
	err := pipe.Check()
	if err != nil {
		return nil, nil, err
	}
	pipe.ConvertCmd()

	m, err := convertVar(tpipe.PipelineId, pipe.Vars)
	if err != nil {
		return nil, nil, err
	}
	replaceStages(pipe.Stages, m)

	number := int64(0)
	_, err = comm.Db.
		SQL("SELECT max(number) FROM t_pipeline_version WHERE pipeline_id = ?", tpipe.PipelineId).
		Get(&number)
	if err != nil {
		return nil, nil, err
	}
	tpv := &model.TPipelineVersion{
		Id:                  utils.NewXid(),
		Uid:                 uid,
		Number:              number + 1,
		Events:              event,
		Sha:                 sha,
		PipelineName:        tp.Name,
		PipelineDisplayName: tp.DisplayName,
		PipelineId:          tpipe.PipelineId,
		Version:             "",
		Content:             tpipe.YmlContent,
		Created:             time.Now(),
		Deleted:             0,
		RepoCloneUrl:        tpipe.Url,
	}
	if len(tvp) > 0 && tvp[0] != nil {
		tpv.Sha = tvp[0].Sha
		tpv.Content = tvp[0].Content
	}
	_, err = comm.Db.InsertOne(tpv)
	if err != nil {
		return nil, nil, err
	}

	tb := &model.TBuild{
		Id:                utils.NewXid(),
		PipelineId:        tpipe.PipelineId,
		PipelineVersionId: tpv.Id,
		Status:            common.BuildStatusPending,
		Created:           time.Now(),
		Version:           "",
	}
	_, err = comm.Db.InsertOne(tb)
	if err != nil {
		return nil, nil, err
	}

	rb := &runtime.Build{
		Id:                tb.Id,
		PipelineId:        tb.PipelineId,
		PipelineVersionId: tb.PipelineVersionId,
		Status:            common.BuildStatusPending,
		Created:           time.Now(),
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
			return nil, nil, err
		}
		for j, step := range stage.Steps {
			cmds, err := json.Marshal(step.Commands)
			if err != nil {
				return nil, nil, err
			}
			djs, err := json.Marshal(step.Waits)
			if err != nil {
				return nil, nil, err
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
				Input:       step.Input,
			}
			for _, v := range step.Artifacts {
				rtp.Artifacts = append(rtp.Artifacts, &runtime.Artifact{
					Scope:      v.Scope,
					Repository: v.Repository,
					Name:       v.Name,
					Path:       v.Path,
				})
			}
			for _, v := range step.UseArtifacts {
				rtp.UseArtifacts = append(rtp.UseArtifacts, &runtime.UseArtifact{
					Scope:      v.Scope,
					Repository: v.Repository,
					Name:       v.Name,
					Path:       v.Path,
					//IsForce:     v.IsForce,
					IsUrl:       v.IsUrl,
					Alias:       v.Alias,
					SourceStage: v.FromStage,
					SourceStep:  v.FromStep,
				})
			}
			_, err = comm.Db.InsertOne(tsp)
			if err != nil {
				return nil, nil, err
			}
			rt.Steps = append(rt.Steps, rtp)
		}
		rb.Stages = append(rb.Stages, rt)
	}
	return tpv, rb, nil
}

func convertVar(pipelineId string, vm map[string]string) (map[string]*runtime.Variables, error) {
	var tVars []*model.TPipelineVar
	err := comm.Db.Where("pipeline_id = ? ", pipelineId).Find(&tVars)
	if err != nil {
		return nil, err
	}
	vms := make(map[string]*runtime.Variables, 0)
	for _, v := range tVars {
		vms[v.Name] = &runtime.Variables{
			Name:   v.Name,
			Value:  v.Value,
			Secret: v.Public == 1,
		}
	}
	for k, v := range vm {
		s, b := replace(v, vms, true)
		vms[k] = &runtime.Variables{
			Name:   k,
			Value:  s,
			Secret: b,
		}
	}

	for k, v := range vms {
		s, _ := replace(v.Value, vms, true)
		vms[k].Value = s
	}
	return vms, nil
}

func replaceStages(stages []*bean.Stage, mVars map[string]*runtime.Variables) {
	for _, stage := range stages {
		replaceStage(stage, mVars)
	}
}
func replaceStage(stage *bean.Stage, mVars map[string]*runtime.Variables) {
	s, _ := replace(stage.Stage, mVars)
	stage.Stage = s
	s, _ = replace(stage.Name, mVars)
	stage.Name = s
	s, _ = replace(stage.DisplayName, mVars)
	stage.DisplayName = s
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
	s, _ := replace(step.Step, mVars)
	step.Step = s
	s, _ = replace(step.Name, mVars)
	step.Name = s
	s, _ = replace(step.DisplayName, mVars)
	step.DisplayName = s
	s, _ = replace(step.Image, mVars)
	step.Image = s
	if step.Env != nil && len(step.Env) > 0 {
		step.Env = replaceMaps(step.Env, mVars)
	}
	if step.Input != nil && len(step.Input) > 0 {
		step.Input = replaceMaps(step.Input, mVars)
	}
}

func replaceMaps(envs map[string]string, mVars map[string]*runtime.Variables) map[string]string {
	m := map[string]string{}
	for k, v := range envs {
		s, _ := replace(v, mVars, true)
		m[k] = s
	}
	return m
}

func replace(s string, mVars map[string]*runtime.Variables, mustShow ...bool) (string, bool) {
	if s == "" {
		return "", false
	}
	conts := s
	ms := false
	if len(mustShow) == 1 && mustShow[0] {
		ms = true
	}
	secret := false
	if common.RegVar.MatchString(s) {
		all := common.RegVar.FindAllStringSubmatch(s, -1)
		for _, v2 := range all {
			rVar, ok := mVars[v2[1]]
			va := ""
			st := false
			if ok {
				st = rVar.Secret
				va = rVar.Value
				if !secret && rVar.Secret {
					secret = true
				}
			}
			if !ms && st {
				conts = strings.ReplaceAll(conts, v2[0], "***")
			} else {
				conts = strings.ReplaceAll(conts, v2[0], va)
			}
		}
	}
	return conts, secret
}
