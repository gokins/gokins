package thirdapi

type (
	RepositoryService interface {
		GetRepos(accessToken, username, types, sort, direction string, page, per_page int) (*RepositoryPage, error)

		DeleteHooks(accessToken, owner, repo, hookId string) error

		CreateWebHooks(accessToken, owner, repo, backUrl, password string) (*RepositoryHook, error)

		GetRepoBranches(accessToken, owner, repo string) ([]*RepositoryBranch, error)

		GetWebHooks(accessToken, owner, repo string, page, per_page int) ([]*RepositoryHook, error)
	}
)
