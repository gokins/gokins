package route

import (
	"fmt"
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
	g.POST("/users", util.GinReqParseJson(c.users))
	g.POST("/save", util.GinReqParseJson(c.save))
	g.POST("/user/edit", util.GinReqParseJson(c.userEdit))
	g.POST("/user/rm", util.GinReqParseJson(c.userRm))
	g.POST("/pipe/add", util.GinReqParseJson(c.pipeAdd))
	g.POST("/pipe/rm", util.GinReqParseJson(c.pipeRm))
}

func (OrgController) list(c *gin.Context, m *hbtp.Map) {
	var ls []*models.TOrgInfo
	q := m.GetString("q")
	pg, _ := m.GetInt("page")

	var err error
	var page *bean.Page
	lgusr := service.GetMidLgUser(c)
	if comm.IsMySQL {
		gen := &bean.PageGen{
			CountCols: "org.id",
			FindCols:  "org.*",
		}
		gen.SQL = `
		select {{select}} from t_org org
		where org.deleted!=1
		`
		if lgusr.Id != "admin" {
			gen.FindCols = "org.*,urg.perm_adm,urg.perm_rw,urg.perm_exec"
			gen.SQL = `
			select {{select}} from t_org org
			LEFT JOIN t_user_org urg on urg.uid=?
			where org.deleted!=1
			and (org.uid=? or org.id=urg.org_id)
			`
			gen.Args = append(gen.Args, lgusr.Id)
			gen.Args = append(gen.Args, lgusr.Id)
		}
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
func (OrgController) users(c *gin.Context, m *hbtp.Map) {
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
	var usrs []*models.TUserOrgInfo
	if comm.IsMySQL {
		ses := comm.Db.SQL(`
		select usr.*,urg.perm_adm,urg.perm_rw,urg.perm_exec,urg.created as join_time from t_user usr
		JOIN t_user_org urg ON urg.org_id=?
		where usr.id=urg.uid
		ORDER BY urg.created ASC
		`, org.Id)
		err := ses.Find(&usrs)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
	}
	var usrsAdm []*models.TUserOrgInfo
	var usrsOtr []*models.TUserOrgInfo
	for _, v := range usrs {
		if v.PermAdm == 1 {
			usrsAdm = append(usrsAdm, v)
		} else {
			usrsOtr = append(usrsOtr, v)
		}
	}
	c.JSON(200, hbtp.Map{
		"adms": usrsAdm,
		"usrs": usrsOtr,
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
	lgusr := service.GetMidLgUser(c)
	if lgusr.Id != "admin" {
		if org.Uid != lgusr.Id {
			urg := &model.TUserOrg{}
			ok, _ = comm.Db.Where("uid=? and org_id=?", lgusr.Id, org.Id).Get(urg)
			if !ok || urg.PermAdm != 1 {
				c.String(405, "no permission")
				return
			}
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
func (OrgController) userEdit(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	uid := m.GetString("uid")
	adm := m.GetBool("adm")
	rw := m.GetBool("rw")
	ex := m.GetBool("ex")
	isadd := m.GetBool("add")
	org := &model.TOrg{}
	ok := service.GetIdOrAid(id, org)
	if !ok || org.Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	usr := &models.TUser{}
	ok = service.GetIdOrAid(uid, usr)
	if !ok {
		c.String(404, "not found user")
		return
	}
	var err error
	ne := &model.TUserOrg{}
	isup, _ := comm.Db.Where("uid=? and org_id=?", usr.Id, org.Id).Get(ne)
	lgusr := service.GetMidLgUser(c)
	if usr.Id == lgusr.Id {
		c.String(511, "can't edit yourself")
		return
	}
	if lgusr.Id != "admin" {
		if adm {
			if org.Uid != lgusr.Id {
				c.String(405, "no permission")
				return
			}
		} else {
			if org.Uid != lgusr.Id {
				urg := &model.TUserOrg{}
				ok, _ = comm.Db.Where("uid=? and org_id=?", lgusr.Id, org.Id).Get(urg)
				if !ok || urg.PermAdm != 1 {
					c.String(405, "no permission")
					return
				}
			}
		}
	}
	if adm {
		ne.PermAdm = 1
	} else {
		ne.PermAdm = 0
	}
	if !isadd {
		if rw {
			ne.PermRw = 1
		} else {
			ne.PermRw = 0
		}
		if ex {
			ne.PermExec = 1
		} else {
			ne.PermExec = 0
		}
	}
	if isup {
		_, err = comm.Db.Cols("perm_adm", "perm_rw", "perm_exec").Where("aid=?", ne.Aid).Update(ne)
	} else {
		ne.Uid = usr.Id
		ne.OrgId = org.Id
		ne.Created = time.Now()
		_, err = comm.Db.InsertOne(ne)
	}
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", ne.Aid))
}

func (OrgController) userRm(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	uid := m.GetString("uid")
	org := &model.TOrg{}
	ok := service.GetIdOrAid(id, org)
	if !ok || org.Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	usr := &models.TUser{}
	ok = service.GetIdOrAid(uid, usr)
	if !ok {
		c.String(404, "not found user")
		return
	}
	ne := &model.TUserOrg{}
	ok, _ = comm.Db.Where("uid=? and org_id=?", usr.Id, org.Id).Get(ne)
	if !ok {
		c.String(404, "not found user org")
		return
	}
	lgusr := service.GetMidLgUser(c)
	if usr.Id == lgusr.Id {
		c.String(511, "can't remove yourself")
		return
	}
	if lgusr.Id != "admin" {
		if org.Uid != lgusr.Id {
			urg := &model.TUserOrg{}
			ok, _ = comm.Db.Where("uid=? and org_id=?", lgusr.Id, org.Id).Get(urg)
			if !ok || urg.PermAdm != 1 {
				c.String(405, "no permission")
				return
			}
		}
	}
	_, err := comm.Db.Where("aid=?", ne.Aid).Delete(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", ne.Aid))
}

func (OrgController) pipeAdd(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	pipeId := m.GetString("pipeId")
	org := &model.TOrg{}
	ok := service.GetIdOrAid(id, org)
	if !ok || org.Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	lgusr := service.GetMidLgUser(c)
	if lgusr.Id != "admin" {
		if org.Uid != lgusr.Id {
			urg := &model.TUserOrg{}
			ok, _ = comm.Db.Where("uid=? and org_id=?", lgusr.Id, org.Id).Get(urg)
			if !ok || urg.PermAdm != 1 {
				c.String(405, "no permission")
				return
			}
		}
	}
	ne := &model.TOrgPipe{}
	ok, _ = comm.Db.Where("org_id=? and pipe_id=?", org.Id, pipeId).Get(ne)
	if ok {
		c.String(511, "pipeline exist")
		return
	}
	ne.OrgId = org.Id
	ne.PipeId = pipeId
	ne.Created = time.Now()
	_, err := comm.Db.InsertOne(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", ne.Aid))
}

func (OrgController) pipeRm(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	pipeId := m.GetString("pipeId")
	org := &model.TOrg{}
	ok := service.GetIdOrAid(id, org)
	if !ok || org.Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	lgusr := service.GetMidLgUser(c)
	if lgusr.Id != "admin" {
		if org.Uid != lgusr.Id {
			urg := &model.TUserOrg{}
			ok, _ = comm.Db.Where("uid=? and org_id=?", lgusr.Id, org.Id).Get(urg)
			if !ok || urg.PermAdm != 1 {
				c.String(405, "no permission")
				return
			}
		}
	}
	ne := &model.TOrgPipe{}
	_, err := comm.Db.Where("org_id=? and pipe_id=?", org.Id, pipeId).Delete(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", ne.Aid))
}
