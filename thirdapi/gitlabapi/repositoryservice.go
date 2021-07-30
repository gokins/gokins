package gitlabapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gokins/gokins/bean/thirdbean"
	"github.com/gokins/gokins/thirdapi"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type RepositoryService struct {
	client *wrapper
}

func (s *RepositoryService) GetRepos(accessToken, username, types, sort, direction string, page, per_page int) (*thirdapi.RepositoryPage, error) {
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGitlabGetRepos, username, types, page, per_page))
	if err != nil {
		logrus.Errorf("Gitlab Api GetRepos Parse err : %v", err)
		return nil, err
	}
	logrus.Debugf("Gitlab Api GetRepos url : %v", parse.String())
	req, err := http.NewRequest("GET", parse.String(), nil)
	if err != nil {
		logrus.Errorf("Gitlab Api GetRepos url :%v Get err : %v", parse, err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := s.client.HttpClient.Do(req)
	if err != nil {
		logrus.Errorf("Gitlab Api GetRepos url :%v Get err : %v", parse, err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Gitlab Api GetRepos url :%v Resp code : %v", parse, resp.StatusCode)
		return nil, errors.New("Gitlab Api GetRepos failed ")
	}
	repos, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logrus.Errorf("Gitlab Api GetRepos ReadAll err : %v", err)
		return nil, err
	}
	var repoList []*thirdbean.ResultGitlabRepo
	err = json.Unmarshal(repos, &repoList)
	if err != nil {
		logrus.Errorf("Gitlab Api GetRepos Unmarshal err : %v", err)
		return nil, err
	}
	tp := resp.Header.Get("X-Total-Pages")
	var totalPages int64 = 0
	if tp != "" {
		totalPages, err = strconv.ParseInt(tp, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	list := convertRepositoryList(repoList)
	rp := &thirdapi.RepositoryPage{
		TotalPages: totalPages,
		Ropes:      list,
	}
	return rp, err
}

func (s *RepositoryService) DeleteHooks(accessToken, owner, repo, hookId string) error {
	escape := url.QueryEscape(owner + "/" + repo)
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGitlabDeleteHooks, escape, hookId))
	if err != nil {
		logrus.Errorf("Gitlab Api DeleteHooks Parse err : %v", err)
		return err
	}
	request, err := http.NewRequest("DELETE", parse.String(), nil)
	if err != nil {
		logrus.Errorf("Gitlab Api DeleteHooks url :%v Get err : %v", parse, err)
		return err
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	logrus.Debugf("Gitlab Api DeleteHooks url : %v", parse)
	resp, err := s.client.HttpClient.Do(request)
	if err != nil {
		logrus.Errorf("Gitlab Api DeleteHooks url :%v Get err : %v", parse, err)
		return err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Gitlab Api DeleteHooks url :%v ReadAll err : %v", parse, err)
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		logrus.Errorf("Gitlab Api DeleteHooks url :%v Resp code : %v", parse, resp.StatusCode)
		return errors.New(string(all))
	}
	return nil
}

// CreateWebHooks
/*
  owner : 仓库所属空间地址(企业、组织或个人的地址path)
  repo : 仓库路径(path)
  backUrl : 回调地址
  password : webhook 密钥
*/
func (s *RepositoryService) CreateWebHooks(accessToken, owner, repo, backUrl, password string) (*thirdapi.RepositoryHook, error) {
	escape := url.QueryEscape(owner + "/" + repo)
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGitlabCreateHooks, escape))
	if err != nil {
		logrus.Errorf("Gitlab Api CreateWebHooks Parse err : %v", err)
		return nil, err
	}
	logrus.Debugf("Gitlab Api CreateWebHooks url : %v", parse)
	m := map[string]interface{}{}
	m["url"] = backUrl
	m["token"] = password
	logrus.Infof("gitlab CreateWebHooks backUrl : %s", backUrl)
	marshal, err := json.Marshal(m)
	if err != nil {
		logrus.Errorf("Gitlab Api CreateWebHooks url :%v Get err : %v", parse, err)
		return nil, err
	}
	logrus.Infof("gitlab CreateWebHooks json : %s", string(marshal))
	request, err := http.NewRequest("POST", parse.String(), bytes.NewBuffer(marshal))
	if err != nil {
		logrus.Errorf("Gitlab Api CreateWebHooks url :%v Get err : %v", parse, err)
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
	resp, err := s.client.HttpClient.Do(request)
	if err != nil {
		logrus.Errorf("Gitlab Api CreateWebHooks url :%v Get err : %v", parse, err)
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Gitlab Api CreateWebHooks url :%v ReadAll err : %v", parse, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		logrus.Errorf("Gitlab Api CreateWebHooks url :%v Resp code : %v", parse, resp.StatusCode)
		return nil, errors.New(string(all))
	}
	k := &thirdbean.ResultGetGitlabHook{}
	err = json.Unmarshal(all, k)
	if err != nil {
		logrus.Errorf("Gitee Api CreateWebHooks url :%v ReadAll err : %v", parse, err)
		return nil, err
	}
	return convertHook(k), nil
}

func (s *RepositoryService) GetRepoBranches(accessToken, owner, repo string) ([]*thirdapi.RepositoryBranch, error) {
	escape := url.QueryEscape(owner + "/" + repo)
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGitlabGetRepoBranches, escape))
	if err != nil {
		logrus.Errorf("Gitlab Api GetRepoBranches Parse err : %v", err)
		return nil, err
	}
	logrus.Debugf("Gitlab Api GetRepoBranches url : %v", parse)
	req, err := http.NewRequest("GET", parse.String(), nil)
	if err != nil {
		logrus.Errorf("Gitlab Api GetRepoBranches url :%v Get err : %v", parse, err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := s.client.HttpClient.Do(req)
	if err != nil {
		logrus.Errorf("Gitlab Api GetRepoBranches url :%v Get err : %v", parse, err)
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Gitlab Api GetRepoBranches url :%v ReadAll err : %v", parse, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Gitlab Api GetRepoBranches url :%v Resp code : %v", parse, resp.StatusCode)
		return nil, errors.New(string(all))
	}
	var branchList []*thirdbean.ResultGitlabRepoBranch
	err = json.Unmarshal(all, &branchList)
	if err != nil {
		logrus.Errorf("RefreshRepos.GetRepoBranches Unmarshal err : %v", err)
		return nil, err
	}
	return convertBranchList(branchList), err
}

func (s *RepositoryService) GetWebHooks(accessToken, owner, repo string, page, per_page int) ([]*thirdapi.RepositoryHook, error) {
	escape := url.QueryEscape(owner + "/" + repo)
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGitlabGetHooks, escape, page, per_page))
	if err != nil {
		logrus.Errorf("Gitlab Api GetWebHooks Parse err : %v", err)
		return nil, err
	}
	logrus.Debugf("Gitlab Api GetWebHooks url : %v", parse)
	req, err := http.NewRequest("GET", parse.String(), nil)
	if err != nil {
		logrus.Errorf("Gitlab Api GetWebHooks url :%v Get err : %v", parse, err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := s.client.HttpClient.Do(req)
	if err != nil {
		logrus.Errorf("Gitlab Api CreateWebHooks url :%vs Get err : %v", parse, err)
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Gitlab Api CreateWebHooks url :%v ReadAll err : %v", parse, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Gitlab Api CreateWebHooks url :%v Resp code : %v", parse, resp.StatusCode)
		return nil, errors.New(string(all))
	}

	hs := make([]*thirdbean.ResultGetGitlabHook, 0)
	err = json.Unmarshal(all, &hs)
	if err != nil {
		return nil, err
	}
	return convertHookList(hs), err
}

func convertRepositoryList(ls []*thirdbean.ResultGitlabRepo) []*thirdapi.Repository {
	repos := make([]*thirdapi.Repository, 0)
	for _, v := range ls {
		repos = append(repos, convertRepository(v))
	}
	return repos
}

func convertRepository(from *thirdbean.ResultGitlabRepo) *thirdapi.Repository {
	return &thirdapi.Repository{
		Id:        strconv.Itoa(from.Id),
		Owner:     from.Owner.Username,
		Name:      from.Name,
		Path:      from.Path,
		Namespace: from.Namespace.Path,
		FullName:  from.PathWithNamespace,
		HtmlURL:   from.WebUrl,
		RepoType:  "gitlab",
	}
}
func convertBranchList(ls []*thirdbean.ResultGitlabRepoBranch) []*thirdapi.RepositoryBranch {
	repos := make([]*thirdapi.RepositoryBranch, 0)
	for _, v := range ls {
		repos = append(repos, convertBranch(v))
	}
	return repos
}

func convertBranch(from *thirdbean.ResultGitlabRepoBranch) *thirdapi.RepositoryBranch {
	return &thirdapi.RepositoryBranch{
		Name: from.Name,
	}
}
func convertHookList(ls []*thirdbean.ResultGetGitlabHook) []*thirdapi.RepositoryHook {
	repos := make([]*thirdapi.RepositoryHook, 0)
	for _, v := range ls {
		repos = append(repos, convertHook(v))
	}
	return repos
}

func convertHook(from *thirdbean.ResultGetGitlabHook) *thirdapi.RepositoryHook {
	return &thirdapi.RepositoryHook{
		Id:        from.Id,
		Url:       from.Url,
		CreatedAt: from.CreatedAt,
	}
}
