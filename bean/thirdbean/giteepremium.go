package thirdbean

import "time"

type ResultGiteePremiumCreateHooks struct {
	Id                  int         `json:"id"`
	Url                 string      `json:"url"`
	CreatedAt           time.Time   `json:"created_at"`
	Password            string      `json:"password"`
	ProjectId           int         `json:"project_id"`
	Result              string      `json:"result"`
	ResultCode          interface{} `json:"result_code"`
	PushEvents          bool        `json:"push_events"`
	TagPushEvents       bool        `json:"tag_push_events"`
	IssuesEvents        bool        `json:"issues_events"`
	NoteEvents          bool        `json:"note_events"`
	MergeRequestsEvents bool        `json:"merge_requests_events"`
}

type ResultGiteePremiumRepo struct {
	Id        int64  `json:"id"`
	FullName  string `json:"full_name"`
	HumanName string `json:"human_name"`
	Url       string `json:"url"`
	Namespace struct {
		Id      int    `json:"id"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Path    string `json:"path"`
		HtmlUrl string `json:"html_url"`
	} `json:"namespace"`
	Path  string `json:"path"`
	Name  string `json:"name"`
	Owner struct {
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		AvatarUrl         string `json:"avatar_url"`
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
	} `json:"owner"`
	Description         string      `json:"description"`
	Private             bool        `json:"private"`
	Public              bool        `json:"public"`
	Internal            bool        `json:"internal"`
	Fork                bool        `json:"fork"`
	HtmlUrl             string      `json:"html_url"`
	SshUrl              string      `json:"ssh_url"`
	ForksUrl            string      `json:"forks_url"`
	KeysUrl             string      `json:"keys_url"`
	CollaboratorsUrl    string      `json:"collaborators_url"`
	HooksUrl            string      `json:"hooks_url"`
	BranchesUrl         string      `json:"branches_url"`
	TagsUrl             string      `json:"tags_url"`
	BlobsUrl            string      `json:"blobs_url"`
	StargazersUrl       string      `json:"stargazers_url"`
	ContributorsUrl     string      `json:"contributors_url"`
	CommitsUrl          string      `json:"commits_url"`
	CommentsUrl         string      `json:"comments_url"`
	IssueCommentUrl     string      `json:"issue_comment_url"`
	IssuesUrl           string      `json:"issues_url"`
	PullsUrl            string      `json:"pulls_url"`
	MilestonesUrl       string      `json:"milestones_url"`
	NotificationsUrl    string      `json:"notifications_url"`
	LabelsUrl           string      `json:"labels_url"`
	ReleasesUrl         string      `json:"releases_url"`
	Recommend           bool        `json:"recommend"`
	Homepage            interface{} `json:"homepage"`
	Language            string      `json:"language"`
	ForksCount          int         `json:"forks_count"`
	StargazersCount     int         `json:"stargazers_count"`
	WatchersCount       int         `json:"watchers_count"`
	DefaultBranch       string      `json:"default_branch"`
	OpenIssuesCount     int         `json:"open_issues_count"`
	HasIssues           bool        `json:"has_issues"`
	HasWiki             bool        `json:"has_wiki"`
	IssueComment        bool        `json:"issue_comment"`
	CanComment          bool        `json:"can_comment"`
	PullRequestsEnabled bool        `json:"pull_requests_enabled"`
	HasPage             bool        `json:"has_page"`
	License             string      `json:"license"`
	Outsourced          bool        `json:"outsourced"`
	ProjectCreator      string      `json:"project_creator"`
	Members             []string    `json:"members"`
	PushedAt            time.Time   `json:"pushed_at"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	Parent              interface{} `json:"parent"`
	Paas                interface{} `json:"paas"`
	Stared              bool        `json:"stared"`
	Watched             bool        `json:"watched"`
	Permission          struct {
		Pull  bool `json:"pull"`
		Push  bool `json:"push"`
		Admin bool `json:"admin"`
	} `json:"permission"`
	Relation        string `json:"relation"`
	AssigneesNumber int    `json:"assignees_number"`
	TestersNumber   int    `json:"testers_number"`
	Assignees       []struct {
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		AvatarUrl         string `json:"avatar_url"`
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
	} `json:"assignees"`
	Testers []struct {
		Id                int    `json:"id"`
		Login             string `json:"login"`
		Name              string `json:"name"`
		AvatarUrl         string `json:"avatar_url"`
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
	} `json:"testers"`
}

type ResultGiteePremiumRepoBranch struct {
	Name   string `json:"name"`
	Commit struct {
		Sha string `json:"sha"`
		Url string `json:"url"`
	} `json:"commit"`
	Protected     bool   `json:"protected"`
	ProtectionUrl string `json:"protection_url"`
}

type ResultGetGiteePremiumHook struct {
	Id                  int       `json:"id"`
	Url                 string    `json:"url"`
	CreatedAt           time.Time `json:"created_at"`
	Password            string    `json:"password"`
	ProjectId           int       `json:"project_id"`
	Result              string    `json:"result"`
	ResultCode          int       `json:"result_code"`
	PushEvents          bool      `json:"push_events"`
	TagPushEvents       bool      `json:"tag_push_events"`
	IssuesEvents        bool      `json:"issues_events"`
	NoteEvents          bool      `json:"note_events"`
	MergeRequestsEvents bool      `json:"merge_requests_events"`
}
