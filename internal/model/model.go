package model

import "time"

type Users struct {
	ID      int64
	Account float64
}
type Transactions struct {
	ID        int64
	Users_ID  int64
	CreatedAt time.Time
}
