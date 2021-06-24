package hook

import (
	"net/http"
)

type (
	WebHook interface {
		Repository() Repository
	}
	WebhookService interface {
		Parse(req *http.Request, fn SecretFunc) (WebHook, error)
	}

	PushHook struct {
		Ref     string
		Repo    Repository
		Before  string
		After   string
		Commit  Commit
		Sender  User
		Commits []Commit
	}
	// BranchHook represents a branch or tag event,
	// eg create and delete github event types.
	BranchHook struct {
		Ref    Reference
		Repo   Repository
		Sender User
	}

	PullRequestHook struct {
		Action      string
		Repo        Repository
		TargetRepo  Repository
		PullRequest PullRequest
		Sender      User
	}

	PullRequestCommentHook struct {
		Action      string
		Repo        Repository
		TargetRepo  Repository
		PullRequest PullRequest
		Comment     Comment
		Sender      User
	}

	SecretFunc func(webhook WebHook) (string, error)
)

func (h *PushHook) Repository() Repository               { return h.Repo }
func (h *BranchHook) Repository() Repository             { return h.Repo }
func (h *PullRequestHook) Repository() Repository        { return h.Repo }
func (h *PullRequestCommentHook) Repository() Repository { return h.Repo }
