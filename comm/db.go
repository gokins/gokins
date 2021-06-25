package comm

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"time"
)

var mainCacheBucket = []byte("mainCacheBucket")

func CacheSet(key string, data []byte, outm ...time.Duration) error {
	if BCache == nil {
		return errors.New("cache not init")
	}
	err := BCache.Update(func(tx *bolt.Tx) error {
		var err error
		bk := tx.Bucket(mainCacheBucket)
		if bk == nil {
			bk, err = tx.CreateBucket(mainCacheBucket)
			if err != nil {
				return err
			}
		}
		if data == nil {
			return bk.Delete([]byte(key))
		}
		buf := &bytes.Buffer{}
		var outms []byte
		if len(outm) > 0 {
			outms = []byte(time.Now().Add(outm[0]).Format(time.RFC3339Nano))
		} else {
			outms = []byte(time.Now().Add(time.Hour).Format(time.RFC3339Nano))
		}
		buf.Write(hbtp.BigIntToByte(int64(len(outms)), 4))
		buf.Write(outms)
		buf.Write(data)
		return bk.Put([]byte(key), buf.Bytes())
	})
	return err
}
func CacheSets(key string, data interface{}, outm ...time.Duration) error {
	if BCache == nil {
		return errors.New("cache not init")
	}
	if data == nil {
		return CacheSet(key, nil)
	}
	bts, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return CacheSet(key, bts, outm...)
}
func parseCacheData(bts []byte) []byte {
	if bts == nil {
		return nil
	}
	ln := int(hbtp.BigByteToInt(bts[:4]))
	tms := string(bts[4 : ln+4])
	outm, err := time.Parse(time.RFC3339Nano, tms)
	if err != nil {
		return nil
	}
	if time.Since(outm).Milliseconds() < 0 {
		return bts[4+ln:]
	}
	return nil
}

var KeyNotFoundErr = errors.New("key not found")
var KeyOutTimeErr = errors.New("key is timeout")

func CacheGet(key string) ([]byte, error) {
	if BCache == nil {
		return nil, errors.New("cache not init")
	}
	var rt []byte
	err := BCache.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket(mainCacheBucket)
		if bk == nil {
			return KeyNotFoundErr
		}
		bts := bk.Get([]byte(key))
		if bts == nil {
			return KeyNotFoundErr
		}
		rt = parseCacheData(bts)
		if rt == nil {
			bk.Delete([]byte(key))
			return KeyOutTimeErr
		}
		return nil
	})
	if time.Since(mainCacheClearTime).Hours() > 30 {
		go mainCacheClear()
	}
	return rt, err
}
func CacheGets(key string, data interface{}) error {
	if BCache == nil {
		return errors.New("cache not init")
	}
	if data == nil {
		return errors.New("data not be nil")
	}
	bts, err := CacheGet(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(bts, data)
}

func CacheFlush() error {
	if BCache == nil {
		return errors.New("cache not init")
	}
	err := BCache.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(mainCacheBucket)
	})
	return err
}

var mainCacheClearTime time.Time

func mainCacheClear() {
	defer func() {
		if err := recover(); err != nil {
			logrus.Errorf("mainCacheClear recover err:%v", err)
		}
	}()

	if BCache == nil {
		return
	}
	/*if time.Now().Hour()!=3|| time.Since(mainCacheClearTime).Hours() < 30 {
		return
	}*/
	mainCacheClearTime = time.Now()
	/*if err := CacheFlush(); err != nil {
		logrus.Errorf("mainCacheClear err:%v", err)
	}*/
	err := BCache.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket(mainCacheBucket)
		if bk == nil {
			return nil
		}
		bk.ForEach(func(k, v []byte) error {
			data := parseCacheData(v)
			if data == nil {
				return bk.Delete(k)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		logrus.Errorf("mainCacheClear err:%v", err)
	}
}
