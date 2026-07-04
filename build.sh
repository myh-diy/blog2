#!/bin/bash
set -e

echo "Building frontend..."
cd frontend
npm run build
cd ..

echo "Copying frontend dist to backend..."
rm -rf backend/frontend-dist
cp -r frontend/dist backend/frontend-dist

echo "Building backend..."
cd backend
go build -o blog-server .
cd ..

echo "Build complete: backend/blog-server"
