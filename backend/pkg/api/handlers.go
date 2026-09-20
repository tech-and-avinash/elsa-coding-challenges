package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"elsa-real-time-quiz/backend/pkg/db"
	"elsa-real-time-quiz/backend/pkg/models"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"db":     "connected",
	})
}

func QuizzesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) == 3 && pathParts[2] != "" {
			// GET /api/quizzes/:id
			id := pathParts[2]
			quiz, err := db.GetQuizByID(id)
			if err != nil || quiz == nil {
				http.Error(w, `{"error":"Quiz not found"}`, http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(quiz)
			return
		}
		if len(pathParts) == 4 && pathParts[2] != "" && pathParts[3] == "sessions" {
			// GET /api/quizzes/:id/sessions
			quizID := pathParts[2]
			sessions, err := db.GetSessionsByQuizID(quizID)
			if err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(sessions)
			return
		}

		quizzes, err := db.GetAllQuizzes()
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(quizzes)

	case http.MethodPost:
		var q models.Quiz
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			http.Error(w, `{"error":"Invalid request payload"}`, http.StatusBadRequest)
			return
		}
		if q.ID == "" {
			http.Error(w, `{"error":"Quiz ID required"}`, http.StatusBadRequest)
			return
		}
		if err := db.SaveQuiz(&q); err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(q)

	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func SessionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	list, err := db.GetAllSessions()
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(list)
}
