# My App

A full-stack web application built with **Go**, **Gin**, **React**, **TypeScript**, and **Vite**.

## Tech Stack

### Backend

* Go
* Gin
* REST API

### Frontend

* React
* TypeScript
* Vite

## Project Structure

```text
my-app/
├── api/        # Go + Gin backend
├── web/        # React + Vite frontend
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

Verify your installations:

```bash
go version
node --version
npm --version
git --version
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

Or run each manually:

```bash
cd api
go run ./cmd/server
```

```bash
cd web
npm run dev
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

## Environment Variables

Environment configuration has not been added yet.

Eventually, backend variables may live in:

```text
api/.env
```

Frontend variables may live in:

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

## Git

Initial commit:

```bash
git add .
git commit -m "Initial Go Gin API and React app"
```

## License

This project does not currently specify a license.

