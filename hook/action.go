package hook

type Action int

const (
	ActionOpen        = "open"
	ActionOpened      = "opened"
	ActionClose       = "close"
	ActionCreate      = "create"
	ActionDelete      = "delete"
	ActionSync        = "sync"
	ActionUpdate      = "update"
	ActionSynchronize = "synchronize"
)
