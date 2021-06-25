package bean

import "time"

type LogOutJson struct {
	Id      string    `json:"id"`
	Content string    `json:"content"`
	Times   time.Time `json:"times"`
	Errs    bool      `json:"errs"`
}
type LogOutJsonRes struct {
	Id      string    `json:"id"`
	Content string    `json:"content"`
	Times   time.Time `json:"times"`
	Errs    bool      `json:"errs"`

	Offset int64 `json:"offset"`
}
