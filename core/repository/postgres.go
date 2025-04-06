package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/toramanomer/passwd-auth-go/core/model"
)

type postgresUserManagementRepository struct {
	db *pgxpool.Pool
}

func (repo *postgresUserManagementRepository) CreatePendingUser(ctx context.Context, user *model.User, ev *model.EmailVerification) error {
	transaction, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer transaction.Rollback(ctx)

	_, err = transaction.Exec(
		ctx,
		`INSERT INTO users
			( id, email, password_hash, email_verified_at, created_at, updated_at )
		VALUES
			( $1, $2, $3, $4, $5, $6 )`,
		user.ID, user.Email, user.PasswordHash, user.EmailVerifiedAt, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	_, err = transaction.Exec(
		ctx,
		`INSERT INTO email_verifications
			( id, user_id, verification_code_hash, strategy, purpose, attempt_count, created_at, expires_at )
		VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 )
		`,
		ev.ID,
		ev.UserID,
		ev.VerificationCode,
		ev.Strategy,
		ev.Purpose,
		ev.AttemptCount,
		ev.CreatedAt,
		ev.ExpiresAt,
	)
	if err != nil {
		return err
	}

	return transaction.Commit(ctx)
}
