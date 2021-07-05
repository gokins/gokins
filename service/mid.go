package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/gokins/model"
)

const LgUserKey = "lguser"

func MidUserCheck(c *gin.Context) {
	usr, ok := CurrUserCache(c)
	if !ok || (!IsAdmin(usr) && usr.Active != 1) {
		c.String(403, "Not Auth")
		c.Abort()
	}
	c.Set(LgUserKey, usr)
	c.Next()
}
func GetMidLgUser(c *gin.Context) *model.TUser {
	usr, ok := c.Get(LgUserKey)
	if !ok {
		return nil
	}
	lguser, ok := usr.(*model.TUser)
	if !ok {
		return nil
	}
	return lguser
}
