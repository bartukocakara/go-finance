package model

import (
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserID string

var NilUserID UserID

type User struct {
	ID        UserID     `json:"id,omitempty" db:"user_id"`
	Email     *string    `json:"email" db:"email"`
	Password  *[]byte    `json:"-" db:"password"`
	CreatedAt *time.Time `json:"-" db:"created_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

var isEmail = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

func (u *User) Verify() error {
	if u.Email == nil || (u.Email != nil && len(*u.Email) == 0) {
		return errors.New("Email is required")
	}
	if !isEmail.MatchString(*u.Email) {
		return errors.New("Email invalid")
	}
	return nil
}

func (u *User) CheckPassword(password string) error {
	if u.Password != nil && len(*u.Password) == 0 {
		return errors.New("Password not set")
	}
	return bcrypt.CompareHashAndPassword(*u.Password, []byte(password))
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
