package model

import (
	"errors"
	"time"
)

type CategoryID string

var NilCategoryID CategoryID

type Category struct {
	ID        CategoryID `json:"id,omitempty" db:"category_id"`
	ParentID  CategoryID `json:"parentID,omitempty" db:"parent_id"`
	UserID    *UserID    `json:"userID,omitempty" db:"user_id"`
	CreatedAt *time.Time `json:"createdAt,omitempty" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt,omitEmpty" db:"deleted_at"`
	Name      *string    `json:"name,omitempty" db:"name"`
}

func (a *Category) Verify() error {
	if a.UserID == nil || len(*a.UserID) == 0 {
		return errors.New("UserID is required")
	}

	if a.Name == nil || len(*a.Name) == 0 {
		return errors.New("Name is required")
	}
	return nil
}
func (a *Category) SetParentID(parentID CategoryID) CategoryID {

	if parentID != NilCategoryID {
		a.ParentID = parentID
	}
	return a.ParentID
}

func (a *Category) SetName(name *string) *string {
	if name != nil {
		a.Name = name
	}
	return a.Name
}
