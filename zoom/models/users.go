// Package models defines data structures that represent Zoom API response
// objects. These types are used as return values by the service methods in the
// client package.
package models

import (
	"time"

	"github.com/Hawkeye-Claims/go-zoom/zoom/enums"
)

// User represents a Zoom user account as returned by the Users API.
type User struct {
	// ID is the unique identifier for the user.
	ID string `json:"id"`
	// CreatedAt is the time the user account was created.
	CreatedAt time.Time `json:"created_at"`
	// Dept is the department the user belongs to.
	Dept string `json:"dept"`
	// Email is the user's email address.
	Email string `json:"email"`
	// FirstName is the user's first name.
	FirstName string `json:"first_name"`
	// LastClientVersion is the version of the Zoom client the user last used.
	LastClientVersion string `json:"last_client_version"`
	// LastLoginTime is the time the user last signed in.
	LastLoginTime time.Time `json:"last_login_time"`
	// LastName is the user's last name.
	LastName string `json:"last_name"`
	// Pmi is the user's Personal Meeting ID.
	Pmi int64 `json:"pmi"`
	// RoleName is the display name of the role assigned to the user.
	RoleName string `json:"role_name"`
	// Timezone is the user's IANA timezone identifier.
	Timezone string `json:"timezone"`
	// Type is the user's license tier.
	Type enums.UserType `json:"type"`
	// UsePmi indicates whether the user's PMI is used for instant meetings.
	UsePmi bool `json:"use_pmi"`
	// DisplayName is the name shown to other participants in meetings.
	DisplayName string `json:"display_name"`
	// AccountID is the ID of the account the user belongs to.
	AccountID string `json:"account_id"`
	// AccountNumber is the numeric identifier of the account.
	AccountNumber int `json:"account_number"`
	// CmsUserID is the user's ID in an external CMS.
	CmsUserID string `json:"cms_user_id"`
	// Company is the company the user is associated with.
	Company string `json:"company"`
	// UserCreatedAt is an alternative creation timestamp.
	UserCreatedAt time.Time `json:"user_created_at"`
	// CustomAttributes contains any account-level custom attributes assigned to
	// the user. Only present when the custom_attributes include_field is requested.
	CustomAttributes []CustomAttributes `json:"custom_attributes"`
	// EmployeeUniqueID is the user's employee ID.
	EmployeeUniqueID string `json:"employee_unique_id"`
	// GroupIds lists the IDs of groups the user is a member of.
	GroupIds []string `json:"group_ids"`
	// DivisionIds lists the division IDs the user is assigned to.
	DivisionIds []string `json:"division_ids"`
	// ImGroupIds lists the IDs of IM groups the user belongs to.
	ImGroupIds []string `json:"im_group_ids"`
	// Jid is the user's Zoom JID (Jabber ID).
	Jid string `json:"jid"`
	// JobTitle is the user's job title.
	JobTitle string `json:"job_title"`
	// CostCenter is the cost center the user is assigned to.
	CostCenter string `json:"cost_center"`
	// Language is the user's preferred display language.
	Language string `json:"language"`
	// Location is the user's office location.
	Location string `json:"location"`
	// LoginTypes lists the authentication methods registered for the user.
	LoginTypes []enums.LoginType `json:"login_types"`
	// Manager is the email address of the user's manager.
	Manager string `json:"manager"`
	// PersonalMeetingURL is the URL for the user's personal meeting room.
	PersonalMeetingURL string `json:"personal_meeting_url"`
	// PhoneNumbers lists the phone numbers associated with the user.
	PhoneNumbers []PhoneNumber `json:"phone_numbers"`
	// PicURL is the URL of the user's profile picture.
	PicURL string `json:"pic_url"`
	// PlanUnitedType is the Zoom United bundle plan assigned to the user.
	PlanUnitedType enums.PlanUnitedType `json:"plan_united_type"`
	// Pronouns is the user's preferred pronouns.
	Pronouns string `json:"pronouns"`
	// PronounsOption is the visibility setting for the user's pronouns.
	PronounsOption int `json:"pronouns_option"`
	// RoleID is the ID of the role assigned to the user.
	RoleID string `json:"role_id"`
	// Status is the activation state of the user account.
	Status enums.UserStatus `json:"status"`
	// VanityURL is the user's custom personal meeting URL alias.
	VanityURL string `json:"vanity_url"`
	// Verified indicates whether the user's email address has been verified (1)
	// or not (0).
	Verified int `json:"verified"`
	// Cluster identifies the Zoom data centre cluster the user is hosted on.
	Cluster string `json:"cluster"`
	// ZoomOneType is the Zoom One bundle tier assigned to the user.
	ZoomOneType enums.ZoomOneType `json:"zoom_one_type,omitempty"`
	// LicenseInfoList contains the add-on licenses assigned to the user.
	LicenseInfoList []LicenseInfo `json:"license_info_list"`
}

// CustomAttributes holds a single custom attribute key-value pair assigned to a
// user. Custom attributes are account-defined metadata fields.
type CustomAttributes struct {
	// Key is the unique identifier of the custom attribute.
	Key string `json:"key"`
	// Name is the human-readable label for the custom attribute.
	Name string `json:"name"`
	// Value is the attribute's value for this user.
	Value string `json:"value"`
}

// PhoneNumber represents a phone number associated with a Zoom user account.
type PhoneNumber struct {
	// ID is the unique identifier for this phone number entry.
	ID string `json:"id"`
	// Code is the country calling code (e.g. "1" for the US).
	Code string `json:"code"`
	// Country is the ISO 3166-1 alpha-2 country code.
	Country string `json:"country"`
	// Label categorises the phone number (mobile, office, home, or fax).
	Label enums.PhoneNumberLabel `json:"label"`
	// Number is the phone number in E.164 format.
	Number string `json:"number"`
	// Verified indicates whether the phone number has been verified.
	Verified bool `json:"verified"`
}

// LicenseInfo describes a single add-on license assigned to a user.
type LicenseInfo struct {
	// LicenseType is the string identifier of the license type.
	LicenseType string `json:"license_type"`
	// LicenseOption is the numeric option value for the license.
	LicenseOption enums.LicenseOption `json:"license_option"`
	// SubscriptionID is the billing subscription ID for this license.
	SubscriptionID string `json:"subscription_id"`
}

// Feature holds feature flags for a user's Zoom account.
type Feature struct {
	// ZoomOneType is the Zoom One bundle tier enabled for the user.
	ZoomOneType enums.ZoomOneType `json:"zoom_one_type,omitempty"`
	// ZoomPhone indicates whether the user has the Zoom Phone add-on enabled.
	ZoomPhone bool `json:"zoom_phone,omitempty"`
}
