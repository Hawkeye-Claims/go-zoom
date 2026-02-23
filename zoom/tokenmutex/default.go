package tokenmutex

import (
	"context"
	"sync"
	"time"
)

type Default struct {
	token        string
	refreshToken string
	expiresAt    time.Time

	lock sync.Mutex
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) Lock(ctx context.Context) error {
	d.lock.Lock()
	return nil
}

func (d *Default) Unlock(context.Context) error {
	d.lock.Unlock()
	return nil
}

func (d *Default) Get(ctx context.Context) (string, error) {
	if len(d.token) == 0 {
		return "", ErrTokenNotExist
	}

	if time.Now().After(d.expiresAt) {
		return "", ErrTokenExpired
	}

	return d.token, nil
}

func (d *Default) Set(ctx context.Context, token string, expiresAt time.Time) error {
	d.token = token
	d.expiresAt = expiresAt

	return nil
}

func (d *Default) GetRefreshToken(ctx context.Context) (string, error) {
	if len(d.refreshToken) == 0 {
		return "", ErrTokenNotExist
	}

	if time.Now().After(d.expiresAt) {
		return "", ErrTokenExpired
	}

	return d.refreshToken, nil
}

func (d *Default) SetRefreshToken(ctx context.Context, refreshToken string) error {
	d.refreshToken = refreshToken

	return nil
}

func (d *Default) Clear(ctx context.Context) error {
	d.token = ""
	d.expiresAt = time.Time{}

	return nil
}
