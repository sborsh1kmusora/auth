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

type CreateInput struct {
	UserInfo        UserInfo
	Password        string
	PasswordConfirm string
}

type UpdateInput struct {
	ID    int64
	Name  *string
	Email *string
}
