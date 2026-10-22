package main

import (
    "log"
    "net/http"

    apphttp "github.com/sebas3261/jokes-api/internal/http"
    "github.com/sebas3261/jokes-api/internal/jokes"
)

func main() {
    repository := jokes.NewMemoryRepository()
    service := jokes.NewService(repository)
    router := apphttp.NewRouter(service)

    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}