package jokes

import "time"

type Joke struct {
    ID        int64     `json:"id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}