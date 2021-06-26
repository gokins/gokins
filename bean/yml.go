package bean

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Pipeline struct {
	Name        string              `yaml:"name,omitempty" json:"name"`
	Version     string              `yaml:"version,omitempty" json:"version"`
	DisplayName string              `yaml:"displayName,omitempty" json:"displayName"`
	Triggers    map[string]*Trigger `yaml:"triggers,omitempty" json:"triggers"`
	Variables   map[string]string   `yaml:"variables,omitempty" json:"variables"`
	Stages      []*Stage            `yaml:"stages,omitempty" json:"stages"`
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

type Step struct {
	Step            string            `yaml:"step" json:"step"`
	DisplayName     string            `yaml:"displayName,omitempty" json:"displayName"`
	Name            string            `yaml:"name,omitempty" json:"name"`
	Environments    map[string]string `yaml:"environments,omitempty" json:"environments"`
	Commands        interface{}       `yaml:"commands,omitempty" json:"commands"`
	DependsOn       []string          `yaml:"dependsOn,omitempty" json:"dependsOn"`
	Image           string            `yaml:"image,omitempty" json:"image"`
	Artifacts       []*Artifact       `yaml:"artifacts,omitempty" json:"artifacts"`
	DependArtifacts []*DependArtifact `yaml:"dependArtifacts,omitempty" json:"dependArtifacts"`
}

type Artifact struct {
	Name       string `yaml:"name,omitempty" json:"name"`
	Scope      string `yaml:"scope,omitempty" json:"scope"`
	Path       string `yaml:"path,omitempty" json:"path"`
	Repository string `yaml:"repository,omitempty" json:"repository"`
	Value      string `yaml:"value,omitempty" json:"value"`
}

type DependArtifact struct {
	BuildName string `yaml:"buildName,omitempty" json:"buildName"`
	StageName string `yaml:"stageName,omitempty" json:"stageName"`
	StepName  string `yaml:"stepName,omitempty" json:"stepName"`

	Type       string `yaml:"type,omitempty" json:"type"`
	Repository string `yaml:"repository,omitempty" json:"repository"`
	Name       string `yaml:"name,omitempty" json:"name"`
	Target     string `yaml:"target,omitempty" json:"target"`
	IsForce    bool   `yaml:"isForce,omitempty" json:"isForce"`

	SourceStage string `yaml:"sourceStage,omitempty" json:"sourceStage"`
	SourceStep  string `yaml:"sourceStep,omitempty" json:"sourceStep"`
}

func (c *Pipeline) ToJson() ([]byte, error) {
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
	return json.Marshal(c)
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
