package database

import (
	"context"
	"financial-app/internal/model"
	"time"

	"github.com/pkg/errors"
)

type TransactionDB interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	ListAllTransactions(ctx context.Context, from, to time.Time) ([]*model.Transaction, error)
	GetTransactionByID(ctx context.Context, transactionID model.TransactionID) (*model.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) error
	DeleteTransaction(ctx context.Context, transactionID model.TransactionID) error
	ListTransactionByCategoryID(ctx context.Context, categoryID model.CategoryID, from, to time.Time) ([]*model.Transaction, error)
	ListTransactionByAccountID(ctx context.Context, accountID model.AccountID, from, to time.Time) ([]*model.Transaction, error)
	ListTransactionByUserID(ctx context.Context, userID model.UserID, from, to time.Time) ([]*model.Transaction, error)
}

var createTransactionQuery = `
	INSERT INTO transactions (user_id, account_id, category_id, date, type, amount, currency, notes)
		VALUES (:user_id, :account_id, :category_id, :date, :type, :amount, :currency, :notes)
	RETURNING transaction_id;
`

func (d *database) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	rows, err := d.conn.NamedQueryContext(ctx, createTransactionQuery, transaction)
	if err != nil {
		return err
	}
	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&transaction.ID); err != nil {
		return err
	}

	return nil
}

var listTransactionByCategoryID = `
	SELECT transaction_id, user_id, account_id, category_id, date, type, amount, currency, notes, created_at, deleted_at
	FROM transactions
	WHERE category_id = $1
		AND deleted_at IS NULL
		AND date > $2
		AND date < $3;
`

func (d *database) ListTransactionByCategoryID(ctx context.Context, categoryID model.CategoryID, from, to time.Time) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	if err := d.conn.SelectContext(ctx, &transactions, listTransactionByCategoryID, categoryID, from, to); err != nil {
		return nil, errors.Wrap(err, "Error: Could not get categories transactions")

	}
	return transactions, nil
}

var listAllTransactionsQuery = `
	SELECT transaction_id, user_id, category_id, account_id, type, amount, currency, notes, date, created_at, deleted_at
		FROM transactions
	WHERE deleted_at IS NULL
		AND date > $1
		AND date < $2
`

func (d *database) ListAllTransactions(ctx context.Context, from, to time.Time) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	if err := d.conn.SelectContext(ctx, &transactions, listAllTransactionsQuery, from, to); err != nil {
		return nil, errors.Wrap(err, "Error: Could not get categories transactions")

	}
	return transactions, nil
}

var getTransactionByIDQuery = `
	SELECT transaction_id, category_id, account_id, type, notes, amount, currency, date, created_at, deleted_at
	FROM transactions
	WHERE transaction_id = $1
		AND deleted_at IS NULL;
`

func (d *database) GetTransactionByID(ctx context.Context, transactionID model.TransactionID) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := d.conn.GetContext(ctx, &transaction, getTransactionByIDQuery, transactionID); err != nil {
		return nil, errors.Wrap(err, "Could not get transaction")
	}

	return &transaction, nil
}

const listTransactionByAccountIDQuery = `
	SELECT transaction_id, user_id, account_id, category_id, date, type, amount, notes, created_at, deleted_at 
	FROM transactions 
	WHERE account_id = $1 
		AND deleted_at IS NULL
		AND date > $2 
		AND date < $3;
`

func (d *database) ListTransactionByAccountID(ctx context.Context, accountID model.AccountID, from, to time.Time) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	if err := d.conn.SelectContext(ctx, &transactions, listTransactionByAccountIDQuery, accountID, from, to); err != nil {
		return nil, errors.Wrap(err, "Error while listing transactions by account id")
	}
	return transactions, nil
}

var listTransactionByUserID = `
	SELECT transaction_id, user_id, account_id, category_id, type, amount, currency, notes, date, created_at, deleted_at
	FROM transactions
	WHERE user_id = $1
		AND deleted_at IS NULL
		AND date > $2
		AND date < $3;
`

func (d *database) ListTransactionByUserID(ctx context.Context, userID model.UserID, from, to time.Time) ([]*model.Transaction, error) {
	var transactions []*model.Transaction
	if err := d.conn.SelectContext(ctx, &transactions, listTransactionByUserID, userID, from, to); err != nil {
		return nil, errors.Wrap(err, "Error: While getting transaction by user ID")
	}
	return transactions, nil
}

var updateTransactionByIDQuery = `
	UPDATE transactions
	SET account_id = :account_id,
		category_id = :category_id,
		date = :date,
		type = :type,
		amount = :amount,
		currency = :currency,
		notes = notes
	WHERE transaction_id = :transaction_id;
`

func (d *database) UpdateTransaction(ctx context.Context, transaction *model.Transaction) error {
	result, err := d.conn.NamedExecContext(ctx, updateTransactionByIDQuery, transaction)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Transaction not found")
	}
	return nil
}

var deleteTransactionByIDQuery = `
	UPDATE transactions
	SET deleted_at = NOW()
	WHERE transaction_id = $1;
`

func (d *database) DeleteTransaction(ctx context.Context, transactionID model.TransactionID) error {
	result, err := d.conn.ExecContext(ctx, deleteTransactionByIDQuery, transactionID)
	if err != nil {
		return errors.Wrap(err, "Error: while deleting transaction")
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.Wrap(err, "Error: Transaction not found")
	}
	return nil
}
