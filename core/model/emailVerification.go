package model

import (
	"time"

	"github.com/google/uuid"
)

type EmailVerification struct {
	ID               string
	UserID           string
	VerificationCode string
	Strategy         string
	Purpose          string
	AttemptCount     int
	CreatedAt        time.Time
	ExpiresAt        time.Time
}

func NewEmailVerification(userId, verificationCode, strategy string) *EmailVerification {
	utcNow := time.Now().UTC()

	return &EmailVerification{
		ID:               uuid.NewString(),
		UserID:           userId,
		VerificationCode: verificationCode,
		Strategy:         strategy,
		Purpose:          "signup",
		AttemptCount:     0,
		CreatedAt:        utcNow,
		ExpiresAt:        utcNow.Add(5 * time.Minute),
	}
}
