package githubapi

//api 路径
const (
	BaseApiGithub = "https://api.github.com"
	/*
	   repos/{owner}/{repo}/contents/{path}
	*/
	ApiGithubCreateFile = "/repos/%s/%s/contents/%s"

	/*
	   /user/repos
	*/
	ApiGithubGetRepos = "/user/repos?type=%v&sort=%v&direction=%v&page=%v&per_page=%v"

	/*
	   {owner}/{repo}/hooks
	*/
	ApiGithubCreateHooks = "/repos/%s/%s/hooks"
	/*
	   /repos/{owner}/{repo}/hooks
	*/
	ApiGithubGetHooks = "/repos/%s/%s/hooks?page=%v&per_page=%v"

	/*
	   /repos/{owner}/{repo}/hooks/{hook_id}
	*/
	ApiGithubDeleteHooks = "/repos/%s/%s/hooks/%v"
	/*
	  repos/{owner}/{repo}/branches
	*/
	ApiGithubGetRepoBranches = "/repos/%s/%s/branches?per_page=100"
)
