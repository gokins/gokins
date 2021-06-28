package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/model"
	"github.com/gokins-main/gokins/util"
	"github.com/sirupsen/logrus"
)

func GetUser(uid string) (*model.TUser, bool) {
	e := &model.TUser{}
	ok, err := comm.Db.Where("id=?", uid).Get(e)
	if err != nil {
		logrus.Errorf("GetUser(%s) err:%v", uid, err)
	}
	return e, ok
}
func FindUserName(name string) (*model.TUser, bool) {
	e := &model.TUser{}
	ok, err := comm.Db.Where("name=?", name).Get(e)
	if err != nil {
		logrus.Errorf("FindUserName(%s) err:%v", name, err)
	}
	return e, ok
}

func GetUserCache(uid string) (*model.TUser, bool) {
	var ok bool
	e := &model.TUser{}
	uids := fmt.Sprintf("user:%s", uid)
	err := comm.CacheGets(uids, e)
	if err == nil {
		return e, true
	}
	e, ok = GetUser(uid)
	if ok {
		comm.CacheSets(uids, e)
	}
	return e, ok
}
func CurrUserCache(c *gin.Context) (*model.TUser, bool) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	tk := util.GetToken(c, comm.Cfg.Server.LoginKey)
	if tk == nil {
		return nil, false
	}
	uid, ok := tk["uid"]
	if !ok {
		return nil, false
	}
	uids, ok := uid.(string)
	if !ok {
		return nil, false
	}
	return GetUserCache(uids)
}
func IsAdmin(usr *model.TUser) bool {
	return usr.Id == "admin"
}
func IsOrgAdmin(uid, orgId string) bool {
	usero, ok := GetUserOrg(uid, orgId)
	if !ok {
		return false
	}
	return usero.PermAdm != 0
}
func GetUsePermRwr(uid, orgId string) int {
	usero, ok := GetUserOrg(uid, orgId)
	if !ok {
		return 0
	}
	return usero.PermRw
}
func HasOrgExec(uid, orgId string) bool {
	usero, ok := GetUserOrg(uid, orgId)
	if !ok {
		return false
	}
	return usero.PermExec != 0
}
func GetUserOrg(uid, orgId string) (*model.TUserOrg, bool) {
	torg := &model.TOrg{}
	ok := GetIdOrAid(orgId, torg)
	if !ok {
		return nil, false
	}
	usero := &model.TUserOrg{}
	get, err := comm.Db.Where("uid =? and org_id =?", uid, torg.Id).Get(usero)
	if err != nil {
		logrus.Debugf("HasOrgExec db err:%v", err)
	}
	if !get {
		return nil, false
	}
	return usero, true
}

type OrgPerm struct {
	lgusr  *model.TUser
	org    *model.TOrg
	usrOrg *model.TUserOrg
}

func NewOrgPerm(lgusr *model.TUser, orgId string) *OrgPerm {
	c := &OrgPerm{lgusr: lgusr}
	org := &model.TOrg{}
	ok := GetIdOrAid(orgId, org)
	if ok && org.Deleted != 1 {
		c.org = org
		usero := &model.TUserOrg{}
		if lgusr != nil {
			ok, _ = comm.Db.Where("uid =? and org_id =?", lgusr.Id, org.Id).Get(usero)
			if ok {
				c.usrOrg = usero
			}
		}
	}
	return c
}
func (c *OrgPerm) IsAdmin() bool {
	if c.lgusr != nil && IsAdmin(c.lgusr) {
		return true
	}
	return false
}
func (c *OrgPerm) IsOrgOwner() bool {
	if c.org != nil && c.lgusr != nil && c.org.Uid == c.lgusr.Id {
		return true
	}
	return false
}
func (c *OrgPerm) IsOrgPublic() bool {
	if c.org != nil && c.org.Public == 1 {
		return true
	}
	return false
}
func (c *OrgPerm) IsOrgAdmin() bool {
	if c.usrOrg != nil && c.usrOrg.PermAdm == 1 {
		return true
	}
	return false
}
func (c *OrgPerm) CanOrgRw() bool {
	if c.usrOrg != nil && c.usrOrg.PermRw == 1 {
		return true
	}
	return false
}
func (c *OrgPerm) CanOrgExec() bool {
	if c.usrOrg != nil && c.usrOrg.PermExec == 1 {
		return true
	}
	return false
}

//LgUser maybe null
func (c *OrgPerm) LgUser() *model.TUser {
	return c.lgusr
}

//Org maybe null
func (c *OrgPerm) Org() *model.TOrg {
	return c.org
}

//UserOrg maybe null
func (c *OrgPerm) UserOrg() *model.TUserOrg {
	return c.usrOrg
}
