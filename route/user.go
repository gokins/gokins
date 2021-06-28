package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/models"
	"github.com/gokins-main/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
)

type UserController struct{}

func (UserController) GetPath() string {
	return "/api/user"
}
func (c *UserController) Routes(g gin.IRoutes) {
	g.POST("/page", util.GinReqParseJson(c.page))
}
func (UserController) page(c *gin.Context, m *hbtp.Map) {
	var ls []*models.TUser
	q := m.GetString("q")
	pg, _ := m.GetInt("page")

	ses := comm.Db.OrderBy("aid ASC")
	if q != "" {
		ses.And("name like ? or nick like ?", "%"+q+"%", "%"+q+"%")
	}

	page, err := comm.FindPage(ses, &ls, pg, 20)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.JSON(200, page)
}
