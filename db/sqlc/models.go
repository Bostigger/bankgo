// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"database/sql"
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID int64 `json:"id"`
	// can be positive or negative
	Amount    int64         `json:"amount"`
	AccountID sql.NullInt64 `json:"account_id"`
	CreatedAt time.Time     `json:"created_at"`
}

type Transfer struct {
	ID int64 `json:"id"`
	// should be only positive
	Amount     int64     `json:"amount"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	CreatedAt  time.Time `json:"created_at"`
}
