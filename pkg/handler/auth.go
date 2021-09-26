package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	exam "github.com/morzik45/test-go"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	user := exam.User{}
	r.ParseForm()
	usernames := r.Form["username"]
	passwords := r.Form["password"]

	if len(usernames) < 1 || len(passwords) < 1 {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		user = exam.User{
			Username: usernames[0],
			Password: passwords[0],
		}
	}

	log.Printf("SingUp user %s with password %s", user.Username, user.Password)

	_, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		log.Printf("Faild on create user: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var user exam.User
	r.ParseForm()
	usernames := r.Form["username"]
	passwords := r.Form["password"]

	if len(usernames) < 1 || len(passwords) < 1 {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		user = exam.User{
			Username: usernames[0],
			Password: passwords[0],
		}
	}

	token, err := h.services.Authorization.GenerateToken(user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) singOut(w http.ResponseWriter, r *http.Request) {
	session, err := h.userIdentity(r)
	if err == http.ErrNoCookie || err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error in userIdentity: %s", err.Error())
		return
	}
	if !session.IsAuthorized {
		log.Println(session.IsAuthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = h.services.Authorization.SessionClose(session.Username, session.SessionToken) // maybe need use only 'id' from session
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Error in sessionClose: %s", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) userIdentity(r *http.Request) (*exam.Authorization, error) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	log.Printf("Login with session token: %s", sessionToken.Value)
	return h.services.Authorization.ParseToken(sessionToken.Value)
}
