package zoom

type LoginType int

const (
	FacebookOAuth    LoginType = 0
	GoogleOAuth      LoginType = 1
	AppleOAuth       LoginType = 24
	MicrosoftOAuth   LoginType = 27
	MobileDevice     LoginType = 97
	RingCentralOAuth LoginType = 98
	APIUser          LoginType = 99
	ZoomWorkEmail    LoginType = 100
	SingleSignOn     LoginType = 101
)

type UserStatus string

const (
	ActiveUser   UserStatus = "active"
	InactiveUser UserStatus = "inactive"
	PendingUser  UserStatus = "pending"
)

type UserFields string

const (
	CustomAttributesField UserFields = "custom_attributes"
	HostKeyField          UserFields = "host_key"
)

type LicenseType string

const (
	ZoomWorkforceManagment   LicenseType = "zoom_workforce_management"
	ZoomComplianceManagement LicenseType = "zoom_compliance_management"
)

type UserPostAction string

const (
	Create     UserPostAction = "create"
	AutoCreate UserPostAction = "auto_create"
	CustCreate UserPostAction = "cust_create"
	SSOCreate  UserPostAction = "sso_create"
)

type UserType int

const (
	Basic       UserType = 1
	Licensed    UserType = 2
	Unnassigned UserType = 3
)

type ZoomOneType int

const (
	ZoomWorkplaceBusinessPlusUSCA         ZoomOneType = 16
	ZoomWorkplaceBusinessPlusUKIR         ZoomOneType = 32
	ZoomWorkplaceBusinessPlusAUNZ         ZoomOneType = 64
	ZoomWorkplaceBusinessPlusJapan        ZoomOneType = 128
	ZoomWorkplaceBusinessPlusGlobalSelect ZoomOneType = 33554432
	ZoomWorkplaceEnterprisePremierUSA     ZoomOneType = 13417728
	ZoomWorkplaceEnterprisePremierAUNZ    ZoomOneType = 1073741824
)
