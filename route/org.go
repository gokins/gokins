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
	g.POST("/rm", util.GinReqParseJson(c.rm))
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
		if !service.IsAdmin(lgusr) {
			/*gen.FindCols = "org.*,urg.perm_adm,urg.perm_rw,urg.perm_exec"
			gen.SQL = `
			select {{select}} from t_org org
			LEFT JOIN t_user_org urg on urg.uid=? and urg.org_id=org.id
			where org.deleted!=1
			and (org.public=1 or org.uid=?)
			`*/
			gen.FindCols = "org.*"
			gen.SQL = `
			select {{select}} from t_org org
			where org.deleted!=1 and
			(
			org.public=1
			or org.uid=?
			or org.id in (select urg.org_id from t_user_org urg where urg.uid=?)
			)
			`
			gen.Args = append(gen.Args, lgusr.Id)
			gen.Args = append(gen.Args, lgusr.Id)
		}
		if q != "" {
			gen.SQL += "\nAND org.name like ? "
			gen.Args = append(gen.Args, "%"+q+"%")
		}
		gen.SQL += "\nORDER BY org.aid DESC"
		page, err = comm.FindPages(gen, &ls, pg, 10)
	}
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	for _, v := range ls {
		usr, ok := service.GetUser(v.Uid)
		if ok {
			v.Nick = usr.Nick
			v.Avat = usr.Avatar
		}
		v.Pipeln, _ = comm.Db.Where("org_id=?", v.Id).Count(model.TOrgPipe{})
		v.Userln, _ = comm.Db.Where("org_id=?", v.Id).Count(model.TUserOrg{})
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
	lgusr := service.GetMidLgUser(c)
	if !service.IsAdmin(lgusr) {
		uf, ok := service.GetUserInfo(lgusr.Id)
		if !ok || uf.PermOrg != 1 {
			c.String(405, "no permission")
			return
		}
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
	perm := service.NewOrgPerm(service.GetMidLgUser(c), org.Id)
	if !perm.CanRead() {
		c.String(405, "no permission")
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
		"perm": hbtp.Map{
			"adm":   perm.IsOrgAdmin(),
			"own":   perm.IsOrgOwner(),
			"read":  perm.CanRead(),
			"write": perm.CanWrite(),
			"exec":  perm.CanExec(),
		},
	})
}
func (OrgController) users(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	if id == "" {
		c.String(500, "param err")
		return
	}
	perm := service.NewOrgPerm(service.GetMidLgUser(c), id)
	if !perm.CanRead() {
		c.String(405, "no permission")
		return
	}
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	var usrs []*models.TUserOrgInfo
	if comm.IsMySQL {
		ses := comm.Db.SQL(`
		select usr.*,urg.perm_adm,urg.perm_rw,urg.perm_exec,urg.perm_down,urg.created as join_time from t_user usr
		JOIN t_user_org urg ON urg.org_id=?
		where usr.id=urg.uid
		ORDER BY urg.created ASC
		`, perm.Org().Id)
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
	perm := service.NewOrgPerm(service.GetMidLgUser(c), id)
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	if !perm.IsOrgAdmin() {
		c.String(405, "no permission")
		return
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
		Where("id=?", perm.Org().Id).Update(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.JSON(200, &bean.IdsRes{
		Id:  ne.Id,
		Aid: ne.Aid,
	})
}
func (OrgController) rm(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	perm := service.NewOrgPerm(service.GetMidLgUser(c), id)
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	if !perm.IsOrgAdmin() {
		c.String(405, "no permission")
		return
	}
	ne := &model.TOrg{
		Deleted:     1,
		DeletedTime: time.Now(),
		Updated:     time.Now(),
	}
	_, err := comm.Db.Cols("deleted", "deleted_time", "updated").
		Where("id=?", perm.Org().Id).Update(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, "ok")
}
func (OrgController) userEdit(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	uid := m.GetString("uid")
	adm := m.GetBool("adm")
	rw := m.GetBool("rw")
	ex := m.GetBool("ex")
	dw := m.GetBool("dw")
	isadd := m.GetBool("add")
	perm := service.NewOrgPerm(service.GetMidLgUser(c), id)
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	usr := &models.TUser{}
	ok := service.GetIdOrAid(uid, usr)
	if !ok {
		c.String(404, "not found user")
		return
	}
	var err error
	ne := &model.TUserOrg{}
	isup, _ := comm.Db.Where("uid=? and org_id=?", usr.Id, perm.Org().Id).Get(ne)
	if usr.Id == perm.LgUser().Id {
		c.String(511, "can't edit yourself")
		return
	}
	if !perm.IsAdmin() {
		if adm {
			if !perm.IsOrgOwner() {
				c.String(405, "no permission")
				return
			}
		} else {
			if !perm.IsOrgAdmin() {
				c.String(405, "no permission")
				return
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
		if dw {
			ne.PermDown = 1
		} else {
			ne.PermDown = 0
		}
	}
	if isup {
		_, err = comm.Db.Cols("perm_adm", "perm_rw", "perm_exec", "perm_down").
			Where("aid=?", ne.Aid).Update(ne)
	} else {
		ne.Uid = usr.Id
		ne.OrgId = perm.Org().Id
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
	perm := service.NewOrgPerm(service.GetMidLgUser(c), id)
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	usr := &models.TUser{}
	ok := service.GetIdOrAid(uid, usr)
	if !ok {
		c.String(404, "not found user")
		return
	}
	ne := &model.TUserOrg{}
	ok, _ = comm.Db.Where("uid=? and org_id=?", usr.Id, perm.Org().Id).Get(ne)
	if !ok {
		c.String(404, "not found user org")
		return
	}
	if usr.Id == perm.LgUser().Id {
		c.String(511, "can't remove yourself")
		return
	}
	if !perm.IsOrgAdmin() {
		c.String(405, "no permission")
		return
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
	perm := service.NewOrgPerm(service.GetMidLgUser(c), id)
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	if !perm.IsOrgAdmin() {
		c.String(405, "no permission")
		return
	}
	ne := &model.TOrgPipe{}
	ok, _ := comm.Db.Where("org_id=? and pipe_id=?", perm.Org().Id, pipeId).Get(ne)
	if ok {
		c.String(511, "pipeline exist")
		return
	}
	ne.OrgId = perm.Org().Id
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
	perm := service.NewOrgPerm(service.GetMidLgUser(c), id)
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	if !perm.IsOrgAdmin() {
		c.String(405, "no permission")
		return
	}
	ne := &model.TOrgPipe{}
	_, err := comm.Db.Where("org_id=? and pipe_id=?", perm.Org().Id, pipeId).Delete(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, fmt.Sprintf("%d", ne.Aid))
}
