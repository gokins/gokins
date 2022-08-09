package bean

type OrgVar struct {
	Aid     int64  `json:"aid"`
	OrgId   string `json:"orgId"`
	Name    string ` json:"name"`
	Value   string ` json:"value"`
	Remarks string ` json:"remarks"`
	Public  bool   ` json:"public"`
}
