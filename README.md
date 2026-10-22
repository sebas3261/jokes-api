# Jokes API

A small REST API built with **Go** to learn the language and explore backend development using its standard library.

This project focuses on understanding Go fundamentals, HTTP handlers, interfaces, dependency injection, concurrency safety, and automated testing.

## Features

- Retrieve a random programming joke.
- Retrieve jokes by ID.
- In-memory data storage.
- Layered architecture with repository and service abstractions.
- JSON-based HTTP responses.
- Health check endpoint.
- Unit and HTTP tests.

## Tech Stack

- **Go**
- **net/http** — HTTP server and routing
- **encoding/json** — JSON serialization
- **sync.RWMutex** — Thread-safe in-memory storage
- **testing** — Unit tests
- **net/http/httptest** — HTTP handler tests

## Project Structure

```text
jokes-api/
├── internal/
│   ├── http/
│   │   ├── handler.go
│   │   └── routes.go
│   └── jokes/
│       ├── model.go
│       ├── repository.go
│       ├── memory_repository.go
│       └── service.go
├── go.mod
└── main.go
```

Test files are located alongside the code they test.

## Getting Started

### Prerequisites

- Go 1.22 or later

### Run locally

Clone the repository:

```bash
git clone https://github.com/sebas3261/jokes-api.git
cd jokes-api
```

Start the server:

```bash
go run .
```

The API will be available at `http://localhost:8080`.

## API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| GET | `/health` | Check API health |
| GET | `/jokes/random` | Retrieve a random joke |

### Example Request

```bash
curl http://localhost:8080/jokes/random
```

Example response:

```json
{
  "id": 2,
  "content": "Why do Java developers wear glasses? Because they don't C#.",
  "created_at": "2026-10-21T21:38:29-05:00"
}
```

## Testing

Run all tests:

```bash
go test ./...
```

Run tests with the race detector:

```bash
go test -race ./...
```

Generate a coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## What I Learned

This project helped me practice:

- Go syntax, structs, pointers, and methods.
- Interfaces and implicit interface implementation.
- Dependency injection without external frameworks.
- Separating HTTP, business logic, and data access.
- Error handling using Go's explicit error returns.
- Concurrency-safe operations with mutexes.
- Writing unit tests and testing HTTP handlers.

## Project Purpose

This is a personal learning project, not a production-ready service. Its primary goal is to build a solid foundation in Go through practical backend development.

And yes, the jokes are intentionally terrible.