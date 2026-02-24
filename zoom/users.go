package zoom

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type UsersServicer interface {
	Get(ctx context.Context, opts ...UsersGetOptions) ([]*User, *http.Response, error)
	Post(ctx context.Context, action UserPostAction, email string, userType UserType, userAttributes *UserAttributes) (*User, *http.Response, error)
	Patch(ctx context.Context, userID string, userAttributes *UserAttributes, opts ...UserPatchOptions) (*User, *http.Response, error)
}

type UsersService struct {
	client *Client
}

var _ UsersServicer = (*UsersService)(nil)

type UsersGetOptions func(*usersGetOptions)

type UserPatchOptions func(*userPatchOptions)

type UserQueryParameters struct {
	LoginType        LoginType `url:"login_type,omitempty"`
	EncryptedEmail   bool      `url:"encrypted_email,omitempty"`
	SearchByUniqueID bool      `url:"search_by_unique_id,omitempty"`
}

type ListUserQueryParameters struct {
	Status        UserStatus  `url:"status,omitempty"`
	RoleID        string      `url:"role_id,omitempty"`
	IncludeFields UserFields  `url:"include_fields,omitempty"`
	LicenseType   LicenseType `url:"license,omitempty"`
}

type UserPatchQueryParameters struct {
	LoginType            LoginType `url:"login_type,omitempty"`
	RemoveTSPCredentials bool      `url:"remove_tsp_credentials,omitempty"`
}

type usersGetOptions struct {
	userId              string `url:"userId,omitempty"`
	queryParameters     *UserQueryParameters
	listQueryParameters *ListUserQueryParameters
}

type userPatchOptions struct {
	queryParameters  *UserPatchQueryParameters
	aboutMe          string
	linkedin         string
	customAttributes map[string]string
}

func WithUserID(userId string) UsersGetOptions {
	return func(o *usersGetOptions) {
		o.userId = userId
	}
}

func WithUserQueryParameters(params *UserQueryParameters) UsersGetOptions {
	return func(o *usersGetOptions) {
		o.queryParameters = params
	}
}

func WithListUserQueryParameters(params *ListUserQueryParameters) UsersGetOptions {
	return func(o *usersGetOptions) {
		o.listQueryParameters = params
	}
}

func WithUserPatchQueryParameters(params *UserPatchQueryParameters) UserPatchOptions {
	return func(o *userPatchOptions) {
		o.queryParameters = params
	}
}

func WithAboutMe(aboutMe string) UserPatchOptions {
	return func(o *userPatchOptions) {
		o.aboutMe = aboutMe
	}
}

func WithLinkedIn(linkedin string) UserPatchOptions {
	return func(o *userPatchOptions) {
		o.linkedin = linkedin
	}
}

type UserAttributes struct {
	DisplayName    string            `json:"display_name,omitempty"`
	DivisionIDs    []string          `json:"division_ids,omitempty"`
	Feature        UserFeature       `json:"feature,omitempty"`
	FirstName      string            `json:"first_name,omitempty"`
	LastName       string            `json:"last_name,omitempty"`
	LicenseInfo    []UserLicenseInfo `json:"license_info,omitempty"`
	Password       string            `json:"password,omitempty"`
	PlanUnitedType string            `json:"plan_united_type,omitempty"`
}

type UserFeature struct {
	ZoomOneType ZoomOneType `json:"zoom_one_type,omitempty"`
	ZoomPhone   bool        `json:"zoom_phone,omitempty"`
}

type UserLicenseInfo struct {
	LicenseOption   int    `json:"license_option"`
	LicenseTypeName string `json:"license_type"`
	SubscriptionID  string `json:"subscription_id"`
}

type User struct {
	AccountID          string                          `json:"account_id"`
	AccountNumber      int                             `json:"account_number"`
	Cluster            string                          `json:"cluster"`
	CMSUserID          string                          `json:"cms_user_id"`
	Company            string                          `json:"company"`
	CostCenter         string                          `json:"cost_center"`
	CustomAttributes   []*UsersListItemCustomAttribute `json:"custom_attributes"`
	Dept               string                          `json:"dept"`
	DisplayName        string                          `json:"display_name"`
	DivisionIDs        []string                        `json:"division_ids"`
	Email              string                          `json:"email"`
	EmployeeUniqueID   string                          `json:"employee_unique_id"`
	FirstName          string                          `json:"first_name"`
	GroupIDs           []string                        `json:"group_ids"`
	ID                 string                          `json:"id"`
	ImGroupIDs         []string                        `json:"im_group_ids"`
	JID                string                          `json:"jid"`
	JobTitle           string                          `json:"job_title"`
	Language           string                          `json:"language"`
	LastClientVersion  string                          `json:"last_client_version"`
	LastLoginTime      time.Time                       `json:"last_login_time"`
	LastName           string                          `json:"last_name"`
	LicenseInfoList    []UserLicenseInfo               `json:"license_info_list"`
	Location           string                          `json:"location"`
	LoginTypes         []LoginType                     `json:"login_types"`
	Manager            string                          `json:"manager"`
	PersonalMeetingURL string                          `json:"personal_meeting_url"`
	PhoneNumbers       []PhoneNumber                   `json:"phone_numbers"`
	PicURL             string                          `json:"pic_url"`
	PlanUnitedType     string                          `json:"plan_united_type"`
	PMI                int64                           `json:"pmi"`
	Pronouns           string                          `json:"pronouns"`
	PronounOption      int                             `json:"pronoun_option"`
	RoleID             string                          `json:"role_id"`
	RoleName           string                          `json:"role_name"`
	Status             UserStatus                      `json:"status"`
	Timezone           string                          `json:"timezone"`
	Type               UserType                        `json:"type"`
	UsePMI             bool                            `json:"use_pmi"`
	UserCreatedAt      time.Time                       `json:"user_created_at"`
	VanityURL          string                          `json:"vanity_url"`
	Verified           int                             `json:"verified"`
	ZoomOneType        ZoomOneType                     `json:"zoom_one_type"`
}

type PhoneNumber struct {
	Code     string `json:"code"`
	Country  string `json:"country"`
	Number   string `json:"number"`
	Verified bool   `json:"verified"`
}

type UsersListItemCustomAttribute struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (u *UsersService) Get(ctx context.Context, opts ...UsersGetOptions) ([]*User, *http.Response, error) {
	options := usersGetOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.queryParameters != nil && options.listQueryParameters != nil {
		return nil, nil, fmt.Errorf("Cannot use both UserQueryParameters and ListUserQueryParameters in the same request")
	}
	query := interface{}(nil)
	if options.queryParameters != nil {
		query = options.queryParameters
	} else if options.listQueryParameters != nil {
		query = options.listQueryParameters
	}

	var users []*User

	type response struct {
		*PaginationResponse
		Users []*User `json:"users"`
	}

	queryResponse := &response{}

	res, err := u.client.request(ctx, http.MethodGet, "/users/", query, nil, queryResponse)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	users = append(users, queryResponse.Users...)

	if options.userId != "" && len(users) == 1 {
		return users, res, nil
	}

	for {
		if queryResponse.NextPageToken == "" {
			break
		}
		nextPageToken := queryResponse.NextPageToken
		queryOption := &PaginationOptions{
			NextPageToken: &nextPageToken,
		}
		res, err = u.client.request(ctx, http.MethodGet, "/users/", queryOption, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		users = append(users, queryResponse.Users...)
	}
	return users, res, nil
}

func (u *UsersService) Post(ctx context.Context, action UserPostAction, email string, userType UserType, userAttributes *UserAttributes) (*User, *http.Response, error) {
	type userInfo struct {
		Email    string   `json:"email"`
		UserType UserType `json:"type"`
		*UserAttributes
	}
	type body struct {
		Action   UserPostAction `json:"action"`
		UserInfo userInfo       `json:"user_info"`
	}
	requestBody := &body{
		Action: action,
		UserInfo: userInfo{
			email,
			userType,
			userAttributes,
		},
	}
	type response struct {
		Email     string   `json:"email"`
		FirstName string   `json:"first_name"`
		ID        string   `json:"id"`
		LastName  string   `json:"last_name"`
		Type      UserType `json:"type"`
	}

	var respBody response

	res, err := u.client.request(ctx, http.MethodPost, "/users/", nil, requestBody, respBody)
	if err != nil {
		return &User{}, res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusCreated {
		return &User{}, res, fmt.Errorf("Expected status code %d, got %d", http.StatusCreated, res.StatusCode)
	}
	createdUser := &User{
		Email:     respBody.Email,
		FirstName: respBody.FirstName,
		ID:        respBody.ID,
		LastName:  respBody.LastName,
		Type:      respBody.Type,
	}
	return createdUser, res, nil
}

func (u *UsersService) Patch(ctx context.Context, userID string, userAttributes *UserAttributes, opts ...UserPatchOptions) (*User, *http.Response, error) {
	options := userPatchOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	type body struct {
		AboutMe    string `json:"about_me,omitempty"`
		LinkedinRL string `json:"linkedin_url,omitempty"`
		*UserAttributes
	}
	requestBody := &body{
		options.aboutMe,
		options.linkedin,
		userAttributes,
	}

	res, err := u.client.request(ctx, http.MethodPatch, fmt.Sprintf("/users/%s", userID), options.queryParameters, requestBody, nil)
	if err != nil {
		return &User{}, res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return &User{}, res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return &User{ID: userID}, res, nil
}
