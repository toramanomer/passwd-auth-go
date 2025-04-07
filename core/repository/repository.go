package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toramanomer/passwd-auth-go/core/model"
)

type UserManagementRepository interface {
	// CreateUser creates a new user with the given email and password
	// along with email verification.
	CreatePendingUser(context.Context, *model.User, *model.EmailVerification) error

	GetUserByEmail(ctx context.Context, email string) (*model.User, error)

	GetUserAndEmailVerification(ctx context.Context, email string) (*model.User, *model.EmailVerification, error)

	CreateEmailVerification(ctx context.Context, ev *model.EmailVerification) error

	VerifyUserEmail(ctx context.Context, userID, evID string, session *model.UserSession) error

	IncrementAttemptCount(ctx context.Context, evID string) error

	DeleteSession(ctx context.Context, sessionId string) error
}

func NewUserManagementRepository(db *pgxpool.Pool) UserManagementRepository {
	return &postgresUserManagementRepository{db}
}
