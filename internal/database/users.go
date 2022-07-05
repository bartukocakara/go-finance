package database

import (
	"context"
	"financial-app/internal/model"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type UsersDB interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, userID model.UserID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

var ErrUserExists = errors.New("User with that email already exists")

const createUserQuery = `
	INSERT INTO users (
		email, password
	)
	VALUES (
		:email, :password
	)
	RETURNING user_id
`

func (d *database) CreateUser(ctx context.Context, user *model.User) error {
	rows, err := d.conn.NamedQueryContext(ctx, createUserQuery, user)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == UniqueViolation {
				if pqError.Constraint == "user_email" {
					return ErrUserExists
				}
			}
		}
		return errors.Wrap(err, "Could not create user")
	}

	rows.Next()
	if err := rows.Scan(&user.ID); err != nil {
		return errors.Wrap(err, "Could not get created user id")
	}
	return nil
}

const getUserByIDQuery = `
	SELECT user_id, email, password, created_at
	FROM users
	WHERE user_id = $1 AND deleted_at IS NULL;
`

func (d *database) GetUserByID(ctx context.Context, userID model.UserID) (*model.User, error) {
	var user model.User
	if err := d.conn.GetContext(ctx, &user, getUserByIDQuery, userID); err != nil {
		return nil, err
	}
	return &user, nil
}

const getUserByEmailQuery = `
	SELECT user_id, email, password, created_at
	FROM users
	WHERE email = $1 AND deleted_at IS NULL;
`

func (d *database) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := d.conn.GetContext(ctx, &user, getUserByEmailQuery, email); err != nil {
		return nil, err
	}
	return &user, nil
}
