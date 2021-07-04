package route

import (
	"github.com/gokins-main/gokins/models"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core/common"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/service"
	"github.com/gokins-main/gokins/util"
)

type LoginController struct{}

func (LoginController) GetPath() string {
	return "/api/lg"
}
func (c *LoginController) Routes(g gin.IRoutes) {
	g.POST("/info", c.info)
	g.POST("/login", util.GinReqParseJson(c.login))
}
func (LoginController) info(c *gin.Context) {
	rt := hbtp.Map{}
	usr, ok := service.CurrUserCache(c)
	if ok {
		usrs := &models.TUser{}
		utils.Struct2Struct(usrs, usr)
		rt["user"] = usrs
		info, _ := service.GetUserInfo(usrs.Id)
		rt["info"] = info
		if service.IsAdmin(usr) {
			info.PermUser = 1
			info.PermOrg = 1
			info.PermPipe = 1
		}
	}
	rt["login"] = ok
	c.JSON(200, rt)
}
func (LoginController) login(c *gin.Context, m *bean.LoginReq) {
	m.Name = strings.TrimSpace(m.Name)
	if m.Name == "" || m.Pass == "" {
		c.String(500, "param err")
		return
	}
	usr, ok := service.FindUserName(m.Name)
	if !ok {
		c.String(404, "not found user")
		return
	}
	if !service.IsAdmin(usr) && usr.Active != 1 {
		c.String(513, "user not active")
		return
	}
	if usr.Pass != utils.Md5String(m.Pass) {
		c.String(511, "password err")
		return
	}
	key := comm.Cfg.Server.LoginKey
	if key == "" {
		c.String(512, "no set login key")
		return
	}
	token, err := util.CreateToken(jwt.MapClaims{
		"uid": usr.Id,
	}, key, time.Hour*24*5)
	if err != nil {
		c.String(500, "create token err:%v", err)
		return
	}
	rt := &bean.LoginRes{
		Token:         token,
		Id:            usr.Id,
		Name:          usr.Name,
		Nick:          usr.Nick,
		Avatar:        usr.Avatar,
		LastLoginTime: usr.LoginTime.Format(common.TimeFmt),
	}
	c.JSON(200, rt)

	usr.LoginTime = time.Now()
	comm.Db.Cols("login_time").Where("id=?", usr.Id).Update(usr)
}
