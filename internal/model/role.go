package model

type Role string

const (
	RoleAdmin Role = "admin"
)

type UserRole struct {
	Role Role `json:"role" db:"role"`
}