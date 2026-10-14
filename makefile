
API_PORT ?= 8081

.PHONY: help up down build restart logs clean backend frontend install \
        build-backend build-frontend lint

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'

# ---------------------------------------------------------------------------
# Docker
# ---------------------------------------------------------------------------
up: ## Build and start backend + frontend (docker compose)
	docker compose up -d --build

down: ## Stop the compose stack
	docker compose down

build: ## Build both images without starting them
	docker compose build

restart: ## Restart the compose stack
	docker compose restart

logs: ## Follow the logs of both containers
	docker compose logs -f

clean: ## Stop the stack and remove the images it built
	docker compose down --rmi local --remove-orphans

# ---------------------------------------------------------------------------
# Native development
# ---------------------------------------------------------------------------
backend: ## Run the Go API natively (listens on config.json port)
	cd client_app && go run .

frontend: ## Run the Next.js dev server natively
	cd client_web && npm run dev

install: ## Install the frontend dependencies
	cd client_web && npm ci

build-backend: ## Compile the Go API into client_app/main
	cd client_app && CGO_ENABLED=1 go build -o main .

build-frontend: ## Production build of the Next.js app
	cd client_web && npm run build

lint: ## Lint the frontend sources
	cd client_web && npm run lint
