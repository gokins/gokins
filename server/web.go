package server

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/route"
	"github.com/gokins-main/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
)

func runWeb() {
	defer func() {
		if err := recover(); err != nil {
			hbtp.Errorf("Web recover:%v", err)
		}
	}()
	comm.WebEgn = gin.Default()
	err := comm.WebEgn.Run(comm.WebHost)
	if err != nil {
		logrus.Errorf("Web err:%v", err)
		//comm.HbtpEgn.Stop()
	}
	comm.Cancel()
	time.Sleep(time.Millisecond * 100)
}

func regApi() {
	if core.Debug {
		comm.WebEgn.Use(util.MidAccessAllowFun)
	}
	util.GinRegController(comm.WebEgn, &route.ApiController{})
	util.GinRegController(comm.WebEgn, &route.LoginController{})
	util.GinRegController(comm.WebEgn, &route.UserController{})
	util.GinRegController(comm.WebEgn, &route.OrgController{})
	util.GinRegController(comm.WebEgn, &route.PipelineController{})
}
func http404(c *gin.Context) {

}
