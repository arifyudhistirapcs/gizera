#!/bin/bash
# ============================================
# ERP SPPG - Production Deployment Script
# Usage: ./deployment/scripts/deploy.sh [tag]
# ============================================
set -e

TAG="${1:-latest}"
DEPLOY_DIR="$(cd "$(dirname "$0")/.." && pwd)"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

log() { echo -e "[$(date '+%H:%M:%S')] $1"; }

log "=== ERP SPPG Production Deploy (tag: $TAG) ==="

# Pre-flight checks
log "Running pre-flight checks..."

if [ ! -f "$DEPLOY_DIR/.env" ]; then
    echo -e "${RED}ERROR: deployment/.env not found. Copy from .env.example${NC}"
    exit 1
fi

if [ ! -f "$DEPLOY_DIR/firebase-credentials.json" ]; then
    echo -e "${RED}ERROR: deployment/firebase-credentials.json not found${NC}"
    exit 1
fi

if [ ! -f "$DEPLOY_DIR/nginx/ssl/cert.pem" ] || [ ! -f "$DEPLOY_DIR/nginx/ssl/key.pem" ]; then
    echo -e "${RED}ERROR: SSL certificates not found in deployment/nginx/ssl/${NC}"
    exit 1
fi

log "${GREEN}✓${NC} Pre-flight checks passed"

# Create required directories
mkdir -p "$DEPLOY_DIR/logs/backend-1" "$DEPLOY_DIR/logs/backend-2" \
         "$DEPLOY_DIR/uploads" "$DEPLOY_DIR/backups" "$DEPLOY_DIR/nginx/logs"

# Build images
log "Building Docker images..."
cd "$(dirname "$DEPLOY_DIR")"
./deployment/scripts/build.sh "$TAG"

# Deploy
log "Starting production services..."
cd "$DEPLOY_DIR"
docker compose -f docker-compose.prod.yml up -d

# Wait for health
log "Waiting for services to be healthy..."
sleep 10

# Check backend health
for i in 1 2; do
    if docker exec "erp-sppg-backend-$i" curl -sf http://localhost:8080/health > /dev/null 2>&1; then
        log "${GREEN}✓${NC} backend-$i healthy"
    else
        log "${RED}✗${NC} backend-$i not healthy yet (may still be starting)"
    fi
done

log ""
log "=== Deployment complete ==="
docker compose -f docker-compose.prod.yml ps
