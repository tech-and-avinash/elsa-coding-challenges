# ELSA Real-Time Quiz - Go Backend & Vue Integration Walkthrough

We have successfully set up a full Go backend server with PostgreSQL storage and real-time WebSockets, containerized using Podman Compose, managed via Makefile targets (`make up`, `make down`), and integrated with the Vue 3 frontend application.

---

## 🛠️ Architecture & Components Created

### 1. Go Backend Server (`backend/`)
- **Main Server Entrypoint**: [`cmd/server/main.go`](file:///d:/playground/elsa-real-time-quiz/backend/cmd/server/main.go)
  - Connects to PostgreSQL, starts the in-memory WebSocket Hub, configures CORS middleware, and listens on port `8080`.
- **Database & Auto-Migrations**: [`pkg/db/db.go`](file:///d:/playground/elsa-real-time-quiz/backend/pkg/db/db.go)
  - PostgreSQL pool setup (`database/sql` + `lib/pq`).
  - Auto-migration of `quizzes` and `sessions` database tables.
  - Auto-seeding default quizzes ("Business English Basics", "Travel Vocabulary") and active demo room `4821`.
- **WebSocket Hub & Handlers**: [`pkg/ws/hub.go`](file:///d:/playground/elsa-real-time-quiz/backend/pkg/ws/hub.go) & [`pkg/ws/client.go`](file:///d:/playground/elsa-real-time-quiz/backend/pkg/ws/client.go)
  - Real-time event engine handling room creation, room joining, session state synchronization, question progression, answer submission, score calculation, distribution tracking, and PostgreSQL persistence.
- **REST API Endpoints**: [`pkg/api/handlers.go`](file:///d:/playground/elsa-real-time-quiz/backend/pkg/api/handlers.go)
  - `/health`: Healthcheck endpoint for container status.
  - `/api/quizzes`: List and create/update quizzes in PostgreSQL.
  - `/api/sessions`: List active and past sessions.

### 2. Containerization & Management (`Makefile`, `compose.yaml`, `Containerfile`)
- **Containerfile**: [`backend/Containerfile`](file:///d:/playground/elsa-real-time-quiz/backend/Containerfile)
  - Multi-stage build (`golang:1.22-alpine` builder stage -> minimal `alpine:3.19` runner).
- **Compose Specification**: [`compose.yaml`](file:///d:/playground/elsa-real-time-quiz/compose.yaml)
  - Orchestrates PostgreSQL 16 container (`db`) with persistent volume `pgdata` and healthcheck, and the Go server container (`backend`).
- **Makefile Commands**: [`Makefile`](file:///d:/playground/elsa-real-time-quiz/Makefile)
  - `make up`: Build and start containers in detached mode (`podman compose up -d --build`).
  - `make down`: Stop and remove containers (`podman compose down`).
  - `make logs`: Stream container logs (`podman compose logs -f`).
  - `make status`: Check running container status (`podman compose ps`).

### 3. Vue 3 Integration
- **Vite Proxy**: [`vite.config.js`](file:///d:/playground/elsa-real-time-quiz/vite.config.js)
  - Proxies `/api` requests to `http://localhost:8080` and `/ws` WebSocket connection to `ws://localhost:8080`.
- **Native WebSocket Client**: [`src/services/socket.js`](file:///d:/playground/elsa-real-time-quiz/src/services/socket.js)
  - Socket event emitter implementation over native browser WebSockets.
- **Pinia Stores**: [`src/stores/quiz.js`](file:///d:/playground/elsa-real-time-quiz/src/stores/quiz.js) & [`src/stores/liveSession.js`](file:///d:/playground/elsa-real-time-quiz/src/stores/liveSession.js)
  - Fetches and saves quiz data from Go REST API and manages real-time quiz room state over WebSockets.

---

## ⚡ Quick Start / How to Run

### 1. Backend Containers
```bash
# Start backend and PostgreSQL database containers
make up

# View container logs
make logs

# Check container status
make status

# Stop containers
make down
```

### 2. Frontend Development Server
```bash
npm run dev
```
Open [http://localhost:5173](http://localhost:5173) in your browser.

---

## 🔍 Verification Results

1. **Containers**:
   - Both `elsa_quiz_db` (PostgreSQL) and `elsa_quiz_backend` (Go backend) built cleanly and passed health checks.
2. **REST API**:
   - `http://localhost:8080/health` returned `{"db":"connected","status":"ok"}`.
   - `http://localhost:8080/api/quizzes` returned default quizzes from PostgreSQL.
3. **WebSockets & UI**:
   - Joined room `4821` at `http://localhost:5173/join/4821` over WebSockets cleanly (`[WebSocket Service] Connected successfully`).
   - Host dashboard at `http://localhost:5173/host/dashboard` populated sessions and quizzes dynamically.
