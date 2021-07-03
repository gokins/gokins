package bean

type PipelineShow struct {
	Id           string `json:"id"`
	Uid          string `json:"uid"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	PipelineType string `json:"pipelineType"`
	YmlContent   string `json:"ymlContent"`
	Url          string `json:"url"`
	Username     string `json:"username"`
	AccessToken  string `json:"accessToken"`
}
