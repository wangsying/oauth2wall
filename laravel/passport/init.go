package passport

import (
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"golang.org/x/oauth2"
)

// Endpoint is passport's OAuth 2.0 endpoint.
var config Config
var oAuth2Config oauth2.Config
var oauthStateString string
var sessionManager *scs.SessionManager

func init() {
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	config = Config{
		AuthURL:      os.Getenv("LARAVEL_PASSPORT_ENDPOINT") + "/oauth/authorize",
		TokenURL:     os.Getenv("LARAVEL_PASSPORT_ENDPOINT") + "/oauth/token",
		UserURL:      os.Getenv("LARAVEL_PASSPORT_ENDPOINT") + "/api/user",
		CallbackURL:  os.Getenv("LARAVEL_PASSPORT_REDIRECT_URL"),
		ClientID:     os.Getenv("LARAVEL_PASSPORT_CLIENT_ID"),
		ClientSecret: os.Getenv("LARAVEL_PASSPORT_CLIENT_SECRET"),
		SessionKey:   os.Getenv("LARAVEL_PASSPORT_SESSION_KEY"),
	}

	oAuth2Config = oauth2.Config{
		RedirectURL:  config.CallbackURL,
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scopes:       []string{"*"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   config.AuthURL,
			TokenURL:  config.TokenURL,
			AuthStyle: oauth2.AuthStyleInHeader,
		},
	}
}
