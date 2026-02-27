package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/TheSlowpes/go-zoom/zoom/enums"
	"github.com/TheSlowpes/go-zoom/zoom/models"
)

type PhoneUsersServicer interface {
	Get(ctx context.Context, opts ...PhoneUserGetOptions) ([]*models.PhoneUser, *http.Response, error)
}

type PhoneUsersService struct {
	client *Client
}

var _ PhoneUsersServicer = (*PhoneUsersService)(nil)

type PhoneUserGetOptions func(*phoneUserGetOptions)

type phoneUserGetOptions struct {
	userID          string
	queryParameters *PhoneUserQueryParameters
}

type PhoneUserQueryParameters struct {
	SiteID      string           `url:"site_id,omitempty"`
	CallingType int              `url:"calling_type,omitempty"`
	Status      enums.UserStatus `url:"status,omitempty"`
	Department  string           `url:"department,omitempty"`
	CostCenter  string           `url:"cost_center,omitempty"`
	Keyword     string           `url:"keyword,omitempty"`
}

func WithPhoneUserID(userId string) PhoneUserGetOptions {
	return func(o *phoneUserGetOptions) {
		o.userID = userId
	}
}

func WithPhoneUserQueryParameters(queryParameters *PhoneUserQueryParameters) PhoneUserGetOptions {
	return func(o *phoneUserGetOptions) {
		o.queryParameters = queryParameters
	}
}

func (u *PhoneUsersService) Get(ctx context.Context, opts ...PhoneUserGetOptions) ([]*models.PhoneUser, *http.Response, error) {
	options := phoneUserGetOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.userID != "" && options.queryParameters != nil {
		return nil, nil, fmt.Errorf("cannot specify both userID and query parameters")
	}

	if options.userID != "" {
		endpoint := fmt.Sprintf("/phone/users/%s", url.PathEscape(options.userID))
		phoneUser := &models.PhoneUser{}
		res, err := u.client.request(ctx, http.MethodGet, endpoint, nil, nil, phoneUser)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		return []*models.PhoneUser{phoneUser}, res, nil
	}

	type response struct {
		*PaginationResponse
		PhoneUsers []*models.PhoneUser `json:"users"`
	}

	queryResponse := &response{}
	var phoneUsers []*models.PhoneUser

	endpoint := "/phone/users"

	res, err := u.client.request(ctx, http.MethodGet, endpoint, options.queryParameters, nil, queryResponse)

	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	phoneUsers = append(phoneUsers, queryResponse.PhoneUsers...)

	type phoneUserPageQuery struct {
		*PhoneUserQueryParameters
		*PaginationOptions
	}
	for {
		if queryResponse.NextPageToken == "" {
			break
		}
		nextPageToken := queryResponse.NextPageToken
		pageQuery := &phoneUserPageQuery{
			PhoneUserQueryParameters: options.queryParameters,
			PaginationOptions: &PaginationOptions{
				NextPageToken: &nextPageToken,
			},
		}
		res, err = u.client.request(ctx, http.MethodGet, endpoint, pageQuery, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		phoneUsers = append(phoneUsers, queryResponse.PhoneUsers...)
	}

	return phoneUsers, res, nil
}
