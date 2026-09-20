package models

import "time"

// Question represents a single quiz question.
type Question struct {
	Q string   `json:"q"`
	A []string `json:"a"`
	C int      `json:"c"`
}

// Quiz represents a quiz template.
type Quiz struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Category        string     `json:"category"`
	TimePerQuestion int        `json:"timePerQuestion"`
	ScoringMode     string     `json:"scoringMode"`
	Status          string     `json:"status"`
	UpdatedAt       string     `json:"updatedAt"`
	Questions       []Question `json:"questions"`
}

// ParticipantAnswer represents an answer submitted by a participant for a question.
type ParticipantAnswer struct {
	SelectedOption int  `json:"selectedOption"`
	IsCorrect      bool `json:"isCorrect"`
	Points         int  `json:"points"`
}

// Participant represents a participant in a live quiz session.
type Participant struct {
	ID           string                       `json:"id,omitempty"`
	Name         string                       `json:"name"`
	Score        int                          `json:"score"`
	CorrectCount int                          `json:"correctCount"`
	SocketID     string                       `json:"socketId,omitempty"`
	Answers      map[int]*ParticipantAnswer   `json:"answers"`
}

// Session represents a live or completed quiz session.
type Session struct {
	Code                 string                  `json:"code"`
	QuizID               string                  `json:"quizId"`
	Title                string                  `json:"title"`
	Questions            []Question              `json:"questions"`
	TimePerQuestion      int                     `json:"timePerQuestion"`
	CurrentQuestionIndex int                     `json:"currentQuestionIndex"`
	Status               string                  `json:"status"` // LOBBY, LIVE, COMPLETED
	Participants         map[string]*Participant `json:"participants"`
	CreatedAt            time.Time               `json:"createdAt,omitempty"`
	UpdatedAt            time.Time               `json:"updatedAt,omitempty"`
}

// SessionSummary represents a high-level summary of a session for REST listing.
type SessionSummary struct {
	Code             string `json:"code"`
	Title            string `json:"title"`
	Status           string `json:"status"`
	ParticipantCount int    `json:"participantCount"`
	UpdatedAt        string `json:"updatedAt"`
}

// WSMessage represents generic incoming/outgoing WebSocket JSON messages.
type WSMessage struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload,omitempty"`
}

// Outgoing payloads for specific events
type RoomStatePayload struct {
	Session        *Session `json:"session"`
	ParticipantID  string   `json:"participantId,omitempty"`
	Distribution   []int    `json:"distribution"`
	AnsweredCount  int      `json:"answeredCount"`
}

type ParticipantsUpdatedPayload struct {
	Participants map[string]*Participant `json:"participants"`
	Count        int                     `json:"count"`
}

type QuestionChangedPayload struct {
	CurrentQuestionIndex int   `json:"currentQuestionIndex"`
	Distribution         []int `json:"distribution"`
	AnsweredCount        int   `json:"answeredCount"`
}

type LiveUpdatePayload struct {
	Participants map[string]*Participant `json:"participants"`
	Distribution []int                   `json:"distribution"`
	AnsweredCount int                    `json:"answeredCount"`
}

type AnswerResultPayload struct {
	QuestionIndex  int  `json:"questionIndex"`
	IsCorrect      bool `json:"isCorrect"`
	CorrectOption  int  `json:"correctOption"`
	MyScore        int  `json:"myScore"`
	MyCorrectCount int  `json:"myCorrectCount"`
}
