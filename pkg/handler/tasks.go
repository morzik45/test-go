package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"text/template"

	exam "github.com/morzik45/test-go"
	"github.com/morzik45/test-go/logger"
)

func renderTestPage(w http.ResponseWriter, r *http.Request, session *exam.Authorization) {
	// Change delims because conflict with vuejs
	if r.Method == "GET" {
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
}

func (h *Handler) saveAnswer(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("Session").(*exam.Authorization)
	response := make(map[string]interface{})
	var data exam.Answer

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		logger.ERROR.Printf("Error in answer params: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.INFO.Printf("User '%s' try save answer: variant: %d, task: %d, answer: %d, test: %d", session.Username, data.VariantID, data.TaskID, data.AnswerID, data.TestID)

	finished, percent, err := h.services.Testing.SaveAnswer(data, session.Username)
	if err != nil {
		logger.ERROR.Printf("Save answer return error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if finished {
		response["percent"] = percent
		response["status"] = "finished"
	} else {
		response["status"] = "ok"
	}
	js, err := json.Marshal(response)
	if err != nil {
		logger.ERROR.Printf("Error in marshaling all variants: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *Handler) getTasksHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := r.Context().Value("Session").(*exam.Authorization)
	variantIDstr := r.URL.Query().Get("variant_id")
	response := make(map[string]interface{})

	if len(variantIDstr) < 1 { // if not variant id return list variants
		variants, err := h.services.Testing.GetAllVariants()
		if err != nil {
			logger.ERROR.Printf("Error in getAllVariants: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			response["variants"] = variants
		}
	} else {
		variantID, errV := strconv.Atoi(variantIDstr)
		taskIDstr := r.URL.Query().Get("task_id")
		taskID, errT := strconv.Atoi(taskIDstr)
		if errT != nil || errV != nil {
			logger.ERROR.Printf("User '%s' request task with invalid variantID: '%s' or taskID: '%s'", session.Username, variantIDstr, taskIDstr)
			http.Error(w, "Invalid variant_id or task_id.", http.StatusInternalServerError)
			return
		}
		task, err := h.services.Testing.GetTaskById(variantID, taskID, session.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.ERROR.Printf("User '%s' request not exist question with variantID: '%s' and taskID: '%s'", session.Username, variantIDstr, taskIDstr)
				http.Error(w, "The question you asked does not exist.", http.StatusInternalServerError)
				return
			} else {
				logger.ERROR.Printf("Unknown error in request variantID '%s' and taskID: '%s' from DB, %s", variantIDstr, taskIDstr, err.Error())
				http.Error(w, fmt.Sprintf("Unknown error in request variantID '%s' and taskID: '%s' from DB", variantIDstr, taskIDstr), http.StatusInternalServerError)
				return
			}
		} else {
			response["question"] = map[string]interface{}{
				"id":         task.Id,
				"variant_id": task.VariantID,
				"test_id":    task.TestID,
				"question":   task.Question,
				"answers":    task.Answers,
			}
			logger.INFO.Printf("User '%s' request question with variantID: '%s' and taskID: '%s'", session.Username, variantIDstr, taskIDstr)
		}
	}
	js, err := json.Marshal(response)
	if err != nil {
		logger.ERROR.Printf("Error in marshaling all variants: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
