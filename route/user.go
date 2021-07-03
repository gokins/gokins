package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/models"
	"github.com/gokins-main/gokins/service"
	"github.com/gokins-main/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"strings"
	"time"
)

type UserController struct{}

func (UserController) GetPath() string {
	return "/api/user"
}
func (c *UserController) Routes(g gin.IRoutes) {
	g.Use(service.MidUserCheck)
	g.POST("/page", util.GinReqParseJson(c.page))
	g.POST("/new", util.GinReqParseJson(c.new))
}
func (UserController) page(c *gin.Context, m *hbtp.Map) {
	var ls []*models.TUser
	q := m.GetString("q")
	pg, _ := m.GetInt("page")

	ses := comm.Db.OrderBy("aid DESC")
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
func (UserController) new(c *gin.Context, m *hbtp.Map) {
	name := strings.TrimSpace(m.GetString("name"))
	nick := strings.TrimSpace(m.GetString("nick"))
	pass := m.GetString("pass")
	//pmUser:=m.GetBool("pmUser")
	//pmOrg:=m.GetBool("pmOrg")
	//pmPipe:=m.GetBool("pmPipe")
	if name == "" || nick == "" || pass == "" {
		c.String(500, "param err")
		return
	}
	lgusr := service.GetMidLgUser(c)
	if !service.IsAdmin(lgusr) {
		uf, ok := service.GetUserInfo(lgusr.Id)
		if !ok || uf.PermUser != 1 {
			c.String(405, "no permission")
			return
		}
	}
	_, ok := service.FindUserName(name)
	if ok {
		c.String(511, "reged")
		return
	}
	ne := &model.TUser{
		Id:        utils.NewXid(),
		Name:      name,
		Pass:      utils.Md5String(pass),
		Nick:      nick,
		Created:   time.Now(),
		LoginTime: time.Now(),
		Active:    1,
	}
	/*if pmUser{
		ne.NewUser=1
	}
	if pmOrg{
		ne.NewOrg=1
	}
	if pmPipe{
		ne.NewPipe=1
	}*/
	_, err := comm.Db.InsertOne(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, ne.Id)
}
