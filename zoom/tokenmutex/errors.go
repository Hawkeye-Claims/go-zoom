// Package tokenmutex provides the TokenMutex interface and a default in-memory
// implementation for thread-safe OAuth token storage used by the Zoom API client.
package tokenmutex

import "errors"

// ErrTokenNotExist is returned by Get or GetRefreshToken when no token has been
// stored yet.
var ErrTokenNotExist = errors.New("token does not exist")

// ErrTokenExpired is returned by Get when a token exists but its expiry time has
// passed. The caller should fetch or refresh the token and store the new value.
var ErrTokenExpired = errors.New("token expired")
