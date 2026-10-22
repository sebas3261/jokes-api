package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apphttp "github.com/sebas3261/jokes-api/internal/http"
	"github.com/sebas3261/jokes-api/internal/jokes"
)

// stubRepository is a manual test double for jokes.Repository used to
// drive the HTTP handlers through a real *jokes.Service, without touching
// the network or an actual data store.
type stubRepository struct {
	getRandomFunc func(ctx context.Context) (jokes.Joke, error)
}

func (s *stubRepository) Create(ctx context.Context, joke jokes.Joke) (jokes.Joke, error) {
	return jokes.Joke{}, errors.New("stubRepository.Create: not implemented")
}

func (s *stubRepository) GetRandom(ctx context.Context) (jokes.Joke, error) {
	return s.getRandomFunc(ctx)
}

func (s *stubRepository) FindByID(ctx context.Context, id int64) (jokes.Joke, error) {
	return jokes.Joke{}, errors.New("stubRepository.FindByID: not implemented")
}

func (s *stubRepository) FindAll(ctx context.Context) ([]jokes.Joke, error) {
	return nil, errors.New("stubRepository.FindAll: not implemented")
}

func (s *stubRepository) Update(ctx context.Context, joke jokes.Joke) (jokes.Joke, error) {
	return jokes.Joke{}, errors.New("stubRepository.Update: not implemented")
}

func (s *stubRepository) Delete(ctx context.Context, id int64) error {
	return errors.New("stubRepository.Delete: not implemented")
}

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	apphttp.HealthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("got Content-Type %q, want %q", ct, "application/json")
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("got status field %q, want %q", body["status"], "ok")
	}
}

func TestGetRandomJokeHandler(t *testing.T) {
	sampleJoke := jokes.Joke{ID: 7, Content: "A joke", CreatedAt: time.Now()}

	tests := []struct {
		name           string
		getRandomFunc  func(ctx context.Context) (jokes.Joke, error)
		wantStatusCode int
		wantJoke       *jokes.Joke
	}{
		{
			name: "repository success returns the joke as JSON",
			getRandomFunc: func(ctx context.Context) (jokes.Joke, error) {
				return sampleJoke, nil
			},
			wantStatusCode: http.StatusOK,
			wantJoke:       &sampleJoke,
		},
		{
			name: "repository error becomes a 500 response",
			getRandomFunc: func(ctx context.Context) (jokes.Joke, error) {
				return jokes.Joke{}, errors.New("no jokes available")
			},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubRepository{getRandomFunc: tt.getRandomFunc}
			service := jokes.NewService(repo)
			handler := apphttp.GetRandomJokeHandler(service)

			req := httptest.NewRequest(http.MethodGet, "/jokes/random", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Fatalf("got status %d, want %d", rec.Code, tt.wantStatusCode)
			}

			if tt.wantJoke == nil {
				if rec.Body.Len() == 0 {
					t.Fatalf("expected an error message in the response body")
				}
				return
			}

			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("got Content-Type %q, want %q", ct, "application/json")
			}

			var got jokes.Joke
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("response body is not valid JSON: %v", err)
			}

			// CreatedAt is compared with Equal, not ==: after a JSON round
			// trip the time.Time is reconstructed with a different
			// internal Location value, so == would report a false
			// mismatch even for the same instant.
			if got.ID != tt.wantJoke.ID || got.Content != tt.wantJoke.Content {
				t.Fatalf("got joke %+v, want %+v", got, *tt.wantJoke)
			}
			if !got.CreatedAt.Equal(tt.wantJoke.CreatedAt) {
				t.Fatalf("got CreatedAt %v, want %v", got.CreatedAt, tt.wantJoke.CreatedAt)
			}
		})
	}
}
