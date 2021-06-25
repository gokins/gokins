package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/util"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetUser(uid int64) (*model.TUser, bool) {
	e := &model.TUser{}
	ok, err := comm.Db.Where("id=?", uid).Get(e)
	if err != nil {
		logrus.Errorf("GetUser(%s) err:%v", uid, err)
	}
	return e, ok
}
func FindUserName(name string) (*model.TUser, bool) {
	e := &model.TUser{}
	ok, err := comm.Db.Where("name=?", name).Get(e)
	if err != nil {
		logrus.Errorf("FindUserName(%s) err:%v", name, err)
	}
	return e, ok
}

func GetUserCache(uid int64) (*model.TUser, bool) {
	var ok bool
	e := &model.TUser{}
	uids := fmt.Sprintf("user:%d", uid)
	err := comm.CacheGets(uids, e)
	if err == nil {
		return e, true
	}
	e, ok = GetUser(uid)
	if ok {
		comm.CacheSets(uids, e)
	}
	return e, ok
}
func CurrUserCache(c *gin.Context) (*model.TUser, bool) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	tk := util.GetToken(c, comm.Cfg.Server.LoginKey)
	if tk == nil {
		return nil, false
	}
	uid, ok := tk["uid"]
	if !ok {
		return nil, false
	}
	uids, ok := uid.(string)
	if !ok {
		return nil, false
	}
	uidi, err := strconv.ParseInt(uids, 10, 64)
	if err != nil {
		return nil, false
	}
	return GetUserCache(uidi)
}
