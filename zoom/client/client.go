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

type Client struct {
	httpClient   *http.Client
	accountID    string
	clientID     string
	clientSecret string
	grantType    string
	redirectURI  string
	code         string
	oauthConf    *oauth2.Config
	state        string
	tokenMutex   TokenMutex

	baseURL string

	Users    *UsersService
	Meetings *MeetingsService
}

type PaginationOptions struct {
	NextPageToken *string `url:"next_page_token,omitempty"`
	PageSize      *int    `url:"page_size,omitempty"`
}

type PaginationResponse struct {
	NextPageToken string `json:"next_page_token"`
	PageCount     int    `json:"page_count"`
	PageSize      int    `json:"page_size"`
	TotalRecords  int    `json:"total_records"`
}

type TokenMutex interface {
	Lock(context.Context) error
	Unlock(context.Context) error
	Get(context.Context) (string, error)
	GetRefreshToken(context.Context) (string, error)
	Set(context.Context, string, time.Time) error
	SetRefreshToken(context.Context, string) error
	Clear(context.Context) error
}

type ClientOption func(*Client)

func WithGrantType(grantType string) ClientOption {
	return func(c *Client) {
		c.grantType = grantType
	}
}

func WithToken(tokenMutex TokenMutex) ClientOption {
	return func(c *Client) {
		c.tokenMutex = tokenMutex
	}
}

func WithRedirectURI(redirectURI string) ClientOption {
	return func(c *Client) {
		c.redirectURI = redirectURI
	}
}

var _ TokenMutex = (*tokenmutex.Default)(nil)

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

type ErrorResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Errors  []FieldError `json:"errors"`
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

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

type authResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func (c *Client) accessToken(ctx context.Context) (string, time.Time, error) {
	query := url.Values{}
	query.Set("account_id", c.accountID)
	query.Set("grant_type", c.grantType)
	switch c.grantType {
	case "authorization_code":
		if len(c.code) == 0 {
			return "", time.Time{}, errors.New("authorization code must be retrieved before using authorization_code grant type")
		}
		query.Set("code", c.code)
		query.Set("redirect_uri", c.redirectURI)
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

func (c *Client) RequestAuthorization() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.state = rand.Text()
		url := c.oauthConf.AuthCodeURL(c.state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})
}

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
		if state != c.state {
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
