package bean

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Pipeline struct {
	Version  string              `yaml:"version,omitempty" json:"version"`
	Triggers map[string]*Trigger `yaml:"triggers,omitempty" json:"triggers"`
	Vars     map[string]string   `yaml:"vars,omitempty" json:"vars"`
	Stages   []*Stage            `yaml:"stages,omitempty" json:"stages"`
}

type Trigger struct {
	AutoCancel     bool       `yaml:"autoCancel,omitempty" json:"autoCancel,omitempty"`
	Timeout        string     `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Branches       *Condition `yaml:"branches,omitempty" json:"branches,omitempty"`
	Tags           *Condition `yaml:"tags,omitempty" json:"tags,omitempty"`
	Paths          *Condition `yaml:"paths,omitempty" json:"paths,omitempty"`
	Notes          *Condition `yaml:"notes,omitempty" json:"notes,omitempty"`
	CommitMessages *Condition `yaml:"commitMessages,omitempty" json:"commitMessages,omitempty"`
}

type Condition struct {
	Include []string `yaml:"include,omitempty" json:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

type Stage struct {
	Stage       string  `yaml:"stage" json:"stage"`
	Name        string  `yaml:"name,omitempty" json:"name"`
	DisplayName string  `yaml:"displayName,omitempty" json:"displayName"`
	Steps       []*Step `yaml:"steps,omitempty" json:"steps"`
}

/*type Input struct {
	Value string `yaml:"value"`
	Required bool `yaml:"required"`
}*/
type Step struct {
	Step         string            `yaml:"step" json:"step"`
	DisplayName  string            `yaml:"displayName,omitempty" json:"displayName"`
	Name         string            `yaml:"name,omitempty" json:"name"`
	Input        map[string]string `yaml:"input,omitempty" json:"input"`
	Env          map[string]string `yaml:"env,omitempty" json:"env"`
	Commands     interface{}       `yaml:"commands,omitempty" json:"commands"`
	Waits        []string          `yaml:"wait,omitempty" json:"wait"`
	Image        string            `yaml:"image,omitempty" json:"image"`
	Artifacts    []*Artifact       `yaml:"artifacts,omitempty" json:"artifacts"`
	UseArtifacts []*UseArtifacts   `yaml:"useArtifacts,omitempty" json:"useArtifacts"`
}

type Artifact struct {
	Scope      string `yaml:"scope,omitempty" json:"scope"`
	Repository string `yaml:"repository,omitempty" json:"repository"`
	Name       string `yaml:"name,omitempty" json:"name"`
	Path       string `yaml:"path,omitempty" json:"path"`
}

type UseArtifacts struct {
	Scope      string `yaml:"scope" json:"scope"`           //archive,pipeline,env
	Repository string `yaml:"repository" json:"repository"` // archive,制品库ID
	Name       string `yaml:"name" json:"name"`             //archive,pipeline,env
	//IsForce    bool   `yaml:"isForce" json:"isForce"`
	IsUrl bool   `yaml:"isUrl" json:"isUrl"`
	Alias string `yaml:"alias" json:"alias"`
	Path  string `yaml:"path" json:"path"` //archive,pipeline

	FromStage string `yaml:"fromStage" json:"sourceStage"` //pipeline
	FromStep  string `yaml:"fromStep" json:"sourceStep"`   //pipeline
}

func (c *Pipeline) ToJson() ([]byte, error) {
	c.ConvertCmd()
	return json.Marshal(c)
}
func (c *Pipeline) ConvertCmd() {
	for _, stage := range c.Stages {
		for _, step := range stage.Steps {
			v := step.Commands
			switch v.(type) {
			case string:
				step.Commands = v.(string)
			case []interface{}:
				ls := make([]string, 0)
				for _, v1 := range v.([]interface{}) {
					ls = append(ls, fmt.Sprintf("%v", v1))
				}
				step.Commands = ls
			default:
				step.Commands = fmt.Sprintf("%v", v)
			}
		}
	}
}

func (c *Pipeline) Check() error {
	stages := make(map[string]map[string]*Step)
	if c.Stages == nil || len(c.Stages) <= 0 {
		return errors.New("stages 为空")
	}
	for _, v := range c.Stages {
		if v.Name == "" {
			return errors.New("stages name 为空")
		}
		if v.Steps == nil || len(v.Steps) <= 0 {
			return errors.New("step 为空")
		}
		if _, ok := stages[v.Name]; ok {
			return errors.New(fmt.Sprintf("build stages.%s 重复", v.Name))
		}
		m := map[string]*Step{}
		stages[v.Name] = m
		for _, e := range v.Steps {
			if strings.TrimSpace(e.Step) == "" {
				return errors.New("step 插件为空")
			}
			if e.Name == "" {
				return errors.New("step name 为空")
			}
			if _, ok := m[e.Name]; ok {
				return errors.New(fmt.Sprintf("steps.%s 重复", e.Name))
			}
			m[e.Name] = e
		}
	}
	return nil
}

//func (c *Pipeline) SkipTriggerRules(events string) bool {
//	if events != "manual" {
//		return true
//	}
//
//	if c.Triggers == nil || len(c.Triggers) <= 0 {
//		logrus.Error("Triggers is empty")
//		return false
//	}
//	switch events {
//	case "push", "pr", "comment":
//	default:
//		logrus.Debugf("not match action:%v", events)
//		return false
//	}
//	v, ok := c.Triggers[events]
//	if !ok {
//		logrus.Debugf("not match action: %v", events)
//		return false
//	}
//	if v == nil {
//		logrus.Debugf("%v trigger is empty",events)
//		return false
//	}
//	if !skipCommitNotes(v.Notes, pb.Info.Note) {
//		return false
//	} else if !skipBranch(v.Branches, pb.Info.Repository.Branch) {
//		return false
//	} else if !skipCommitMessages(v.CommitMessages, pb.Info.CommitMessage) {
//		return false
//	} else {
//		logrus.Debugf("%v skip", c.Name)
//		return true
//	}
//}
