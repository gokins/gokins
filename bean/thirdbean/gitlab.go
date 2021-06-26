package thirdbean

import "time"

type ResultGitlabRepo struct {
	Id                int           `json:"id"`
	Description       string        `json:"description"`
	Name              string        `json:"name"`
	NameWithNamespace string        `json:"name_with_namespace"`
	Path              string        `json:"path"`
	PathWithNamespace string        `json:"path_with_namespace"`
	CreatedAt         time.Time     `json:"created_at"`
	DefaultBranch     string        `json:"default_branch"`
	TagList           []interface{} `json:"tag_list"`
	Topics            []interface{} `json:"topics"`
	SshUrlToRepo      string        `json:"ssh_url_to_repo"`
	HttpUrlToRepo     string        `json:"http_url_to_repo"`
	WebUrl            string        `json:"web_url"`
	ReadmeUrl         string        `json:"readme_url"`
	AvatarUrl         interface{}   `json:"avatar_url"`
	ForksCount        int           `json:"forks_count"`
	StarCount         int           `json:"star_count"`
	LastActivityAt    time.Time     `json:"last_activity_at"`
	Namespace         struct {
		Id        int         `json:"id"`
		Name      string      `json:"name"`
		Path      string      `json:"path"`
		Kind      string      `json:"kind"`
		FullPath  string      `json:"full_path"`
		ParentId  interface{} `json:"parent_id"`
		AvatarUrl string      `json:"avatar_url"`
		WebUrl    string      `json:"web_url"`
	} `json:"namespace"`
	ContainerRegistryImagePrefix string `json:"container_registry_image_prefix"`
	Links                        struct {
		Self          string `json:"self"`
		Issues        string `json:"issues"`
		MergeRequests string `json:"merge_requests"`
		RepoBranches  string `json:"repo_branches"`
		Labels        string `json:"labels"`
		Events        string `json:"events"`
		Members       string `json:"members"`
	} `json:"_links"`
	PackagesEnabled bool   `json:"packages_enabled"`
	EmptyRepo       bool   `json:"empty_repo"`
	Archived        bool   `json:"archived"`
	Visibility      string `json:"visibility"`
	Owner           struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		State     string `json:"state"`
		AvatarUrl string `json:"avatar_url"`
		WebUrl    string `json:"web_url"`
	} `json:"owner"`
	ResolveOutdatedDiffDiscussions bool `json:"resolve_outdated_diff_discussions"`
	ContainerExpirationPolicy      struct {
		Cadence       string      `json:"cadence"`
		Enabled       bool        `json:"enabled"`
		KeepN         int         `json:"keep_n"`
		OlderThan     string      `json:"older_than"`
		NameRegex     string      `json:"name_regex"`
		NameRegexKeep interface{} `json:"name_regex_keep"`
		NextRunAt     time.Time   `json:"next_run_at"`
	} `json:"container_expiration_policy"`
	IssuesEnabled                             bool          `json:"issues_enabled"`
	MergeRequestsEnabled                      bool          `json:"merge_requests_enabled"`
	WikiEnabled                               bool          `json:"wiki_enabled"`
	JobsEnabled                               bool          `json:"jobs_enabled"`
	SnippetsEnabled                           bool          `json:"snippets_enabled"`
	ContainerRegistryEnabled                  bool          `json:"container_registry_enabled"`
	ServiceDeskEnabled                        bool          `json:"service_desk_enabled"`
	ServiceDeskAddress                        string        `json:"service_desk_address"`
	CanCreateMergeRequestIn                   bool          `json:"can_create_merge_request_in"`
	IssuesAccessLevel                         string        `json:"issues_access_level"`
	RepositoryAccessLevel                     string        `json:"repository_access_level"`
	MergeRequestsAccessLevel                  string        `json:"merge_requests_access_level"`
	ForkingAccessLevel                        string        `json:"forking_access_level"`
	WikiAccessLevel                           string        `json:"wiki_access_level"`
	BuildsAccessLevel                         string        `json:"builds_access_level"`
	SnippetsAccessLevel                       string        `json:"snippets_access_level"`
	PagesAccessLevel                          string        `json:"pages_access_level"`
	OperationsAccessLevel                     string        `json:"operations_access_level"`
	AnalyticsAccessLevel                      string        `json:"analytics_access_level"`
	EmailsDisabled                            interface{}   `json:"emails_disabled"`
	SharedRunnersEnabled                      bool          `json:"shared_runners_enabled"`
	LfsEnabled                                bool          `json:"lfs_enabled"`
	CreatorId                                 int           `json:"creator_id"`
	ImportStatus                              string        `json:"import_status"`
	OpenIssuesCount                           int           `json:"open_issues_count"`
	CiDefaultGitDepth                         int           `json:"ci_default_git_depth"`
	CiForwardDeploymentEnabled                bool          `json:"ci_forward_deployment_enabled"`
	PublicJobs                                bool          `json:"public_jobs"`
	BuildTimeout                              int           `json:"build_timeout"`
	AutoCancelPendingPipelines                string        `json:"auto_cancel_pending_pipelines"`
	BuildCoverageRegex                        interface{}   `json:"build_coverage_regex"`
	CiConfigPath                              string        `json:"ci_config_path"`
	SharedWithGroups                          []interface{} `json:"shared_with_groups"`
	OnlyAllowMergeIfPipelineSucceeds          bool          `json:"only_allow_merge_if_pipeline_succeeds"`
	AllowMergeOnSkippedPipeline               interface{}   `json:"allow_merge_on_skipped_pipeline"`
	RestrictUserDefinedVariables              bool          `json:"restrict_user_defined_variables"`
	RequestAccessEnabled                      bool          `json:"request_access_enabled"`
	OnlyAllowMergeIfAllDiscussionsAreResolved bool          `json:"only_allow_merge_if_all_discussions_are_resolved"`
	RemoveSourceBranchAfterMerge              bool          `json:"remove_source_branch_after_merge"`
	PrintingMergeRequestLinkEnabled           bool          `json:"printing_merge_request_link_enabled"`
	MergeMethod                               string        `json:"merge_method"`
	SuggestionCommitMessage                   interface{}   `json:"suggestion_commit_message"`
	AutoDevopsEnabled                         bool          `json:"auto_devops_enabled"`
	AutoDevopsDeployStrategy                  string        `json:"auto_devops_deploy_strategy"`
	AutocloseReferencedIssues                 bool          `json:"autoclose_referenced_issues"`
	ExternalAuthorizationClassificationLabel  string        `json:"external_authorization_classification_label"`
	RequirementsEnabled                       bool          `json:"requirements_enabled"`
	SecurityAndComplianceEnabled              bool          `json:"security_and_compliance_enabled"`
	ComplianceFrameworks                      []interface{} `json:"compliance_frameworks"`
	Permissions                               struct {
		ProjectAccess struct {
			AccessLevel       int `json:"access_level"`
			NotificationLevel int `json:"notification_level"`
		} `json:"project_access"`
		GroupAccess interface{} `json:"group_access"`
	} `json:"permissions"`
}

type ResultGitlabRepoBranch struct {
	Name               string `json:"name"`
	Merged             bool   `json:"merged"`
	Protected          bool   `json:"protected"`
	Default            bool   `json:"default"`
	DevelopersCanPush  bool   `json:"developers_can_push"`
	DevelopersCanMerge bool   `json:"developers_can_merge"`
	CanPush            bool   `json:"can_push"`
	WebUrl             string `json:"web_url"`
	Commit             struct {
		AuthorEmail    string    `json:"author_email"`
		AuthorName     string    `json:"author_name"`
		AuthoredDate   time.Time `json:"authored_date"`
		CommittedDate  time.Time `json:"committed_date"`
		CommitterEmail string    `json:"committer_email"`
		CommitterName  string    `json:"committer_name"`
		Id             string    `json:"id"`
		ShortId        string    `json:"short_id"`
		Title          string    `json:"title"`
		Message        string    `json:"message"`
		ParentIds      []string  `json:"parent_ids"`
	} `json:"commit"`
}

type ResultGetGitlabHook struct {
	Id                       int       `json:"id"`
	Url                      string    `json:"url"`
	CreatedAt                time.Time `json:"created_at"`
	PushEvents               bool      `json:"push_events"`
	TagPushEvents            bool      `json:"tag_push_events"`
	MergeRequestsEvents      bool      `json:"merge_requests_events"`
	RepositoryUpdateEvents   bool      `json:"repository_update_events"`
	EnableSslVerification    bool      `json:"enable_ssl_verification"`
	ProjectId                int       `json:"project_id"`
	IssuesEvents             bool      `json:"issues_events"`
	ConfidentialIssuesEvents bool      `json:"confidential_issues_events"`
	NoteEvents               bool      `json:"note_events"`
	ConfidentialNoteEvents   bool      `json:"confidential_note_events"`
	PipelineEvents           bool      `json:"pipeline_events"`
	WikiPageEvents           bool      `json:"wiki_page_events"`
	DeploymentEvents         bool      `json:"deployment_events"`
	JobEvents                bool      `json:"job_events"`
	ReleasesEvents           bool      `json:"releases_events"`
	PushEventsBranchFilter   string    `json:"push_events_branch_filter"`
}
