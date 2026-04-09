# ============================================
# ERP SPPG - Makefile
# ============================================

.PHONY: help dev dev-build dev-down dev-logs dev-ps \
        build push deploy prod-up prod-down prod-logs \
        db-only backend-only web-only pwa-only \
        db-migrate db-backup db-restore \
        clean test lint

# Default
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ============================================
# Development (docker-compose.yml)
# ============================================

dev: ## Start all services (dev)
	docker compose up -d

dev-build: ## Build & start all services (dev)
	docker compose up -d --build

dev-down: ## Stop all services (dev)
	docker compose down

dev-reset: ## Stop all + remove volumes (reset DB)
	docker compose down -v

dev-logs: ## Tail logs all services
	docker compose logs -f

dev-ps: ## Show running containers
	docker compose ps

# Individual services
db-only: ## Start only PostgreSQL + Redis
	docker compose up -d postgres redis

backend-only: ## Rebuild & restart backend only
	docker compose up -d --build backend

web-only: ## Rebuild & restart web only
	docker compose up -d --build web

pwa-only: ## Rebuild & restart pwa only
	docker compose up -d --build pwa

# ============================================
# Build Docker Images
# ============================================

TAG ?= latest

build: ## Build all Docker images (TAG=latest)
	chmod +x deployment/scripts/build.sh
	./deployment/scripts/build.sh $(TAG)

push: ## Build & push to registry (DOCKER_REGISTRY=xxx TAG=latest)
	DOCKER_REGISTRY=$(DOCKER_REGISTRY) ./deployment/scripts/build.sh $(TAG)

# ============================================
# Production (deployment/docker-compose.prod.yml)
# ============================================

prod-up: ## Start production stack
	docker compose -f deployment/docker-compose.prod.yml up -d

prod-down: ## Stop production stack
	docker compose -f deployment/docker-compose.prod.yml down

prod-logs: ## Tail production logs
	docker compose -f deployment/docker-compose.prod.yml logs -f

prod-ps: ## Show production containers
	docker compose -f deployment/docker-compose.prod.yml ps

prod-restart: ## Restart production backend
	docker compose -f deployment/docker-compose.prod.yml restart backend-1 backend-2

# ============================================
# Database
# ============================================

db-migrate: ## Run database migrations
	docker compose exec backend ./server migrate

db-backup: ## Manual database backup
	docker compose exec postgres pg_dump -U $${DB_USER:-erp_sppg_user} -d $${DB_NAME:-erp_sppg} --format=custom > backups/backup_$$(date +%Y%m%d_%H%M%S).sql

db-restore: ## Restore database (FILE=path/to/backup.sql)
	docker compose exec -T postgres pg_restore -U $${DB_USER:-erp_sppg_user} -d $${DB_NAME:-erp_sppg} --clean --if-exists < $(FILE)

db-shell: ## Open psql shell
	docker compose exec postgres psql -U $${DB_USER:-erp_sppg_user} -d $${DB_NAME:-erp_sppg}

redis-shell: ## Open redis-cli
	docker compose exec redis redis-cli

# ============================================
# Testing
# ============================================

test-backend: ## Run backend tests
	docker compose exec backend go test ./... -v

test-web: ## Run web tests
	docker compose exec web npm run test

test-pwa: ## Run pwa tests
	docker compose exec pwa npm run test

# ============================================
# Cleanup
# ============================================

clean: ## Remove all containers, images, volumes
	docker compose down -v --rmi local
	docker image prune -f
