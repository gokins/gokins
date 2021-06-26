package giteepremiumapi

//api 路径
const (
	BaseApiGiteePremium = "https://gitee.com/api/v5"

	/*
	   https://gitee.com/api/v5/repos/{owner}/{repo}/contents/{path}
	*/
	ApiGiteePremiumCreateFile = "/repos/%s/%s/contents/%v"

	/*
	   https://gitee.com/api/v5/user/repos?access_token={access_token}&visibility={visibility}&affiliation={affiliation}&type={type}&sort={sort}&direction={direction}&q={1}&page={page}&per_page={per_page}
	*/
	ApiGiteePremiumGetRepos = "/user/repos?access_token=%s&type=%v&sort=%v&direction=%v&page=%v&per_page=%v"

	/*
	   https://gitee.com/api/v5/repos/{owner}/{repo}/hooks
	*/
	ApiGiteePremiumCreateHooks = "/repos/%s/%s/hooks"
	/*
	   https://gitee.com/api/v5/repos/{owner}/{repo}/hooks
	*/
	ApiGiteePremiumGetHooks = "/repos/%s/%s/hooks?access_token=%s&page=%v&per_page=%v"

	/*
	   https://gitee.com/api/v5/repos/{owner}/{repo}/hooks/{id}
	*/
	ApiGiteePremiumDeleteHooks = "/repos/%s/%s/hooks/%v?access_token=%s"
	/*
	  https://gitee.com/api/v5/repos/{owner}/{repo}/branches
	*/
	ApiGiteePremiumGetRepoBranches = "/repos/%s/%s/branches?access_token=%s"
)
