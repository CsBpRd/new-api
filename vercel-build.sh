#!/bin/bash
set -e

echo "=== Building frontend ==="
cd web/classic

# Install dependencies using pnpm (project uses catalog: protocol)
if ! command -v pnpm &> /dev/null; then
    echo "Installing pnpm..."
    npm install -g pnpm
fi

pnpm install --frozen-lockfile

# Build frontend
pnpm run build

cd ../..

# Create public directory for static files
mkdir -p public
cp -r web/classic/dist/* public/

echo "=== Frontend build complete ==="
echo "=== Building Go API ==="

# Create api directory
mkdir -p api

# Copy the Go source files needed
cp main.go api/
cp -r common api/
cp -r constant api/
cp -r controller api/
cp -r middleware api/
cp -r model api/
cp -r oauth api/
cp -r relay api/
cp -r router api/
cp -r service api/
cp -r setting api/
cp -r types api/
cp -r i18n api/
cp -r logger api/
cp -r pkg api/
cp go.mod api/
cp go.sum api/

echo "=== Build complete ==="
