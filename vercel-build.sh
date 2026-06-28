#!/bin/bash
set -e

echo "=== Installing pnpm ==="
npm install -g pnpm

echo "=== Building frontend ==="
cd web

# Install dependencies in workspace root
pnpm install --frozen-lockfile

# Build classic template
cd classic
pnpm run build
cd ..

# Copy built files to public directory
mkdir -p ../public
cp -r classic/dist/* ../public/

echo "=== Frontend build complete ==="
