package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/shmoulana/xeroapi/internal/api/handlers"
	"github.com/shmoulana/xeroapi/internal/config"
	"golang.org/x/oauth2"
)

const (
	EnvXeroClientID     = "XERO_CLIENT_ID"
	EnvXeroClientSecret = "XERO_CLIENT_SECRET"
	EnvXeroAuthUrl      = "XERO_AUTH_URL"
	EnvXeroTokenUrl     = "XERO_TOKEN_URL"
	EnvXeroCallbackAPI  = "XERO_CALLBACK_API"
)

func NewApp(sig chan os.Signal) *http.Server {
	oAuthScopes := []string{
		"openid",
		"profile",
		"email",
		"accounting.transactions",
		"accounting.settings",
		"offline_access",
	}
	xeroOauth2Config := &oauth2.Config{
		ClientID:     os.Getenv(EnvXeroClientID),
		ClientSecret: os.Getenv(EnvXeroClientSecret),
		Scopes:       oAuthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  os.Getenv(EnvXeroAuthUrl),
			TokenURL: os.Getenv(EnvXeroTokenUrl),
		},
		RedirectURL: fmt.Sprintf("http://localhost:%s%s", os.Getenv(config.EnvXeroAPISvrPort), os.Getenv(EnvXeroCallbackAPI)),
	}
	return &http.Server{
		Addr:    ":" + os.Getenv(config.EnvXeroAPISvrPort),
		Handler: NewRouter(handlers.NewHandler(xeroOauth2Config)),
	}
}
