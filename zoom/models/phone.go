package models

import (
	"time"

	"github.com/TheSlowpes/go-zoom/zoom/enums"
)

type CallHistory struct {
	ID                string              `json:"id"`
	CallHistoryUUID   string              `json:"call_history_uuid"`
	CallID            string              `json:"call_id"`
	ConnectType       enums.ConnectType   `json:"connect_type"`
	CallType          string              `json:"call_type"`
	Direction         enums.Direction     `json:"direction"`
	International     bool                `json:"international"`
	HideCallerID      bool                `json:"hide_caller_id"`
	EndToEnd          bool                `json:"end_to_end"`
	CallerExtID       string              `json:"caller_ext_id"`
	CallerName        string              `json:"caller_name"`
	CallerDidNumber   string              `json:"caller_did_number"`
	CallerExtNumber   string              `json:"caller_ext_number"`
	CallerEmail       string              `json:"caller_email"`
	CallerExtType     enums.ExtensionType `json:"caller_ext_type"`
	CalleeExtID       string              `json:"callee_ext_id"`
	CalleeName        string              `json:"callee_name"`
	CalleeEmail       string              `json:"callee_email"`
	CalleeDidNumber   string              `json:"callee_did_number"`
	CalleeExtNumber   string              `json:"callee_ext_number"`
	CalleeExtType     enums.ExtensionType `json:"callee_ext_type"`
	Department        string              `json:"department"`
	CostCenter        string              `json:"cost_center"`
	SiteID            string              `json:"site_id"`
	GroupID           string              `json:"group_id"`
	SiteName          string              `json:"site_name"`
	StartTime         time.Time           `json:"start_time"`
	AnswerTime        time.Time           `json:"answer_time"`
	EndTime           time.Time           `json:"end_time"`
	CallPath          []CallPath          `json:"call_path"`
	CallerAccountCode string              `json:"caller_account_code"`
	CalleeAccountCode string              `json:"callee_account_code"`
	CallElements      []CallElement       `json:"call_elements"`
}

type CallPath struct {
	ID                   string              `json:"id"`
	CallID               string              `json:"call_id"`
	ConnectType          enums.ConnectType   `json:"connect_type"`
	CallType             enums.CallType      `json:"call_type"`
	Direction            enums.Direction     `json:"direction"`
	HideCallerID         bool                `json:"hide_caller_id"`
	EndToEnd             bool                `json:"end_to_end"`
	CallerExtID          string              `json:"caller_ext_id"`
	CallerName           string              `json:"caller_name"`
	CallerEmail          string              `json:"caller_email"`
	CallerDidNumber      string              `json:"caller_did_number"`
	CallerExtNumber      string              `json:"caller_ext_number"`
	CallerExtType        enums.ExtensionType `json:"caller_ext_type"`
	CallerNumberType     enums.NumberType    `json:"caller_number_type"`
	CallerDeviceType     string              `json:"caller_device_type"`
	CallerCountryIsoCode string              `json:"caller_country_iso_code"`
	CallerCountryCode    string              `json:"caller_country_code"`
	CalleeExtID          string              `json:"callee_ext_id"`
	CalleeName           string              `json:"callee_name"`
	CalleeDidNumber      string              `json:"callee_did_number"`
	CalleeExtNumber      string              `json:"callee_ext_number"`
	CalleeEmail          string              `json:"callee_email"`
	CalleeExtType        enums.ExtensionType `json:"callee_ext_type"`
	CalleeNumberType     enums.NumberType    `json:"callee_number_type"`
	CalleeDeviceType     string              `json:"callee_device_type"`
	CalleeCountryIsoCode string              `json:"callee_country_iso_code"`
	CalleeCountryCode    string              `json:"callee_country_code"`
	ClientCode           string              `json:"client_code"`
	Department           string              `json:"department"`
	CostCenter           string              `json:"cost_center"`
	SiteID               string              `json:"site_id"`
	GroupID              string              `json:"group_id"`
	SiteName             string              `json:"site_name"`
	StartTime            time.Time           `json:"start_time"`
	AnswerTime           time.Time           `json:"answer_time"`
	EndTime              time.Time           `json:"end_time"`
	Event                string              `json:"event"`
	Result               enums.CallResult    `json:"result"`
	ResultReason         string              `json:"result_reason"`
	DevicePrivateIP      string              `json:"device_private_ip"`
	DevicePublicIP       string              `json:"device_public_ip"`
	OperatorExtNumber    string              `json:"operator_ext_number"`
	OperatorExtID        string              `json:"operator_ext_id"`
	OperatorExtType      string              `json:"operator_ext_type"`
	OperatorName         string              `json:"operator_name"`
	PressKey             string              `json:"press_key"`
	Segment              int                 `json:"segment"`
	Node                 int                 `json:"node"`
	IsNode               int                 `json:"is_node"`
	RecordingID          string              `json:"recording_id"`
	RecordingType        string              `json:"recording_type"`
	HoldTime             int                 `json:"hold_time"`
	WaitingTime          int                 `json:"waiting_time"`
	VoicemailID          string              `json:"voicemail_id"`
}

type CallElement struct {
	CallElementID        string    `json:"call_element_id"`
	CallID               string    `json:"call_id"`
	ConnectType          string    `json:"connect_type"`
	CallType             string    `json:"call_type"`
	Direction            string    `json:"direction"`
	HideCallerID         bool      `json:"hide_caller_id"`
	EndToEnd             bool      `json:"end_to_end"`
	CallerExtID          string    `json:"caller_ext_id"`
	CallerName           string    `json:"caller_name"`
	CallerEmail          string    `json:"caller_email"`
	CallerDidNumber      string    `json:"caller_did_number"`
	CallerExtNumber      string    `json:"caller_ext_number"`
	CallerExtType        string    `json:"caller_ext_type"`
	CallerNumberType     string    `json:"caller_number_type"`
	CallerDeviceType     string    `json:"caller_device_type"`
	CallerCountryIsoCode string    `json:"caller_country_iso_code"`
	CallerCountryCode    string    `json:"caller_country_code"`
	CalleeExtID          string    `json:"callee_ext_id"`
	CalleeName           string    `json:"callee_name"`
	CalleeDidNumber      string    `json:"callee_did_number"`
	CalleeExtNumber      string    `json:"callee_ext_number"`
	CalleeEmail          string    `json:"callee_email"`
	CalleeExtType        string    `json:"callee_ext_type"`
	CalleeNumberType     string    `json:"callee_number_type"`
	CalleeDeviceType     string    `json:"callee_device_type"`
	CalleeCountryIsoCode string    `json:"callee_country_iso_code"`
	CalleeCountryCode    string    `json:"callee_country_code"`
	ClientCode           string    `json:"client_code"`
	Department           string    `json:"department"`
	CostCenter           string    `json:"cost_center"`
	SiteID               string    `json:"site_id"`
	GroupID              string    `json:"group_id"`
	SiteName             string    `json:"site_name"`
	StartTime            time.Time `json:"start_time"`
	AnswerTime           time.Time `json:"answer_time"`
	EndTime              time.Time `json:"end_time"`
	Event                string    `json:"event"`
	Result               string    `json:"result"`
	ResultReason         string    `json:"result_reason"`
	DevicePrivateIP      string    `json:"device_private_ip"`
	DevicePublicIP       string    `json:"device_public_ip"`
	OperatorExtNumber    string    `json:"operator_ext_number"`
	OperatorExtID        string    `json:"operator_ext_id"`
	OperatorExtType      string    `json:"operator_ext_type"`
	OperatorName         string    `json:"operator_name"`
	PressKey             string    `json:"press_key"`
	Segment              int       `json:"segment"`
	Node                 int       `json:"node"`
	IsNode               int       `json:"is_node"`
	RecordingID          string    `json:"recording_id"`
	RecordingType        string    `json:"recording_type"`
	HoldTime             int       `json:"hold_time"`
	WaitingTime          int       `json:"waiting_time"`
	VoicemailID          string    `json:"voicemail_id"`
}

type AICallSummary struct {
	AICallSummaryID    string    `json:"ai_call_summary_id"`
	AccountID          string    `json:"account_id,omitempty"`
	CallID             string    `json:"call_id"`
	CallLogIDs         []string  `json:"call_log_ids,omitempty"`
	UserID             string    `json:"user_id"`
	CallSummaryRate    string    `json:"call_summary_rate,omitempty"`
	TranscriptLanguage string    `json:"transcript_language,omitempty"`
	CallSummary        string    `json:"call_summary,omitempty"`
	NextSteps          string    `json:"next_steps,omitempty"`
	DetailedSummary    string    `json:"detailed_summary,omitempty"`
	CreatedTime        time.Time `json:"created_time"`
	ModifiedTime       time.Time `json:"modified_time"`
	Edited             bool      `json:"edited"`
	Deleted            bool      `json:"deleted,omitempty"`
}

type PhoneAccountSettings struct {
	BillingAccount          BillingAccount          `json:"billing_account"`
	BYOC                    BYOC                    `json:"byoc"`
	Country                 Country                 `json:"country"`
	MultiplePartyConference MultiplePartyConference `json:"multiple_party_conference"`
	MultipleSites           MultipleSites           `json:"multiple_sites"`
	ShowDeviceIPForCallLog  ShowDeviceIPForCallLog  `json:"show_device_ip_for_call_log"`
}

type BillingAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BYOC struct {
	Enable bool `json:"enable"`
}

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type MultiplePartyConference struct {
	Enable bool `json:"enable"`
}

type MultipleSites struct {
	Enabled  bool `json:"enabled"`
	SiteCode bool `json:"site_code,omitempty"`
}

type ShowDeviceIPForCallLog struct {
	Enable bool `json:"enable"`
}

type CallRecording struct {
	AutoDeletePolicy      string         `json:"auto_delete_policy"`
	CallID                string         `json:"call_id"`
	CallLogID             string         `json:"call_log_id"`
	CallElementID         string         `json:"call_element_id"`
	CalleeName            string         `json:"callee_name"`
	CalleeNumber          string         `json:"callee_number"`
	CalleeNumberType      int            `json:"callee_number_type"`
	CallerName            string         `json:"caller_name"`
	CallerNumber          string         `json:"caller_number"`
	CallerNumberType      int            `json:"caller_number_type"`
	OutgoingBy            RecordingUser  `json:"outgoing_by"`
	AcceptedBy            RecordingUser  `json:"accepted_by"`
	DateTime              time.Time      `json:"date_time"`
	DisclaimerStatus      int            `json:"disclaimer_status"`
	Direction             string         `json:"direction"`
	DownloadURL           string         `json:"download_url"`
	Duration              int            `json:"duration"`
	EndTime               time.Time      `json:"end_time"`
	ID                    string         `json:"id"`
	MeetingUUID           string         `json:"meeting_uuid"`
	Owner                 RecordingOwner `json:"owner"`
	RecordingType         string         `json:"recording_type"`
	Site                  RecordingSite  `json:"site"`
	TranscriptDownloadURL string         `json:"transcript_download_url"`
	AutoDeleteEnable      bool           `json:"auto_delete_enable"`
	CallerAccountCode     string         `json:"caller_account_code"`
	CalleeAccountCode     string         `json:"callee_account_code"`
}

type RecordingUser struct {
	Name            string `json:"name"`
	ExtensionNumber string `json:"extension_number"`
}

type RecordingOwner struct {
	ExtensionNumber      int       `json:"extension_number"`
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Type                 string    `json:"type"`
	ExtensionStatus      string    `json:"extension_status"`
	ExtensionDeletedTime time.Time `json:"extension_deleted_time"`
}

type RecordingSite struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
