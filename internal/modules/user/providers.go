package user

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/twitterv2"
)

func InitSessionStore() {
	secret := os.Getenv("SESSION_SECRET")
	store := sessions.NewCookieStore([]byte(secret))
	store.Options = &sessions.Options{
		Path:     "/",
		Domain:   "localhost",
		MaxAge:   86400 * 30,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	gothic.Store = store
}

func InitOAuthProviders() {
	// LinkedIn (OIDC)
	linkedinWellKnown := "https://www.linkedin.com/oauth/.well-known/openid-configuration"
	linkedinProvider, err := openidConnect.New(
		os.Getenv("LINKEDIN_CLIENT_ID"),
		os.Getenv("LINKEDIN_CLIENT_SECRET"),
		"http://localhost:8080/v1/auth/social/linkedin/callback",
		linkedinWellKnown,
		"openid",
		"profile",
		"email",
		"w_member_social",
	)
	if err != nil {
		log.Fatal("Failed to initialize LinkedIn OIDC:", err)
	}
	linkedinProvider.SetName("linkedin")

	// ---------- Twitter ----------
	twitterProvider := twitterv2.New(
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_SECRET_KEY"),
		"http://localhost:8080/v1/auth/social/twitter/callback",
	)
	twitterProvider.SetName("twitter")

	// ---------- Facebook ----------
	facebookProvider := facebook.New(
		os.Getenv("FACEBOOK_APP_ID"),
		os.Getenv("FACEBOOK_APP_SECRET"),
		"https://bijugate-superdevilishly-malinda.ngrok-free.dev/v1/auth/social/facebook/callback",
		"pages_show_list",
		"public_profile",
		"email",
		"pages_manage_posts",
		"pages_read_engagement",
	)
	facebookProvider.SetName("facebook")

	// ---------- Instagram ----------
	instagramProvider := facebook.New(
		os.Getenv("FACEBOOK_APP_ID"),
		os.Getenv("FACEBOOK_APP_SECRET"),
		"https://bijugate-superdevilishly-malinda.ngrok-free.dev/v1/auth/social/instagram/callback",
		"instagram_basic",
		"instagram_manage_messages",
		"instagram_manage_comments",
		"instagram_content_publish",
		"instagram_manage_insights",
	)
	instagramProvider.SetName("instagram")

	// ---------- Google ----------
	googleProvider := google.New(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		"http://localhost:8080/v1/auth/social/google/callback",
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/business.manage",
	)
	googleProvider.SetPrompt("consent", "select_account")
	googleProvider.SetAccessType("offline")
	googleProvider.SetName("google")

	// ---------- Register all ----------
	goth.UseProviders(
		facebookProvider,
		instagramProvider,
		twitterProvider,
		linkedinProvider,
		googleProvider,
	)
}
