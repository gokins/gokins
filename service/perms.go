package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gokins/gokins/model"
)

func CheckPermission(uid string, perms string) bool {
	usr, ok := GetUser(uid)
	if !ok {
		return false
	}
	return CheckUPermission(usr, perms)
}
func CheckUPermission(usr *model.TUser, perms string) bool {
	if usr == nil {
		return false
	}
	if perms == "common" {
		return true
	} else if perms == "admin" && usr.Name == "admin" {
		return true
	}
	return false
}
func CheckCurrPermission(c *gin.Context, perms string) bool {
	usr, ok := CurrUserCache(c)
	if !ok {
		return false
	}
	return CheckUPermission(usr, perms)
}
