package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserInfo struct {
	Name    string
	Email   string
	IsAdmin bool
}
