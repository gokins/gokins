package comm

import "time"

type LogOutJson struct {
	Id      string    `json:"id"`
	Content string    `json:"content"`
	Times   time.Time `json:"times"`
	Errs    bool      `json:"errs"`
}
