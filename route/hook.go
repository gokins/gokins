package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/gokins/engine"
	"github.com/gokins-main/gokins/service"
)

type ReposController struct {
}

func (ReposController) GetPath() string {
	return "/api/repos"
}
func (c *ReposController) Routes(g gin.IRoutes) {
	g.POST("/hooks/:hookType", c.hooks)
}

func (ReposController) hooks(c *gin.Context) {
	hookType := c.Param("hookType")
	if hookType == "" {
		c.JSON(200, gin.H{
			"msg": "hook类型为空",
		})
		return
	}
	rb, err := service.Parse(c.Request, hookType)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	engine.Mgr.BuildEgn().Put(rb)
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}
