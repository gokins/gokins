package thirdapi

import "time"

type Repository struct {
	Id        string `json:"id"`
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Namespace string `json:"namespace"`
	FullName  string `json:"fullName"`
	HtmlURL   string `json:"htmlURL"`
	RepoType  string `json:"-"`
}

type RepositoryPage struct {
	TotalPages int64
	Ropes      []*Repository
}

type RepositoryBranch struct {
	Name string `json:"name"`
}

type RepositoryHook struct {
	Id        int       `json:"id"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}
