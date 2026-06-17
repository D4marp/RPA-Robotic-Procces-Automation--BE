# RPA Backend (Go)

REST API and Task Scheduler for the Robotic Process Automation (RPA) module.

## Stack
- **Language**: Go 1.22
- **Framework**: Gin Gonic
- **ORM**: GORM
- **Database**: SQLite (default for development/local) / MySQL (production)
- **Scheduler**: robfig/cron/v3

## Quick Start

1. Copy env file:
   ```bash
   cp .env.example .env
   ```
2. Download dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run cmd/main.go
   ```
   The backend API will start on `http://localhost:8080`.
