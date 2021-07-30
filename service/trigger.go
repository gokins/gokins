package service

import (
	"errors"
	"github.com/gokins/gokins/model"
)

func TriggerPerm(tt *model.TTrigger) error {
	lgus := &model.TUser{
		Id: tt.Uid,
	}
	perm := NewPipePerm(lgus, tt.PipelineId)
	if perm.Pipeline() == nil {
		return errors.New("流水线不存在")
	}
	if !IsAdmin(lgus) && !perm.CanWrite() {
		return errors.New("触发器创建者没有权限")
	}
	return nil
}
