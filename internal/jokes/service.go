package jokes

import (
    "context"
    "errors"
)

type Service struct {
    repository Repository
}

func NewService(repository Repository) *Service {
    return &Service{
        repository: repository,
    }
}

func (s *Service) GetJokeByID(ctx context.Context, id int64) (Joke, error) {
    if id <= 0 {
        return Joke{}, errors.New("invalid joke ID")
    }

    return s.repository.FindByID(ctx, id)
}

func (s *Service) GetRandomJoke(ctx context.Context) (Joke, error){
	return s.repository.GetRandom(ctx)
}