#!/bin/bash
# PostPilot — LinkedIn Post Automation
# Usage: chmod +x start.sh && ./start.sh

echo "========================================"
echo "  PostPilot — LinkedIn Post Automation"
echo "========================================"

if ! command -v go &> /dev/null; then
  echo "ERROR: Go not installed. https://go.dev/dl/"
  exit 1
fi

[ ! -f ".env" ] && echo "ERROR: .env not found!" && exit 1

mkdir -p data frontend
go mod tidy

echo "Starting at http://localhost:8081"
echo "Open http://localhost:8081 in browser"
echo ""
go run cmd/server/main.go
