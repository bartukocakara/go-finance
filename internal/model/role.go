package model

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type UserRole struct {
	Role Role `json:"role" db:"role"`
}
