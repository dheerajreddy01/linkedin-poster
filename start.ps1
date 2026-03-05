# PostPilot — Windows Startup
Write-Host "PostPilot — LinkedIn Post Automation" -ForegroundColor Cyan
if (-not (Get-Command go -ErrorAction SilentlyContinue)) { Write-Host "ERROR: Go not installed" -ForegroundColor Red; exit 1 }
if (-not (Test-Path ".env")) { Write-Host "ERROR: .env not found!" -ForegroundColor Red; exit 1 }
New-Item -ItemType Directory -Force -Path "data" | Out-Null
go mod tidy
Write-Host "Starting at http://localhost:8081" -ForegroundColor Green
go run cmd\server\main.go
