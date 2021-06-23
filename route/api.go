package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/gokins/engine"
	"github.com/gokins-main/gokins/util"
	"github.com/gokins-main/gokins/yml"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ApiController struct{}

func (ApiController) GetPath() string {
	return "/api"
}
func (c *ApiController) Routes(g gin.IRoutes) {
	g.POST("/builds", util.GinReqParseJson(c.test))
}
func (ApiController) test(c *gin.Context) {
	all, err := ioutil.ReadAll(c.Request.Body)
	y := &yml.YML{}
	err = yaml.Unmarshal(all, y)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	marshal, err := yaml.Marshal(y)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	//TODO insert db

	b := &runtime.Build{}
	err = yaml.Unmarshal(marshal, b)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	engine.Mgr.BuildEgn().Put(b)
	c.JSON(200, gin.H{
		"msg": b,
	})
}
