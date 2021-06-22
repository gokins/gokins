package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/gokins/comm"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
	"time"
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
		time.Sleep(time.Millisecond * 100)
		comm.HbtpEgn.Stop()
		comm.Cancel()
	}
}
