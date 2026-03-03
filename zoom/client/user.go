package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Hawkeye-Claims/go-zoom/zoom/enums"
	"github.com/Hawkeye-Claims/go-zoom/zoom/models"
)

// UsersServicers is the interface implemented by UsersService. It declares the
// complete set of user management operations available through the Zoom Users
// API.
type UsersServicers interface {
	// Get retrieves one or more users. When WithUserId is supplied it returns a
	// single user; otherwise it lists all users and paginates automatically.
	Get(ctx context.Context, opts ...UserGetOptions) ([]*models.User, *http.Response, error)
	// Create creates a new Zoom user with the given action and attributes.
	Create(ctx context.Context, action enums.UserCreateAction, userAttributes UserAttributes) (*models.User, *http.Response, error)
	// Update patches an existing user identified by userId with the supplied
	// attributes.
	Update(ctx context.Context, userId string, userAttributes *UserUpdateAttributes, opts ...UserPatchOptions) (*http.Response, error)
	// Delete removes a user identified by userId from the account.
	Delete(ctx context.Context, userId string, opts ...UserDeleteOptions) (*http.Response, error)
}

// UsersService implements UsersServicers and provides access to the Zoom Users
// API endpoints.
type UsersService struct {
	client *Client
}

// Compile-time assertion that UsersService satisfies the UsersServicers
// interface.
var _ UsersServicers = (*UsersService)(nil)

// UserGetOptions is a functional option for configuring a user Get request.
type UserGetOptions func(*usersGetOptions)

// UserPatchOptions is a functional option for configuring a user Update
// (PATCH) request.
type UserPatchOptions func(*usersPatchOptions)

// UserDeleteOptions is a functional option for configuring a user Delete
// request.
type UserDeleteOptions func(*usersDeleteOptions)

// usersGetOptions holds the resolved configuration for a user Get call.
type usersGetOptions struct {
	userId              string `url:"userId,omitempty"`
	queryParameters     *UserQueryParameters
	listQueryParameters *ListUserQueryParameters
}

// usersPatchOptions holds the resolved configuration for a user Update call.
type usersPatchOptions struct {
	queryParameters *UserPatchQueryParameters
}

// usersDeleteOptions holds the resolved configuration for a user Delete call.
type usersDeleteOptions struct {
	queryParameters *UserDeleteQueryParameters
}

// UserPatchQueryParameters holds optional query parameters for a user Update
// (PATCH) request.
type UserPatchQueryParameters struct {
	// LoginType specifies the login method used to identify the user.
	LoginType enums.LoginType `url:"login_type,omitempty"`
	// RemoveTSPCredentials, when true, removes the user's TSP credentials.
	RemoveTSPCredentials bool `url:"remove_tsp_credentials,omitempty"`
}

// UserQueryParameters holds optional query parameters for a single-user Get
// request.
type UserQueryParameters struct {
	// LoginType specifies the login method used to identify the user.
	LoginType enums.LoginType `url:"login_type,omitempty"`
	// EncryptedEmail, when true, treats the supplied user ID as an encrypted
	// email address.
	EncryptedEmail bool `url:"encrypted_email,omitempty"`
	// SearchByUniqueID, when true, looks up the user by their unique ID rather
	// than email.
	SearchByUniqueID bool `url:"search_by_unique_id,omitempty"`
}

// UserDeleteQueryParameters holds optional query parameters for a user Delete
// request, including options to transfer the deleted user's data to another
// account.
type UserDeleteQueryParameters struct {
	// EncryptedEmail, when true, treats the user ID as an encrypted email.
	EncryptedEmail bool `url:"encrypted_email,omitempty"`
	// Action specifies the delete action to perform (e.g. disassociate or
	// delete).
	Action enums.UserDeleteAction `url:"action,omitempty"`
	// TransferEmail is the email address of the user who will receive
	// transferred data.
	TransferEmail string `url:"transfer_email,omitempty"`
	// TransferMeeting, when true, transfers the deleted user's meetings.
	TransferMeeting bool `url:"transfer_meeting,omitempty"`
	// TransferWebinar, when true, transfers the deleted user's webinars.
	TransferWebinar bool `url:"transfer_webinar,omitempty"`
	// TransferRecording, when true, transfers the deleted user's recordings.
	TransferRecording bool `url:"transfer_recording,omitempty"`
	// TransferWhiteboard, when true, transfers the deleted user's whiteboards.
	TransferWhiteboard bool `url:"transfer_whiteboard,omitempty"`
	// TransferClipfiles, when true, transfers the deleted user's clip files.
	TransferClipfiles bool `url:"transfer_clipfiles,omitempty"`
	// TransferNotes, when true, transfers the deleted user's notes.
	TransferNotes bool `url:"transfer_notes,omitempty"`
	// TransferVisitors, when true, transfers the deleted user's visitor data.
	TransferVisitors bool `url:"transfer_visitors,omitempty"`
	// TransferDocs, when true, transfers the deleted user's documents.
	TransferDocs bool `url:"transfer_docs,omitempty"`
}

// ListUserQueryParameters holds optional query parameters for listing users.
type ListUserQueryParameters struct {
	// Status filters the list to users with the given status.
	Status enums.UserStatus `url:"status,omitempty"`
	// RoleID filters the list to users belonging to the specified role.
	RoleID string `url:"role_id,omitempty"`
	// IncludeFields specifies additional fields to include in the response.
	IncludeFields enums.UserFields `url:"include_fields,omitempty"`
}

// WithUserId returns a UserGetOptions that fetches a single user by their
// Zoom user ID or email address.
func WithUserId(userId string) UserGetOptions {
	return func(o *usersGetOptions) {
		o.userId = userId
	}
}

// WithUserQueryParameters returns a UserGetOptions that attaches the given
// query parameters to a single-user Get request.
func WithUserQueryParameters(params *UserQueryParameters) UserGetOptions {
	return func(o *usersGetOptions) {
		o.queryParameters = params
	}
}

// WithListUserQueryParameters returns a UserGetOptions that attaches the given
// query parameters to a list-users request.
func WithListUserQueryParameters(params *ListUserQueryParameters) UserGetOptions {
	return func(o *usersGetOptions) {
		o.listQueryParameters = params
	}
}

// WithUserPatchQueryParameters returns a UserPatchOptions that attaches the
// given query parameters to a user Update request.
func WithUserPatchQueryParameters(params *UserPatchQueryParameters) UserPatchOptions {
	return func(o *usersPatchOptions) {
		o.queryParameters = params
	}
}

// WithUserDeleteQueryParameters returns a UserDeleteOptions that attaches the
// given query parameters to a user Delete request.
func WithUserDeleteQueryParameters(params *UserDeleteQueryParameters) UserDeleteOptions {
	return func(o *usersDeleteOptions) {
		o.queryParameters = params
	}
}

// Get retrieves Zoom users. When WithUserId is provided it returns a slice
// containing only that user. Without a user ID it lists all users, following
// pagination automatically so the returned slice contains every user.
// It is an error to supply both WithUserQueryParameters and
// WithListUserQueryParameters simultaneously.
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

// UserAttributes holds the fields required and optional when creating a new
// Zoom user.
type UserAttributes struct {
	// Email is the user's email address. Required.
	Email string `json:"email"`
	// UserType specifies the user's license type. Required.
	UserType enums.UserType `json:"type"`
	// DisplayName is the user's display name.
	DisplayName string `json:"display_name,omitempty"`
	// DivisionIds lists the division IDs the user belongs to.
	DivisionIds []string `json:"division_ids,omitempty"`
	// Feature holds the user's feature entitlements.
	Feature models.Feature `json:"feature"`
	// FirstName is the user's first name.
	FirstName string `json:"first_name,omitempty"`
	// LastName is the user's last name.
	LastName string `json:"last_name,omitempty"`
	// LicenseInfoList specifies additional licence assignments for the user.
	LicenseInfoList []models.LicenseInfo `json:"license_info_list,omitempty"`
	// Password is the user's initial password (where applicable).
	Password string `json:"password,omitempty"`
	// PlanUnitedType specifies the Zoom United plan type for the user.
	PlanUnitedType enums.PlanUnitedType `json:"plan_united_type,omitempty"`
}

// Create creates a new Zoom user using the specified action and user
// attributes. Both Email and UserType in userAttributes are required. Returns
// the newly created user.
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

// UserUpdateAttributes holds the fields that can be patched on an existing
// Zoom user. All fields are optional.
type UserUpdateAttributes struct {
	// CMSUserID is the user's CMS user ID.
	CMSUserID string `json:"cms_user_id,omitempty"`
	// Company is the user's company name.
	Company string `json:"company,omitempty"`
	// CostCenter is the cost center associated with the user.
	CostCenter string `json:"cost_center,omitempty"`
	// CustomAttributes contains custom key/value attributes for the user.
	CustomAttributes []models.CustomAttributes `json:"custom_attributes,omitempty"`
	// Dept is the user's department.
	Dept string `json:"dept,omitempty"`
	// DivisionIds lists the division IDs the user should be assigned to.
	DivisionIds []string `json:"division_ids,omitempty"`
	// FirstName is the user's first name.
	FirstName string `json:"first_name,omitempty"`
	// GroupId is the ID of the group the user should be placed in.
	GroupId string `json:"group_id,omitempty"`
	// HostKey is the user's host key.
	HostKey string `json:"host_key,omitempty"`
	// JobTitle is the user's job title.
	JobTitle string `json:"job_title,omitempty"`
	// Language is the user's preferred language (e.g. "en-US").
	Language string `json:"language,omitempty"`
	// LastName is the user's last name.
	LastName string `json:"last_name,omitempty"`
	// AboutMe is a brief biography or description for the user.
	AboutMe string `json:"about_me,omitempty"`
	// DisplayName is the user's public display name.
	DisplayName string `json:"display_name,omitempty"`
	// Feature holds the user's updated feature entitlements.
	Feature *models.Feature `json:"feature,omitempty"`
	// LicenseInfoList specifies updated licence assignments for the user.
	LicenseInfoList []models.LicenseInfo `json:"license_info_list,omitempty"`
	// LinkedinURL is the URL of the user's LinkedIn profile.
	LinkedinURL string `json:"linkedin_url,omitempty"`
	// Location is the user's office location.
	Location string `json:"location,omitempty"`
	// Manager is the email or user ID of the user's manager.
	Manager string `json:"manager,omitempty"`
	// PhoneNumbers is the list of phone numbers associated with the user.
	PhoneNumbers []models.PhoneNumber `json:"phone_numbers,omitempty"`
	// PlanUnitedType specifies the updated Zoom United plan type.
	PlanUnitedType enums.PlanUnitedType `json:"plan_united_type,omitempty"`
	// PMI is the user's Personal Meeting ID.
	PMI int64 `json:"pmi,omitempty"`
	// Pronouns are the user's preferred pronouns.
	Pronouns string `json:"pronouns,omitempty"`
	// PronounsOption controls whether and how pronouns are displayed.
	PronounsOption int `json:"pronouns_option,omitempty"`
	// Timezone is the user's time zone (e.g. "America/New_York").
	Timezone string `json:"timezone,omitempty"`
	// Type is the user's updated licence type.
	Type enums.UserType `json:"type,omitempty"`
	// UsePMI, when true, uses the Personal Meeting ID for instant meetings.
	UsePMI bool `json:"use_pmi,omitempty"`
	// VanityName is the user's personal meeting URL vanity name.
	VanityName string `json:"vanity_name,omitempty"`
	// ZoomOneType specifies the Zoom One plan type for the user.
	ZoomOneType enums.ZoomOneType `json:"zoom_one_type,omitempty"`
}

// Update patches an existing Zoom user identified by userId with the supplied
// attributes. Optional query parameters can be supplied via UserPatchOptions.
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

// Delete removes a Zoom user identified by userId. Optional query parameters
// (e.g. data-transfer settings) can be supplied via UserDeleteOptions.
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
