package comm

import (
	"context"
	"github.com/gin-gonic/gin"
	"xorm.io/xorm"
)

var (
	Ctx  context.Context
	cncl context.CancelFunc
)
var (
	Cfg    = Config{}
	Db     *xorm.Engine
	WebEgn = gin.Default()
	//HbtpEgn  *hbtp.Engine
	Installed = false
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
