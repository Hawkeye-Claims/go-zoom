package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Hawkeye-Claims/go-zoom/zoom/enums"
	"github.com/Hawkeye-Claims/go-zoom/zoom/models"
)

// PhoneUsersServicer is the interface implemented by PhoneUsersService. It
// declares the operations available for managing Zoom Phone user profiles.
type PhoneUsersServicer interface {
	// Get retrieves Zoom Phone user profiles. Provide WithPhoneUserID to fetch
	// a single user by ID, or no option to list all phone users (paginated
	// automatically).
	Get(ctx context.Context, opts ...PhoneUserGetOptions) ([]*models.PhoneUser, *http.Response, error)
	// GetProfileSetting retrieves the phone profile settings for the user
	// identified by userId.
	GetProfileSetting(ctx context.Context, userId string) (*models.PhoneUserSettings, *http.Response, error)
}

// PhoneUsersService implements PhoneUsersServicer and provides access to Zoom
// Phone user profile API endpoints.
type PhoneUsersService struct {
	client *Client
}

// Compile-time assertion that PhoneUsersService satisfies the
// PhoneUsersServicer interface.
var _ PhoneUsersServicer = (*PhoneUsersService)(nil)

// PhoneUserGetOptions is a functional option for configuring a phone user Get
// request.
type PhoneUserGetOptions func(*phoneUserGetOptions)

// phoneUserGetOptions holds the resolved configuration for a phone user Get
// call.
type phoneUserGetOptions struct {
	userID          string
	queryParameters *PhoneUserQueryParameters
}

// PhoneUserQueryParameters holds optional query parameters for listing Zoom
// Phone users.
type PhoneUserQueryParameters struct {
	// SiteID filters results to users belonging to the specified site.
	SiteID string `url:"site_id,omitempty"`
	// CallingType filters results by the user's calling plan type.
	CallingType int `url:"calling_type,omitempty"`
	// Status filters results by user status (e.g. active or inactive).
	Status enums.UserStatus `url:"status,omitempty"`
	// Department filters results to users in the specified department.
	Department string `url:"department,omitempty"`
	// CostCenter filters results to users associated with the specified cost
	// center.
	CostCenter string `url:"cost_center,omitempty"`
	// Keyword filters results by a search keyword (e.g. name or extension
	// number).
	Keyword string `url:"keyword,omitempty"`
}

// WithPhoneUserID returns a PhoneUserGetOptions that fetches the Zoom Phone
// profile for the user identified by userId.
func WithPhoneUserID(userId string) PhoneUserGetOptions {
	return func(o *phoneUserGetOptions) {
		o.userID = userId
	}
}

// WithPhoneUserQueryParameters returns a PhoneUserGetOptions that attaches the
// given query parameters to a phone user list request.
func WithPhoneUserQueryParameters(queryParameters *PhoneUserQueryParameters) PhoneUserGetOptions {
	return func(o *phoneUserGetOptions) {
		o.queryParameters = queryParameters
	}
}

// Get retrieves Zoom Phone user profiles. Provide WithPhoneUserID to fetch a
// single user by their ID, or no option to list all phone users (pagination is
// followed automatically). It is an error to specify both a user ID and query
// parameters simultaneously.
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

// GetProfileSetting retrieves the Zoom Phone profile settings for the user
// identified by userId.
func (u *PhoneUsersService) GetProfileSetting(ctx context.Context, userId string) (*models.PhoneUserSettings, *http.Response, error) {
	phoneUserSettings := &models.PhoneUserSettings{}
	res, err := u.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/users/%s/settings", url.PathEscape(userId)), nil, nil, phoneUserSettings)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return phoneUserSettings, res, nil
}
