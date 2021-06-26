package service

import (
	"errors"
	"fmt"
	"github.com/gokins-main/core/runtime"
	"github.com/gokins-main/gokins/comm"
	"github.com/gokins-main/gokins/hook"
	"github.com/gokins-main/gokins/hook/gitea"
	"github.com/gokins-main/gokins/hook/gitee"
	"github.com/gokins-main/gokins/hook/github"
	"github.com/gokins-main/gokins/hook/gitlab"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Parse(req *http.Request, hookType string) (*runtime.Build, error) {
	fn := func(webhook hook.WebHook) (string, error) {
		//TODO get HookSecret
		//repository := webhook.Repository()
		//repo := new(pipeline.TRepo)
		//db := comm.DBMain.GetDB()
		//ok, err := db.
		//	Where("openid=? and deleted != 1 and active != 0",
		//		repository.RepoOpenid).
		//	Get(repo)
		//if err != nil {
		//	return "", err
		//}
		//if !ok {
		//	return "", errors.New("仓库不存在! ")
		//}
		return "repo.HookSecret", nil
	}
	h, err := parseHook(hookType, req, fn)
	if err != nil {
		return nil, err
	}
	var rb *runtime.Build
	switch c := h.(type) {
	case *hook.PullRequestHook:
		fmt.Sprintf("%v", c)
	case *hook.PullRequestCommentHook:
	case *hook.PushHook:
	default:
		return nil, errors.New("Build Parse Failed ")
	}

	//marshal, err := json.Marshal(pb)
	//if err != nil {
	//	return nil, err
	//}
	//TODO insert webhook table
	//repository := h.Repository()
	//repo, b := checkRepo(repository)
	//if !b {
	//	return nil, fmt.Errorf("%v 仓库不存在", repository.Name)
	//}
	//pb.Info.Repository.Id = repo.Id
	//th := &pipeline.THook{
	//	Id:       utils.NewXid(),
	//	HookType: req.Header.Get(common.GITEE_EVENT),
	//	Type:     common.GITEE,
	//	Snapshot: string(marshal),
	//	Status:   common.PIPELINE_VERSION_STATUS_OK,
	//	Msg:      "",
	//}
	//db := comm.DBMain.GetDB()
	//_, err = db.Insert(th)
	//if err != nil {
	//	return nil, err
	//}
	//logrus.Debugf("Parse preBuild %s ", string(marshal))
	return rb, nil
}

func parseHook(hookType string, req *http.Request, fn hook.SecretFunc) (hook.WebHook, error) {
	switch hookType {
	case "gitee", "giteepremium":
		return gitee.Parse(req, fn)
	case "github":
		return github.Parse(req, fn)
	case "gitlab":
		return gitlab.Parse(req, fn)
	case "gitea":
		//TODO 获取SourceHost
		//appinfo, err := GetsParamOAuthKey()
		//if err != nil {
		//	core.Log.Errorf("convertPullRequestURL.GetsParamOAuthKey err : %v", err)
		//	return nil, err
		//}
		//k, info := appinfo.DefAppInfo()
		cl, err := comm.GetThirdApi("k", "info.SourceHost")
		if err != nil {
			logrus.Errorf("convertPullRequestURL.GetThirdApi err : %v", err)
			return nil, err
		}
		return gitea.Parse(req, fn, cl)
	default:
		return nil, fmt.Errorf("未知的webhook类型")
	}
}
