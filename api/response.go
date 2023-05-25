package api

import "time"

type ResponseBook struct {
	Id         uint64    `json:"id"`
	Title      string    `json:"title"`
	Authors    []string  `json:"authors"`
	Publisher  string    `json:"publisher"`
	Isbn       string   `json:"isbn"`
	Price      int64     `json:"price"`
	Quantity   int64     `json:"quantity"`
	Created_by string    `json:"created_by"`
	Created_at time.Time `json:"created_at"`
}

type ResponseUser struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	Fullname          string    `json:"fullname"`
	HashedPassword    string    `json:"hashed_password"`
	CreatedAt         time.Time `json:"created_at"`
}
