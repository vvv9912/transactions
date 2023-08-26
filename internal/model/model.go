package model

import "time"

type Users struct {
	ID      int64
	UserID  int64
	Account float64
}
type Transactions struct {
	ID                int64
	NumberTransaction string
	UserID            int64
	CreatedAt         time.Time
}
type Caches struct {
	ID                int64
	NumberTransaction string
	Status            int8
}
