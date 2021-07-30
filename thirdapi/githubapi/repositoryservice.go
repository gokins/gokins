package githubapi

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
	"strings"
)

type RepositoryService struct {
	client *wrapper
}

// GetRepos
/*
  visibility : 公开(public)、私有(private)或者所有(all)，默认: 所有(all)
  affiliation : owner(授权用户拥有的仓库)、collaborator(授权用户为仓库成员)、organization_member(授权用户为仓库所在组织并有访问仓库权限)、
    enterprise_member(授权用户所在企业并有访问仓库权限)、admin(所有有权限的，包括所管理的组织中所有仓库、所管理的企业的所有仓库)。 可以用逗号分隔符组合。
    如: owner, organization_member 或 owner, collaborator, organization_member
  type : 筛选用户仓库: 其创建(owner)、个人(personal)、其为成员(member)、公开(public)、私有(private)，不能与 visibility 或 affiliation 参数一并使用，否则会报 422 错误
  sort : 排序方式: 创建时间(created)，更新时间(updated)，最后推送时间(pushed)，仓库所属与名称(full_name)。默认: full_name
  direction : 如果sort参数为full_name，用升序(asc)。否则降序(desc)
  q : 搜索关键字
  page : 当前的页码
  per_page : 每页的数量，最大为 100
*/
func (s *RepositoryService) GetRepos(accessToken, username, types, sort, direction string, page, per_page int) (*thirdapi.RepositoryPage, error) {
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGithubGetRepos, types, sort, direction, page, per_page))
	if err != nil {
		logrus.Errorf("Github Api GetRepos Parse err : %v", err)
		return nil, err
	}
	logrus.Debugf("Github Api GetRepos url : %v", parse.String())
	req, err := http.NewRequest("GET", parse.String(), nil)
	if err != nil {
		logrus.Errorf("Github Api GetRepos url :%v Get err : %v", parse, err)
		return nil, err
	}
	req.Header.Set("Authorization", "token "+accessToken)
	resp, err := s.client.HttpClient.Do(req)
	if err != nil {
		logrus.Errorf("Github Api GetRepos url :%v Get err : %v", parse, err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Github Api GetRepos url :%v Resp code : %v", parse, resp.StatusCode)
		return nil, errors.New("Github Api GetRepos failed ")
	}
	repos, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logrus.Errorf("Github Api GetRepos ReadAll err : %v", err)
		return nil, err
	}
	var repoList []*thirdbean.ResultGithubRepo
	err = json.Unmarshal(repos, &repoList)
	if err != nil {
		logrus.Errorf("Github Api GetRepos Unmarshal err : %v", err)
		return nil, err
	}
	lk := resp.Header.Get("Link")
	var totalPages int64 = 1
	if lk != "" {
		splits := strings.Split(lk, ", ")
		for _, v := range splits {
			if strings.Contains(v, `rel="last"`) {
				replace := strings.Replace(strings.Replace(v, "<", "", -1), ">;", "", -1)
				p, errs := url.Parse(replace)
				if errs != nil {
					logrus.Errorf("Github Api GetRepos url Parse err : %v", errs)
				}
				get := p.Query().Get("page")
				totalPages, err = strconv.ParseInt(get, 10, 64)
				if err != nil {
					return nil, err
				}
			}
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
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGithubDeleteHooks, owner, repo, hookId))
	if err != nil {
		logrus.Errorf("Github Api DeleteHooks Parse err : %v", err)
		return err
	}
	request, err := http.NewRequest("DELETE", parse.String(), nil)
	if err != nil {
		logrus.Errorf("Github Api DeleteHooks url :%v Get err : %v", parse, err)
		return err
	}
	request.Header.Set("Authorization", "token "+accessToken)
	logrus.Debugf("Github Api DeleteHooks url : %v", parse)
	resp, err := s.client.HttpClient.Do(request)
	if err != nil {
		logrus.Errorf("Github Api DeleteHooks url :%v Get err : %v", parse, err)
		return err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Github Api DeleteHooks url :%v ReadAll err : %v", parse, err)
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		logrus.Errorf("Github Api DeleteHooks url :%v Resp code : %v", parse, resp.StatusCode)
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
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGithubCreateHooks, owner, repo))
	if err != nil {
		logrus.Errorf("Github Api CreateWebHooks Parse err : %v", err)
		return nil, err
	}
	logrus.Debugf("Github Api CreateWebHooks url : %s", parse.String())
	m := map[string]interface{}{}
	m["url"] = backUrl
	m["content_type"] = "json"
	m["secret"] = password
	m["token"] = accessToken
	obj := map[string]map[string]interface{}{}
	obj["config"] = m
	marshal, err := json.Marshal(obj)
	if err != nil {
		logrus.Errorf("Github Api CreateWebHooks url :%v Get err : %v", parse, err)
		return nil, err
	}
	logrus.Debugf("CreateWebHooks json %s", string(marshal))
	request, err := http.NewRequest("POST", parse.String(), bytes.NewBuffer(marshal))
	if err != nil {
		logrus.Errorf("Github Api CreateWebHooks url :%v Get err : %v", parse, err)
		return nil, err
	}
	request.Header.Set("Authorization", "token "+accessToken)
	request.Header.Set("Content-Type", "application/json")
	resp, err := s.client.HttpClient.Do(request)
	if err != nil {
		logrus.Errorf("Github Api CreateWebHooks url :%v Get err : %v", parse, err)
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Github Api CreateWebHooks url :%v ReadAll err : %v", parse, err)
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		logrus.Errorf("Github Api CreateWebHooks url :%v Resp code : %v", parse, resp.StatusCode)
		return nil, errors.New(string(all))
	}
	k := &thirdbean.ResultGetGithubHook{}
	err = json.Unmarshal(all, k)
	if err != nil {
		logrus.Errorf("Gitee Api CreateWebHooks url :%v ReadAll err : %v", parse, err)
		return nil, err
	}
	return convertHook(k), nil
}

func (s *RepositoryService) GetRepoBranches(accessToken, owner, repo string) ([]*thirdapi.RepositoryBranch, error) {
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGithubGetRepoBranches, owner, repo))
	if err != nil {
		logrus.Errorf("Github Api GetRepoBranches Parse err : %v", err)
		return nil, err
	}
	logrus.Debugf("Github Api GetRepoBranches url : %v", parse)
	req, err := http.NewRequest("GET", parse.String(), nil)
	if err != nil {
		logrus.Errorf("Github Api GetRepoBranches url :%v Get err : %v", parse, err)
		return nil, err
	}
	req.Header.Set("Authorization", "token "+accessToken)
	resp, err := s.client.HttpClient.Do(req)
	if err != nil {
		logrus.Errorf("Github Api GetRepoBranches url :%v Get err : %v", parse, err)
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Github Api GetRepoBranches url :%v ReadAll err : %v", parse, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Github Api GetRepoBranches url :%v Resp code : %v", parse, resp.StatusCode)
		return nil, errors.New(string(all))
	}
	var branchList []*thirdbean.ResultGithubRepoBranch
	err = json.Unmarshal(all, &branchList)
	if err != nil {
		logrus.Errorf("RefreshRepos.GetRepoBranches Unmarshal err : %v", err)
		return nil, err
	}
	return convertBranchList(branchList), err
}

func (s *RepositoryService) GetWebHooks(accessToken, owner, repo string, page, per_page int) ([]*thirdapi.RepositoryHook, error) {
	parse, err := s.client.BaseURL.Parse(s.client.BaseURL.String() + fmt.Sprintf(ApiGithubGetHooks, owner, repo, page, per_page))
	if err != nil {
		logrus.Errorf("Github Api GetWebHooks Parse err : %v", err)
		return nil, err
	}
	logrus.Debugf("Github Api GetWebHooks url : %v", parse)
	req, err := http.NewRequest("GET", parse.String(), nil)
	if err != nil {
		logrus.Errorf("Github Api GetWebHooks url :%v Get err : %v", parse, err)
		return nil, err
	}
	req.Header.Set("Authorization", "token "+accessToken)
	resp, err := s.client.HttpClient.Do(req)
	if err != nil {
		logrus.Errorf("Github Api CreateWebHooks url :%vs Get err : %v", parse, err)
		return nil, err
	}
	defer resp.Body.Close()
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Github Api CreateWebHooks url :%v ReadAll err : %v", parse, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Github Api CreateWebHooks url :%v Resp code : %v", parse, resp.StatusCode)
		return nil, errors.New(string(all))
	}

	hs := make([]*thirdbean.ResultGetGithubHook, 0)
	err = json.Unmarshal(all, &hs)
	if err != nil {
		return nil, err
	}
	return convertHookList(hs), err
}

func convertRepositoryList(ls []*thirdbean.ResultGithubRepo) []*thirdapi.Repository {
	repos := make([]*thirdapi.Repository, 0)
	for _, v := range ls {
		repos = append(repos, convertRepository(v))
	}
	return repos
}

func convertRepository(from *thirdbean.ResultGithubRepo) *thirdapi.Repository {
	return &thirdapi.Repository{
		Id:        strconv.Itoa(from.Id),
		Owner:     from.Owner.Login,
		Name:      from.Name,
		Path:      from.Name,
		Namespace: from.Owner.Login,
		FullName:  from.FullName,
		HtmlURL:   from.HtmlUrl,
		RepoType:  "github",
	}
}
func convertBranchList(ls []*thirdbean.ResultGithubRepoBranch) []*thirdapi.RepositoryBranch {
	repos := make([]*thirdapi.RepositoryBranch, 0)
	for _, v := range ls {
		repos = append(repos, convertBranch(v))
	}
	return repos
}

func convertBranch(from *thirdbean.ResultGithubRepoBranch) *thirdapi.RepositoryBranch {
	return &thirdapi.RepositoryBranch{
		Name: from.Name,
	}
}
func convertHookList(ls []*thirdbean.ResultGetGithubHook) []*thirdapi.RepositoryHook {
	repos := make([]*thirdapi.RepositoryHook, 0)
	for _, v := range ls {
		repos = append(repos, convertHook(v))
	}
	return repos
}

func convertHook(from *thirdbean.ResultGetGithubHook) *thirdapi.RepositoryHook {
	return &thirdapi.RepositoryHook{
		Id:        from.Id,
		Url:       from.Config.Url,
		CreatedAt: from.CreatedAt,
	}
}
