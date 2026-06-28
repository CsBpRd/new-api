#!/bin/bash
set -e

echo "=== Installing pnpm ==="
npm install -g pnpm

echo "=== Building frontend ==="
cd web
pnpm install --frozen-lockfile
pnpm --filter react-template run build
cd ..

mkdir -p public
cp -r web/classic/dist/* public/

echo "=== Frontend build complete ==="
