package engine

import (
	"fmt"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/runner/runners"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"strconv"
	"time"
)

type HbtpRunner struct {
}

func (HbtpRunner) AuthFun() hbtp.AuthFun {
	return func(c *hbtp.Context) bool {
		cmds := c.Command()
		times := c.Args().Get("times")
		random := c.Args().Get("random")
		sign := c.Args().Get("sign")
		if cmds == "" || len(times) <= 5 || len(random) < 20 || sign == "" {
			c.ResString(hbtp.ResStatusAuth, "auth param err")
			return false
		}
		signs := utils.Md5String(cmds + random + times + comm.Cfg.Server.Secret)
		if sign != signs {
			println("token err:" + sign)
			c.ResString(hbtp.ResStatusAuth, "token err:"+sign)
			return false
		}
		tm, err := strconv.ParseInt(times, 10, 64)
		if err != nil {
			c.ResString(hbtp.ResStatusAuth, "times err:"+err.Error())
			return false
		}
		tms := time.Unix(tm, 0)
		/*if err != nil {
			c.ResString(hbtp.ResStatusAuth, "times err:"+err.Error())
			return false
		}*/
		hbtp.Debugf("HbtpRunnerAuth parse.times:%s", tms.Format(common.TimeFmt))
		return true
	}
}
func (HbtpRunner) ServerInfo(c *hbtp.Context) {
	rts, err := Mgr.brun.ServerInfo()
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResJson(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) PullJob(c *hbtp.Context, m *runners.ReqPullJob) {
	rts, err := Mgr.brun.PullJob(m.Name, m.Plugs)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResJson(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) CheckCancel(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	c.ResString(hbtp.ResStatusOk, fmt.Sprintf("%t", Mgr.brun.CheckCancel(buildId)))
}
func (HbtpRunner) Update(c *hbtp.Context, m *runners.UpdateJobInfo) {
	err := Mgr.brun.Update(m)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, "ok")
}
func (HbtpRunner) UpdateCmd(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	jobId := c.ReqHeader().GetString("jobId")
	cmdId := c.ReqHeader().GetString("cmdId")
	fs, err := c.ReqHeader().GetInt("fs")
	code, _ := c.ReqHeader().GetInt("code")
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	err = Mgr.brun.UpdateCmd(buildId, jobId, cmdId, int(fs), int(code))
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, "ok")
}
func (HbtpRunner) PushOutLine(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	jobId := c.ReqHeader().GetString("jobId")
	cmdId := c.ReqHeader().GetString("cmdId")
	bs := c.ReqHeader().GetString("bs")
	iserr := c.ReqHeader().GetBool("iserr")
	err := Mgr.brun.PushOutLine(buildId, jobId, cmdId, bs, iserr)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, "ok")
}
func (HbtpRunner) FindJobId(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	stgNm := c.ReqHeader().GetString("stgNm")
	stpNm := c.ReqHeader().GetString("stpNm")
	rts, ok := Mgr.brun.FindJobId(buildId, stgNm, stpNm)
	if !ok {
		c.ResString(hbtp.ResStatusNotFound, "")
		return
	}
	c.ResString(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) ReadDir(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	pth := c.ReqHeader().GetString("pth")
	fs, err := c.ReqHeader().GetInt("fs")
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	rts, err := Mgr.brun.ReadDir(int(fs), buildId, pth)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResJson(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) ReadFile(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	pth := c.ReqHeader().GetString("pth")
	fs, err := c.ReqHeader().GetInt("fs")
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	flsz, flr, err := Mgr.brun.ReadFile(int(fs), buildId, pth)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	defer flr.Close()
	c.ResString(hbtp.ResStatusOk, fmt.Sprintf("%d", flsz))
	bts := make([]byte, 10240)
	for !hbtp.EndContext(comm.Ctx) {
		n, err := flr.Read(bts)
		if n > 0 {
			_, err = c.Conn().Write(bts[:n])
			if err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
}
func (HbtpRunner) GetEnv(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	jobId := c.ReqHeader().GetString("jobId")
	key := c.ReqHeader().GetString("key")
	rts, ok := Mgr.brun.GetEnv(buildId, jobId, key)
	if !ok {
		c.ResString(hbtp.ResStatusNotFound, "")
		return
	}
	c.ResString(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) GenEnv(c *hbtp.Context, env utils.EnvVal) {
	buildId := c.ReqHeader().GetString("buildId")
	jobId := c.ReqHeader().GetString("jobId")
	err := Mgr.brun.GenEnv(buildId, jobId, env)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, "ok")
}
func (HbtpRunner) UploadFile(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	jobId := c.ReqHeader().GetString("jobId")
	dir := c.ReqHeader().GetString("dir")
	pth := c.ReqHeader().GetString("pth")
	fs, err := c.ReqHeader().GetInt("fs")
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	flw, err := Mgr.brun.UploadFile(int(fs), buildId, jobId, dir, pth)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	defer flw.Close()
	c.ResString(hbtp.ResStatusOk, "ok")

	bts := make([]byte, 10240)
	for !hbtp.EndContext(comm.Ctx) {
		n, err := c.Conn().Read(bts)
		if n > 0 {
			_, err = flw.Write(bts[:n])
			if err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
}
func (HbtpRunner) FindArtVersionId(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	idnt := c.ReqHeader().GetString("idnt")
	name := c.ReqHeader().GetString("name")
	rts, err := Mgr.brun.FindArtVersionId(buildId, idnt, name)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) NewArtVersionId(c *hbtp.Context) {
	buildId := c.ReqHeader().GetString("buildId")
	idnt := c.ReqHeader().GetString("idnt")
	name := c.ReqHeader().GetString("name")
	rts, err := Mgr.brun.NewArtVersionId(buildId, idnt, name)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, rts)
}
