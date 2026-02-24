package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/TheSlowpes/go-zoom/zoom/enums"
	"github.com/TheSlowpes/go-zoom/zoom/models"
)

type UsersServicers interface {
	Get(ctx context.Context, opts ...UserGetOptions) ([]*models.User, *http.Response, error)
}

type UserService struct {
	client *Client
}

var _ UsersServicers = (*UserService)(nil)

type UserGetOptions func(*usersGetOptions)

type usersGetOptions struct {
	userId              string `url:"userId,omitempty"`
	queryParameters     *UserQueryParameters
	listQueryParameters *ListUserQueryParameters
}

type UserQueryParameters struct {
	LoginType        enums.LoginType `url:"login_type,omitempty"`
	EncryptedEmail   bool            `url:"encrypted_email,omitempty"`
	SearchByUniqueID bool            `url:"search_by_unique_id,omitempty"`
}

type ListUserQueryParameters struct {
	Status        enums.UserStatus `url:"status,omitempty"`
	RoleID        string           `url:"role_id,omitempty"`
	IncludeFields enums.UserFields `url:"include_fields,omitempty"`
}

func WithUserId(userId string) UserGetOptions {
	return func(o *usersGetOptions) {
		o.userId = userId
	}
}

func WithUserQueryParameters(params *UserQueryParameters) UserGetOptions {
	return func(o *usersGetOptions) {
		o.queryParameters = params
	}
}

func WithListUserQueryParameters(params *ListUserQueryParameters) UserGetOptions {
	return func(o *usersGetOptions) {
		o.listQueryParameters = params
	}
}

func (u *UserService) Get(ctx context.Context, opts ...UserGetOptions) ([]*models.User, *http.Response, error) {
	options := usersGetOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.queryParameters != nil && options.listQueryParameters != nil {
		return nil, nil, errors.New("Cannot use both UserQueryParameters and ListUserQueryParameters at the same time")
	}

	query := interface{}(nil)
	if options.queryParameters != nil {
		query = options.queryParameters
	} else if options.listQueryParameters != nil {
		query = options.listQueryParameters
	}

	var users []*models.User

	type response struct {
		*PaginationResponse
		Users []*models.User `json:"users"`
	}

	queryResponse := &response{}

	res, err := u.client.request(ctx, http.MethodGet, "/users/", query, nil, queryResponse)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	users = append(users, queryResponse.Users...)

}
