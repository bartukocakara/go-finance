package model

import (
	"errors"
	"time"
)

type MerchantID string

var NilMerchantID MerchantID

type Merchant struct {
	ID        MerchantID `json:"id,omitempty" db:"merchant_id"`
	UserID    *UserID    `json:"userID,omitempty" db:"user_id"`
	CreatedAt *time.Time `json:"createdAt,omitempty" db:"created_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
	Name      *string    `json:"name,omitempty" db:"name"`
}

func (m *Merchant) Verify() error {
	if m.UserID == nil || len(*m.UserID) == 0 {
		return errors.New("UserID is required")
	}
	if m.Name == nil || len(*m.Name) == 0 {
		return errors.New("Name is required")
	}

	return nil
}
