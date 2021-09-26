package exam

import (
	"database/sql"
	"time"
)

type User struct {
	Id       int    `json:"-"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Authorization struct {
	Id           int          `json:"-"`
	Username     string       `json:"username"`
	IsAuthorized bool         `json:"is_authorized"`
	LoginAt      time.Time    `json:"login_at"`
	LogoutAt     sql.NullTime `json:"logout_at"`
	SessionToken string       `json:"session_token"`
}
