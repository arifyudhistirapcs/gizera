#!/bin/bash
# ============================================
# ERP SPPG - Docker Development Setup
# Otomatis setup semua yang dibutuhkan untuk
# menjalankan sistem via Docker
# ============================================
set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "=========================================="
echo "ERP SPPG - Docker Setup"
echo "=========================================="
echo ""

# Check Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker tidak ditemukan. Install dulu: https://docs.docker.com/get-docker/${NC}"
    exit 1
fi
echo -e "${GREEN}✓${NC} Docker: $(docker --version | awk '{print $3}')"

# Check Docker Compose
if ! docker compose version &> /dev/null; then
    echo -e "${RED}✗ Docker Compose tidak ditemukan.${NC}"
    exit 1
fi
echo -e "${GREEN}✓${NC} Docker Compose: $(docker compose version --short)"

echo ""

# 1. Setup root .env
if [ ! -f .env ]; then
    cp .env.example .env
    echo -e "${GREEN}✓${NC} Created .env (root)"
else
    echo -e "${YELLOW}→${NC} .env (root) sudah ada, skip"
fi

# 2. Setup backend .env
if [ ! -f backend/.env ]; then
    cp backend/.env.example backend/.env
    echo -e "${GREEN}✓${NC} Created backend/.env"
else
    echo -e "${YELLOW}→${NC} backend/.env sudah ada, skip"
fi

# 3. Setup web .env
if [ ! -f web/.env ]; then
    cp web/.env.example web/.env
    echo -e "${GREEN}✓${NC} Created web/.env"
else
    echo -e "${YELLOW}→${NC} web/.env sudah ada, skip"
fi

# 4. Setup pwa .env
if [ ! -f pwa/.env ]; then
    cp pwa/.env.example pwa/.env
    echo -e "${GREEN}✓${NC} Created pwa/.env"
else
    echo -e "${YELLOW}→${NC} pwa/.env sudah ada, skip"
fi

# 5. Firebase credentials
if [ ! -f backend/firebase-credentials.json ]; then
    echo -e "${YELLOW}⚠${NC}  backend/firebase-credentials.json belum ada"
    echo "   → Buka https://console.firebase.google.com/"
    echo "   → Project Settings > Service Accounts > Generate New Private Key"
    echo "   → Simpan sebagai backend/firebase-credentials.json"
    echo "   → Tanpa file ini, backend jalan tapi KDS realtime, dashboard sync,"
    echo "     monitoring, dan notifications TIDAK akan berfungsi"
    echo ""
    echo -e "${YELLOW}⚠${NC}  Jangan lupa isi Firebase config di web/.env dan pwa/.env juga:"
    echo "   → Firebase Console > Project Settings > General > Your Apps > Web App"
    echo "   → Copy apiKey, authDomain, databaseURL, projectId, dll"
else
    echo -e "${GREEN}✓${NC} Firebase credentials found"
fi

# 6. Create required directories
mkdir -p backend/logs backend/uploads backups
echo -e "${GREEN}✓${NC} Directories created (logs, uploads, backups)"

echo ""
echo "=========================================="
echo "Setup selesai! Jalankan dengan:"
echo "=========================================="
echo ""
echo "  docker compose up -d --build"
echo ""
echo "Atau pakai Makefile:"
echo "  make dev-build"
echo ""
echo "Akses:"
echo "  Backend API  : http://localhost:8080"
echo "  Web Dashboard: http://localhost:3001"
echo "  PWA Mobile   : http://localhost:3002"
echo "  PostgreSQL   : localhost:5432"
echo "  Redis        : localhost:6379"
echo ""
