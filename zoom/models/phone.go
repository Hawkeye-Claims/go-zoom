package models

import (
	"time"

	"github.com/TheSlowpes/go-zoom/zoom/enums"
)

// CallHistory represents a single call log record as returned by the Zoom Phone
// call history endpoints.
type CallHistory struct {
	// ID is the unique identifier for the call log entry.
	ID string `json:"id,omitempty"`
	// CallHistoryUUID is the UUID of the call history record.
	CallHistoryUUID string `json:"call_history_uuid,omitempty"`
	// CallID is the unique identifier for the call.
	CallID string `json:"call_id,omitempty"`
	// ConnectType indicates whether the call was internal or external.
	ConnectType enums.ConnectType `json:"connect_type,omitempty"`
	// CallType is the type of call.
	CallType string `json:"call_type,omitempty"`
	// Direction indicates whether the call was inbound or outbound.
	Direction enums.Direction `json:"direction,omitempty"`
	// International indicates whether the call was an international call.
	International bool `json:"international"`
	// HideCallerID indicates whether the caller's number was withheld.
	HideCallerID bool `json:"hide_caller_id"`
	// EndToEnd indicates whether the call used end-to-end encryption.
	EndToEnd bool `json:"end_to_end"`
	// CallerExtID is the extension ID of the caller.
	CallerExtID string `json:"caller_ext_id,omitempty"`
	// CallerName is the display name of the caller.
	CallerName string `json:"caller_name,omitempty"`
	// CallerDidNumber is the DID number of the caller.
	CallerDidNumber string `json:"caller_did_number,omitempty"`
	// CallerExtNumber is the extension number of the caller.
	CallerExtNumber string `json:"caller_ext_number,omitempty"`
	// CallerEmail is the email address of the caller.
	CallerEmail string `json:"caller_email,omitempty"`
	// CallerExtType is the extension type of the caller.
	CallerExtType enums.ExtensionType `json:"caller_ext_type,omitempty"`
	// CalleeExtID is the extension ID of the callee.
	CalleeExtID string `json:"callee_ext_id,omitempty"`
	// CalleeName is the display name of the callee.
	CalleeName string `json:"callee_name,omitempty"`
	// CalleeEmail is the email address of the callee.
	CalleeEmail string `json:"callee_email,omitempty"`
	// CalleeDidNumber is the DID number of the callee.
	CalleeDidNumber string `json:"callee_did_number,omitempty"`
	// CalleeExtNumber is the extension number of the callee.
	CalleeExtNumber string `json:"callee_ext_number,omitempty"`
	// CalleeExtType is the extension type of the callee.
	CalleeExtType enums.ExtensionType `json:"callee_ext_type,omitempty"`
	// Department is the department associated with this call log entry.
	Department string `json:"department,omitempty"`
	// CostCenter is the cost center associated with this call log entry.
	CostCenter string `json:"cost_center,omitempty"`
	// SiteID is the ID of the site associated with this call log entry.
	SiteID string `json:"site_id,omitempty"`
	// GroupID is the ID of the group associated with this call log entry.
	GroupID string `json:"group_id,omitempty"`
	// SiteName is the name of the site associated with this call log entry.
	SiteName string `json:"site_name,omitempty"`
	// StartTime is the time the call started.
	StartTime time.Time `json:"start_time"`
	// AnswerTime is the time the call was answered.
	AnswerTime time.Time `json:"answer_time"`
	// EndTime is the time the call ended.
	EndTime time.Time `json:"end_time"`
	// CallPath lists the routing hops the call traversed.
	CallPath []CallPath `json:"call_path"`
	// CallerAccountCode is the account code assigned to the caller.
	CallerAccountCode string `json:"caller_account_code,omitempty"`
	// CalleeAccountCode is the account code assigned to the callee.
	CalleeAccountCode string `json:"callee_account_code,omitempty"`
	// CallElements lists the raw call elements associated with the call.
	CallElements []CallElement `json:"call_elements"`
}

// CallPath represents a single routing hop within a phone call, capturing the
// caller and callee details at each segment of the call's path.
type CallPath struct {
	// ID is the unique identifier for this call path segment.
	ID string `json:"id,omitempty"`
	// CallID is the unique identifier for the call.
	CallID string `json:"call_id,omitempty"`
	// ConnectType indicates whether this hop was internal or external.
	ConnectType enums.ConnectType `json:"connect_type,omitempty"`
	// CallType classifies the call as general or emergency.
	CallType enums.CallType `json:"call_type,omitempty"`
	// Direction indicates whether this hop was inbound or outbound.
	Direction enums.Direction `json:"direction,omitempty"`
	// HideCallerID indicates whether the caller's number was withheld.
	HideCallerID bool `json:"hide_caller_id"`
	// EndToEnd indicates whether this hop used end-to-end encryption.
	EndToEnd bool `json:"end_to_end"`
	// CallerExtID is the extension ID of the caller for this hop.
	CallerExtID string `json:"caller_ext_id,omitempty"`
	// CallerName is the display name of the caller.
	CallerName string `json:"caller_name,omitempty"`
	// CallerEmail is the email address of the caller.
	CallerEmail string `json:"caller_email,omitempty"`
	// CallerDidNumber is the DID number of the caller.
	CallerDidNumber string `json:"caller_did_number,omitempty"`
	// CallerExtNumber is the extension number of the caller.
	CallerExtNumber string `json:"caller_ext_number,omitempty"`
	// CallerExtType is the extension type of the caller.
	CallerExtType enums.ExtensionType `json:"caller_ext_type,omitempty"`
	// CallerNumberType classifies the caller's number.
	CallerNumberType enums.NumberType `json:"caller_number_type,omitempty"`
	// CallerDeviceType is the device type used by the caller.
	CallerDeviceType string `json:"caller_device_type,omitempty"`
	// CallerCountryIsoCode is the ISO country code for the caller's number.
	CallerCountryIsoCode string `json:"caller_country_iso_code,omitempty"`
	// CallerCountryCode is the country calling code for the caller's number.
	CallerCountryCode string `json:"caller_country_code,omitempty"`
	// CalleeExtID is the extension ID of the callee for this hop.
	CalleeExtID string `json:"callee_ext_id,omitempty"`
	// CalleeName is the display name of the callee.
	CalleeName string `json:"callee_name,omitempty"`
	// CalleeDidNumber is the DID number of the callee.
	CalleeDidNumber string `json:"callee_did_number,omitempty"`
	// CalleeExtNumber is the extension number of the callee.
	CalleeExtNumber string `json:"callee_ext_number,omitempty"`
	// CalleeEmail is the email address of the callee.
	CalleeEmail string `json:"callee_email,omitempty"`
	// CalleeExtType is the extension type of the callee.
	CalleeExtType enums.ExtensionType `json:"callee_ext_type,omitempty"`
	// CalleeNumberType classifies the callee's number.
	CalleeNumberType enums.NumberType `json:"callee_number_type,omitempty"`
	// CalleeDeviceType is the device type used by the callee.
	CalleeDeviceType string `json:"callee_device_type,omitempty"`
	// CalleeCountryIsoCode is the ISO country code for the callee's number.
	CalleeCountryIsoCode string `json:"callee_country_iso_code,omitempty"`
	// CalleeCountryCode is the country calling code for the callee's number.
	CalleeCountryCode string `json:"callee_country_code,omitempty"`
	// ClientCode is an optional annotation code added to the call log.
	ClientCode string `json:"client_code,omitempty"`
	// Department is the department associated with this call path segment.
	Department string `json:"department,omitempty"`
	// CostCenter is the cost center associated with this call path segment.
	CostCenter string `json:"cost_center,omitempty"`
	// SiteID is the site ID associated with this call path segment.
	SiteID string `json:"site_id,omitempty"`
	// GroupID is the group ID associated with this call path segment.
	GroupID string `json:"group_id,omitempty"`
	// SiteName is the site name associated with this call path segment.
	SiteName string `json:"site_name,omitempty"`
	// StartTime is the time this path segment started.
	StartTime time.Time `json:"start_time"`
	// AnswerTime is the time this path segment was answered.
	AnswerTime time.Time `json:"answer_time"`
	// EndTime is the time this path segment ended.
	EndTime time.Time `json:"end_time"`
	// Event describes the routing event that triggered this hop (e.g. transfer).
	Event string `json:"event,omitempty"`
	// Result is the outcome of this call path segment.
	Result enums.CallResult `json:"result,omitempty"`
	// ResultReason provides additional detail about the call result.
	ResultReason string `json:"result_reason,omitempty"`
	// DevicePrivateIP is the private IP address of the device used for this hop.
	DevicePrivateIP string `json:"device_private_ip,omitempty"`
	// DevicePublicIP is the public IP address of the device used for this hop.
	DevicePublicIP string `json:"device_public_ip,omitempty"`
	// OperatorExtNumber is the extension number of the operator involved.
	OperatorExtNumber string `json:"operator_ext_number,omitempty"`
	// OperatorExtID is the extension ID of the operator involved.
	OperatorExtID string `json:"operator_ext_id,omitempty"`
	// OperatorExtType is the extension type of the operator involved.
	OperatorExtType string `json:"operator_ext_type,omitempty"`
	// OperatorName is the display name of the operator involved.
	OperatorName string `json:"operator_name,omitempty"`
	// PressKey is the IVR key the caller pressed during this hop.
	PressKey string `json:"press_key,omitempty"`
	// Segment is the segment number of this call path entry.
	Segment int `json:"segment"`
	// Node is the node number within the segment.
	Node int `json:"node"`
	// IsNode indicates whether this entry represents a node (1) or not (0).
	IsNode int `json:"is_node"`
	// RecordingID is the ID of the recording associated with this hop.
	RecordingID string `json:"recording_id,omitempty"`
	// RecordingType describes the type of recording (e.g. "cloud").
	RecordingType string `json:"recording_type,omitempty"`
	// HoldTime is the total hold time in seconds for this hop.
	HoldTime int `json:"hold_time"`
	// WaitingTime is the total waiting time in seconds for this hop.
	WaitingTime int `json:"waiting_time"`
	// VoicemailID is the ID of the voicemail left during this hop.
	VoicemailID string `json:"voicemail_id,omitempty"`
}

// CallElement is a raw call element record. It mirrors the fields of CallPath
// but uses plain string types for enum fields, providing looser typing for
// values that may not match the defined enum constants.
type CallElement struct {
	// CallElementID is the unique identifier for this call element.
	CallElementID string `json:"call_element_id,omitempty"`
	// CallID is the unique identifier for the call.
	CallID string `json:"call_id,omitempty"`
	// ConnectType indicates whether the call was internal or external.
	ConnectType string `json:"connect_type,omitempty"`
	// CallType is the type of the call.
	CallType string `json:"call_type,omitempty"`
	// Direction indicates whether the call was inbound or outbound.
	Direction string `json:"direction,omitempty"`
	// HideCallerID indicates whether the caller's number was withheld.
	HideCallerID bool `json:"hide_caller_id"`
	// EndToEnd indicates whether the call used end-to-end encryption.
	EndToEnd bool `json:"end_to_end"`
	// CallerExtID is the extension ID of the caller.
	CallerExtID string `json:"caller_ext_id,omitempty"`
	// CallerName is the display name of the caller.
	CallerName string `json:"caller_name,omitempty"`
	// CallerEmail is the email address of the caller.
	CallerEmail string `json:"caller_email,omitempty"`
	// CallerDidNumber is the DID number of the caller.
	CallerDidNumber string `json:"caller_did_number,omitempty"`
	// CallerExtNumber is the extension number of the caller.
	CallerExtNumber string `json:"caller_ext_number,omitempty"`
	// CallerExtType is the extension type of the caller.
	CallerExtType string `json:"caller_ext_type,omitempty"`
	// CallerNumberType is the number type for the caller.
	CallerNumberType string `json:"caller_number_type,omitempty"`
	// CallerDeviceType is the device type used by the caller.
	CallerDeviceType string `json:"caller_device_type,omitempty"`
	// CallerCountryIsoCode is the ISO country code for the caller's number.
	CallerCountryIsoCode string `json:"caller_country_iso_code,omitempty"`
	// CallerCountryCode is the calling code for the caller's country.
	CallerCountryCode string `json:"caller_country_code,omitempty"`
	// CalleeExtID is the extension ID of the callee.
	CalleeExtID string `json:"callee_ext_id,omitempty"`
	// CalleeName is the display name of the callee.
	CalleeName string `json:"callee_name,omitempty"`
	// CalleeDidNumber is the DID number of the callee.
	CalleeDidNumber string `json:"callee_did_number,omitempty"`
	// CalleeExtNumber is the extension number of the callee.
	CalleeExtNumber string `json:"callee_ext_number,omitempty"`
	// CalleeEmail is the email address of the callee.
	CalleeEmail string `json:"callee_email,omitempty"`
	// CalleeExtType is the extension type of the callee.
	CalleeExtType string `json:"callee_ext_type,omitempty"`
	// CalleeNumberType is the number type for the callee.
	CalleeNumberType string `json:"callee_number_type,omitempty"`
	// CalleeDeviceType is the device type used by the callee.
	CalleeDeviceType string `json:"callee_device_type,omitempty"`
	// CalleeCountryIsoCode is the ISO country code for the callee's number.
	CalleeCountryIsoCode string `json:"callee_country_iso_code,omitempty"`
	// CalleeCountryCode is the calling code for the callee's country.
	CalleeCountryCode string `json:"callee_country_code,omitempty"`
	// ClientCode is an optional annotation code added to the call log.
	ClientCode string `json:"client_code,omitempty"`
	// Department is the department associated with this call element.
	Department string `json:"department,omitempty"`
	// CostCenter is the cost center associated with this call element.
	CostCenter string `json:"cost_center,omitempty"`
	// SiteID is the site ID associated with this call element.
	SiteID string `json:"site_id,omitempty"`
	// GroupID is the group ID associated with this call element.
	GroupID string `json:"group_id,omitempty"`
	// SiteName is the site name associated with this call element.
	SiteName string `json:"site_name,omitempty"`
	// StartTime is the time this call element started.
	StartTime time.Time `json:"start_time"`
	// AnswerTime is the time this call element was answered.
	AnswerTime time.Time `json:"answer_time"`
	// EndTime is the time this call element ended.
	EndTime time.Time `json:"end_time"`
	// Event describes the routing event for this call element.
	Event string `json:"event,omitempty"`
	// Result is the outcome of this call element.
	Result string `json:"result,omitempty"`
	// ResultReason provides additional detail about the result.
	ResultReason string `json:"result_reason,omitempty"`
	// DevicePrivateIP is the private IP of the device for this element.
	DevicePrivateIP string `json:"device_private_ip,omitempty"`
	// DevicePublicIP is the public IP of the device for this element.
	DevicePublicIP string `json:"device_public_ip,omitempty"`
	// OperatorExtNumber is the extension number of the operator involved.
	OperatorExtNumber string `json:"operator_ext_number,omitempty"`
	// OperatorExtID is the extension ID of the operator involved.
	OperatorExtID string `json:"operator_ext_id,omitempty"`
	// OperatorExtType is the extension type of the operator involved.
	OperatorExtType string `json:"operator_ext_type,omitempty"`
	// OperatorName is the display name of the operator involved.
	OperatorName string `json:"operator_name,omitempty"`
	// PressKey is the IVR key pressed during this element.
	PressKey string `json:"press_key,omitempty"`
	// Segment is the segment number of this call element.
	Segment int `json:"segment"`
	// Node is the node number within the segment.
	Node int `json:"node"`
	// IsNode indicates whether this entry represents a node (1) or not (0).
	IsNode int `json:"is_node"`
	// RecordingID is the ID of the recording associated with this element.
	RecordingID string `json:"recording_id,omitempty"`
	// RecordingType describes the type of recording for this element.
	RecordingType string `json:"recording_type,omitempty"`
	// HoldTime is the total hold time in seconds.
	HoldTime int `json:"hold_time"`
	// WaitingTime is the total waiting time in seconds.
	WaitingTime int `json:"waiting_time"`
	// VoicemailID is the ID of the voicemail left during this element.
	VoicemailID string `json:"voicemail_id,omitempty"`
}

// AICallSummary represents an AI-generated summary for a Zoom Phone call.
type AICallSummary struct {
	// AICallSummaryID is the unique identifier for this summary.
	AICallSummaryID string `json:"ai_call_summary_id,omitempty"`
	// AccountID is the account the summary belongs to.
	AccountID string `json:"account_id,omitempty"`
	// CallID is the unique identifier of the summarised call.
	CallID string `json:"call_id,omitempty"`
	// CallLogIDs lists the call log IDs included in the summary.
	CallLogIDs []string `json:"call_log_ids,omitempty"`
	// UserID is the ID of the user the summary is associated with.
	UserID string `json:"user_id,omitempty"`
	// CallSummaryRate is a quality or confidence rating for the summary.
	CallSummaryRate string `json:"call_summary_rate,omitempty"`
	// TranscriptLanguage is the language of the call transcript.
	TranscriptLanguage string `json:"transcript_language,omitempty"`
	// CallSummary is the high-level summary of the call.
	CallSummary string `json:"call_summary,omitempty"`
	// NextSteps lists the action items identified from the call.
	NextSteps string `json:"next_steps,omitempty"`
	// DetailedSummary is a more detailed narrative summary of the call.
	DetailedSummary string `json:"detailed_summary,omitempty"`
	// CreatedTime is the time the summary was generated.
	CreatedTime time.Time `json:"created_time"`
	// ModifiedTime is the time the summary was last modified.
	ModifiedTime time.Time `json:"modified_time"`
	// Edited indicates whether a human has manually edited the summary.
	Edited bool `json:"edited"`
	// Deleted indicates whether the summary has been deleted.
	Deleted bool `json:"deleted,omitempty"`
}

// PhoneAccountSettings contains account-level Zoom Phone configuration.
type PhoneAccountSettings struct {
	// BillingAccount is the billing account associated with the phone service.
	BillingAccount BillingAccount `json:"billing_account"`
	// BYOC contains Bring Your Own Carrier settings.
	BYOC BYOC `json:"byoc"`
	// Country is the primary country for the phone account.
	Country Country `json:"country"`
	// MultiplePartyConference contains multi-party conference settings.
	MultiplePartyConference MultiplePartyConference `json:"multiple_party_conference"`
	// MultipleSites contains multiple sites configuration.
	MultipleSites MultipleSites `json:"multiple_sites"`
	// ShowDeviceIPForCallLog controls display of device IP in call logs.
	ShowDeviceIPForCallLog ShowDeviceIPForCallLog `json:"show_device_ip_for_call_log"`
}

// BillingAccount identifies the billing account linked to the phone service.
type BillingAccount struct {
	// ID is the unique identifier of the billing account.
	ID string `json:"id,omitempty"`
	// Name is the display name of the billing account.
	Name string `json:"name,omitempty"`
}

// BYOC holds Bring Your Own Carrier configuration.
type BYOC struct {
	// Enable indicates whether BYOC is enabled for the account.
	Enable bool `json:"enable"`
}

// Country represents a country configuration entry within phone settings.
type Country struct {
	// Code is the two-letter ISO 3166-1 alpha-2 country code.
	Code string `json:"code,omitempty"`
	// CountryCode is the international dialling code (e.g. "1" for US).
	CountryCode string `json:"country_code,omitempty"`
	// Name is the human-readable country name.
	Name string `json:"name,omitempty"`
}

// MultiplePartyConference holds the multi-party conference feature setting.
type MultiplePartyConference struct {
	// Enable indicates whether multi-party conference calls are enabled.
	Enable bool `json:"enable"`
}

// MultipleSites holds the multiple-sites feature configuration.
type MultipleSites struct {
	// Enabled indicates whether multiple sites are enabled for the account.
	Enabled bool `json:"enabled"`
	// SiteCode indicates whether site codes are enabled.
	SiteCode bool `json:"site_code,omitempty"`
}

// ShowDeviceIPForCallLog controls whether device IP addresses are shown in call
// logs.
type ShowDeviceIPForCallLog struct {
	// Enable indicates whether device IP addresses are shown in call logs.
	Enable bool `json:"enable"`
}

// CallRecording represents a Zoom Phone call recording.
type CallRecording struct {
	// AutoDeletePolicy describes the auto-deletion policy for the recording.
	AutoDeletePolicy string `json:"auto_delete_policy,omitempty"`
	// CallID is the unique identifier of the call that was recorded.
	CallID string `json:"call_id,omitempty"`
	// CallLogID is the call log ID associated with this recording.
	CallLogID string `json:"call_log_id,omitempty"`
	// CallElementID is the call element ID associated with this recording.
	CallElementID string `json:"call_element_id,omitempty"`
	// CalleeName is the display name of the callee.
	CalleeName string `json:"callee_name,omitempty"`
	// CalleeNumber is the phone number of the callee.
	CalleeNumber string `json:"callee_number,omitempty"`
	// CalleeNumberType is the number type of the callee.
	CalleeNumberType int `json:"callee_number_type"`
	// CallerName is the display name of the caller.
	CallerName string `json:"caller_name,omitempty"`
	// CallerNumber is the phone number of the caller.
	CallerNumber string `json:"caller_number,omitempty"`
	// CallerNumberType is the number type of the caller.
	CallerNumberType int `json:"caller_number_type"`
	// OutgoingBy identifies the participant who placed the call.
	OutgoingBy RecordingUser `json:"outgoing_by"`
	// AcceptedBy identifies the participant who answered the call.
	AcceptedBy RecordingUser `json:"accepted_by"`
	// DateTime is the date and time the recording was made.
	DateTime time.Time `json:"date_time"`
	// DisclaimerStatus indicates the recording disclaimer acceptance status.
	DisclaimerStatus int `json:"disclaimer_status"`
	// Direction indicates whether the call was inbound or outbound.
	Direction string `json:"direction,omitempty"`
	// DownloadURL is the URL from which the recording file can be downloaded.
	DownloadURL string `json:"download_url,omitempty"`
	// Duration is the duration of the recording in seconds.
	Duration int `json:"duration"`
	// EndTime is the time the recording ended.
	EndTime time.Time `json:"end_time"`
	// ID is the unique identifier of the recording.
	ID string `json:"id,omitempty"`
	// MeetingUUID is the UUID of the meeting, if the call was a Zoom meeting.
	MeetingUUID string `json:"meeting_uuid,omitempty"`
	// Owner identifies the Zoom Phone extension that owns the recording.
	Owner RecordingOwner `json:"owner"`
	// RecordingType describes the type of recording.
	RecordingType string `json:"recording_type,omitempty"`
	// Site is the site the recording is associated with.
	Site RecordingSite `json:"site"`
	// TranscriptDownloadURL is the URL from which the transcript can be downloaded.
	TranscriptDownloadURL string `json:"transcript_download_url,omitempty"`
	// AutoDeleteEnable indicates whether auto-deletion is enabled for this
	// recording.
	AutoDeleteEnable bool `json:"auto_delete_enable"`
	// CallerAccountCode is the account code for the caller.
	CallerAccountCode string `json:"caller_account_code,omitempty"`
	// CalleeAccountCode is the account code for the callee.
	CalleeAccountCode string `json:"callee_account_code,omitempty"`
}

// RecordingUser identifies a call participant in the context of a recording.
type RecordingUser struct {
	// Name is the display name of the participant.
	Name string `json:"name,omitempty"`
	// ExtensionNumber is the extension number of the participant.
	ExtensionNumber string `json:"extension_number,omitempty"`
}

// RecordingOwner describes the Zoom Phone extension that owns a recording.
type RecordingOwner struct {
	// ExtensionNumber is the extension number of the owner.
	ExtensionNumber int `json:"extension_number"`
	// ID is the unique identifier of the owner.
	ID string `json:"id,omitempty"`
	// Name is the display name of the owner.
	Name string `json:"name,omitempty"`
	// Type is the extension type of the owner.
	Type string `json:"type,omitempty"`
	// ExtensionStatus is the current status of the owner's extension.
	ExtensionStatus string `json:"extension_status,omitempty"`
	// ExtensionDeletedTime is the time the owner's extension was deleted, if
	// applicable.
	ExtensionDeletedTime time.Time `json:"extension_deleted_time"`
}

// RecordingSite identifies the Zoom Phone site associated with a recording.
type RecordingSite struct {
	// ID is the unique identifier of the site.
	ID string `json:"id,omitempty"`
	// Name is the display name of the site.
	Name string `json:"name,omitempty"`
}

// RecordingTranscript contains the full transcript of a phone call recording,
// including speaker-attributed timeline segments.
type RecordingTranscript struct {
	// Type identifies the transcript format type.
	Type string `json:"type,omitempty"`
	// Ver is the transcript format version.
	Ver string `json:"ver,omitempty"`
	// RecordingID is the unique identifier of the associated recording.
	RecordingID string `json:"recording_id,omitempty"`
	// MeetingID is the meeting ID if the call was a Zoom meeting.
	MeetingID string `json:"meeting_id,omitempty"`
	// AccountID is the account the transcript belongs to.
	AccountID string `json:"account_id,omitempty"`
	// HostID is the user ID of the meeting or call host.
	HostID string `json:"host_id,omitempty"`
	// RecordingStart is the time the recording started.
	RecordingStart time.Time `json:"recording_start"`
	// RecordingEnd is the time the recording ended.
	RecordingEnd time.Time `json:"recording_end"`
	// Timeline lists the ordered transcript segments with speaker attribution.
	Timeline []TranscriptTimeline `json:"timeline"`
}

// TranscriptTimeline represents a single spoken segment within a call
// transcript, attributed to one or more speakers.
type TranscriptTimeline struct {
	// Text is the normalised text of the spoken segment.
	Text string `json:"text,omitempty"`
	// RawText is the unprocessed text of the spoken segment.
	RawText string `json:"raw_text,omitempty"`
	// TS is the start timestamp of this segment.
	TS time.Time `json:"ts"`
	// EndTS is the end timestamp of this segment.
	EndTS time.Time `json:"end_ts"`
	// Users lists the speakers in this segment.
	Users []TranscriptUser `json:"users"`
	// UserID is the primary speaker's user ID.
	UserID string `json:"userId,omitempty"`
	// UserIDs lists all speaker user IDs in this segment.
	UserIDs []string `json:"userIds"`
	// ChannelMark identifies the audio channel for this segment.
	ChannelMark int `json:"channelMark"`
}

// TranscriptUser represents a speaker within a transcript timeline segment.
type TranscriptUser struct {
	// Username is the display name of the speaker.
	Username string `json:"username,omitempty"`
	// MultiplePeople indicates whether multiple people are speaking in this
	// segment.
	MultiplePeople bool `json:"multiple_people"`
	// UserID is the Zoom user ID of the speaker.
	UserID string `json:"user_id,omitempty"`
	// ZoomUserID is an alternative Zoom user identifier.
	ZoomUserID string `json:"zoom_userid,omitempty"`
	// AvatarURL is the URL of the speaker's profile picture.
	AvatarURL string `json:"avatar_url,omitempty"`
	// ClientType is the numeric client type identifier.
	ClientType int `json:"client_type"`
	// EmailAddress is the email address of the speaker.
	EmailAddress string `json:"email_address,omitempty"`
	// ChannelMark identifies the audio channel for this speaker.
	ChannelMark int `json:"channel_mark"`
	// Pronoun is the speaker's preferred pronouns.
	Pronoun string `json:"pronoun,omitempty"`
	// EnableSpeakerDiarization indicates whether speaker diarization is enabled.
	EnableSpeakerDiarization string `json:"enable_speaker_diarization,omitempty"`
	// TS is a numeric timestamp offset for this speaker entry.
	TS int `json:"ts"`
	// IsAutoGenerated indicates whether this speaker entry was auto-generated.
	IsAutoGenerated bool `json:"is_auto_generated"`
}

// PhoneUser represents a Zoom Phone user as returned by the phone users
// endpoint.
type PhoneUser struct {
	// ActivationStatus is the activation state of the phone user.
	ActivationStatus string `json:"activation_status,omitempty"`
	// CallingPlans lists the calling plans assigned to the user.
	CallingPlans []CallingPlan `json:"calling_plans"`
	// CostCenter is the cost center the user is assigned to.
	CostCenter string `json:"cost_center,omitempty"`
	// Department is the department the user is assigned to.
	Department string `json:"department,omitempty"`
	// Email is the user's Zoom account email address.
	Email string `json:"email,omitempty"`
	// EmergencyAddress is the user's registered emergency (E911) address.
	EmergencyAddress EmergencyAddress `json:"emergency_address"`
	// ExtensionID is the unique identifier of the user's extension.
	ExtensionID string `json:"extension_id,omitempty"`
	// ExtensionNumber is the user's phone extension number.
	ExtensionNumber int `json:"extension_number"`
	// ID is the unique identifier of the phone user.
	ID string `json:"id,omitempty"`
	// PhoneNumbers lists the phone numbers assigned to the user.
	PhoneNumbers []PhoneNumber `json:"phone_numbers"`
	// PhoneUserID is an alternative identifier for the phone user.
	PhoneUserID string `json:"phone_user_id,omitempty"`
	// Policy contains the per-user phone policy configuration.
	Policy PhoneUserPolicy `json:"policy"`
	// SiteAdmin indicates whether the user is an administrator for their site.
	SiteAdmin bool `json:"site_admin"`
	// SiteID is the ID of the site the user belongs to.
	SiteID string `json:"site_id,omitempty"`
	// Site contains the site ID and name the user belongs to.
	Site PhoneUserSite `json:"site"`
	// Status is the current status of the phone user.
	Status string `json:"status,omitempty"`
}

// PhoneUserSite identifies the Zoom Phone site a user belongs to.
type PhoneUserSite struct {
	// ID is the unique identifier of the site.
	ID string `json:"id,omitempty"`
	// Name is the display name of the site.
	Name string `json:"name,omitempty"`
}

// EmergencyAddress is a physical address registered for emergency (E911)
// services for a Zoom Phone user.
type EmergencyAddress struct {
	// AddressLine1 is the first line of the street address.
	AddressLine1 string `json:"address_line1,omitempty"`
	// AddressLine2 is the second line of the street address (suite, unit, etc.).
	AddressLine2 string `json:"address_line2,omitempty"`
	// City is the city of the address.
	City string `json:"city,omitempty"`
	// Country is the ISO 3166-1 alpha-2 country code.
	Country string `json:"country,omitempty"`
	// ID is the unique identifier of this emergency address record.
	ID string `json:"id,omitempty"`
	// StateCode is the state or province code.
	StateCode string `json:"state_code,omitempty"`
	// Zip is the postal code.
	Zip string `json:"zip,omitempty"`
}

// CallingPlan represents a Zoom Phone calling plan assigned to a user.
type CallingPlan struct {
	// Type is the numeric type identifier for the calling plan.
	Type int `json:"type"`
	// BillingAccountID is the billing account the plan is charged to.
	BillingAccountID string `json:"billing_account_id,omitempty"`
	// BillingAccountName is the display name of the billing account.
	BillingAccountName string `json:"billing_account_name,omitempty"`
	// BillingSubscriptionID is the billing subscription ID for this plan.
	BillingSubscriptionID string `json:"billing_subscription_id,omitempty"`
	// BillingSubscriptionName is the display name of the billing subscription.
	BillingSubscriptionName string `json:"billing_subscription_name,omitempty"`
}

// PhoneUserPolicy holds the complete phone policy configuration for an
// individual Zoom Phone user.
type PhoneUserPolicy struct {
	// AdHocCallRecording configures manual (on-demand) call recording.
	AdHocCallRecording AdHocCallRecording `json:"ad_hoc_call_recording"`
	// AdHocCallRecordingAccessMembers lists members who can access ad-hoc
	// recordings.
	AdHocCallRecordingAccessMembers []PolicyAccessMember `json:"ad_hoc_call_recording_access_members"`
	// AutoCallRecording configures automatic call recording.
	AutoCallRecording AutoCallRecordingPolicy `json:"auto_call_recording"`
	// AutoCallRecordingAccessMembers lists members who can access auto-recordings.
	AutoCallRecordingAccessMembers []*PolicyAccessMember `json:"auto_call_recording_access_members"`
	// CallOverflow configures call overflow behaviour.
	CallOverflow CallOverflowPolicy `json:"call_overflow"`
	// CallPark configures call park settings.
	CallPark CallParkPolicy `json:"call_park"`
	// CallTransferring configures call transfer restrictions.
	CallTransferring CallTransferPolicy `json:"call_transferring"`
	// Delegation indicates whether call delegation is enabled.
	Delegation bool `json:"delegation"`
	// ElevateToMeeting allows elevating a phone call to a Zoom meeting.
	ElevateToMeeting bool `json:"elevate_to_meeting"`
	// EmergencyAddressManagement configures emergency address management.
	EmergencyAddressManagement AddressManagementPolicy `json:"emergency_address_management"`
	// EmergencyCallsToPsap controls whether emergency calls are routed to the
	// public safety answering point.
	EmergencyCallsToPsap bool `json:"emergency_calls_to_psap"`
	// CallHandlingForwardingToOtherUsers configures forwarding restrictions.
	CallHandlingForwardingToOtherUsers CallForwardPolicy `json:"call_handling_forwarding_to_other_users"`
	// HandOffToRoom configures hand-off to a Zoom Room policy state.
	HandOffToRoom *PolicyState `json:"hand_off_to_room,omitempty"`
	// InternationalCalling enables making international calls.
	InternationalCalling bool `json:"international_calling"`
	// MobileSwitchToCarrier configures mobile carrier switching policy.
	MobileSwitchToCarrier *PolicyState `json:"mobile_switch_to_carrier,omitempty"`
	// SelectOutboundCallerID configures outbound caller ID selection.
	SelectOutboundCallerID OutboundCallerIDPolicy `json:"select_outbound_caller_id"`
	// SMS configures SMS and MMS policy.
	SMS SMSPolicy `json:"sms"`
	// Voicemail configures voicemail policy.
	Voicemail VoicemailPolicy `json:"voicemail"`
	// VoicemailAccessMembers lists members who can access the user's voicemail.
	VoicemailAccessMembers []*PolicyAccessMember `json:"voicemail_access_members"`
	// ZoomPhoneOnMobile configures Zoom Phone on mobile device policy.
	ZoomPhoneOnMobile ZoomPhoneOnMobilePolicy `json:"zoom_phone_on_mobile"`
	// PersonalAudioLibrary configures personal audio library policy.
	PersonalAudioLibrary PersonalAudioLibraryPolicy `json:"personal_audio_library"`
	// VoicemailTranscription configures voicemail-to-text transcription policy.
	VoicemailTranscription *PolicyState `json:"voicemail_transcription,omitempty"`
	// VoicemailNotificationByEmail configures voicemail email notification policy.
	VoicemailNotificationByEmail VoicemailNotificationByEmailPolicy `json:"voicemail_notification_by_email"`
	// SharedVoicemailNotificationByEmail configures shared voicemail email
	// notification policy.
	SharedVoicemailNotificationByEmail *PolicyState `json:"shared_voicemail_notification_by_email,omitempty"`
	// CheckVoicemailsOverPhone configures phone-based voicemail check policy.
	CheckVoicemailsOverPhone *PolicyState `json:"check_voicemails_over_phone,omitempty"`
	// AudioIntercom configures the audio intercom feature policy.
	AudioIntercom *PolicyState `json:"audio_intercom,omitempty"`
	// PeerToPeerMedia configures peer-to-peer media policy.
	PeerToPeerMedia *PolicyState `json:"peer_to_peer_media,omitempty"`
	// E2EEncryption configures end-to-end encryption policy.
	E2EEncryption *PolicyState `json:"e2e_encryption,omitempty"`
	// OutboundCalling configures outbound calling policy.
	OutboundCalling *PolicyState `json:"outbound_calling,omitempty"`
	// OutboundSms configures outbound SMS policy.
	OutboundSms *PolicyState `json:"outbound_sms,omitempty"`
	// AllowEndUserEditCallHandling allows users to edit their own call handling.
	AllowEndUserEditCallHandling *PolicyState `json:"allow_end_user_edit_call_handling,omitempty"`
	// VoicemailIntentBasedPrioritization configures intent-based voicemail
	// prioritisation.
	VoicemailIntentBasedPrioritization *PolicyState `json:"voicemail_intent_based_prioritization,omitempty"`
	// VoicemailTasks configures voicemail task extraction policy.
	VoicemailTasks *PolicyState `json:"voicemail_tasks,omitempty"`
	// ZoomPhoneOnDesktop configures Zoom Phone on desktop device policy.
	ZoomPhoneOnDesktop ZoomPhoneOnDesktopPolicy `json:"zoom_phone_on_desktop"`
}

// PhoneUserSettings contains per-user Zoom Phone profile settings.
type PhoneUserSettings struct {
	// AreaCode is the user's default area code.
	AreaCode string `json:"area_code,omitempty"`
	// AudioPromptLanguage is the language used for audio prompts.
	AudioPromptLanguage string `json:"audio_prompt_language,omitempty"`
	// CompanyNumber is the user's company phone number.
	CompanyNumber string `json:"company_number,omitempty"`
	// Country is the user's primary country for phone settings.
	Country Country `json:"country"`
	// Delegation contains the user's call delegation configuration.
	Delegation DelegationSettings `json:"delegation"`
	// DeskPhone contains the user's desk phone configuration.
	DeskPhone DeskPhoneSettings `json:"desk_phone"`
	// ExtensionNumber is the user's phone extension number.
	ExtensionNumber int `json:"extension_number"`
	// MusicOnHoldID is the ID of the music-on-hold audio file.
	MusicOnHoldID string `json:"music_on_hold_id,omitempty"`
	// OutboundCaller is the user's current default outbound caller ID.
	OutboundCaller OutboundCallerIDSetting `json:"outbound_caller"`
	// OutboundCallerIds lists all available outbound caller ID options for the
	// user.
	OutboundCallerIds []OutboundCallerIDSetting `json:"outbound_caller_ids"`
	// Status is the current status of the phone user.
	Status string `json:"status,omitempty"`
	// VoiceMail lists the members who have access to the user's voicemail.
	VoiceMail []PolicyAccessMember `json:"voice_mail"`
	// Intercom contains the user's audio intercom configuration.
	Intercom IntercomSettings `json:"intercom"`
	// AutoCallRecordingAccessMembers lists members with access to auto-recordings.
	AutoCallRecordingAccessMembers []PolicyAccessMember `json:"auto_call_recording_access_members"`
	// AdHocCallRecordingAccessMembers lists members with access to ad-hoc
	// recordings.
	AdHocCallRecordingAccessMembers []PolicyAccessMember `json:"ad_hoc_call_recording_access_members"`
	// SharedLinesCallSetting contains the user's shared line configuration.
	SharedLinesCallSetting SharedLineCallSetting `json:"shared_lines_call_setting"`
}

// DeskPhoneSettings contains the physical desk phone configuration for a user.
type DeskPhoneSettings struct {
	// Devices lists the desk phone devices assigned to the user.
	Devices []Device `json:"devices"`
	// KeysPositions configures the key layout on the desk phone.
	KeysPositions KeysPositions `json:"keys_positions"`
	// PhoneScreenLock indicates whether the phone screen lock is enabled.
	PhoneScreenLock bool `json:"phone_screen_lock"`
	// PinCode is the screen lock PIN code.
	PinCode string `json:"pin_code,omitempty"`
}

// KeysPositions configures the key layout on a desk phone.
type KeysPositions struct {
	// PrimaryNumber is the phone number assigned to the primary key position.
	PrimaryNumber string `json:"primary_number,omitempty"`
}

// DelegationSettings configures call delegation for a user, allowing assistants
// to make and receive calls on the user's behalf.
type DelegationSettings struct {
	// Assistants lists the extensions that are authorised as call delegates.
	Assistants []PhoneUserExtension `json:"assistants"`
	// Privacy enables privacy mode for delegation.
	Privacy bool `json:"privacy"`
	// Privileges is the bitmask of privileges granted to assistants.
	Privileges int `json:"privileges"`
	// Locked indicates whether the delegation settings are locked by an admin.
	Locked bool `json:"locked"`
}

// OutboundCallerIDSetting represents a single outbound caller ID option
// available to a user.
type OutboundCallerIDSetting struct {
	// IsDefault indicates whether this is the user's current default caller ID.
	IsDefault bool `json:"is_default"`
	// Name is the display name for this caller ID option.
	Name string `json:"name,omitempty"`
	// Number is the phone number for this caller ID option.
	Number string `json:"number,omitempty"`
}

// DevicePolicy holds call-control and hot-desking policy states for a device.
type DevicePolicy struct {
	// CallControl is the policy state for call control on the device.
	CallControl PolicyStatus `json:"call_control"`
	// HotDesking is the policy state for hot-desking on the device.
	HotDesking PolicyStatus `json:"hot_desking"`
}

// PolicyStatus holds the status string for a device-level policy.
type PolicyStatus struct {
	// Status is the current status value for the policy.
	Status string `json:"status,omitempty"`
}

// AudioIntercomSettings configures a single audio intercom entry for a user.
type AudioIntercomSettings struct {
	*PhoneUserExtension
	// Status is the status of the intercom connection.
	Status string `json:"status,omitempty"`
	// DeviceID is the ID of the device used for the intercom.
	DeviceID string `json:"device_id,omitempty"`
	// DeviceStatus is the current status of the intercom device.
	DeviceStatus string `json:"device_status,omitempty"`
}

// IntercomSettings contains the audio intercom configuration for a user.
type IntercomSettings struct {
	// AudioIntercoms lists the configured audio intercom entries.
	AudioIntercoms []AudioIntercomSettings `json:"audio_intercoms"`
	// Device is the primary device used for audio intercom.
	Device Device `json:"device"`
}

// Device represents a physical or software phone device assigned to a user.
type Device struct {
	// DeviceType is the type identifier of the device.
	DeviceType string `json:"device_type,omitempty"`
	// DisplayName is the display name of the device.
	DisplayName string `json:"display_name,omitempty"`
	// ID is the unique identifier of the device.
	ID string `json:"id,omitempty"`
	// Policy contains the call-control and hot-desking policy for the device.
	Policy DevicePolicy `json:"policy"`
	// Name is the name of the device model.
	Name string `json:"name,omitempty"`
	// Status is the current online/offline status of the device.
	Status string `json:"status,omitempty"`
	// MacAddress is the MAC address of the physical device.
	MacAddress string `json:"mac_address,omitempty"`
	// PrivateIP is the private IP address of the device.
	PrivateIP string `json:"private_ip,omitempty"`
	// PublicIP is the public IP address of the device.
	PublicIP string `json:"public_ip,omitempty"`
}

// SharedLineCallSetting contains the shared line appearances and groups
// configuration for a user.
type SharedLineCallSetting struct {
	// SharedLineAppearances configures shared line appearances with executives.
	SharedLineAppearances SharedLineAppearances `json:"shared_line_appearances"`
	// SharedLineGroups configures shared line group memberships.
	SharedLineGroups SharedLineGroups `json:"shared_line_groups"`
}

// SharedLineAppearances configures the shared line appearance feature, where
// one user's calls appear on another user's phone.
type SharedLineAppearances struct {
	// Executives lists the executives whose lines appear on this user's phone.
	Executives []Executive `json:"executives"`
}

// Executive represents an executive whose line appears on an assistant's phone
// via the shared line appearance feature.
type Executive struct {
	// UserID is the unique identifier of the executive.
	UserID string `json:"user_id,omitempty"`
	// DisplayName is the display name of the executive.
	DisplayName string `json:"display_name,omitempty"`
	// ReceiveCalls indicates whether the assistant receives calls for this
	// executive.
	ReceiveCalls bool `json:"receive_calls"`
	// AllowOptOut allows the assistant to opt out of receiving calls for this
	// executive.
	AllowOptOut bool `json:"allow_opt_out"`
}

// SharedLineGroups configures shared line group membership for a user.
type SharedLineGroups struct {
	// ReceiveCalls indicates whether the user receives calls for all shared line
	// groups.
	ReceiveCalls bool `json:"receive_calls"`
	// SharedLineGroup lists the individual shared line group memberships.
	SharedLineGroup []SharedLineGroupSetting `json:"shared_line_group"`
}

// SharedLineGroupSetting describes a single shared line group membership.
type SharedLineGroupSetting struct {
	// SlgID is the unique identifier of the shared line group.
	SlgID string `json:"slg_id,omitempty"`
	// DisplayName is the display name of the shared line group.
	DisplayName string `json:"display_name,omitempty"`
	// ReceiveCalls indicates whether the user receives calls for this group.
	ReceiveCalls bool `json:"receive_calls"`
}

// AdHocCallRecording configures the manual (on-demand) call recording policy
// for a user.
type AdHocCallRecording struct {
	*PolicyState
	// RecordingStartPrompt plays an audio prompt when recording starts.
	RecordingStartPrompt bool `json:"recording_start_prompt"`
	// RecordingTranscription enables transcription of ad-hoc recordings.
	RecordingTranscription bool `json:"recording_transcription"`
	// PlayRecordingBeepTone configures the beep tone played during recording.
	PlayRecordingBeepTone PlayRecordingBeepTone `json:"play_recording_beep_tone"`
}

// PlayRecordingBeepTone configures the beep tone that plays to indicate a call
// is being recorded.
type PlayRecordingBeepTone struct {
	*PolicyState
	// PlayBeepVolume is the volume level of the beep tone.
	PlayBeepVolume int `json:"play_beep_volume"`
	// PlayBeepTimeInterval is the interval in seconds between beep tones.
	PlayBeepTimeInterval int `json:"play_beep_time_interval"`
	// PlayBeepMember specifies which call participants hear the beep tone.
	PlayBeepMember string `json:"play_beep_member,omitempty"`
}

// PolicyState represents the common enable/lock/modified state shared by many
// Zoom Phone policy types.
type PolicyState struct {
	// Enable indicates whether the policy feature is currently active.
	Enable bool `json:"enable"`
	// Locked indicates whether the policy is locked and cannot be changed by the
	// user.
	Locked bool `json:"locked"`
	// LockedBy identifies the admin or level that locked the policy.
	LockedBy string `json:"locked_by,omitempty"`
	// Modified indicates whether the policy has been changed from its default
	// value.
	Modified bool `json:"modified"`
}

// AutoCallRecordingPolicy configures automatic call recording for a user.
type AutoCallRecordingPolicy struct {
	*PolicyState
	// AllowStopResumeRecording allows the user to stop and resume recording.
	AllowStopResumeRecording bool `json:"allow_stop_resume_recording"`
	// DisconnectOnRecordingFailure disconnects the call if recording fails.
	DisconnectOnRecordingFailure bool `json:"disconnect_on_recording_failure"`
	// RecordingCalls specifies which call types are automatically recorded
	// ("inbound", "outbound", "both").
	RecordingCalls string `json:"recording_calls,omitempty"`
	// RecordingTranscription enables transcription of auto-recorded calls.
	RecordingTranscription bool `json:"recording_transcription"`
	// PlayRecordingBeepTone configures the beep tone played during recording.
	PlayRecordingBeepTone PlayRecordingBeepTone `json:"play_recording_beep_tone"`
	// InboundAudioNotification configures the consent prompt played to inbound
	// callers.
	InboundAudioNotification *RecordingNotificationPolicy `json:"inbound_audio_notification,omitempty"`
	// OutboundAudioNotification configures the consent prompt played to outbound
	// callees.
	OutboundAudioNotification *RecordingNotificationPolicy `json:"outbound_audio_notification,omitempty"`
}

// PolicyAccessMember describes a user or extension that has been granted access
// to shared recordings or voicemail.
type PolicyAccessMember struct {
	// AccessUserID is the user ID of the member being granted access.
	AccessUserID string `json:"access_user_id,omitempty"`
	// AccessUserType is the type of user being granted access.
	AccessUserType string `json:"access_user_type,omitempty"`
	// AllowDelete indicates whether the member is permitted to delete records.
	AllowDelete bool `json:"allow_delete"`
	// Delete is the current delete permission state.
	Delete bool `json:"delete"`
	// AllowDownload indicates whether the member is permitted to download records.
	AllowDownload bool `json:"allow_download"`
	// Download is the current download permission state.
	Download bool `json:"download"`
	// AllowSharing indicates whether the member is permitted to share records.
	AllowSharing bool `json:"allow_sharing"`
	// SharedID is the ID of the share, if the record has been shared.
	SharedID string `json:"shared_id,omitempty"`
}

// RecordingNotificationPolicy configures audio consent prompts played to call
// participants when automatic recording begins.
type RecordingNotificationPolicy struct {
	// RecordingStartPrompt plays an audio notification when recording starts.
	RecordingStartPrompt bool `json:"recording_start_prompt"`
	// RecordingExplicitConsent requires the participant to press a key to consent
	// before the call continues.
	RecordingExplicitConsent bool `json:"recording_explicit_consent"`
}

// CallOverflowPolicy configures call overflow behaviour when a user is
// unavailable.
type CallOverflowPolicy struct {
	*PolicyState
	// CallOverflowType specifies the forwarding restriction type for overflow.
	CallOverflowType enums.CallForwardingType `json:"call_overflow_type,omitempty"`
}

// CallParkPolicy configures call park settings for a user.
type CallParkPolicy struct {
	*PolicyState
	// CallNotPickedUpAction is the action taken when a parked call is not
	// retrieved within the expiration period.
	CallNotPickedUpAction enums.CallNotPickedUpAction `json:"call_not_picked_up_action,omitempty"`
	// ExpirationPeriod is the number of seconds before an unanswered parked call
	// triggers the CallNotPickedUpAction.
	ExpirationPeriod int `json:"expiration_period"`
	// ForwardTo is the extension to forward the call to when the expiration period
	// elapses and the action is ForwardToExtension.
	ForwardTo PhoneUserExtension `json:"forward_to"`
}

// PhoneUserExtension is a reference to a Zoom Phone extension, used in policy
// and settings contexts.
type PhoneUserExtension struct {
	// DisplayName is the display name of the extension.
	DisplayName string `json:"display_name,omitempty"`
	// ExtensionID is the unique identifier of the extension.
	ExtensionID string `json:"extension_id,omitempty"`
	// ExtensionNumber is the numeric extension number. The type is any because the
	// API may return it as either a string or an integer.
	ExtensionNumber any `json:"extension_number,omitempty"`
	// ExtensionType is the type of the extension.
	ExtensionType string `json:"extension_type,omitempty"`
	// ID is an alternative unique identifier for the extension.
	ID string `json:"id,omitempty"`
}

// CallTransferPolicy configures call transfer restrictions for a user.
type CallTransferPolicy struct {
	*PolicyState
	// CallTransferringType specifies the transfer destination restriction level.
	CallTransferringType enums.CallTransferType `json:"call_transferring_type,omitempty"`
}

// AddressManagementPolicy configures emergency address management for a user.
type AddressManagementPolicy struct {
	*PolicyState
	// PromptDefaultAddress prompts the user to confirm their emergency address
	// when making a call.
	PromptDefaultAddress bool `json:"prompt_default_address"`
}

// CallForwardPolicy configures call forwarding destination restrictions for a
// user.
type CallForwardPolicy struct {
	*PolicyState
	// CallForwardingType specifies the forwarding destination restriction level.
	CallForwardingType enums.CallForwardingType `json:"call_forwarding_type,omitempty"`
}

// OutboundCallerIDPolicy configures outbound caller ID visibility options for a
// user.
type OutboundCallerIDPolicy struct {
	*PolicyState
	// AllowHideOutboundCallerID permits the user to hide their outbound caller ID.
	AllowHideOutboundCallerID bool `json:"allow_hide_outbound_caller_id"`
}

// SMSPolicy configures SMS and MMS permissions for a user.
type SMSPolicy struct {
	*PolicyState
	// InternationalSms enables sending international SMS and MMS messages.
	InternationalSms bool `json:"international_sms"`
	// InternationalSmsCountries lists the countries to which international SMS is
	// allowed.
	InternationalSmsCountries []string `json:"international_sms_countries"`
	// AllowCopy allows copying SMS and MMS message content.
	AllowCopy bool `json:"allow_copy"`
	// AllowPaste allows pasting content into SMS and MMS messages.
	AllowPaste bool `json:"allow_paste"`
}

// VoicemailPolicy configures voicemail features for a user.
type VoicemailPolicy struct {
	*PolicyState
	// AllowTranscription enables voicemail-to-text transcription.
	AllowTranscription bool `json:"allow_transcription"`
	// AllowVideomail enables video voicemail.
	AllowVideomail bool `json:"allow_videomail"`
}

// ZoomPhoneOnMobilePolicy configures which mobile client applications are
// permitted to use Zoom Phone calling and SMS features.
type ZoomPhoneOnMobilePolicy struct {
	*PolicyState
	// AllowCallingClients lists the mobile client identifiers permitted for voice
	// calls.
	AllowCallingClients []string `json:"allow_calling_clients"`
	// AllowSmsMmsClients lists the mobile client identifiers permitted for SMS and
	// MMS.
	AllowSmsMmsClients []string `json:"allow_sms_mms_clients"`
}

// PersonalAudioLibraryPolicy configures what a user can customise in their
// personal audio library.
type PersonalAudioLibraryPolicy struct {
	*PolicyState
	// AllowMusicOnHoldCustomization allows the user to set a custom music-on-hold
	// audio file.
	AllowMusicOnHoldCustomization bool `json:"allow_music_on_hold_customization"`
	// AllowVoicemailAndMessageGreetingCustomization allows the user to set custom
	// voicemail and message greeting audio files.
	AllowVoicemailAndMessageGreetingCustomization bool `json:"allow_voicemail_and_message_greeting_customization"`
}

// VoicemailNotificationByEmailPolicy configures email notifications sent when a
// new voicemail is received.
type VoicemailNotificationByEmailPolicy struct {
	*PolicyState
	// IncludeVoicemailFile attaches the voicemail audio file to the email
	// notification.
	IncludeVoicemailFile bool `json:"include_voicemail_file"`
	// IncludeVoicemailTranscription attaches the voicemail transcription to the
	// email notification.
	IncludeVoicemailTranscription bool `json:"include_voicemail_transcription"`
}

// ZoomPhoneOnDesktopPolicy configures which desktop client applications are
// permitted to use Zoom Phone calling and SMS features.
type ZoomPhoneOnDesktopPolicy struct {
	*PolicyState
	// AllowCallingClients lists the desktop client identifiers permitted for voice
	// calls.
	AllowCallingClients []string `json:"allow_calling_clients"`
	// AllowSmsMmsClients lists the desktop client identifiers permitted for SMS
	// and MMS.
	AllowSmsMmsClients []string `json:"allow_sms_mms_clients"`
}
