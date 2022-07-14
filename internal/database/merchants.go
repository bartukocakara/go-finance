package database

import (
	"context"
	"financial-app/internal/model"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type MerchantDB interface {
	CreateMerchant(ctx context.Context, merchant *model.Merchant) error
	ListMerchantsByUserID(ctx context.Context, userID model.UserID) ([]*model.Merchant, error)
	UpdateMerchant(ctx context.Context, merchant *model.Merchant) error
	GetMerchantByID(ctx context.Context, merchantID model.MerchantID) (*model.Merchant, error)
	DeleteMerchantByID(ctx context.Context, merchantID *model.MerchantID) error
}

var MerchantNameExists = errors.New("Merchant Name already using")

const createMerchantQuery = `
	INSERT INTO merchants (user_id, name)
		VALUES (:user_id, :name)
	RETURNING merchant_id;
`

func (d *database) CreateMerchant(ctx context.Context, merchant *model.Merchant) error {
	rows, err := d.conn.NamedQueryContext(ctx, createMerchantQuery, merchant)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == UniqueViolation {
				if pqError.Constraint == "merchant_name" {
					return MerchantNameExists
				}
			}
		}
		return errors.Wrap(err, "Could not create Merchant")
	}
	rows.Next()
	if err := rows.Scan(&merchant.ID); err != nil {
		return errors.Wrap(err, "Could not get created Merchant ID")
	}

	return nil
}

const listMerchantByUserIDQuery = `
	SELECT merchant_id, user_id, name, created_at, deleted_at 
	FROM merchants   
	WHERE user_id = $1 AND deleted_at IS NULL;
`

func (d *database) ListMerchantsByUserID(ctx context.Context, userID model.UserID) ([]*model.Merchant, error) {
	var merchants []*model.Merchant
	if err := d.conn.SelectContext(ctx, &merchants, listMerchantByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "Could not get users merchants")
	}
	return merchants, nil
}

const getMerchantByIDQuery = `
	SELECT merchant_id, user_id, name, created_at, deleted_at 
	FROM merchants   
	WHERE merchant_id = $1 AND deleted_at IS NULL;
`

func (d *database) GetMerchantByID(ctx context.Context, merchantID model.MerchantID) (*model.Merchant, error) {
	var merchant model.Merchant
	if err := d.conn.GetContext(ctx, &merchant, getMerchantByIDQuery, merchantID); err != nil {
		return nil, errors.Wrap(err, "could not get merchant")
	}

	return &merchant, nil
}

const updateMerchantQuery = `
	UPDATE merchants 
	SET name = :name 
	WHERE merchant_id = :merchant_id;
`

func (d *database) UpdateMerchant(ctx context.Context, merchant *model.Merchant) error {
	rows, err := d.conn.NamedExecContext(ctx, updateMerchantQuery, merchant)
	if rows != nil {
		defer rows.RowsAffected()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == UniqueViolation {
				if pqError.Constraint == "merchant_name" {
					return MerchantNameExists
				}
			}
		}
		return errors.Wrap(err, "Could not update Merchant")
	}

	return nil
}

const deleteMerchantByIDQuery = `
	DELETE FROM merchants
	WHERE merchant_id = $1
	AND deleted_at IS NULL;
`

func (d *database) DeleteMerchantByID(ctx context.Context, merchantID *model.MerchantID) error {
	if _, err := d.conn.ExecContext(ctx, deleteMerchantByIDQuery, merchantID); err != nil {
		return err
	}

	return nil
}
