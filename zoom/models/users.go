package models

import (
	"time"

	"github.com/TheSlowpes/go-zoom/zoom/enums"
)

type User struct {
	ID                 string               `json:"id"`
	CreatedAt          time.Time            `json:"created_at"`
	Dept               string               `json:"dept"`
	Email              string               `json:"email"`
	FirstName          string               `json:"first_name"`
	LastClientVersion  string               `json:"last_client_version"`
	LastLoginTime      time.Time            `json:"last_login_time"`
	LastName           string               `json:"last_name"`
	Pmi                int64                `json:"pmi"`
	RoleName           string               `json:"role_name"`
	Timezone           string               `json:"timezone"`
	Type               enums.UserType       `json:"type"`
	UsePmi             bool                 `json:"use_pmi"`
	DisplayName        string               `json:"display_name"`
	AccountID          string               `json:"account_id"`
	AccountNumber      int                  `json:"account_number"`
	CmsUserID          string               `json:"cms_user_id"`
	Company            string               `json:"company"`
	UserCreatedAt      time.Time            `json:"user_created_at"`
	CustomAttributes   []CustomAttributes   `json:"custom_attributes"`
	EmployeeUniqueID   string               `json:"employee_unique_id"`
	GroupIds           []string             `json:"group_ids"`
	DivisionIds        []string             `json:"division_ids"`
	ImGroupIds         []string             `json:"im_group_ids"`
	Jid                string               `json:"jid"`
	JobTitle           string               `json:"job_title"`
	CostCenter         string               `json:"cost_center"`
	Language           string               `json:"language"`
	Location           string               `json:"location"`
	LoginTypes         []enums.LoginType    `json:"login_types"`
	Manager            string               `json:"manager"`
	PersonalMeetingURL string               `json:"personal_meeting_url"`
	PhoneNumbers       []PhoneNumber        `json:"phone_numbers"`
	PicURL             string               `json:"pic_url"`
	PlanUnitedType     enums.PlanUnitedType `json:"plan_united_type"`
	Pronouns           string               `json:"pronouns"`
	PronounsOption     int                  `json:"pronouns_option"`
	RoleID             string               `json:"role_id"`
	Status             enums.UserStatus     `json:"status"`
	VanityURL          string               `json:"vanity_url"`
	Verified           int                  `json:"verified"`
	Cluster            string               `json:"cluster"`
	ZoomOneType        enums.ZoomOneType    `json:"zoom_one_type,omitempty"`
	LicenseInfoList    []LicenseInfo        `json:"license_info_list"`
}

type CustomAttributes struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PhoneNumber struct {
	ID       string                 `json:"id"`
	Code     string                 `json:"code"`
	Country  string                 `json:"country"`
	Label    enums.PhoneNumberLabel `json:"label"`
	Number   string                 `json:"number"`
	Verified bool                   `json:"verified"`
}

type LicenseInfo struct {
	LicenseType    string              `json:"license_type"`
	LicenseOption  enums.LicenseOption `json:"license_option"`
	SubscriptionID string              `json:"subscription_id"`
}

type Feature struct {
	ZoomOneType enums.ZoomOneType `json:"zoom_one_type,omitempty"`
	ZoomPhone   bool              `json:"zoom_phone,omitempty"`
}
