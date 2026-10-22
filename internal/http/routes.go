package http

import (
    "net/http"

    "github.com/sebas3261/jokes-api/internal/jokes"
)

func NewRouter(service *jokes.Service) *http.ServeMux {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /health", HealthHandler)
    mux.HandleFunc("GET /jokes/random", GetRandomJokeHandler(service))

    return mux
}