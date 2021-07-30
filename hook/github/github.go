package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gokins/gokins/hook"
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
	data, err := ioutil.ReadAll(
		io.LimitReader(req.Body, 10000000),
	)
	if err != nil {
		return nil, err
	}
	var wb hook.WebHook
	switch req.Header.Get(hook.GITHUB_EVENT) {
	case hook.GITHUB_EVENT_PUSH:
		wb, err = parsePushHook(data)
	case hook.GITHUB_EVENT_ISSUE_COMMENT:
		wb, err = parseCommentHook(data)
	case hook.GITHUB_EVENT_PR:
		wb, err = parsePullRequestHook(data)
	default:
		return nil, errors.New(fmt.Sprintf("hook含有未知的header:%v", req.Header.Get(hook.GITHUB_EVENT)))
	}
	if err != nil {
		return nil, err
	}
	sig := req.Header.Get("X-Hub-Signature")
	if !validatePrefix(data, []byte(secret), sig) {
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
	parts := strings.Split(signature, "=")
	if len(parts) != 2 {
		return false
	}
	switch parts[0] {
	case "sha1":
		return Validate(sha1.New, message, key, parts[1])
	case "sha256":
		return Validate(sha256.New, message, key, parts[1])
	default:
		return false
	}
}

func validate(h func() hash.Hash, message, key, signature []byte) bool {
	mac := hmac.New(h, key)
	mac.Write(message)
	sum := mac.Sum(nil)
	return hmac.Equal(signature, sum)
}

func parseCommentHook(data []byte) (*hook.PullRequestCommentHook, error) {
	gp := new(githubCommentHook)
	err := json.Unmarshal(data, gp)
	if err != nil {
		return nil, err
	}
	return convertCommentHook(gp)
}

func parsePullRequestHook(data []byte) (*hook.PullRequestHook, error) {
	gp := new(githubPRHook)
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
	gp := new(githubPushHook)
	err := json.Unmarshal(data, gp)
	if err != nil {
		return nil, err
	}
	return convertPushHook(gp), nil
}
func convertPullRequestHook(gp *githubPRHook) *hook.PullRequestHook {
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
			GitHttpURL:  gp.PullRequest.Head.Repo.GitUrl,
			GitShhURL:   gp.PullRequest.Head.Repo.SshUrl,
			GitSvnURL:   gp.PullRequest.Head.Repo.SvnUrl,
			GitURL:      gp.PullRequest.Head.Repo.GitUrl,
			HtmlURL:     gp.PullRequest.Head.Repo.HtmlUrl,
			SshURL:      gp.PullRequest.Head.Repo.SshUrl,
			SvnURL:      gp.PullRequest.Head.Repo.SvnUrl,
			Name:        gp.PullRequest.Head.Repo.Name,
			Private:     gp.PullRequest.Head.Repo.Private,
			URL:         gp.PullRequest.Head.Repo.Url,
			Owner:       gp.PullRequest.Head.Repo.Owner.Login,
			RepoType:    "github",
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
			GitHttpURL:  gp.PullRequest.Base.Repo.GitUrl,
			GitShhURL:   gp.PullRequest.Base.Repo.SshUrl,
			GitSvnURL:   gp.PullRequest.Base.Repo.SvnUrl,
			GitURL:      gp.PullRequest.Base.Repo.GitUrl,
			HtmlURL:     gp.PullRequest.Base.Repo.HtmlUrl,
			SshURL:      gp.PullRequest.Base.Repo.SshUrl,
			SvnURL:      gp.PullRequest.Base.Repo.SvnUrl,
			Name:        gp.PullRequest.Base.Repo.Name,
			Private:     gp.PullRequest.Base.Repo.Private,
			URL:         gp.PullRequest.Base.Repo.Url,
			Owner:       gp.PullRequest.Base.Repo.Owner.Login,
			RepoType:    "github",
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
func convertPushHook(gp *githubPushHook) *hook.PushHook {
	branch := gp.Ref
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
			CreatedAt:   time.Unix(gp.Repository.CreatedAt, 0),
			Branch:      branch,
			Description: gp.Repository.Description,
			FullName:    gp.Repository.FullName,
			GitHttpURL:  gp.Repository.GitUrl,
			GitShhURL:   gp.Repository.SshUrl,
			GitSvnURL:   gp.Repository.SvnUrl,
			GitURL:      gp.Repository.GitUrl,
			HtmlURL:     gp.Repository.HtmlUrl,
			SshURL:      gp.Repository.SshUrl,
			SvnURL:      gp.Repository.SvnUrl,
			Name:        gp.Repository.Name,
			Private:     gp.Repository.Private,
			URL:         gp.Repository.Url,
			Owner:       gp.Repository.Owner.Login,
			RepoType:    "github",
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
func convertPullRequestURL(u string) (*githubPullRequestURL, error) {
	client := &http.Client{
		Timeout: time.Second * 8,
	}
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		logrus.Errorf("github convertPullRequestURL CommentsUrl err %v", err)
		return nil, err
	}
	res, err := client.Do(request)
	if err != nil {
		logrus.Errorf("github convertPullRequestURL Do err %v", err)
		return nil, err
	}
	all, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		logrus.Errorf("github convertPullRequestURL ReadAll err %v", err)
		return nil, err
	}
	requestURL := &githubPullRequestURL{}
	err = json.Unmarshal(all, requestURL)
	if err != nil {
		logrus.Errorf("github convertPullRequestURL Unmarshal err %v", err)
		return nil, err
	}
	return requestURL, nil
}
func convertCommentHook(gp *githubCommentHook) (*hook.PullRequestCommentHook, error) {
	ul := gp.Issue.PullRequest.Url
	pullRequestHook, err := convertPullRequestURL(ul)
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
			GitHttpURL: pullRequestHook.Head.Repo.GitUrl,
			GitShhURL:  pullRequestHook.Head.Repo.SshUrl,
			GitSvnURL:  pullRequestHook.Head.Repo.SvnUrl,
			GitURL:     pullRequestHook.Head.Repo.GitUrl,
			HtmlURL:    pullRequestHook.Head.Repo.HtmlUrl,
			SshURL:     pullRequestHook.Head.Repo.SshUrl,
			SvnURL:     pullRequestHook.Head.Repo.SvnUrl,
			Name:       pullRequestHook.Head.Repo.Name,
			Private:    pullRequestHook.Head.Repo.Private,
			URL:        pullRequestHook.Head.Repo.Url,
			Owner:      pullRequestHook.Head.Repo.Owner.Login,
			RepoType:   "github",
			RepoOpenid: strconv.Itoa(gp.Repository.Id),
		},
		TargetRepo: hook.Repository{
			Ref:        pullRequestHook.Base.Ref,
			Sha:        pullRequestHook.Base.Sha,
			CloneURL:   pullRequestHook.Base.Repo.CloneUrl,
			CreatedAt:  pullRequestHook.Base.Repo.CreatedAt,
			Branch:     pullRequestHook.Base.Ref,
			FullName:   pullRequestHook.Base.Repo.FullName,
			GitHttpURL: pullRequestHook.Base.Repo.GitUrl,
			GitShhURL:  pullRequestHook.Base.Repo.SshUrl,
			GitSvnURL:  pullRequestHook.Base.Repo.SvnUrl,
			GitURL:     pullRequestHook.Base.Repo.GitUrl,
			HtmlURL:    pullRequestHook.Base.Repo.HtmlUrl,
			SshURL:     pullRequestHook.Base.Repo.SshUrl,
			SvnURL:     pullRequestHook.Base.Repo.SvnUrl,
			Name:       pullRequestHook.Base.Repo.Name,
			Private:    pullRequestHook.Base.Repo.Private,
			URL:        pullRequestHook.Base.Repo.Url,
			Owner:      pullRequestHook.Base.Repo.Owner.Login,
			RepoType:   "github",
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

type githubPullRequestURL struct {
	Url      string `json:"url"`
	Id       int64  `json:"id"`
	NodeId   string `json:"node_id"`
	HtmlUrl  string `json:"html_url"`
	DiffUrl  string `json:"diff_url"`
	PatchUrl string `json:"patch_url"`
	IssueUrl string `json:"issue_url"`
	Number   int64  `json:"number"`
	State    string `json:"state"`
	Locked   bool   `json:"locked"`
	Title    string `json:"title"`
	User     struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"user"`
	Body               string        `json:"body"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	ClosedAt           interface{}   `json:"closed_at"`
	MergedAt           interface{}   `json:"merged_at"`
	MergeCommitSha     string        `json:"merge_commit_sha"`
	Assignee           interface{}   `json:"assignee"`
	Assignees          []interface{} `json:"assignees"`
	RequestedReviewers []interface{} `json:"requested_reviewers"`
	RequestedTeams     []interface{} `json:"requested_teams"`
	Labels             []interface{} `json:"labels"`
	Milestone          interface{}   `json:"milestone"`
	Draft              bool          `json:"draft"`
	CommitsUrl         string        `json:"commits_url"`
	ReviewCommentsUrl  string        `json:"review_comments_url"`
	ReviewCommentUrl   string        `json:"review_comment_url"`
	CommentsUrl        string        `json:"comments_url"`
	StatusesUrl        string        `json:"statuses_url"`
	Head               struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Sha   string `json:"sha"`
		User  struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"user"`
		Repo struct {
			Id       int    `json:"id"`
			NodeId   string `json:"node_id"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			Private  bool   `json:"private"`
			Owner    struct {
				Login             string `json:"login"`
				Id                int    `json:"id"`
				NodeId            string `json:"node_id"`
				AvatarUrl         string `json:"avatar_url"`
				GravatarId        string `json:"gravatar_id"`
				Url               string `json:"url"`
				HtmlUrl           string `json:"html_url"`
				FollowersUrl      string `json:"followers_url"`
				FollowingUrl      string `json:"following_url"`
				GistsUrl          string `json:"gists_url"`
				StarredUrl        string `json:"starred_url"`
				SubscriptionsUrl  string `json:"subscriptions_url"`
				OrganizationsUrl  string `json:"organizations_url"`
				ReposUrl          string `json:"repos_url"`
				EventsUrl         string `json:"events_url"`
				ReceivedEventsUrl string `json:"received_events_url"`
				Type              string `json:"type"`
				SiteAdmin         bool   `json:"site_admin"`
			} `json:"owner"`
			HtmlUrl          string      `json:"html_url"`
			Description      interface{} `json:"description"`
			Fork             bool        `json:"fork"`
			Url              string      `json:"url"`
			ForksUrl         string      `json:"forks_url"`
			KeysUrl          string      `json:"keys_url"`
			CollaboratorsUrl string      `json:"collaborators_url"`
			TeamsUrl         string      `json:"teams_url"`
			HooksUrl         string      `json:"hooks_url"`
			IssueEventsUrl   string      `json:"issue_events_url"`
			EventsUrl        string      `json:"events_url"`
			AssigneesUrl     string      `json:"assignees_url"`
			BranchesUrl      string      `json:"branches_url"`
			TagsUrl          string      `json:"tags_url"`
			BlobsUrl         string      `json:"blobs_url"`
			GitTagsUrl       string      `json:"git_tags_url"`
			GitRefsUrl       string      `json:"git_refs_url"`
			TreesUrl         string      `json:"trees_url"`
			StatusesUrl      string      `json:"statuses_url"`
			LanguagesUrl     string      `json:"languages_url"`
			StargazersUrl    string      `json:"stargazers_url"`
			ContributorsUrl  string      `json:"contributors_url"`
			SubscribersUrl   string      `json:"subscribers_url"`
			SubscriptionUrl  string      `json:"subscription_url"`
			CommitsUrl       string      `json:"commits_url"`
			GitCommitsUrl    string      `json:"git_commits_url"`
			CommentsUrl      string      `json:"comments_url"`
			IssueCommentUrl  string      `json:"issue_comment_url"`
			ContentsUrl      string      `json:"contents_url"`
			CompareUrl       string      `json:"compare_url"`
			MergesUrl        string      `json:"merges_url"`
			ArchiveUrl       string      `json:"archive_url"`
			DownloadsUrl     string      `json:"downloads_url"`
			IssuesUrl        string      `json:"issues_url"`
			PullsUrl         string      `json:"pulls_url"`
			MilestonesUrl    string      `json:"milestones_url"`
			NotificationsUrl string      `json:"notifications_url"`
			LabelsUrl        string      `json:"labels_url"`
			ReleasesUrl      string      `json:"releases_url"`
			DeploymentsUrl   string      `json:"deployments_url"`
			CreatedAt        time.Time   `json:"created_at"`
			UpdatedAt        time.Time   `json:"updated_at"`
			PushedAt         time.Time   `json:"pushed_at"`
			GitUrl           string      `json:"git_url"`
			SshUrl           string      `json:"ssh_url"`
			CloneUrl         string      `json:"clone_url"`
			SvnUrl           string      `json:"svn_url"`
			Homepage         interface{} `json:"homepage"`
			Size             int         `json:"size"`
			StargazersCount  int         `json:"stargazers_count"`
			WatchersCount    int         `json:"watchers_count"`
			Language         string      `json:"language"`
			HasIssues        bool        `json:"has_issues"`
			HasProjects      bool        `json:"has_projects"`
			HasDownloads     bool        `json:"has_downloads"`
			HasWiki          bool        `json:"has_wiki"`
			HasPages         bool        `json:"has_pages"`
			ForksCount       int         `json:"forks_count"`
			MirrorUrl        interface{} `json:"mirror_url"`
			Archived         bool        `json:"archived"`
			Disabled         bool        `json:"disabled"`
			OpenIssuesCount  int         `json:"open_issues_count"`
			License          interface{} `json:"license"`
			Forks            int         `json:"forks"`
			OpenIssues       int         `json:"open_issues"`
			Watchers         int         `json:"watchers"`
			DefaultBranch    string      `json:"default_branch"`
		} `json:"repo"`
	} `json:"head"`
	Base struct {
		Label string `json:"label"`
		Ref   string `json:"ref"`
		Sha   string `json:"sha"`
		User  struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"user"`
		Repo struct {
			Id       int    `json:"id"`
			NodeId   string `json:"node_id"`
			Name     string `json:"name"`
			FullName string `json:"full_name"`
			Private  bool   `json:"private"`
			Owner    struct {
				Login             string `json:"login"`
				Id                int    `json:"id"`
				NodeId            string `json:"node_id"`
				AvatarUrl         string `json:"avatar_url"`
				GravatarId        string `json:"gravatar_id"`
				Url               string `json:"url"`
				HtmlUrl           string `json:"html_url"`
				FollowersUrl      string `json:"followers_url"`
				FollowingUrl      string `json:"following_url"`
				GistsUrl          string `json:"gists_url"`
				StarredUrl        string `json:"starred_url"`
				SubscriptionsUrl  string `json:"subscriptions_url"`
				OrganizationsUrl  string `json:"organizations_url"`
				ReposUrl          string `json:"repos_url"`
				EventsUrl         string `json:"events_url"`
				ReceivedEventsUrl string `json:"received_events_url"`
				Type              string `json:"type"`
				SiteAdmin         bool   `json:"site_admin"`
			} `json:"owner"`
			HtmlUrl          string      `json:"html_url"`
			Description      interface{} `json:"description"`
			Fork             bool        `json:"fork"`
			Url              string      `json:"url"`
			ForksUrl         string      `json:"forks_url"`
			KeysUrl          string      `json:"keys_url"`
			CollaboratorsUrl string      `json:"collaborators_url"`
			TeamsUrl         string      `json:"teams_url"`
			HooksUrl         string      `json:"hooks_url"`
			IssueEventsUrl   string      `json:"issue_events_url"`
			EventsUrl        string      `json:"events_url"`
			AssigneesUrl     string      `json:"assignees_url"`
			BranchesUrl      string      `json:"branches_url"`
			TagsUrl          string      `json:"tags_url"`
			BlobsUrl         string      `json:"blobs_url"`
			GitTagsUrl       string      `json:"git_tags_url"`
			GitRefsUrl       string      `json:"git_refs_url"`
			TreesUrl         string      `json:"trees_url"`
			StatusesUrl      string      `json:"statuses_url"`
			LanguagesUrl     string      `json:"languages_url"`
			StargazersUrl    string      `json:"stargazers_url"`
			ContributorsUrl  string      `json:"contributors_url"`
			SubscribersUrl   string      `json:"subscribers_url"`
			SubscriptionUrl  string      `json:"subscription_url"`
			CommitsUrl       string      `json:"commits_url"`
			GitCommitsUrl    string      `json:"git_commits_url"`
			CommentsUrl      string      `json:"comments_url"`
			IssueCommentUrl  string      `json:"issue_comment_url"`
			ContentsUrl      string      `json:"contents_url"`
			CompareUrl       string      `json:"compare_url"`
			MergesUrl        string      `json:"merges_url"`
			ArchiveUrl       string      `json:"archive_url"`
			DownloadsUrl     string      `json:"downloads_url"`
			IssuesUrl        string      `json:"issues_url"`
			PullsUrl         string      `json:"pulls_url"`
			MilestonesUrl    string      `json:"milestones_url"`
			NotificationsUrl string      `json:"notifications_url"`
			LabelsUrl        string      `json:"labels_url"`
			ReleasesUrl      string      `json:"releases_url"`
			DeploymentsUrl   string      `json:"deployments_url"`
			CreatedAt        time.Time   `json:"created_at"`
			UpdatedAt        time.Time   `json:"updated_at"`
			PushedAt         time.Time   `json:"pushed_at"`
			GitUrl           string      `json:"git_url"`
			SshUrl           string      `json:"ssh_url"`
			CloneUrl         string      `json:"clone_url"`
			SvnUrl           string      `json:"svn_url"`
			Homepage         interface{} `json:"homepage"`
			Size             int         `json:"size"`
			StargazersCount  int         `json:"stargazers_count"`
			WatchersCount    int         `json:"watchers_count"`
			Language         string      `json:"language"`
			HasIssues        bool        `json:"has_issues"`
			HasProjects      bool        `json:"has_projects"`
			HasDownloads     bool        `json:"has_downloads"`
			HasWiki          bool        `json:"has_wiki"`
			HasPages         bool        `json:"has_pages"`
			ForksCount       int         `json:"forks_count"`
			MirrorUrl        interface{} `json:"mirror_url"`
			Archived         bool        `json:"archived"`
			Disabled         bool        `json:"disabled"`
			OpenIssuesCount  int         `json:"open_issues_count"`
			License          interface{} `json:"license"`
			Forks            int         `json:"forks"`
			OpenIssues       int         `json:"open_issues"`
			Watchers         int         `json:"watchers"`
			DefaultBranch    string      `json:"default_branch"`
		} `json:"repo"`
	} `json:"base"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Html struct {
			Href string `json:"href"`
		} `json:"html"`
		Issue struct {
			Href string `json:"href"`
		} `json:"issue"`
		Comments struct {
			Href string `json:"href"`
		} `json:"comments"`
		ReviewComments struct {
			Href string `json:"href"`
		} `json:"review_comments"`
		ReviewComment struct {
			Href string `json:"href"`
		} `json:"review_comment"`
		Commits struct {
			Href string `json:"href"`
		} `json:"commits"`
		Statuses struct {
			Href string `json:"href"`
		} `json:"statuses"`
	} `json:"_links"`
	AuthorAssociation   string      `json:"author_association"`
	AutoMerge           interface{} `json:"auto_merge"`
	ActiveLockReason    interface{} `json:"active_lock_reason"`
	Merged              bool        `json:"merged"`
	Mergeable           bool        `json:"mergeable"`
	Rebaseable          bool        `json:"rebaseable"`
	MergeableState      string      `json:"mergeable_state"`
	MergedBy            interface{} `json:"merged_by"`
	Comments            int         `json:"comments"`
	ReviewComments      int         `json:"review_comments"`
	MaintainerCanModify bool        `json:"maintainer_can_modify"`
	Commits             int         `json:"commits"`
	Additions           int         `json:"additions"`
	Deletions           int         `json:"deletions"`
	ChangedFiles        int         `json:"changed_files"`
}

type githubPushHook struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Repository struct {
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Name              string `json:"name"`
			Email             string `json:"email"`
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HtmlUrl          string      `json:"html_url"`
		Description      string      `json:"description"`
		Fork             bool        `json:"fork"`
		Url              string      `json:"url"`
		ForksUrl         string      `json:"forks_url"`
		KeysUrl          string      `json:"keys_url"`
		CollaboratorsUrl string      `json:"collaborators_url"`
		TeamsUrl         string      `json:"teams_url"`
		HooksUrl         string      `json:"hooks_url"`
		IssueEventsUrl   string      `json:"issue_events_url"`
		EventsUrl        string      `json:"events_url"`
		AssigneesUrl     string      `json:"assignees_url"`
		BranchesUrl      string      `json:"branches_url"`
		TagsUrl          string      `json:"tags_url"`
		BlobsUrl         string      `json:"blobs_url"`
		GitTagsUrl       string      `json:"git_tags_url"`
		GitRefsUrl       string      `json:"git_refs_url"`
		TreesUrl         string      `json:"trees_url"`
		StatusesUrl      string      `json:"statuses_url"`
		LanguagesUrl     string      `json:"languages_url"`
		StargazersUrl    string      `json:"stargazers_url"`
		ContributorsUrl  string      `json:"contributors_url"`
		SubscribersUrl   string      `json:"subscribers_url"`
		SubscriptionUrl  string      `json:"subscription_url"`
		CommitsUrl       string      `json:"commits_url"`
		GitCommitsUrl    string      `json:"git_commits_url"`
		CommentsUrl      string      `json:"comments_url"`
		IssueCommentUrl  string      `json:"issue_comment_url"`
		ContentsUrl      string      `json:"contents_url"`
		CompareUrl       string      `json:"compare_url"`
		MergesUrl        string      `json:"merges_url"`
		ArchiveUrl       string      `json:"archive_url"`
		DownloadsUrl     string      `json:"downloads_url"`
		IssuesUrl        string      `json:"issues_url"`
		PullsUrl         string      `json:"pulls_url"`
		MilestonesUrl    string      `json:"milestones_url"`
		NotificationsUrl string      `json:"notifications_url"`
		LabelsUrl        string      `json:"labels_url"`
		ReleasesUrl      string      `json:"releases_url"`
		DeploymentsUrl   string      `json:"deployments_url"`
		CreatedAt        int64       `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		PushedAt         int         `json:"pushed_at"`
		GitUrl           string      `json:"git_url"`
		SshUrl           string      `json:"ssh_url"`
		CloneUrl         string      `json:"clone_url"`
		SvnUrl           string      `json:"svn_url"`
		Homepage         interface{} `json:"homepage"`
		Size             int         `json:"size"`
		StargazersCount  int         `json:"stargazers_count"`
		WatchersCount    int         `json:"watchers_count"`
		Language         string      `json:"language"`
		HasIssues        bool        `json:"has_issues"`
		HasProjects      bool        `json:"has_projects"`
		HasDownloads     bool        `json:"has_downloads"`
		HasWiki          bool        `json:"has_wiki"`
		HasPages         bool        `json:"has_pages"`
		ForksCount       int         `json:"forks_count"`
		MirrorUrl        interface{} `json:"mirror_url"`
		Archived         bool        `json:"archived"`
		Disabled         bool        `json:"disabled"`
		OpenIssuesCount  int         `json:"open_issues_count"`
		License          interface{} `json:"license"`
		Forks            int         `json:"forks"`
		OpenIssues       int         `json:"open_issues"`
		Watchers         int         `json:"watchers"`
		DefaultBranch    string      `json:"default_branch"`
		Stargazers       int         `json:"stargazers"`
		MasterBranch     string      `json:"master_branch"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Sender struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"sender"`
	Created bool        `json:"created"`
	Deleted bool        `json:"deleted"`
	Forced  bool        `json:"forced"`
	BaseRef interface{} `json:"base_ref"`
	Compare string      `json:"compare"`
	Commits []struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []interface{} `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []string      `json:"modified"`
	} `json:"commits"`
	HeadCommit struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []interface{} `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []string      `json:"modified"`
	} `json:"head_commit"`
}
type githubPRHook struct {
	Action      string `json:"action"`
	Number      int64  `json:"number"`
	PullRequest struct {
		Url      string `json:"url"`
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		HtmlUrl  string `json:"html_url"`
		DiffUrl  string `json:"diff_url"`
		PatchUrl string `json:"patch_url"`
		IssueUrl string `json:"issue_url"`
		Number   int    `json:"number"`
		State    string `json:"state"`
		Locked   bool   `json:"locked"`
		Title    string `json:"title"`
		User     struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"user"`
		Body               string        `json:"body"`
		CreatedAt          time.Time     `json:"created_at"`
		UpdatedAt          time.Time     `json:"updated_at"`
		ClosedAt           interface{}   `json:"closed_at"`
		MergedAt           interface{}   `json:"merged_at"`
		MergeCommitSha     interface{}   `json:"merge_commit_sha"`
		Assignee           interface{}   `json:"assignee"`
		Assignees          []interface{} `json:"assignees"`
		RequestedReviewers []interface{} `json:"requested_reviewers"`
		RequestedTeams     []interface{} `json:"requested_teams"`
		Labels             []interface{} `json:"labels"`
		Milestone          interface{}   `json:"milestone"`
		Draft              bool          `json:"draft"`
		CommitsUrl         string        `json:"commits_url"`
		ReviewCommentsUrl  string        `json:"review_comments_url"`
		ReviewCommentUrl   string        `json:"review_comment_url"`
		CommentsUrl        string        `json:"comments_url"`
		StatusesUrl        string        `json:"statuses_url"`
		Head               struct {
			Label string `json:"label"`
			Ref   string `json:"ref"`
			Sha   string `json:"sha"`
			User  struct {
				Login             string `json:"login"`
				Id                int    `json:"id"`
				NodeId            string `json:"node_id"`
				AvatarUrl         string `json:"avatar_url"`
				GravatarId        string `json:"gravatar_id"`
				Url               string `json:"url"`
				HtmlUrl           string `json:"html_url"`
				FollowersUrl      string `json:"followers_url"`
				FollowingUrl      string `json:"following_url"`
				GistsUrl          string `json:"gists_url"`
				StarredUrl        string `json:"starred_url"`
				SubscriptionsUrl  string `json:"subscriptions_url"`
				OrganizationsUrl  string `json:"organizations_url"`
				ReposUrl          string `json:"repos_url"`
				EventsUrl         string `json:"events_url"`
				ReceivedEventsUrl string `json:"received_events_url"`
				Type              string `json:"type"`
				SiteAdmin         bool   `json:"site_admin"`
			} `json:"user"`
			Repo struct {
				Id       int    `json:"id"`
				NodeId   string `json:"node_id"`
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				Private  bool   `json:"private"`
				Owner    struct {
					Login             string `json:"login"`
					Id                int    `json:"id"`
					NodeId            string `json:"node_id"`
					AvatarUrl         string `json:"avatar_url"`
					GravatarId        string `json:"gravatar_id"`
					Url               string `json:"url"`
					HtmlUrl           string `json:"html_url"`
					FollowersUrl      string `json:"followers_url"`
					FollowingUrl      string `json:"following_url"`
					GistsUrl          string `json:"gists_url"`
					StarredUrl        string `json:"starred_url"`
					SubscriptionsUrl  string `json:"subscriptions_url"`
					OrganizationsUrl  string `json:"organizations_url"`
					ReposUrl          string `json:"repos_url"`
					EventsUrl         string `json:"events_url"`
					ReceivedEventsUrl string `json:"received_events_url"`
					Type              string `json:"type"`
					SiteAdmin         bool   `json:"site_admin"`
				} `json:"owner"`
				HtmlUrl             string      `json:"html_url"`
				Description         string      `json:"description"`
				Fork                bool        `json:"fork"`
				Url                 string      `json:"url"`
				ForksUrl            string      `json:"forks_url"`
				KeysUrl             string      `json:"keys_url"`
				CollaboratorsUrl    string      `json:"collaborators_url"`
				TeamsUrl            string      `json:"teams_url"`
				HooksUrl            string      `json:"hooks_url"`
				IssueEventsUrl      string      `json:"issue_events_url"`
				EventsUrl           string      `json:"events_url"`
				AssigneesUrl        string      `json:"assignees_url"`
				BranchesUrl         string      `json:"branches_url"`
				TagsUrl             string      `json:"tags_url"`
				BlobsUrl            string      `json:"blobs_url"`
				GitTagsUrl          string      `json:"git_tags_url"`
				GitRefsUrl          string      `json:"git_refs_url"`
				TreesUrl            string      `json:"trees_url"`
				StatusesUrl         string      `json:"statuses_url"`
				LanguagesUrl        string      `json:"languages_url"`
				StargazersUrl       string      `json:"stargazers_url"`
				ContributorsUrl     string      `json:"contributors_url"`
				SubscribersUrl      string      `json:"subscribers_url"`
				SubscriptionUrl     string      `json:"subscription_url"`
				CommitsUrl          string      `json:"commits_url"`
				GitCommitsUrl       string      `json:"git_commits_url"`
				CommentsUrl         string      `json:"comments_url"`
				IssueCommentUrl     string      `json:"issue_comment_url"`
				ContentsUrl         string      `json:"contents_url"`
				CompareUrl          string      `json:"compare_url"`
				MergesUrl           string      `json:"merges_url"`
				ArchiveUrl          string      `json:"archive_url"`
				DownloadsUrl        string      `json:"downloads_url"`
				IssuesUrl           string      `json:"issues_url"`
				PullsUrl            string      `json:"pulls_url"`
				MilestonesUrl       string      `json:"milestones_url"`
				NotificationsUrl    string      `json:"notifications_url"`
				LabelsUrl           string      `json:"labels_url"`
				ReleasesUrl         string      `json:"releases_url"`
				DeploymentsUrl      string      `json:"deployments_url"`
				CreatedAt           time.Time   `json:"created_at"`
				UpdatedAt           time.Time   `json:"updated_at"`
				PushedAt            time.Time   `json:"pushed_at"`
				GitUrl              string      `json:"git_url"`
				SshUrl              string      `json:"ssh_url"`
				CloneUrl            string      `json:"clone_url"`
				SvnUrl              string      `json:"svn_url"`
				Homepage            interface{} `json:"homepage"`
				Size                int         `json:"size"`
				StargazersCount     int         `json:"stargazers_count"`
				WatchersCount       int         `json:"watchers_count"`
				Language            string      `json:"language"`
				HasIssues           bool        `json:"has_issues"`
				HasProjects         bool        `json:"has_projects"`
				HasDownloads        bool        `json:"has_downloads"`
				HasWiki             bool        `json:"has_wiki"`
				HasPages            bool        `json:"has_pages"`
				ForksCount          int         `json:"forks_count"`
				MirrorUrl           interface{} `json:"mirror_url"`
				Archived            bool        `json:"archived"`
				Disabled            bool        `json:"disabled"`
				OpenIssuesCount     int         `json:"open_issues_count"`
				License             interface{} `json:"license"`
				Forks               int         `json:"forks"`
				OpenIssues          int         `json:"open_issues"`
				Watchers            int         `json:"watchers"`
				DefaultBranch       string      `json:"default_branch"`
				AllowSquashMerge    bool        `json:"allow_squash_merge"`
				AllowMergeCommit    bool        `json:"allow_merge_commit"`
				AllowRebaseMerge    bool        `json:"allow_rebase_merge"`
				DeleteBranchOnMerge bool        `json:"delete_branch_on_merge"`
			} `json:"repo"`
		} `json:"head"`
		Base struct {
			Label string `json:"label"`
			Ref   string `json:"ref"`
			Sha   string `json:"sha"`
			User  struct {
				Login             string `json:"login"`
				Id                int    `json:"id"`
				NodeId            string `json:"node_id"`
				AvatarUrl         string `json:"avatar_url"`
				GravatarId        string `json:"gravatar_id"`
				Url               string `json:"url"`
				HtmlUrl           string `json:"html_url"`
				FollowersUrl      string `json:"followers_url"`
				FollowingUrl      string `json:"following_url"`
				GistsUrl          string `json:"gists_url"`
				StarredUrl        string `json:"starred_url"`
				SubscriptionsUrl  string `json:"subscriptions_url"`
				OrganizationsUrl  string `json:"organizations_url"`
				ReposUrl          string `json:"repos_url"`
				EventsUrl         string `json:"events_url"`
				ReceivedEventsUrl string `json:"received_events_url"`
				Type              string `json:"type"`
				SiteAdmin         bool   `json:"site_admin"`
			} `json:"user"`
			Repo struct {
				Id       int    `json:"id"`
				NodeId   string `json:"node_id"`
				Name     string `json:"name"`
				FullName string `json:"full_name"`
				Private  bool   `json:"private"`
				Owner    struct {
					Login             string `json:"login"`
					Id                int    `json:"id"`
					NodeId            string `json:"node_id"`
					AvatarUrl         string `json:"avatar_url"`
					GravatarId        string `json:"gravatar_id"`
					Url               string `json:"url"`
					HtmlUrl           string `json:"html_url"`
					FollowersUrl      string `json:"followers_url"`
					FollowingUrl      string `json:"following_url"`
					GistsUrl          string `json:"gists_url"`
					StarredUrl        string `json:"starred_url"`
					SubscriptionsUrl  string `json:"subscriptions_url"`
					OrganizationsUrl  string `json:"organizations_url"`
					ReposUrl          string `json:"repos_url"`
					EventsUrl         string `json:"events_url"`
					ReceivedEventsUrl string `json:"received_events_url"`
					Type              string `json:"type"`
					SiteAdmin         bool   `json:"site_admin"`
				} `json:"owner"`
				HtmlUrl             string      `json:"html_url"`
				Description         string      `json:"description"`
				Fork                bool        `json:"fork"`
				Url                 string      `json:"url"`
				ForksUrl            string      `json:"forks_url"`
				KeysUrl             string      `json:"keys_url"`
				CollaboratorsUrl    string      `json:"collaborators_url"`
				TeamsUrl            string      `json:"teams_url"`
				HooksUrl            string      `json:"hooks_url"`
				IssueEventsUrl      string      `json:"issue_events_url"`
				EventsUrl           string      `json:"events_url"`
				AssigneesUrl        string      `json:"assignees_url"`
				BranchesUrl         string      `json:"branches_url"`
				TagsUrl             string      `json:"tags_url"`
				BlobsUrl            string      `json:"blobs_url"`
				GitTagsUrl          string      `json:"git_tags_url"`
				GitRefsUrl          string      `json:"git_refs_url"`
				TreesUrl            string      `json:"trees_url"`
				StatusesUrl         string      `json:"statuses_url"`
				LanguagesUrl        string      `json:"languages_url"`
				StargazersUrl       string      `json:"stargazers_url"`
				ContributorsUrl     string      `json:"contributors_url"`
				SubscribersUrl      string      `json:"subscribers_url"`
				SubscriptionUrl     string      `json:"subscription_url"`
				CommitsUrl          string      `json:"commits_url"`
				GitCommitsUrl       string      `json:"git_commits_url"`
				CommentsUrl         string      `json:"comments_url"`
				IssueCommentUrl     string      `json:"issue_comment_url"`
				ContentsUrl         string      `json:"contents_url"`
				CompareUrl          string      `json:"compare_url"`
				MergesUrl           string      `json:"merges_url"`
				ArchiveUrl          string      `json:"archive_url"`
				DownloadsUrl        string      `json:"downloads_url"`
				IssuesUrl           string      `json:"issues_url"`
				PullsUrl            string      `json:"pulls_url"`
				MilestonesUrl       string      `json:"milestones_url"`
				NotificationsUrl    string      `json:"notifications_url"`
				LabelsUrl           string      `json:"labels_url"`
				ReleasesUrl         string      `json:"releases_url"`
				DeploymentsUrl      string      `json:"deployments_url"`
				CreatedAt           time.Time   `json:"created_at"`
				UpdatedAt           time.Time   `json:"updated_at"`
				PushedAt            time.Time   `json:"pushed_at"`
				GitUrl              string      `json:"git_url"`
				SshUrl              string      `json:"ssh_url"`
				CloneUrl            string      `json:"clone_url"`
				SvnUrl              string      `json:"svn_url"`
				Homepage            interface{} `json:"homepage"`
				Size                int         `json:"size"`
				StargazersCount     int         `json:"stargazers_count"`
				WatchersCount       int         `json:"watchers_count"`
				Language            string      `json:"language"`
				HasIssues           bool        `json:"has_issues"`
				HasProjects         bool        `json:"has_projects"`
				HasDownloads        bool        `json:"has_downloads"`
				HasWiki             bool        `json:"has_wiki"`
				HasPages            bool        `json:"has_pages"`
				ForksCount          int         `json:"forks_count"`
				MirrorUrl           interface{} `json:"mirror_url"`
				Archived            bool        `json:"archived"`
				Disabled            bool        `json:"disabled"`
				OpenIssuesCount     int         `json:"open_issues_count"`
				License             interface{} `json:"license"`
				Forks               int         `json:"forks"`
				OpenIssues          int         `json:"open_issues"`
				Watchers            int         `json:"watchers"`
				DefaultBranch       string      `json:"default_branch"`
				AllowSquashMerge    bool        `json:"allow_squash_merge"`
				AllowMergeCommit    bool        `json:"allow_merge_commit"`
				AllowRebaseMerge    bool        `json:"allow_rebase_merge"`
				DeleteBranchOnMerge bool        `json:"delete_branch_on_merge"`
			} `json:"repo"`
		} `json:"base"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Html struct {
				Href string `json:"href"`
			} `json:"html"`
			Issue struct {
				Href string `json:"href"`
			} `json:"issue"`
			Comments struct {
				Href string `json:"href"`
			} `json:"comments"`
			ReviewComments struct {
				Href string `json:"href"`
			} `json:"review_comments"`
			ReviewComment struct {
				Href string `json:"href"`
			} `json:"review_comment"`
			Commits struct {
				Href string `json:"href"`
			} `json:"commits"`
			Statuses struct {
				Href string `json:"href"`
			} `json:"statuses"`
		} `json:"_links"`
		AuthorAssociation   string      `json:"author_association"`
		AutoMerge           interface{} `json:"auto_merge"`
		ActiveLockReason    interface{} `json:"active_lock_reason"`
		Merged              bool        `json:"merged"`
		Mergeable           interface{} `json:"mergeable"`
		Rebaseable          interface{} `json:"rebaseable"`
		MergeableState      string      `json:"mergeable_state"`
		MergedBy            interface{} `json:"merged_by"`
		Comments            int         `json:"comments"`
		ReviewComments      int         `json:"review_comments"`
		MaintainerCanModify bool        `json:"maintainer_can_modify"`
		Commits             int         `json:"commits"`
		Additions           int         `json:"additions"`
		Deletions           int         `json:"deletions"`
		ChangedFiles        int         `json:"changed_files"`
	} `json:"pull_request"`
	Repository struct {
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HtmlUrl          string      `json:"html_url"`
		Description      string      `json:"description"`
		Fork             bool        `json:"fork"`
		Url              string      `json:"url"`
		ForksUrl         string      `json:"forks_url"`
		KeysUrl          string      `json:"keys_url"`
		CollaboratorsUrl string      `json:"collaborators_url"`
		TeamsUrl         string      `json:"teams_url"`
		HooksUrl         string      `json:"hooks_url"`
		IssueEventsUrl   string      `json:"issue_events_url"`
		EventsUrl        string      `json:"events_url"`
		AssigneesUrl     string      `json:"assignees_url"`
		BranchesUrl      string      `json:"branches_url"`
		TagsUrl          string      `json:"tags_url"`
		BlobsUrl         string      `json:"blobs_url"`
		GitTagsUrl       string      `json:"git_tags_url"`
		GitRefsUrl       string      `json:"git_refs_url"`
		TreesUrl         string      `json:"trees_url"`
		StatusesUrl      string      `json:"statuses_url"`
		LanguagesUrl     string      `json:"languages_url"`
		StargazersUrl    string      `json:"stargazers_url"`
		ContributorsUrl  string      `json:"contributors_url"`
		SubscribersUrl   string      `json:"subscribers_url"`
		SubscriptionUrl  string      `json:"subscription_url"`
		CommitsUrl       string      `json:"commits_url"`
		GitCommitsUrl    string      `json:"git_commits_url"`
		CommentsUrl      string      `json:"comments_url"`
		IssueCommentUrl  string      `json:"issue_comment_url"`
		ContentsUrl      string      `json:"contents_url"`
		CompareUrl       string      `json:"compare_url"`
		MergesUrl        string      `json:"merges_url"`
		ArchiveUrl       string      `json:"archive_url"`
		DownloadsUrl     string      `json:"downloads_url"`
		IssuesUrl        string      `json:"issues_url"`
		PullsUrl         string      `json:"pulls_url"`
		MilestonesUrl    string      `json:"milestones_url"`
		NotificationsUrl string      `json:"notifications_url"`
		LabelsUrl        string      `json:"labels_url"`
		ReleasesUrl      string      `json:"releases_url"`
		DeploymentsUrl   string      `json:"deployments_url"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		PushedAt         time.Time   `json:"pushed_at"`
		GitUrl           string      `json:"git_url"`
		SshUrl           string      `json:"ssh_url"`
		CloneUrl         string      `json:"clone_url"`
		SvnUrl           string      `json:"svn_url"`
		Homepage         interface{} `json:"homepage"`
		Size             int         `json:"size"`
		StargazersCount  int         `json:"stargazers_count"`
		WatchersCount    int         `json:"watchers_count"`
		Language         string      `json:"language"`
		HasIssues        bool        `json:"has_issues"`
		HasProjects      bool        `json:"has_projects"`
		HasDownloads     bool        `json:"has_downloads"`
		HasWiki          bool        `json:"has_wiki"`
		HasPages         bool        `json:"has_pages"`
		ForksCount       int         `json:"forks_count"`
		MirrorUrl        interface{} `json:"mirror_url"`
		Archived         bool        `json:"archived"`
		Disabled         bool        `json:"disabled"`
		OpenIssuesCount  int         `json:"open_issues_count"`
		License          interface{} `json:"license"`
		Forks            int         `json:"forks"`
		OpenIssues       int         `json:"open_issues"`
		Watchers         int         `json:"watchers"`
		DefaultBranch    string      `json:"default_branch"`
	} `json:"repository"`
	Sender struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"sender"`
}
type githubCommentHook struct {
	Action string `json:"action"`
	Issue  struct {
		Url           string `json:"url"`
		RepositoryUrl string `json:"repository_url"`
		LabelsUrl     string `json:"labels_url"`
		CommentsUrl   string `json:"comments_url"`
		EventsUrl     string `json:"events_url"`
		HtmlUrl       string `json:"html_url"`
		Id            int    `json:"id"`
		NodeId        string `json:"node_id"`
		Number        int    `json:"number"`
		Title         string `json:"title"`
		User          struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"user"`
		Labels            []interface{} `json:"labels"`
		State             string        `json:"state"`
		Locked            bool          `json:"locked"`
		Assignee          interface{}   `json:"assignee"`
		Assignees         []interface{} `json:"assignees"`
		Milestone         interface{}   `json:"milestone"`
		Comments          int           `json:"comments"`
		CreatedAt         time.Time     `json:"created_at"`
		UpdatedAt         time.Time     `json:"updated_at"`
		ClosedAt          interface{}   `json:"closed_at"`
		AuthorAssociation string        `json:"author_association"`
		ActiveLockReason  interface{}   `json:"active_lock_reason"`
		PullRequest       struct {
			Url      string `json:"url"`
			HtmlUrl  string `json:"html_url"`
			DiffUrl  string `json:"diff_url"`
			PatchUrl string `json:"patch_url"`
		} `json:"pull_request"`
		Body                  string      `json:"body"`
		PerformedViaGithubApp interface{} `json:"performed_via_github_app"`
	} `json:"issue"`
	Comment struct {
		Url      string `json:"url"`
		HtmlUrl  string `json:"html_url"`
		IssueUrl string `json:"issue_url"`
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		User     struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"user"`
		CreatedAt             time.Time   `json:"created_at"`
		UpdatedAt             time.Time   `json:"updated_at"`
		AuthorAssociation     string      `json:"author_association"`
		Body                  string      `json:"body"`
		PerformedViaGithubApp interface{} `json:"performed_via_github_app"`
	} `json:"comment"`
	Repository struct {
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		Owner    struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		HtmlUrl          string      `json:"html_url"`
		Description      string      `json:"description"`
		Fork             bool        `json:"fork"`
		Url              string      `json:"url"`
		ForksUrl         string      `json:"forks_url"`
		KeysUrl          string      `json:"keys_url"`
		CollaboratorsUrl string      `json:"collaborators_url"`
		TeamsUrl         string      `json:"teams_url"`
		HooksUrl         string      `json:"hooks_url"`
		IssueEventsUrl   string      `json:"issue_events_url"`
		EventsUrl        string      `json:"events_url"`
		AssigneesUrl     string      `json:"assignees_url"`
		BranchesUrl      string      `json:"branches_url"`
		TagsUrl          string      `json:"tags_url"`
		BlobsUrl         string      `json:"blobs_url"`
		GitTagsUrl       string      `json:"git_tags_url"`
		GitRefsUrl       string      `json:"git_refs_url"`
		TreesUrl         string      `json:"trees_url"`
		StatusesUrl      string      `json:"statuses_url"`
		LanguagesUrl     string      `json:"languages_url"`
		StargazersUrl    string      `json:"stargazers_url"`
		ContributorsUrl  string      `json:"contributors_url"`
		SubscribersUrl   string      `json:"subscribers_url"`
		SubscriptionUrl  string      `json:"subscription_url"`
		CommitsUrl       string      `json:"commits_url"`
		GitCommitsUrl    string      `json:"git_commits_url"`
		CommentsUrl      string      `json:"comments_url"`
		IssueCommentUrl  string      `json:"issue_comment_url"`
		ContentsUrl      string      `json:"contents_url"`
		CompareUrl       string      `json:"compare_url"`
		MergesUrl        string      `json:"merges_url"`
		ArchiveUrl       string      `json:"archive_url"`
		DownloadsUrl     string      `json:"downloads_url"`
		IssuesUrl        string      `json:"issues_url"`
		PullsUrl         string      `json:"pulls_url"`
		MilestonesUrl    string      `json:"milestones_url"`
		NotificationsUrl string      `json:"notifications_url"`
		LabelsUrl        string      `json:"labels_url"`
		ReleasesUrl      string      `json:"releases_url"`
		DeploymentsUrl   string      `json:"deployments_url"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		PushedAt         time.Time   `json:"pushed_at"`
		GitUrl           string      `json:"git_url"`
		SshUrl           string      `json:"ssh_url"`
		CloneUrl         string      `json:"clone_url"`
		SvnUrl           string      `json:"svn_url"`
		Homepage         interface{} `json:"homepage"`
		Size             int         `json:"size"`
		StargazersCount  int         `json:"stargazers_count"`
		WatchersCount    int         `json:"watchers_count"`
		Language         string      `json:"language"`
		HasIssues        bool        `json:"has_issues"`
		HasProjects      bool        `json:"has_projects"`
		HasDownloads     bool        `json:"has_downloads"`
		HasWiki          bool        `json:"has_wiki"`
		HasPages         bool        `json:"has_pages"`
		ForksCount       int         `json:"forks_count"`
		MirrorUrl        interface{} `json:"mirror_url"`
		Archived         bool        `json:"archived"`
		Disabled         bool        `json:"disabled"`
		OpenIssuesCount  int         `json:"open_issues_count"`
		License          interface{} `json:"license"`
		Forks            int         `json:"forks"`
		OpenIssues       int         `json:"open_issues"`
		Watchers         int         `json:"watchers"`
		DefaultBranch    string      `json:"default_branch"`
	} `json:"repository"`
	Sender struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"sender"`
}
