package jokes_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sebas3261/jokes-api/internal/jokes"
)

func TestMemoryRepository_FindByID(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{name: "seeded joke 1 exists", id: 1, wantErr: false},
		{name: "seeded joke 2 exists", id: 2, wantErr: false},
		{name: "seeded joke 3 exists", id: 3, wantErr: false},
		{name: "unknown id returns error", id: 999, wantErr: true},
		{name: "zero id returns error", id: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joke, err := repo.FindByID(ctx, tt.id)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for id %d, got nil", tt.id)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if joke.ID != tt.id {
				t.Fatalf("got joke with ID %d, want %d", joke.ID, tt.id)
			}
		})
	}
}

func TestMemoryRepository_FindAll(t *testing.T) {
	repo := jokes.NewMemoryRepository()

	all, err := repo.FindAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("got %d jokes, want 3", len(all))
	}

	seen := map[int64]bool{}
	for _, j := range all {
		seen[j.ID] = true
	}
	for _, id := range []int64{1, 2, 3} {
		if !seen[id] {
			t.Errorf("seeded joke %d missing from FindAll result", id)
		}
	}
}

func TestMemoryRepository_FindAll_EmptyRepository(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	for _, id := range []int64{1, 2, 3} {
		if err := repo.Delete(ctx, id); err != nil {
			t.Fatalf("setup: failed to delete seeded joke %d: %v", id, err)
		}
	}

	all, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error on empty repository: %v", err)
	}
	if len(all) != 0 {
		t.Fatalf("got %d jokes, want 0", len(all))
	}
}

func TestMemoryRepository_Create(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	before := time.Now()
	created, err := repo.Create(ctx, jokes.Joke{Content: "new joke"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.ID != 4 {
		t.Fatalf("got ID %d, want 4 (first free ID after the 3 seeded jokes)", created.ID)
	}
	if created.Content != "new joke" {
		t.Fatalf("got content %q, want %q", created.Content, "new joke")
	}
	if created.CreatedAt.Before(before) {
		t.Fatalf("CreatedAt %v should not be before the call started at %v", created.CreatedAt, before)
	}

	fetched, err := repo.FindByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("created joke should be retrievable: %v", err)
	}
	if fetched != created {
		t.Fatalf("fetched joke %+v does not match created joke %+v", fetched, created)
	}
}

func TestMemoryRepository_Create_AssignsSequentialIDs(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	first, err := repo.Create(ctx, jokes.Joke{Content: "first"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := repo.Create(ctx, jokes.Joke{Content: "second"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if second.ID != first.ID+1 {
		t.Fatalf("got IDs %d then %d, want sequential IDs", first.ID, second.ID)
	}
}

func TestMemoryRepository_Create_IgnoresCallerSuppliedID(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	created, err := repo.Create(ctx, jokes.Joke{ID: 999, Content: "spoofed id"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created.ID != 4 {
		t.Fatalf("Create should assign its own sequential ID, got %d, want 4", created.ID)
	}
}

func TestMemoryRepository_Update(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	updated, err := repo.Update(ctx, jokes.Joke{ID: 1, Content: "updated content"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Content != "updated content" {
		t.Fatalf("got content %q, want %q", updated.Content, "updated content")
	}

	fetched, err := repo.FindByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fetched.Content != "updated content" {
		t.Fatalf("update was not persisted, got content %q", fetched.Content)
	}
}

func TestMemoryRepository_Update_UnknownIDReturnsError(t *testing.T) {
	repo := jokes.NewMemoryRepository()

	_, err := repo.Update(context.Background(), jokes.Joke{ID: 999, Content: "ghost"})
	if err == nil {
		t.Fatalf("expected error when updating a joke that does not exist")
	}
}

func TestMemoryRepository_Delete(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	if err := repo.Delete(ctx, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := repo.FindByID(ctx, 1); err == nil {
		t.Fatalf("expected joke 1 to be gone after deletion")
	}
}

func TestMemoryRepository_Delete_UnknownIDReturnsError(t *testing.T) {
	repo := jokes.NewMemoryRepository()

	err := repo.Delete(context.Background(), 999)
	if err == nil {
		t.Fatalf("expected error when deleting a joke that does not exist")
	}
}

func TestMemoryRepository_GetRandom_SingleJoke(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	// Narrow the store down to a single known joke so the result is deterministic.
	if err := repo.Delete(ctx, 2); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := repo.Delete(ctx, 3); err != nil {
		t.Fatalf("setup: %v", err)
	}

	for i := 0; i < 10; i++ {
		got, err := repo.GetRandom(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID != 1 {
			t.Fatalf("got joke ID %d, want 1 (the only remaining joke)", got.ID)
		}
	}
}

func TestMemoryRepository_GetRandom_DistributesAcrossAllJokes(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	seen := map[int64]bool{}
	const attempts = 200
	for i := 0; i < attempts; i++ {
		got, err := repo.GetRandom(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got.ID < 1 || got.ID > 3 {
			t.Fatalf("got unexpected joke ID %d", got.ID)
		}
		seen[got.ID] = true
	}

	if len(seen) < 2 {
		t.Fatalf("GetRandom returned only ID(s) %v across %d attempts; expected more variety", seen, attempts)
	}
}

func TestMemoryRepository_GetRandom_EmptyRepository(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	for _, id := range []int64{1, 2, 3} {
		if err := repo.Delete(ctx, id); err != nil {
			t.Fatalf("setup: failed to delete seeded joke %d: %v", id, err)
		}
	}

	_, err := repo.GetRandom(ctx)
	if err == nil {
		t.Fatalf("expected error when repository is empty")
	}
}

func TestMemoryRepository_ConcurrentAccess(t *testing.T) {
	repo := jokes.NewMemoryRepository()
	ctx := context.Background()

	const workers = 50
	var wg sync.WaitGroup
	wg.Add(workers * 3)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, _ = repo.Create(ctx, jokes.Joke{Content: "concurrent"})
		}()
		go func() {
			defer wg.Done()
			_, _ = repo.GetRandom(ctx)
		}()
		go func() {
			defer wg.Done()
			_, _ = repo.FindAll(ctx)
		}()
	}

	wg.Wait()
}
