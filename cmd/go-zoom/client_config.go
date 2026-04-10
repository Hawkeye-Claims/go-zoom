package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
)

// clientFlags holds shared auth/client options used by all command groups.
type clientFlags struct {
	accountID    string
	clientID     string
	clientSecret string
	grantType    string
	redirectURI  string
	timeout      time.Duration
}

// bindClientFlags registers common client flags and returns the target config.
func (c *cli) bindClientFlags(fs *flag.FlagSet, includeGrantType bool, includeRedirectURI bool) *clientFlags {
	cfg := &clientFlags{}

	fs.StringVar(&cfg.accountID, "account-id", os.Getenv(envZoomAccountID), "Zoom account ID (or set ZOOM_ACCOUNT_ID)")
	fs.StringVar(&cfg.clientID, "client-id", os.Getenv(envZoomClientID), "Zoom client ID (or set ZOOM_CLIENT_ID)")
	fs.StringVar(&cfg.clientSecret, "client-secret", os.Getenv(envZoomClientSecret), "Zoom client secret (or set ZOOM_CLIENT_SECRET)")
	fs.DurationVar(&cfg.timeout, "http-timeout", 30*time.Second, "HTTP client timeout")

	if includeGrantType {
		fs.StringVar(&cfg.grantType, "grant-type", "account_credentials", "OAuth grant type: account_credentials or authorization_code")
	}
	if includeRedirectURI {
		fs.StringVar(&cfg.redirectURI, "redirect-uri", "", "OAuth redirect URI (required for authorization_code)")
	}

	return cfg
}

// newServiceClient builds a client for non-interactive service commands.
func (c *cli) newServiceClient(cfg *clientFlags) (*client.Client, error) {
	zoomClient, err := c.newClient(cfg)
	if err != nil {
		return nil, err
	}
	if cfg.grantType == "authorization_code" {
		return nil, errors.New("service commands support grant-type=account_credentials; use auth authorization-code for interactive OAuth")
	}
	return zoomClient, nil
}

// newClient constructs an SDK client from validated CLI flags.
func (c *cli) newClient(cfg *clientFlags) (*client.Client, error) {
	if strings.TrimSpace(cfg.accountID) == "" {
		return nil, fmt.Errorf("missing required account ID (use --account-id or %s)", envZoomAccountID)
	}
	if strings.TrimSpace(cfg.clientID) == "" {
		return nil, fmt.Errorf("missing required client ID (use --client-id or %s)", envZoomClientID)
	}
	if strings.TrimSpace(cfg.clientSecret) == "" {
		return nil, fmt.Errorf("missing required client secret (use --client-secret or %s)", envZoomClientSecret)
	}

	grantType := strings.TrimSpace(cfg.grantType)
	if grantType == "" {
		grantType = "account_credentials"
	}
	if grantType != "account_credentials" && grantType != "authorization_code" {
		return nil, fmt.Errorf("unsupported grant type %q", grantType)
	}

	opts := []client.ClientOption{client.WithGrantType(grantType)}
	if grantType == "authorization_code" {
		if strings.TrimSpace(cfg.redirectURI) == "" {
			return nil, errors.New("redirect URI is required when using authorization_code grant type")
		}
		opts = append(opts, client.WithRedirectURI(cfg.redirectURI))
	}

	zoomClient, err := client.NewClient(
		&http.Client{Timeout: cfg.timeout},
		cfg.accountID,
		cfg.clientID,
		cfg.clientSecret,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return zoomClient, nil
}
