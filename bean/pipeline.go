package bean

type NewPipeline struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"displayName"`
	Content     string            `json:"content"`
	OrgId       string            `json:"orgId"`
	AccessToken string            `json:"accessToken"`
	Url         string            `json:"url"`
	Username    string            `json:"username"`
	Vars        []*NewPipelineVar `json:"vars"`
}

type NewPipelineVar struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Remarks string `json:"remarks"`
	Public  bool   `json:"public"`
}

func (p *NewPipeline) Check() bool {
	if p.Name == "" || p.Content == "" {
		return false
	}
	return true
}
