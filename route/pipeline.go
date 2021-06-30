package route

import (
	"github.com/gokins-main/gokins/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/service"
	"github.com/gokins-main/gokins/util"
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
	g.POST("/deleted", util.GinReqParseJson(c.deleted))
	g.POST("/info", util.GinReqParseJson(c.info))
	g.POST("/save", util.GinReqParseJson(c.save))
	g.POST("/run", util.GinReqParseJson(c.run))
	g.POST("/pipelineVersions", util.GinReqParseJson(c.pipelineVersions))
	g.POST("/pipelineVersion", util.GinReqParseJson(c.pipelineVersion))
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
	ls := make([]*model.TPipeline, 0)
	var err error
	var page *bean.Page
	if comm.IsMySQL {
		gen := &bean.PageGen{
			CountCols: "top.pipe_id",
			FindCols:  "pipe.*",
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
		page, err = comm.FindPages(gen, &ls, pg, 20)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, page)
}
func (PipelineController) getPipelines(c *gin.Context, m *hbtp.Map) {
	q := m.GetString("q")
	pg, _ := m.GetInt("page")
	usr := service.GetMidLgUser(c)
	ls := make([]*model.TPipeline, 0)
	var err error
	var page *bean.Page
	if comm.IsMySQL {
		session := comm.Db.NewSession()
		session.Where("deleted != 1")
		session.Desc("id")

		if !service.IsAdmin(usr) {
			session.And("uid = ?", usr.Id)
		}
		if q != "" {
			session.And("name = ?", q)
		}

		page, err = comm.FindPage(session, &ls, pg)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
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
	js, err := y.ToJson()
	if err != nil {
		c.String(500, "yaml ToJson err:"+err.Error())
		return
	}
	pipeline := &model.TPipeline{
		Name:        name,
		DisplayName: displayName,
		JsonContent: string(js),
		YmlContent:  content,
		Url:         ul,
		Username:    username,
		AccessToken: accessToken,
	}
	_, err = comm.Db.Cols("name , display_name,json_content,yml_content,url,username,access_token").Where("id = ?", pipelineId).Update(pipeline)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}
	c.String(http.StatusOK, "ok")
}
func (PipelineController) deleted(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	if id == "" {
		c.String(500, "param err")
		return
	}
	pipe := &model.TPipeline{}
	ok, _ := comm.Db.Where("id=? and deleted != 1", id).Get(pipe)
	if !ok {
		c.String(404, "未找到流水线信息")
		return
	}
	usr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(usr, id)
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
func (PipelineController) new(c *gin.Context, m *hbtp.Map) {
	name := m.GetString("name")
	content := m.GetString("content")
	orgId := m.GetString("orgId")
	accessToken := m.GetString("accessToken")
	ul := m.GetString("url")
	username := m.GetString("username")
	displayName := m.GetString("displayName")
	if name == "" || content == "" {
		c.String(500, "param err")
		return
	}
	y := &bean.Pipeline{}
	err := yaml.Unmarshal([]byte(content), y)
	if err != nil {
		c.String(500, "yaml Unmarshal err:"+err.Error())
		return
	}
	err = y.Check()
	if err != nil {
		c.String(500, "yaml Check err:"+err.Error())
		return
	}
	js, err := y.ToJson()
	if err != nil {
		c.String(500, "yaml ToJson err:"+err.Error())
		return
	}
	usr := service.GetMidLgUser(c)
	perm := service.NewOrgPerm(usr, orgId)
	if perm.Org() != nil && !perm.CanWrite() {
		c.String(405, "No Auth")
		return
	}
	pipeline := &model.TPipeline{
		Id:           utils.NewXid(),
		Uid:          usr.Id,
		Name:         name,
		DisplayName:  displayName,
		PipelineType: "",
		JsonContent:  string(js),
		YmlContent:   content,
		Url:          ul,
		Username:     username,
		AccessToken:  accessToken,
	}
	_, err = comm.Db.InsertOne(pipeline)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
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
	pipe := &model.TPipeline{}
	ok, _ := comm.Db.Where("id=? and deleted != 1", id).Get(pipe)
	if !ok {
		c.String(404, "未找到流水线信息")
		return
	}
	usr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(usr, id)
	if !perm.CanRead() {
		c.String(405, "No Auth")
		return
	}
	c.JSON(200, pipe)
}

func (PipelineController) run(c *gin.Context, m *hbtp.Map) {
	pipelineId := m.GetString("pipelineId")
	if pipelineId == "" {
		c.String(500, "param err")
		return
	}
	usr := service.GetMidLgUser(c)
	perm := service.NewPipePerm(usr, pipelineId)
	if !perm.CanExec() {
		c.String(405, "No Auth")
		return
	}
	tvp, err := service.Run(pipelineId)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	c.JSON(200, tvp)
}

func (PipelineController) pipelineVersions(c *gin.Context, m *hbtp.Map) {
	pipelineId := m.GetString("pipelineId")
	pg, _ := m.GetInt("page")
	usr := service.GetMidLgUser(c)
	ls := make([]*model.TPipelineVersion, 0)
	var page *bean.Page
	var err error
	if pipelineId != "" {
		pipe := &model.TPipeline{}
		ok, _ := comm.Db.Where("id=? and deleted != 1", pipelineId).Get(pipe)
		if !ok {
			c.String(404, "未找到流水线信息")
			return
		}
		where := comm.Db.Where("pipeline_id = ? and deleted != 1", pipelineId).Desc("id")
		page, err = comm.FindPage(where, &ls, pg)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
		perm := service.NewPipePerm(usr, pipelineId)
		if !perm.CanRead() {
			c.String(405, "No Auth")
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
			where := comm.Db.In("pipeline_id", tpipeIds).Desc("id")
			page, err = comm.FindPage(where, &ls, pg)
			if err != nil {
				c.String(500, "db err:"+err.Error())
				return
			}
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
	c.JSON(200, hbtp.Map{
		"build": build,
		"pv":    pv,
		"pipe":  pipeShow,
	})
}
