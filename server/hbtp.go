package server

import (
	"github.com/gokins-main/gokins/comm"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
)

func runHbtp() {
	defer func() {
		if err := recover(); err != nil {
			hbtp.Errorf("Hbtp recover:%v", err)
		}
	}()
	comm.HbtpEgn = hbtp.NewEngine(comm.Ctx)
	err := comm.HbtpEgn.Run(comm.HbtpHost)
	if err != nil {
		logrus.Errorf("Hbtp err:%v", err)
	}
}
