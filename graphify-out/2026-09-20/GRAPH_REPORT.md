# Graph Report - elsa-real-time-quiz  (2026-09-20)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 200 nodes · 310 edges · 8 communities (7 shown, 1 thin omitted)
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- db.go
- package.json
- models.go
- index.js
- HostQuizEditorView.vue
- HostSessionView.vue
- ParticipantApp.vue
- elsa-real-time-quiz/backend

## God Nodes (most connected - your core abstractions)
1. `Hub` - 11 edges
2. `vue` - 11 edges
3. `Session` - 10 edges
4. `WsClient` - 9 edges
5. `vue-router` - 9 edges
6. `ServeWS()` - 7 edges
7. `Client` - 7 edges
8. `QuizzesHandler()` - 6 edges
9. `scripts` - 6 edges
10. `main()` - 5 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `ServeWS()`  [EXTRACTED]
  backend/cmd/server/main.go → backend/pkg/ws/client.go
- `main()` --calls--> `NewHub()`  [EXTRACTED]
  backend/cmd/server/main.go → backend/pkg/ws/hub.go
- `QuizzesHandler()` --calls--> `GetAllQuizzes()`  [EXTRACTED]
  backend/pkg/api/handlers.go → backend/pkg/db/db.go
- `GetQuizByID()` --references--> `Quiz`  [EXTRACTED]
  backend/pkg/db/db.go → backend/pkg/models/models.go
- `SaveQuiz()` --references--> `Quiz`  [EXTRACTED]
  backend/pkg/db/db.go → backend/pkg/models/models.go

## Import Cycles
- None detected.

## Communities (8 total, 1 thin omitted)

### Community 0 - "db.go"
Cohesion: 0.09
Nodes (32): corsMiddleware(), main(), HealthHandler(), QuizzesHandler(), SessionsHandler(), GetAllSessions(), GetQuizByID(), InitDB() (+24 more)

### Community 1 - "package.json"
Cohesion: 0.06
Nodes (33): dependencies, express, pinia, socket.io, socket.io-client, vue, vue-router, devDependencies (+25 more)

### Community 2 - "models.go"
Cohesion: 0.10
Nodes (24): GetAllQuizzes(), GetSessionByCode(), SaveSession(), Quiz, Session, ServeWS(), getDistribution(), NewHub() (+16 more)

### Community 3 - "index.js"
Cohesion: 0.09
Nodes (21): pinia, vue, src_assets_main, app, router, routes, useAuthStore, authStore (+13 more)

### Community 4 - "HostQuizEditorView.vue"
Cohesion: 0.10
Nodes (11): vue-router, useQuizStore, quizStore, router, quiz, quizStore, route, router (+3 more)

### Community 5 - "HostSessionView.vue"
Cohesion: 0.12
Nodes (11): getSocket(), WsClient, useLiveSessionStore, averageScore, copied, currentAccuracy, isLastQuestion, joinUrl (+3 more)

### Community 6 - "ParticipantApp.vue"
Cohesion: 0.14
Nodes (9): currentStage, displayName, inputCode, inputName, isSelectedCorrect, progressPercent, route, sessionStore (+1 more)

## Knowledge Gaps
- **69 isolated node(s):** `express`, `pinia`, `socket.io`, `socket.io-client`, `vue` (+64 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 108 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **1 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `vue` connect `index.js` to `package.json`, `HostQuizEditorView.vue`, `HostSessionView.vue`, `ParticipantApp.vue`?**
  _High betweenness centrality (0.126) - this node is a cross-community bridge._
- **Why does `vue-router` connect `HostQuizEditorView.vue` to `package.json`, `index.js`, `HostSessionView.vue`, `ParticipantApp.vue`?**
  _High betweenness centrality (0.103) - this node is a cross-community bridge._
- **What connects `express`, `pinia`, `socket.io` to the rest of the system?**
  _69 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `db.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08858858858858859 - nodes in this community are weakly interconnected._
- **Should `package.json` be split into smaller, more focused modules?**
  _Cohesion score 0.057057057057057055 - nodes in this community are weakly interconnected._
- **Should `models.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10483870967741936 - nodes in this community are weakly interconnected._
- **Should `index.js` be split into smaller, more focused modules?**
  _Cohesion score 0.08669354838709678 - nodes in this community are weakly interconnected._