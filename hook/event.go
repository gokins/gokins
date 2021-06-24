package hook

//触发事件
const (
	// EVENTS_TYPE_COMMENT 评论事件
	EVENTS_TYPE_COMMENT = "comment"
	// EVENTS_TYPE_PR pull request事件
	EVENTS_TYPE_PR = "pr"
	// EVENTS_TYPE_PUSH push事件
	EVENTS_TYPE_PUSH = "push"
	// EVENTS_TYPE_BUILD 手动运行事件
	EVENTS_TYPE_BUILD = "build"
	// EVENTS_TYPE_REBUILD 手动重新构建
	EVENTS_TYPE_REBUILD = "rebuild"
)

const (
	GITEE_EVENT                   = "X-Gitee-Event"
	GITEE_EVENT_PUSH              = "Push Hook"
	GITEE_EVENT_NOTE              = "Note Hook"
	GITEE_EVENT_PR                = "Merge Request Hook"
	GITEE_EVENT_PR_ACTION_OPEN    = "open"
	GITEE_EVENT_PR_ACTION_UPDATE  = "update"
	GITEE_EVENT_PR_ACTION_COMMENT = "comment"
)
const (
	GITHUB_EVENT                   = "X-GitHub-Event"
	GITHUB_EVENT_ISSUE_COMMENT     = "issue_comment"
	GITHUB_EVENT_PUSH              = "push"
	GITHUB_EVENT_PR                = "pull_request"
	GITHUB_EVENT_PR_ACTION_OPEN    = "open"
	GITHUB_EVENT_PR_ACTION_UPDATE  = "update"
	GITHUB_EVENT_PR_ACTION_COMMENT = "comment"
)

const (
	GITLAB_EVENT      = "X-Gitlab-Event"
	GITLAB_EVENT_PUSH = "Push Hook"
	GITLAB_EVENT_PR   = "Merge Request Hook"
	GITLAB_EVENT_NOTE = "Note Hook"
)

const (
	GITEA_EVENT      = "X-Gitea-Event"
	GITEA_EVENT_PUSH = "push"
	GITEA_EVENT_PR   = "pull_request"
	GITEA_EVENT_NOTE = "issue_comment"
)
