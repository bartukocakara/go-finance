package model

import "errors"

type DeviceID string

type Session struct {
	UserID       UserID   `db:"user_id"`
	DeviceID     DeviceID `db:"device_id"`
	RefreshToken string   `db:"refresh_token"`
	ExpiresAt    int64    `db:"expires_at"`
}

type SessionData struct {
	DeviceID DeviceID `db:"deviceID"`
}

func (s *SessionData) Verify() error {
	if len(s.DeviceID) == 0 {
		return errors.New("Device Id is required")
	}
	return nil
}
