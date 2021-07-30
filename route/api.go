package route

import (
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/gokins/core/runtime"
	"github.com/gokins/core/utils"
	"github.com/gokins/gokins/bean"
	"github.com/gokins/gokins/engine"
	"github.com/gokins/gokins/util"
	"gopkg.in/yaml.v3"
)

type ApiController struct{}

func (ApiController) GetPath() string {
	return "/api"
}
func (c *ApiController) Routes(g gin.IRoutes) {
	g.Any("/", c.hello)
	g.POST("/builds", util.GinReqParseJson(c.test))
}
func (ApiController) hello(c *gin.Context) {
	c.String(200, "hello world")
}
func (ApiController) test(c *gin.Context) {
	all, err := ioutil.ReadAll(c.Request.Body)
	y := &bean.Pipeline{}
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
	b := &runtime.Build{}
	err = yaml.Unmarshal(marshal, b)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err,
		})
		return
	}
	err = prebuild(b)
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

func prebuild(b *runtime.Build) error {
	if b == nil {
		return errors.New("build is empty")
	}
	if b.Stages == nil || len(b.Stages) <= 0 {
		return errors.New("stages is empty")
	}
	pipelineId := utils.NewXid()
	buildId := utils.NewXid()
	b.Id = buildId
	b.PipelineId = pipelineId
	b.Repo = &runtime.Repository{
		Name:     "",
		Token:    "",
		Sha:      "",
		CloneURL: "https://gitee.com/SuperHeroJim/gokins-test.git",
	}
	for _, stage := range b.Stages {
		stage.Id = utils.NewXid()
		stage.BuildId = buildId
		for _, step := range stage.Steps {
			step.Id = utils.NewXid()
			step.StageId = stage.Id
			step.BuildId = buildId
		}
	}
	return nil
}
