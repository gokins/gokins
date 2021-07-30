package route

import (
	"fmt"
	"github.com/gokins/gokins/engine"
	"github.com/gokins/gokins/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokins/core/utils"
	"github.com/gokins/gokins/bean"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/model"
	"github.com/gokins/gokins/service"
	"github.com/gokins/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"gopkg.in/yaml.v3"
)

type PipelineController struct{}

func (PipelineController) GetPath() string {
	return "/api/pipeline"
}
func (c *PipelineController) Routes(g gin.IRoutes) {
	g.Use(service.MidUserCheck)
	g.POST("/org/pipelines", util.GinReqParseJson(c.orgPipelines))
	g.POST("/pipelines", util.GinReqParseJson(c.getPipelines))
	g.POST("/new", util.GinReqParseJson(c.new))
	g.POST("/delete", util.GinReqParseJson(c.delete))
	g.POST("/info", util.GinReqParseJson(c.info))
	g.POST("/save", util.GinReqParseJson(c.save))
	g.POST("/run", util.GinReqParseJson(c.run))
	g.POST("/copy", util.GinReqParseJson(c.copy))
	g.POST("/rebuild", util.GinReqParseJson(c.rebuild))
	g.POST("/pipelineVersions", util.GinReqParseJson(c.pipelineVersions))
	g.POST("/pipelineVersion", util.GinReqParseJson(c.pipelineVersion))
	g.POST("/search/sha", util.GinReqParseJson(c.searchSha))
	g.POST("/vars", util.GinReqParseJson(c.vars))
	g.POST("/var/save", util.GinReqParseJson(c.varSave))
	g.POST("/var/del", util.GinReqParseJson(c.varDel))
}
func (PipelineController) orgPipelines(c *gin.Context, m *hbtp.Map) {
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
	ls := make([]*models.TPipeline, 0)
	var err error
	var page *bean.Page
	if comm.IsMySQL {
		gen := &bean.PageGen{
			CountCols: "DISTINCT(pipe.id),pipe.id",
			FindCols:  "DISTINCT(pipe.id),pipe.*",
		}
		gen.SQL = `
			select {{select}} from t_pipeline pipe 
			LEFT JOIN t_org_pipe top on pipe.id = top.pipe_id 
			where top.org_id = ? and pipe.deleted != 1
		    `
		gen.Args = append(gen.Args, perm.Org().Id)
		if q != "" {
			gen.SQL += "\nAND pipe.name like ? "
			gen.Args = append(gen.Args, "%"+q+"%")
		}
		gen.SQL += "\nORDER BY pipe.id DESC"
		page, err = comm.FindPages(gen, &ls, pg, 10)
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
		last := &models.RunBuild{}
		v.Buildln, _ = comm.Db.Where("pipeline_id=?", v.Id).Count(last)
		if v.Buildln > 0 {
			ok, _ = comm.Db.Where("pipeline_id=?", v.Id).OrderBy("created DESC").Get(last)
			if ok {
				v.Build = last
			}
		}
	}
	c.JSON(http.StatusOK, page)
}
func (PipelineController) getPipelines(c *gin.Context, m *hbtp.Map) {
	q := m.GetString("q")
	pg, _ := m.GetInt("page")
	lgusr := service.GetMidLgUser(c)
	ls := make([]*models.TPipeline, 0)
	var err error
	var page *bean.Page
	if comm.IsMySQL {
		gen := &bean.PageGen{
			CountCols: "pipe.id",
			FindCols:  "pipe.*",
		}
		gen.SQL = `
			select {{select}} from t_pipeline pipe where pipe.deleted != 1 `
		if !service.IsAdmin(lgusr) {
			gen.SQL = gen.SQL + ` and pipe.uid = ? `
			gen.Args = append(gen.Args, lgusr.Id)
		}
		if q != "" {
			gen.SQL += "\nAND pipe.name like ? "
			gen.Args = append(gen.Args, "%"+q+"%")
		}
		gen.SQL += "\nORDER BY pipe.id DESC"
		page, err = comm.FindPages(gen, &ls, pg, 10)
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
		last := &models.RunBuild{}
		v.Buildln, _ = comm.Db.Where("pipeline_id=?", v.Id).Count(last)
		if v.Buildln > 0 {
			ok, _ = comm.Db.Where("pipeline_id=?", v.Id).OrderBy("created DESC").Get(last)
			if ok {
				v.Build = last
			}
		}
	}
	c.JSON(http.StatusOK, page)
}

func (PipelineController) save(c *gin.Context, m *hbtp.Map) {
	name := m.GetString("name")
	content := m.GetString("content")
	pipelineId := m.GetString("pipelineId")
	accessToken := m.GetString("accessToken")
	ul := m.GetString("url")
	username := m.GetString("username")
	displayName := m.GetString("displayName")
	if pipelineId == "" {
		c.String(500, "param err")
		return
	}
	usr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(usr, pipelineId)
	if !perm.CanWrite() {
		c.String(405, "No Auth")
		return
	}
	y := &bean.Pipeline{}
	err := yaml.Unmarshal([]byte(content), y)
	err = y.Check()
	if err != nil {
		c.String(500, "yaml Check err:"+err.Error())
		return
	}
	pipeline := &model.TPipeline{
		Name:        name,
		DisplayName: displayName,
	}
	_, err = comm.Db.Cols("name,display_name").Where("id = ?", pipelineId).Update(pipeline)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	tpc := &model.TPipelineConf{
		YmlContent:  content,
		Url:         ul,
		Username:    username,
		AccessToken: accessToken,
	}
	_, err = comm.Db.Cols("yml_content,url,username,access_token").Where("pipeline_id = ?", pipelineId).Update(tpc)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}
func (PipelineController) delete(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	if id == "" {
		c.String(500, "param err")
		return
	}
	usr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(usr, id)
	if perm.Pipeline() == nil || perm.Pipeline().Deleted == 1 {
		c.String(404, "未找到流水线信息")
		return
	}
	if !perm.CanWrite() {
		c.String(405, "No Auth")
		return
	}
	tp := &model.TPipeline{
		Deleted:     1,
		DeletedTime: time.Now(),
	}
	_, err := comm.Db.Cols("deleted").Where("id = ?", id).Update(tp)
	if err != nil {
		c.String(500, "TPipeline Update db err:"+err.Error())
		return
	}
	version := &model.TPipelineVersion{
		Deleted: 1,
	}
	_, err = comm.Db.Cols("deleted").Where("pipeline_id = ?", id).Update(version)
	if err != nil {
		c.String(500, "TPipeline Update db err:"+err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}
func (PipelineController) new(c *gin.Context, npipe *bean.NewPipeline) {
	if !npipe.Check() {
		c.String(500, "param err")
		return
	}
	y := &bean.Pipeline{}
	err := yaml.Unmarshal([]byte(npipe.Content), y)
	if err != nil {
		c.String(500, "yaml Unmarshal err:"+err.Error())
		return
	}
	err = y.Check()
	if err != nil {
		c.String(500, "yaml Check err:"+err.Error())
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewOrgPerm(lgusr, npipe.OrgId)
	if npipe.OrgId != "" && perm.Org() == nil {
		c.String(404, "组织不存在")
		return
	}
	if !perm.IsAdmin() {
		uf, ok := service.GetUserInfo(lgusr.Id)
		if !ok || uf.PermPipe != 1 {
			c.String(405, "no permission")
			return
		}
		if perm.Org() != nil && !perm.CanWrite() {
			c.String(405, "No Auth")
			return
		}
	}
	pipeline := &model.TPipeline{
		Id:           utils.NewXid(),
		Uid:          lgusr.Id,
		Name:         npipe.Name,
		DisplayName:  npipe.DisplayName,
		PipelineType: "",
	}
	_, err = comm.Db.InsertOne(pipeline)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	tpc := &model.TPipelineConf{
		PipelineId:  pipeline.Id,
		YmlContent:  npipe.Content,
		Url:         npipe.Url,
		Username:    npipe.Username,
		AccessToken: npipe.AccessToken,
	}
	_, err = comm.Db.InsertOne(tpc)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	if npipe.Vars != nil && len(npipe.Vars) > 0 {
		for _, v := range npipe.Vars {
			pipelineVar := &model.TPipelineVar{}
			err = utils.Struct2Struct(pipelineVar, v)
			if err != nil {
				c.String(500, "model err:"+err.Error())
				return
			}
			pipelineVar.Uid = lgusr.Id
			pipelineVar.PipelineId = pipeline.Id
			if v.Public {
				pipelineVar.Public = 1
			}
			_, err = comm.Db.InsertOne(pipelineVar)
			if err != nil {
				c.String(500, "db err:"+err.Error())
				return
			}
		}
	}
	if perm.Org() != nil {
		top := &model.TOrgPipe{
			OrgId:   perm.Org().Id,
			PipeId:  pipeline.Id,
			Created: time.Now(),
			Public:  0,
		}
		_, err = comm.Db.InsertOne(top)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}

	}
	c.JSON(http.StatusOK, pipeline)
}

func (PipelineController) info(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	if id == "" {
		c.String(500, "param err")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(lgusr, id)
	if perm.Pipeline() == nil || perm.Pipeline().Deleted == 1 {
		c.String(404, "未找到流水线信息")
		return
	}
	if !perm.CanRead() {
		c.String(405, "No Auth")
		return
	}
	pipe := &models.TPipelineInfo{}
	ok, _ := comm.Db.Where("id=? and deleted != 1", id).Get(pipe)
	if !ok {
		c.String(404, "未找到流水线信息")
		return
	}
	tpc := &model.TPipelineConf{}
	_, err := comm.Db.Where("pipeline_id=?", pipe.Id).Get(tpc)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	pipe.YmlContent = tpc.YmlContent
	pipe.Url = tpc.Url
	s := "***"
	if perm.CanWrite() {
		pipe.Username = tpc.Username
		pipe.AccessToken = tpc.AccessToken
	} else {
		pipe.Username = s
		pipe.AccessToken = s
	}
	c.JSON(200, hbtp.Map{
		"pipe": pipe,
		"perm": hbtp.Map{
			"read":  perm.CanRead(),
			"write": perm.CanWrite(),
			"exec":  perm.CanExec(),
		},
	})
}

func (PipelineController) run(c *gin.Context, m *hbtp.Map) {
	pipelineId := m.GetString("pipelineId")
	sha := m.GetString("sha")
	if pipelineId == "" {
		c.String(500, "param err")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(lgusr, pipelineId)
	if perm.Pipeline() == nil || perm.Pipeline().Deleted == 1 {
		c.String(404, "未找到流水线信息")
		return
	}
	if !perm.CanExec() {
		c.String(405, "No Auth")
		return
	}
	tvp, rb, err := service.Run(lgusr.Id, pipelineId, sha, "run")
	if err != nil {
		c.String(500, err.Error())
		return
	}
	engine.Mgr.BuildEgn().Put(rb)
	c.JSON(200, tvp)
}

func (PipelineController) copy(c *gin.Context, m *hbtp.Map) {
	pipelineId := m.GetString("pipelineId")
	if pipelineId == "" {
		c.String(500, "param err")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(lgusr, pipelineId)
	if perm.Pipeline() == nil || perm.Pipeline().Deleted == 1 {
		c.String(404, "未找到流水线信息")
		return
	}
	if !perm.CanRead() {
		c.String(405, "No Auth")
		return
	}
	if !perm.IsAdmin() {
		uf, ok := service.GetUserInfo(lgusr.Id)
		if !ok || uf.PermPipe != 1 {
			c.String(405, "no permission")
			return
		}
	}
	pipe := &model.TPipeline{
		Id:           utils.NewXid(),
		Uid:          lgusr.Id,
		Name:         fmt.Sprintf("%s_copy", perm.Pipeline().Name),
		DisplayName:  perm.Pipeline().DisplayName,
		PipelineType: perm.Pipeline().PipelineType,
	}
	_, err := comm.Db.InsertOne(pipe)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}

	tpc := &model.TPipelineConf{}
	_, err = comm.Db.Where("pipeline_id=?", perm.Pipeline().Id).Get(tpc)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	ne := &model.TPipelineConf{
		PipelineId:  pipe.Id,
		Url:         tpc.Url,
		AccessToken: tpc.AccessToken,
		YmlContent:  tpc.YmlContent,
		Username:    tpc.Username,
	}
	_, err = comm.Db.InsertOne(ne)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.JSON(200, pipe)
}
func (PipelineController) rebuild(c *gin.Context, m *hbtp.Map) {
	pipelineVersionId := m.GetString("pipelineVersionId")
	if pipelineVersionId == "" {
		c.String(500, "param err")
		return
	}
	tvp := &model.TPipelineVersion{}
	ok, _ := comm.Db.Where("id=? and deleted != 1", pipelineVersionId).Get(tvp)
	if !ok {
		c.String(404, "构建记录不存在")
		return
	}
	lgusr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(lgusr, tvp.PipelineId)
	if perm.Pipeline() == nil || perm.Pipeline().Deleted == 1 {
		c.String(404, "未找到流水线信息")
		return
	}
	if !perm.CanExec() {
		c.String(405, "No Permission")
		return
	}
	tvp, rb, err := service.ReBuild(lgusr.Id, tvp)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	engine.Mgr.BuildEgn().Put(rb)
	c.JSON(200, tvp)
}

func (PipelineController) pipelineVersions(c *gin.Context, m *hbtp.Map) {
	pipelineId := m.GetString("pipelineId")
	pg, _ := m.GetInt("page")
	usr := service.GetMidLgUser(c)
	ls := make([]*models.TPipelineVersion, 0)
	var page *bean.Page
	var err error
	if pipelineId != "" {
		perm := service.NewPipePerm(usr, pipelineId)
		if perm.Pipeline() == nil || perm.Pipeline().Deleted == 1 {
			c.String(404, "未找到流水线信息")
			return
		}
		if !perm.CanRead() {
			c.String(405, "No Auth")
			return
		}
		where := comm.Db.Where("pipeline_id = ? and deleted != 1", pipelineId).Desc("id")
		page, err = comm.FindPage(where, &ls, pg)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
	} else {
		if service.IsAdmin(usr) {
			where := comm.Db.Where(" deleted != 1").Desc("id")
			page, err = comm.FindPage(where, &ls, pg)
			if err != nil {
				c.String(500, "db err:"+err.Error())
				return
			}
		} else {
			tpipeIds := []string{}
			err = comm.Db.Table(&model.TPipeline{}).Cols("id").Where("uid = ? and deleted != 1", usr.Id).Find(&tpipeIds)
			if err != nil {
				c.String(500, "db err:"+err.Error())
				return
			}
			if len(tpipeIds) <= 0 {
				c.JSON(200, page)
				return
			}
			where := comm.Db.In("pipeline_id", tpipeIds).Where("deleted != 1").Desc("id")
			page, err = comm.FindPage(where, &ls, pg, 20)
			if err != nil {
				c.String(500, "db err:"+err.Error())
				return
			}
		}
	}

	for _, v := range ls {
		last := &models.RunBuild{}
		ok, _ := comm.Db.Where("pipeline_version_id=?", v.Id).OrderBy("created DESC").Get(last)
		if ok {
			v.Build = last
		}
	}

	c.JSON(200, page)

}
func (PipelineController) pipelineVersion(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	if id == "" {
		c.String(500, "param err")
		return
	}
	pv := &model.TPipelineVersion{}
	ok, _ := comm.Db.Where("id=?", id).Get(pv)
	if !ok {
		c.String(404, "not found pv")
		return
	}
	usr := &models.TUser{}
	service.GetIdOrAid(pv.Uid, usr)
	build := &models.RunBuild{}
	ok, _ = comm.Db.Where("pipeline_version_id=?", pv.Id).Get(build)
	if !ok {
		c.String(404, "not found build")
		return
	}
	perm := service.NewPipePerm(service.GetMidLgUser(c), pv.PipelineId)
	if perm.Pipeline() == nil {
		c.String(404, "not found pipe")
		return
	}
	if !perm.CanRead() {
		c.String(405, "no permission")
		return
	}

	pipeShow := &bean.PipelineShow{}
	err := utils.Struct2Struct(pipeShow, perm.Pipeline())
	if err != nil {
		c.String(405, "conv err:%v", err)
		return
	}
	pinfo := &model.TPipelineConf{}
	ok, _ = comm.Db.Where("pipeline_id=?", perm.Pipeline().Id).Get(pinfo)
	if ok {
		pipeShow.Url = pinfo.Url
	}
	c.JSON(200, hbtp.Map{
		"build": build,
		"pv":    pv,
		"usr":   usr,
		"pipe":  pipeShow,
		"perm": hbtp.Map{
			"read":  perm.CanRead(),
			"write": perm.CanWrite(),
			"exec":  perm.CanExec(),
		},
	})
}
func (PipelineController) searchSha(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	q := m.GetString("q")
	if id == "" {
		c.String(500, "param err")
		return
	}
	perm := service.NewPipePerm(service.GetMidLgUser(c), id)
	if perm.Pipeline() == nil {
		c.String(404, "not found pipe")
		return
	}
	if !perm.CanRead() {
		c.String(405, "no permission")
		return
	}
	shas := []string{}
	session := comm.Db.Table("t_pipeline_version").
		Distinct("sha").Cols("sha").
		Where("pipeline_id = ?", id).Desc("sha")
	if q != "" {
		session.And("sha like '%" + q + "%'")
	}
	err := session.Find(&shas)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	res := make([]map[string]string, 0)
	for _, sha := range shas {
		if sha == "" {
			continue
		}
		m2 := map[string]string{}
		m2["name"] = sha
		res = append(res, m2)
	}
	c.JSON(200, res)
}
func (PipelineController) vars(c *gin.Context, m *hbtp.Map) {
	pipelineId := m.GetString("pipelineId")
	q := m.GetString("q")
	pg, _ := m.GetInt("page")
	if pipelineId == "" {
		c.String(500, "param err")
		return
	}
	perm := service.NewPipePerm(service.GetMidLgUser(c), pipelineId)
	if perm.Pipeline() == nil {
		c.String(404, "not found pipe")
		return
	}
	if !perm.CanRead() {
		c.String(405, "no permission")
		return
	}
	ls := make([]*models.TPipelineVar, 0)
	var page *bean.Page
	var err error
	session := comm.Db.Where("pipeline_id = ?", pipelineId)
	if q != "" {
		session.And("(name like '%" + q + "%' or value like '%" + q + "'%)")
	}
	page, err = comm.FindPage(session, &ls, pg)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	if !perm.CanWrite() {
		lss := make([]*models.TPipelineVar, 0)
		for _, v := range ls {
			if v.Public != 0 {
				v.Value = "***"
			}
			lss = append(lss, v)
		}
		page.Data = lss
	}
	c.JSON(200, page)
}
func (PipelineController) varSave(c *gin.Context, pv *bean.PipelineVar) {
	if pv.Value == "" || pv.Name == "" || pv.PipelineId == "" {
		c.String(500, "param err")
		return
	}
	perm := service.NewPipePerm(service.GetMidLgUser(c), pv.PipelineId)
	if perm.Pipeline() == nil {
		c.String(404, "not found pipe")
		return
	}
	if !perm.CanWrite() {
		c.String(405, "no permission")
		return
	}
	pipelineVar := &model.TPipelineVar{}
	err := utils.Struct2Struct(pipelineVar, pv)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	if pv.Public {
		pipelineVar.Public = 1
	}
	tpv := &model.TPipelineVar{}
	ok, err := comm.Db.Where("pipeline_id = ? and name = ?", pv.PipelineId, pv.Name).Get(tpv)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	if pv.Aid > 0 {
		if ok && tpv.Aid != pv.Aid {
			c.String(500, "变量名重复")
			return
		}
		_, err = comm.Db.Cols("name,value,remarks,public").Where("aid = ?", pv.Aid).Update(pipelineVar)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
		c.String(200, "ok")
		return
	}
	if ok {
		c.String(500, "变量名重复")
		return
	}
	_, err = comm.Db.InsertOne(pipelineVar)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, "ok")
}
func (PipelineController) varDel(c *gin.Context, m *hbtp.Map) {
	aId, err := m.GetInt("aid")
	if err != nil || aId <= 0 {
		c.String(500, "param err")
		return
	}
	pipelineVar := &model.TPipelineVar{}
	ok, _ := comm.Db.Where("aid = ? ", aId).Get(pipelineVar)
	if !ok {
		c.String(404, "not found pipe_var")
		return
	}
	perm := service.NewPipePerm(service.GetMidLgUser(c), pipelineVar.PipelineId)
	if perm.Pipeline() == nil {
		c.String(404, "not found pipe")
		return
	}
	if !perm.CanWrite() {
		c.String(405, "no permission")
		return
	}
	_, err = comm.Db.Where("aid = ?", aId).Delete(pipelineVar)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(200, "ok")
}
