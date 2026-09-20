package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"elsa-real-time-quiz/backend/pkg/db"
	"elsa-real-time-quiz/backend/pkg/models"
)

type Hub struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	rooms      map[string]map[*Client]bool // roomCode -> clients
	sessions   map[string]*models.Session  // roomCode -> in-memory active session
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]map[*Client]bool),
		sessions:   make(map[string]*models.Session),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.RoomCode != "" && h.rooms[client.RoomCode] != nil {
					delete(h.rooms[client.RoomCode], client)
					if len(h.rooms[client.RoomCode]) == 0 {
						delete(h.rooms, client.RoomCode)
					}
				}
				close(client.Send)
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) BroadcastToRoom(roomCode string, event string, payload interface{}) {
	data, err := json.Marshal(models.WSMessage{
		Event:   event,
		Payload: payload,
	})
	if err != nil {
		log.Printf("[Hub] Error marshaling broadcast message: %v", err)
		return
	}

	h.mu.RLock()
	clients := h.rooms[roomCode]
	for client := range clients {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
	h.mu.RUnlock()
}

func (h *Hub) SendToClient(client *Client, event string, payload interface{}) {
	data, err := json.Marshal(models.WSMessage{
		Event:   event,
		Payload: payload,
	})
	if err != nil {
		log.Printf("[Hub] Error marshaling client message: %v", err)
		return
	}

	select {
	case client.Send <- data:
	default:
		h.mu.Lock()
		close(client.Send)
		delete(h.clients, client)
		h.mu.Unlock()
	}
}

func (h *Hub) GetOrCreateSession(code string) (*models.Session, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if sess, ok := h.sessions[code]; ok {
		return sess, nil
	}

	// Try loading from Database
	sess, err := db.GetSessionByCode(code)
	if err == nil && sess != nil {
		if sess.Participants == nil {
			sess.Participants = make(map[string]*models.Participant)
		}
		h.sessions[code] = sess
		return sess, nil
	}

	// Fallback create with default quiz 1
	quizzes, err := db.GetAllQuizzes()
	var quiz models.Quiz
	if err == nil && len(quizzes) > 0 {
		quiz = quizzes[0]
	} else {
		quiz = models.Quiz{
			ID:    "quiz-1",
			Title: "Business English Basics",
			Questions: []models.Question{
				{Q: "What does “concise” mean?", A: []string{"Brief and clear", "Difficult to understand", "Very detailed", "Informal"}, C: 0},
			},
		}
	}

	sess = &models.Session{
		Code:                 code,
		QuizID:               quiz.ID,
		Title:                quiz.Title,
		Questions:            quiz.Questions,
		TimePerQuestion:      quiz.TimePerQuestion,
		CurrentQuestionIndex: 0,
		Status:               "LOBBY",
		Participants:         make(map[string]*models.Participant),
	}
	if sess.TimePerQuestion == 0 {
		sess.TimePerQuestion = 20
	}

	h.sessions[code] = sess
	_ = db.SaveSession(sess)
	return sess, nil
}

func getDistribution(session *models.Session, qIndex int) ([]int, int) {
	dist := []int{0, 0, 0, 0}
	answeredCount := 0

	for _, p := range session.Participants {
		if p.Answers != nil {
			if ans, exists := p.Answers[qIndex]; exists {
				if ans.SelectedOption >= 0 && ans.SelectedOption < 4 {
					dist[ans.SelectedOption]++
				}
				answeredCount++
			}
		}
	}
	return dist, answeredCount
}

func (h *Hub) HandleEvent(client *Client, event string, rawData json.RawMessage) {
	switch event {
	case "create_session":
		var req struct {
			QuizID     string `json:"quizId"`
			CustomCode string `json:"customCode"`
		}
		_ = json.Unmarshal(rawData, &req)

		quiz, err := db.GetQuizByID(req.QuizID)
		if err != nil || quiz == nil {
			quizzes, _ := db.GetAllQuizzes()
			if len(quizzes) > 0 {
				quiz = &quizzes[0]
			} else {
				quiz = &models.Quiz{
					ID:    "quiz-1",
					Title: "Business English Basics",
				}
			}
		}

		code := req.CustomCode
		if code == "" {
			rand.Seed(time.Now().UnixNano())
			code = fmt.Sprintf("%04d", rand.Intn(9000)+1000)
		}

		session := &models.Session{
			Code:                 code,
			QuizID:               quiz.ID,
			Title:                quiz.Title,
			Questions:            quiz.Questions,
			TimePerQuestion:      quiz.TimePerQuestion,
			CurrentQuestionIndex: 0,
			Status:               "LOBBY",
			Participants:         make(map[string]*models.Participant),
		}

		h.mu.Lock()
		h.sessions[code] = session
		if h.rooms[code] == nil {
			h.rooms[code] = make(map[*Client]bool)
		}
		h.rooms[code][client] = true
		client.RoomCode = code
		h.mu.Unlock()

		_ = db.SaveSession(session)
		h.SendToClient(client, "session_created", session)

	case "join_room":
		var req struct {
			Code   string `json:"code"`
			Name   string `json:"name"`
			IsHost bool   `json:"isHost"`
		}
		_ = json.Unmarshal(rawData, &req)

		roomCode := req.Code
		if roomCode == "" {
			roomCode = "4821"
		}

		session, err := h.GetOrCreateSession(roomCode)
		if err != nil {
			log.Printf("[Hub] Error getting session %s: %v", roomCode, err)
			return
		}

		h.mu.Lock()
		if h.rooms[roomCode] == nil {
			h.rooms[roomCode] = make(map[*Client]bool)
		}
		h.rooms[roomCode][client] = true
		client.RoomCode = roomCode

		var participantID string
		if !req.IsHost && req.Name != "" {
			existingID := ""
			for id, p := range session.Participants {
				if strings.EqualFold(p.Name, req.Name) {
					existingID = id
					break
				}
			}

			if existingID != "" {
				participantID = existingID
			} else {
				participantID = fmt.Sprintf("p_%d", time.Now().UnixNano()%1000000)
				session.Participants[participantID] = &models.Participant{
					ID:           participantID,
					Name:         req.Name,
					Score:        0,
					CorrectCount: 0,
					Answers:      make(map[int]*models.ParticipantAnswer),
				}
			}
		}
		h.mu.Unlock()

		dist, answeredCount := getDistribution(session, session.CurrentQuestionIndex)

		// Send initial state to client
		h.SendToClient(client, "room_state", models.RoomStatePayload{
			Session:       session,
			ParticipantID: participantID,
			Distribution:  dist,
			AnsweredCount: answeredCount,
		})

		// Broadcast participants update to room
		h.BroadcastToRoom(roomCode, "participants_updated", models.ParticipantsUpdatedPayload{
			Participants: session.Participants,
			Count:        len(session.Participants),
		})

		_ = db.SaveSession(session)

	case "start_session":
		var req struct {
			Code string `json:"code"`
		}
		_ = json.Unmarshal(rawData, &req)

		h.mu.Lock()
		session, ok := h.sessions[req.Code]
		if ok {
			session.Status = "LIVE"
			session.CurrentQuestionIndex = 0
		}
		h.mu.Unlock()

		if session != nil {
			_ = db.SaveSession(session)
			h.BroadcastToRoom(req.Code, "session_started", map[string]interface{}{
				"status":               "LIVE",
				"currentQuestionIndex": 0,
			})
		}

	case "next_question":
		var req struct {
			Code  string `json:"code"`
			Index int    `json:"index"`
		}
		_ = json.Unmarshal(rawData, &req)

		h.mu.Lock()
		session, ok := h.sessions[req.Code]
		if ok && req.Index < len(session.Questions) {
			session.CurrentQuestionIndex = req.Index
		}
		h.mu.Unlock()

		if session != nil {
			_ = db.SaveSession(session)
			dist, answeredCount := getDistribution(session, session.CurrentQuestionIndex)
			h.BroadcastToRoom(req.Code, "question_changed", models.QuestionChangedPayload{
				CurrentQuestionIndex: session.CurrentQuestionIndex,
				Distribution:         dist,
				AnsweredCount:        answeredCount,
			})
		}

	case "submit_answer":
		var req struct {
			Code           string `json:"code"`
			Name           string `json:"name"`
			QuestionIndex  int    `json:"questionIndex"`
			SelectedOption int    `json:"selectedOption"`
		}
		_ = json.Unmarshal(rawData, &req)

		h.mu.Lock()
		session, ok := h.sessions[req.Code]
		var participant *models.Participant
		var isCorrect bool
		var q models.Question

		if ok && req.QuestionIndex < len(session.Questions) {
			q = session.Questions[req.QuestionIndex]
			isCorrect = req.SelectedOption == q.C

			for _, p := range session.Participants {
				if strings.EqualFold(p.Name, req.Name) {
					participant = p
					break
				}
			}

			if participant != nil {
				if participant.Answers == nil {
					participant.Answers = make(map[int]*models.ParticipantAnswer)
				}
				if _, exists := participant.Answers[req.QuestionIndex]; !exists {
					points := 0
					if isCorrect {
						points = 1000
						participant.Score += points
						participant.CorrectCount++
					}
					participant.Answers[req.QuestionIndex] = &models.ParticipantAnswer{
						SelectedOption: req.SelectedOption,
						IsCorrect:      isCorrect,
						Points:         points,
					}
				}
			}
		}
		h.mu.Unlock()

		if session != nil && participant != nil {
			_ = db.SaveSession(session)
			dist, answeredCount := getDistribution(session, req.QuestionIndex)

			// Send result to answer submitter
			h.SendToClient(client, "answer_result", models.AnswerResultPayload{
				QuestionIndex:  req.QuestionIndex,
				IsCorrect:      isCorrect,
				CorrectOption:  q.C,
				MyScore:        participant.Score,
				MyCorrectCount: participant.CorrectCount,
			})

			// Broadcast live update to room
			h.BroadcastToRoom(req.Code, "live_update", models.LiveUpdatePayload{
				Participants: session.Participants,
				Distribution: dist,
				AnsweredCount: answeredCount,
			})
		}

	case "end_session":
		var req struct {
			Code string `json:"code"`
		}
		_ = json.Unmarshal(rawData, &req)

		h.mu.Lock()
		session, ok := h.sessions[req.Code]
		if ok {
			session.Status = "COMPLETED"
		}
		h.mu.Unlock()

		if session != nil {
			_ = db.SaveSession(session)
			h.BroadcastToRoom(req.Code, "session_ended", map[string]interface{}{
				"status":       "COMPLETED",
				"participants": session.Participants,
			})
		}
	}
}
