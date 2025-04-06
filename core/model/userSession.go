package model

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type UserSession struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
	IPAddress string
	UserAgent string
}

func NewUserSession(userID, ip, userAgent string) *UserSession {
	sessionId := make([]byte, 64)
	rand.Read(sessionId)
	utcNow := time.Now().UTC()

	return &UserSession{
		ID:        base64.RawURLEncoding.EncodeToString(sessionId),
		UserID:    userID,
		CreatedAt: utcNow,
		UpdatedAt: utcNow,
		ExpiresAt: utcNow.Add(10 * time.Minute),
		IPAddress: ip,
		UserAgent: userAgent,
	}
}
