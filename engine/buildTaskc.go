package engine

import (
	"errors"
	"fmt"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/runner/runners"
	"github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func (c *BuildTask) check() bool {
	if c.build.Repo == nil {
		c.status(common.BuildEventCheckParam, "repo param err")
		return false
	}
	if c.build.Repo.CloneURL == "" {
		c.status(common.BuildEventCheckParam, "repo param err:clone url")
		return false
	}
	c.repoPath = c.build.Repo.CloneURL
	s, err := os.Stat(c.repoPath)
	if err == nil && s.IsDir() {
		c.isClone = false
		c.repoPaths = c.repoPath
	} else {
		if !common.RegUrl.MatchString(c.build.Repo.CloneURL) {
			c.status(common.BuildEventCheckParam, "repo param err:clone url")
			return false
		}
		c.isClone = true
	}
	if c.build.Stages == nil || len(c.build.Stages) <= 0 {
		c.build.Event = common.BuildEventCheckParam
		c.build.Error = "build Stages is empty"
		return false
	}
	stages := make(map[string]*taskStage)
	for _, v := range c.build.Stages {
		if v.BuildId != c.build.Id {
			c.build.Event = common.BuildEventCheckParam
			c.build.Error = fmt.Sprintf("Stage Build id err:%s/%s", v.BuildId, c.build.Id)
			return false
		}
		if v.Name == "" {
			c.build.Event = common.BuildEventCheckParam
			c.build.Error = "build Stage name is empty"
			return false
		}
		if v.Steps == nil || len(v.Steps) <= 0 {
			c.build.Event = common.BuildEventCheckParam
			c.build.Error = "build Stages is empty"
			return false
		}
		if _, ok := stages[v.Name]; ok {
			c.build.Event = common.BuildEventCheckParam
			c.build.Error = fmt.Sprintf("build Stages.%s is repeat", v.Name)
			return false
		}
		vs := &taskStage{
			stage: v,
			jobs:  make(map[string]*jobSync),
		}
		stages[v.Name] = vs
		for _, e := range v.Steps {
			if e.BuildId != c.build.Id {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = fmt.Sprintf("Job Build id err:%s/%s", v.BuildId, c.build.Id)
				return false
			}
			if e.StageId != v.Id {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = fmt.Sprintf("Job Stage id err:%s/%s", v.BuildId, c.build.Id)
				return false
			}
			e.Step = strings.TrimSpace(e.Step)
			if e.Step == "" {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = "build Step Plugin is empty"
				return false
			}
			if e.Name == "" {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = "build Step name is empty"
				return false
			}
			if _, ok := vs.jobs[e.Name]; ok {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = fmt.Sprintf("build Job.%s is repeat", e.Name)
				return false
			}
			job := &jobSync{
				task:  c,
				step:  e,
				cmdmp: make(map[string]*cmdSync),
			}
			err = c.genRunjob(v, job)
			if err != nil {
				c.build.Event = common.BuildEventCheckParam
				c.build.Error = fmt.Sprintf("build Job.%s Commands err", e.Name)
				return false
			}
			vs.RLock()
			vs.jobs[e.Name] = job
			vs.RUnlock()
			c.joblk.Lock()
			c.jobs[e.Id] = job
			c.joblk.Unlock()
		}
	}
	/*for _,v:=range stages{
		for _,e:=range v.jobs{
			err:=Mgr.jobEgn.Put(e)
			if err!=nil{
				c.build.Event = common.BuildEventPutJob
				c.build.Error=err.Error()
				return false
			}
		}
	}*/

	for k, v := range stages {
		c.stages[k] = v
	}
	return true
}

func (c *BuildTask) genRunjob(stage *runtime.Stage, job *jobSync) (rterr error) {
	defer func() {
		if err := recover(); err != nil {
			rterr = fmt.Errorf("recover:%v", err)
			logrus.Warnf("BuildTask genRunjob recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
		}
	}()
	runjb := &runners.RunJob{
		Id:           job.step.Id,
		StageId:      job.step.StageId,
		BuildId:      job.step.BuildId,
		StageName:    stage.Name,
		Step:         job.step.Step,
		Name:         job.step.Name,
		Input:        job.step.Input,
		Env:          job.step.Env,
		Artifacts:    job.step.Artifacts,
		UseArtifacts: job.step.UseArtifacts,
	}
	var err error
	switch job.step.Commands.(type) {
	case string:
		c.appendcmds(runjb, job.step.Commands.(string))
	case []interface{}:
		err = c.gencmds(runjb, job.step.Commands.([]interface{}))
	case []string:
		var ls []interface{}
		ts := job.step.Commands.([]string)
		for _, v := range ts {
			ls = append(ls, v)
		}
		err = c.gencmds(runjb, ls)
	default:
		err = errors.New("commands format err")
	}
	if err != nil {
		return err
	}
	if len(runjb.Commands) <= 0 {
		return errors.New("command format empty")
	}
	job.runjb = runjb
	for i, v := range runjb.Commands {
		job.cmdmp[v.Id] = &cmdSync{
			cmd:    v,
			status: common.BuildStatusPending,
		}
		cmd := &model.TCmdLine{
			Id: v.Id,
			//GroupId: v.Gid,
			BuildId: job.step.BuildId,
			StepId:  job.step.Id,
			Status:  common.BuildStatusPending,
			Num:     i + 1,
			Content: v.Conts,
			Created: time.Now(),
		}
		vls := common.RegVar.FindAllStringSubmatch(v.Conts, -1)
		for _, zs := range vls {
			k := zs[1]
			if k == "" {
				continue
			}
			vas := ""
			secret := false
			va, ok := c.build.Vars[k]
			if ok {
				vas = va.Value
				secret = va.Secret
			}
			v.Conts = strings.ReplaceAll(v.Conts, zs[0], vas)
			if secret {
				cmd.Content = strings.ReplaceAll(cmd.Content, zs[0], "***")
			} else {
				cmd.Content = strings.ReplaceAll(cmd.Content, zs[0], vas)
			}
		}
		_, err = comm.Db.InsertOne(cmd)
		if err != nil {
			comm.Db.Where("build_id=? and step_id=?", cmd.BuildId, cmd.StepId).Delete(cmd)
			return err
		}
	}
	return nil
}
func (c *BuildTask) appendcmds(runjb *runners.RunJob, conts string) {
	m := &runners.CmdContent{
		Id:    utils.NewXid(),
		Conts: conts,
	}
	logrus.Debugf("append cmd(%d)-%s", len(runjb.Commands), m.Conts)
	//job.Commands[m.Id] = m
	runjb.Commands = append(runjb.Commands, m)
}
func (c *BuildTask) gencmds(runjb *runners.RunJob, cmds []interface{}) (rterr error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("BuildTask gencmds recover:%v", err)
			logrus.Warnf("BuildTask stack:%s", string(debug.Stack()))
			rterr = fmt.Errorf("%v", err)
		}
	}()
	for _, v := range cmds {
		switch v.(type) {
		case string:
			//gid := utils.NewXid()
			//grp:=&hbtpBean.CmdGroupJson{Id: utils.NewXid()}
			c.appendcmds(runjb, v.(string))
		case []interface{}:
			//gid := utils.NewXid()
			for _, v1 := range v.([]interface{}) {
				c.appendcmds(runjb, fmt.Sprintf("%v", v1))
			}
		case map[interface{}]interface{}:
			for _, v1 := range v.(map[interface{}]interface{}) {
				//gid := utils.NewXid()
				switch v1.(type) {
				case string:
					c.appendcmds(runjb, fmt.Sprintf("%v", v1))
				case []interface{}:
					for _, v2 := range v1.([]interface{}) {
						c.appendcmds(runjb, fmt.Sprintf("%v", v2))
					}
				}
			}
		case map[string]interface{}:
			for _, v1 := range v.(map[string]interface{}) {
				//gid := utils.NewXid()
				switch v1.(type) {
				case string:
					c.appendcmds(runjb, fmt.Sprintf("%v", v1))
				case []interface{}:
					for _, v2 := range v1.([]interface{}) {
						c.appendcmds(runjb, fmt.Sprintf("%v", v2))
					}
				}
			}
		}
	}
	return nil
}
