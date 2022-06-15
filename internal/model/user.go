package model

import "time"

type UserID string

type NilUserID UserID

type User struct {
	ID        UserID     `json:"id,omitempty"`
	Email     *string    `json:"email" db:"email"`
	Password  *[]byte    `json:"-" db:"password"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}
