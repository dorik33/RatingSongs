package users

import "time"

type Users struct {
	ID        int       `json:id`
	Username  string    `json:username`
	Email     string    `json:email`
	CreatedAt time.Time `json:created_at`
}
