package enums

type LoginType int

const (
	FacebookOAuth    LoginType = 0
	GoogleOAuth      LoginType = 1
	PhoneNumberAuth  LoginType = 11
	WeChatAuth       LoginType = 21
	AlipayAuth       LoginType = 23
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

type UserCreateAction string

const (
	Create     UserCreateAction = "create"
	AutoCreate UserCreateAction = "auto_create"
	CustCreate UserCreateAction = "cust_create"
	SSOCreate  UserCreateAction = "sso_create"
)

type UserDeleteAction string

const (
	Disassociate UserDeleteAction = "disassociate"
	Delete       UserDeleteAction = "delete"
)

type LicenseOption int

type UserType LicenseOption

const (
	Basic    UserType = 1
	Licensed UserType = 2
	// Unnassigned is deprecated; use Unassigned instead.
	Unassigned  UserType = 3
	Unnassigned UserType = Unassigned
)

type ZoomOneType LicenseOption

const (
	ZoomWorkplaceEnterprise               ZoomOneType = 4
	ZoomWorkplaceEnterprisePlus           ZoomOneType = 8
	ZoomWorkplaceBusinessPlusUSCA         ZoomOneType = 16
	ZoomWorkplaceBusinessPlusUKIR         ZoomOneType = 32
	ZoomWorkplaceBusinessPlusAUNZ         ZoomOneType = 64
	ZoomWorkplaceBusinessPlusJapan        ZoomOneType = 128
	ZoomWorkplaceBusinessPlusGlobalSelect ZoomOneType = 33554432
	ZoomWorkplaceEnterprisePremierUSA     ZoomOneType = 13417728
	ZoomWorkplaceEnterprisePremierAUNZ    ZoomOneType = 1073741824
	ZoomWorkplaceEnterprisePremierUKIR    ZoomOneType = 536870912
	ZoomWorkplaceEnterprisePremierJapan   ZoomOneType = 268435456
	ZoomWorkplaceProPlus                  ZoomOneType = 4398046511104
)

type PhoneNumberLabel string

const (
	MobileLabel PhoneNumberLabel = "Mobile"
	OfficeLabel PhoneNumberLabel = "Office"
	HomeLabel   PhoneNumberLabel = "Home"
	FaxLabel    PhoneNumberLabel = "Fax"
)

type PlanUnitedType string

const (
	ZoomUnitedProUnitedUSCA         PlanUnitedType = "1"
	ZoomUnitedProUnitedUKIR         PlanUnitedType = "2"
	ZoomUnitedProUnitedAUNZ         PlanUnitedType = "4"
	ZoomUnitedProUnitedGlobalSelect PlanUnitedType = "8"
	ZoomUnitedProZoomPhonePro       PlanUnitedType = "16"
	ZoomUnitedBizUnitedUSCA         PlanUnitedType = "32"
	ZoomUnitedBizUnitedUKIR         PlanUnitedType = "64"
	ZoomUnitedBizUnitedAUNZ         PlanUnitedType = "128"
	ZoomUnitedBizUnitedGlobalSelect PlanUnitedType = "256"
	ZoomUnitedBizUnitedZoomPhonePro PlanUnitedType = "512"
	ZoomUnitedEntUnitedUSCA         PlanUnitedType = "1024"
	ZoomUnitedEntUnitedUKIR         PlanUnitedType = "2048"
	ZoomUnitedEntUnitedAUNZ         PlanUnitedType = "4096"
	ZoomUnitedEntUnitedGlobalSelect PlanUnitedType = "8192"
	ZoomUnitedEntUnitedZoomPhonePro PlanUnitedType = "16384"
	ZoomUnitedProUnitedJP           PlanUnitedType = "32768"
	ZoomUnitedBizUnitedJP           PlanUnitedType = "65536"
	ZoomUnitedEntUnitedJP           PlanUnitedType = "131072"
)
