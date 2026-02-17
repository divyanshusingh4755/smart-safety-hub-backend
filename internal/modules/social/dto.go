package social

import "time"

type SocialAccountResponse struct {
	ID         string    `json:"id"`
	Platform   Platform  `json:"platform"`
	ProviderID string    `json:"provider_id"`
	ExpiresAt  time.Time `json:"expires_at"`
}

type SocialResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Platform     Platform  `json:"platform"`
	ProviderID   string    `json:"provider_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken *string   `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ListOfSocialConnectResponse struct {
	SocialResponse
	IsExpired bool `json:"is_expired"`
}

func ToResponse(s *Social) SocialAccountResponse {
	return SocialAccountResponse{
		ID:         s.ID,
		Platform:   s.Platform,
		ProviderID: s.ProviderID,
		ExpiresAt:  s.ExpiresAt,
	}
}

type CreatePostRequest struct {
	Platform    string `json:"platform"`
	PostContent string `json:"post_content"`
	ProductID   string `json:"product_id"`
}
