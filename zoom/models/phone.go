package models

import (
	"time"

	"github.com/TheSlowpes/go-zoom/zoom/enums"
)

type CallHistory struct {
	ID                string              `json:"id,omitempty"`
	CallHistoryUUID   string              `json:"call_history_uuid,omitempty"`
	CallID            string              `json:"call_id,omitempty"`
	ConnectType       enums.ConnectType   `json:"connect_type,omitempty"`
	CallType          string              `json:"call_type,omitempty"`
	Direction         enums.Direction     `json:"direction,omitempty"`
	International     bool                `json:"international"`
	HideCallerID      bool                `json:"hide_caller_id"`
	EndToEnd          bool                `json:"end_to_end"`
	CallerExtID       string              `json:"caller_ext_id,omitempty"`
	CallerName        string              `json:"caller_name,omitempty"`
	CallerDidNumber   string              `json:"caller_did_number,omitempty"`
	CallerExtNumber   string              `json:"caller_ext_number,omitempty"`
	CallerEmail       string              `json:"caller_email,omitempty"`
	CallerExtType     enums.ExtensionType `json:"caller_ext_type,omitempty"`
	CalleeExtID       string              `json:"callee_ext_id,omitempty"`
	CalleeName        string              `json:"callee_name,omitempty"`
	CalleeEmail       string              `json:"callee_email,omitempty"`
	CalleeDidNumber   string              `json:"callee_did_number,omitempty"`
	CalleeExtNumber   string              `json:"callee_ext_number,omitempty"`
	CalleeExtType     enums.ExtensionType `json:"callee_ext_type,omitempty"`
	Department        string              `json:"department,omitempty"`
	CostCenter        string              `json:"cost_center,omitempty"`
	SiteID            string              `json:"site_id,omitempty"`
	GroupID           string              `json:"group_id,omitempty"`
	SiteName          string              `json:"site_name,omitempty"`
	StartTime         time.Time           `json:"start_time"`
	AnswerTime        time.Time           `json:"answer_time"`
	EndTime           time.Time           `json:"end_time"`
	CallPath          []CallPath          `json:"call_path"`
	CallerAccountCode string              `json:"caller_account_code,omitempty"`
	CalleeAccountCode string              `json:"callee_account_code,omitempty"`
	CallElements      []CallElement       `json:"call_elements"`
}

type CallPath struct {
	ID                   string              `json:"id,omitempty"`
	CallID               string              `json:"call_id,omitempty"`
	ConnectType          enums.ConnectType   `json:"connect_type,omitempty"`
	CallType             enums.CallType      `json:"call_type,omitempty"`
	Direction            enums.Direction     `json:"direction,omitempty"`
	HideCallerID         bool                `json:"hide_caller_id"`
	EndToEnd             bool                `json:"end_to_end"`
	CallerExtID          string              `json:"caller_ext_id,omitempty"`
	CallerName           string              `json:"caller_name,omitempty"`
	CallerEmail          string              `json:"caller_email,omitempty"`
	CallerDidNumber      string              `json:"caller_did_number,omitempty"`
	CallerExtNumber      string              `json:"caller_ext_number,omitempty"`
	CallerExtType        enums.ExtensionType `json:"caller_ext_type,omitempty"`
	CallerNumberType     enums.NumberType    `json:"caller_number_type,omitempty"`
	CallerDeviceType     string              `json:"caller_device_type,omitempty"`
	CallerCountryIsoCode string              `json:"caller_country_iso_code,omitempty"`
	CallerCountryCode    string              `json:"caller_country_code,omitempty"`
	CalleeExtID          string              `json:"callee_ext_id,omitempty"`
	CalleeName           string              `json:"callee_name,omitempty"`
	CalleeDidNumber      string              `json:"callee_did_number,omitempty"`
	CalleeExtNumber      string              `json:"callee_ext_number,omitempty"`
	CalleeEmail          string              `json:"callee_email,omitempty"`
	CalleeExtType        enums.ExtensionType `json:"callee_ext_type,omitempty"`
	CalleeNumberType     enums.NumberType    `json:"callee_number_type,omitempty"`
	CalleeDeviceType     string              `json:"callee_device_type,omitempty"`
	CalleeCountryIsoCode string              `json:"callee_country_iso_code,omitempty"`
	CalleeCountryCode    string              `json:"callee_country_code,omitempty"`
	ClientCode           string              `json:"client_code,omitempty"`
	Department           string              `json:"department,omitempty"`
	CostCenter           string              `json:"cost_center,omitempty"`
	SiteID               string              `json:"site_id,omitempty"`
	GroupID              string              `json:"group_id,omitempty"`
	SiteName             string              `json:"site_name,omitempty"`
	StartTime            time.Time           `json:"start_time"`
	AnswerTime           time.Time           `json:"answer_time"`
	EndTime              time.Time           `json:"end_time"`
	Event                string              `json:"event,omitempty"`
	Result               enums.CallResult    `json:"result,omitempty"`
	ResultReason         string              `json:"result_reason,omitempty"`
	DevicePrivateIP      string              `json:"device_private_ip,omitempty"`
	DevicePublicIP       string              `json:"device_public_ip,omitempty"`
	OperatorExtNumber    string              `json:"operator_ext_number,omitempty"`
	OperatorExtID        string              `json:"operator_ext_id,omitempty"`
	OperatorExtType      string              `json:"operator_ext_type,omitempty"`
	OperatorName         string              `json:"operator_name,omitempty"`
	PressKey             string              `json:"press_key,omitempty"`
	Segment              int                 `json:"segment"`
	Node                 int                 `json:"node"`
	IsNode               int                 `json:"is_node"`
	RecordingID          string              `json:"recording_id,omitempty"`
	RecordingType        string              `json:"recording_type,omitempty"`
	HoldTime             int                 `json:"hold_time"`
	WaitingTime          int                 `json:"waiting_time"`
	VoicemailID          string              `json:"voicemail_id,omitempty"`
}

type CallElement struct {
	CallElementID        string    `json:"call_element_id,omitempty"`
	CallID               string    `json:"call_id,omitempty"`
	ConnectType          string    `json:"connect_type,omitempty"`
	CallType             string    `json:"call_type,omitempty"`
	Direction            string    `json:"direction,omitempty"`
	HideCallerID         bool      `json:"hide_caller_id"`
	EndToEnd             bool      `json:"end_to_end"`
	CallerExtID          string    `json:"caller_ext_id,omitempty"`
	CallerName           string    `json:"caller_name,omitempty"`
	CallerEmail          string    `json:"caller_email,omitempty"`
	CallerDidNumber      string    `json:"caller_did_number,omitempty"`
	CallerExtNumber      string    `json:"caller_ext_number,omitempty"`
	CallerExtType        string    `json:"caller_ext_type,omitempty"`
	CallerNumberType     string    `json:"caller_number_type,omitempty"`
	CallerDeviceType     string    `json:"caller_device_type,omitempty"`
	CallerCountryIsoCode string    `json:"caller_country_iso_code,omitempty"`
	CallerCountryCode    string    `json:"caller_country_code,omitempty"`
	CalleeExtID          string    `json:"callee_ext_id,omitempty"`
	CalleeName           string    `json:"callee_name,omitempty"`
	CalleeDidNumber      string    `json:"callee_did_number,omitempty"`
	CalleeExtNumber      string    `json:"callee_ext_number,omitempty"`
	CalleeEmail          string    `json:"callee_email,omitempty"`
	CalleeExtType        string    `json:"callee_ext_type,omitempty"`
	CalleeNumberType     string    `json:"callee_number_type,omitempty"`
	CalleeDeviceType     string    `json:"callee_device_type,omitempty"`
	CalleeCountryIsoCode string    `json:"callee_country_iso_code,omitempty"`
	CalleeCountryCode    string    `json:"callee_country_code,omitempty"`
	ClientCode           string    `json:"client_code,omitempty"`
	Department           string    `json:"department,omitempty"`
	CostCenter           string    `json:"cost_center,omitempty"`
	SiteID               string    `json:"site_id,omitempty"`
	GroupID              string    `json:"group_id,omitempty"`
	SiteName             string    `json:"site_name,omitempty"`
	StartTime            time.Time `json:"start_time"`
	AnswerTime           time.Time `json:"answer_time"`
	EndTime              time.Time `json:"end_time"`
	Event                string    `json:"event,omitempty"`
	Result               string    `json:"result,omitempty"`
	ResultReason         string    `json:"result_reason,omitempty"`
	DevicePrivateIP      string    `json:"device_private_ip,omitempty"`
	DevicePublicIP       string    `json:"device_public_ip,omitempty"`
	OperatorExtNumber    string    `json:"operator_ext_number,omitempty"`
	OperatorExtID        string    `json:"operator_ext_id,omitempty"`
	OperatorExtType      string    `json:"operator_ext_type,omitempty"`
	OperatorName         string    `json:"operator_name,omitempty"`
	PressKey             string    `json:"press_key,omitempty"`
	Segment              int       `json:"segment"`
	Node                 int       `json:"node"`
	IsNode               int       `json:"is_node"`
	RecordingID          string    `json:"recording_id,omitempty"`
	RecordingType        string    `json:"recording_type,omitempty"`
	HoldTime             int       `json:"hold_time"`
	WaitingTime          int       `json:"waiting_time"`
	VoicemailID          string    `json:"voicemail_id,omitempty"`
}

type AICallSummary struct {
	AICallSummaryID    string    `json:"ai_call_summary_id,omitempty"`
	AccountID          string    `json:"account_id,omitempty"`
	CallID             string    `json:"call_id,omitempty"`
	CallLogIDs         []string  `json:"call_log_ids,omitempty"`
	UserID             string    `json:"user_id,omitempty"`
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
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type BYOC struct {
	Enable bool `json:"enable"`
}

type Country struct {
	Code        string `json:"code,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Name        string `json:"name,omitempty"`
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
	AutoDeletePolicy      string         `json:"auto_delete_policy,omitempty"`
	CallID                string         `json:"call_id,omitempty"`
	CallLogID             string         `json:"call_log_id,omitempty"`
	CallElementID         string         `json:"call_element_id,omitempty"`
	CalleeName            string         `json:"callee_name,omitempty"`
	CalleeNumber          string         `json:"callee_number,omitempty"`
	CalleeNumberType      int            `json:"callee_number_type"`
	CallerName            string         `json:"caller_name,omitempty"`
	CallerNumber          string         `json:"caller_number,omitempty"`
	CallerNumberType      int            `json:"caller_number_type"`
	OutgoingBy            RecordingUser  `json:"outgoing_by"`
	AcceptedBy            RecordingUser  `json:"accepted_by"`
	DateTime              time.Time      `json:"date_time"`
	DisclaimerStatus      int            `json:"disclaimer_status"`
	Direction             string         `json:"direction,omitempty"`
	DownloadURL           string         `json:"download_url,omitempty"`
	Duration              int            `json:"duration"`
	EndTime               time.Time      `json:"end_time"`
	ID                    string         `json:"id,omitempty"`
	MeetingUUID           string         `json:"meeting_uuid,omitempty"`
	Owner                 RecordingOwner `json:"owner"`
	RecordingType         string         `json:"recording_type,omitempty"`
	Site                  RecordingSite  `json:"site"`
	TranscriptDownloadURL string         `json:"transcript_download_url,omitempty"`
	AutoDeleteEnable      bool           `json:"auto_delete_enable"`
	CallerAccountCode     string         `json:"caller_account_code,omitempty"`
	CalleeAccountCode     string         `json:"callee_account_code,omitempty"`
}

type RecordingUser struct {
	Name            string `json:"name,omitempty"`
	ExtensionNumber string `json:"extension_number,omitempty"`
}

type RecordingOwner struct {
	ExtensionNumber      int       `json:"extension_number"`
	ID                   string    `json:"id,omitempty"`
	Name                 string    `json:"name,omitempty"`
	Type                 string    `json:"type,omitempty"`
	ExtensionStatus      string    `json:"extension_status,omitempty"`
	ExtensionDeletedTime time.Time `json:"extension_deleted_time"`
}

type RecordingSite struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type RecordingTranscript struct {
	Type           string               `json:"type,omitempty"`
	Ver            string               `json:"ver,omitempty"`
	RecordingID    string               `json:"recording_id,omitempty"`
	MeetingID      string               `json:"meeting_id,omitempty"`
	AccountID      string               `json:"account_id,omitempty"`
	HostID         string               `json:"host_id,omitempty"`
	RecordingStart time.Time            `json:"recording_start"`
	RecordingEnd   time.Time            `json:"recording_end"`
	Timeline       []TranscriptTimeline `json:"timeline"`
}

type TranscriptTimeline struct {
	Text        string           `json:"text,omitempty"`
	RawText     string           `json:"raw_text,omitempty"`
	TS          time.Time        `json:"ts"`
	EndTS       time.Time        `json:"end_ts"`
	Users       []TranscriptUser `json:"users"`
	UserID      string           `json:"userId,omitempty"`
	UserIDs     []string         `json:"userIds"`
	ChannelMark int              `json:"channelMark"`
}

type TranscriptUser struct {
	Username                 string `json:"username,omitempty"`
	MultiplePeople           bool   `json:"multiple_people"`
	UserID                   string `json:"user_id,omitempty"`
	ZoomUserID               string `json:"zoom_userid,omitempty"`
	AvatarURL                string `json:"avatar_url,omitempty"`
	ClientType               int    `json:"client_type"`
	EmailAddress             string `json:"email_address,omitempty"`
	ChannelMark              int    `json:"channel_mark"`
	Pronoun                  string `json:"pronoun,omitempty"`
	EnableSpeakerDiarization string `json:"enable_speaker_diarization,omitempty"`
	TS                       int    `json:"ts"`
	IsAutoGenerated          bool   `json:"is_auto_generated"`
}

type PhoneUser struct {
	ActivationStatus string           `json:"activation_status,omitempty"`
	CallingPlans     []CallingPlan    `json:"calling_plans"`
	CostCenter       string           `json:"cost_center,omitempty"`
	Department       string           `json:"department,omitempty"`
	Email            string           `json:"email,omitempty"`
	EmergencyAddress EmergencyAddress `json:"emergency_address"`
	ExtensionID      string           `json:"extension_id,omitempty"`
	ExtensionNumber  int              `json:"extension_number"`
	ID               string           `json:"id,omitempty"`
	PhoneNumbers     []PhoneNumber    `json:"phone_numbers"`
	PhoneUserID      string           `json:"phone_user_id,omitempty"`
	Policy           PhoneUserPolicy  `json:"policy"`
	SiteAdmin        bool             `json:"site_admin"`
	SiteID           string           `json:"site_id,omitempty"`
	Site             PhoneUserSite    `json:"site"`
	Status           string           `json:"status,omitempty"`
}

type PhoneUserSite struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type EmergencyAddress struct {
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city,omitempty"`
	Country      string `json:"country,omitempty"`
	ID           string `json:"id,omitempty"`
	StateCode    string `json:"state_code,omitempty"`
	Zip          string `json:"zip,omitempty"`
}

type CallingPlan struct {
	Type                    int    `json:"type"`
	BillingAccountID        string `json:"billing_account_id,omitempty"`
	BillingAccountName      string `json:"billing_account_name,omitempty"`
	BillingSubscriptionID   string `json:"billing_subscription_id,omitempty"`
	BillingSubscriptionName string `json:"billing_subscription_name,omitempty"`
}

type PhoneUserPolicy struct {
	AdHocCallRecording                 AdHocCallRecording                 `json:"ad_hoc_call_recording"`
	AdHocCallRecordingAccessMembers    []PolicyAccessMember               `json:"ad_hoc_call_recording_access_members"`
	AutoCallRecording                  AutoCallRecordingPolicy            `json:"auto_call_recording"`
	AutoCallRecordingAccessMembers     []*PolicyAccessMember              `json:"auto_call_recording_access_members"`
	CallOverflow                       CallOverflowPolicy                 `json:"call_overflow"`
	CallPark                           CallParkPolicy                     `json:"call_park"`
	CallTransferring                   CallTransferPolicy                 `json:"call_transferring"`
	Delegation                         bool                               `json:"delegation"`
	ElevateToMeeting                   bool                               `json:"elevate_to_meeting"`
	EmergencyAddressManagement         AddressManagementPolicy            `json:"emergency_address_management"`
	EmergencyCallsToPsap               bool                               `json:"emergency_calls_to_psap"`
	CallHandlingForwardingToOtherUsers CallForwardPolicy                  `json:"call_handling_forwarding_to_other_users"`
	HandOffToRoom                      *PolicyState                       `json:"hand_off_to_room,omitempty"`
	InternationalCalling               bool                               `json:"international_calling"`
	MobileSwitchToCarrier              *PolicyState                       `json:"mobile_switch_to_carrier,omitempty"`
	SelectOutboundCallerID             OutboundCallerIDPolicy             `json:"select_outbound_caller_id"`
	SMS                                SMSPolicy                          `json:"sms"`
	Voicemail                          VoicemailPolicy                    `json:"voicemail"`
	VoicemailAccessMembers             []*PolicyAccessMember              `json:"voicemail_access_members"`
	ZoomPhoneOnMobile                  ZoomPhoneOnMobilePolicy            `json:"zoom_phone_on_mobile"`
	PersonalAudioLibrary               PersonalAudioLibraryPolicy         `json:"personal_audio_library"`
	VoicemailTranscription             *PolicyState                       `json:"voicemail_transcription,omitempty"`
	VoicemailNotificationByEmail       VoicemailNotificationByEmailPolicy `json:"voicemail_notification_by_email"`
	SharedVoicemailNotificationByEmail *PolicyState                       `json:"shared_voicemail_notification_by_email,omitempty"`
	CheckVoicemailsOverPhone           *PolicyState                       `json:"check_voicemails_over_phone,omitempty"`
	AudioIntercom                      *PolicyState                       `json:"audio_intercom,omitempty"`
	PeerToPeerMedia                    *PolicyState                       `json:"peer_to_peer_media,omitempty"`
	E2EEncryption                      *PolicyState                       `json:"e2e_encryption,omitempty"`
	OutboundCalling                    *PolicyState                       `json:"outbound_calling,omitempty"`
	OutboundSms                        *PolicyState                       `json:"outbound_sms,omitempty"`
	AllowEndUserEditCallHandling       *PolicyState                       `json:"allow_end_user_edit_call_handling,omitempty"`
	VoicemailIntentBasedPrioritization *PolicyState                       `json:"voicemail_intent_based_prioritization,omitempty"`
	VoicemailTasks                     *PolicyState                       `json:"voicemail_tasks,omitempty"`
	ZoomPhoneOnDesktop                 ZoomPhoneOnDesktopPolicy           `json:"zoom_phone_on_desktop"`
}

type PhoneUserSettings struct {
	AreaCode                        string                    `json:"area_code,omitempty"`
	AudioPromptLanguage             string                    `json:"audio_prompt_language,omitempty"`
	CompanyNumber                   string                    `json:"company_number,omitempty"`
	Country                         Country                   `json:"country"`
	Delegation                      DelegationSettings        `json:"delegation"`
	DeskPhone                       DeskPhoneSettings         `json:"desk_phone"`
	ExtensionNumber                 int                       `json:"extension_number"`
	MusicOnHoldID                   string                    `json:"music_on_hold_id,omitempty"`
	OutboundCaller                  OutboundCallerIDSetting   `json:"outbound_caller"`
	OutboundCallerIds               []OutboundCallerIDSetting `json:"outbound_caller_ids"`
	Status                          string                    `json:"status,omitempty"`
	VoiceMail                       []PolicyAccessMember      `json:"voice_mail"`
	Intercom                        IntercomSettings          `json:"intercom"`
	AutoCallRecordingAccessMembers  []PolicyAccessMember      `json:"auto_call_recording_access_members"`
	AdHocCallRecordingAccessMembers []PolicyAccessMember      `json:"ad_hoc_call_recording_access_members"`
	SharedLinesCallSetting          SharedLineCallSetting     `json:"shared_lines_call_setting"`
}

type DeskPhoneSettings struct {
	Devices         []Device      `json:"devices"`
	KeysPositions   KeysPositions `json:"keys_positions"`
	PhoneScreenLock bool          `json:"phone_screen_lock"`
	PinCode         string        `json:"pin_code,omitempty"`
}

type KeysPositions struct {
	PrimaryNumber string `json:"primary_number,omitempty"`
}

type DelegationSettings struct {
	Assistants []PhoneUserExtension `json:"assistants"`
	Privacy    bool                 `json:"privacy"`
	Privileges int                  `json:"privileges"`
	Locked     bool                 `json:"locked"`
}

type OutboundCallerIDSetting struct {
	IsDefault bool   `json:"is_default"`
	Name      string `json:"name,omitempty"`
	Number    string `json:"number,omitempty"`
}

type DevicePolicy struct {
	CallControl PolicyStatus `json:"call_control"`
	HotDesking  PolicyStatus `json:"hot_desking"`
}

type PolicyStatus struct {
	Status string `json:"status,omitempty"`
}

type AudioIntercomSettings struct {
	*PhoneUserExtension
	Status       string `json:"status,omitempty"`
	DeviceID     string `json:"device_id,omitempty"`
	DeviceStatus string `json:"device_status,omitempty"`
}

type IntercomSettings struct {
	AudioIntercoms []AudioIntercomSettings `json:"audio_intercoms"`
	Device         Device                  `json:"device"`
}

type Device struct {
	DeviceType  string       `json:"device_type,omitempty"`
	DisplayName string       `json:"display_name,omitempty"`
	ID          string       `json:"id,omitempty"`
	Policy      DevicePolicy `json:"policy"`
	Name        string       `json:"name,omitempty"`
	Status      string       `json:"status,omitempty"`
	MacAddress  string       `json:"mac_address,omitempty"`
	PrivateIP   string       `json:"private_ip,omitempty"`
	PublicIP    string       `json:"public_ip,omitempty"`
}

type SharedLineCallSetting struct {
	SharedLineAppearances SharedLineAppearances `json:"shared_line_appearances"`
	SharedLineGroups      SharedLineGroups      `json:"shared_line_groups"`
}

type SharedLineAppearances struct {
	Executives []Executive `json:"executives"`
}

type Executive struct {
	UserID       string `json:"user_id,omitempty"`
	DisplayName  string `json:"display_name,omitempty"`
	ReceiveCalls bool   `json:"receive_calls"`
	AllowOptOut  bool   `json:"allow_opt_out"`
}

type SharedLineGroups struct {
	ReceiveCalls    bool                     `json:"receive_calls"`
	SharedLineGroup []SharedLineGroupSetting `json:"shared_line_group"`
}

type SharedLineGroupSetting struct {
	SlgID        string `json:"slg_id,omitempty"`
	DisplayName  string `json:"display_name,omitempty"`
	ReceiveCalls bool   `json:"receive_calls"`
}

type AdHocCallRecording struct {
	*PolicyState
	RecordingStartPrompt   bool                  `json:"recording_start_prompt"`
	RecordingTranscription bool                  `json:"recording_transcription"`
	PlayRecordingBeepTone  PlayRecordingBeepTone `json:"play_recording_beep_tone"`
}

type PlayRecordingBeepTone struct {
	*PolicyState
	PlayBeepVolume       int    `json:"play_beep_volume"`
	PlayBeepTimeInterval int    `json:"play_beep_time_interval"`
	PlayBeepMember       string `json:"play_beep_member,omitempty"`
}

type PolicyState struct {
	Enable   bool   `json:"enable"`
	Locked   bool   `json:"locked"`
	LockedBy string `json:"locked_by,omitempty"`
	Modified bool   `json:"modified"`
}

type AutoCallRecordingPolicy struct {
	*PolicyState
	AllowStopResumeRecording     bool                         `json:"allow_stop_resume_recording"`
	DisconnectOnRecordingFailure bool                         `json:"disconnect_on_recording_failure"`
	RecordingCalls               string                       `json:"recording_calls,omitempty"`
	RecordingTranscription       bool                         `json:"recording_transcription"`
	PlayRecordingBeepTone        PlayRecordingBeepTone        `json:"play_recording_beep_tone"`
	InboundAudioNotification     *RecordingNotificationPolicy `json:"inbound_audio_notification,omitempty"`
	OutboundAudioNotification    *RecordingNotificationPolicy `json:"outbound_audio_notification,omitempty"`
}

type PolicyAccessMember struct {
	AccessUserID   string `json:"access_user_id,omitempty"`
	AccessUserType string `json:"access_user_type,omitempty"`
	AllowDelete    bool   `json:"allow_delete"`
	Delete         bool   `json:"delete"`
	AllowDownload  bool   `json:"allow_download"`
	Download       bool   `json:"download"`
	AllowSharing   bool   `json:"allow_sharing"`
	SharedID       string `json:"shared_id,omitempty"`
}

type RecordingNotificationPolicy struct {
	RecordingStartPrompt     bool `json:"recording_start_prompt"`
	RecordingExplicitConsent bool `json:"recording_explicit_consent"`
}

type CallOverflowPolicy struct {
	*PolicyState
	CallOverflowType enums.CallForwardingType `json:"call_overflow_type,omitempty"`
}

type CallParkPolicy struct {
	*PolicyState
	CallNotPickedUpAction enums.CallNotPickedUpAction `json:"call_not_picked_up_action,omitempty"`
	ExpirationPeriod      int                         `json:"expiration_period"`
	ForwardTo             PhoneUserExtension          `json:"forward_to"`
}

type PhoneUserExtension struct {
	DisplayName     string `json:"display_name,omitempty"`
	ExtensionID     string `json:"extension_id,omitempty"`
	ExtensionNumber any    `json:"extension_number,omitempty"`
	ExtensionType   string `json:"extension_type,omitempty"`
	ID              string `json:"id,omitempty"`
}

type CallTransferPolicy struct {
	*PolicyState
	CallTransferringType enums.CallTransferType `json:"call_transferring_type,omitempty"`
}

type AddressManagementPolicy struct {
	*PolicyState
	PromptDefaultAddress bool `json:"prompt_default_address"`
}

type CallForwardPolicy struct {
	*PolicyState
	CallForwardingType enums.CallForwardingType `json:"call_forwarding_type,omitempty"`
}

type OutboundCallerIDPolicy struct {
	*PolicyState
	AllowHideOutboundCallerID bool `json:"allow_hide_outbound_caller_id"`
}

type SMSPolicy struct {
	*PolicyState
	InternationalSms          bool     `json:"international_sms"`
	InternationalSmsCountries []string `json:"international_sms_countries"`
	AllowCopy                 bool     `json:"allow_copy"`
	AllowPaste                bool     `json:"allow_paste"`
}

type VoicemailPolicy struct {
	*PolicyState
	AllowTranscription bool `json:"allow_transcription"`
	AllowVideomail     bool `json:"allow_videomail"`
}

type ZoomPhoneOnMobilePolicy struct {
	*PolicyState
	AllowCallingClients []string `json:"allow_calling_clients"`
	AllowSmsMmsClients  []string `json:"allow_sms_mms_clients"`
}

type PersonalAudioLibraryPolicy struct {
	*PolicyState
	AllowMusicOnHoldCustomization                 bool `json:"allow_music_on_hold_customization"`
	AllowVoicemailAndMessageGreetingCustomization bool `json:"allow_voicemail_and_message_greeting_customization"`
}

type VoicemailNotificationByEmailPolicy struct {
	*PolicyState
	IncludeVoicemailFile          bool `json:"include_voicemail_file"`
	IncludeVoicemailTranscription bool `json:"include_voicemail_transcription"`
}

type ZoomPhoneOnDesktopPolicy struct {
	*PolicyState
	AllowCallingClients []string `json:"allow_calling_clients"`
	AllowSmsMmsClients  []string `json:"allow_sms_mms_clients"`
}
