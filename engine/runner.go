package engine

import (
	"errors"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/runner/runners"
)

type shellRunner struct {
}

func (c *shellRunner) PullJob(plugs []string) (*runtime.Step, error) {
	v := Mgr.jobEgn.Pull(plugs)
	if v == nil {
		return nil, errors.New("not found")
	}
	return v, nil
}
func (c *shellRunner) Update(m *runners.UpdateJobInfo) error {
	Mgr.jobEgn.joblk.RLock()
	step, ok := Mgr.jobEgn.jobs[m.Id]
	Mgr.jobEgn.joblk.RUnlock()
	if !ok {
		return errors.New("not found")
	}
	step.Lock()
	step.step.Status = m.Status
	step.step.Error = m.Error
	step.step.ExitCode = m.ExitCode
	step.Unlock()
	return nil
}
