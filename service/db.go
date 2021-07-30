package service

import "github.com/gokins/gokins/comm"

func GetIdOrAid(id interface{}, e interface{}) bool {
	if id == nil || e == nil {
		return false
	}
	switch id.(type) {
	case string:
		ids := id.(string)
		if ids == "" {
			return false
		}
	}
	ok, _ := comm.Db.Where("id=?", id).Get(e)
	if !ok {
		ok, _ = comm.Db.Where("aid=?", id).Get(e)
	}
	return ok
}
