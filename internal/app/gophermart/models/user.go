package models

import "time"

type User struct {
	UserId    int       `json:"user_id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
