package service

import (
	"encoding/json"
	"errors"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/model"
	"github.com/sirupsen/logrus"
	"time"
)

func FindParam(key string) (*model.TParam, bool) {
	e := &model.TParam{}
	ok, err := comm.Db.Where("name=?", key).Get(e)
	if err != nil {
		logrus.Errorf("FindParam(%s) err:%v", key, err)
	}
	return e, ok
}
func SetParam(key string, data []byte, tit ...string) error {
	var err error
	db := comm.Db
	e, ok := FindParam(key)
	if len(tit) > 0 {
		e.Title = tit[0]
	}
	e.Data = string(data)
	if ok && e.Aid > 0 {
		_, err = db.Cols("title", "data").Where("aid=?", e.Aid).Update(e)
	} else {
		e.Name = key
		e.Times = time.Now()
		_, err = db.Insert(e)
	}
	return err
}
func SetsParam(key string, data interface{}, tit ...string) error {
	if data == nil {
		return errors.New("data is nil")
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return SetParam(key, bts, tit...)
}

func GetParam(key string) ([]byte, error) {
	e, ok := FindParam(key)
	if ok {
		return []byte(e.Data), nil
	}
	return nil, errors.New("not found param")
}
func GetsParam(key string, data interface{}) error {
	if data == nil {
		return errors.New("data is nil")
	}
	bts, err := GetParam(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bts, data)
}

func GetsParamCache(key string, data interface{}, outm ...time.Duration) error {
	err := comm.CacheGets(key, data)
	if err == nil {
		return nil
	}
	err = GetsParam(key, data)
	if err == nil {
		errs := comm.CacheSets(key, data, outm...)
		if errs != nil {
			logrus.Errorf("GetsParamCache.CacheSets(%s) err:%v", key, errs)
		}
	}
	return err
}
