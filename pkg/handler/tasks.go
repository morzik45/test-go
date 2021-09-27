package handler

import (
	"encoding/json"
	"net/http"

	exam "github.com/morzik45/test-go"
	"github.com/morzik45/test-go/logger"
)

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
