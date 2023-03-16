package models

import "time"

type Order struct {
	ID        int       `json:"id"`
	OrderID   string    `json:"order_id"`
	UserID    int       `json:"user_id,omitempty"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
