package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/service"
	"github.com/gokins-main/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
)

type OrgController struct{}

func (OrgController) GetPath() string {
	return "/api/org"
}
func (c *OrgController) Routes(g gin.IRoutes) {
	g.Use(service.MidUserCheck)
	g.POST("/new", util.GinReqParseJson(c.new))
}
func (OrgController) new(c *gin.Context, m *hbtp.Map) {
	name := m.GetString("name")
	desc := m.GetString("desc")
	pub := m.GetBool("public")
	if name == "" {
		c.String(500, "param err")
		return
	}
	usr := service.GetMidLgUser(c)
	ne := &model.TOrg{
		Id:      utils.NewXid(),
		Uid:     usr.Id,
		Name:    name,
		Desc:    desc,
		Created: time.Now(),
		Updated: time.Now(),
	}
	if pub {
		ne.Public = 1
	}
	_, err := comm.Db.InsertOne(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.JSON(200, &bean.IdsRes{
		Id:  ne.Id,
		Aid: ne.Aid,
	})
}
