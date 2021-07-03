package bean

type PipelineVar struct {
	Aid        int64  `json:"aid"`
	PipelineId string `json:"pipelineId"`
	Name       string ` json:"name"`
	Value      string ` json:"value"`
	Remarks    string ` json:"remarks"`
	Public     bool   ` json:"public"`
}
