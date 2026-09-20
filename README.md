# ELSA Real-Time Vocabulary Quiz

A real-time vocabulary quiz application with separate host and participant experiences. The frontend is built with Vue, the live quiz backend is written in Go, and PostgreSQL stores persistent application data.

## Stack

- Vue 3, Vite, Pinia, and Vue Router
- Go 1.22 with WebSocket support
- PostgreSQL 16
- Podman and Podman Compose
- Make for common development commands

## Run With Podman

### Prerequisites

- Podman
- Podman Compose, or a compatible Compose provider
- GNU Make

Start the complete stack from the repository root:

```bash
make up
```

The command builds and starts the frontend, Go backend, and PostgreSQL containers in the background.

Open the application at:

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432

The demo session code is `4821`.

## Make Commands

```bash
make up        # Build and start all services
make down      # Stop and remove containers
make logs      # Follow service logs
make status    # Show container status
make restart   # Restart all services
make db-shell  # Open a PostgreSQL shell
```

To use a different Compose command, override `COMPOSE`:

```bash
make COMPOSE="podman-compose" up
```

## Services

| Service | Purpose | Host port |
| --- | --- | ---: |
| `frontend` | Vue/Vite development server | `5173` |
| `backend` | Go HTTP and WebSocket API | `8080` |
| `db` | PostgreSQL persistence | `5432` |

The frontend proxies `/api` and `/ws` requests to the backend. Inside the Compose network, it uses `http://backend:8080`; in local Vite development, it defaults to `http://localhost:8080`.

## Local Frontend Development

Install dependencies and run Vite directly when the frontend container is not needed:

```bash
npm install
npm run dev
```

For a production build:

```bash
npm run build
```

## Project Structure

```text
src/                  Vue application, routes, stores, and WebSocket client
backend/              Go API, WebSocket hub, database, and models
frontend/Containerfile Frontend development container image
compose.yaml          Podman Compose service definitions
Makefile              Container lifecycle commands
```

## Documentation

See [SDD.md](SDD.md) for the system design and real-time data-flow details.