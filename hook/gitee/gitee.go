package gitee

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gokins-main/gokins/hook"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func Parse(req *http.Request, secret string) (hook.WebHook, error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("WebhookService Parse err:%+v", err)
			logrus.Warnf("%s", string(debug.Stack()))
		}
	}()
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}
	var wb hook.WebHook
	switch req.Header.Get(hook.GITEE_EVENT) {
	case hook.GITEE_EVENT_PUSH:
		wb, err = parsePushHook(data)
	case hook.GITEE_EVENT_NOTE:
		wb, err = parseCommentHook(data)
	case hook.GITEE_EVENT_PR:
		wb, err = parsePullRequestHook(data)
	// case "pull_request_review_comment":
	// case "issues":
	// case "issue_comment":
	default:
		return nil, errors.New(fmt.Sprintf("hook含有未知的header:%v", req.Header.Get(hook.GITEE_EVENT)))
	}
	if err != nil {
		return nil, err
	}
	sig := req.Header.Get("X-Gitee-Token")
	if secret != sig {
		return wb, errors.New("密钥不正确")
	}
	return wb, nil
}

func parseCommentHook(data []byte) (*hook.PullRequestCommentHook, error) {
	gp := new(giteeCommentHook)
	err := json.Unmarshal(data, gp)
	if err != nil {
		return nil, err
	}
	return convertCommentHook(gp), nil
}

func parsePullRequestHook(data []byte) (*hook.PullRequestHook, error) {
	gp := new(giteePRHook)
	err := json.Unmarshal(data, gp)
	if err != nil {
		return nil, err
	}
	if gp.Action != "" {
		if gp.Action != hook.ActionOpen && gp.Action != hook.ActionUpdate {
			return nil, fmt.Errorf("action is %v", gp.Action)
		}
	}
	return convertPullRequestHook(gp), nil
}

func parsePushHook(data []byte) (*hook.PushHook, error) {
	gp := new(giteePushHook)
	err := json.Unmarshal(data, gp)
	if err != nil {
		return nil, err
	}
	return convertPushHook(gp), nil
}
func convertPullRequestHook(gp *giteePRHook) *hook.PullRequestHook {
	return &hook.PullRequestHook{
		Action: gp.Action,
		Repo: hook.Repository{
			Ref:         gp.PullRequest.Head.Ref,
			Sha:         gp.PullRequest.Head.Sha,
			CloneURL:    gp.SourceRepo.Repository.CloneUrl,
			CreatedAt:   gp.SourceRepo.Repository.CreatedAt,
			Branch:      gp.PullRequest.Head.Ref,
			Description: gp.SourceRepo.Repository.Description,
			FullName:    gp.SourceRepo.Repository.FullName,
			GitHttpURL:  gp.SourceRepo.Repository.GitHttpUrl,
			GitShhURL:   gp.SourceRepo.Repository.GitSshUrl,
			GitSvnURL:   gp.SourceRepo.Repository.GitSvnUrl,
			GitURL:      gp.SourceRepo.Repository.GitUrl,
			HtmlURL:     gp.SourceRepo.Repository.HtmlUrl,
			SshURL:      gp.SourceRepo.Repository.SshUrl,
			SvnURL:      gp.SourceRepo.Repository.SvnUrl,
			Name:        gp.SourceRepo.Repository.Name,
			Private:     gp.SourceRepo.Repository.Private,
			URL:         gp.SourceRepo.Repository.Url,
			Owner:       gp.SourceRepo.Repository.Owner.UserName,
			RepoType:    "gitee",
			RepoOpenid:  strconv.Itoa(gp.Repository.Id),
		},
		TargetRepo: hook.Repository{
			Ref:         gp.PullRequest.Base.Ref,
			Sha:         gp.PullRequest.Base.Sha,
			CloneURL:    gp.TargetRepo.Repository.CloneUrl,
			CreatedAt:   gp.TargetRepo.Repository.CreatedAt,
			Branch:      gp.PullRequest.Base.Ref,
			Description: gp.TargetRepo.Repository.Description,
			FullName:    gp.TargetRepo.Repository.FullName,
			GitHttpURL:  gp.TargetRepo.Repository.GitHttpUrl,
			GitShhURL:   gp.TargetRepo.Repository.GitSshUrl,
			GitSvnURL:   gp.TargetRepo.Repository.GitSvnUrl,
			GitURL:      gp.TargetRepo.Repository.GitUrl,
			HtmlURL:     gp.TargetRepo.Repository.HtmlUrl,
			SshURL:      gp.TargetRepo.Repository.SshUrl,
			SvnURL:      gp.TargetRepo.Repository.SvnUrl,
			Name:        gp.TargetRepo.Repository.Name,
			Private:     gp.TargetRepo.Repository.Private,
			URL:         gp.TargetRepo.Repository.Url,
			Owner:       gp.TargetRepo.Repository.Owner.UserName,
			RepoType:    "gitee",
			RepoOpenid:  "",
		},
		PullRequest: hook.PullRequest{
			Number: gp.PullRequest.Number,
			Body:   gp.PullRequest.Body,
			Title:  gp.PullRequest.Title,
			Base: hook.Reference{
				Name: gp.PullRequest.Base.Ref,
				Path: gp.PullRequest.Base.Repo.Path,
				Sha:  gp.PullRequest.Base.Sha,
			},
			Head: hook.Reference{
				Name: gp.PullRequest.Head.Ref,
				Path: gp.PullRequest.Head.Repo.Path,
				Sha:  gp.PullRequest.Head.Sha,
			},
			Author: hook.User{
				UserName: gp.PullRequest.User.UserName,
			},
			Created: time.Time{},
			Updated: time.Time{},
		},
		Sender: hook.User{},
	}
}
func convertPushHook(gp *giteePushHook) *hook.PushHook {
	branch := ""
	if gp.Ref != "" {
		if len(strings.Split(gp.Ref, "/")) > 2 {
			branch = strings.Split(gp.Ref, "/")[2]
		}
	}
	return &hook.PushHook{
		Ref: gp.Ref,
		Repo: hook.Repository{
			Ref:         gp.Ref,
			Sha:         gp.After,
			CloneURL:    gp.Repository.CloneUrl,
			CreatedAt:   gp.Repository.CreatedAt,
			Branch:      branch,
			Description: gp.Repository.Description,
			FullName:    gp.Repository.FullName,
			GitHttpURL:  gp.Repository.GitHttpUrl,
			GitShhURL:   gp.Repository.GitSshUrl,
			GitSvnURL:   gp.Repository.GitSvnUrl,
			GitURL:      gp.Repository.GitUrl,
			HtmlURL:     gp.Repository.HtmlUrl,
			SshURL:      gp.Repository.SshUrl,
			SvnURL:      gp.Repository.SvnUrl,
			Name:        gp.Repository.Name,
			Private:     gp.Repository.Private,
			URL:         gp.Repository.Url,
			Owner:       gp.Repository.Owner.Username,
			RepoType:    "gitee",
			RepoOpenid:  strconv.FormatInt(gp.Repository.Id, 10),
		},
		Before: gp.Before,
		After:  gp.After,
		Commit: hook.Commit{
			Message: gp.HeadCommit.Message,
			Link:    gp.HeadCommit.Url,
		},
		Sender: hook.User{
			UserName: gp.User.UserName,
		},
	}
}
func convertCommentHook(gp *giteeCommentHook) *hook.PullRequestCommentHook {
	return &hook.PullRequestCommentHook{
		Action: gp.Action,
		Repo: hook.Repository{
			Ref:         gp.PullRequest.Head.Ref,
			Sha:         gp.PullRequest.Head.Sha,
			CloneURL:    gp.PullRequest.Head.Repo.CloneUrl,
			CreatedAt:   gp.PullRequest.Head.Repo.CreatedAt,
			Branch:      gp.PullRequest.Head.Ref,
			Description: gp.PullRequest.Head.Repo.Description,
			FullName:    gp.PullRequest.Head.Repo.FullName,
			GitHttpURL:  gp.PullRequest.Head.Repo.GitHttpUrl,
			GitShhURL:   gp.PullRequest.Head.Repo.GitSshUrl,
			GitSvnURL:   gp.PullRequest.Head.Repo.GitSvnUrl,
			GitURL:      gp.PullRequest.Head.Repo.GitUrl,
			HtmlURL:     gp.PullRequest.Head.Repo.HtmlUrl,
			SshURL:      gp.PullRequest.Head.Repo.SshUrl,
			SvnURL:      gp.PullRequest.Head.Repo.SvnUrl,
			Name:        gp.PullRequest.Head.Repo.Name,
			Private:     gp.PullRequest.Head.Repo.Private,
			URL:         gp.PullRequest.Head.Repo.Url,
			Owner:       gp.PullRequest.Head.Repo.Owner.UserName,
			RepoType:    "gitee",
			RepoOpenid:  strconv.Itoa(gp.Repository.Id),
		},
		TargetRepo: hook.Repository{
			Ref:         gp.PullRequest.Base.Ref,
			Sha:         gp.PullRequest.Base.Sha,
			CloneURL:    gp.PullRequest.Base.Repo.CloneUrl,
			CreatedAt:   gp.PullRequest.Base.Repo.CreatedAt,
			Branch:      gp.PullRequest.Base.Ref,
			Description: gp.PullRequest.Base.Repo.Description,
			FullName:    gp.PullRequest.Base.Repo.FullName,
			GitHttpURL:  gp.PullRequest.Base.Repo.GitHttpUrl,
			GitShhURL:   gp.PullRequest.Base.Repo.GitSshUrl,
			GitSvnURL:   gp.PullRequest.Base.Repo.GitSvnUrl,
			GitURL:      gp.PullRequest.Base.Repo.GitUrl,
			HtmlURL:     gp.PullRequest.Base.Repo.HtmlUrl,
			SshURL:      gp.PullRequest.Base.Repo.SshUrl,
			SvnURL:      gp.PullRequest.Base.Repo.SvnUrl,
			Name:        gp.PullRequest.Base.Repo.Name,
			Private:     gp.PullRequest.Base.Repo.Private,
			URL:         gp.PullRequest.Base.Repo.Url,
			Owner:       gp.PullRequest.Base.Repo.Owner.UserName,
			RepoType:    "gitee",
			RepoOpenid:  "",
		},
		PullRequest: hook.PullRequest{
			Number: gp.PullRequest.Number,
			Body:   gp.PullRequest.Body,
			Title:  gp.PullRequest.Title,
			Base: hook.Reference{
				Name: gp.PullRequest.Base.Ref,
				Path: gp.PullRequest.Base.Repo.Path,
				Sha:  gp.PullRequest.Base.Sha,
			},
			Head: hook.Reference{
				Name: gp.PullRequest.Head.Ref,
				Path: gp.PullRequest.Head.Repo.Path,
				Sha:  gp.PullRequest.Head.Sha,
			},
			Author: hook.User{
				UserName: gp.PullRequest.User.UserName,
			},
			Created: time.Time{},
			Updated: time.Time{},
		},
		Comment: hook.Comment{
			Body: gp.Note,
			Author: hook.User{
				UserName: gp.Author.UserName,
			},
		},
		Sender: hook.User{
			UserName: gp.PullRequest.User.UserName,
		},
	}
}

type giteePushHook struct {
	Ref                string `json:"ref"`
	Before             string `json:"before"`
	After              string `json:"after"`
	TotalCommitsCount  int    `json:"total_commits_count"`
	CommitsMoreThanTen bool   `json:"commits_more_than_ten"`
	Created            bool   `json:"created"`
	Deleted            bool   `json:"deleted"`
	Compare            string `json:"compare"`
	Commits            []struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		ParentIds []string  `json:"parent_ids"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Time     time.Time `json:"time"`
			Id       int       `json:"id"`
			Name     string    `json:"name"`
			Email    string    `json:"email"`
			Username string    `json:"username"`
			UserName string    `json:"user_name"`
			Url      string    `json:"url"`
		} `json:"author"`
		Committer struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"committer"`
		Distinct bool          `json:"distinct"`
		Added    []interface{} `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []string      `json:"modified"`
	} `json:"commits"`
	HeadCommit struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		ParentIds []string  `json:"parent_ids"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Time     time.Time `json:"time"`
			Id       int       `json:"id"`
			Name     string    `json:"name"`
			Email    string    `json:"email"`
			Username string    `json:"username"`
			UserName string    `json:"user_name"`
			Url      string    `json:"url"`
		} `json:"author"`
		Committer struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"committer"`
		Distinct bool          `json:"distinct"`
		Added    []interface{} `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []string      `json:"modified"`
	} `json:"head_commit"`
	Repository struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		FullName string `json:"full_name"`
		Owner    struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"owner"`
		Private           bool        `json:"private"`
		HtmlUrl           string      `json:"html_url"`
		Url               string      `json:"url"`
		Description       string      `json:"description"`
		Fork              bool        `json:"fork"`
		CreatedAt         time.Time   `json:"created_at"`
		UpdatedAt         time.Time   `json:"updated_at"`
		PushedAt          time.Time   `json:"pushed_at"`
		GitUrl            string      `json:"git_url"`
		SshUrl            string      `json:"ssh_url"`
		CloneUrl          string      `json:"clone_url"`
		SvnUrl            string      `json:"svn_url"`
		GitHttpUrl        string      `json:"git_http_url"`
		GitSshUrl         string      `json:"git_ssh_url"`
		GitSvnUrl         string      `json:"git_svn_url"`
		Homepage          interface{} `json:"homepage"`
		StargazersCount   int         `json:"stargazers_count"`
		WatchersCount     int         `json:"watchers_count"`
		ForksCount        int         `json:"forks_count"`
		Language          string      `json:"language"`
		HasIssues         bool        `json:"has_issues"`
		HasWiki           bool        `json:"has_wiki"`
		HasPages          bool        `json:"has_pages"`
		License           string      `json:"license"`
		OpenIssuesCount   int         `json:"open_issues_count"`
		DefaultBranch     string      `json:"default_branch"`
		Namespace         string      `json:"namespace"`
		NameWithNamespace string      `json:"name_with_namespace"`
		PathWithNamespace string      `json:"path_with_namespace"`
	} `json:"repository"`
	Project struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		FullName string `json:"full_name"`
		Owner    struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"owner"`
		Private           bool        `json:"private"`
		HtmlUrl           string      `json:"html_url"`
		Url               string      `json:"url"`
		Description       string      `json:"description"`
		Fork              bool        `json:"fork"`
		CreatedAt         time.Time   `json:"created_at"`
		UpdatedAt         time.Time   `json:"updated_at"`
		PushedAt          time.Time   `json:"pushed_at"`
		GitUrl            string      `json:"git_url"`
		SshUrl            string      `json:"ssh_url"`
		CloneUrl          string      `json:"clone_url"`
		SvnUrl            string      `json:"svn_url"`
		GitHttpUrl        string      `json:"git_http_url"`
		GitSshUrl         string      `json:"git_ssh_url"`
		GitSvnUrl         string      `json:"git_svn_url"`
		Homepage          interface{} `json:"homepage"`
		StargazersCount   int         `json:"stargazers_count"`
		WatchersCount     int         `json:"watchers_count"`
		ForksCount        int         `json:"forks_count"`
		Language          string      `json:"language"`
		HasIssues         bool        `json:"has_issues"`
		HasWiki           bool        `json:"has_wiki"`
		HasPages          bool        `json:"has_pages"`
		License           string      `json:"license"`
		OpenIssuesCount   int         `json:"open_issues_count"`
		DefaultBranch     string      `json:"default_branch"`
		Namespace         string      `json:"namespace"`
		NameWithNamespace string      `json:"name_with_namespace"`
		PathWithNamespace string      `json:"path_with_namespace"`
	} `json:"project"`
	UserId   int    `json:"user_id"`
	UserName string `json:"user_name"`
	User     struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		UserName string `json:"user_name"`
		Url      string `json:"url"`
	} `json:"user"`
	Pusher struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		UserName string `json:"user_name"`
		Url      string `json:"url"`
	} `json:"pusher"`
	Sender struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		UserName  string `json:"user_name"`
		Url       string `json:"url"`
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
		Type      string `json:"type"`
		SiteAdmin bool   `json:"site_admin"`
	} `json:"sender"`
	Enterprise interface{} `json:"enterprise"`
	HookName   string      `json:"hook_name"`
	HookId     int         `json:"hook_id"`
	HookUrl    string      `json:"hook_url"`
	Password   string      `json:"password"`
	Timestamp  string      `json:"timestamp"`
	Sign       string      `json:"sign"`
}

type giteePRHook struct {
	Action      string `json:"action"`
	ActionDesc  string `json:"action_desc"`
	PullRequest struct {
		Id                 int           `json:"id"`
		Number             int64         `json:"number"`
		State              string        `json:"state"`
		HtmlUrl            string        `json:"html_url"`
		DiffUrl            string        `json:"diff_url"`
		PatchUrl           string        `json:"patch_url"`
		Title              string        `json:"title"`
		Body               string        `json:"body"`
		Labels             []interface{} `json:"labels"`
		Languages          []interface{} `json:"languages"`
		CreatedAt          time.Time     `json:"created_at"`
		UpdatedAt          time.Time     `json:"updated_at"`
		ClosedAt           interface{}   `json:"closed_at"`
		MergedAt           interface{}   `json:"merged_at"`
		MergeCommitSha     string        `json:"merge_commit_sha"`
		MergeReferenceName string        `json:"merge_reference_name"`
		User               struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"user"`
		Assignee  interface{} `json:"assignee"`
		Assignees []struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"assignees"`
		Tester  interface{} `json:"tester"`
		Testers []struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"testers"`
		NeedTest   bool        `json:"need_test"`
		NeedReview bool        `json:"need_review"`
		Milestone  interface{} `json:"milestone"`
		Head       struct {
			Label string `json:"label"`
			Ref   string `json:"ref"`
			Sha   string `json:"sha"`
			User  struct {
				Id        int    `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Username  string `json:"username"`
				UserName  string `json:"user_name"`
				Url       string `json:"url"`
				Login     string `json:"login"`
				AvatarUrl string `json:"avatar_url"`
				HtmlUrl   string `json:"html_url"`
				Type      string `json:"type"`
				SiteAdmin bool   `json:"site_admin"`
			} `json:"user"`
			Repo struct {
				Id       int    `json:"id"`
				Name     string `json:"name"`
				Path     string `json:"path"`
				FullName string `json:"full_name"`
				Owner    struct {
					Id        int    `json:"id"`
					Name      string `json:"name"`
					Email     string `json:"email"`
					Username  string `json:"username"`
					UserName  string `json:"user_name"`
					Url       string `json:"url"`
					Login     string `json:"login"`
					AvatarUrl string `json:"avatar_url"`
					HtmlUrl   string `json:"html_url"`
					Type      string `json:"type"`
					SiteAdmin bool   `json:"site_admin"`
				} `json:"owner"`
				Private           bool        `json:"private"`
				HtmlUrl           string      `json:"html_url"`
				Url               string      `json:"url"`
				Description       string      `json:"description"`
				Fork              bool        `json:"fork"`
				CreatedAt         time.Time   `json:"created_at"`
				UpdatedAt         time.Time   `json:"updated_at"`
				PushedAt          time.Time   `json:"pushed_at"`
				GitUrl            string      `json:"git_url"`
				SshUrl            string      `json:"ssh_url"`
				CloneUrl          string      `json:"clone_url"`
				SvnUrl            string      `json:"svn_url"`
				GitHttpUrl        string      `json:"git_http_url"`
				GitSshUrl         string      `json:"git_ssh_url"`
				GitSvnUrl         string      `json:"git_svn_url"`
				Homepage          interface{} `json:"homepage"`
				StargazersCount   int         `json:"stargazers_count"`
				WatchersCount     int         `json:"watchers_count"`
				ForksCount        int         `json:"forks_count"`
				Language          string      `json:"language"`
				HasIssues         bool        `json:"has_issues"`
				HasWiki           bool        `json:"has_wiki"`
				HasPages          bool        `json:"has_pages"`
				License           string      `json:"license"`
				OpenIssuesCount   int         `json:"open_issues_count"`
				DefaultBranch     string      `json:"default_branch"`
				Namespace         string      `json:"namespace"`
				NameWithNamespace string      `json:"name_with_namespace"`
				PathWithNamespace string      `json:"path_with_namespace"`
			} `json:"repo"`
		} `json:"head"`
		Base struct {
			Label string `json:"label"`
			Ref   string `json:"ref"`
			Sha   string `json:"sha"`
			User  struct {
				Id        int    `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Username  string `json:"username"`
				UserName  string `json:"user_name"`
				Url       string `json:"url"`
				Login     string `json:"login"`
				AvatarUrl string `json:"avatar_url"`
				HtmlUrl   string `json:"html_url"`
				Type      string `json:"type"`
				SiteAdmin bool   `json:"site_admin"`
			} `json:"user"`
			Repo struct {
				Id       int    `json:"id"`
				Name     string `json:"name"`
				Path     string `json:"path"`
				FullName string `json:"full_name"`
				Owner    struct {
					Id        int    `json:"id"`
					Name      string `json:"name"`
					Email     string `json:"email"`
					Username  string `json:"username"`
					UserName  string `json:"user_name"`
					Url       string `json:"url"`
					Login     string `json:"login"`
					AvatarUrl string `json:"avatar_url"`
					HtmlUrl   string `json:"html_url"`
					Type      string `json:"type"`
					SiteAdmin bool   `json:"site_admin"`
				} `json:"owner"`
				Private           bool        `json:"private"`
				HtmlUrl           string      `json:"html_url"`
				Url               string      `json:"url"`
				Description       string      `json:"description"`
				Fork              bool        `json:"fork"`
				CreatedAt         time.Time   `json:"created_at"`
				UpdatedAt         time.Time   `json:"updated_at"`
				PushedAt          time.Time   `json:"pushed_at"`
				GitUrl            string      `json:"git_url"`
				SshUrl            string      `json:"ssh_url"`
				CloneUrl          string      `json:"clone_url"`
				SvnUrl            string      `json:"svn_url"`
				GitHttpUrl        string      `json:"git_http_url"`
				GitSshUrl         string      `json:"git_ssh_url"`
				GitSvnUrl         string      `json:"git_svn_url"`
				Homepage          interface{} `json:"homepage"`
				StargazersCount   int         `json:"stargazers_count"`
				WatchersCount     int         `json:"watchers_count"`
				ForksCount        int         `json:"forks_count"`
				Language          string      `json:"language"`
				HasIssues         bool        `json:"has_issues"`
				HasWiki           bool        `json:"has_wiki"`
				HasPages          bool        `json:"has_pages"`
				License           string      `json:"license"`
				OpenIssuesCount   int         `json:"open_issues_count"`
				DefaultBranch     string      `json:"default_branch"`
				Namespace         string      `json:"namespace"`
				NameWithNamespace string      `json:"name_with_namespace"`
				PathWithNamespace string      `json:"path_with_namespace"`
			} `json:"repo"`
		} `json:"base"`
		Merged      bool   `json:"merged"`
		Mergeable   bool   `json:"mergeable"`
		MergeStatus string `json:"merge_status"`
		UpdatedBy   struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"updated_by"`
		Comments     int `json:"comments"`
		Commits      int `json:"commits"`
		Additions    int `json:"additions"`
		Deletions    int `json:"deletions"`
		ChangedFiles int `json:"changed_files"`
	} `json:"pull_request"`
	Number         int           `json:"number"`
	Iid            int           `json:"iid"`
	Title          string        `json:"title"`
	Body           string        `json:"body"`
	Languages      []interface{} `json:"languages"`
	State          string        `json:"state"`
	MergeStatus    string        `json:"merge_status"`
	MergeCommitSha string        `json:"merge_commit_sha"`
	Url            string        `json:"url"`
	SourceBranch   string        `json:"source_branch"`
	SourceRepo     struct {
		Project struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			Path     string `json:"path"`
			FullName string `json:"full_name"`
			Owner    struct {
				Id        int    `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Username  string `json:"username"`
				UserName  string `json:"user_name"`
				Url       string `json:"url"`
				Login     string `json:"login"`
				AvatarUrl string `json:"avatar_url"`
				HtmlUrl   string `json:"html_url"`
				Type      string `json:"type"`
				SiteAdmin bool   `json:"site_admin"`
			} `json:"owner"`
			Private           bool        `json:"private"`
			HtmlUrl           string      `json:"html_url"`
			Url               string      `json:"url"`
			Description       string      `json:"description"`
			Fork              bool        `json:"fork"`
			CreatedAt         time.Time   `json:"created_at"`
			UpdatedAt         time.Time   `json:"updated_at"`
			PushedAt          time.Time   `json:"pushed_at"`
			GitUrl            string      `json:"git_url"`
			SshUrl            string      `json:"ssh_url"`
			CloneUrl          string      `json:"clone_url"`
			SvnUrl            string      `json:"svn_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitSvnUrl         string      `json:"git_svn_url"`
			Homepage          interface{} `json:"homepage"`
			StargazersCount   int         `json:"stargazers_count"`
			WatchersCount     int         `json:"watchers_count"`
			ForksCount        int         `json:"forks_count"`
			Language          string      `json:"language"`
			HasIssues         bool        `json:"has_issues"`
			HasWiki           bool        `json:"has_wiki"`
			HasPages          bool        `json:"has_pages"`
			License           string      `json:"license"`
			OpenIssuesCount   int         `json:"open_issues_count"`
			DefaultBranch     string      `json:"default_branch"`
			Namespace         string      `json:"namespace"`
			NameWithNamespace string      `json:"name_with_namespace"`
			PathWithNamespace string      `json:"path_with_namespace"`
		} `json:"project"`
		Repository struct {
			Id       int64  `json:"id"`
			Name     string `json:"name"`
			Path     string `json:"path"`
			FullName string `json:"full_name"`
			Owner    struct {
				Id        int    `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Username  string `json:"username"`
				UserName  string `json:"user_name"`
				Url       string `json:"url"`
				Login     string `json:"login"`
				AvatarUrl string `json:"avatar_url"`
				HtmlUrl   string `json:"html_url"`
				Type      string `json:"type"`
				SiteAdmin bool   `json:"site_admin"`
			} `json:"owner"`
			Private           bool        `json:"private"`
			HtmlUrl           string      `json:"html_url"`
			Url               string      `json:"url"`
			Description       string      `json:"description"`
			Fork              bool        `json:"fork"`
			CreatedAt         time.Time   `json:"created_at"`
			UpdatedAt         time.Time   `json:"updated_at"`
			PushedAt          time.Time   `json:"pushed_at"`
			GitUrl            string      `json:"git_url"`
			SshUrl            string      `json:"ssh_url"`
			CloneUrl          string      `json:"clone_url"`
			SvnUrl            string      `json:"svn_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitSvnUrl         string      `json:"git_svn_url"`
			Homepage          interface{} `json:"homepage"`
			StargazersCount   int         `json:"stargazers_count"`
			WatchersCount     int         `json:"watchers_count"`
			ForksCount        int         `json:"forks_count"`
			Language          string      `json:"language"`
			HasIssues         bool        `json:"has_issues"`
			HasWiki           bool        `json:"has_wiki"`
			HasPages          bool        `json:"has_pages"`
			License           string      `json:"license"`
			OpenIssuesCount   int         `json:"open_issues_count"`
			DefaultBranch     string      `json:"default_branch"`
			Namespace         string      `json:"namespace"`
			NameWithNamespace string      `json:"name_with_namespace"`
			PathWithNamespace string      `json:"path_with_namespace"`
		} `json:"repository"`
	} `json:"source_repo"`
	TargetBranch string `json:"target_branch"`
	TargetRepo   struct {
		Project struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			Path     string `json:"path"`
			FullName string `json:"full_name"`
			Owner    struct {
				Id        int    `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Username  string `json:"username"`
				UserName  string `json:"user_name"`
				Url       string `json:"url"`
				Login     string `json:"login"`
				AvatarUrl string `json:"avatar_url"`
				HtmlUrl   string `json:"html_url"`
				Type      string `json:"type"`
				SiteAdmin bool   `json:"site_admin"`
			} `json:"owner"`
			Private           bool        `json:"private"`
			HtmlUrl           string      `json:"html_url"`
			Url               string      `json:"url"`
			Description       string      `json:"description"`
			Fork              bool        `json:"fork"`
			CreatedAt         time.Time   `json:"created_at"`
			UpdatedAt         time.Time   `json:"updated_at"`
			PushedAt          time.Time   `json:"pushed_at"`
			GitUrl            string      `json:"git_url"`
			SshUrl            string      `json:"ssh_url"`
			CloneUrl          string      `json:"clone_url"`
			SvnUrl            string      `json:"svn_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitSvnUrl         string      `json:"git_svn_url"`
			Homepage          interface{} `json:"homepage"`
			StargazersCount   int         `json:"stargazers_count"`
			WatchersCount     int         `json:"watchers_count"`
			ForksCount        int         `json:"forks_count"`
			Language          string      `json:"language"`
			HasIssues         bool        `json:"has_issues"`
			HasWiki           bool        `json:"has_wiki"`
			HasPages          bool        `json:"has_pages"`
			License           string      `json:"license"`
			OpenIssuesCount   int         `json:"open_issues_count"`
			DefaultBranch     string      `json:"default_branch"`
			Namespace         string      `json:"namespace"`
			NameWithNamespace string      `json:"name_with_namespace"`
			PathWithNamespace string      `json:"path_with_namespace"`
		} `json:"project"`
		Repository struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			Path     string `json:"path"`
			FullName string `json:"full_name"`
			Owner    struct {
				Id        int    `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Username  string `json:"username"`
				UserName  string `json:"user_name"`
				Url       string `json:"url"`
				Login     string `json:"login"`
				AvatarUrl string `json:"avatar_url"`
				HtmlUrl   string `json:"html_url"`
				Type      string `json:"type"`
				SiteAdmin bool   `json:"site_admin"`
			} `json:"owner"`
			Private           bool        `json:"private"`
			HtmlUrl           string      `json:"html_url"`
			Url               string      `json:"url"`
			Description       string      `json:"description"`
			Fork              bool        `json:"fork"`
			CreatedAt         time.Time   `json:"created_at"`
			UpdatedAt         time.Time   `json:"updated_at"`
			PushedAt          time.Time   `json:"pushed_at"`
			GitUrl            string      `json:"git_url"`
			SshUrl            string      `json:"ssh_url"`
			CloneUrl          string      `json:"clone_url"`
			SvnUrl            string      `json:"svn_url"`
			GitHttpUrl        string      `json:"git_http_url"`
			GitSshUrl         string      `json:"git_ssh_url"`
			GitSvnUrl         string      `json:"git_svn_url"`
			Homepage          interface{} `json:"homepage"`
			StargazersCount   int         `json:"stargazers_count"`
			WatchersCount     int         `json:"watchers_count"`
			ForksCount        int         `json:"forks_count"`
			Language          string      `json:"language"`
			HasIssues         bool        `json:"has_issues"`
			HasWiki           bool        `json:"has_wiki"`
			HasPages          bool        `json:"has_pages"`
			License           string      `json:"license"`
			OpenIssuesCount   int         `json:"open_issues_count"`
			DefaultBranch     string      `json:"default_branch"`
			Namespace         string      `json:"namespace"`
			NameWithNamespace string      `json:"name_with_namespace"`
			PathWithNamespace string      `json:"path_with_namespace"`
		} `json:"repository"`
	} `json:"target_repo"`
	Project struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		FullName string `json:"full_name"`
		Owner    struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"owner"`
		Private           bool        `json:"private"`
		HtmlUrl           string      `json:"html_url"`
		Url               string      `json:"url"`
		Description       string      `json:"description"`
		Fork              bool        `json:"fork"`
		CreatedAt         time.Time   `json:"created_at"`
		UpdatedAt         time.Time   `json:"updated_at"`
		PushedAt          time.Time   `json:"pushed_at"`
		GitUrl            string      `json:"git_url"`
		SshUrl            string      `json:"ssh_url"`
		CloneUrl          string      `json:"clone_url"`
		SvnUrl            string      `json:"svn_url"`
		GitHttpUrl        string      `json:"git_http_url"`
		GitSshUrl         string      `json:"git_ssh_url"`
		GitSvnUrl         string      `json:"git_svn_url"`
		Homepage          interface{} `json:"homepage"`
		StargazersCount   int         `json:"stargazers_count"`
		WatchersCount     int         `json:"watchers_count"`
		ForksCount        int         `json:"forks_count"`
		Language          string      `json:"language"`
		HasIssues         bool        `json:"has_issues"`
		HasWiki           bool        `json:"has_wiki"`
		HasPages          bool        `json:"has_pages"`
		License           string      `json:"license"`
		OpenIssuesCount   int         `json:"open_issues_count"`
		DefaultBranch     string      `json:"default_branch"`
		Namespace         string      `json:"namespace"`
		NameWithNamespace string      `json:"name_with_namespace"`
		PathWithNamespace string      `json:"path_with_namespace"`
	} `json:"project"`
	Repository struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		FullName string `json:"full_name"`
		Owner    struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"owner"`
		Private           bool        `json:"private"`
		HtmlUrl           string      `json:"html_url"`
		Url               string      `json:"url"`
		Description       string      `json:"description"`
		Fork              bool        `json:"fork"`
		CreatedAt         time.Time   `json:"created_at"`
		UpdatedAt         time.Time   `json:"updated_at"`
		PushedAt          time.Time   `json:"pushed_at"`
		GitUrl            string      `json:"git_url"`
		SshUrl            string      `json:"ssh_url"`
		CloneUrl          string      `json:"clone_url"`
		SvnUrl            string      `json:"svn_url"`
		GitHttpUrl        string      `json:"git_http_url"`
		GitSshUrl         string      `json:"git_ssh_url"`
		GitSvnUrl         string      `json:"git_svn_url"`
		Homepage          interface{} `json:"homepage"`
		StargazersCount   int         `json:"stargazers_count"`
		WatchersCount     int         `json:"watchers_count"`
		ForksCount        int         `json:"forks_count"`
		Language          string      `json:"language"`
		HasIssues         bool        `json:"has_issues"`
		HasWiki           bool        `json:"has_wiki"`
		HasPages          bool        `json:"has_pages"`
		License           string      `json:"license"`
		OpenIssuesCount   int         `json:"open_issues_count"`
		DefaultBranch     string      `json:"default_branch"`
		Namespace         string      `json:"namespace"`
		NameWithNamespace string      `json:"name_with_namespace"`
		PathWithNamespace string      `json:"path_with_namespace"`
	} `json:"repository"`
	Author struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		UserName  string `json:"user_name"`
		Url       string `json:"url"`
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
		Type      string `json:"type"`
		SiteAdmin bool   `json:"site_admin"`
	} `json:"author"`
	UpdatedBy struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		UserName  string `json:"user_name"`
		Url       string `json:"url"`
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
		Type      string `json:"type"`
		SiteAdmin bool   `json:"site_admin"`
	} `json:"updated_by"`
	Sender struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		UserName  string `json:"user_name"`
		Url       string `json:"url"`
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
		Type      string `json:"type"`
		SiteAdmin bool   `json:"site_admin"`
	} `json:"sender"`
	TargetUser interface{} `json:"target_user"`
	Enterprise interface{} `json:"enterprise"`
	HookName   string      `json:"hook_name"`
	HookId     int         `json:"hook_id"`
	HookUrl    string      `json:"hook_url"`
	Password   string      `json:"password"`
	Timestamp  string      `json:"timestamp"`
	Sign       string      `json:"sign"`
}

type giteeCommentHook struct {
	Action  string `json:"action"`
	Comment struct {
		Id   int    `json:"id"`
		Body string `json:"body"`
		User struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"user"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		HtmlUrl   string    `json:"html_url"`
	} `json:"comment"`
	Repository struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		FullName string `json:"full_name"`
		Owner    struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"owner"`
		Private           bool        `json:"private"`
		HtmlUrl           string      `json:"html_url"`
		Url               string      `json:"url"`
		Description       string      `json:"description"`
		Fork              bool        `json:"fork"`
		CreatedAt         time.Time   `json:"created_at"`
		UpdatedAt         time.Time   `json:"updated_at"`
		PushedAt          time.Time   `json:"pushed_at"`
		GitUrl            string      `json:"git_url"`
		SshUrl            string      `json:"ssh_url"`
		CloneUrl          string      `json:"clone_url"`
		SvnUrl            string      `json:"svn_url"`
		GitHttpUrl        string      `json:"git_http_url"`
		GitSshUrl         string      `json:"git_ssh_url"`
		GitSvnUrl         string      `json:"git_svn_url"`
		Homepage          interface{} `json:"homepage"`
		StargazersCount   int         `json:"stargazers_count"`
		WatchersCount     int         `json:"watchers_count"`
		ForksCount        int         `json:"forks_count"`
		Language          string      `json:"language"`
		HasIssues         bool        `json:"has_issues"`
		HasWiki           bool        `json:"has_wiki"`
		HasPages          bool        `json:"has_pages"`
		License           string      `json:"license"`
		OpenIssuesCount   int         `json:"open_issues_count"`
		DefaultBranch     string      `json:"default_branch"`
		Namespace         string      `json:"namespace"`
		NameWithNamespace string      `json:"name_with_namespace"`
		PathWithNamespace string      `json:"path_with_namespace"`
	} `json:"repository"`
	Project struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		FullName string `json:"full_name"`
		Owner    struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"owner"`
		Private           bool        `json:"private"`
		HtmlUrl           string      `json:"html_url"`
		Url               string      `json:"url"`
		Description       string      `json:"description"`
		Fork              bool        `json:"fork"`
		CreatedAt         time.Time   `json:"created_at"`
		UpdatedAt         time.Time   `json:"updated_at"`
		PushedAt          time.Time   `json:"pushed_at"`
		GitUrl            string      `json:"git_url"`
		SshUrl            string      `json:"ssh_url"`
		CloneUrl          string      `json:"clone_url"`
		SvnUrl            string      `json:"svn_url"`
		GitHttpUrl        string      `json:"git_http_url"`
		GitSshUrl         string      `json:"git_ssh_url"`
		GitSvnUrl         string      `json:"git_svn_url"`
		Homepage          interface{} `json:"homepage"`
		StargazersCount   int         `json:"stargazers_count"`
		WatchersCount     int         `json:"watchers_count"`
		ForksCount        int         `json:"forks_count"`
		Language          string      `json:"language"`
		HasIssues         bool        `json:"has_issues"`
		HasWiki           bool        `json:"has_wiki"`
		HasPages          bool        `json:"has_pages"`
		License           string      `json:"license"`
		OpenIssuesCount   int         `json:"open_issues_count"`
		DefaultBranch     string      `json:"default_branch"`
		Namespace         string      `json:"namespace"`
		NameWithNamespace string      `json:"name_with_namespace"`
		PathWithNamespace string      `json:"path_with_namespace"`
	} `json:"project"`
	Author struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		UserName  string `json:"user_name"`
		Url       string `json:"url"`
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
		Type      string `json:"type"`
		SiteAdmin bool   `json:"site_admin"`
	} `json:"author"`
	Sender struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		UserName  string `json:"user_name"`
		Url       string `json:"url"`
		Login     string `json:"login"`
		AvatarUrl string `json:"avatar_url"`
		HtmlUrl   string `json:"html_url"`
		Type      string `json:"type"`
		SiteAdmin bool   `json:"site_admin"`
	} `json:"sender"`
	Url           string      `json:"url"`
	Note          string      `json:"note"`
	NoteableType  string      `json:"noteable_type"`
	NoteableId    int         `json:"noteable_id"`
	Title         string      `json:"title"`
	PerIid        string      `json:"per_iid"`
	ShortCommitId interface{} `json:"short_commit_id"`
	Enterprise    interface{} `json:"enterprise"`
	PullRequest   struct {
		Id                 int           `json:"id"`
		Number             int64         `json:"number"`
		State              string        `json:"state"`
		HtmlUrl            string        `json:"html_url"`
		DiffUrl            string        `json:"diff_url"`
		PatchUrl           string        `json:"patch_url"`
		Title              string        `json:"title"`
		Body               string        `json:"body"`
		Labels             []interface{} `json:"labels"`
		Languages          []interface{} `json:"languages"`
		CreatedAt          time.Time     `json:"created_at"`
		UpdatedAt          time.Time     `json:"updated_at"`
		ClosedAt           interface{}   `json:"closed_at"`
		MergedAt           interface{}   `json:"merged_at"`
		MergeCommitSha     string        `json:"merge_commit_sha"`
		MergeReferenceName string        `json:"merge_reference_name"`
		User               struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"user"`
		Assignee  interface{} `json:"assignee"`
		Assignees []struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"assignees"`
		Tester  interface{} `json:"tester"`
		Testers []struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"testers"`
		NeedTest   bool        `json:"need_test"`
		NeedReview bool        `json:"need_review"`
		Milestone  interface{} `json:"milestone"`
		Head       struct {
			Label string `json:"label"`
			Ref   string `json:"ref"`
			Sha   string `json:"sha"`
			User  struct {
				Id        int    `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Username  string `json:"username"`
				UserName  string `json:"user_name"`
				Url       string `json:"url"`
				Login     string `json:"login"`
				AvatarUrl string `json:"avatar_url"`
				HtmlUrl   string `json:"html_url"`
				Type      string `json:"type"`
				SiteAdmin bool   `json:"site_admin"`
			} `json:"user"`
			Repo struct {
				Id       int64  `json:"id"`
				Name     string `json:"name"`
				Path     string `json:"path"`
				FullName string `json:"full_name"`
				Owner    struct {
					Id        int    `json:"id"`
					Name      string `json:"name"`
					Email     string `json:"email"`
					Username  string `json:"username"`
					UserName  string `json:"user_name"`
					Url       string `json:"url"`
					Login     string `json:"login"`
					AvatarUrl string `json:"avatar_url"`
					HtmlUrl   string `json:"html_url"`
					Type      string `json:"type"`
					SiteAdmin bool   `json:"site_admin"`
				} `json:"owner"`
				Private           bool        `json:"private"`
				HtmlUrl           string      `json:"html_url"`
				Url               string      `json:"url"`
				Description       string      `json:"description"`
				Fork              bool        `json:"fork"`
				CreatedAt         time.Time   `json:"created_at"`
				UpdatedAt         time.Time   `json:"updated_at"`
				PushedAt          time.Time   `json:"pushed_at"`
				GitUrl            string      `json:"git_url"`
				SshUrl            string      `json:"ssh_url"`
				CloneUrl          string      `json:"clone_url"`
				SvnUrl            string      `json:"svn_url"`
				GitHttpUrl        string      `json:"git_http_url"`
				GitSshUrl         string      `json:"git_ssh_url"`
				GitSvnUrl         string      `json:"git_svn_url"`
				Homepage          interface{} `json:"homepage"`
				StargazersCount   int         `json:"stargazers_count"`
				WatchersCount     int         `json:"watchers_count"`
				ForksCount        int         `json:"forks_count"`
				Language          string      `json:"language"`
				HasIssues         bool        `json:"has_issues"`
				HasWiki           bool        `json:"has_wiki"`
				HasPages          bool        `json:"has_pages"`
				License           string      `json:"license"`
				OpenIssuesCount   int         `json:"open_issues_count"`
				DefaultBranch     string      `json:"default_branch"`
				Namespace         string      `json:"namespace"`
				NameWithNamespace string      `json:"name_with_namespace"`
				PathWithNamespace string      `json:"path_with_namespace"`
			} `json:"repo"`
		} `json:"head"`
		Base struct {
			Label string `json:"label"`
			Ref   string `json:"ref"`
			Sha   string `json:"sha"`
			User  struct {
				Id        int    `json:"id"`
				Name      string `json:"name"`
				Email     string `json:"email"`
				Username  string `json:"username"`
				UserName  string `json:"user_name"`
				Url       string `json:"url"`
				Login     string `json:"login"`
				AvatarUrl string `json:"avatar_url"`
				HtmlUrl   string `json:"html_url"`
				Type      string `json:"type"`
				SiteAdmin bool   `json:"site_admin"`
			} `json:"user"`
			Repo struct {
				Id       int    `json:"id"`
				Name     string `json:"name"`
				Path     string `json:"path"`
				FullName string `json:"full_name"`
				Owner    struct {
					Id        int    `json:"id"`
					Name      string `json:"name"`
					Email     string `json:"email"`
					Username  string `json:"username"`
					UserName  string `json:"user_name"`
					Url       string `json:"url"`
					Login     string `json:"login"`
					AvatarUrl string `json:"avatar_url"`
					HtmlUrl   string `json:"html_url"`
					Type      string `json:"type"`
					SiteAdmin bool   `json:"site_admin"`
				} `json:"owner"`
				Private           bool        `json:"private"`
				HtmlUrl           string      `json:"html_url"`
				Url               string      `json:"url"`
				Description       string      `json:"description"`
				Fork              bool        `json:"fork"`
				CreatedAt         time.Time   `json:"created_at"`
				UpdatedAt         time.Time   `json:"updated_at"`
				PushedAt          time.Time   `json:"pushed_at"`
				GitUrl            string      `json:"git_url"`
				SshUrl            string      `json:"ssh_url"`
				CloneUrl          string      `json:"clone_url"`
				SvnUrl            string      `json:"svn_url"`
				GitHttpUrl        string      `json:"git_http_url"`
				GitSshUrl         string      `json:"git_ssh_url"`
				GitSvnUrl         string      `json:"git_svn_url"`
				Homepage          interface{} `json:"homepage"`
				StargazersCount   int         `json:"stargazers_count"`
				WatchersCount     int         `json:"watchers_count"`
				ForksCount        int         `json:"forks_count"`
				Language          string      `json:"language"`
				HasIssues         bool        `json:"has_issues"`
				HasWiki           bool        `json:"has_wiki"`
				HasPages          bool        `json:"has_pages"`
				License           string      `json:"license"`
				OpenIssuesCount   int         `json:"open_issues_count"`
				DefaultBranch     string      `json:"default_branch"`
				Namespace         string      `json:"namespace"`
				NameWithNamespace string      `json:"name_with_namespace"`
				PathWithNamespace string      `json:"path_with_namespace"`
			} `json:"repo"`
		} `json:"base"`
		Merged      bool   `json:"merged"`
		Mergeable   bool   `json:"mergeable"`
		MergeStatus string `json:"merge_status"`
		UpdatedBy   struct {
			Id        int    `json:"id"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			Username  string `json:"username"`
			UserName  string `json:"user_name"`
			Url       string `json:"url"`
			Login     string `json:"login"`
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			Type      string `json:"type"`
			SiteAdmin bool   `json:"site_admin"`
		} `json:"updated_by"`
		Comments     int `json:"comments"`
		Commits      int `json:"commits"`
		Additions    int `json:"additions"`
		Deletions    int `json:"deletions"`
		ChangedFiles int `json:"changed_files"`
	} `json:"pull_request"`
	HookName  string `json:"hook_name"`
	HookId    int    `json:"hook_id"`
	HookUrl   string `json:"hook_url"`
	Password  string `json:"password"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}
