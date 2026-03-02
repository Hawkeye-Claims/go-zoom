package tokenmutex

import (
	"context"
	"sync"
	"time"
)

// Default is an in-memory, mutex-backed implementation of the client.TokenMutex
// interface. It is used automatically when no custom token store is supplied to
// client.NewClient. It is safe for concurrent use by a single process but is not
// suitable for horizontally-scaled deployments — use a shared store (e.g. Redis)
// in that case.
type Default struct {
	token        string
	refreshToken string
	expiresAt    time.Time

	lock sync.Mutex
}

// NewDefault returns a new, empty Default token store.
func NewDefault() *Default {
	return &Default{}
}

// Lock acquires the underlying mutex so the caller can read and conditionally
// update the token atomically. The ctx parameter is accepted for interface
// compatibility but is not used.
func (d *Default) Lock(ctx context.Context) error {
	d.lock.Lock()
	return nil
}

// Unlock releases the underlying mutex. It must be called after every successful
// Lock call.
func (d *Default) Unlock(context.Context) error {
	d.lock.Unlock()
	return nil
}

// Get returns the stored access token. It returns ErrTokenNotExist when no token
// has been stored yet, and ErrTokenExpired when the stored token has passed its
// expiry time.
func (d *Default) Get(ctx context.Context) (string, error) {
	if len(d.token) == 0 {
		return "", ErrTokenNotExist
	}

	if time.Now().After(d.expiresAt) {
		return "", ErrTokenExpired
	}

	return d.token, nil
}

// Set stores the access token together with its expiry time, replacing any
// previously stored value.
func (d *Default) Set(ctx context.Context, token string, expiresAt time.Time) error {
	d.token = token
	d.expiresAt = expiresAt

	return nil
}

// GetRefreshToken returns the stored refresh token. It returns ErrTokenNotExist
// when no refresh token has been stored yet.
func (d *Default) GetRefreshToken(ctx context.Context) (string, error) {
	if len(d.refreshToken) == 0 {
		return "", ErrTokenNotExist
	}

	return d.refreshToken, nil
}

// SetRefreshToken stores a refresh token, replacing any previously stored value.
func (d *Default) SetRefreshToken(ctx context.Context, refreshToken string) error {
	d.refreshToken = refreshToken

	return nil
}

// Clear removes the stored access token and resets its expiry time. It is called
// automatically by the client when a 401 Unauthorized response is received, so
// the next request will fetch a fresh token.
func (d *Default) Clear(ctx context.Context) error {
	d.token = ""
	d.expiresAt = time.Time{}

	return nil
}
