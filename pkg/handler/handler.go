package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/morzik45/test-go/logger"
	"github.com/morzik45/test-go/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.root)
	mux.HandleFunc("/add", h.signUp) // FIXME: not in task, for dev. dont forget to delete
	return mux
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	session, err := h.userIdentity(r)
	if err == http.ErrNoCookie || err == sql.ErrNoRows {
		if r.Method == "POST" {
			h.signIn(w, r)
			return
		} else {
			logger.ERROR.Printf("Invalid session token: %s", err.Error())
			http.Error(w, "Authorization is required!", http.StatusUnauthorized)
			return
		}
	} else if err != nil {
		logger.ERROR.Printf("Error on user identity: %s", err.Error())
		http.Error(w, "Error on user identity", http.StatusInternalServerError)
		return
	}
	if !session.IsAuthorized {
		logger.INFO.Printf("Request from user '%s', but session closed at %s, try login", session.Username, session.LogoutAt.Time.Format(time.RFC3339))
		if r.Method == "POST" {
			h.signIn(w, r)
			return
		} else {
			http.Error(w, "Session is closed, relogin is required!", http.StatusUnauthorized)
			return
		}
	}
	if r.Method == "PUT" {
		h.singOut(w, r, session)
		return
	}
	if r.Method == "GET" {
		h.getAllVariants(w, r, session)
		return
	}

}
