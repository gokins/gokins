package bean

type IdsRes struct {
	Id  string `json:"id"`
	Aid int64  `json:"aid"`
}
type LoginReq struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}
type LoginRes struct {
	Token         string `json:"token"`
	Id            string `json:"id"`
	Name          string `json:"name"`
	Nick          string `json:"nick"`
	Avatar        string `json:"avatar"`
	LastLoginTime string `json:"lastLoginTime"`
}

type LgInfoRes struct {
	Login     bool   `json:"login"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	Nick      string `json:"nick"`
	Avatar    string `json:"avatar"`
	LoginTime string `json:"login_time"`
	RegTime   string `json:"reg_time"`
}
