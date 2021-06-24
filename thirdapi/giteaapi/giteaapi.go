package giteaapi

//api 路径
const (
	BaseApiGitea = "https://api.gitea.com"
	/*
	   https://gitea.com/api/v5/repos/{owner}/{repo}/contents/{path}
	*/
	ApiGiteaCreateFile = "/repos/%s/%s/contents/%s"

	/*
	   /api/v1/users/jimbirthday/repos
	*/
	ApiGiteaGetRepos = "/user/repos?page=%v&limit=%v"

	/*
	*repos/{owner}/{repo}/hooks
	 */
	ApiGiteaCreateHooks = "/repos/%s/%s/hooks"
	/*
	   /repos/{owner}/{repo}/hooks
	*/
	ApiGiteaGetHooks = "/repos/%s/%s/hooks?page=%v&limit=%v"

	/*
	   /repos/{owner}/{repo}/hooks/{hook_id}
	*/
	ApiGiteaDeleteHooks = "/repos/%s/%s/hooks/%v"
	/*
	 https://gitea.com/api/v1/repos/{owner}/{repo}/branches
	*/
	ApiGiteaGetRepoBranches = "/repos/%s/%s/branches?page=1&limit=30"
	/*
	 */
	ApiGiteaGetPullRequest = "/repos/%s/%s/pulls/%v"
)
