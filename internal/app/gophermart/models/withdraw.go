package models

import "time"

type Withdraw struct {
	ID          int       `json:"-"`
	UserID      int       `json:"-"`
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
	CreatedAt   time.Time `json:"-"`
}

type WithdrawalRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}
