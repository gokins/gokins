package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/model"
	"github.com/gokins/gokins/service"
	"github.com/gokins/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
)

type PipelineVersionController struct{}

func (PipelineVersionController) GetPath() string {
	return "/api/pipelineVersion"
}
func (c *PipelineVersionController) Routes(g gin.IRoutes) {
	g.Use(service.MidUserCheck)
	g.POST("/delete", util.GinReqParseJson(c.delete))
}

func (PipelineVersionController) delete(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	if id == "" {
		c.String(500, "param err")
		return
	}
	tpv := &model.TPipelineVersion{}
	ok, _ := comm.Db.Where("id = ? ", id).Get(tpv)
	if !ok {
		c.String(404, "not found pipe_var")
		return
	}
	perm := service.NewPipePerm(service.GetMidLgUser(c), tpv.PipelineId)
	if perm.Pipeline() == nil {
		c.String(404, "not found pipe")
		return
	}
	if !perm.CanWrite() {
		c.String(405, "no permission")
		return
	}
	tpv.Deleted = 1
	_, err := comm.Db.Cols("deleted").Where("id = ?", tpv.Id).Update(tpv)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, "ok")
}
