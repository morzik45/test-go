package handler

import (
	"net/http"
	"runtime"

	exam "github.com/morzik45/test-go"
	"github.com/morzik45/test-go/logger"
	"github.com/morzik45/test-go/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/", panicRecovery(h.root))
	mux.HandleFunc("/api", panicRecovery(h.api))
	mux.HandleFunc("/add", h.signUp) // FIXME: not in task, for dev. dont forget to delete
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	muxWithAuth := authContext(h.services, mux)
	return muxWithAuth
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sessionRaw := r.Context().Value("Session")
		session, ok := sessionRaw.(*exam.Authorization)
		if ok && session.IsAuthorized {
			renderTestPage(w, r, session)
		} else {
			renderLoginForm(w)
		}
	}
}

func (h *Handler) api(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value("Session").(*exam.Authorization)

	if r.Method == "POST" {
		if !ok || !session.IsAuthorized {
			h.signIn(w, r)
			return
		} else {
			h.saveAnswer(w, r)
			return
		}

	}

	if r.Method == "PUT" {
		h.singOut(w, r)
		return
	}

	if r.Method == "GET" {
		if ok && session.IsAuthorized {
			h.getTasksHandler(w, r)
			return
		} else {
			logger.ERROR.Println("Unauthorized access attempt")
			http.Error(w, "Authorization is required!", http.StatusUnauthorized)
			return
		}
	}
}

func panicRecovery(h func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				n := runtime.Stack(buf, false)
				buf = buf[:n]

				logger.ERROR.Printf("recovering from err %v\n %s", err, buf)
				w.Write([]byte(`{"error":"our server got panic"}`))
			}
		}()

		h(w, r)
	}
}
