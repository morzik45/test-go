package handler

import (
	"net/http"

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
	mux.HandleFunc("/add", h.signUp)
	return mux
}

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.signIn(w, r)
	case "PUT":
		h.singOut(w, r)
	}
}
