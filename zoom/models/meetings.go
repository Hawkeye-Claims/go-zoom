package models

import (
	"time"

	"github.com/Hawkeye-Claims/go-zoom/zoom/enums"
)

// Meeting represents a Zoom meeting as returned by the Meetings API.
type Meeting struct {
	// AssistantID is the ID of a scheduling assistant, if one was used.
	AssistantID string `json:"assistant_id"`
	// HostEmail is the email address of the meeting host.
	HostEmail string `json:"host_email"`
	// HostID is the user ID of the meeting host.
	HostID string `json:"host_id"`
	// ID is the meeting ID (numeric).
	ID int64 `json:"id"`
	// UUID is the unique meeting UUID for this particular occurrence.
	UUID string `json:"uuid"`
	// RegistrationURL is the URL attendees use to register for the meeting.
	RegistrationURL string `json:"registration_url"`
	// Agenda is the meeting's description or agenda text.
	Agenda string `json:"agenda"`
	// CreatedAt is the time the meeting was created.
	CreatedAt time.Time `json:"created_at"`
	// Duration is the meeting duration in minutes.
	Duration int `json:"duration"`
	// EncryptedPassword is the encrypted meeting password.
	EncryptedPassword string `json:"encrypted_password"`
	// PstnPassword is the numeric password for PSTN dial-in participants.
	PstnPassword string `json:"pstn_password"`
	// H323Password is the H.323/SIP room system password.
	H323Password string `json:"h323_password"`
	// JoinURL is the URL participants use to join the meeting.
	JoinURL string `json:"join_url"`
	// ChatJoinURL is the URL used to join the meeting's persistent chat channel.
	ChatJoinURL string `json:"chat_join_url"`
	// Occurrences lists individual occurrences for a recurring meeting.
	Occurrences []MeetingOccurrence `json:"occurrences"`
	// Password is the plain-text meeting password.
	Password string `json:"password"`
	// Pmi is the Personal Meeting ID used for this meeting, if any.
	Pmi string `json:"pmi"`
	// PreSchedule indicates whether the meeting is pre-scheduled.
	PreSchedule bool `json:"pre_schedule"`
	// Recurrence describes the recurrence rules for a recurring meeting.
	Recurrence MeetingRecurrence `json:"recurrence"`
	// Settings contains the full set of meeting configuration options.
	Settings MeetingSettings `json:"settings"`
	// StartTime is the scheduled start time of the meeting.
	StartTime time.Time `json:"start_time"`
	// StartURL is the URL the host uses to start the meeting.
	StartURL string `json:"start_url"`
	// Status is the current status of the meeting (e.g. "waiting", "started").
	Status string `json:"status"`
	// TemplateID is the ID of the meeting template used to create this meeting.
	TemplateID string `json:"template_id"`
	// Timezone is the IANA timezone identifier for the meeting's start time.
	Timezone string `json:"timezone"`
	// Topic is the meeting title.
	Topic string `json:"topic"`
	// TrackingFields contains custom tracking key-value pairs attached to the meeting.
	TrackingFields []MeetingTrackingField `json:"tracking_fields"`
	// Type is the meeting type (instant, scheduled, recurring, etc.).
	Type enums.MeetingType `json:"type"`
	// DynamicHostKey is the host key generated dynamically for the meeting.
	DynamicHostKey string `json:"dynamic_host_key"`
	// CreationSource describes how the meeting was created.
	CreationSource string `json:"creation_source"`
}

// MeetingOccurrence represents a single scheduled occurrence of a recurring meeting.
type MeetingOccurrence struct {
	// Duration is the occurrence duration in minutes.
	Duration int `json:"duration"`
	// OccurrenceID is the unique identifier for this occurrence.
	OccurrenceID string `json:"occurrence_id"`
	// StartTime is the scheduled start time for this occurrence.
	StartTime time.Time `json:"start_time"`
	// Status is the occurrence status (e.g. "available", "deleted").
	Status string `json:"status"`
}

// MeetingRecurrence describes the repeat schedule for a recurring meeting.
type MeetingRecurrence struct {
	// EndDateTime is the date and time at which the recurrence ends.
	EndDateTime time.Time `json:"end_date_time"`
	// EndTimes is the number of occurrences after which the recurrence ends.
	// A value of 0 means it ends by EndDateTime instead.
	EndTimes int `json:"end_times"`
	// MonthlyDay is the day of the month for a monthly recurrence (1–31).
	MonthlyDay int `json:"monthly_day"`
	// MonthlyWeek is the week of the month for a monthly-by-day recurrence
	// (-1 = last week, 1–4 = first through fourth week).
	MonthlyWeek int `json:"monthly_week"`
	// MonthlyWeekDay is the day of the week for a monthly-by-day recurrence
	// (1 = Sunday, 7 = Saturday).
	MonthlyWeekDay int `json:"monthly_week_day"`
	// RepeatInterval is how often the recurrence repeats (e.g. every 2 weeks).
	RepeatInterval int `json:"repeat_interval"`
	// Type is the recurrence cadence (daily, weekly, or monthly).
	Type enums.MeetingRecurrenceType `json:"type"`
	// WeeklyDays is a comma-separated list of day numbers for a weekly recurrence
	// (1 = Sunday, 7 = Saturday).
	WeeklyDays string `json:"weekly_days"`
}

// MeetingSettings holds the full set of configuration options for a meeting.
type MeetingSettings struct {
	// AdditionalDataCenterRegions lists additional data centre regions to use for
	// the meeting.
	AdditionalDataCenterRegions []string `json:"additional_data_center_regions"`
	// AllowMultipleDevices permits participants to join from multiple devices
	// simultaneously.
	AllowMultipleDevices bool `json:"allow_multiple_devices"`
	// AlternativeHosts is a comma-separated list of alternative host email
	// addresses or user IDs.
	AlternativeHosts string `json:"alternative_hosts"`
	// AlternativeHostsEmailNotification controls whether alternative hosts receive
	// an email notification when they are assigned.
	AlternativeHostsEmailNotification bool `json:"alternative_hosts_email_notification"`
	// AlternativeHostUpdatePolls allows alternative hosts to add or edit polls.
	AlternativeHostUpdatePolls bool `json:"alternative_host_update_polls"`
	// AlternativeHostManageMeetingSummary allows alternative hosts to manage AI
	// meeting summaries.
	AlternativeHostManageMeetingSummary bool `json:"alternative_host_manage_meeting_summary"`
	// AlternativeHostManageCloudRecording allows alternative hosts to manage cloud
	// recordings.
	AlternativeHostManageCloudRecording bool `json:"alternative_host_manage_cloud_recording"`
	// ApprovalType controls registration approval (0 = automatic, 1 = manual,
	// 2 = no registration required).
	ApprovalType int `json:"approval_type"`
	// ApprovedOrDeniedCountriesOrRegions restricts meeting access by country or
	// region.
	ApprovedOrDeniedCountriesOrRegions MeetingApprovedOrDeniedCountriesOrRegions `json:"approved_or_denied_countries_or_regions"`
	// Audio specifies the audio connection options ("both", "telephony", "voip",
	// "thirdParty").
	Audio string `json:"audio"`
	// AudioConferenceInfo contains third-party audio conference details.
	AudioConferenceInfo string `json:"audio_conference_info"`
	// AuthenticationDomains is a comma-separated list of email domains allowed to
	// join when authentication is required.
	AuthenticationDomains string `json:"authentication_domains"`
	// AuthenticationException lists participants who are exempt from the meeting's
	// authentication requirement.
	AuthenticationException []MeetingAuthenticationException `json:"authentication_exception"`
	// AuthenticationName is the display name of the authentication profile in use.
	AuthenticationName string `json:"authentication_name"`
	// AuthenticationOption is the ID of the authentication profile to apply.
	AuthenticationOption string `json:"authentication_option"`
	// AutoRecording controls automatic cloud recording ("local", "cloud", "none").
	AutoRecording string `json:"auto_recording"`
	// AutoAddRecordingToVideoManagement controls automatic addition of recordings
	// to Zoom Video Management channels.
	AutoAddRecordingToVideoManagement AutoAddRecordingToVideoManagement `json:"auto_add_recording_to_video_management"`
	// BreakoutRoom configures breakout rooms for the meeting.
	BreakoutRoom BreakoutRoom `json:"breakout_room"`
	// CalendarType specifies the calendar integration type (1 = Zoom, 2 = Google,
	// 3 = Exchange).
	CalendarType int `json:"calendar_type"`
	// CloseRegistration closes registration after the meeting start time.
	CloseRegistration bool `json:"close_registration"`
	// ContactEmail is the host's contact email shown on the registration page.
	ContactEmail string `json:"contact_email"`
	// ContactName is the host's contact name shown on the registration page.
	ContactName string `json:"contact_name"`
	// CustomKeys contains up to ten custom key-value pairs for the meeting.
	CustomKeys []CustomKey `json:"custom_keys"`
	// EmailNotification controls whether email notifications are sent to
	// registrants and attendees.
	EmailNotification bool `json:"email_notification"`
	// EncryptionType specifies the end-to-end encryption type.
	EncryptionType string `json:"encryption_type"`
	// FocusMode enables focus mode at the start of the meeting.
	FocusMode bool `json:"focus_mode"`
	// GlobalDialInCountries lists country codes for which global dial-in numbers
	// are provided.
	GlobalDialInCountries []string `json:"global_dial_in_countries"`
	// GlobalDialInNumbers lists the PSTN dial-in numbers available for the
	// meeting.
	GlobalDialInNumbers []GlobalDialInNumber `json:"global_dial_in_numbers"`
	// HostVideo starts the meeting with the host's video on.
	HostVideo bool `json:"host_video"`
	// JbhTime is the number of minutes before the start time that participants
	// can join before the host (join before host).
	JbhTime int `json:"jbh_time"`
	// JoinBeforeHost allows participants to join before the host arrives.
	JoinBeforeHost bool `json:"join_before_host"`
	// QuestionAndAnswer configures the Q&A feature for the meeting.
	QuestionAndAnswer QuestionAndAnswer `json:"question_and_answer"`
	// LanguageInterpretation configures language interpretation channels.
	LanguageInterpretation LanguageInterpretation `json:"language_interpretation"`
	// SignLanguageInterpretation configures sign language interpretation.
	SignLanguageInterpretation LanguageInterpretation `json:"sign_language_interpretation"`
	// MeetingAuthentication requires authenticated sign-in to join the meeting.
	MeetingAuthentication bool `json:"meeting_authentication"`
	// MuteUponEntry mutes all participants when they join the meeting.
	MuteUponEntry bool `json:"mute_upon_entry"`
	// ParticipantVideo starts participants' video when they join.
	ParticipantVideo bool `json:"participant_video"`
	// PrivateMeeting hides the meeting from public search results.
	PrivateMeeting bool `json:"private_meeting"`
	// RegistrantsConfirmationEmail sends a confirmation email to registrants.
	RegistrantsConfirmationEmail bool `json:"registrants_confirmation_email"`
	// RegistrantsEmailNotification sends an email to registrants before the meeting.
	RegistrantsEmailNotification bool `json:"registrants_email_notification"`
	// RegistrationType controls registration for recurring meetings (1 = register
	// once for all occurrences, 2 = register for each occurrence).
	RegistrationType int `json:"registration_type"`
	// ShowShareButton shows the social share button in the Zoom client for this
	// meeting.
	ShowShareButton bool `json:"show_share_button"`
	// ShowJoinInfo shows the join information panel during the meeting.
	ShowJoinInfo bool `json:"show_join_info"`
	// UsePmi indicates whether the host's Personal Meeting ID is used instead of a
	// generated meeting ID.
	UsePmi bool `json:"use_pmi"`
	// WaitingRoom enables the waiting room so the host can admit participants
	// individually.
	WaitingRoom bool `json:"waiting_room"`
	// WaitingRoomOptions configures waiting room behaviour.
	WaitingRoomOptions WaitingRoomOptions `json:"waiting_room_options"`
	// Watermark adds a watermark to shared screens and participant video.
	Watermark bool `json:"watermark"`
	// HostSaveVideoOrder allows the host to save the participant video order.
	HostSaveVideoOrder bool `json:"host_save_video_order"`
	// InternalMeeting restricts the meeting to users within the same account.
	InternalMeeting bool `json:"internal_meeting"`
	// MeetingInvitees lists specific participants to invite to the meeting.
	MeetingInvitees []MeetingInvitee `json:"meeting_invitees"`
	// ContinuousMeetingChat configures the persistent chat channel attached to the
	// meeting.
	ContinuousMeetingChat ContinuousMeetingChat `json:"continuous_meeting_chat"`
	// ParticipantFocusedMeeting enables participant-focused meeting mode.
	ParticipantFocusedMeeting bool `json:"participant_focused_meeting"`
	// PushChangeToCalendar pushes meeting changes to the host's calendar.
	PushChangeToCalendar bool `json:"push_change_to_calendar"`
	// Resources lists resources (e.g. Zoom Rooms) associated with the meeting.
	Resources []Resource `json:"resources"`
	// AutoStartMeetingSummary automatically starts an AI meeting summary.
	AutoStartMeetingSummary bool `json:"auto_start_meeting_summary"`
	// WhoWillReceiveSummary controls who receives the AI-generated meeting summary.
	WhoWillReceiveSummary int `json:"who_will_receive_summary"`
	// AutoStartAiCompanionQuestions automatically starts AI Companion questions.
	AutoStartAiCompanionQuestions bool `json:"auto_start_ai_companion_questions"`
	// WhoCanAskQuestions controls who can submit questions to AI Companion.
	WhoCanAskQuestions int `json:"who_can_ask_questions"`
	// SummaryTemplateID is the ID of the summary template to use.
	SummaryTemplateID string `json:"summary_template_id"`
	// DeviceTesting enables device testing before the meeting starts.
	DeviceTesting bool `json:"device_testing"`
	// RequestPermissionToUnmuteParticipants requires host approval before a
	// participant can unmute themselves.
	RequestPermissionToUnmuteParticipants bool `json:"request_permission_to_unmute_participants"`
	// AllowHostControlParticipantMuteState allows the host to control whether
	// participants can unmute themselves.
	AllowHostControlParticipantMuteState bool `json:"allow_host_control_participant_mute_state"`
	// DisableParticipantVideo prevents participants from enabling their cameras.
	DisableParticipantVideo bool `json:"disable_participant_video"`
	// EmailInAttendeeReport includes email addresses in the attendee report.
	EmailInAttendeeReport bool `json:"email_in_attendee_report"`
}

// MeetingTrackingField represents a custom tracking key-value pair attached to
// a meeting for reporting purposes.
type MeetingTrackingField struct {
	// Field is the name of the tracking field.
	Field string `json:"field"`
	// Value is the value of the tracking field for this meeting.
	Value string `json:"value"`
	// Visible controls whether the field is visible to participants.
	Visible bool `json:"visible"`
}

// MeetingApprovedOrDeniedCountriesOrRegions restricts or permits meeting access
// based on the participant's country or region.
type MeetingApprovedOrDeniedCountriesOrRegions struct {
	// ApprovedList contains ISO country codes that are allowed to join.
	ApprovedList []string `json:"approved_list"`
	// DeniedList contains ISO country codes that are blocked from joining.
	DeniedList []string `json:"denied_list"`
	// Enable turns the country/region restriction on or off.
	Enable bool `json:"enable"`
	// Method is the restriction method ("approve" or "deny").
	Method string `json:"method"`
}

// MeetingAuthenticationException represents a participant who is exempt from
// the meeting's authentication requirement.
type MeetingAuthenticationException struct {
	// Email is the participant's email address.
	Email string `json:"email"`
	// Name is the participant's display name.
	Name string `json:"name"`
	// JoinURL is the personalised join link for this participant.
	JoinURL string `json:"join_url"`
}

// AutoAddRecordingToVideoManagement controls whether cloud recordings are
// automatically published to Zoom Video Management channels.
type AutoAddRecordingToVideoManagement struct {
	// Enable turns automatic publishing on or off.
	Enable bool `json:"enable"`
	// Channels lists the Video Management channels to publish recordings to.
	Channels []Channel `json:"channels"`
}

// Channel represents a Zoom Video Management channel.
type Channel struct {
	// ChannelID is the unique identifier of the channel.
	ChannelID string `json:"channel_id"`
	// Name is the display name of the channel.
	Name string `json:"name"`
}

// BreakoutRoom configures breakout room behaviour for a meeting.
type BreakoutRoom struct {
	// Enable turns breakout rooms on or off.
	Enable bool `json:"enable"`
	// Rooms lists the pre-configured breakout rooms and their participant
	// assignments.
	Rooms []Room `json:"rooms"`
}

// Room represents a single breakout room within a meeting.
type Room struct {
	// Name is the display name of the breakout room.
	Name string `json:"name"`
	// Participants lists the email addresses of participants pre-assigned to this
	// room.
	Participants []string `json:"participants"`
}

// CustomKey is an arbitrary key-value pair that can be attached to a meeting for
// custom tracking purposes.
type CustomKey struct {
	// Key is the custom field name.
	Key string `json:"key"`
	// Value is the custom field value.
	Value string `json:"value"`
}

// GlobalDialInNumber contains the dial-in details for a specific PSTN location.
type GlobalDialInNumber struct {
	// City is the city associated with this dial-in number.
	City string `json:"city"`
	// Country is the ISO 3166-1 alpha-2 country code.
	Country string `json:"country"`
	// CountryName is the human-readable country name.
	CountryName string `json:"country_name"`
	// Number is the PSTN phone number.
	Number string `json:"number"`
	// Type is the number type (e.g. "toll", "tollfree").
	Type string `json:"type"`
}

// QuestionAndAnswer configures the Q&A feature for a meeting.
type QuestionAndAnswer struct {
	// Enable turns the Q&A feature on or off.
	Enable bool `json:"enable"`
	// AllowSubmitQuestions allows attendees to submit questions.
	AllowSubmitQuestions bool `json:"allow_submit_questions"`
	// AllowAnonymousQuestions allows attendees to submit questions anonymously.
	AllowAnonymousQuestions bool `json:"allow_anonymous_questions"`
	// QuestionVisibility controls who can see submitted questions ("all" or
	// "attendee").
	QuestionVisibility string `json:"question_visibility"`
	// AttendeesCanComment allows attendees to comment on questions.
	AttendeesCanComment bool `json:"attendees_can_comment"`
	// AttendeesCanUpvote allows attendees to upvote questions.
	AttendeesCanUpvote bool `json:"attendees_can_upvote"`
}

// LanguageInterpretation configures interpretation channels for a meeting,
// supporting both spoken language and sign language interpretation.
type LanguageInterpretation struct {
	// Enable turns language interpretation on or off.
	Enable bool `json:"enable"`
	// Interpreters lists the assigned interpreters and their language pairs.
	Interpreters []Interpreter `json:"interpreters"`
}

// Interpreter represents a participant who is designated as a language or sign
// language interpreter for a meeting.
type Interpreter struct {
	// Email is the interpreter's Zoom account email address.
	Email string `json:"email"`
	// InterpreterLanguages is a comma-separated pair of language codes (e.g.
	// "US,FR") for spoken language interpreters.
	InterpreterLanguages string `json:"interpreter_languages,omitempty"`
	// SignLanguage is the sign language the interpreter will interpret (e.g. "ASL")
	// for sign language interpreters.
	SignLanguage string `json:"sign_language,omitempty"`
}

// WaitingRoomOptions configures how the waiting room behaves for a meeting.
type WaitingRoomOptions struct {
	// Mode controls when the waiting room is active ("always", "not_coworker",
	// "have_domain", "not_listed").
	Mode string `json:"mode"`
	// WhoGoesToWaitingRoom specifies which participants are placed in the waiting
	// room ("everyone", "not_coworker", "have_domain", "not_listed").
	WhoGoesToWaitingRoom string `json:"who_goes_to_waiting_room"`
}

// MeetingInvitee represents a participant who has been specifically invited to a
// meeting.
type MeetingInvitee struct {
	// Email is the invitee's email address.
	Email string `json:"email"`
	// InternalUser indicates whether the invitee belongs to the same Zoom account.
	InternalUser bool `json:"internal_user"`
}

// ContinuousMeetingChat configures the persistent chat channel that stays active
// between meeting occurrences.
type ContinuousMeetingChat struct {
	// Enable turns the continuous meeting chat channel on or off.
	Enable bool `json:"enable"`
	// ChannelID is the ID of the associated persistent chat channel.
	ChannelID string `json:"channel_id"`
}

// Resource represents an additional resource (such as a Zoom Room) associated
// with a meeting.
type Resource struct {
	// ResourceType is the type of resource (e.g. "zoom_room").
	ResourceType string `json:"resource_type"`
	// ResourceID is the unique identifier of the resource.
	ResourceID string `json:"resource_id"`
	// PermissionLevel is the permission level granted to this resource.
	PermissionLevel string `json:"permission_level"`
}

// MeetingSummary represents an AI-generated summary for a Zoom meeting, as
// returned by the meeting summaries endpoint.
type MeetingSummary struct {
	// MeetingEndTime is the time the meeting ended.
	MeetingEndTime time.Time `json:"meeting_end_time"`
	// MeetingHostEmail is the email address of the meeting host.
	MeetingHostEmail string `json:"meeting_host_email"`
	// MeetingHostID is the user ID of the meeting host.
	MeetingHostID string `json:"meeting_host_id"`
	// MeetingID is the numeric meeting ID.
	MeetingID int `json:"meeting_id"`
	// MeetingStartTime is the time the meeting started.
	MeetingStartTime time.Time `json:"meeting_start_time"`
	// MeetingTopic is the title of the meeting.
	MeetingTopic string `json:"meeting_topic"`
	// MeetingUUID is the unique UUID of the meeting occurrence.
	MeetingUUID string `json:"meeting_uuid"`
	// SummaryContent is the full text of the AI-generated summary.
	SummaryContent string `json:"summary_content,omitempty"`
	// SummaryCreatedTime is the time the summary was generated.
	SummaryCreatedTime time.Time `json:"summary_created_time"`
	// SummaryDocURL is the URL of the summary document, if available.
	SummaryDocURL string `json:"summary_doc_url,omitempty"`
	// SummaryEndTime is the end of the time range covered by the summary.
	SummaryEndTime time.Time `json:"summary_end_time"`
	// SummaryLastModifiedTime is the time the summary was last edited.
	SummaryLastModifiedTime time.Time `json:"summary_last_modified_time"`
	// SummaryLastModifiedUserEmail is the email of the last user to edit the
	// summary.
	SummaryLastModifiedUserEmail string `json:"summary_last_modified_user_email"`
	// SummaryLastModifiedUserID is the ID of the last user to edit the summary.
	SummaryLastModifiedUserID string `json:"summary_last_modified_user_id"`
	// SummaryStartTime is the start of the time range covered by the summary.
	SummaryStartTime time.Time `json:"summary_start_time"`
	// SummaryTitle is the title of the generated summary document.
	SummaryTitle string `json:"summary_title,omitempty"`
}
