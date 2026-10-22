package jokes_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sebas3261/jokes-api/internal/jokes"
)

// mockRepository is a hand-written test double for jokes.Repository.
// Each method delegates to an optional function field, so each test
// configures only the behavior it actually needs.
type mockRepository struct {
	findByIDFunc  func(ctx context.Context, id int64) (jokes.Joke, error)
	getRandomFunc func(ctx context.Context) (jokes.Joke, error)

	findByIDCalls  int
	getRandomCalls int
}

func (m *mockRepository) Create(ctx context.Context, joke jokes.Joke) (jokes.Joke, error) {
	return jokes.Joke{}, errors.New("mockRepository.Create: not implemented")
}

func (m *mockRepository) GetRandom(ctx context.Context) (jokes.Joke, error) {
	m.getRandomCalls++
	if m.getRandomFunc == nil {
		return jokes.Joke{}, errors.New("mockRepository.GetRandom: not implemented")
	}
	return m.getRandomFunc(ctx)
}

func (m *mockRepository) FindByID(ctx context.Context, id int64) (jokes.Joke, error) {
	m.findByIDCalls++
	if m.findByIDFunc == nil {
		return jokes.Joke{}, errors.New("mockRepository.FindByID: not implemented")
	}
	return m.findByIDFunc(ctx, id)
}

func (m *mockRepository) FindAll(ctx context.Context) ([]jokes.Joke, error) {
	return nil, errors.New("mockRepository.FindAll: not implemented")
}

func (m *mockRepository) Update(ctx context.Context, joke jokes.Joke) (jokes.Joke, error) {
	return jokes.Joke{}, errors.New("mockRepository.Update: not implemented")
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	return errors.New("mockRepository.Delete: not implemented")
}

func TestService_GetJokeByID(t *testing.T) {
	sampleJoke := jokes.Joke{ID: 1, Content: "Why did the chicken cross the road?", CreatedAt: time.Now()}

	tests := []struct {
		name          string
		id            int64
		findByIDFunc  func(ctx context.Context, id int64) (jokes.Joke, error)
		wantJoke      jokes.Joke
		wantErr       bool
		wantRepoCalls int
	}{
		{
			name:          "zero id is rejected before reaching the repository",
			id:            0,
			wantErr:       true,
			wantRepoCalls: 0,
		},
		{
			name:          "negative id is rejected before reaching the repository",
			id:            -5,
			wantErr:       true,
			wantRepoCalls: 0,
		},
		{
			name: "existing id returns the joke from the repository",
			id:   1,
			findByIDFunc: func(ctx context.Context, id int64) (jokes.Joke, error) {
				return sampleJoke, nil
			},
			wantJoke:      sampleJoke,
			wantErr:       false,
			wantRepoCalls: 1,
		},
		{
			name: "repository error is propagated unchanged",
			id:   99,
			findByIDFunc: func(ctx context.Context, id int64) (jokes.Joke, error) {
				return jokes.Joke{}, errors.New("joke not found")
			},
			wantErr:       true,
			wantRepoCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{findByIDFunc: tt.findByIDFunc}
			service := jokes.NewService(repo)

			got, err := service.GetJokeByID(context.Background(), tt.id)

			if tt.wantErr && err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && got != tt.wantJoke {
				t.Fatalf("got joke %+v, want %+v", got, tt.wantJoke)
			}
			if repo.findByIDCalls != tt.wantRepoCalls {
				t.Fatalf("repository.FindByID called %d times, want %d", repo.findByIDCalls, tt.wantRepoCalls)
			}
		})
	}
}

func TestService_GetRandomJoke(t *testing.T) {
	sampleJoke := jokes.Joke{ID: 2, Content: "A random joke", CreatedAt: time.Now()}

	tests := []struct {
		name          string
		getRandomFunc func(ctx context.Context) (jokes.Joke, error)
		wantJoke      jokes.Joke
		wantErr       bool
	}{
		{
			name: "success returns whatever the repository returns",
			getRandomFunc: func(ctx context.Context) (jokes.Joke, error) {
				return sampleJoke, nil
			},
			wantJoke: sampleJoke,
		},
		{
			name: "repository error (e.g. empty store) is propagated",
			getRandomFunc: func(ctx context.Context) (jokes.Joke, error) {
				return jokes.Joke{}, errors.New("no jokes available")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{getRandomFunc: tt.getRandomFunc}
			service := jokes.NewService(repo)

			got, err := service.GetRandomJoke(context.Background())

			if tt.wantErr && err == nil {
				t.Fatalf("expected an error, got nil")
			}
			if !tt.wantErr {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got != tt.wantJoke {
					t.Fatalf("got joke %+v, want %+v", got, tt.wantJoke)
				}
			}
			if repo.getRandomCalls != 1 {
				t.Fatalf("repository.GetRandom called %d times, want 1", repo.getRandomCalls)
			}
		})
	}
}
