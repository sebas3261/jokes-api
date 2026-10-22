package jokes

import (
    "context"
    "errors"
    "sync"
	"math/rand/v2"
	"time"
)

type MemoryRepository struct {
    mu     sync.RWMutex
    jokes  map[int64]Joke
    nextID int64
}

func NewMemoryRepository() *MemoryRepository {
    return &MemoryRepository{
        jokes: map[int64]Joke{
            1: {
                ID:      1,
                Content: "Why do programmers prefer dark mode? Because light attracts bugs.",
                CreatedAt: time.Now(),
            },
            2: {
                ID:      2,
                Content: "Why do Java developers wear glasses? Because they don't C#.",
                CreatedAt: time.Now(),
            },
            3: {
                ID:      3,
                Content: "There are only 10 types of people: those who understand binary and those who don't.",
                CreatedAt: time.Now(),
            },
        },
        nextID: 4,
    }
}

func (r *MemoryRepository) FindByID(ctx context.Context, id int64) (Joke, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    joke, exists := r.jokes[id]

    if !exists {
        return Joke{}, errors.New("joke not found")
    }

    return joke, nil
}

func (r *MemoryRepository) Create(ctx context.Context, joke Joke) (Joke, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    joke.ID = r.nextID
    joke.CreatedAt = time.Now()
    r.jokes[joke.ID] = joke
    r.nextID++

    return joke, nil
}

func (r *MemoryRepository) FindAll(ctx context.Context) ([]Joke, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    result := make([]Joke, 0, len(r.jokes))
    for _, joke := range r.jokes {
        result = append(result, joke)
    }

    return result, nil
}

func (r *MemoryRepository) Update(ctx context.Context, joke Joke) (Joke, error) {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.jokes[joke.ID]; !exists {
        return Joke{}, errors.New("joke not found")
    }

    r.jokes[joke.ID] = joke

    return joke, nil
}

func (r *MemoryRepository) Delete(ctx context.Context, id int64) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.jokes[id]; !exists {
        return errors.New("joke not found")
    }

    delete(r.jokes, id)

    return nil
}

func (r *MemoryRepository) GetRandom(ctx context.Context) (Joke, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    if len(r.jokes) == 0 {
        return Joke{}, errors.New("no jokes available")
    }

    index := rand.IntN(len(r.jokes))

    for _, joke := range r.jokes {
        if index == 0 {
            return joke, nil
        }
        index--
    }

    return Joke{}, errors.New("no jokes available")
}