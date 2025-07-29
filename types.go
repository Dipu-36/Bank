package main

import (
	"math/rand"
	"time"
)

// CreateAccountRequest represents the expected input from a client
// when creating a new account via the API,
// It includes only the fields that user is allowed to set
type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"llast_name"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(1000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(100000)),
		CreatedAt: time.Now().UTC(),
	}
}
