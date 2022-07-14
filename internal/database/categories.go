package database

import (
	"context"
	"financial-app/internal/model"

	"github.com/pkg/errors"
)

type CategoryDB interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, category *model.Category) error
	GetCategoryByID(ctx context.Context, categoryID model.CategoryID) (*model.Category, error)
	ListCategoriesByUserID(ctx context.Context, userID model.UserID) ([]*model.Category, error)
	DeleteCategoryByID(ctx context.Context, categoryID model.CategoryID) (bool, error)
}

var createCategoryQuery = `
	INSERT INTO categories (parent_id, user_id, name)
		VALUES (:parent_id, :user_id, :name)
	RETURNING category_id;
`

func (d *database) CreateCategory(ctx context.Context, category *model.Category) error {
	rows, err := d.conn.NamedQueryContext(ctx, createCategoryQuery, category)
	if err != nil {
		return err
	}

	defer rows.Close()
	rows.Next()
	if err := rows.Scan(&category.ID); err != nil {
		return err
	}

	return nil
}

var updateCategoryQuery = `
	UPDATE categories
	SET parent_id = :parent_id,
		name = :name
	WHERE category_id = :category_id;
`

func (d *database) UpdateCategory(ctx context.Context, category *model.Category) error {
	result, err := d.conn.NamedExecContext(ctx, updateCategoryQuery, category)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil || row == 0 {
		return errors.New("Category not found")
	}

	return nil
}

var getCategoryByIDQuery = `
	SELECT category_id, parent_id, user_id, name, created_at, deleted_at
	FROM categories
	WHERE category_id = $1
		AND deleted_at IS NULL;
`

func (d *database) GetCategoryByID(ctx context.Context, categoryID model.CategoryID) (*model.Category, error) {
	var category model.Category
	if err := d.conn.GetContext(ctx, &category, getCategoryByIDQuery, categoryID); err != nil {
		return nil, errors.Wrap(err, "Could not get category")
	}
	return &category, nil
}

var listCategoriesByUserIDQuery = `
	SELECT category_id, parent_id, user_id, name, created_at, deleted_at
	FROM categories
	WHERE user_id = $1 AND deleted_at IS NULL;
`

func (d *database) ListCategoriesByUserID(ctx context.Context, userID model.UserID) ([]*model.Category, error) {
	var categories []*model.Category
	if err := d.conn.SelectContext(ctx, &categories, listCategoriesByUserIDQuery, userID); err != nil {
		return nil, errors.Wrap(err, "Could not get users categories")
	}
	return categories, nil
}

var deleteCategoryByIDQuery = `
	DELETE FROM categories
	WHERE category_id = $1
		AND deleted_at IS NULL;
`

func (db *database) DeleteCategoryByID(ctx context.Context, categoryId model.CategoryID) (bool, error) {
	result, err := db.conn.ExecContext(ctx, deleteCategoryByIDQuery, categoryId)
	if err != nil {
		return false, err
	}

	row, err := result.RowsAffected()
	if err != nil || row == 0 {
		return false, err
	}

	return true, nil
}
