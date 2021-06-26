package route

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/core/utils"
	"github.com/gokins-main/gokins/bean"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/service"
	"github.com/gokins-main/gokins/util"
	hbtp "github.com/mgr9525/HyperByte-Transfer-Protocol"
	"gopkg.in/yaml.v3"
	"net/http"
	"time"
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
	g.POST("/info", util.GinReqParseJson(c.info))
}
func (PipelineController) orgPipelines(c *gin.Context, m *hbtp.Map) {
	orgId := m.GetString("orgId")
	q := m.GetString("q")
	pg, _ := m.GetInt("page")
	if orgId == "" {
		c.String(500, "param err")
		return
	}
	usr := service.GetMidLgUser(c)
	ls := make([]*model.TPipeline, 0)
	var err error
	var page *bean.Page
	if comm.IsMySQL {
		gen := &bean.PageGen{
			CountCols: "top.pipe_id",
			FindCols:  "pipe.*",
		}
		if usr.Id == "admin" {
			gen.SQL = `
			select {{select}} from t_pipeline pipe 
			LEFT JOIN t_org_pipe top on pipe.id = top.pipe_id 
			where top.org_id = ? 
		    `
			gen.Args = append(gen.Args, orgId)
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
			c.JSON(http.StatusOK, page)
			return
		}

		usero := &model.TUserOrg{}
		get, err := comm.Db.Where("uid =? and org_id =?", usr.Id, orgId).Get(usero)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
		if !get {
			c.JSON(http.StatusOK, gin.H{
				"data": ls,
			})
			return
		}
		if usero.PermAdm != 0 {
			gen.SQL = `
			select {{select}} from t_pipeline pipe 
			LEFT JOIN t_org_pipe top on pipe.id = top.pipe_id 
			where top.org_id = ?  
		    `
		} else if usero.PermRw != 0 {
			gen.SQL = `
			select {{select}} from t_pipeline pipe
			LEFT JOIN t_org_pipe top on pipe.id = top.pipe_id
			where (top.org_id = ?  or t_org_pipe.public != 0 )
		    `
		} else {
			c.JSON(http.StatusOK, page)
			return
		}

		gen.Args = append(gen.Args, orgId)
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
		if q != "" {
			session.Where("name = ?", q)
		}
		session.Desc("id")
		if usr.Id != "admin" {
			session.Where("create_user_id = ?", usr.Id)
		}
		page, err = comm.FindPage(session, &ls, pg)
		if err != nil {
			c.String(500, "db err:"+err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, page)
}

func (PipelineController) new(c *gin.Context, m *hbtp.Map) {
	name := m.GetString("name")
	content := m.GetString("content")
	orgId := m.GetString("orgId")
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
	json, err := y.ToJson()
	if err != nil {
		c.String(500, "yaml ToJson err:"+err.Error())
		return
	}
	usr := service.GetMidLgUser(c)
	pipeline := &model.TPipeline{
		Id:           utils.NewXid(),
		Name:         name,
		DisplayName:  y.DisplayName,
		PipelineType: "",
		JsonContent:  string(json),
		CreateUserId: usr.Id,
	}
	_, err = comm.Db.InsertOne(pipeline)
	if err != nil {
		c.String(500, "db err:"+err.Error())
		return
	}

	if orgId != "" {
		top := &model.TOrgPipe{
			OrgId:   orgId,
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
	c.JSON(http.StatusOK, "ok")
}

func (PipelineController) info(c *gin.Context, m *hbtp.Map) {
	id := m.GetString("id")
	if id == "" {
		c.String(500, "param err")
		return
	}
	pipe := &model.TPipeline{}
	ok, _ := comm.Db.Where("id=?", id).Get(pipe)
	if !ok {
		c.String(404, "not found org")
		return
	}
	c.JSON(200, pipe)
}
