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
}

func NewUserManagementRepository(db *pgxpool.Pool) UserManagementRepository {
	return &postgresUserManagementRepository{db}
}
