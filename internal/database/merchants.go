package database

import (
	"context"
	"financial-app/internal/model"
)

type MerchantDB interface {
	CreateMerchant(ctx context.Context, merchant *model.Merchant) error
}

const createMerchantQuery = `
	INSERT INTO merchants (user_id, name)
		VALUES (:user_id, name)
	RETURNING merchant_id;
`

func (d *database) CreateMerchant(ctx context.Context, merchant *model.Merchant) error {
	return nil
}
