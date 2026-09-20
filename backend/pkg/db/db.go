package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"elsa-real-time-quiz/backend/pkg/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes PostgreSQL connection pool and runs database schema migrations.
func InitDB(databaseURL string) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Retry loop for Postgres startup in container environment
	for i := 0; i < 10; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		log.Printf("[DB] Waiting for database connection... retry %d/10", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	log.Println("[DB] Successfully connected to PostgreSQL")

	if err := migrateSchema(); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	if err := seedDefaultData(); err != nil {
		log.Printf("[DB] Warning: seeding default data: %v", err)
	}

	return DB, nil
}

func migrateSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS quizzes (
		id VARCHAR(64) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		category VARCHAR(100),
		time_per_question INT DEFAULT 20,
		scoring_mode VARCHAR(100) DEFAULT 'Speed & Accuracy',
		status VARCHAR(50) DEFAULT 'Ready',
		updated_at VARCHAR(100),
		questions JSONB NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS sessions (
		code VARCHAR(32) PRIMARY KEY,
		quiz_id VARCHAR(64) REFERENCES quizzes(id) ON DELETE CASCADE,
		title VARCHAR(255) NOT NULL,
		time_per_question INT DEFAULT 20,
		current_question_index INT DEFAULT 0,
		status VARCHAR(50) DEFAULT 'LOBBY',
		questions JSONB NOT NULL,
		participants JSONB NOT NULL DEFAULT '{}'::jsonb,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := DB.Exec(schema)
	if err != nil {
		return err
	}
	log.Println("[DB] Database schema migration completed successfully")
	return nil
}

func seedDefaultData() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM quizzes").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // Already seeded
	}

	defaultQuizzes := []models.Quiz{
		{
			ID:              "quiz-1",
			Title:           "Business English Basics",
			Description:     "Everyday vocabulary for meetings, email, and workplace communication.",
			Category:        "Vocabulary · General",
			TimePerQuestion: 20,
			ScoringMode:     "Speed & Accuracy",
			Status:          "Ready",
			UpdatedAt:       "Today",
			Questions: []models.Question{
				{Q: "What does “concise” mean?", A: []string{"Brief and clear", "Difficult to understand", "Very detailed", "Informal"}, C: 0},
				{Q: "Which word means “to postpone”?", A: []string{"Confirm", "Defer", "Resolve", "Proceed"}, C: 1},
				{Q: "Choose the best phrase for a meeting follow-up.", A: []string{"Please find below", "Let's circle back", "I am agree", "Do the needful"}, C: 1},
				{Q: "What is the opposite of “expand”?", A: []string{"Increase", "Extend", "Reduce", "Explain"}, C: 2},
				{Q: "Which word means “necessary or essential”?", A: []string{"Optional", "Relevant", "Mandatory", "Casual"}, C: 2},
			},
		},
		{
			ID:              "quiz-2",
			Title:           "Travel Vocabulary",
			Description:     "Essential phrases for airports, hotels, and directions abroad.",
			Category:        "Vocabulary · B1",
			TimePerQuestion: 20,
			ScoringMode:     "Speed & Accuracy",
			Status:          "Ready",
			UpdatedAt:       "Yesterday",
			Questions: []models.Question{
				{Q: "What is an itinerary?", A: []string{"A travel plan", "A luggage bag", "A passport visa", "A flight delay"}, C: 0},
				{Q: "Which phrase is used at hotel check-in?", A: []string{"I have a reservation", "I want a refund", "Where is the exit", "Boarding pass please"}, C: 0},
				{Q: "What does 'layover' mean?", A: []string{"A stopover between flights", "Overbooked seat", "Lost luggage", "Direct flight"}, C: 0},
			},
		},
	}

	for _, q := range defaultQuizzes {
		if err := SaveQuiz(&q); err != nil {
			log.Printf("[DB] Error seeding quiz %s: %v", q.ID, err)
		}
	}

	// Seed room 4821
	initialParticipants := map[string]*models.Participant{
		"p-priya":  {Name: "Priya", Score: 980, CorrectCount: 5, Answers: make(map[int]*models.ParticipantAnswer)},
		"p-daniel": {Name: "Daniel", Score: 920, CorrectCount: 5, Answers: make(map[int]*models.ParticipantAnswer)},
		"p-maya":   {Name: "Maya", Score: 860, CorrectCount: 4, Answers: make(map[int]*models.ParticipantAnswer)},
		"p-noah":   {Name: "Noah", Score: 780, CorrectCount: 4, Answers: make(map[int]*models.ParticipantAnswer)},
		"p-sara":   {Name: "Sara", Score: 720, CorrectCount: 4, Answers: make(map[int]*models.ParticipantAnswer)},
		"p-liam":   {Name: "Liam", Score: 680, CorrectCount: 3, Answers: make(map[int]*models.ParticipantAnswer)},
		"p-emma":   {Name: "Emma", Score: 640, CorrectCount: 3, Answers: make(map[int]*models.ParticipantAnswer)},
		"p-ravi":   {Name: "Ravi", Score: 600, CorrectCount: 3, Answers: make(map[int]*models.ParticipantAnswer)},
	}

	demoSession := &models.Session{
		Code:                 "4821",
		QuizID:               "quiz-1",
		Title:                "Business English Basics",
		Questions:            defaultQuizzes[0].Questions,
		TimePerQuestion:      20,
		CurrentQuestionIndex: 0,
		Status:               "LOBBY",
		Participants:         initialParticipants,
	}

	if err := SaveSession(demoSession); err != nil {
		log.Printf("[DB] Error seeding demo session 4821: %v", err)
	}

	log.Println("[DB] Initial default data seeded successfully")
	return nil
}

// GetAllQuizzes retrieves all quizzes from the database.
func GetAllQuizzes() ([]models.Quiz, error) {
	rows, err := DB.Query("SELECT id, title, description, category, time_per_question, scoring_mode, status, updated_at, questions FROM quizzes ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Quiz
	for rows.Next() {
		var q models.Quiz
		var qJSON []byte
		err := rows.Scan(&q.ID, &q.Title, &q.Description, &q.Category, &q.TimePerQuestion, &q.ScoringMode, &q.Status, &q.UpdatedAt, &qJSON)
		if err != nil {
			return nil, err
		}
		if len(qJSON) > 0 {
			_ = json.Unmarshal(qJSON, &q.Questions)
		}
		list = append(list, q)
	}
	return list, nil
}

// GetQuizByID retrieves a specific quiz by ID.
func GetQuizByID(id string) (*models.Quiz, error) {
	row := DB.QueryRow("SELECT id, title, description, category, time_per_question, scoring_mode, status, updated_at, questions FROM quizzes WHERE id = $1", id)
	var q models.Quiz
	var qJSON []byte
	err := row.Scan(&q.ID, &q.Title, &q.Description, &q.Category, &q.TimePerQuestion, &q.ScoringMode, &q.Status, &q.UpdatedAt, &qJSON)
	if err != nil {
		return nil, err
	}
	if len(qJSON) > 0 {
		_ = json.Unmarshal(qJSON, &q.Questions)
	}
	return &q, nil
}

// SaveQuiz inserts or updates a quiz in the database.
func SaveQuiz(q *models.Quiz) error {
	qJSON, err := json.Marshal(q.Questions)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO quizzes (id, title, description, category, time_per_question, scoring_mode, status, updated_at, questions)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (id) DO UPDATE SET
		title = EXCLUDED.title,
		description = EXCLUDED.description,
		category = EXCLUDED.category,
		time_per_question = EXCLUDED.time_per_question,
		scoring_mode = EXCLUDED.scoring_mode,
		status = EXCLUDED.status,
		updated_at = EXCLUDED.updated_at,
		questions = EXCLUDED.questions;
	`
	_, err = DB.Exec(query, q.ID, q.Title, q.Description, q.Category, q.TimePerQuestion, q.ScoringMode, q.Status, q.UpdatedAt, qJSON)
	return err
}

// GetAllSessions retrieves summaries for all sessions.
func GetAllSessions() ([]models.SessionSummary, error) {
	rows, err := DB.Query("SELECT code, title, status, participants, updated_at FROM sessions ORDER BY updated_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SessionSummary
	for rows.Next() {
		var code, title, status, updatedAtStr string
		var pJSON []byte
		var updatedAt time.Time
		err := rows.Scan(&code, &title, &status, &pJSON, &updatedAt)
		if err != nil {
			// fallback if scan fails on string vs timestamp
			continue
		}
		var participants map[string]interface{}
		if len(pJSON) > 0 {
			_ = json.Unmarshal(pJSON, &participants)
		}
		if updatedAt.IsZero() {
			updatedAtStr = "Recently"
		} else {
			updatedAtStr = updatedAt.Format("2006-01-02 15:04")
		}

		list = append(list, models.SessionSummary{
			Code:             code,
			Title:            title,
			Status:           status,
			ParticipantCount: len(participants),
			UpdatedAt:        updatedAtStr,
		})
	}
	return list, nil
}

// GetSessionsByQuizID retrieves all sessions for a specific quiz.
func GetSessionsByQuizID(quizID string) ([]models.SessionSummary, error) {
	rows, err := DB.Query("SELECT code, title, status, participants, updated_at FROM sessions WHERE quiz_id = $1 ORDER BY updated_at DESC", quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.SessionSummary
	for rows.Next() {
		var code, title, status, updatedAtStr string
		var pJSON []byte
		var updatedAt time.Time
		err := rows.Scan(&code, &title, &status, &pJSON, &updatedAt)
		if err != nil {
			// fallback if scan fails on string vs timestamp
			continue
		}
		var participants map[string]interface{}
		if len(pJSON) > 0 {
			_ = json.Unmarshal(pJSON, &participants)
		}
		if updatedAt.IsZero() {
			updatedAtStr = "Recently"
		} else {
			updatedAtStr = updatedAt.Format("2006-01-02 15:04")
		}

		list = append(list, models.SessionSummary{
			Code:             code,
			Title:            title,
			Status:           status,
			ParticipantCount: len(participants),
			UpdatedAt:        updatedAtStr,
		})
	}
	return list, nil
}

// GetSessionByCode retrieves a full session from DB.
func GetSessionByCode(code string) (*models.Session, error) {
	row := DB.QueryRow("SELECT code, quiz_id, title, time_per_question, current_question_index, status, questions, participants FROM sessions WHERE code = $1", code)
	var s models.Session
	var qJSON, pJSON []byte
	err := row.Scan(&s.Code, &s.QuizID, &s.Title, &s.TimePerQuestion, &s.CurrentQuestionIndex, &s.Status, &qJSON, &pJSON)
	if err != nil {
		return nil, err
	}
	if len(qJSON) > 0 {
		_ = json.Unmarshal(qJSON, &s.Questions)
	}
	if len(pJSON) > 0 {
		_ = json.Unmarshal(pJSON, &s.Participants)
	}
	return &s, nil
}

// SaveSession updates or creates a session record in PostgreSQL.
func SaveSession(s *models.Session) error {
	qJSON, err := json.Marshal(s.Questions)
	if err != nil {
		return err
	}
	pJSON, err := json.Marshal(s.Participants)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO sessions (code, quiz_id, title, time_per_question, current_question_index, status, questions, participants, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP)
	ON CONFLICT (code) DO UPDATE SET
		quiz_id = EXCLUDED.quiz_id,
		title = EXCLUDED.title,
		time_per_question = EXCLUDED.time_per_question,
		current_question_index = EXCLUDED.current_question_index,
		status = EXCLUDED.status,
		questions = EXCLUDED.questions,
		participants = EXCLUDED.participants,
		updated_at = CURRENT_TIMESTAMP;
	`
	_, err = DB.Exec(query, s.Code, s.QuizID, s.Title, s.TimePerQuestion, s.CurrentQuestionIndex, s.Status, qJSON, pJSON)
	return err
}
