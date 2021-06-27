package service

import "github.com/gokins-main/gokins/comm"

func GetOrg(id interface{}, org interface{}) bool {
	if id == nil || org == nil {
		return false
	}
	ok, _ := comm.Db.Where("id=?", id).Get(org)
	if !ok {
		ok, _ = comm.Db.Where("aid=?", id).Get(org)
	}
	return ok
}
