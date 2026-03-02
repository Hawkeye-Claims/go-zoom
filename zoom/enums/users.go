// Package enums defines enumeration types and their constants for the Zoom API.
package enums

// LoginType identifies the authentication method used by a Zoom account.
type LoginType int

const (
	// FacebookOAuth represents login via Facebook OAuth.
	FacebookOAuth LoginType = 0
	// GoogleOAuth represents login via Google OAuth.
	GoogleOAuth LoginType = 1
	// PhoneNumberAuth represents login via phone number.
	PhoneNumberAuth LoginType = 11
	// WeChatAuth represents login via WeChat.
	WeChatAuth LoginType = 21
	// AlipayAuth represents login via Alipay.
	AlipayAuth LoginType = 23
	// AppleOAuth represents login via Apple OAuth.
	AppleOAuth LoginType = 24
	// MicrosoftOAuth represents login via Microsoft OAuth.
	MicrosoftOAuth LoginType = 27
	// MobileDevice represents login from a mobile device.
	MobileDevice LoginType = 97
	// RingCentralOAuth represents login via RingCentral OAuth.
	RingCentralOAuth LoginType = 98
	// APIUser represents an API-only user.
	APIUser LoginType = 99
	// ZoomWorkEmail represents login via a Zoom work email address.
	ZoomWorkEmail LoginType = 100
	// SingleSignOn represents login via SSO.
	SingleSignOn LoginType = 101
)

// UserStatus represents the activation state of a Zoom user account.
type UserStatus string

const (
	// ActiveUser identifies an active, fully enabled user account.
	ActiveUser UserStatus = "active"
	// InactiveUser identifies a deactivated user account.
	InactiveUser UserStatus = "inactive"
	// PendingUser identifies a user whose registration is not yet complete.
	PendingUser UserStatus = "pending"
)

// UserFields enumerates optional fields that can be included in a user list
// response via the include_fields query parameter.
type UserFields string

const (
	// CustomAttributesField requests inclusion of custom attribute data.
	CustomAttributesField UserFields = "custom_attributes"
	// HostKeyField requests inclusion of the host key.
	HostKeyField UserFields = "host_key"
)

// LicenseType identifies a Zoom add-on license assigned to a user.
type LicenseType string

const (
	// ZoomWorkforceManagement identifies the Zoom Workforce Management license.
	ZoomWorkforceManagement LicenseType = "zoom_workforce_management"
	// ZoomComplianceManagement identifies the Zoom Compliance Management license.
	ZoomComplianceManagement LicenseType = "zoom_compliance_management"
)

// UserCreateAction specifies how a user should be provisioned when calling the
// create user endpoint.
type UserCreateAction string

const (
	// Create sends an activation email to the new user.
	Create UserCreateAction = "create"
	// AutoCreate creates the user without sending an activation email.
	AutoCreate UserCreateAction = "auto_create"
	// CustCreate creates a user managed by a customer account.
	CustCreate UserCreateAction = "cust_create"
	// SSOCreate creates a user who will sign in via SSO.
	SSOCreate UserCreateAction = "sso_create"
)

// UserDeleteAction controls what happens to the user's data when they are
// removed from the account.
type UserDeleteAction string

const (
	// Disassociate removes the user from the account but preserves their data.
	Disassociate UserDeleteAction = "disassociate"
	// Delete permanently deletes the user and their data.
	Delete UserDeleteAction = "delete"
)

// LicenseOption is the underlying integer type for license-related enumerations.
type LicenseOption int

// UserType describes the license tier of a Zoom user.
type UserType LicenseOption

const (
	// Basic is a free Zoom user with limited features.
	Basic UserType = 1
	// Licensed is a paid Zoom user with full feature access.
	Licensed UserType = 2
	// Unassigned is a user without an assigned license type.
	Unassigned UserType = 3
)

// ZoomOneType represents the Zoom One bundle tier assigned to a user.
type ZoomOneType LicenseOption

const (
	// ZoomWorkplaceEnterprise is the Zoom Workplace Enterprise tier.
	ZoomWorkplaceEnterprise ZoomOneType = 4
	// ZoomWorkplaceEnterprisePlus is the Zoom Workplace Enterprise Plus tier.
	ZoomWorkplaceEnterprisePlus ZoomOneType = 8
	// ZoomWorkplaceBusinessPlusUSCA is the Zoom Workplace Business Plus plan for US/Canada.
	ZoomWorkplaceBusinessPlusUSCA ZoomOneType = 16
	// ZoomWorkplaceBusinessPlusUKIR is the Zoom Workplace Business Plus plan for UK/Ireland.
	ZoomWorkplaceBusinessPlusUKIR ZoomOneType = 32
	// ZoomWorkplaceBusinessPlusAUNZ is the Zoom Workplace Business Plus plan for Australia/New Zealand.
	ZoomWorkplaceBusinessPlusAUNZ ZoomOneType = 64
	// ZoomWorkplaceBusinessPlusJapan is the Zoom Workplace Business Plus plan for Japan.
	ZoomWorkplaceBusinessPlusJapan ZoomOneType = 128
	// ZoomWorkplaceBusinessPlusGlobalSelect is the Zoom Workplace Business Plus Global Select plan.
	ZoomWorkplaceBusinessPlusGlobalSelect ZoomOneType = 33554432
	// ZoomWorkplaceEnterprisePremierUSA is the Zoom Workplace Enterprise Premier plan for USA.
	ZoomWorkplaceEnterprisePremierUSA ZoomOneType = 13417728
	// ZoomWorkplaceEnterprisePremierAUNZ is the Zoom Workplace Enterprise Premier plan for Australia/New Zealand.
	ZoomWorkplaceEnterprisePremierAUNZ ZoomOneType = 1073741824
	// ZoomWorkplaceEnterprisePremierUKIR is the Zoom Workplace Enterprise Premier plan for UK/Ireland.
	ZoomWorkplaceEnterprisePremierUKIR ZoomOneType = 536870912
	// ZoomWorkplaceEnterprisePremierJapan is the Zoom Workplace Enterprise Premier plan for Japan.
	ZoomWorkplaceEnterprisePremierJapan ZoomOneType = 268435456
	// ZoomWorkplaceProPlus is the Zoom Workplace Pro Plus tier.
	ZoomWorkplaceProPlus ZoomOneType = 4398046511104
)

// PhoneNumberLabel describes the category of a phone number associated with a
// user account.
type PhoneNumberLabel string

const (
	// MobileLabel indicates a mobile phone number.
	MobileLabel PhoneNumberLabel = "Mobile"
	// OfficeLabel indicates an office phone number.
	OfficeLabel PhoneNumberLabel = "Office"
	// HomeLabel indicates a home phone number.
	HomeLabel PhoneNumberLabel = "Home"
	// FaxLabel indicates a fax number.
	FaxLabel PhoneNumberLabel = "Fax"
)

// PlanUnitedType identifies a Zoom United bundle plan that combines meetings and
// phone in a single subscription.
type PlanUnitedType string

const (
	// ZoomUnitedProUnitedUSCA is the Zoom United Pro plan for US/Canada.
	ZoomUnitedProUnitedUSCA PlanUnitedType = "1"
	// ZoomUnitedProUnitedUKIR is the Zoom United Pro plan for UK/Ireland.
	ZoomUnitedProUnitedUKIR PlanUnitedType = "2"
	// ZoomUnitedProUnitedAUNZ is the Zoom United Pro plan for Australia/New Zealand.
	ZoomUnitedProUnitedAUNZ PlanUnitedType = "4"
	// ZoomUnitedProUnitedGlobalSelect is the Zoom United Pro Global Select plan.
	ZoomUnitedProUnitedGlobalSelect PlanUnitedType = "8"
	// ZoomUnitedProZoomPhonePro is the Zoom United Pro + Zoom Phone Pro bundle.
	ZoomUnitedProZoomPhonePro PlanUnitedType = "16"
	// ZoomUnitedBizUnitedUSCA is the Zoom United Business plan for US/Canada.
	ZoomUnitedBizUnitedUSCA PlanUnitedType = "32"
	// ZoomUnitedBizUnitedUKIR is the Zoom United Business plan for UK/Ireland.
	ZoomUnitedBizUnitedUKIR PlanUnitedType = "64"
	// ZoomUnitedBizUnitedAUNZ is the Zoom United Business plan for Australia/New Zealand.
	ZoomUnitedBizUnitedAUNZ PlanUnitedType = "128"
	// ZoomUnitedBizUnitedGlobalSelect is the Zoom United Business Global Select plan.
	ZoomUnitedBizUnitedGlobalSelect PlanUnitedType = "256"
	// ZoomUnitedBizUnitedZoomPhonePro is the Zoom United Business + Zoom Phone Pro bundle.
	ZoomUnitedBizUnitedZoomPhonePro PlanUnitedType = "512"
	// ZoomUnitedEntUnitedUSCA is the Zoom United Enterprise plan for US/Canada.
	ZoomUnitedEntUnitedUSCA PlanUnitedType = "1024"
	// ZoomUnitedEntUnitedUKIR is the Zoom United Enterprise plan for UK/Ireland.
	ZoomUnitedEntUnitedUKIR PlanUnitedType = "2048"
	// ZoomUnitedEntUnitedAUNZ is the Zoom United Enterprise plan for Australia/New Zealand.
	ZoomUnitedEntUnitedAUNZ PlanUnitedType = "4096"
	// ZoomUnitedEntUnitedGlobalSelect is the Zoom United Enterprise Global Select plan.
	ZoomUnitedEntUnitedGlobalSelect PlanUnitedType = "8192"
	// ZoomUnitedEntUnitedZoomPhonePro is the Zoom United Enterprise + Zoom Phone Pro bundle.
	ZoomUnitedEntUnitedZoomPhonePro PlanUnitedType = "16384"
	// ZoomUnitedProUnitedJP is the Zoom United Pro plan for Japan.
	ZoomUnitedProUnitedJP PlanUnitedType = "32768"
	// ZoomUnitedBizUnitedJP is the Zoom United Business plan for Japan.
	ZoomUnitedBizUnitedJP PlanUnitedType = "65536"
	// ZoomUnitedEntUnitedJP is the Zoom United Enterprise plan for Japan.
	ZoomUnitedEntUnitedJP PlanUnitedType = "131072"
)
