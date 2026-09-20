package main

import (
	"log"
	"net/http"
	"os"

	"elsa-real-time-quiz/backend/pkg/api"
	"elsa-real-time-quiz/backend/pkg/db"
	"elsa-real-time-quiz/backend/pkg/ws"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://elsa:elsapass@localhost:5432/elsa_quiz?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Database connection
	log.Println("[Server] Connecting to PostgreSQL...")
	_, err := db.InitDB(dbURL)
	if err != nil {
		log.Fatalf("[Server] Fatal error initializing database: %v", err)
	}

	// Start WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	mux := http.NewServeMux()

	// REST API Routes
	mux.HandleFunc("/health", corsMiddleware(api.HealthHandler))
	mux.HandleFunc("/api/quizzes", corsMiddleware(api.QuizzesHandler))
	mux.HandleFunc("/api/quizzes/", corsMiddleware(api.QuizzesHandler))
	mux.HandleFunc("/api/sessions", corsMiddleware(api.SessionsHandler))

	// WebSocket Route
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(hub, w, r)
	})

	log.Printf("[Server] ELSA Real-Time Quiz Go Backend listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("[Server] Server stopped unexpectedly: %v", err)
	}
}
