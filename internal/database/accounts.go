package database

import (
	"context"
	"financial-app/internal/model"

	"github.com/pkg/errors"
)

type AccountDB interface {
	CreateAccount(ctx context.Context, account *model.Account) error
	ListAccountsByUserID(ctx context.Context, userID model.UserID) ([]*model.Account, error)
	GetAccountByID(ctx context.Context, accountID model.AccountID) (*model.Account, error)
	UpdateAccountByID(ctx context.Context, account *model.Account) error
	DeleteAccountByID(ctx context.Context, accountID model.AccountID) (bool, error)
}

const createAccountQuery = `
	INSERT INTO accounts (user_id, start_balance, account_type, account_name, currency)
		VALUES( :user_id, :start_balance, :account_type, :account_name, :currency)
	RETURNING account_id;
`

func (d *database) CreateAccount(ctx context.Context, account *model.Account) error {
	rows, err := d.conn.NamedQueryContext(ctx, createAccountQuery, account)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&account.ID); err != nil {
		return err
	}

	return nil
}

var listAccountsByUserIDQuery = `
	SELECT account_id, user_id, account_name, start_balance, currency, created_at, deleted_at
	FROM accounts
	WHERE user_id = $1 AND deleted_at IS NULL;
`

func (d *database) ListAccountsByUserID(ctx context.Context, userID model.UserID) ([]*model.Account, error) {
	var accountList []*model.Account
	if err := d.conn.SelectContext(ctx, &accountList, listAccountsByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "Could not get account lists")
	}
	return accountList, nil
}

var getAccountByID = `
	SELECT account_id, user_id, start_balance, account_type, currency, created_at, deleted_at
		FROM accounts
	WHERE account_id = $1;
`

func (d *database) GetAccountByID(ctx context.Context, accountID model.AccountID) (*model.Account, error) {
	var account model.Account
	if err := d.conn.GetContext(ctx, &account, getAccountByID, accountID); err != nil {
		return nil, errors.Wrap(err, "Could not get account")
	}

	return &account, nil
}

var updateAccountByID = `
	UPDATE accounts
	SET start_balance = :start_balance,
		account_type  = :account_type,
		account_name  = :account_name,
		currency 	  = :currency
	WHERE account_id  = :account_id;
`

func (d *database) UpdateAccountByID(ctx context.Context, account *model.Account) error {
	result, err := d.conn.NamedExecContext(ctx, updateAccountByID, account)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Account not found")
	}

	return nil
}

var deleteAccountByID = `
	UPDATE accounts
	SET deleted_at = NOW()
	WHERE account_id = $1;
`

func (d *database) DeleteAccountByID(ctx context.Context, accountID model.AccountID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteAccountByID, accountID)
	if err != nil {
		return false, err
	}

	row, err := result.RowsAffected()
	if err != nil || row == 0 {
		return false, err
	}
	return true, nil
}
