package hook

import "time"

type Repository struct {
	Id          string    `json:"id"`
	Ref         string    `json:"ref"`
	Sha         string    `json:"sha"`
	CloneURL    string    `json:"cloneURL"`
	CreatedAt   time.Time `json:"createdAt"`
	Branch      string    `json:"branch"`
	Description string    `json:"description"`
	FullName    string    `json:"fullName"`
	GitHttpURL  string    `json:"gitHttpURL"`
	GitShhURL   string    `json:"gitSshURL"`
	GitSvnURL   string    `json:"gitSvnURL"`
	GitURL      string    `json:"gitURL"`
	HtmlURL     string    `json:"htmlURL"`
	SshURL      string    `json:"sshURL"`
	SvnURL      string    `json:"svnURL"`
	Name        string    `json:"name"`
	Private     bool      `json:"private"`
	URL         string    `json:"url"`
	Owner       string    `json:"owner"`
	RepoType    string    `json:"repoType"`
	RepoOpenid  string    `json:"repoOpenid"`
}
