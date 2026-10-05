package model

import "time"

type Data struct {
	Id        int    `json:"id"`
	Operation string `json:"operation"`
	Amount    int    `json:"amount"`
}
type Event struct {
	ID        uint64    `json:"id"`
	Operation string    `json:"operation"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Currency struct {
	USD float32
	EUR float32
	RUB float32
}
