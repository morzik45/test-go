package handler

import (
	"encoding/json"
	"net/http"
	"path"
	"text/template"

	"github.com/lib/pq"
	exam "github.com/morzik45/test-go"
	"github.com/morzik45/test-go/logger"
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
			logger.ERROR.Printf("Singup without username or password: %s", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		user = exam.User{
			Username: usernames[0],
			Password: passwords[0],
		}
	}

	_, err := h.services.Authorization.CreateUser(user)
	if err != nil {
		pqe, ok := err.(*pq.Error)
		if ok && string(pqe.Code) == "23505" {
			logger.INFO.Printf("Try create user with username %s, but username already exist", user.Username)
		} else {
			logger.ERROR.Printf("Faild on create user: %s", err.Error())
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	logger.INFO.Printf("SingUp user %s", user.Username)
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
			logger.ERROR.Printf("Singin without username or password: %s", err.Error())
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
		logger.ERROR.Printf("Faild on login user: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		HttpOnly: true,
	})
	logger.INFO.Printf("SingIn user %s", user.Username)
	http.Redirect(w, r, "/test", http.StatusFound)
}

func (h *Handler) singOut(w http.ResponseWriter, r *http.Request, s *exam.Authorization) {
	_, err := h.services.Authorization.SessionClose(s.Username, s.SessionToken) // maybe need use only 'id' from session
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.ERROR.Printf("Error in sessionClose: %s", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
	})
	logger.INFO.Printf("Session closed by user %s", s.Username)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) userIdentity(r *http.Request) (*exam.Authorization, error) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	return h.services.Authorization.ParseToken(sessionToken.Value)
}

func renderLoginForm(w http.ResponseWriter) {
	fp := path.Join("static", "login.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		logger.ERROR.Printf("Error on render login form template: %s", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, ""); err != nil {
		logger.ERROR.Printf("Error on render login form template: %s", err.Error())

		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
