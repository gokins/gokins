package route

import (
	"github.com/gin-gonic/gin"
)

type ApiController struct{}

func (ApiController) GetPath() string {
	return "/api/"
}
func (c *ApiController) Routes(g gin.IRoutes) {
	g.POST("/url", c.test)
}
func (ApiController) test(c *gin.Context) {
}
