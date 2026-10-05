package model

import "time"

type Event struct {
	ID        uint64    `json:"id"`
	Operation string    `json:"operation"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
