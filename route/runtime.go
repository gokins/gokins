package route

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gokins/core/common"
	"github.com/gokins/gokins/bean"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/engine"
	"github.com/gokins/gokins/model"
	"github.com/gokins/gokins/models"
	"github.com/gokins/gokins/service"
	"github.com/gokins/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"os"
	"path/filepath"
)

type RuntimeController struct{}

func (RuntimeController) GetPath() string {
	return "/api/runtime"
}
func (c *RuntimeController) Routes(g gin.IRoutes) {
	g.Use(service.MidUserCheck)
	g.POST("/stages", util.GinReqParseJson(c.stages))
	g.POST("/cmds", util.GinReqParseJson(c.cmds))
	g.POST("/build", util.GinReqParseJson(c.build))
	g.POST("/cancel", util.GinReqParseJson(c.cancel))
	g.POST("/logs", util.GinReqParseJson(c.logs))
}
func (RuntimeController) stages(c *gin.Context, m *hbtp.Map) {
	pvId := m.GetString("pvId")
	if pvId == "" {
		c.String(500, "param err")
		return
	}
	var ls []*models.RunStage
	err := comm.Db.Where("pipeline_version_id=?", pvId).OrderBy("sort ASC").Find(&ls)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	ids := make([]string, 0)
	stages := map[string]*models.RunStage{}
	steps := map[string]*models.RunStep{}
	for _, v := range ls {
		v.Stepids = make([]string, 0)
		var spls []*models.RunStep
		err := comm.Db.Where("stage_id=?", v.Id).OrderBy("sort ASC").Find(&spls)
		if err == nil {
			ids = append(ids, v.Id)
			stages[v.Id] = v
			for _, step := range spls {
				if step.Id != "" {
					json.Unmarshal([]byte(step.Waits), &step.Waitings)
					v.Stepids = append(v.Stepids, step.Id)
					steps[step.Id] = step
				}
			}
		}
	}
	c.JSON(200, hbtp.Map{
		"ids":    ids,
		"stages": stages,
		"steps":  steps,
	})
}
func (RuntimeController) cmds(c *gin.Context, m *hbtp.Map) {
	stepId := m.GetString("stepId")
	if stepId == "" {
		c.String(500, "param err")
		return
	}
	var ls []*model.TCmdLine
	err := comm.Db.Where("step_id=?", stepId).OrderBy("num ASC").Find(&ls)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.JSON(200, hbtp.Map{
		"stepId": stepId,
		"cmds":   ls,
	})
	/*ids:=make([]string,0)
	cmds := map[string]*model.TCmdLine{}
	for _, v := range ls {
		ids = append(ids, v.Id)
		cmds[v.Id] = v
	}
	c.JSON(200, hbtp.Map{
		//"ids":    ids,
		"cmds": ls,
	})*/
}
func (RuntimeController) build(c *gin.Context, m *hbtp.Map) {
	bdid := m.GetString("buildId")
	if bdid == "" {
		c.String(500, "param err")
		return
	}
	v, ok := engine.Mgr.BuildEgn().Get(bdid)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	show, ok := v.Show()
	if !ok {
		c.String(404, "Not Found")
		return
	}
	c.JSON(200, hbtp.Map{
		"workpgss": v.WorkProgress(),
		"show":     show,
	})
}
func (RuntimeController) cancel(c *gin.Context, m *hbtp.Map) {
	bdid := m.GetString("buildId")
	if bdid == "" {
		c.String(500, "param err")
		return
	}
	build := &model.TBuild{}
	ok, _ := comm.Db.Where("id=?", bdid).Get(build)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	perm := service.NewPipePerm(service.GetMidLgUser(c), build.PipelineId)
	if !perm.CanExec() {
		c.String(405, "No Permission")
		return
	}
	v, ok := engine.Mgr.BuildEgn().Get(build.Id)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	v.Cancel()
	c.String(200, "ok")
}
func (RuntimeController) logs(c *gin.Context, m *hbtp.Map) {
	buildId := m.GetString("buildId")
	stepId := m.GetString("stepId")
	offset, _ := m.GetInt("offset")
	limit, _ := m.GetInt("limit")
	if stepId == "" {
		c.String(500, "param err")
		return
	}
	/*tstp := &model.TStep{}
	ok, _ := comm.Db.Where("id=?", stepId).Get(tstp)
	if !ok {
		c.String(404, "Not Found")
		return
	}*/
	dir := filepath.Join(comm.WorkPath, common.PathBuild, buildId, common.PathJobs, stepId)
	logpth := filepath.Join(dir, "build.log")
	fl, err := os.Open(logpth)
	if err != nil {
		c.String(404, "Not Found File")
		return
	}
	defer fl.Close()
	off := offset
	if offset > 0 {
		off, err = fl.Seek(offset, 0)
		if err != nil {
			c.String(510, "err:%v", err)
			return
		}
	}

	var lastoff int64
	ls := make([]*bean.LogOutJsonRes, 0)
	bts := make([]byte, 1024*5)
	linebuf := &bytes.Buffer{}
	for !hbtp.EndContext(c) {
		rn, err := fl.Read(bts)
		if rn > 0 {
			for i := 0; i < rn; i++ {
				off++
				b := bts[i]
				if linebuf == nil && b == '{' {
					linebuf.Reset()
				}
				if linebuf != nil {
					if b == '\n' {
						e := &bean.LogOutJsonRes{}
						err := json.Unmarshal(linebuf.Bytes(), e)
						linebuf.Reset()
						if err == nil {
							/*if e.Type == hbtpBean.TypeCmdLogLineSys {
								continue
							}*/
							e.Offset = off - 1
							ls = append(ls, e)
							lastoff = e.Offset
						}
						if limit > 0 && limit >= int64(len(ls)) {
							break
						}
					} else {
						linebuf.WriteByte(b)
					}
				}
			}
		}
		if err != nil {
			break
		}
	}
	c.JSON(200, hbtp.Map{
		"stepId":  stepId,
		"lastoff": lastoff,
		"logs":    ls,
	})
}
