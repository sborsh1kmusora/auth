package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64        `db:"id"`
	UserInfo  UserInfo     `db:""`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type UserInfo struct {
	Username string `db:"username"`
	Password string `db:"password"`
	Role     string `db:"role"`
}

type UpdateUser struct {
	ID       int64
	Username *string
	Password *string
}
