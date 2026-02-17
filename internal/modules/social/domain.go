package social

import "time"

type Platform string

const (
	LINKEDIN  Platform = "LINKEDIN"
	INSTAGRAM Platform = "INSTAGRAM"
	TWITTER   Platform = "TWITTER"
	FACEBOOK  Platform = "FACEBOOK"
	GOOGLE    Platform = "GOOGLE"
)

type Social struct {
	ID           string    `db:"id"`
	UserID       string    `db:"user_id"`
	Platform     Platform  `db:"platform"`
	ProviderID   string    `db:"provider_id"`
	AccessToken  string    `db:"access_token"`
	RefreshToken *string   `db:"refresh_token"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
