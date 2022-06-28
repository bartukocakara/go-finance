package model

import "fmt"

type Credentials struct {
	SessionData
	Email    string `json:"email"`
	Password string `json:"password"`

	// TODO : We  will add google and facebook login
}

type Principal struct {
	UserID UserID `json:"UserID,omitempty"`
}

var NilPrincipal Principal

func (p Principal) String() string {
	if p.UserID != "" {
		return fmt.Sprintf("UserID[%s]", p.UserID)
	}
	return "(none)"
}
