package gitea

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gokins-main/gokins/hook"
	"github.com/sirupsen/logrus"
	"hash"
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
	logrus.Debugf("Gitea Parse ~")
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}
	var wb hook.WebHook
	switch req.Header.Get(hook.GITEA_EVENT) {
	case hook.GITEA_EVENT_PUSH:
		wb, err = parsePushHook(data)
	case hook.GITEA_EVENT_NOTE:
		wb, err = parseCommentHook(data)
	case hook.GITEA_EVENT_PR:
		wb, err = parsePullRequestHook(data)
	default:
		return nil, errors.New(fmt.Sprintf("hook含有未知的header:%v", req.Header.Get(hook.GITEA_EVENT)))
	}
	if err != nil {
		return nil, err
	}
	sig := req.Header.Get("X-Gitea-Signature")
	if !validatePrefix(data, []byte(secret), sig) {
		logrus.Debugf("Gitea validatePrefix failed")
		return wb, errors.New("密钥不正确")
	}
	return wb, nil
}

func Validate(h func() hash.Hash, message, key []byte, signature string) bool {
	decoded, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return validate(h, message, key, decoded)
}

func validatePrefix(message, key []byte, signature string) bool {
	if Validate(sha1.New, message, key, signature) {
		return true
	}
	if Validate(sha256.New, message, key, signature) {
		return true
	}
	return false
}

func validate(h func() hash.Hash, message, key, signature []byte) bool {
	mac := hmac.New(h, key)
	mac.Write(message)
	sum := mac.Sum(nil)
	return hmac.Equal(signature, sum)
}

func parseCommentHook(data []byte) (*hook.PullRequestCommentHook, error) {
	gp := new(giteaCommentHook)
	err := json.Unmarshal(data, gp)
	if err != nil {
		return nil, err
	}
	if !gp.IsPull {
		return nil, errors.New("not pull_request comment")
	}
	return convertCommentHook(gp)
}

func parsePullRequestHook(data []byte) (*hook.PullRequestHook, error) {
	gp := new(giteaPRHook)
	err := json.Unmarshal(data, gp)
	if err != nil {
		return nil, err
	}
	if gp.Action != "" {
		if gp.Action == hook.ActionSynchronize {
			gp.Action = hook.ActionUpdate
		} else if gp.Action == hook.ActionOpened {
			gp.Action = hook.ActionOpen
		} else {
			return nil, fmt.Errorf("action is %v", gp.Action)
		}
	}
	return convertPullRequestHook(gp), nil
}

func parsePushHook(data []byte) (*hook.PushHook, error) {
	gp := new(giteaPushHook)
	err := json.Unmarshal(data, gp)
	if err != nil {
		return nil, err
	}
	return convertPushHook(gp), nil
}
func convertPushHook(gp *giteaPushHook) *hook.PushHook {
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
			GitHttpURL:  gp.Repository.CloneUrl,
			GitShhURL:   gp.Repository.SshUrl,
			GitURL:      gp.Repository.CloneUrl,
			HtmlURL:     gp.Repository.HtmlUrl,
			SshURL:      gp.Repository.SshUrl,
			Name:        gp.Repository.Name,
			Private:     gp.Repository.Private,
			URL:         gp.Repository.HtmlUrl,
			Owner:       gp.Repository.Owner.Login,
			RepoType:    "gitea",
			RepoOpenid:  strconv.Itoa(gp.Repository.Id),
		},
		Before: gp.Before,
		After:  gp.After,
		Commit: hook.Commit{
			Message: gp.Commits[0].Message,
			Link:    gp.Commits[0].Url,
		},
		Sender: hook.User{
			UserName: gp.Sender.Login,
		},
	}
}
func convertPullRequestHook(gp *giteaPRHook) *hook.PullRequestHook {
	return &hook.PullRequestHook{
		Action: gp.Action,
		Repo: hook.Repository{
			Ref:         gp.PullRequest.Head.Ref,
			Sha:         gp.PullRequest.Head.Sha,
			CloneURL:    gp.PullRequest.Head.Repo.CloneUrl,
			CreatedAt:   gp.PullRequest.Head.Repo.CreatedAt,
			Branch:      gp.PullRequest.Head.Ref,
			Description: gp.PullRequest.Head.Repo.Description,
			FullName:    gp.PullRequest.Head.Repo.FullName,
			GitHttpURL:  gp.PullRequest.Head.Repo.CloneUrl,
			GitShhURL:   gp.PullRequest.Head.Repo.SshUrl,
			GitURL:      gp.PullRequest.Head.Repo.CloneUrl,
			HtmlURL:     gp.PullRequest.Head.Repo.HtmlUrl,
			SshURL:      gp.PullRequest.Head.Repo.SshUrl,
			Name:        gp.PullRequest.Head.Repo.Name,
			Private:     gp.PullRequest.Head.Repo.Private,
			URL:         gp.PullRequest.Head.Repo.HtmlUrl,
			Owner:       gp.PullRequest.Head.Repo.Owner.Login,
			RepoType:    "gitea",
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
			GitHttpURL:  gp.PullRequest.Base.Repo.CloneUrl,
			GitShhURL:   gp.PullRequest.Base.Repo.SshUrl,
			GitURL:      gp.PullRequest.Base.Repo.CloneUrl,
			HtmlURL:     gp.PullRequest.Base.Repo.HtmlUrl,
			SshURL:      gp.PullRequest.Base.Repo.SshUrl,
			Name:        gp.PullRequest.Base.Repo.Name,
			Private:     gp.PullRequest.Base.Repo.Private,
			URL:         gp.PullRequest.Base.Repo.HtmlUrl,
			Owner:       gp.PullRequest.Base.Repo.Owner.Login,
			RepoType:    "gitea",
			RepoOpenid:  "",
		},
		PullRequest: hook.PullRequest{
			Number: gp.Number,
			Body:   gp.PullRequest.Body,
			Title:  gp.PullRequest.Title,
			Base: hook.Reference{
				Name: gp.PullRequest.Base.Ref,
				Path: gp.PullRequest.Base.Repo.Name,
				Sha:  gp.PullRequest.Base.Sha,
			},
			Head: hook.Reference{
				Name: gp.PullRequest.Head.Ref,
				Path: gp.PullRequest.Head.Repo.Name,
				Sha:  gp.PullRequest.Head.Sha,
			},
			Author: hook.User{
				UserName: gp.PullRequest.User.Login,
			},
			Created: time.Time{},
			Updated: time.Time{},
		},
		Sender: hook.User{},
	}
}
func convertPullRequestURL(gc *giteaCommentHook) (*giteaPullRequestURL, error) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("convertPullRequestURL err:%+v", err)
			logrus.Warnf("%s", string(debug.Stack()))
		}
	}()
	//TODO 请求gitea api 获取pr 详情
	//client := cl.Repositories.(*giteaapi.RepositoryService)
	//tk := &pipeline.TUserToken{}
	//ok, _ := comm.DBMain.GetDB().Where("openid=?", gc.Sender.Id).Get(tk)
	//if !ok {
	//	logrus.Error("convertPullRequestURL not found user token")
	//	return nil, errors.New("convertPullRequestURL not found user token")
	//}
	//
	//quest, err := client.GetPullQuest(tk.AccessToken, gc.Repository.Owner.Login, gc.Repository.Name, gc.Issue.Number)
	//if err != nil {
	//	logrus.Errorf("convertPullRequestURL.GetPullQuest err : %v", err)
	//	return nil, err
	//}
	//requestURL := &giteaPullRequestURL{}
	//err = json.Unmarshal(quest, requestURL)
	//if err != nil {
	//	logrus.Errorf("gitea convertPullRequestURL Unmarshal err %v", err)
	//	return nil, err
	//}
	return nil, nil
}

func convertCommentHook(gp *giteaCommentHook) (*hook.PullRequestCommentHook, error) {
	pullRequestHook, err := convertPullRequestURL(gp)
	if err != nil {
		return nil, err
	}
	return &hook.PullRequestCommentHook{
		Action: hook.EVENTS_TYPE_COMMENT,
		Repo: hook.Repository{
			Ref:        pullRequestHook.Head.Ref,
			Sha:        pullRequestHook.Head.Sha,
			CloneURL:   pullRequestHook.Head.Repo.CloneUrl,
			CreatedAt:  pullRequestHook.Head.Repo.CreatedAt,
			Branch:     pullRequestHook.Head.Ref,
			FullName:   pullRequestHook.Head.Repo.FullName,
			GitHttpURL: pullRequestHook.Head.Repo.CloneUrl,
			GitShhURL:  pullRequestHook.Head.Repo.SshUrl,
			GitURL:     pullRequestHook.Head.Repo.CloneUrl,
			HtmlURL:    pullRequestHook.Head.Repo.HtmlUrl,
			SshURL:     pullRequestHook.Head.Repo.SshUrl,
			Name:       pullRequestHook.Head.Repo.Name,
			Private:    pullRequestHook.Head.Repo.Private,
			URL:        pullRequestHook.Head.Repo.HtmlUrl,
			Owner:      pullRequestHook.Head.Repo.Owner.Login,
			RepoType:   "gitea",
			RepoOpenid: strconv.Itoa(gp.Repository.Id),
		},
		TargetRepo: hook.Repository{
			Ref:        pullRequestHook.Base.Ref,
			Sha:        pullRequestHook.Base.Sha,
			CloneURL:   pullRequestHook.Base.Repo.CloneUrl,
			CreatedAt:  pullRequestHook.Base.Repo.CreatedAt,
			Branch:     pullRequestHook.Base.Ref,
			FullName:   pullRequestHook.Base.Repo.FullName,
			GitHttpURL: pullRequestHook.Base.Repo.CloneUrl,
			GitShhURL:  pullRequestHook.Base.Repo.SshUrl,
			GitURL:     pullRequestHook.Base.Repo.CloneUrl,
			HtmlURL:    pullRequestHook.Base.Repo.HtmlUrl,
			SshURL:     pullRequestHook.Base.Repo.SshUrl,
			Name:       pullRequestHook.Base.Repo.Name,
			Private:    pullRequestHook.Base.Repo.Private,
			URL:        pullRequestHook.Base.Repo.HtmlUrl,
			Owner:      pullRequestHook.Base.Repo.Owner.Login,
			RepoType:   "gitea",
			RepoOpenid: strconv.Itoa(gp.Repository.Id),
		},
		PullRequest: hook.PullRequest{
			Number: pullRequestHook.Number,
			Body:   pullRequestHook.Body,
			Title:  pullRequestHook.Title,
			Base: hook.Reference{
				Name: pullRequestHook.Base.Repo.Name,
				Path: pullRequestHook.Base.Repo.Name,
				Sha:  pullRequestHook.Base.Sha,
			},
			Head: hook.Reference{
				Name: pullRequestHook.Head.Repo.Name,
				Path: pullRequestHook.Head.Repo.Name,
				Sha:  pullRequestHook.Head.Sha,
			},
			Author: hook.User{
				UserName: pullRequestHook.User.Login,
			},
			Created: time.Time{},
			Updated: time.Time{},
		},
		Comment: hook.Comment{
			Body: gp.Comment.Body,
			Author: hook.User{
				UserName: gp.Comment.User.Login,
			},
		},
		Sender: hook.User{
			UserName: gp.Sender.Login,
		},
	}, nil
}

type giteaPushHook struct {
	Secret     string `json:"secret"`
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	CompareUrl string `json:"compare_url"`
	Commits    []struct {
		Id      string `json:"id"`
		Message string `json:"message"`
		Url     string `json:"url"`
		Author  struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Verification interface{} `json:"verification"`
		Timestamp    time.Time   `json:"timestamp"`
		Added        interface{} `json:"added"`
		Removed      interface{} `json:"removed"`
		Modified     interface{} `json:"modified"`
	} `json:"commits"`
	HeadCommit interface{} `json:"head_commit"`
	Repository struct {
		Id    int `json:"id"`
		Owner struct {
			Id            int       `json:"id"`
			Login         string    `json:"login"`
			FullName      string    `json:"full_name"`
			Email         string    `json:"email"`
			AvatarUrl     string    `json:"avatar_url"`
			Language      string    `json:"language"`
			IsAdmin       bool      `json:"is_admin"`
			LastLogin     time.Time `json:"last_login"`
			Created       time.Time `json:"created"`
			Restricted    bool      `json:"restricted"`
			Active        bool      `json:"active"`
			ProhibitLogin bool      `json:"prohibit_login"`
			Location      string    `json:"location"`
			Website       string    `json:"website"`
			Description   string    `json:"description"`
			Username      string    `json:"username"`
		} `json:"owner"`
		Name            string      `json:"name"`
		FullName        string      `json:"full_name"`
		Description     string      `json:"description"`
		Empty           bool        `json:"empty"`
		Private         bool        `json:"private"`
		Fork            bool        `json:"fork"`
		Template        bool        `json:"template"`
		Parent          interface{} `json:"parent"`
		Mirror          bool        `json:"mirror"`
		Size            int         `json:"size"`
		HtmlUrl         string      `json:"html_url"`
		SshUrl          string      `json:"ssh_url"`
		CloneUrl        string      `json:"clone_url"`
		OriginalUrl     string      `json:"original_url"`
		Website         string      `json:"website"`
		StarsCount      int         `json:"stars_count"`
		ForksCount      int         `json:"forks_count"`
		WatchersCount   int         `json:"watchers_count"`
		OpenIssuesCount int         `json:"open_issues_count"`
		OpenPrCounter   int         `json:"open_pr_counter"`
		ReleaseCounter  int         `json:"release_counter"`
		DefaultBranch   string      `json:"default_branch"`
		Archived        bool        `json:"archived"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		Permissions     struct {
			Admin bool `json:"admin"`
			Push  bool `json:"push"`
			Pull  bool `json:"pull"`
		} `json:"permissions"`
		HasIssues       bool `json:"has_issues"`
		InternalTracker struct {
			EnableTimeTracker                bool `json:"enable_time_tracker"`
			AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
			EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
		} `json:"internal_tracker"`
		HasWiki                   bool   `json:"has_wiki"`
		HasPullRequests           bool   `json:"has_pull_requests"`
		HasProjects               bool   `json:"has_projects"`
		IgnoreWhitespaceConflicts bool   `json:"ignore_whitespace_conflicts"`
		AllowMergeCommits         bool   `json:"allow_merge_commits"`
		AllowRebase               bool   `json:"allow_rebase"`
		AllowRebaseExplicit       bool   `json:"allow_rebase_explicit"`
		AllowSquashMerge          bool   `json:"allow_squash_merge"`
		DefaultMergeStyle         string `json:"default_merge_style"`
		AvatarUrl                 string `json:"avatar_url"`
		Internal                  bool   `json:"internal"`
		MirrorInterval            string `json:"mirror_interval"`
	} `json:"repository"`
	Pusher struct {
		Id            int       `json:"id"`
		Login         string    `json:"login"`
		FullName      string    `json:"full_name"`
		Email         string    `json:"email"`
		AvatarUrl     string    `json:"avatar_url"`
		Language      string    `json:"language"`
		IsAdmin       bool      `json:"is_admin"`
		LastLogin     time.Time `json:"last_login"`
		Created       time.Time `json:"created"`
		Restricted    bool      `json:"restricted"`
		Active        bool      `json:"active"`
		ProhibitLogin bool      `json:"prohibit_login"`
		Location      string    `json:"location"`
		Website       string    `json:"website"`
		Description   string    `json:"description"`
		Username      string    `json:"username"`
	} `json:"pusher"`
	Sender struct {
		Id            int       `json:"id"`
		Login         string    `json:"login"`
		FullName      string    `json:"full_name"`
		Email         string    `json:"email"`
		AvatarUrl     string    `json:"avatar_url"`
		Language      string    `json:"language"`
		IsAdmin       bool      `json:"is_admin"`
		LastLogin     time.Time `json:"last_login"`
		Created       time.Time `json:"created"`
		Restricted    bool      `json:"restricted"`
		Active        bool      `json:"active"`
		ProhibitLogin bool      `json:"prohibit_login"`
		Location      string    `json:"location"`
		Website       string    `json:"website"`
		Description   string    `json:"description"`
		Username      string    `json:"username"`
	} `json:"sender"`
}
type giteaPRHook struct {
	Secret      string `json:"secret"`
	Action      string `json:"action"`
	Number      int64  `json:"number"`
	PullRequest struct {
		Id     int    `json:"id"`
		Url    string `json:"url"`
		Number int    `json:"number"`
		User   struct {
			Id            int       `json:"id"`
			Login         string    `json:"login"`
			FullName      string    `json:"full_name"`
			Email         string    `json:"email"`
			AvatarUrl     string    `json:"avatar_url"`
			Language      string    `json:"language"`
			IsAdmin       bool      `json:"is_admin"`
			LastLogin     time.Time `json:"last_login"`
			Created       time.Time `json:"created"`
			Restricted    bool      `json:"restricted"`
			Active        bool      `json:"active"`
			ProhibitLogin bool      `json:"prohibit_login"`
			Location      string    `json:"location"`
			Website       string    `json:"website"`
			Description   string    `json:"description"`
			Username      string    `json:"username"`
		} `json:"user"`
		Title          string        `json:"title"`
		Body           string        `json:"body"`
		Labels         []interface{} `json:"labels"`
		Milestone      interface{}   `json:"milestone"`
		Assignee       interface{}   `json:"assignee"`
		Assignees      interface{}   `json:"assignees"`
		State          string        `json:"state"`
		IsLocked       bool          `json:"is_locked"`
		Comments       int           `json:"comments"`
		HtmlUrl        string        `json:"html_url"`
		DiffUrl        string        `json:"diff_url"`
		PatchUrl       string        `json:"patch_url"`
		Mergeable      bool          `json:"mergeable"`
		Merged         bool          `json:"merged"`
		MergedAt       interface{}   `json:"merged_at"`
		MergeCommitSha interface{}   `json:"merge_commit_sha"`
		MergedBy       interface{}   `json:"merged_by"`
		Base           struct {
			Label  string `json:"label"`
			Ref    string `json:"ref"`
			Sha    string `json:"sha"`
			RepoId int    `json:"repo_id"`
			Repo   struct {
				Id    int `json:"id"`
				Owner struct {
					Id            int       `json:"id"`
					Login         string    `json:"login"`
					FullName      string    `json:"full_name"`
					Email         string    `json:"email"`
					AvatarUrl     string    `json:"avatar_url"`
					Language      string    `json:"language"`
					IsAdmin       bool      `json:"is_admin"`
					LastLogin     time.Time `json:"last_login"`
					Created       time.Time `json:"created"`
					Restricted    bool      `json:"restricted"`
					Active        bool      `json:"active"`
					ProhibitLogin bool      `json:"prohibit_login"`
					Location      string    `json:"location"`
					Website       string    `json:"website"`
					Description   string    `json:"description"`
					Username      string    `json:"username"`
				} `json:"owner"`
				Name            string      `json:"name"`
				FullName        string      `json:"full_name"`
				Description     string      `json:"description"`
				Empty           bool        `json:"empty"`
				Private         bool        `json:"private"`
				Fork            bool        `json:"fork"`
				Template        bool        `json:"template"`
				Parent          interface{} `json:"parent"`
				Mirror          bool        `json:"mirror"`
				Size            int         `json:"size"`
				HtmlUrl         string      `json:"html_url"`
				SshUrl          string      `json:"ssh_url"`
				CloneUrl        string      `json:"clone_url"`
				OriginalUrl     string      `json:"original_url"`
				Website         string      `json:"website"`
				StarsCount      int         `json:"stars_count"`
				ForksCount      int         `json:"forks_count"`
				WatchersCount   int         `json:"watchers_count"`
				OpenIssuesCount int         `json:"open_issues_count"`
				OpenPrCounter   int         `json:"open_pr_counter"`
				ReleaseCounter  int         `json:"release_counter"`
				DefaultBranch   string      `json:"default_branch"`
				Archived        bool        `json:"archived"`
				CreatedAt       time.Time   `json:"created_at"`
				UpdatedAt       time.Time   `json:"updated_at"`
				Permissions     struct {
					Admin bool `json:"admin"`
					Push  bool `json:"push"`
					Pull  bool `json:"pull"`
				} `json:"permissions"`
				HasIssues       bool `json:"has_issues"`
				InternalTracker struct {
					EnableTimeTracker                bool `json:"enable_time_tracker"`
					AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
					EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
				} `json:"internal_tracker"`
				HasWiki                   bool   `json:"has_wiki"`
				HasPullRequests           bool   `json:"has_pull_requests"`
				HasProjects               bool   `json:"has_projects"`
				IgnoreWhitespaceConflicts bool   `json:"ignore_whitespace_conflicts"`
				AllowMergeCommits         bool   `json:"allow_merge_commits"`
				AllowRebase               bool   `json:"allow_rebase"`
				AllowRebaseExplicit       bool   `json:"allow_rebase_explicit"`
				AllowSquashMerge          bool   `json:"allow_squash_merge"`
				DefaultMergeStyle         string `json:"default_merge_style"`
				AvatarUrl                 string `json:"avatar_url"`
				Internal                  bool   `json:"internal"`
				MirrorInterval            string `json:"mirror_interval"`
			} `json:"repo"`
		} `json:"base"`
		Head struct {
			Label  string `json:"label"`
			Ref    string `json:"ref"`
			Sha    string `json:"sha"`
			RepoId int    `json:"repo_id"`
			Repo   struct {
				Id    int `json:"id"`
				Owner struct {
					Id            int       `json:"id"`
					Login         string    `json:"login"`
					FullName      string    `json:"full_name"`
					Email         string    `json:"email"`
					AvatarUrl     string    `json:"avatar_url"`
					Language      string    `json:"language"`
					IsAdmin       bool      `json:"is_admin"`
					LastLogin     time.Time `json:"last_login"`
					Created       time.Time `json:"created"`
					Restricted    bool      `json:"restricted"`
					Active        bool      `json:"active"`
					ProhibitLogin bool      `json:"prohibit_login"`
					Location      string    `json:"location"`
					Website       string    `json:"website"`
					Description   string    `json:"description"`
					Username      string    `json:"username"`
				} `json:"owner"`
				Name            string      `json:"name"`
				FullName        string      `json:"full_name"`
				Description     string      `json:"description"`
				Empty           bool        `json:"empty"`
				Private         bool        `json:"private"`
				Fork            bool        `json:"fork"`
				Template        bool        `json:"template"`
				Parent          interface{} `json:"parent"`
				Mirror          bool        `json:"mirror"`
				Size            int         `json:"size"`
				HtmlUrl         string      `json:"html_url"`
				SshUrl          string      `json:"ssh_url"`
				CloneUrl        string      `json:"clone_url"`
				OriginalUrl     string      `json:"original_url"`
				Website         string      `json:"website"`
				StarsCount      int         `json:"stars_count"`
				ForksCount      int         `json:"forks_count"`
				WatchersCount   int         `json:"watchers_count"`
				OpenIssuesCount int         `json:"open_issues_count"`
				OpenPrCounter   int         `json:"open_pr_counter"`
				ReleaseCounter  int         `json:"release_counter"`
				DefaultBranch   string      `json:"default_branch"`
				Archived        bool        `json:"archived"`
				CreatedAt       time.Time   `json:"created_at"`
				UpdatedAt       time.Time   `json:"updated_at"`
				Permissions     struct {
					Admin bool `json:"admin"`
					Push  bool `json:"push"`
					Pull  bool `json:"pull"`
				} `json:"permissions"`
				HasIssues       bool `json:"has_issues"`
				InternalTracker struct {
					EnableTimeTracker                bool `json:"enable_time_tracker"`
					AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
					EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
				} `json:"internal_tracker"`
				HasWiki                   bool   `json:"has_wiki"`
				HasPullRequests           bool   `json:"has_pull_requests"`
				HasProjects               bool   `json:"has_projects"`
				IgnoreWhitespaceConflicts bool   `json:"ignore_whitespace_conflicts"`
				AllowMergeCommits         bool   `json:"allow_merge_commits"`
				AllowRebase               bool   `json:"allow_rebase"`
				AllowRebaseExplicit       bool   `json:"allow_rebase_explicit"`
				AllowSquashMerge          bool   `json:"allow_squash_merge"`
				DefaultMergeStyle         string `json:"default_merge_style"`
				AvatarUrl                 string `json:"avatar_url"`
				Internal                  bool   `json:"internal"`
				MirrorInterval            string `json:"mirror_interval"`
			} `json:"repo"`
		} `json:"head"`
		MergeBase string      `json:"merge_base"`
		DueDate   interface{} `json:"due_date"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt time.Time   `json:"updated_at"`
		ClosedAt  interface{} `json:"closed_at"`
	} `json:"pull_request"`
	Repository struct {
		Id    int `json:"id"`
		Owner struct {
			Id            int       `json:"id"`
			Login         string    `json:"login"`
			FullName      string    `json:"full_name"`
			Email         string    `json:"email"`
			AvatarUrl     string    `json:"avatar_url"`
			Language      string    `json:"language"`
			IsAdmin       bool      `json:"is_admin"`
			LastLogin     time.Time `json:"last_login"`
			Created       time.Time `json:"created"`
			Restricted    bool      `json:"restricted"`
			Active        bool      `json:"active"`
			ProhibitLogin bool      `json:"prohibit_login"`
			Location      string    `json:"location"`
			Website       string    `json:"website"`
			Description   string    `json:"description"`
			Username      string    `json:"username"`
		} `json:"owner"`
		Name            string      `json:"name"`
		FullName        string      `json:"full_name"`
		Description     string      `json:"description"`
		Empty           bool        `json:"empty"`
		Private         bool        `json:"private"`
		Fork            bool        `json:"fork"`
		Template        bool        `json:"template"`
		Parent          interface{} `json:"parent"`
		Mirror          bool        `json:"mirror"`
		Size            int         `json:"size"`
		HtmlUrl         string      `json:"html_url"`
		SshUrl          string      `json:"ssh_url"`
		CloneUrl        string      `json:"clone_url"`
		OriginalUrl     string      `json:"original_url"`
		Website         string      `json:"website"`
		StarsCount      int         `json:"stars_count"`
		ForksCount      int         `json:"forks_count"`
		WatchersCount   int         `json:"watchers_count"`
		OpenIssuesCount int         `json:"open_issues_count"`
		OpenPrCounter   int         `json:"open_pr_counter"`
		ReleaseCounter  int         `json:"release_counter"`
		DefaultBranch   string      `json:"default_branch"`
		Archived        bool        `json:"archived"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		Permissions     struct {
			Admin bool `json:"admin"`
			Push  bool `json:"push"`
			Pull  bool `json:"pull"`
		} `json:"permissions"`
		HasIssues       bool `json:"has_issues"`
		InternalTracker struct {
			EnableTimeTracker                bool `json:"enable_time_tracker"`
			AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
			EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
		} `json:"internal_tracker"`
		HasWiki                   bool   `json:"has_wiki"`
		HasPullRequests           bool   `json:"has_pull_requests"`
		HasProjects               bool   `json:"has_projects"`
		IgnoreWhitespaceConflicts bool   `json:"ignore_whitespace_conflicts"`
		AllowMergeCommits         bool   `json:"allow_merge_commits"`
		AllowRebase               bool   `json:"allow_rebase"`
		AllowRebaseExplicit       bool   `json:"allow_rebase_explicit"`
		AllowSquashMerge          bool   `json:"allow_squash_merge"`
		DefaultMergeStyle         string `json:"default_merge_style"`
		AvatarUrl                 string `json:"avatar_url"`
		Internal                  bool   `json:"internal"`
		MirrorInterval            string `json:"mirror_interval"`
	} `json:"repository"`
	Sender struct {
		Id            int       `json:"id"`
		Login         string    `json:"login"`
		FullName      string    `json:"full_name"`
		Email         string    `json:"email"`
		AvatarUrl     string    `json:"avatar_url"`
		Language      string    `json:"language"`
		IsAdmin       bool      `json:"is_admin"`
		LastLogin     time.Time `json:"last_login"`
		Created       time.Time `json:"created"`
		Restricted    bool      `json:"restricted"`
		Active        bool      `json:"active"`
		ProhibitLogin bool      `json:"prohibit_login"`
		Location      string    `json:"location"`
		Website       string    `json:"website"`
		Description   string    `json:"description"`
		Username      string    `json:"username"`
	} `json:"sender"`
	Review interface{} `json:"review"`
}
type giteaCommentHook struct {
	Secret string `json:"secret"`
	Action string `json:"action"`
	Issue  struct {
		Id      int    `json:"id"`
		Url     string `json:"url"`
		HtmlUrl string `json:"html_url"`
		Number  int    `json:"number"`
		User    struct {
			Id            int       `json:"id"`
			Login         string    `json:"login"`
			FullName      string    `json:"full_name"`
			Email         string    `json:"email"`
			AvatarUrl     string    `json:"avatar_url"`
			Language      string    `json:"language"`
			IsAdmin       bool      `json:"is_admin"`
			LastLogin     time.Time `json:"last_login"`
			Created       time.Time `json:"created"`
			Restricted    bool      `json:"restricted"`
			Active        bool      `json:"active"`
			ProhibitLogin bool      `json:"prohibit_login"`
			Location      string    `json:"location"`
			Website       string    `json:"website"`
			Description   string    `json:"description"`
			Username      string    `json:"username"`
		} `json:"user"`
		OriginalAuthor   string        `json:"original_author"`
		OriginalAuthorId int           `json:"original_author_id"`
		Title            string        `json:"title"`
		Body             string        `json:"body"`
		Ref              string        `json:"ref"`
		Labels           []interface{} `json:"labels"`
		Milestone        interface{}   `json:"milestone"`
		Assignee         interface{}   `json:"assignee"`
		Assignees        interface{}   `json:"assignees"`
		State            string        `json:"state"`
		IsLocked         bool          `json:"is_locked"`
		Comments         int           `json:"comments"`
		CreatedAt        time.Time     `json:"created_at"`
		UpdatedAt        time.Time     `json:"updated_at"`
		ClosedAt         interface{}   `json:"closed_at"`
		DueDate          interface{}   `json:"due_date"`
		PullRequest      interface{}   `json:"pull_request"`
		Repository       struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			Owner    string `json:"owner"`
			FullName string `json:"full_name"`
		} `json:"repository"`
	} `json:"issue"`
	Comment struct {
		Id             int    `json:"id"`
		HtmlUrl        string `json:"html_url"`
		PullRequestUrl string `json:"pull_request_url"`
		IssueUrl       string `json:"issue_url"`
		User           struct {
			Id            int       `json:"id"`
			Login         string    `json:"login"`
			FullName      string    `json:"full_name"`
			Email         string    `json:"email"`
			AvatarUrl     string    `json:"avatar_url"`
			Language      string    `json:"language"`
			IsAdmin       bool      `json:"is_admin"`
			LastLogin     time.Time `json:"last_login"`
			Created       time.Time `json:"created"`
			Restricted    bool      `json:"restricted"`
			Active        bool      `json:"active"`
			ProhibitLogin bool      `json:"prohibit_login"`
			Location      string    `json:"location"`
			Website       string    `json:"website"`
			Description   string    `json:"description"`
			Username      string    `json:"username"`
		} `json:"user"`
		OriginalAuthor   string    `json:"original_author"`
		OriginalAuthorId int       `json:"original_author_id"`
		Body             string    `json:"body"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedAt        time.Time `json:"updated_at"`
	} `json:"comment"`
	Repository struct {
		Id    int `json:"id"`
		Owner struct {
			Id            int       `json:"id"`
			Login         string    `json:"login"`
			FullName      string    `json:"full_name"`
			Email         string    `json:"email"`
			AvatarUrl     string    `json:"avatar_url"`
			Language      string    `json:"language"`
			IsAdmin       bool      `json:"is_admin"`
			LastLogin     time.Time `json:"last_login"`
			Created       time.Time `json:"created"`
			Restricted    bool      `json:"restricted"`
			Active        bool      `json:"active"`
			ProhibitLogin bool      `json:"prohibit_login"`
			Location      string    `json:"location"`
			Website       string    `json:"website"`
			Description   string    `json:"description"`
			Username      string    `json:"username"`
		} `json:"owner"`
		Name            string      `json:"name"`
		FullName        string      `json:"full_name"`
		Description     string      `json:"description"`
		Empty           bool        `json:"empty"`
		Private         bool        `json:"private"`
		Fork            bool        `json:"fork"`
		Template        bool        `json:"template"`
		Parent          interface{} `json:"parent"`
		Mirror          bool        `json:"mirror"`
		Size            int         `json:"size"`
		HtmlUrl         string      `json:"html_url"`
		SshUrl          string      `json:"ssh_url"`
		CloneUrl        string      `json:"clone_url"`
		OriginalUrl     string      `json:"original_url"`
		Website         string      `json:"website"`
		StarsCount      int         `json:"stars_count"`
		ForksCount      int         `json:"forks_count"`
		WatchersCount   int         `json:"watchers_count"`
		OpenIssuesCount int         `json:"open_issues_count"`
		OpenPrCounter   int         `json:"open_pr_counter"`
		ReleaseCounter  int         `json:"release_counter"`
		DefaultBranch   string      `json:"default_branch"`
		Archived        bool        `json:"archived"`
		CreatedAt       time.Time   `json:"created_at"`
		UpdatedAt       time.Time   `json:"updated_at"`
		Permissions     struct {
			Admin bool `json:"admin"`
			Push  bool `json:"push"`
			Pull  bool `json:"pull"`
		} `json:"permissions"`
		HasIssues       bool `json:"has_issues"`
		InternalTracker struct {
			EnableTimeTracker                bool `json:"enable_time_tracker"`
			AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
			EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
		} `json:"internal_tracker"`
		HasWiki                   bool   `json:"has_wiki"`
		HasPullRequests           bool   `json:"has_pull_requests"`
		HasProjects               bool   `json:"has_projects"`
		IgnoreWhitespaceConflicts bool   `json:"ignore_whitespace_conflicts"`
		AllowMergeCommits         bool   `json:"allow_merge_commits"`
		AllowRebase               bool   `json:"allow_rebase"`
		AllowRebaseExplicit       bool   `json:"allow_rebase_explicit"`
		AllowSquashMerge          bool   `json:"allow_squash_merge"`
		DefaultMergeStyle         string `json:"default_merge_style"`
		AvatarUrl                 string `json:"avatar_url"`
		Internal                  bool   `json:"internal"`
		MirrorInterval            string `json:"mirror_interval"`
	} `json:"repository"`
	Sender struct {
		Id            int       `json:"id"`
		Login         string    `json:"login"`
		FullName      string    `json:"full_name"`
		Email         string    `json:"email"`
		AvatarUrl     string    `json:"avatar_url"`
		Language      string    `json:"language"`
		IsAdmin       bool      `json:"is_admin"`
		LastLogin     time.Time `json:"last_login"`
		Created       time.Time `json:"created"`
		Restricted    bool      `json:"restricted"`
		Active        bool      `json:"active"`
		ProhibitLogin bool      `json:"prohibit_login"`
		Location      string    `json:"location"`
		Website       string    `json:"website"`
		Description   string    `json:"description"`
		Username      string    `json:"username"`
	} `json:"sender"`
	IsPull bool `json:"is_pull"`
}
type giteaPullRequestURL struct {
	Id     int    `json:"id"`
	Url    string `json:"url"`
	Number int64  `json:"number"`
	User   struct {
		Id            int       `json:"id"`
		Login         string    `json:"login"`
		FullName      string    `json:"full_name"`
		Email         string    `json:"email"`
		AvatarUrl     string    `json:"avatar_url"`
		Language      string    `json:"language"`
		IsAdmin       bool      `json:"is_admin"`
		LastLogin     time.Time `json:"last_login"`
		Created       time.Time `json:"created"`
		Restricted    bool      `json:"restricted"`
		Active        bool      `json:"active"`
		ProhibitLogin bool      `json:"prohibit_login"`
		Location      string    `json:"location"`
		Website       string    `json:"website"`
		Description   string    `json:"description"`
		Username      string    `json:"username"`
	} `json:"user"`
	Title          string        `json:"title"`
	Body           string        `json:"body"`
	Labels         []interface{} `json:"labels"`
	Milestone      interface{}   `json:"milestone"`
	Assignee       interface{}   `json:"assignee"`
	Assignees      interface{}   `json:"assignees"`
	State          string        `json:"state"`
	IsLocked       bool          `json:"is_locked"`
	Comments       int           `json:"comments"`
	HtmlUrl        string        `json:"html_url"`
	DiffUrl        string        `json:"diff_url"`
	PatchUrl       string        `json:"patch_url"`
	Mergeable      bool          `json:"mergeable"`
	Merged         bool          `json:"merged"`
	MergedAt       interface{}   `json:"merged_at"`
	MergeCommitSha interface{}   `json:"merge_commit_sha"`
	MergedBy       interface{}   `json:"merged_by"`
	Base           struct {
		Label  string `json:"label"`
		Ref    string `json:"ref"`
		Sha    string `json:"sha"`
		RepoId int    `json:"repo_id"`
		Repo   struct {
			Id    int `json:"id"`
			Owner struct {
				Id            int       `json:"id"`
				Login         string    `json:"login"`
				FullName      string    `json:"full_name"`
				Email         string    `json:"email"`
				AvatarUrl     string    `json:"avatar_url"`
				Language      string    `json:"language"`
				IsAdmin       bool      `json:"is_admin"`
				LastLogin     time.Time `json:"last_login"`
				Created       time.Time `json:"created"`
				Restricted    bool      `json:"restricted"`
				Active        bool      `json:"active"`
				ProhibitLogin bool      `json:"prohibit_login"`
				Location      string    `json:"location"`
				Website       string    `json:"website"`
				Description   string    `json:"description"`
				Username      string    `json:"username"`
			} `json:"owner"`
			Name            string      `json:"name"`
			FullName        string      `json:"full_name"`
			Description     string      `json:"description"`
			Empty           bool        `json:"empty"`
			Private         bool        `json:"private"`
			Fork            bool        `json:"fork"`
			Template        bool        `json:"template"`
			Parent          interface{} `json:"parent"`
			Mirror          bool        `json:"mirror"`
			Size            int         `json:"size"`
			HtmlUrl         string      `json:"html_url"`
			SshUrl          string      `json:"ssh_url"`
			CloneUrl        string      `json:"clone_url"`
			OriginalUrl     string      `json:"original_url"`
			Website         string      `json:"website"`
			StarsCount      int         `json:"stars_count"`
			ForksCount      int         `json:"forks_count"`
			WatchersCount   int         `json:"watchers_count"`
			OpenIssuesCount int         `json:"open_issues_count"`
			OpenPrCounter   int         `json:"open_pr_counter"`
			ReleaseCounter  int         `json:"release_counter"`
			DefaultBranch   string      `json:"default_branch"`
			Archived        bool        `json:"archived"`
			CreatedAt       time.Time   `json:"created_at"`
			UpdatedAt       time.Time   `json:"updated_at"`
			Permissions     struct {
				Admin bool `json:"admin"`
				Push  bool `json:"push"`
				Pull  bool `json:"pull"`
			} `json:"permissions"`
			HasIssues       bool `json:"has_issues"`
			InternalTracker struct {
				EnableTimeTracker                bool `json:"enable_time_tracker"`
				AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
				EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
			} `json:"internal_tracker"`
			HasWiki                   bool   `json:"has_wiki"`
			HasPullRequests           bool   `json:"has_pull_requests"`
			HasProjects               bool   `json:"has_projects"`
			IgnoreWhitespaceConflicts bool   `json:"ignore_whitespace_conflicts"`
			AllowMergeCommits         bool   `json:"allow_merge_commits"`
			AllowRebase               bool   `json:"allow_rebase"`
			AllowRebaseExplicit       bool   `json:"allow_rebase_explicit"`
			AllowSquashMerge          bool   `json:"allow_squash_merge"`
			DefaultMergeStyle         string `json:"default_merge_style"`
			AvatarUrl                 string `json:"avatar_url"`
			Internal                  bool   `json:"internal"`
			MirrorInterval            string `json:"mirror_interval"`
		} `json:"repo"`
	} `json:"base"`
	Head struct {
		Label  string `json:"label"`
		Ref    string `json:"ref"`
		Sha    string `json:"sha"`
		RepoId int    `json:"repo_id"`
		Repo   struct {
			Id    int `json:"id"`
			Owner struct {
				Id            int       `json:"id"`
				Login         string    `json:"login"`
				FullName      string    `json:"full_name"`
				Email         string    `json:"email"`
				AvatarUrl     string    `json:"avatar_url"`
				Language      string    `json:"language"`
				IsAdmin       bool      `json:"is_admin"`
				LastLogin     time.Time `json:"last_login"`
				Created       time.Time `json:"created"`
				Restricted    bool      `json:"restricted"`
				Active        bool      `json:"active"`
				ProhibitLogin bool      `json:"prohibit_login"`
				Location      string    `json:"location"`
				Website       string    `json:"website"`
				Description   string    `json:"description"`
				Username      string    `json:"username"`
			} `json:"owner"`
			Name            string      `json:"name"`
			FullName        string      `json:"full_name"`
			Description     string      `json:"description"`
			Empty           bool        `json:"empty"`
			Private         bool        `json:"private"`
			Fork            bool        `json:"fork"`
			Template        bool        `json:"template"`
			Parent          interface{} `json:"parent"`
			Mirror          bool        `json:"mirror"`
			Size            int         `json:"size"`
			HtmlUrl         string      `json:"html_url"`
			SshUrl          string      `json:"ssh_url"`
			CloneUrl        string      `json:"clone_url"`
			OriginalUrl     string      `json:"original_url"`
			Website         string      `json:"website"`
			StarsCount      int         `json:"stars_count"`
			ForksCount      int         `json:"forks_count"`
			WatchersCount   int         `json:"watchers_count"`
			OpenIssuesCount int         `json:"open_issues_count"`
			OpenPrCounter   int         `json:"open_pr_counter"`
			ReleaseCounter  int         `json:"release_counter"`
			DefaultBranch   string      `json:"default_branch"`
			Archived        bool        `json:"archived"`
			CreatedAt       time.Time   `json:"created_at"`
			UpdatedAt       time.Time   `json:"updated_at"`
			Permissions     struct {
				Admin bool `json:"admin"`
				Push  bool `json:"push"`
				Pull  bool `json:"pull"`
			} `json:"permissions"`
			HasIssues       bool `json:"has_issues"`
			InternalTracker struct {
				EnableTimeTracker                bool `json:"enable_time_tracker"`
				AllowOnlyContributorsToTrackTime bool `json:"allow_only_contributors_to_track_time"`
				EnableIssueDependencies          bool `json:"enable_issue_dependencies"`
			} `json:"internal_tracker"`
			HasWiki                   bool   `json:"has_wiki"`
			HasPullRequests           bool   `json:"has_pull_requests"`
			HasProjects               bool   `json:"has_projects"`
			IgnoreWhitespaceConflicts bool   `json:"ignore_whitespace_conflicts"`
			AllowMergeCommits         bool   `json:"allow_merge_commits"`
			AllowRebase               bool   `json:"allow_rebase"`
			AllowRebaseExplicit       bool   `json:"allow_rebase_explicit"`
			AllowSquashMerge          bool   `json:"allow_squash_merge"`
			DefaultMergeStyle         string `json:"default_merge_style"`
			AvatarUrl                 string `json:"avatar_url"`
			Internal                  bool   `json:"internal"`
			MirrorInterval            string `json:"mirror_interval"`
		} `json:"repo"`
	} `json:"head"`
	MergeBase string      `json:"merge_base"`
	DueDate   interface{} `json:"due_date"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	ClosedAt  interface{} `json:"closed_at"`
}
