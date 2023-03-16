package models

import "time"

type Order struct {
	Id        int       `json:"id"`
	OrderId   string    `json:"order_id"`
	UserId    int       `json:"user_id,omitempty"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}
