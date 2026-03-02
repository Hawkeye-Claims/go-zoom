// Package client provides a Zoom REST API client with support for OAuth 2.0
// authentication, automatic token refresh, and access to the Users, Meetings,
// and Phone service groups.
package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/TheSlowpes/go-zoom/zoom/tokenmutex"
	querystring "github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
)

const (
	zoomTokenURL = "https://zoom.us/oauth/token"
	zoomAuthUrl  = "https://zoom.us/oauth/authorize"
	zoomBaseURL  = "https://api.zoom.us/v2"
)

// Client is the top-level Zoom API client. It manages OAuth 2.0 tokens and
// exposes service groups (Users, Meetings, Phone) for interacting with the
// respective Zoom API resources.
type Client struct {
	httpClient   *http.Client
	accountID    string
	clientID     string
	clientSecret string
	grantType    string
	redirectURI  string
	oauthConf    *oauth2.Config
	stateMap     sync.Map
	tokenMutex   TokenMutex

	baseURL string

	// Users provides access to Zoom user management operations.
	Users *UsersService
	// Meetings provides access to Zoom meeting management operations.
	Meetings *MeetingsService
	// Phone provides access to Zoom Phone operations.
	Phone *PhoneService
}

// PaginationOptions holds cursor-based pagination parameters that can be
// appended to list requests to retrieve subsequent pages of results.
type PaginationOptions struct {
	// NextPageToken is the token returned by a prior list response, used to
	// fetch the next page of results.
	NextPageToken *string `url:"next_page_token,omitempty"`
	// PageSize is the maximum number of records to return per page.
	PageSize *int `url:"page_size,omitempty"`
}

// PaginationResponse contains the pagination metadata returned by Zoom list
// endpoints.
type PaginationResponse struct {
	// NextPageToken is the token to supply in the next request to retrieve the
	// following page. An empty string indicates the last page.
	NextPageToken string `json:"next_page_token"`
	// PageCount is the total number of pages available.
	PageCount int `json:"page_count"`
	// PageSize is the number of records returned in the current page.
	PageSize int `json:"page_size"`
	// TotalRecords is the total number of records across all pages.
	TotalRecords int `json:"total_records"`
}

// TokenMutex is the interface used by Client to store and retrieve OAuth
// access tokens in a concurrency-safe manner. Implement this interface to
// provide a custom token storage backend (e.g. Redis, a database, etc.).
// The default in-memory implementation is tokenmutex.Default.
type TokenMutex interface {
	// Lock acquires an exclusive lock before reading or writing token state.
	Lock(context.Context) error
	// Unlock releases the exclusive lock acquired by Lock.
	Unlock(context.Context) error
	// Get returns the current access token, or an error if no valid token
	// exists (e.g. tokenmutex.ErrTokenNotExist or tokenmutex.ErrTokenExpired).
	Get(context.Context) (string, error)
	// GetRefreshToken returns the stored OAuth refresh token.
	GetRefreshToken(context.Context) (string, error)
	// Set stores a new access token together with its expiry time.
	Set(context.Context, string, time.Time) error
	// SetRefreshToken stores a new refresh token.
	SetRefreshToken(context.Context, string) error
	// Clear removes any stored access token, forcing the next request to
	// re-authenticate.
	Clear(context.Context) error
}

// ClientOption is a functional option for configuring a Client at construction
// time via NewClient.
type ClientOption func(*Client)

// WithGrantType returns a ClientOption that overrides the OAuth 2.0 grant type
// used by the client. Supported values are "account_credentials" (default) and
// "authorization_code".
func WithGrantType(grantType string) ClientOption {
	return func(c *Client) {
		c.grantType = grantType
	}
}

// WithToken returns a ClientOption that sets a custom TokenMutex
// implementation on the client, replacing the default in-memory store.
func WithToken(tokenMutex TokenMutex) ClientOption {
	return func(c *Client) {
		c.tokenMutex = tokenMutex
	}
}

// WithRedirectURI returns a ClientOption that sets the OAuth 2.0 redirect URI.
// This is required when using the "authorization_code" grant type.
func WithRedirectURI(redirectURI string) ClientOption {
	return func(c *Client) {
		c.redirectURI = redirectURI
	}
}

// Compile-time assertion that tokenmutex.Default satisfies the TokenMutex
// interface.
var _ TokenMutex = (*tokenmutex.Default)(nil)

// NewClient creates and returns a new Zoom API Client. The httpClient is used
// for all outbound HTTP requests; pass nil to use http.DefaultClient behaviour
// (a non-nil *http.Client must be provided). accountID, clientID, and
// clientSecret are the Zoom OAuth application credentials. Additional
// behaviour can be customised via opts.
//
// Returns an error if the grant type is "authorization_code" and no redirect
// URI has been supplied via WithRedirectURI.
func NewClient(httpClient *http.Client, accountID, clientID, clientSecret string, opts ...ClientOption) (*Client, error) {
	var c = &Client{
		httpClient:   httpClient,
		accountID:    accountID,
		clientID:     clientID,
		clientSecret: clientSecret,
		grantType:    "account_credentials",
		baseURL:      zoomBaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.tokenMutex == nil {
		tokenMutex := tokenmutex.NewDefault()
		c.tokenMutex = tokenMutex
	}
	if c.grantType == "authorization_code" && len(c.redirectURI) == 0 {
		return nil, errors.New("redirect URI must be provided when using authorization_code grant type")
	}
	c.oauthConf = &oauth2.Config{
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		RedirectURL:  c.redirectURI,
		Endpoint: oauth2.Endpoint{
			AuthURL:  zoomAuthUrl,
			TokenURL: zoomTokenURL,
		},
	}
	c.Users = &UsersService{c}
	c.Meetings = &MeetingsService{c}

	return c, nil
}

// ErrorResponse represents an error body returned by the Zoom API. It
// implements the error interface, returning the human-readable Message field.
type ErrorResponse struct {
	// Code is the Zoom-specific numeric error code.
	Code int `json:"code"`
	// Message is a human-readable description of the error.
	Message string `json:"message"`
	// Errors contains field-level validation errors, if any.
	Errors []FieldError `json:"errors"`
}

// Error implements the error interface for ErrorResponse, returning the
// human-readable API error message.
func (e *ErrorResponse) Error() string {
	return e.Message
}

// FieldError describes a validation error for a specific request field
// returned by the Zoom API.
type FieldError struct {
	// Field is the name of the request field that failed validation.
	Field string `json:"field"`
	// Message describes why the field value was rejected.
	Message string `json:"message"`
}

// request is the internal method used by all service methods to send
// authenticated HTTP requests to the Zoom API. It acquires a valid access
// token (refreshing or fetching one as needed), encodes query parameters, and
// decodes the JSON response body into out. It returns the raw *http.Response
// alongside any error.
func (c *Client) request(ctx context.Context, method string, path string, query any, body any, out any) (*http.Response, error) {
	err := c.tokenMutex.Lock(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error locking token mutex: %w", err)
	}

	token, err := c.tokenMutex.Get(ctx)
	if err != nil {
		if !errors.Is(err, tokenmutex.ErrTokenNotExist) && !errors.Is(err, tokenmutex.ErrTokenExpired) {
			err = c.tokenMutex.Unlock(ctx)
			if err != nil {
				return nil, fmt.Errorf("Error unlocking token mutex: %w", err)
			}

			return nil, fmt.Errorf("Error getting token from mutex: %w", err)
		}

		var expiresAt time.Time
		switch c.grantType {
		case "authorization_code":
			refreshToken, err := c.tokenMutex.GetRefreshToken(ctx)
			if err != nil || len(refreshToken) == 0 {
				token, expiresAt, err = c.accessToken(ctx)
				if err != nil {
					err = c.tokenMutex.Unlock(ctx)
					if err != nil {
						return nil, fmt.Errorf("Error unlocking token mutex: %w", err)
					}
					return nil, fmt.Errorf("Error getting access token: %w", err)
				}
			} else {
				token, expiresAt, err = c.refreshToken(ctx, refreshToken)
				if err != nil {
					err = c.tokenMutex.Unlock(ctx)
					if err != nil {
						return nil, fmt.Errorf("Error unlocking token mutex: %w", err)
					}

					return nil, fmt.Errorf("Error refreshing access token: %w", err)
				}
			}
		case "account_credentials":
			token, expiresAt, err = c.accessToken(ctx)
			if err != nil {
				err = c.tokenMutex.Unlock(ctx)
				if err != nil {
					return nil, fmt.Errorf("Error unlocking token mutex: %w", err)
				}

				return nil, fmt.Errorf("Error getting access token: %w", err)
			}
		}

		err = c.tokenMutex.Set(ctx, token, expiresAt)
		if err != nil {
			err = c.tokenMutex.Unlock(ctx)
			if err != nil {
				return nil, fmt.Errorf("Error unlocking token mutex: %w", err)
			}

			return nil, fmt.Errorf("Error setting token in mutex: %w", err)
		}
	}

	err = c.tokenMutex.Unlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error unlocking token mutex: %w", err)
	}

	q, err := querystring.Values(query)
	if err != nil {
		return nil, fmt.Errorf("Error encoding query parameters: %w", err)
	}

	u := fmt.Sprintf("%s%s", c.baseURL, path)
	if len(q) > 0 {
		u = fmt.Sprintf("%s?%s", u, q.Encode())
	}

	reader := bytes.NewReader(nil)
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling request body: %w", err)
		}

		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, reader)
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error making HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusIMUsed {
		if res.StatusCode == http.StatusUnauthorized {
			err = c.tokenMutex.Clear(ctx)
			if err != nil {
				return nil, fmt.Errorf("Error clearing token mutex: %w", err)
			}
		}

		errRes := &ErrorResponse{}
		err = json.NewDecoder(res.Body).Decode(errRes)
		if err != nil {
			return res, fmt.Errorf("Error decoding error response body: %w", err)
		}

		return res, fmt.Errorf("Zoom API error (status %d): %v", res.StatusCode, errRes)
	}

	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return res, fmt.Errorf("Error decoding response body: %w", err)
		}
	}

	return res, nil
}

// authResponse is the JSON body returned by the Zoom token endpoint for both
// initial access-token requests and token refreshes.
type authResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// accessToken fetches a new access token from the Zoom token endpoint using
// the configured grant type. For "account_credentials" it uses Basic
// authentication with the client ID and secret. For "authorization_code" an
// error is returned because the initial exchange must go through the OAuth
// callback handler. The returned time.Time is the computed token expiry (five
// minutes before the server-reported TTL).
func (c *Client) accessToken(ctx context.Context) (string, time.Time, error) {
	query := url.Values{}
	query.Set("account_id", c.accountID)
	query.Set("grant_type", c.grantType)
	switch c.grantType {
	case "authorization_code":
		return "", time.Time{}, errors.New("authorization_code grant type requires authentication via the OAuth callback handler")
	case "account_credentials":

	default:
		return "", time.Time{}, fmt.Errorf("Unsupported grant type: %s", c.grantType)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s?%s", zoomTokenURL, query.Encode()), nil)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Error creating HTTP request: %w", err)
	}

	auth := base64.URLEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Error doing HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("Received non-200 status code: %d", res.StatusCode)
	}

	authRes := &authResponse{}
	err = json.NewDecoder(res.Body).Decode(authRes)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Error decoding response body: %w", err)
	}

	expiresAt := authRes.ExpiresIn - 300

	if authRes.RefreshToken != "" {
		if err = c.tokenMutex.SetRefreshToken(ctx, authRes.RefreshToken); err != nil {
			return "", time.Time{}, fmt.Errorf("Error setting refresh token in mutex: %w", err)
		}
	}

	return authRes.AccessToken, time.Now().Add(time.Duration(expiresAt) * time.Second), nil
}

// refreshToken exchanges a refresh token for a new access token using the
// Zoom token endpoint. It stores any new refresh token returned in the
// response. The returned time.Time is the computed token expiry (five minutes
// before the server-reported TTL).
func (c *Client) refreshToken(ctx context.Context, refreshToken string) (string, time.Time, error) {
	query := url.Values{}
	query.Set("grant_type", "refresh_token")
	query.Set("refresh_token", refreshToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s?%s", zoomTokenURL, query.Encode()), nil)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Error creating HTTP request: %w", err)
	}

	auth := base64.URLEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Error doing HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("Received non-200 status code: %d", res.StatusCode)
	}

	authRes := &authResponse{}
	err = json.NewDecoder(res.Body).Decode(authRes)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("Error decoding response body: %w", err)
	}

	expiresAt := authRes.ExpiresIn - 300

	if authRes.RefreshToken != "" {
		if err = c.tokenMutex.SetRefreshToken(ctx, authRes.RefreshToken); err != nil {
			return "", time.Time{}, fmt.Errorf("Error setting refresh token in mutex: %w", err)
		}
	}

	return authRes.AccessToken, time.Now().Add(time.Duration(expiresAt) * time.Second), nil
}

// RequestAuthorization returns an http.Handler that initiates the OAuth 2.0
// authorization code flow by redirecting the user to the Zoom authorization
// page. A random CSRF state value is generated, stored internally, and
// included in the redirect URL. Use HandleOAuthCallback to complete the flow.
func (c *Client) RequestAuthorization() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		state := rand.Text()
		c.stateMap.Store(state, struct{}{})
		url := c.oauthConf.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

// HandleOAuthCallback returns an http.Handler that completes the OAuth 2.0
// authorization code flow. It validates the CSRF state parameter, exchanges
// the authorization code for tokens, and stores the resulting access token in
// the TokenMutex. Mount this handler at the redirect URI registered with your
// Zoom OAuth application.
func (c *Client) HandleOAuthCallback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := c.tokenMutex.Lock(ctx)
		if err != nil {
			http.Error(w, "Error locking token mutex", http.StatusInternalServerError)
			return
		}
		defer c.tokenMutex.Unlock(ctx)

		state := r.FormValue("state")
		if _, ok := c.stateMap.LoadAndDelete(state); !ok {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		code := r.FormValue("code")
		token, err := c.oauthConf.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "oauthConf.Exchange() failed", http.StatusInternalServerError)
			return
		}
		expiresAt := token.Expiry.Add(-5 * time.Minute)
		err = c.tokenMutex.Set(ctx, token.AccessToken, expiresAt)
		if err != nil {
			http.Error(w, "Error setting token in mutex", http.StatusInternalServerError)
			return
		}
	})
}
