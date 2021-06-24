package hook

import "time"

type (
	// PullRequest represents a repository pull request.
	PullRequest struct {
		Number  int64
		Title   string
		Body    string
		Base    Reference
		Head    Reference
		Author  User
		Created time.Time
		Updated time.Time
	}

	Comment struct {
		Body   string
		Author User
	}
)
