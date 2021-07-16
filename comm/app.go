package comm

import (
	"context"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"xorm.io/xorm"
)

var (
	Ctx  context.Context
	cncl context.CancelFunc
)
var (
	Cfg       = Config{}
	Db        *xorm.Engine
	BCache    *bolt.DB
	WebEgn    *gin.Engine
	HbtpEgn   *hbtp.Engine
	IsMySQL   = false
	Installed = false
	NotUpPass = false
	WorkPath  = ""
	WebHost   = ""
	//HbtpHost = ""
)

func init() {
	Ctx, cncl = context.WithCancel(context.Background())
}
func Cancel() {
	if cncl != nil {
		cncl()
	}
}
