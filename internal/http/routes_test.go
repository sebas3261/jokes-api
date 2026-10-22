package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/sebas3261/jokes-api/internal/http"
	"github.com/sebas3261/jokes-api/internal/jokes"
)

func TestNewRouter(t *testing.T) {
	repo := &stubRepository{
		getRandomFunc: func(ctx context.Context) (jokes.Joke, error) {
			return jokes.Joke{ID: 1, Content: "hi"}, nil
		},
	}
	service := jokes.NewService(repo)
	router := apphttp.NewRouter(service)

	tests := []struct {
		name           string
		method         string
		path           string
		wantStatusCode int
	}{
		{name: "health endpoint is wired up", method: http.MethodGet, path: "/health", wantStatusCode: http.StatusOK},
		{name: "random joke endpoint is wired up", method: http.MethodGet, path: "/jokes/random", wantStatusCode: http.StatusOK},
		{name: "unknown path returns 404", method: http.MethodGet, path: "/does-not-exist", wantStatusCode: http.StatusNotFound},
		{name: "wrong method on a known path returns 405", method: http.MethodPost, path: "/jokes/random", wantStatusCode: http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("got status %d, want %d", rec.Code, tt.wantStatusCode)
			}
		})
	}
}
