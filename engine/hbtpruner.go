package engine

import (
	"encoding/json"
	"fmt"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/runner/runners"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
)

type HbtpRunner struct {
}

func (HbtpRunner) AuthFun() hbtp.AuthFun {
	return nil
}
func (HbtpRunner) ServerInfo(c *hbtp.Context) {
	c.ResJson(hbtp.ResStatusOk, Mgr.brun.ServerInfo())
}
func (HbtpRunner) PullJob(c *hbtp.Context, plugs []string) {
	rts, err := Mgr.brun.PullJob(plugs)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResJson(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) CheckCancel(c *hbtp.Context, buildId string) {
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
func (HbtpRunner) UpdateCmd(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	jobId := m.GetString("jobId")
	cmdId := m.GetString("cmdId")
	fs, err := m.GetInt("fs")
	code, _ := m.GetInt("code")
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
func (HbtpRunner) PushOutLine(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	jobId := m.GetString("jobId")
	cmdId := m.GetString("cmdId")
	bs := m.GetString("bs")
	iserr := m.GetBool("iserr")
	err := Mgr.brun.PushOutLine(buildId, jobId, cmdId, bs, iserr)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, "ok")
}
func (HbtpRunner) FindJobId(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	stgNm := m.GetString("stgNm")
	stpNm := m.GetString("stpNm")
	rts, ok := Mgr.brun.FindJobId(buildId, stgNm, stpNm)
	if !ok {
		c.ResString(hbtp.ResStatusNotFound, "")
		return
	}
	c.ResString(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) ReadDir(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	pth := m.GetString("pth")
	fs, err := m.GetInt("fs")
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
func (HbtpRunner) ReadFile(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	pth := m.GetString("pth")
	fs, err := m.GetInt("fs")
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
	c.ResJson(hbtp.ResStatusOk, fmt.Sprintf("%d", flsz))
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
func (HbtpRunner) GetEnv(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	jobId := m.GetString("jobId")
	key := m.GetString("key")
	rts, ok := Mgr.brun.GetEnv(buildId, jobId, key)
	if !ok {
		c.ResString(hbtp.ResStatusNotFound, "")
		return
	}
	c.ResString(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) GenEnv(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	jobId := m.GetString("jobId")
	env, ok := m.Get("env")
	if !ok {
		c.ResString(hbtp.ResStatusNotFound, "")
		return
	}
	envs := utils.EnvVal{}
	bts, err := json.Marshal(env)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	err = json.Unmarshal(bts, &envs)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	err = Mgr.brun.GenEnv(buildId, jobId, envs)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, "ok")
}
func (HbtpRunner) UploadFile(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	jobId := m.GetString("jobId")
	dir := m.GetString("dir")
	pth := m.GetString("pth")
	fs, err := m.GetInt("fs")
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
	ln := int64(0)
	bts := make([]byte, 10240)
	for !hbtp.EndContext(comm.Ctx) {
		n, err := c.Conn().Read(bts)
		if n > 0 {
			ln += int64(n)
			_, err = flw.Write(bts[:n])
			if err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	c.ResJson(hbtp.ResStatusOk, fmt.Sprintf("%d", ln))
}
func (HbtpRunner) FindArtVersionId(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	idnt := m.GetString("idnt")
	name := m.GetString("name")
	rts, err := Mgr.brun.FindArtVersionId(buildId, idnt, name)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, rts)
}
func (HbtpRunner) NewArtVersionId(c *hbtp.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	idnt := m.GetString("idnt")
	name := m.GetString("name")
	rts, err := Mgr.brun.NewArtVersionId(buildId, idnt, name)
	if err != nil {
		c.ResString(hbtp.ResStatusErr, err.Error())
		return
	}
	c.ResString(hbtp.ResStatusOk, rts)
}
