package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              string
	Email           string
	PasswordHash    string
	EmailVerifiedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func NewUser(email, passwordHash string) *User {
	utcNow := time.Now().UTC()
	return &User{
		ID:              uuid.NewString(),
		Email:           email,
		PasswordHash:    passwordHash,
		EmailVerifiedAt: nil,
		CreatedAt:       utcNow,
		UpdatedAt:       utcNow,
	}
}
