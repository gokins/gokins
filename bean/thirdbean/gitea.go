package thirdbean

import "time"

type ResultGiteaRepo struct {
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
}
type ResultGiteaRepoBranch struct {
	Name   string `json:"name"`
	Commit struct {
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
		Verification struct {
			Verified  bool        `json:"verified"`
			Reason    string      `json:"reason"`
			Signature string      `json:"signature"`
			Signer    interface{} `json:"signer"`
			Payload   string      `json:"payload"`
		} `json:"verification"`
		Timestamp time.Time   `json:"timestamp"`
		Added     interface{} `json:"added"`
		Removed   interface{} `json:"removed"`
		Modified  interface{} `json:"modified"`
	} `json:"commit"`
	Protected                     bool          `json:"protected"`
	RequiredApprovals             int           `json:"required_approvals"`
	EnableStatusCheck             bool          `json:"enable_status_check"`
	StatusCheckContexts           []interface{} `json:"status_check_contexts"`
	UserCanPush                   bool          `json:"user_can_push"`
	UserCanMerge                  bool          `json:"user_can_merge"`
	EffectiveBranchProtectionName string        `json:"effective_branch_protection_name"`
}
type ResultGetGiteaHook struct {
	Id     int    `json:"id"`
	Type   string `json:"type"`
	Config struct {
		ContentType string `json:"content_type"`
		Url         string `json:"url"`
	} `json:"config"`
	Events    []string  `json:"events"`
	Active    bool      `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
