# CashLenX Server

Go backend providing a CLI and REST API for personal finance management.

## Features
- Track income and expenses
- Category-based organization
- Date range queries and summaries
- Excel import/export
- MongoDB and MySQL support
- Structured logging and validation

## Project Structure
```
cashlenx-server/
├── cmd/                 # CLI commands (Cobra)
├── controller/          # HTTP controllers
├── service/             # Business logic
├── mapper/              # Database mappers
├── model/               # Data models
├── middleware/          # HTTP middleware
├── util/                # Utilities
├── docker/              # Database init assets
├── docs/                # Server documentation
└── main.go              # Entry point
```

## Quick Start

### 1) Database
```bash
cd cashlenx-server
docker compose up -d mongodb
# or
docker compose --profile mysql up -d mysql
```

### 2) Configure
```bash
cd cashlenx-server
cp .env.sample .env
export $(cat .env | xargs)
```

### 3) Run
```bash
# API server
go run main.go server start -p 8080

# CLI examples
go run main.go cash outcome -c "Food" -a 45.50 -d "Lunch"
go run main.go cash income -c "Salary" -a 5000
go run main.go cash summary -f 2024-01-01 -t 2024-01-31
```

## REST API
- `POST /api/cash/outcome`
- `POST /api/cash/income`
- `GET /api/cash/{id}`
- `GET /api/cash/date/{date}`
- `DELETE /api/cash/{id}`
- `DELETE /api/cash/date/{date}`
- `GET /api/health`
- `GET /api/version`

See `docs/api.md` for detailed endpoints.

## Documentation
- `docs/cli.md` — CLI command reference
- `docs/api.md` — REST API reference
- `docs/testing.md` — Testing guide
- `docs/deployment_guide.md` — Deployment guide

## Build and Test
```bash
go build -o cashlenx main.go
go test ./...
```

## Technology
- Go 1.21+
- Cobra, Gorilla Mux, Zap, Excelize
- MongoDB and MySQL drivers

## License
See `LICENSE` for details.
