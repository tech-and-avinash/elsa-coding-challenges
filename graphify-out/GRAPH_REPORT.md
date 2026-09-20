# Graph Report - elsa-real-time-quiz  (2026-09-21)

## Corpus Check
- 35 files · ~15,614 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 5 file(s) not represented in the graph (top: (none) 4, .css 1)

## Summary
- 291 nodes · 402 edges · 15 communities (9 shown, 6 thin omitted)
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
- System Design Document — Real-Time Vocabulary Quiz
- rules/graphify.md
- workflows/graphify.md
- CLAUDE.md
- copilot-instructions.md
- AI Collaboration

## God Nodes (most connected - your core abstractions)
1. `System Design Document — Real-Time Vocabulary Quiz` - 15 edges
2. `vue` - 13 edges
3. `Hub` - 11 edges
4. `Session` - 10 edges
5. `vue-router` - 9 edges
6. `WsClient` - 9 edges
7. `AI Collaboration` - 9 edges
8. `QuizzesHandler()` - 7 edges
9. `Client` - 7 edges
10. `ServeWS()` - 7 edges

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

## Communities (15 total, 6 thin omitted)

### Community 0 - "hub.go"
Cohesion: 0.08
Nodes (30): corsMiddleware(), main(), HealthHandler(), SessionsHandler(), ServeWS(), getDistribution(), NewHub(), go_pkg_elsa_real_time_quiz_backend_pkg_api (+22 more)

### Community 1 - "package.json"
Cohesion: 0.06
Nodes (33): dependencies, express, pinia, socket.io, socket.io-client, vue, vue-router, devDependencies (+25 more)

### Community 2 - "db.go"
Cohesion: 0.12
Nodes (27): QuizzesHandler(), GetAllQuizzes(), GetAllSessions(), GetQuizByID(), GetSessionByCode(), GetSessionsByQuizID(), InitDB(), migrateSchema() (+19 more)

### Community 3 - "index.js"
Cohesion: 0.05
Nodes (34): pinia, vue, vue-router, src_assets_main, app, router, routes, useAuthStore (+26 more)

### Community 4 - "Go Backend (`backend/`)"
Cohesion: 0.08
Nodes (23): Automated Tests & Container Verification, Containerization & Service Management, Go Backend (`backend/`), Manual Verification, [MODIFY] [liveSession.js](file:///d:/playground/elsa-real-time-quiz/src/stores/liveSession.js), [MODIFY] [package.json](file:///d:/playground/elsa-real-time-quiz/package.json), [MODIFY] [quiz.js](file:///d:/playground/elsa-real-time-quiz/src/stores/quiz.js), [MODIFY] [vite.config.js](file:///d:/playground/elsa-real-time-quiz/vite.config.js) (+15 more)

### Community 5 - "ParticipantApp.vue"
Cohesion: 0.08
Nodes (19): getSocket(), useLiveSessionStore, averageScore, copied, currentAccuracy, isLastQuestion, joinUrl, responsePercentage (+11 more)

### Community 6 - "🛠️ Architecture & Components Created"
Cohesion: 0.20
Nodes (9): 1. Backend Containers, 1. Go Backend Server (`backend/`), 2. Containerization & Management (`Makefile`, `compose.yaml`, `Containerfile`), 2. Frontend Development Server, 3. Vue 3 Integration, 🛠️ Architecture & Components Created, ELSA Real-Time Quiz - Go Backend & Vue Integration Walkthrough, ⚡ Quick Start / How to Run (+1 more)

### Community 9 - "System Design Document — Real-Time Vocabulary Quiz"
Cohesion: 0.08
Nodes (25): 10. Maintainability, 11. AI Collaboration, 12. Testing, 13. Scope and Trade-offs, 14. Summary, 1. HTML prototype, 1. Overview, 2. Approach (+17 more)

### Community 14 - "AI Collaboration"
Cohesion: 0.11
Nodes (17): 1. Product Exploration and Prototype — ChatGPT, 2. HTML to Vue.js — Antigravity, 3. Frontend Refinement — Antigravity, OpenCode and Copilot, 4. Backend Setup and Real-Time Development — Antigravity, OpenCode, Copilot and Devin, 5. AI-Assisted Real-Time Logic, 6. Verification of AI-Assisted Code, 7. Why Multiple AI Tools?, 8. What I Learned (+9 more)

## Knowledge Gaps
- **131 isolated node(s):** `elsa-real-time-quiz/backend`, `WSMessage`, `QuestionChangedPayload`, `AnswerResultPayload`, `name` (+126 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 176 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `vue` connect `index.js` to `package.json`, `ParticipantApp.vue`?**
  _High betweenness centrality (0.103) - this node is a cross-community bridge._
- **Why does `WsClient` connect `WsClient` to `ParticipantApp.vue`?**
  _High betweenness centrality (0.097) - this node is a cross-community bridge._
- **Why does `[MODIFY] [socket.js](file:///d:/playground/elsa-real-time-quiz/src/services/socket.js)` connect `WsClient` to `Go Backend (`backend/`)`?**
  _High betweenness centrality (0.076) - this node is a cross-community bridge._
- **What connects `elsa-real-time-quiz/backend`, `WSMessage`, `QuestionChangedPayload` to the rest of the system?**
  _131 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `hub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08205128205128205 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.057057057057057055 - nodes in this community are weakly interconnected._
- **Should `db.go` be split into smaller, more focused modules?**
  _Cohesion score 0.11954022988505747 - nodes in this community are weakly interconnected._