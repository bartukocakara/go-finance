package database

import (
	"context"
	"financial-app/internal/model"

	"github.com/pkg/errors"
)

type UserRoleDB interface {
	GrantRole(ctx context.Context, UserID model.UserID, role model.Role) error
	GetRolesByUser(ctx context.Context, userID model.UserID) ([]*model.UserRole, error)
}

const grantUserRoleQuery = `
	INSERT INTO user_roles (user_id, role)
		VALUES ($1, $2);
`

func (d *database) GrantRole(ctx context.Context, userID model.UserID, role model.Role) error {
	if _, err := d.conn.ExecContext(ctx, grantUserRoleQuery, userID, role); err != nil {
		return errors.Wrap(err, "Could not grant user role")
	}

	return nil
}

const getRolesByUserIDQuery = `
	SELECT role
	FROM user_roles
	WHERE user_id = $1;
`

func (d *database) GetRolesByUser(ctx context.Context, userID model.UserID) ([]*model.UserRole, error) {
	var roles []*model.UserRole
	if err := d.conn.SelectContext(ctx, &roles, getRolesByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "Could not get user roles")
	}
	return roles, nil
}
