package route

import (
	"fmt"
	"github.com/gokins/core/common"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokins/core/utils"
	"github.com/gokins/gokins/bean"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/model"
	"github.com/gokins/gokins/models"
	"github.com/gokins/gokins/service"
	"github.com/gokins/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"github.com/sirupsen/logrus"
)

type ArtifactController struct{}

func (ArtifactController) GetPath() string {
	return "/api/art"
}
func (c *ArtifactController) Routes(g gin.IRoutes) {
	g.Use(service.MidUserCheck)
	g.POST("/org-list", util.GinReqParseJson(c.orgList))
	g.POST("/info", util.GinReqParseJson(c.info))
	g.POST("/edit", util.GinReqParseJson(c.edit))
	g.POST("/rm", util.GinReqParseJson(c.rm))
	g.POST("/package/list", util.GinReqParseJson(c.packageList))
	g.POST("/version/list", util.GinReqParseJson(c.versionList))
	g.POST("/version/infos", util.GinReqParseJson(c.versionInfos))
	g.POST("/version/url", util.GinReqParseJson(c.versionUrl))
	g.POST("/version/save", util.GinReqParseJson(c.versionSave))
	g.POST("/version/rm", util.GinReqParseJson(c.versionRm))
}
func (ArtifactController) orgList(c *gin.Context, m *hbtp.Map) {
	orgId := m.GetString("orgId")
	q := m.GetString("q")
	pg, _ := m.GetInt("page")
	if orgId == "" {
		c.String(500, "param err")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewOrgPerm(lgusr, orgId)
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	if !perm.CanRead() {
		c.String(405, "No Auth")
		return
	}
	ls := make([]*models.TArtifactory, 0)
	var err error
	var page *bean.Page
	if comm.IsMySQL {
		gen := &bean.PageGen{
			CountCols: "art.aid",
			FindCols:  "art.*",
		}
		gen.SQL = `
			select {{select}} from t_artifactory art 
			where art.deleted != 1 and art.org_id=?
		    `
		gen.Args = append(gen.Args, perm.Org().Id)
		if q != "" {
			gen.SQL += "\nAND art.name like ? "
			gen.Args = append(gen.Args, "%"+q+"%")
		}
		gen.SQL += "\nORDER BY art.aid DESC"
		page, err = comm.FindPages(gen, &ls, pg, 20)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
	}
	for _, v := range ls {
		usr, ok := service.GetUser(v.Uid)
		if ok {
			v.Nick = usr.Nick
			v.Avat = usr.Avatar
		}
		e := &model.TArtifactPackage{}
		v.Artln, _ = comm.Db.Where("repo_id=?", v.Id).Count(e)
	}
	c.JSON(http.StatusOK, page)
}
func (ArtifactController) info(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	arty := &model.TArtifactory{}
	ok := service.GetIdOrAid(id, arty)
	if !ok || arty.Deleted == 1 {
		c.String(404, "not found art")
		return
	}
	perm := service.NewOrgPerm(service.GetMidLgUser(c), arty.OrgId)
	if !perm.CanRead() {
		c.String(405, "no permission")
		return
	}
	usr := &models.TUser{}
	ok = service.GetIdOrAid(arty.Uid, usr)
	if !ok {
		c.String(404, "not found user?")
		return
	}
	c.JSON(200, hbtp.Map{
		"arty": arty,
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
func (ArtifactController) edit(c *gin.Context, m *hbtp.Map) {
	orgId := m.GetString("orgId")
	id := m.GetString("id")
	name := strings.TrimSpace(m.GetString("name"))
	desc := strings.TrimSpace(m.GetString("desc"))
	disabled := m.GetBool("disabled")
	if name == "" {
		c.String(500, "param err")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewOrgPerm(lgusr, orgId)
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	if !perm.CanWrite() {
		c.String(405, "No Permission")
		return
	}
	var err error
	ne := &model.TArtifactory{}
	isup := service.GetIdOrAid(id, ne)
	ne.Name = name
	ne.Desc = desc
	if disabled {
		ne.Disabled = 1
	} else {
		ne.Disabled = 0
	}
	ne.Updated = time.Now()
	if isup {
		if ne.OrgId != perm.Org().Id {
			c.String(405, "No Permission")
			return
		}
		_, err = comm.Db.Cols("name", "desc", "disabled", "updated").Where("id=?", ne.Id).Update(ne)
	} else {
		ne.Id = utils.NewXid()
		ne.Uid = lgusr.Id
		ne.OrgId = perm.Org().Id
		ne.Created = time.Now()

		ln := 0
		ne.Identifier = strings.ToLower(utils.RandomString(8))
		for !hbtp.EndContext(c) {
			ln++
			n, _ := comm.Db.Where("identifier=?", ne.Identifier).Count(ne)
			if n <= 0 {
				break
			}
			i := 8
			if ln >= 9 {
				i = 11
			} else if ln >= 6 {
				i = 10
			} else if ln >= 3 {
				i = 9
			}
			ne.Identifier = strings.ToLower(utils.RandomString(i))
		}
		_, err = comm.Db.InsertOne(ne)
	}
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, ne.Id)
}
func (ArtifactController) rm(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	art := &model.TArtifactory{}
	ok := service.GetIdOrAid(id, art)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewOrgPerm(lgusr, art.OrgId)
	if perm.Org() == nil || perm.Org().Deleted == 1 {
		c.String(404, "not found org")
		return
	}
	if !perm.CanWrite() {
		c.String(405, "No Permission")
		return
	}
	art.Deleted = 1
	art.DeletedTime = time.Now()
	art.Updated = time.Now()
	_, err := comm.Db.Cols("deleted", "deleted_time", "updated").Where("id=?", art.Id).Update(art)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, art.Id)
}
func (ArtifactController) packageList(c *gin.Context, m *hbtp.Map) {
	repoId := m.GetString("repoId")
	if repoId == "" {
		c.String(500, "param err")
		return
	}
	q := m.GetString("q")
	pg, _ := m.GetInt("page")
	var ls []*models.TArtifactPackage
	ses := comm.Db.Where("deleted!=1 and repo_id=?", repoId).OrderBy("aid DESC")
	if q != "" {
		qs := "%" + q + "%"
		ses.And("name like ? or display_name like ?", qs, qs)
	}
	page, err := comm.FindPage(ses, &ls, pg, 20)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	for _, v := range ls {
		/*usr, ok := service.GetUser(v.Uid)
		if ok {
			v.Nick = usr.Nick
			v.Avat = usr.Avatar
		}*/
		e := &model.TArtifactVersion{}
		v.Verln, _ = comm.Db.Where("package_id=?", v.Id).Count(e)
	}
	c.JSON(200, page)
}
func (ArtifactController) versionList(c *gin.Context, m *hbtp.Map) {
	packId := m.GetString("packId")
	if packId == "" {
		c.String(500, "param err")
		return
	}
	q := m.GetString("q")
	pg, _ := m.GetInt("page")
	var ls []*model.TArtifactVersion
	ses := comm.Db.Where("package_id=?", packId).OrderBy("aid DESC")
	if q != "" {
		qs := "%" + q + "%"
		ses.And("name like ? or display_name like ?", qs, qs)
	}
	page, err := comm.FindPage(ses, &ls, pg, 10)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	/*for _, v := range ls {
		usr, ok := service.GetUser(v.Uid)
		if ok {
			v.Nick = usr.Nick
			v.Avat = usr.Avatar
		}
		e := &model.TArtifactVersion{}
		v.Verln, _ = comm.Db.Where("package_id=?", v.Id).Count(e)
	}*/
	c.JSON(200, page)
}
func (ArtifactController) versionInfos(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	artv := &models.TArtifactVersion{}
	ok := service.GetIdOrAid(id, artv)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	arty := &model.TArtifactory{}
	ok = service.GetIdOrAid(artv.RepoId, arty)
	if !ok || arty.Deleted == 1 {
		c.String(404, "Not Found repo")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewOrgPerm(lgusr, arty.OrgId)
	if !perm.CanRead() {
		c.String(405, "No Permission")
		return
	}
	err := artv.ReadFiles()
	if err != nil {
		//c.String(511, "Files is err")
		//return
		logrus.Debugf("files err:%v", err)
	}
	c.JSON(200, hbtp.Map{
		"info": artv,
		"perm": hbtp.Map{
			"read":  perm.CanRead(),
			"write": perm.CanWrite(),
		},
	})
}
func (ArtifactController) versionUrl(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	pth := m.GetString("path")
	artv := &models.TArtifactVersion{}
	ok := service.GetIdOrAid(id, artv)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	arty := &model.TArtifactory{}
	ok = service.GetIdOrAid(artv.RepoId, arty)
	if !ok || arty.Deleted == 1 {
		c.String(404, "Not Found repo")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewOrgPerm(lgusr, arty.OrgId)
	if !perm.CanDownload() {
		c.String(405, "No Permission")
		return
	}

	tms := time.Now().Format(time.RFC3339Nano)
	random := utils.RandomString(20)
	sign := utils.Md5String(artv.Id + tms + random + comm.Cfg.Server.DownToken)
	ul := fmt.Sprintf("%s/api/art/pub/down/%s/%s?times=%s&random=%s&sign=%s",
		comm.Cfg.Server.Host, artv.Id, pth, url.QueryEscape(tms), random, sign)
	c.JSON(200, hbtp.Map{
		"id":     artv.Id,
		"times":  tms,
		"random": random,
		"sign":   sign,
		"url":    ul,
	})
}
func (ArtifactController) versionSave(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	artv := &models.TArtifactVersion{}
	ok := service.GetIdOrAid(id, artv)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	arty := &model.TArtifactory{}
	ok = service.GetIdOrAid(artv.RepoId, arty)
	if !ok || arty.Deleted == 1 {
		c.String(404, "Not Found repo")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewOrgPerm(lgusr, arty.OrgId)
	if !perm.CanWrite() {
		c.String(405, "No Permission")
		return
	}
	artv.Version = strings.TrimSpace(m.GetString("version"))
	artv.Desc = strings.TrimSpace(m.GetString("desc"))
	if m.GetBool("ispre") {
		artv.Preview = 1
	} else {
		artv.Preview = 0
	}
	artv.Updated = time.Now()
	_, err := comm.Db.Cols("version", "desc", "preview", "updated").Where("id=?", artv.Id).Update(artv)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, artv.Id)
}
func (ArtifactController) versionRm(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	artv := &models.TArtifactVersion{}
	ok := service.GetIdOrAid(id, artv)
	if !ok {
		c.String(404, "Not Found")
		return
	}
	arty := &model.TArtifactory{}
	ok = service.GetIdOrAid(artv.RepoId, arty)
	if !ok || arty.Deleted == 1 {
		c.String(404, "Not Found repo")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewOrgPerm(lgusr, arty.OrgId)
	if !perm.CanWrite() {
		c.String(405, "No Permission")
		return
	}
	_, err := comm.Db.Where("id=?", artv.Id).Delete(artv)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	fls := filepath.Join(comm.WorkPath, common.PathArtifacts, artv.Id)
	os.RemoveAll(fls)
	c.String(200, artv.Id)
}
