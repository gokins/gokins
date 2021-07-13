package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/engine"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/service"
)

type HookController struct {
}

func (HookController) GetPath() string {
	return "/hook"
}
func (c *HookController) Routes(g gin.IRoutes) {
	g.POST("/:triggerId", c.hooks)
}

func (HookController) hooks(c *gin.Context) {
	triggerId := c.Param("triggerId")
	if triggerId == "" {
		c.String(500, "param err")
		return
	}
	tt := &model.TTrigger{}
	ok, _ := comm.Db.Where("id = ? and enabled != 0", triggerId).Get(tt)
	if !ok {
		c.String(404, "触发器不存在或者未激活")
		return
	}
	rb, err := service.TriggerHook(tt, c.Request)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	engine.Mgr.BuildEgn().Put(rb)
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}
