# Raceday Checklist

Raceday Checklist is a full-stack web app for managing RC car race-day prep.
It is meant to be a repeatable routine checklist for the tasks and maintenance
needed to keep RC car(s) ready throughout an event.

The checklist is organized around the rhythm of a race day:

* Before practice
* Before qualifying
* Mid-qualifying
* Before the main

The goal is to make it easy to track what still needs attention between runs:
car setup, maintenance, charging, tires, transponders, tools, and other race-day
tasks.

## Tech Stack

### Backend

* Go
* Gin
* REST API
* MySQL database

### Frontend

* React
* TypeScript
* Vite
* Pico CSS
* SCSS

## Project Structure

```text
raceday-checklist/
├── api/        # Go + Gin backend
├── web/        # React + TypeScript + Vite frontend
├── Makefile
├── .gitignore
└── README.md
```

## Getting Started

### Prerequisites

Make sure you have the following installed:

* Go
* Node.js
* npm
* Git
* MySQL

Verify your installations:

```bash
go version
node --version
npm --version
git --version
mysql --version
```

## Backend Setup

From the project root:

```bash
cd api
go mod tidy
go run ./cmd/server
```

The API will run at:

```text
http://localhost:8080
```

Test the health endpoint:

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{
  "status": "ok"
}
```

## Frontend Setup

From the project root:

```bash
cd web
npm install
npm run dev
```

The frontend will run at:

```text
http://localhost:5173
```

## Development

Run the backend:

```bash
make api
```

Run the frontend:

```bash
make web
```

Run backend tests:

```bash
make test
```

Format backend code:

```bash
make fmt
```

## API Endpoints

### Health Check

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

### Checklist API

```http
GET /api/checklist
```

Response:

```json
{
  "items": [
    {
      "id": "fuel",
      "title": "Fuel and fluids checked",
      "category": "Car",
      "done": false
    }
  ]
}
```

### Hello API

```http
GET /api/hello
```

Response:

```json
{
  "message": "Hello from Go + Gin"
}
```

## Database

The backend is intended to connect to a MySQL database for persisted checklist
data. Backend configuration is loaded from environment variables. For local
development, place them in:

```text
api/.env
```

Use `api/.env.example` as the starting point:

```dotenv
PORT=8080

DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=raceday
DB_PASSWORD=raceday_password
DB_NAME=raceday_checklist
```

`api/.env` is ignored by Git so local database credentials are not committed.

Frontend variables may live in:

```text
web/.env
```

```text
web/.env
```

## Build

### Build Frontend

```bash
cd web
npm run build
```

The production frontend build will be output to:

```text
web/dist
```

### Build Backend

```bash
cd api
go build -o bin/server ./cmd/server
```

Run the built backend:

```bash
./bin/server
```

## License

This project does not currently specify a license.
