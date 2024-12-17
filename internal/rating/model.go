package rating

import "time"

type Rating struct {
	ID          int       `json:"id"`
	SongID      int       `json:"song_id"`
	Rating      int       `json:"rating"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
