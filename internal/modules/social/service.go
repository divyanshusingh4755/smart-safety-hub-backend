package social

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/markbates/goth"
	"go.uber.org/zap"
)

type SocialService struct {
	logger *zap.Logger
	repo   *Repository
}

func NewSocialService(logger *zap.Logger, repo *Repository) *SocialService {
	return &SocialService{
		logger: logger,
		repo:   repo,
	}
}

func (s *SocialService) LinkExternalAccount(ctx context.Context, UserID string, gothUser goth.User) (*SocialResponse, error) {
	expiryUTC := gothUser.ExpiresAt.UTC().Round(0)

	if gothUser.ExpiresAt.IsZero() {
		expiryUTC = time.Now().AddDate(100, 0, 0).UTC().Round(0)
	}

	newSocial := &SocialResponse{
		UserID:       UserID,
		Platform:     Platform(strings.ToUpper(gothUser.Provider)),
		ProviderID:   gothUser.UserID,
		AccessToken:  gothUser.AccessToken,
		RefreshToken: nil,
		ExpiresAt:    expiryUTC,
	}

	// Handles oAuth 2.0 Refresh Tokens
	if gothUser.RefreshToken != "" {
		newSocial.RefreshToken = &gothUser.RefreshToken
	}

	// Handles Twitter OAuth 1.0a "secret"
	if strings.ToLower(gothUser.Provider) == "twitter" && gothUser.AccessTokenSecret != "" {
		newSocial.RefreshToken = &gothUser.AccessTokenSecret
	}

	if err := s.repo.Upsert(ctx, newSocial); err != nil {
		return nil, err
	}
	return newSocial, nil
}

func (s *SocialService) HandleListConnections(ctx context.Context, userID string) ([]ListOfSocialConnectResponse, error) {
	response, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("Social Serivce: failed to get connections for user %s: %w", userID, err)
	}

	connections := make([]ListOfSocialConnectResponse, 0, len(response))
	now := time.Now().UTC().Truncate(time.Second)

	for _, resp := range response {
		expiry := resp.ExpiresAt.UTC().Truncate(time.Second)
		data := SocialResponse{
			ID:           resp.ID,
			UserID:       resp.UserID,
			Platform:     resp.Platform,
			ProviderID:   resp.ProviderID,
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    resp.ExpiresAt,
			CreatedAt:    resp.CreatedAt,
			UpdatedAt:    resp.UpdatedAt,
		}
		c := ListOfSocialConnectResponse{
			SocialResponse: data,
			IsExpired:      now.After(expiry),
		}
		connections = append(connections, c)
	}

	return connections, nil
}

func (s *SocialService) CreatePost(ctx context.Context, userID string, req CreatePostRequest) error {
	// Get token from DB
	account, err := s.repo.GetByPlatform(ctx, userID, strings.ToUpper(req.Platform))
	if err != nil {
		return fmt.Errorf("Social account not found: %w", err)
	}

	// Route to specific API handler
	switch strings.ToLower(req.Platform) {
	case "twitter":
		return s.postToTwitter(account.AccessToken, *account.RefreshToken, req.PostContent)
	case "linkedin":
		return s.postToLinkedIn(account.ProviderID, account.AccessToken, req.PostContent)
	default:
		return fmt.Errorf("Platform %s not supported for posting yet", req.Platform)
	}
}

func (s *SocialService) postToTwitter(accessToken string, accessSecret string, content string) error {
	// Setup the OAuth 1.0a Config
	config := oauth1.NewConfig(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_SECRET_KEY"),
	)

	// Setup the user's specific token and secret
	credentials := oauth1.NewToken(accessToken, accessSecret)

	// Create the authorized client
	httpClient := config.Client(oauth1.NoContext, credentials)

	// Endpoint for Twitter v2
	url := "https://api.twitter.com/2/tweets"

	// Prepare JSON Body
	payload := map[string]string{"text": content}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error while marshalling twitter json: %v", err)
	}

	// AUTHORIZED CLIENT TO POST
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error while sending data to twitter: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("twitter api error: %s - %s", resp.Status, string(respBody))
	}

	return nil
}

func (s *SocialService) postToLinkedIn(personID string, token string, content string) error {
	url := "https://api.linkedin.com/v2/ugcPosts"
	// LinkedIn required a specific JSON structure
	payload := map[string]interface{}{
		"author":         "urn:li:person:" + personID,
		"lifecycleState": "PUBLISHED",
		"specificContent": map[string]interface{}{
			"com.linkedin.ugc.ShareContent": map[string]interface{}{
				"shareCommentary": map[string]interface{}{
					"text": content,
				},
				"shareMediaCategory": "NONE",
			},
		},
		"visibility": map[string]interface{}{
			"com.linkedin.ugc.MemberNetworkVisibility": "PUBLIC",
		},
	}

	body, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("Error while marshelling linkedin json: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Error while send data to linkedin: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 201 {
		return fmt.Errorf("linkedin api error: %v", resp.Status)
	}
	return nil
}
