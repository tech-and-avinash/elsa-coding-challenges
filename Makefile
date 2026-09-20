COMPOSE ?= podman compose

.PHONY: help up down logs status restart db-shell build

help:
	@echo "Available commands:"
	@echo "  make up       - Build and start the frontend, Go backend, and PostgreSQL containers in background"
	@echo "  make down     - Stop and remove containers"
	@echo "  make logs     - View and follow container logs"
	@echo "  make status   - Show running container status"
	@echo "  make restart  - Restart all backend containers"
	@echo "  make db-shell - Access PostgreSQL shell inside db container"

up:
	@echo "Starting the ELSA Quiz stack with Podman..."
	$(COMPOSE) up -d --build

down:
	@echo "Stopping ELSA Quiz backend services..."
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f

status:
	$(COMPOSE) ps

restart:
	$(COMPOSE) restart

db-shell:
	podman exec -it elsa_quiz_db psql -U elsa -d elsa_quiz
