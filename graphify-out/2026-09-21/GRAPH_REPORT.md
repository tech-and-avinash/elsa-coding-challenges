# Graph Report - elsa-real-time-quiz  (2026-09-20)

## Corpus Check
- 33 files · ~12,227 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 4 file(s) not represented in the graph (top: (none) 3, .css 1)

## Summary
- 244 nodes · 352 edges · 14 communities (8 shown, 6 thin omitted)
- Extraction: 99% EXTRACTED · 1% INFERRED · 0% AMBIGUOUS · INFERRED: 3 edges (avg confidence: 0.95)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- hub.go
- package.json
- db.go
- index.js
- Go Backend (`backend/`)
- ParticipantApp.vue
- 🛠️ Architecture & Components Created
- elsa-real-time-quiz/backend
- WsClient
- dependencies
- rules/graphify.md
- workflows/graphify.md
- CLAUDE.md
- copilot-instructions.md

## God Nodes (most connected - your core abstractions)
1. `vue` - 12 edges
2. `Hub` - 11 edges
3. `Session` - 10 edges
4. `vue-router` - 9 edges
5. `WsClient` - 9 edges
6. `Client` - 7 edges
7. `ServeWS()` - 7 edges
8. `Go Backend (`backend/`)` - 7 edges
9. `QuizzesHandler()` - 6 edges
10. `scripts` - 6 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `InitDB()`  [EXTRACTED]
  backend/cmd/server/main.go → backend/pkg/db/db.go
- `SessionsHandler()` --calls--> `GetAllSessions()`  [EXTRACTED]
  backend/pkg/api/handlers.go → backend/pkg/db/db.go
- `getDistribution()` --references--> `Session`  [EXTRACTED]
  backend/pkg/ws/hub.go → backend/pkg/models/models.go
- `Hub` --references--> `Session`  [EXTRACTED]
  backend/pkg/ws/hub.go → backend/pkg/models/models.go
- `main()` --calls--> `ServeWS()`  [EXTRACTED]
  backend/cmd/server/main.go → backend/pkg/ws/client.go

## Import Cycles
- None detected.

## Communities (14 total, 6 thin omitted)

### Community 0 - "hub.go"
Cohesion: 0.08
Nodes (30): corsMiddleware(), main(), HealthHandler(), SessionsHandler(), ServeWS(), getDistribution(), NewHub(), go_pkg_elsa_real_time_quiz_backend_pkg_api (+22 more)

### Community 1 - "package.json"
Cohesion: 0.07
Nodes (26): devDependencies, concurrently, vite, @vitejs/plugin-vue, name, private, scripts, build (+18 more)

### Community 2 - "db.go"
Cohesion: 0.12
Nodes (26): QuizzesHandler(), GetAllQuizzes(), GetAllSessions(), GetQuizByID(), GetSessionByCode(), InitDB(), migrateSchema(), SaveQuiz() (+18 more)

### Community 3 - "index.js"
Cohesion: 0.05
Nodes (33): pinia, vue, vue-router, src_assets_main, app, router, routes, useAuthStore (+25 more)

### Community 4 - "Go Backend (`backend/`)"
Cohesion: 0.08
Nodes (23): Automated Tests & Container Verification, Containerization & Service Management, Go Backend (`backend/`), Manual Verification, [MODIFY] [liveSession.js](file:///d:/playground/elsa-real-time-quiz/src/stores/liveSession.js), [MODIFY] [package.json](file:///d:/playground/elsa-real-time-quiz/package.json), [MODIFY] [quiz.js](file:///d:/playground/elsa-real-time-quiz/src/stores/quiz.js), [MODIFY] [vite.config.js](file:///d:/playground/elsa-real-time-quiz/vite.config.js) (+15 more)

### Community 5 - "ParticipantApp.vue"
Cohesion: 0.08
Nodes (19): getSocket(), useLiveSessionStore, averageScore, copied, currentAccuracy, isLastQuestion, joinUrl, responsePercentage (+11 more)

### Community 6 - "🛠️ Architecture & Components Created"
Cohesion: 0.20
Nodes (9): 1. Backend Containers, 1. Go Backend Server (`backend/`), 2. Containerization & Management (`Makefile`, `compose.yaml`, `Containerfile`), 2. Frontend Development Server, 3. Vue 3 Integration, 🛠️ Architecture & Components Created, ELSA Real-Time Quiz - Go Backend & Vue Integration Walkthrough, ⚡ Quick Start / How to Run (+1 more)

### Community 9 - "dependencies"
Cohesion: 0.29
Nodes (7): dependencies, express, pinia, socket.io, socket.io-client, vue, vue-router

## Knowledge Gaps
- **97 isolated node(s):** `elsa-real-time-quiz/backend`, `WSMessage`, `QuestionChangedPayload`, `AnswerResultPayload`, `name` (+92 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 140 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `vue` connect `index.js` to `package.json`, `ParticipantApp.vue`?**
  _High betweenness centrality (0.138) - this node is a cross-community bridge._
- **Why does `WsClient` connect `WsClient` to `ParticipantApp.vue`?**
  _High betweenness centrality (0.136) - this node is a cross-community bridge._
- **Why does `[MODIFY] [socket.js](file:///d:/playground/elsa-real-time-quiz/src/services/socket.js)` connect `WsClient` to `Go Backend (`backend/`)`?**
  _High betweenness centrality (0.107) - this node is a cross-community bridge._
- **What connects `elsa-real-time-quiz/backend`, `WSMessage`, `QuestionChangedPayload` to the rest of the system?**
  _97 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `hub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08205128205128205 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.07126436781609195 - nodes in this community are weakly interconnected._
- **Should `db.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1206896551724138 - nodes in this community are weakly interconnected._