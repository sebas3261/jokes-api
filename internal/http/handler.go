package http

import (
    "encoding/json"
    "net/http"

    "github.com/sebas3261/jokes-api/internal/jokes"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`{"status":"ok"}`))
}

func GetRandomJokeHandler(service *jokes.Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        joke, err := service.GetRandomJoke(r.Context())

        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(joke)
    }
}