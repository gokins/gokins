package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins/core/utils"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/model"
	"github.com/gokins/gokins/models"
	"github.com/gokins/gokins/service"
	"github.com/gokins/gokins/util"
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
	g.POST("/info", util.GinReqParseJson(c.info))
	g.POST("/upinfo", util.GinReqParseJson(c.upinfo))
	g.POST("/upass", util.GinReqParseJson(c.upass))
	g.POST("/active", util.GinReqParseJson(c.active))
	g.POST("/perm", util.GinReqParseJson(c.perm))
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

func (UserController) info(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	if id == "" {
		c.String(500, "param err")
		return
	}
	usr := &models.TUser{}
	ok := service.GetIdOrAid(id, usr)
	if !ok {
		c.String(404, "not found user")
		return
	}
	uinfo, ok := service.GetUserInfo(usr.Id)
	c.JSON(200, hbtp.Map{
		"user": usr,
		"info": uinfo,
	})
}
func (UserController) upinfo(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	nick := strings.TrimSpace(m.GetString("nick"))
	phone := m.GetString("phone")
	email := m.GetString("email")
	remark := m.GetString("remark")
	if id == "" || nick == "" {
		c.String(500, "param err")
		return
	}
	usr := &models.TUser{}
	ok := service.GetIdOrAid(id, usr)
	if !ok {
		c.String(404, "not found user")
		return
	}
	lgusr := service.GetMidLgUser(c)
	if !service.IsAdmin(lgusr) && usr.Id != lgusr.Id {
		c.String(405, "is not you")
		return
	}
	uinfo, isup := service.GetUserInfo(usr.Id)
	usr.Nick = nick
	_, err := comm.Db.Cols("nick").Where("id=?", usr.Id).Update(usr)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	uinfo.Phone = phone
	uinfo.Email = email
	uinfo.Remark = remark
	if isup {
		_, err = comm.Db.Cols("phone", "email", "remark").
			Where("id=?", usr.Id).Update(uinfo)
	} else {
		uinfo.Id = usr.Id
		_, err = comm.Db.InsertOne(uinfo)
	}
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	service.ClearUserCache(usr.Id)
	c.String(200, usr.Id)
}
func (UserController) upass(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	olds := m.GetString("olds")
	pass := m.GetString("pass")
	if id == "" || pass == "" {
		c.String(500, "param err")
		return
	}
	lgusr := service.GetMidLgUser(c)
	usr := &model.TUser{}
	if id == lgusr.Id {
		usr = lgusr
	} else {
		ok := service.GetIdOrAid(id, usr)
		if !ok {
			c.String(404, "not found user")
			return
		}
	}

	if comm.NotUpPass && !service.IsAdmin(lgusr) {
		c.String(513, "can't update")
		return
	}
	if usr.Id == lgusr.Id {
		if olds == "" {
			c.String(511, "param err1")
			return
		}
		if usr.Pass != utils.Md5String(olds) {
			c.String(512, "old pass err")
			return
		}
	} else if !service.IsAdmin(lgusr) {
		c.String(405, "is not admin")
		return
	}

	usr.Pass = utils.Md5String(pass)
	_, err := comm.Db.Cols("pass").Where("id=?", usr.Id).Update(usr)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	service.ClearUserCache(usr.Id)
	c.String(200, usr.Id)
}
func (UserController) active(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	act := m.GetString("act")
	if id == "" || act == "" {
		c.String(500, "param err")
		return
	}
	lgusr := service.GetMidLgUser(c)
	if !service.IsAdmin(lgusr) {
		c.String(405, "is not admin")
		return
	}
	usr := &model.TUser{}
	ok := service.GetIdOrAid(id, usr)
	if !ok {
		c.String(404, "not found user")
		return
	}
	if act == "1" {
		usr.Active = 1
	} else {
		usr.Active = 0
	}
	_, err := comm.Db.Cols("active").Where("id=?", usr.Id).Update(usr)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	service.ClearUserCache(usr.Id)
	c.String(200, usr.Id)
}
func (UserController) perm(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	permUser := m.GetBool("permUser")
	permOrg := m.GetBool("permOrg")
	permPipe := m.GetBool("permPipe")
	if id == "" {
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
	usr := &models.TUser{}
	ok := service.GetIdOrAid(id, usr)
	if !ok {
		c.String(404, "not found user")
		return
	}
	uinfo, isup := service.GetUserInfo(usr.Id)
	if permUser {
		uinfo.PermUser = 1
	} else {
		uinfo.PermUser = 0
	}
	if permOrg {
		uinfo.PermOrg = 1
	} else {
		uinfo.PermOrg = 0
	}
	if permPipe {
		uinfo.PermPipe = 1
	} else {
		uinfo.PermPipe = 0
	}
	var err error
	if isup {
		_, err = comm.Db.Cols("perm_user", "perm_org", "perm_pipe").
			Where("id=?", usr.Id).Update(uinfo)
	} else {
		uinfo.Id = usr.Id
		_, err = comm.Db.InsertOne(uinfo)
	}
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	service.ClearUserCache(usr.Id)
	c.String(200, usr.Id)
}
