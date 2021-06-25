package route

import (
	"fmt"
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
	rt := &bean.LgInfoRes{}
	usr, ok := service.CurrUserCache(c)
	if ok {
		rt.Login = true
		rt.Id = usr.Id
		rt.Name = usr.Name
		rt.Nick = usr.Nick
		rt.Avatar = usr.Avatar
		rt.LoginTime = usr.LoginTime.Format(common.TimeFmt)
		rt.RegTime = usr.Created.Format(common.TimeFmt)
	}
	c.JSON(200, rt)
}
func (LoginController) login(c *gin.Context, m *bean.LoginReq) {
	if m.Name == "" || m.Pass == "" {
		c.String(500, "param err")
		return
	}
	usr, ok := service.FindUserName(m.Name)
	if !ok {
		c.String(404, "not found user")
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
		"uid": fmt.Sprintf("%d", usr.Id),
	}, key, time.Hour*24*5)
	if err != nil {
		c.String(500, "create token err:%v", err)
		return
	}
	rt := &bean.LoginRes{
		Token:         token,
		Id:            fmt.Sprintf("%d", usr.Id),
		Name:          usr.Name,
		Nick:          usr.Nick,
		Avatar:        usr.Avatar,
		LastLoginTime: usr.LoginTime.Format(common.TimeFmt),
	}
	c.JSON(200, rt)

	usr.LoginTime = time.Now()
	comm.Db.Cols("login_time").Where("id=?", usr.Id).Update(usr)
}
