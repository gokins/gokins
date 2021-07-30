package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gokins/gokins/comm"
	"github.com/gokins/gokins/model"
	"github.com/gokins/gokins/util"
	"github.com/sirupsen/logrus"
)

func GetUser(uid string) (*model.TUser, bool) {
	if uid == "" {
		return nil, false
	}
	e := &model.TUser{}
	ok, err := comm.Db.Where("id=?", uid).Get(e)
	if err != nil {
		logrus.Errorf("GetUser(%s) err:%v", uid, err)
	}
	return e, ok
}
func GetUserInfo(uid string) (*model.TUserInfo, bool) {
	if uid == "" {
		return nil, false
	}
	e := &model.TUserInfo{Id: uid}
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

func ClearUserCache(uid string) {
	if uid == "" {
		return
	}
	uids := fmt.Sprintf("user:%s", uid)
	comm.CacheSet(uids, nil)
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
	ok := false
	if orgId != "" {
		ok = GetIdOrAid(orgId, org)
	}
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
	if c.IsAdmin() || c.IsOrgOwner() {
		return true
	}
	if c.usrOrg != nil && c.usrOrg.PermAdm == 1 {
		return true
	}
	return false
}
func (c *OrgPerm) CanRead() bool {
	if c.IsOrgPublic() || c.IsOrgAdmin() {
		return true
	}
	return c.usrOrg != nil
}
func (c *OrgPerm) CanWrite() bool {
	if c.IsOrgAdmin() {
		return true
	}
	if c.usrOrg != nil && c.usrOrg.PermRw == 1 {
		return true
	}
	return false
}
func (c *OrgPerm) CanDownload() bool {
	if c.IsOrgAdmin() {
		return true
	}
	if c.usrOrg != nil && c.usrOrg.PermDown == 1 {
		return true
	}
	return false
}
func (c *OrgPerm) CanExec() bool {
	if c.IsOrgAdmin() {
		return true
	}
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

type UserPipeOrgPerm struct {
	OrgId     string `xorm:"org_id"`
	OrgName   string `xorm:"org_name"`
	OrgUid    string `xorm:"org_uid"`
	OrgPublic int    `xorm:"org_public"`
	OpPublic  int    `xorm:"op_public"`
	CurUid    string `xorm:"cur_uid"`
	PermAdm   int    `xorm:"perm_adm"`
	PermRw    int    `xorm:"perm_rw"`
	PermExec  int    `xorm:"perm_exec"`
}
type PipePerm struct {
	lgusr *model.TUser
	pipe  *model.TPipeline
	perms []*UserPipeOrgPerm
}

func NewPipePerm(lgusr *model.TUser, pipeId string) *PipePerm {
	c := &PipePerm{lgusr: lgusr}
	pipe := &model.TPipeline{}
	ok := false
	if pipeId != "" {
		ok, _ = comm.Db.Where("id=?", pipeId).Get(pipe)
	}
	if ok {
		c.pipe = pipe
		if comm.IsMySQL && lgusr != nil {
			ses := comm.Db.SQL(`
select org.id as org_id,org.name as org_name,org.uid as org_uid,org.public as org_public,op.public as op_public,
uo.uid as cur_uid,uo.perm_adm,uo.perm_rw,uo.perm_exec,uo.perm_down
from t_org org
JOIN t_org_pipe op ON op.pipe_id=? and org.id=op.org_id
LEFT JOIN t_user_org uo ON uo.uid=? and org.id=uo.org_id
where org.deleted!=1 or org.public=1
			`, pipe.Id, lgusr.Id)
			ses.Find(&c.perms)
		}
	}
	return c
}
func (c *PipePerm) IsAdmin() bool {
	if c.lgusr != nil && IsAdmin(c.lgusr) {
		return true
	}
	return false
}
func (c *PipePerm) IsPipeOwner() bool {
	if c.pipe != nil && c.lgusr != nil && c.pipe.Uid == c.lgusr.Id {
		return true
	}
	return false
}
func (c *PipePerm) CanRead() bool {
	if c.IsAdmin() || c.IsPipeOwner() {
		return true
	}
	for _, v := range c.perms {
		if c.lgusr != nil && v.OrgUid == c.lgusr.Id {
			return true
		}
		if v.OrgPublic == 1 {
			return true
		}
		if v.CurUid != "" {
			return true
		}
	}
	return false
}
func (c *PipePerm) CanWrite() bool {
	if c.IsAdmin() || c.IsPipeOwner() {
		return true
	}
	for _, v := range c.perms {
		if c.lgusr != nil && v.OrgUid == c.lgusr.Id {
			return true
		}
		if v.CurUid != "" && (v.PermAdm == 1 || v.PermRw == 1) {
			return true
		}
	}
	return false
}
func (c *PipePerm) CanExec() bool {
	if c.IsAdmin() || c.IsPipeOwner() {
		return true
	}
	for _, v := range c.perms {
		if c.lgusr != nil && v.OrgUid == c.lgusr.Id {
			return true
		}
		if v.CurUid != "" && (v.PermAdm == 1 || v.PermExec == 1) {
			return true
		}
	}
	return false
}

//LgUser maybe null
func (c *PipePerm) LgUser() *model.TUser {
	return c.lgusr
}

//Pipeline maybe null
func (c *PipePerm) Pipeline() *model.TPipeline {
	return c.pipe
}
