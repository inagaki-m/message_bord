package model

import "time"

type MessageInfo struct {
	Name       string    `db:"name"`
	Message    string    `db:"message"`
	CreateTime time.Time `db:"createTime"`
}
