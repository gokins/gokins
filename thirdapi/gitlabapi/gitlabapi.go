package gitlabapi

//api 路径
const (
	BaseApiGitlab = "https://gitlab.com/api/v4"
	/*
	  repos/{owner}/{repo}/contents/{path}
	*/
	ApiGitlabCreateFile = "/repos/%s/%s/contents/%s"

	/*
	  /users/:user_id/projects
	*/
	ApiGitlabGetRepos = "/users/%s/projects?type=%v&page=%v&per_page=%v"

	/*
	  repos/{owner}/{repo}/hooks
	*/
	ApiGitlabCreateHooks = "/projects/%s/hooks"
	/*
	   /projects/:id/hooks
	*/
	ApiGitlabGetHooks = "/projects/%s/hooks?page=%v&per_page=%v"

	/*
	   /projects/:id/hooks/:hook_id
	*/
	ApiGitlabDeleteHooks = "/projects/%s/hooks/%v"
	/*
	  /projects/:id/repository/branches
	*/
	ApiGitlabGetRepoBranches = "/projects/%v/repository/branches"
)
