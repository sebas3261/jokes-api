package jokes

import "context"

type Repository interface {
    Create(ctx context.Context, joke Joke) (Joke, error)
	GetRandom(ctx context.Context) (Joke, error)
    FindByID(ctx context.Context, id int64) (Joke, error)
    FindAll(ctx context.Context) ([]Joke, error)
    Update(ctx context.Context, joke Joke) (Joke, error)
    Delete(ctx context.Context, id int64) error
}