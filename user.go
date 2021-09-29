package exam

import (
	"database/sql"
	"time"
)

type User struct {
	Id       int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Authorization struct {
	Id           int          `json:"-"`
	Username     string       `json:"username"`
	IsAuthorized bool         `json:"is_authorized"`
	LoginAt      time.Time    `json:"login_at"`
	LogoutAt     sql.NullTime `json:"logout_at"`
	SessionToken string       `json:"session_token"`
}
