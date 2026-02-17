package social

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/smart-safety-hub/backend/shared"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Upsert(ctx context.Context, s *SocialResponse) error {
	query := `INSERT INTO social_accounts (user_id, platform, provider_id, access_token, refresh_token, expires_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW())
	ON CONFLICT (user_id, platform)
	DO UPDATE SET
		access_token = EXCLUDED.access_token,
		refresh_token = EXCLUDED.refresh_token,
		expires_at = EXCLUDED.expires_at,
		updated_at = NOW()
	RETURNING id, created_at`

	return r.db.QueryRowContext(ctx, query, s.UserID, strings.ToUpper(string(s.Platform)), s.ProviderID, s.AccessToken, s.RefreshToken, s.ExpiresAt).Scan(&s.ID, &s.CreatedAt)
}

func (r *Repository) GetByUserID(ctx context.Context, userID string) ([]Social, error) {
	query := `SELECT id, user_id, platform, provider_id, access_token, refresh_token, expires_at, updated_at FROM social_accounts WHERE user_id = $1`
	var socials []Social
	if err := r.db.SelectContext(ctx, &socials, query, userID); err != nil {
		return nil, shared.PostgresError(err)
	}
	return socials, nil
}

func (r *Repository) GetByPlatform(ctx context.Context, userID, platform string) (*Social, error) {
	query := `SELECT id, user_id, platform, provider_id, access_token, refresh_token, expires_at, updated_at FROM social_accounts WHERE user_id = $1 AND platform = $2 LIMIT 1`
	var social Social
	err := r.db.QueryRowContext(ctx, query, userID, strings.ToUpper(platform)).Scan(
		&social.ID,
		&social.UserID,
		&social.Platform,
		&social.ProviderID,
		&social.AccessToken,
		&social.RefreshToken,
		&social.ExpiresAt,
		&social.UpdatedAt,
	)

	if err != nil {
		return nil, shared.PostgresError(err)
	}
	return &social, nil
}
