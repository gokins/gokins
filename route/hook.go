package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/engine"
	"github.com/gokins/gokins/model"
	"github.com/gokins/gokins/service"
	"github.com/gokins/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
)

type HookController struct {
}

func (HookController) GetPath() string {
	return "/trigger"
}
func (c *HookController) Routes(g gin.IRoutes) {
	g.POST("/hook/:triggerId", c.hooks)
	g.POST("/web/:triggerId", util.GinReqParseJson(c.web))
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

func (HookController) web(c *gin.Context, m *hbtp.Map) {
	triggerId := c.Param("triggerId")
	secret := m.GetString("secret")
	if triggerId == "" || secret == "" {
		c.String(500, "param err")
		return
	}
	tt := &model.TTrigger{}
	ok, _ := comm.Db.Where("id = ? and enabled != 0", triggerId).Get(tt)
	if !ok {
		c.String(404, "触发器不存在或者未激活")
		return
	}
	rb, err := service.TriggerWeb(tt, secret)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	engine.Mgr.BuildEgn().Put(rb)
	c.JSON(200, gin.H{
		"msg": "ok",
	})
}
