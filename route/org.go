package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/models"
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
	g.POST("/list", util.GinReqParseJson(c.list))
	g.POST("/new", util.GinReqParseJson(c.new))
	g.POST("/info", util.GinReqParseJson(c.info))
	g.POST("/save", util.GinReqParseJson(c.save))
}

func (OrgController) list(c *gin.Context, m *hbtp.Map) {
	var ls []*models.TOrgInfo
	usr := service.GetMidLgUser(c)
	q := m.GetString("q")
	pg, _ := m.GetInt("page")

	var err error
	var page *bean.Page
	if comm.IsMySQL {
		gen := &bean.PageGen{
			CountCols: "org.id",
			FindCols:  "org.*,urg.perm_adm,urg.perm_rw,urg.perm_exec",
		}
		gen.SQL = `
		select {{select}} from t_org org
		LEFT JOIN t_user_org urg on urg.uid=?
		where deleted!=1
		and (org.uid=? or org.id=urg.org_id)
		`
		gen.Args = append(gen.Args, usr.Id)
		gen.Args = append(gen.Args, usr.Id)
		if q != "" {
			gen.SQL += "\nAND org.name like ? "
			gen.Args = append(gen.Args, "%"+q+"%")
		}
		gen.SQL += "\nORDER BY org.aid DESC"
		page, err = comm.FindPages(gen, &ls, pg, 20)
	}
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.JSON(200, page)
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

func (OrgController) info(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	if id == "" {
		c.String(500, "param err")
		return
	}
	org := &models.TOrg{}
	ok := service.GetIdOrAid(id, org)
	if !ok || org.Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	usr := &models.TUser{}
	ok = service.GetIdOrAid(org.Uid, usr)
	if !ok {
		c.String(404, "not found user?")
		return
	}
	c.JSON(200, hbtp.Map{
		"org":  org,
		"user": usr,
	})
}
func (OrgController) save(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	name := m.GetString("name")
	desc := m.GetString("desc")
	pub := m.GetBool("public")
	if name == "" {
		c.String(500, "param err")
		return
	}
	org := &model.TOrg{}
	ok := service.GetIdOrAid(id, org)
	if !ok || org.Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	usr := service.GetMidLgUser(c)
	if org.Uid != usr.Id {
		urg := &model.TUserOrg{}
		ok, _ = comm.Db.Where("uid=? and org_id=?", usr.Id, org.Id).Get(urg)
		if !ok || urg.PermAdm != 1 {
			c.String(405, "no permission")
			return
		}
	}
	ne := &model.TOrg{
		Name:    name,
		Desc:    desc,
		Updated: time.Now(),
	}
	if pub {
		ne.Public = 1
	}
	_, err := comm.Db.Cols("name", "desc", "public", "updated").
		Where("id=?", org.Id).Update(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.JSON(200, &bean.IdsRes{
		Id:  ne.Id,
		Aid: ne.Aid,
	})
}
