package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/gokins/engine"
	"github.com/gokins-main/gokins/service"
	"github.com/gokins-main/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
)

type RuntimeController struct{}

func (RuntimeController) GetPath() string {
	return "/api/runtime"
}
func (c *RuntimeController) Routes(g gin.IRoutes) {
	g.Use(service.MidUserCheck)
	g.POST("/build", util.GinReqParseJson(c.build))
	g.POST("/cancel", util.GinReqParseJson(c.cancel))
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
	c.JSON(200, show)
}
func (RuntimeController) cancel(c *gin.Context, m *hbtp.Map) {
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
	v.Cancel()
	c.String(200, "ok")
}
