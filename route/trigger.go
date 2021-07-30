package route

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gokins/core/utils"
	"github.com/gokins/gokins/bean"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/engine"
	"github.com/gokins/gokins/model"
	"github.com/gokins/gokins/models"
	"github.com/gokins/gokins/service"
	"github.com/gokins/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"time"
)

type TriggerController struct{}

func (TriggerController) GetPath() string {
	return "/api/trigger"
}
func (c *TriggerController) Routes(g gin.IRoutes) {
	g.Use(service.MidUserCheck)
	g.POST("/triggers", util.GinReqParseJson(c.triggers))
	g.POST("/save", util.GinReqParseJson(c.save))
	g.POST("/delete", util.GinReqParseJson(c.delete))
	g.POST("/runs", util.GinReqParseJson(c.runs))
}

func (TriggerController) triggers(c *gin.Context, m *hbtp.Map) {
	pipelineId := m.GetString("pipelineId")
	types := m.GetString("types")
	q := m.GetString("q")
	pg, _ := m.GetInt("page")
	if pipelineId == "" {
		c.String(500, "param err")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(lgusr, pipelineId)
	if perm.Pipeline() == nil {
		c.String(404, "流水线不存在")
		return
	}
	if !perm.IsAdmin() {
		if !perm.CanRead() {
			c.String(405, "No Auth")
			return
		}
	}
	ls := make([]*models.TTrigger, 0)
	session := comm.Db.NewSession()
	if pipelineId != "" {
		session.And("pipeline_id = ?", pipelineId)
	}
	if types != "" {
		session.And("types = ?", types)
	}
	if q != "" {
		session.And("name like '%" + q + "%'")
	}
	page, err := comm.FindPage(session, &ls, pg)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	for _, v := range ls {
		usr, ok := service.GetUser(v.Uid)
		if ok {
			v.Nick = usr.Nick
			v.Avat = usr.Avatar
		}
		_ = json.Unmarshal([]byte(v.Params), &v.Param)
	}
	ms := map[string]interface{}{}
	ms["page"] = page
	ms["host"] = comm.Cfg.Server.Host
	c.JSON(200, ms)
}

func (TriggerController) save(c *gin.Context, tp *bean.TriggerParam) {
	if err := tp.Check(); err != nil {
		c.String(500, err.Error())
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(lgusr, tp.PipelineId)
	if perm.Pipeline() == nil {
		c.String(404, "流水线不存在")
		return
	}
	if !perm.IsAdmin() && !perm.CanWrite() {
		c.String(405, "No Auth")
		return
	}
	tt := &model.TTrigger{}
	err := utils.Struct2Struct(tt, tp)
	if err != nil {
		c.String(500, "Struct2Struct err:"+err.Error())
		return
	}
	if tp.Enabled {
		tt.Enabled = 1
	}
	if tp.Id == "" {
		tt.Id = utils.NewXid()
		tt.Created = time.Now()
		tt.Uid = lgusr.Id
		_, err = comm.Db.InsertOne(tt)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
	} else {
		tt.Updated = time.Now()
		_, err = comm.Db.Cols("name,desc,params,types,enabled,updated").Where("id =?", tt.Id).Update(tt)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
	}
	if tt.Types == "timer" {
		engine.Mgr.TimerEng().Refresh(tt.Id)
	}
	c.JSON(200, "ok")
}

func (TriggerController) delete(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	tt := &model.TTrigger{}
	ok, _ := comm.Db.Where("id = ?", id).Get(tt)
	if !ok {
		c.String(404, "触发器不存在")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(lgusr, tt.PipelineId)
	if perm.Pipeline() == nil {
		c.String(404, "流水线不存在")
		return
	}
	if !perm.IsAdmin() && !perm.CanWrite() {
		c.String(405, "No Auth")
		return
	}
	_, err := comm.Db.Where("id = ?", tt.Id).Delete(tt)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	tr := model.TTriggerRun{}
	_, err = comm.Db.Where("tid = ?", tt.Id).Delete(tr)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}

	if tt.Types == "timer" {
		engine.Mgr.TimerEng().Delete(tt.Id)
	}
	c.JSON(200, "ok")
}

func (TriggerController) runs(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	pg, _ := m.GetInt("page")
	tt := &model.TTrigger{}
	ok, _ := comm.Db.Where("id = ?", id).Get(tt)
	if !ok {
		c.String(404, "触发器不存在")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(lgusr, tt.PipelineId)
	if perm.Pipeline() == nil {
		c.String(404, "流水线不存在")
		return
	}
	if !perm.IsAdmin() && !perm.CanRead() {
		c.String(405, "No Auth")
		return
	}
	var ls []*models.TTriggerRun
	session := comm.Db.Where("tid = ?", tt.Id).Desc("created")
	page, err := comm.FindPage(session, &ls, pg)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	for _, v := range ls {
		if v.Error != "" || v.PipeVersionId == "" {
			continue
		}
		rpv := &models.RunPipelineVersion{}
		ok, _ = comm.Db.Table("t_pipeline_version").
			Where("t_pipeline_version.id = ?", v.PipeVersionId).
			Join("left", "t_build", "t_build.pipeline_version_id = ?", v.PipeVersionId).
			Get(rpv)
		if ok {
			v.Number = rpv.Number
			v.PipelineName = rpv.PipelineName
			v.PipelineDisplayName = rpv.PipelineDisplayName
			v.BStatus = rpv.Status
		}
	}
	c.JSON(200, page)
}
