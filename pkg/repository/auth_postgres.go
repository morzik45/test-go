package repository

import (
	"database/sql"
	"errors"
	"fmt"

	exam "github.com/morzik45/test-go"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user exam.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) values ($1, $2) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (exam.User, error) {
	var id int
	var user exam.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	row := r.db.QueryRow(query, username, password)
	if err := row.Scan(&id); err != nil {
		return user, err
	}
	user = exam.User{
		Id:       id,
		Username: username,
		Password: password,
	}
	return user, nil
}

func (r *AuthPostgres) LoginUser(username, session_token string) (int, error) {
	closeOldSessionsQuery := fmt.Sprintf("UPDATE %s SET is_authorized=FALSE, logout_at=NOW() WHERE username=$1 AND is_authorized=TRUE;", authorizationsTable)
	r.db.Exec(closeOldSessionsQuery, username)
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, session_token) values ($1, $2) RETURNING id", authorizationsTable)
	row := r.db.QueryRow(query, username, session_token)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) LogoutUser(username, session_token string) (int, error) {
	var id int
	query := fmt.Sprintf("UPDATE %s SET is_authorized=FALSE, logout_at=NOW() WHERE username=$1 AND is_authorized=TRUE AND session_token=$2 RETURNING id;", authorizationsTable)
	row := r.db.QueryRow(query, username, session_token)
	err := row.Scan(&id)
	switch err {
	case nil:
		return id, nil
	case sql.ErrNoRows:
		return 0, errors.New("session not found")
	default:
		return 0, err
	}
}

func (r *AuthPostgres) ParseToken(session_token string) (*exam.Authorization, error) {
	auth := exam.Authorization{SessionToken: session_token}
	query := fmt.Sprintf("SELECT id, username, is_authorized, login_at, logout_at FROM %s WHERE session_token = $1;", authorizationsTable)
	row := r.db.QueryRow(query, auth.SessionToken)
	if err := row.Scan(&auth.Id, &auth.Username, &auth.IsAuthorized, &auth.LoginAt, &auth.LogoutAt); err != nil {
		return nil, err
	}
	return &auth, nil
}
