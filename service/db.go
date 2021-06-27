package service

import "github.com/gokins-main/gokins/comm"

func GetIdOrAid(id interface{}, e interface{}) bool {
	if id == nil || e == nil {
		return false
	}
	ok, _ := comm.Db.Where("id=?", id).Get(e)
	if !ok {
		ok, _ = comm.Db.Where("aid=?", id).Get(e)
	}
	return ok
}
