package model

import (
	"errors"
	"time"
)

type TransactionID string

var NilTransactionID TransactionID

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

var CurrencyTypes = []string{"USD", "EUR"}

type Transaction struct {
	ID         TransactionID `json:"id" db:"transaction_id"`
	UserID     *UserID       `json:"user_id" db:"user_id"`
	AccountID  *AccountID    `json:"account_id" db:"account_id"`
	CategoryID *CategoryID   `json:"category_id" db:"category_id"`

	Date     *time.Time       `json:"date" db:"date"`
	Type     *TransactionType `json:"type" db:"type"`
	Amount   *int64           `json:"amount" db:"amount"`
	Currency *string          `json:"currency" db:"currency"`
	Notes    *string          `json:"notes" db:"notes"`

	CreatedAt *time.Time `json:"createdAt,omitempty" db:"created_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

func (t *Transaction) Verify() error {
	if t.UserID == nil || len(*t.UserID) == 0 {
		return errors.New("UserID is required")
	}

	if t.AccountID == nil || len(*t.AccountID) == 0 {
		return errors.New("AccountID is required")
	}

	if t.CategoryID == nil || len(*t.CategoryID) == 0 {
		return errors.New("CategoryID is required")
	}

	if t.Date == nil {
		return errors.New("Date is required")
	}

	if t.Type == nil || len(*t.Type) == 0 {
		return errors.New("Type is required")
	}

	if t.Amount == nil {
		return errors.New("Amount is required")
	}

	if t.Currency == nil {
		return errors.New("Currency is required")
	}

	return nil
}

func (t *Transaction) SetUserID(userID *UserID) *UserID {
	t.UserID = userID
	return t.UserID
}

func (t *Transaction) SetAccountID(AccountID *AccountID) *AccountID {
	t.AccountID = AccountID
	return t.AccountID
}

func (t *Transaction) SetCategoryID(CategoryID *CategoryID) *CategoryID {
	t.CategoryID = CategoryID
	return t.CategoryID
}

func (t *Transaction) SetDate() *time.Time {
	Date := time.Now()
	t.Date = &Date
	return t.Date
}

func (t *Transaction) SetAmount(Amount *int64) *int64 {
	t.Amount = Amount
	return t.Amount
}

func (t *Transaction) SetCurrency(Currency *string) *string {
	t.Currency = Currency
	return t.Currency
}

func (t *Transaction) SetType(TransactionType *TransactionType) *TransactionType {
	t.Type = TransactionType
	return t.Type
}

// func (t *Transaction) VerifyCurrency(Currency string) (string, error){
// 	for _, currency := range CurrencyTypes {
// 		if Currency != currency {
// 			return "", errors.New("Wrong currency Type")
// 		}
// 	}

// 	return t.Currency, nil
// }

func (t *Transaction) SetNotes(Notes *string) *string {
	if Notes != nil {
		t.Notes = Notes
	}
	return t.Notes
}

func (t *Transaction) SetAll(userID *UserID, transaction *Transaction) *Transaction {
	t.UserID = transaction.SetUserID(userID)
	t.AccountID = transaction.SetAccountID(transaction.AccountID)
	t.CategoryID = transaction.SetCategoryID(transaction.CategoryID)
	t.Date = transaction.SetDate()
	t.Notes = transaction.SetNotes(transaction.Notes)
	t.Amount = transaction.SetAmount(transaction.Amount)
	t.Currency = transaction.SetCurrency(transaction.Currency)
	return t
}
