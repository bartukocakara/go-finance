package model

import (
	"errors"
	"time"
)

type AccountID string

var NilAccountID AccountID

type AccountType string

const (
	Cash   AccountType = "cash"
	Credit AccountType = "credit"
)

type Account struct {
	ID           AccountID    `json:"id,omitempty" db:"account_id"`
	UserID       *UserID      `json:"UserID,omitempty" db:"user_id"`
	Name         *string      `json:"account_name,omitempty" db:"account_name"`
	Type         *AccountType `json:"type,omitempty" db:"account_type"`
	StartBalance *float64     `json:"start_balance,omitempty" db:"start_balance"`
	Currency     *string      `json:"currency,omitempty" db:"currency"`
	CreatedAt    *time.Time   `json:"createdAt,omitempty" db:"created_at"`
	DeletedAt    *time.Time   `json:"deletedAt,omitempty" db:"deleted_at"`
}

func (a *Account) Verify() error {
	if a.UserID == nil || len(*a.UserID) == 0 {
		return errors.New("UserID is required")
	}

	if a.Name == nil || len(*a.Name) == 0 {
		return errors.New("Name is required")
	}

	if a.Type == nil || len(*a.Type) == 0 {
		return errors.New("Type is required")
	}

	if a.StartBalance == nil {
		return errors.New("Start Balance is required")
	}

	if a.Currency == nil || len(*a.Currency) == 0 {
		return errors.New("Currency is required")
	}

	return nil
}

func (a *Account) SetName(name *string) *string {
	if a.Name != nil {
		a.Name = name
	}
	return a.Name
}

func (a *Account) SetStartBalance(startBalance *float64) *float64 {
	if a.StartBalance != nil {
		a.StartBalance = startBalance
	}
	return a.StartBalance
}

func (a *Account) SetCurrency(currency *string) *string {
	if a.Currency != nil || len(*a.Currency) == 0 {
		a.Currency = currency
	}
	return a.Currency
}

func (a *Account) SetType(Type *AccountType) *AccountType {
	if a.Type == nil || len(*a.Type) == 0 {
		a.Type = Type
	}

	return a.Type
}

func (a *Account) SetAll(account *Account) *Account {
	a.Name = account.SetName(account.Name)
	a.StartBalance = account.SetStartBalance(account.StartBalance)
	a.Currency = account.SetCurrency(account.Currency)
	a.Type = account.SetType(account.Type)

	return a
}
