package repository

import (
	"context"
	"fmt"
	"time"

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

func (repo *postgresUserManagementRepository) GetUserAndEmailVerification(ctx context.Context, email string) (*model.User, *model.EmailVerification, error) {
	var (
		user model.User
		ev   model.EmailVerification
	)

	err := repo.db.QueryRow(
		ctx,
		`SELECT users.id, users.email_verified_at,
				ev.id, ev.verification_code_hash,
				ev.strategy, ev.attempt_count,
				ev.expires_at
		FROM email_verifications AS ev
		JOIN users ON users.id = ev.user_id
		WHERE
			users.email = $1 AND
			ev.purpose = 'signup' AND
			ev.expires_at > now()
			ORDER BY ev.created_at DESC
		LIMIT 1`,
		email,
	).Scan(
		&user.ID, &user.EmailVerifiedAt,
		&ev.ID, &ev.VerificationCode,
		&ev.Strategy, &ev.AttemptCount,
		&ev.ExpiresAt,
	)

	if err != nil {
		return nil, nil, err
	}

	return &user, &ev, nil
}

func (repo *postgresUserManagementRepository) VerifyUserEmail(ctx context.Context, userID, evID string, session *model.UserSession) error {
	transaction, err := repo.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("transaction Begin: %v", err)
	}
	defer transaction.Rollback(ctx)

	// Set the user verified
	_, err = transaction.Exec(
		ctx,
		"UPDATE users SET email_verified_at = $1 WHERE id = $2",
		time.Now().UTC(), userID,
	)
	if err != nil {
		return fmt.Errorf("UPDATE users: %v", err)
	}

	// Delete email verification item
	_, err = transaction.Exec(
		ctx,
		"DELETE FROM email_verifications WHERE id = $1",
		evID,
	)
	if err != nil {
		return fmt.Errorf("DELETE FROM email_verifications: %v", err)
	}

	// Create user session
	_, err = transaction.Exec(
		ctx,
		`INSERT INTO user_sessions ( id, user_id, created_at, updated_at, expires_at, ip_address, user_agent )
		VALUES ( $1, $2, $3, $4, $5, $6, $7 )`,
		session.ID, session.UserID, session.CreatedAt, session.UpdatedAt, session.ExpiresAt, session.IPAddress, session.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("INSERT INTO user_sessions: %v", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return fmt.Errorf("transaction commit: %v", err)
	}

	return nil
}

func (repo *postgresUserManagementRepository) IncrementAttemptCount(ctx context.Context, evID string) error {
	_, err := repo.db.Exec(
		context.Background(),
		`UPDATE email_verifications SET attempt_count = attempt_count + 1 WHERE id = $1`,
		evID,
	)
	return err
}

func (repo *postgresUserManagementRepository) DeleteSession(ctx context.Context, sessionId string) error {
	_, err := repo.db.Exec(ctx, "DELETE FROM user_sessions WHERE id = $1", sessionId)
	return err
}
