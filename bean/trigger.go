package bean

import "errors"

type TriggerParam struct {
	Id         string `json:"id"`
	PipelineId string `json:"pipelineId"`
	Types      string `json:"types"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Params     string `json:"params"`
	Enabled    bool   ` json:"enabled"`
}

func (c *TriggerParam) Check() error {
	if c.PipelineId == "" {
		return errors.New("流水线ID不能为空")
	}
	if c.Types == "" {
		return errors.New("触发器类型不能为空")
	}
	if c.Name == "" {
		return errors.New("触发器名称不能为空")
	}
	if c.Params == "" {
		return errors.New("触发器参数不能为空")
	}
	return nil
}
