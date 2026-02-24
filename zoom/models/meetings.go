package models

import "time"

type Meeting struct {
	AssistantID       string                 `json:"assistant_id"`
	HostEmail         string                 `json:"host_email"`
	HostID            string                 `json:"host_id"`
	ID                int64                  `json:"id"`
	UUID              string                 `json:"uuid"`
	RegistrationURL   string                 `json:"registration_url"`
	Agenda            string                 `json:"agenda"`
	CreatedAt         time.Time              `json:"created_at"`
	Duration          int                    `json:"duration"`
	EncryptedPassword string                 `json:"encrypted_password"`
	PstnPassword      string                 `json:"pstn_password"`
	H323Password      string                 `json:"h323_password"`
	JoinURL           string                 `json:"join_url"`
	ChatJoinURL       string                 `json:"chat_join_url"`
	Occurrences       []MeetingOccurrence    `json:"occurrences"`
	Password          string                 `json:"password"`
	Pmi               string                 `json:"pmi"`
	PreSchedule       bool                   `json:"pre_schedule"`
	Recurrence        MeetingRecurrence      `json:"recurrence"`
	Settings          MeetingSettings        `json:"settings"`
	StartTime         time.Time              `json:"start_time"`
	StartURL          string                 `json:"start_url"`
	Status            string                 `json:"status"`
	TemplateID        string                 `json:"template_id"`
	Timezone          string                 `json:"timezone"`
	Topic             string                 `json:"topic"`
	TrackingFields    []MeetingTrackingField `json:"tracking_fields"`
	Type              int                    `json:"type"`
	DynamicHostKey    string                 `json:"dynamic_host_key"`
	CreationSource    string                 `json:"creation_source"`
}

type MeetingOccurrence struct {
	Duration     int       `json:"duration"`
	OccurrenceID string    `json:"occurrence_id"`
	StartTime    time.Time `json:"start_time"`
	Status       string    `json:"status"`
}

type MeetingRecurrence struct {
	EndDateTime    time.Time `json:"end_date_time"`
	EndTimes       int       `json:"end_times"`
	MonthlyDay     int       `json:"monthly_day"`
	MonthlyWeek    int       `json:"monthly_week"`
	MonthlyWeekDay int       `json:"monthly_week_day"`
	RepeatInterval int       `json:"repeat_interval"`
	Type           int       `json:"type"`
	WeeklyDays     string    `json:"weekly_days"`
}

type MeetingSettings struct {
	AdditionalDataCenterRegions           []string                                  `json:"additional_data_center_regions"`
	AllowMultipleDevices                  bool                                      `json:"allow_multiple_devices"`
	AlternativeHosts                      string                                    `json:"alternative_hosts"`
	AlternativeHostsEmailNotification     bool                                      `json:"alternative_hosts_email_notification"`
	AlternativeHostUpdatePolls            bool                                      `json:"alternative_host_update_polls"`
	AlternativeHostManageMeetingSummary   bool                                      `json:"alternative_host_manage_meeting_summary"`
	AlternativeHostManageCloudRecording   bool                                      `json:"alternative_host_manage_cloud_recording"`
	ApprovalType                          int                                       `json:"approval_type"`
	ApprovedOrDeniedCountriesOrRegions    MeetingApprovedOrDeniedCountriesOrRegions `json:"approved_or_denied_countries_or_regions"`
	Audio                                 string                                    `json:"audio"`
	AudioConferenceInfo                   string                                    `json:"audio_conference_info"`
	AuthenticationDomains                 string                                    `json:"authentication_domains"`
	AuthenticationException               []MeetingAuthenticationException          `json:"authentication_exception"`
	AuthenticationName                    string                                    `json:"authentication_name"`
	AuthenticationOption                  string                                    `json:"authentication_option"`
	AutoRecording                         string                                    `json:"auto_recording"`
	AutoAddRecordingToVideoManagement     AutoAddRecordingToVideoManagement         `json:"auto_add_recording_to_video_management"`
	BreakoutRoom                          BreakoutRoom                              `json:"breakout_room"`
	CalendarType                          int                                       `json:"calendar_type"`
	CloseRegistration                     bool                                      `json:"close_registration"`
	ContactEmail                          string                                    `json:"contact_email"`
	ContactName                           string                                    `json:"contact_name"`
	CustomKeys                            []CustomKey                               `json:"custom_keys"`
	EmailNotification                     bool                                      `json:"email_notification"`
	EncryptionType                        string                                    `json:"encryption_type"`
	FocusMode                             bool                                      `json:"focus_mode"`
	GlobalDialInCountries                 []string                                  `json:"global_dial_in_countries"`
	GlobalDialInNumbers                   []GlobalDialInNumber                      `json:"global_dial_in_numbers"`
	HostVideo                             bool                                      `json:"host_video"`
	JbhTime                               int                                       `json:"jbh_time"`
	JoinBeforeHost                        bool                                      `json:"join_before_host"`
	QuestionAndAnswer                     QuestionAndAnswer                         `json:"question_and_answer"`
	LanguageInterpretation                LanguageInterpretation                    `json:"language_interpretation"`
	SignLanguageInterpretation            LanguageInterpretation                    `json:"sign_language_interpretation"`
	MeetingAuthentication                 bool                                      `json:"meeting_authentication"`
	MuteUponEntry                         bool                                      `json:"mute_upon_entry"`
	ParticipantVideo                      bool                                      `json:"participant_video"`
	PrivateMeeting                        bool                                      `json:"private_meeting"`
	RegistrantsConfirmationEmail          bool                                      `json:"registrants_confirmation_email"`
	RegistrantsEmailNotification          bool                                      `json:"registrants_email_notification"`
	RegistrationType                      int                                       `json:"registration_type"`
	ShowShareButton                       bool                                      `json:"show_share_button"`
	ShowJoinInfo                          bool                                      `json:"show_join_info"`
	UsePmi                                bool                                      `json:"use_pmi"`
	WaitingRoom                           bool                                      `json:"waiting_room"`
	WaitingRoomOptions                    WaitingRoomOptions                        `json:"waiting_room_options"`
	Watermark                             bool                                      `json:"watermark"`
	HostSaveVideoOrder                    bool                                      `json:"host_save_video_order"`
	InternalMeeting                       bool                                      `json:"internal_meeting"`
	MeetingInvitees                       []MeetingInvitee                          `json:"meeting_invitees"`
	ContinuousMeetingChat                 ContinuousMeetingChat                     `json:"continuous_meeting_chat"`
	ParticipantFocusedMeeting             bool                                      `json:"participant_focused_meeting"`
	PushChangeToCalendar                  bool                                      `json:"push_change_to_calendar"`
	Resources                             []Resource                                `json:"resources"`
	AutoStartMeetingSummary               bool                                      `json:"auto_start_meeting_summary"`
	WhoWillReceiveSummary                 int                                       `json:"who_will_receive_summary"`
	AutoStartAiCompanionQuestions         bool                                      `json:"auto_start_ai_companion_questions"`
	WhoCanAskQuestions                    int                                       `json:"who_can_ask_questions"`
	SummaryTemplateID                     string                                    `json:"summary_template_id"`
	DeviceTesting                         bool                                      `json:"device_testing"`
	RequestPermissionToUnmuteParticipants bool                                      `json:"request_permission_to_unmute_participants"`
	AllowHostControlParticipantMuteState  bool                                      `json:"allow_host_control_participant_mute_state"`
	DisableParticipantVideo               bool                                      `json:"disable_participant_video"`
	EmailInAttendeeReport                 bool                                      `json:"email_in_attendee_report"`
}

type MeetingTrackingField struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Visible bool   `json:"visible"`
}

type MeetingApprovedOrDeniedCountriesOrRegions struct {
	ApprovedList []string `json:"approved_list"`
	DeniedList   []string `json:"denied_list"`
	Enable       bool     `json:"enable"`
	Method       string   `json:"method"`
}

type MeetingAuthenticationException struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	JoinURL string `json:"join_url"`
}

type AutoAddRecordingToVideoManagement struct {
	Enable   bool      `json:"enable"`
	Channels []Channel `json:"channels"`
}

type Channel struct {
	ChannelID string `json:"channel_id"`
	Name      string `json:"name"`
}

type BreakoutRoom struct {
	Enable bool   `json:"enable"`
	Rooms  []Room `json:"rooms"`
}

type Room struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

type CustomKey struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GlobalDialInNumber struct {
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryName string `json:"country_name"`
	Number      string `json:"number"`
	Type        string `json:"type"`
}

type QuestionAndAnswer struct {
	Enable                  bool   `json:"enable"`
	AllowSubmitQuestions    bool   `json:"allow_submit_questions"`
	AllowAnonymousQuestions bool   `json:"allow_anonymous_questions"`
	QuestionVisibility      string `json:"question_visibility"`
	AttendeesCanComment     bool   `json:"attendees_can_comment"`
	AttendeesCanUpvote      bool   `json:"attendees_can_upvote"`
}

type LanguageInterpretation struct {
	Enable       bool          `json:"enable"`
	Interpreters []Interpreter `json:"interpreters"`
}

type Interpreter struct {
	Email                string `json:"email"`
	InterpreterLanguages string `json:"interpreter_languages,omitempty"`
	SignLanguage         string `json:"sign_language,omitempty"`
}

type WaitingRoomOptions struct {
	Mode                 string `json:"mode"`
	WhoGoesToWaitingRoom string `json:"who_goes_to_waiting_room"`
}

type MeetingInvitee struct {
	Email        string `json:"email"`
	InternalUser bool   `json:"internal_user"`
}

type ContinuousMeetingChat struct {
	Enable    bool   `json:"enable"`
	ChannelID string `json:"channel_id"`
}

type Resource struct {
	ResourceType    string `json:"resource_type"`
	ResourceID      string `json:"resource_id"`
	PermissionLevel string `json:"permission_level"`
}
