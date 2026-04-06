#!/bin/bash
# ============================================
# ERP SPPG - Build All Docker Images
# Usage: ./deployment/scripts/build.sh [tag]
# ============================================
set -e

TAG="${1:-latest}"
REGISTRY="${DOCKER_REGISTRY:-}"

echo "=== Building ERP SPPG Docker Images (tag: $TAG) ==="

# Build backend
echo ""
echo ">>> Building backend..."
docker build -t erp-sppg-backend:$TAG ./backend

# Build web dashboard
echo ""
echo ">>> Building web dashboard..."
docker build -t erp-sppg-web:$TAG ./web

# Build PWA mobile
echo ""
echo ">>> Building PWA mobile..."
docker build -t erp-sppg-pwa:$TAG ./pwa

echo ""
echo "=== All images built successfully ==="
echo ""
docker images | grep erp-sppg

# Optional: push to registry
if [ -n "$REGISTRY" ]; then
  echo ""
  echo ">>> Pushing to registry: $REGISTRY"
  for img in backend web pwa; do
    docker tag erp-sppg-$img:$TAG $REGISTRY/erp-sppg-$img:$TAG
    docker push $REGISTRY/erp-sppg-$img:$TAG
  done
  echo ">>> Push complete"
fi
