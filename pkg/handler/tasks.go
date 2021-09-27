package handler

import (
	"encoding/json"
	"net/http"
	"path"
	"text/template"

	exam "github.com/morzik45/test-go"
	"github.com/morzik45/test-go/logger"
)

func (h *Handler) testPage(w http.ResponseWriter, r *http.Request) {
	// Change delims because conflict with vuejs
	tmpl, err := template.New("test.html").Delims("{#", "#}").ParseFiles(path.Join("static", "test.html"))
	if err != nil {
		logger.ERROR.Printf("Error on render list template: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, ""); err != nil {
		logger.ERROR.Printf("Error on render list template: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) getAllVariants(w http.ResponseWriter, r *http.Request, a *exam.Authorization) {
	variants, err := h.services.Testing.GetAllVariants()
	if err != nil {
		logger.ERROR.Printf("Error in getAllVariants: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(variants)
	if err != nil {
		logger.ERROR.Printf("Error in marshaling all variants: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
