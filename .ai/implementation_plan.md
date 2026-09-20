# Setup Backend (Go, WebSocket, Postgres, Podman) & Vue Integration Plan

This document outlines the implementation plan for creating a robust Go backend server with PostgreSQL storage and WebSocket real-time communication, containerized using Podman & Compose with Makefile targets (`make up`, `make down`), and integrating it with the existing Vue 3 frontend application.

## User Review Required

> [!IMPORTANT]
> The backend will run on port `8080` inside Podman containers (Postgres on port `5432`). `podman compose` (backed by Docker Compose / Podman on Windows) will orchestrate the environment.
> The Vue dev server (`vite`) will proxy `/api` and `/ws` to `http://localhost:8080`.

## Open Questions

None. The existing Node/Socket.io mock server (`server.js`) contract was analyzed and will be fully translated into native Go WebSockets + REST API handlers with PostgreSQL persistence.

---

## Proposed Changes

### Go Backend (`backend/`)

#### [NEW] [go.mod](file:///d:/playground/elsa-real-time-quiz/backend/go.mod)
- Define Go module `elsa-real-time-quiz/backend`.
- Include dependencies: `github.com/gorilla/websocket`, `github.com/lib/pq`.

#### [NEW] [models.go](file:///d:/playground/elsa-real-time-quiz/backend/pkg/models/models.go)
- Struct definitions for `Quiz`, `Question`, `Session`, `Participant`, `Answer`, and WebSocket event message payloads.

#### [NEW] [db.go](file:///d:/playground/elsa-real-time-quiz/backend/pkg/db/db.go)
- PostgreSQL connection pool setup using `database/sql` + `lib/pq`.
- Automatic table creation migrations (`quizzes`, `sessions`, `participants`, `answers`).
- Automatic seeding of default quizzes ("Business English Basics", "Travel Vocabulary") and initial active demo room (`4821`).

#### [NEW] [hub.go](file:///d:/playground/elsa-real-time-quiz/backend/pkg/ws/hub.go) & [client.go](file:///d:/playground/elsa-real-time-quiz/backend/pkg/ws/client.go)
- In-memory WebSocket hub & room manager.
- Handles socket upgrades, message parsing, room broadcasts, and DB state sync for:
  - `join_room`
  - `create_session`
  - `start_session`
  - `next_question`
  - `submit_answer`
  - `end_session`

#### [NEW] [handlers.go](file:///d:/playground/elsa-real-time-quiz/backend/pkg/api/handlers.go)
- REST HTTP handlers:
  - `GET /health`
  - `GET /api/quizzes` & `POST /api/quizzes` & `GET /api/quizzes/:id`
  - `GET /api/sessions`
  - `GET /ws` (WebSocket connection handler)

#### [NEW] [main.go](file:///d:/playground/elsa-real-time-quiz/backend/cmd/server/main.go)
- Application entrypoint initializing DB connection, WS Hub, HTTP router, CORS headers, and listening on `PORT` (default 8080).

---

### Containerization & Service Management

#### [NEW] [Containerfile](file:///d:/playground/elsa-real-time-quiz/backend/Containerfile)
- Multi-stage build for Go application:
  - Build stage: `golang:1.22-alpine` compiling binary.
  - Final image: `alpine:latest` with SSL certs and binary runner.

#### [NEW] [compose.yaml](file:///d:/playground/elsa-real-time-quiz/compose.yaml)
- Define `db` (PostgreSQL 16 container with healthcheck) and `backend` (Go backend container dependent on `db`).

#### [NEW] [Makefile](file:///d:/playground/elsa-real-time-quiz/Makefile)
- Targets:
  - `up`: Runs `podman compose up -d --build` (with fallback logic for `podman-compose` / `docker compose`).
  - `down`: Runs `podman compose down`.
  - `logs`: Tail container logs.
  - `status`: Show container status.
  - `restart`: Restart backend containers.

---

### Vue Frontend Integration

#### [MODIFY] [vite.config.js](file:///d:/playground/elsa-real-time-quiz/vite.config.js)
- Update server proxy target from port 3001 to `http://localhost:8080` for `/api` and `/ws`.

#### [MODIFY] [socket.js](file:///d:/playground/elsa-real-time-quiz/src/services/socket.js)
- Implement WebSocket client wrapper providing standard event emitter syntax (`socket.on`, `socket.off`, `socket.emit`) over native browser WebSocket `ws://.../ws`.

#### [MODIFY] [quiz.js](file:///d:/playground/elsa-real-time-quiz/src/stores/quiz.js)
- Connect Pinia quiz store to fetch/save quizzes via Go REST API `/api/quizzes`.

#### [MODIFY] [liveSession.js](file:///d:/playground/elsa-real-time-quiz/src/stores/liveSession.js)
- Connect live session store to Go WebSocket service for host & participant workflows.

#### [MODIFY] [package.json](file:///d:/playground/elsa-real-time-quiz/package.json)
- Update dev scripts to reflect containerized backend workflow.

---

## Verification Plan

### Automated Tests & Container Verification
- Test container startup: `make up`
- Check container logs: `make logs`
- Verify database connection & seeding via `/health` and `/api/quizzes` endpoints:
  - `curl http://localhost:8080/health`
  - `curl http://localhost:8080/api/quizzes`
  - `curl http://localhost:8080/api/sessions`

### Manual Verification
1. Start containers with `make up`.
2. Launch Vue frontend (`npm run dev:client`).
3. Open Host dashboard (`http://localhost:5173/host/dashboard`), create/select quiz, start a room session.
4. Open Participant view (`http://localhost:5173/join/4821`), submit answers, verify real-time leaderboard update, answer distributions, and question transitions over WebSockets.
5. Tear down containers with `make down`.
