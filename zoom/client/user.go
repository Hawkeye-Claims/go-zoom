package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/TheSlowpes/go-zoom/zoom/enums"
	"github.com/TheSlowpes/go-zoom/zoom/models"
)

type UsersServicers interface {
	Get(ctx context.Context, opts ...UserGetOptions) ([]*models.User, *http.Response, error)
	Create(ctx context.Context, action enums.UserCreateAction, userAttributes UserAttributes) (*models.User, *http.Response, error)
	Update(ctx context.Context, userId string, userAttributes *UserUpdateAttributes, opts ...UserPatchOptions) (*http.Response, error)
	Delete(ctx context.Context, userId string, opts ...UserDeleteOptions) (*http.Response, error)
}

type UsersService struct {
	client *Client
}

var _ UsersServicers = (*UsersService)(nil)

type UserGetOptions func(*usersGetOptions)

type UserPatchOptions func(*usersPatchOptions)

type UserDeleteOptions func(*usersDeleteOptions)

type usersGetOptions struct {
	userId              string `url:"userId,omitempty"`
	queryParameters     *UserQueryParameters
	listQueryParameters *ListUserQueryParameters
}

type usersPatchOptions struct {
	queryParameters *UserPatchQueryParameters
}

type usersDeleteOptions struct {
	queryParameters *UserDeleteQueryParameters
}

type UserPatchQueryParameters struct {
	LoginType            enums.LoginType `url:"login_type,omitempty"`
	RemoveTSPCredentials bool            `url:"remove_tsp_credentials,omitempty"`
}

type UserQueryParameters struct {
	LoginType        enums.LoginType `url:"login_type,omitempty"`
	EncryptedEmail   bool            `url:"encrypted_email,omitempty"`
	SearchByUniqueID bool            `url:"search_by_unique_id,omitempty"`
}

type UserDeleteQueryParameters struct {
	EncryptedEmail     bool                   `url:"encrypted_email,omitempty"`
	Action             enums.UserDeleteAction `url:"action,omitempty"`
	TransferEmail      string                 `url:"transfer_email,omitempty"`
	TransferMeeting    bool                   `url:"transfer_meeting,omitempty"`
	TransferWebinar    bool                   `url:"transfer_webinar,omitempty"`
	TransferRecording  bool                   `url:"transfer_recording,omitempty"`
	TransferWhiteboard bool                   `url:"transfer_whiteboard,omitempty"`
	TransferClipfiles  bool                   `url:"transfer_clipfiles,omitempty"`
	TransferNotes      bool                   `url:"transfer_notes,omitempty"`
	TransferVisitors   bool                   `url:"transfer_visitors,omitempty"`
	TransferDocs       bool                   `url:"transfer_docs,omitempty"`
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

func WithUserPatchQueryParameters(params *UserPatchQueryParameters) UserPatchOptions {
	return func(o *usersPatchOptions) {
		o.queryParameters = params
	}
}

func WithUserDeleteQueryParameters(params *UserDeleteQueryParameters) UserDeleteOptions {
	return func(o *usersDeleteOptions) {
		o.queryParameters = params
	}
}

func (u *UsersService) Get(ctx context.Context, opts ...UserGetOptions) ([]*models.User, *http.Response, error) {
	options := usersGetOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.queryParameters != nil && options.listQueryParameters != nil {
		return nil, nil, errors.New("Cannot use both UserQueryParameters and ListUserQueryParameters at the same time")
	}

	query := any(nil)
	if options.queryParameters != nil {
		query = options.queryParameters
	} else if options.listQueryParameters != nil {
		query = options.listQueryParameters
	}

	if options.userId != "" {
		endpoint := fmt.Sprintf("/users/%s", url.PathEscape(options.userId))
		user := &models.User{}
		res, err := u.client.request(ctx, http.MethodGet, endpoint, query, nil, user)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		return []*models.User{user}, res, nil
	}

	type response struct {
		*PaginationResponse
		Users []*models.User `json:"users"`
	}

	queryResponse := &response{}
	var users []*models.User

	endpoint := "/users/"

	res, err := u.client.request(ctx, http.MethodGet, endpoint, query, nil, queryResponse)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	users = append(users, queryResponse.Users...)

	type usersListPageQuery struct {
		*ListUserQueryParameters
		*PaginationOptions
	}
	for {
		if queryResponse.NextPageToken == "" {
			break
		}
		nextPageToken := queryResponse.NextPageToken
		pageQuery := &usersListPageQuery{
			ListUserQueryParameters: options.listQueryParameters,
			PaginationOptions: &PaginationOptions{
				NextPageToken: &nextPageToken,
			},
		}
		res, err = u.client.request(ctx, http.MethodGet, "/users/", pageQuery, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		users = append(users, queryResponse.Users...)
	}

	return users, res, nil
}

type UserAttributes struct {
	Email           string               `json:"email"`
	UserType        enums.UserType       `json:"type"`
	DisplayName     string               `json:"display_name,omitempty"`
	DivisionIds     []string             `json:"division_ids,omitempty"`
	Feature         models.Feature       `json:"feature"`
	FirstName       string               `json:"first_name,omitempty"`
	LastName        string               `json:"last_name,omitempty"`
	LicenseInfoList []models.LicenseInfo `json:"license_info_list,omitempty"`
	Password        string               `json:"password,omitempty"`
	PlanUnitedType  enums.PlanUnitedType `json:"plan_united_type,omitempty"`
}

func (u *UsersService) Create(ctx context.Context, action enums.UserCreateAction, userAttributes UserAttributes) (*models.User, *http.Response, error) {
	if userAttributes.Email == "" || userAttributes.UserType == 0 {
		return &models.User{}, nil, errors.New("Email and UserType are required fields")
	}
	type body struct {
		Action   enums.UserCreateAction `json:"action"`
		UserInfo UserAttributes         `json:"user_info"`
	}
	requestBody := &body{
		Action:   action,
		UserInfo: userAttributes,
	}

	var response models.User

	res, err := u.client.request(ctx, http.MethodPost, "/users/", nil, requestBody, &response)
	if err != nil {
		return &models.User{}, res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusCreated {
		return &models.User{}, res, fmt.Errorf("Expected status code %d, got %d", http.StatusCreated, res.StatusCode)
	}
	return &response, res, nil
}

type UserUpdateAttributes struct {
	CMSUserID        string                    `json:"cms_user_id,omitempty"`
	Company          string                    `json:"company,omitempty"`
	CostCenter       string                    `json:"cost_center,omitempty"`
	CustomAttributes []models.CustomAttributes `json:"custom_attributes,omitempty"`
	Dept             string                    `json:"dept,omitempty"`
	DivisionIds      []string                  `json:"division_ids,omitempty"`
	FirstName        string                    `json:"first_name,omitempty"`
	GroupId          string                    `json:"group_id,omitempty"`
	HostKey          string                    `json:"host_key,omitempty"`
	JobTitle         string                    `json:"job_title,omitempty"`
	Language         string                    `json:"language,omitempty"`
	LastName         string                    `json:"last_name,omitempty"`
	AboutMe          string                    `json:"about_me,omitempty"`
	DisplayName      string                    `json:"display_name,omitempty"`
	Feature          models.Feature            `json:"feature"`
	LicenseInfoList  []models.LicenseInfo      `json:"license_info_list,omitempty"`
	LinkedinURL      string                    `json:"linkedin_url,omitempty"`
	Location         string                    `json:"location,omitempty"`
	Manager          string                    `json:"manager,omitempty"`
	PhoneNumbers     []models.PhoneNumber      `json:"phone_numbers"`
	PlanUnitedType   enums.PlanUnitedType      `json:"plan_united_type,omitempty"`
	PMI              int64                     `json:"pmi,omitempty"`
	Pronouns         string                    `json:"pronouns,omitempty"`
	PronounsOption   int                       `json:"pronouns_option,omitempty"`
	Timezone         string                    `json:"timezone,omitempty"`
	Type             enums.UserType            `json:"type,omitempty"`
	UsePMI           bool                      `json:"use_pmi,omitempty"`
	VanityName       string                    `json:"vanity_name,omitempty"`
	ZoomOneType      enums.ZoomOneType         `json:"zoom_one_type,omitempty"`
}

func (u *UsersService) Update(ctx context.Context, userId string, userAttributes *UserUpdateAttributes, opts ...UserPatchOptions) (*http.Response, error) {
	options := usersPatchOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	res, err := u.client.request(ctx, http.MethodPatch, fmt.Sprintf("/users/%s", url.PathEscape(userId)), options.queryParameters, userAttributes, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}

func (u *UsersService) Delete(ctx context.Context, userId string, opts ...UserDeleteOptions) (*http.Response, error) {
	options := usersDeleteOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	res, err := u.client.request(ctx, http.MethodDelete, fmt.Sprintf("/users/%s", url.PathEscape(userId)), options.queryParameters, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}
