package auth

import "financial-app/internal/model"

// We will have 3 permission type for now
type PermissionType string

const (
	// User has 'admin' role
	Admin PermissionType = "admin"

	// User is logged in (we have user id in principal)
	Member PermissionType = "member"

	// User is logged in and user id passed to API is the same
	MemberIsTarget PermissionType = "memberIsTarget"

	Any PermissionType = "anonym"
)

// we have to create function for each type

// Admin
var adminOnly = func(roles []*model.UserRole) bool {
	for _, role := range roles {
		switch role.Role {
		case model.RoleAdmin:
			return true
		}
	}
	return false
}

// Logged in  user(TokenUserID)
var member = func(principal model.Principal) bool {
	return principal.UserID != ""
}

var memberIsTarget = func(userID model.UserID, principal model.Principal) bool {
	if userID == "" || principal.UserID == "" {
		return false
	}

	if userID != principal.UserID {
		return false
	}

	return true
}

var any = func() bool {
	return true
}
