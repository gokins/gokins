package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/models"
	"github.com/gokins/gokins/util"
)

type YmlController struct{}

func (YmlController) GetPath() string {
	return "/api/yml"
}
func (c *YmlController) Routes(g gin.IRoutes) {
	g.POST("/templates", util.GinReqParseJson(c.templates))
	g.POST("/plugins", util.GinReqParseJson(c.plugins))
}
func (YmlController) templates(c *gin.Context) {
	ls := make([]*models.TYmlTemplate, 0)
	comm.Db.Where("deleted != 1").Find(&ls)
	c.JSON(200, ls)
}

func (YmlController) plugins(c *gin.Context) {
	ls := make([]*models.TYmlPlugin, 0)
	comm.Db.Where("deleted != 1").Find(&ls)
	c.JSON(200, ls)
}
