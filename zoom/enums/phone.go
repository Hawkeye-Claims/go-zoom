package enums

type ConnectType string

const (
	Internal ConnectType = "internal"
	External ConnectType = "external"
)

type Direction string

const (
	Inbound  Direction = "inbound"
	Outbound Direction = "outbound"
)

type NumberType string

const (
	ZoomPSTN                NumberType = "zoom_pstn"
	ZoomTollFreeNumber      NumberType = "zoom_toll_free_number"
	ExternalPSTN            NumberType = "external_pstn"
	BYOC                    NumberType = "byoc"
	BYOP                    NumberType = "byop"
	ThirdPartyContactCenter NumberType = "3rd_party_contact_center"
	ZoomServiceNumber       NumberType = "zoom_service_number"
	ExternalServiceNumber   NumberType = "external_service_number"
	ZoomContactCenter       NumberType = "zoom_contact_center"
	MeetingPhoneNumber      NumberType = "meeting_phone_number"
	MeetingID               NumberType = "meeting_id"
	AnonymousNumber         NumberType = "anonymous_number"
	ZoomRevenueAccelerator  NumberType = "zoom_revenue_accelerator"
)

type CallType string

const (
	General   CallType = "general"
	Emergency CallType = "emergency"
)

type ExtensionType string

const (
	User             ExtensionType = "user"
	CallQueue        ExtensionType = "call_queue"
	AutoReceptionist ExtensionType = "auto_receptionist"
	CommonArea       ExtensionType = "common_area"
	ZoomRoom         ExtensionType = "zoom_room"
	CiscoRoom        ExtensionType = "cisco_room"
	SharedLineGroup  ExtensionType = "shared_line_group"
	GroupCallPickup  ExtensionType = "group_call_pickup"
	ExternalContact  ExtensionType = "external_contact"
)

type CallResult string

const (
	Answered           CallResult = "answered"
	Accepted           CallResult = "accepted"
	PickedUp           CallResult = "picked_up"
	Connected          CallResult = "connected"
	Succeeded          CallResult = "succeeded"
	Voicemail          CallResult = "voicemail"
	HangUp             CallResult = "hang_up"
	Canceled           CallResult = "canceled"
	CallFailed         CallResult = "call_failed"
	Unconnected        CallResult = "unconnected"
	Rejected           CallResult = "rejected"
	Busy               CallResult = "busy"
	RingTimeout        CallResult = "ring_timeout"
	Overflow           CallResult = "overflow"
	NoAnswer           CallResult = "no_answer"
	InvalidKey         CallResult = "invalid_key"
	InvalidOperation   CallResult = "invalid_operation"
	Abandoned          CallResult = "abandoned"
	SystemBlocked      CallResult = "system_blocked"
	ServiceUnavailable CallResult = "service_unavailable"
)

type TimeType string

const (
	StartTime TimeType = "start_time"
	EndTime   TimeType = "end_time"
)

type RecordingStatus string

const (
	Recorded    RecordingStatus = "recorded"
	NonRecorded RecordingStatus = "non_recorded"
)
