package enums

// ConnectType indicates whether a phone call was routed internally within the
// Zoom account or externally through the PSTN.
type ConnectType string

const (
	// Internal indicates the call was between two extensions on the same account.
	Internal ConnectType = "internal"
	// External indicates the call involved a number outside the account.
	External ConnectType = "external"
)

// Direction indicates whether a call was received or placed by the user.
type Direction string

const (
	// Inbound indicates the user received the call.
	Inbound Direction = "inbound"
	// Outbound indicates the user placed the call.
	Outbound Direction = "outbound"
)

// NumberType classifies the type of phone number involved in a call.
type NumberType string

const (
	// ZoomPSTN is a Zoom-provided PSTN number.
	ZoomPSTN NumberType = "zoom_pstn"
	// ZoomTollFreeNumber is a Zoom-provided toll-free number.
	ZoomTollFreeNumber NumberType = "zoom_toll_free_number"
	// ExternalPSTN is an external PSTN number not managed by Zoom.
	ExternalPSTN NumberType = "external_pstn"
	// BYOC is a Bring Your Own Carrier number.
	BYOC NumberType = "byoc"
	// BYOP is a Bring Your Own PSTN number.
	BYOP NumberType = "byop"
	// ThirdPartyContactCenter is a number associated with a third-party contact center.
	ThirdPartyContactCenter NumberType = "3rd_party_contact_center"
	// ZoomServiceNumber is an internal Zoom service number.
	ZoomServiceNumber NumberType = "zoom_service_number"
	// ExternalServiceNumber is an external service number.
	ExternalServiceNumber NumberType = "external_service_number"
	// ZoomContactCenter is a Zoom Contact Center number.
	ZoomContactCenter NumberType = "zoom_contact_center"
	// MeetingPhoneNumber is a dial-in number for a Zoom meeting.
	MeetingPhoneNumber NumberType = "meeting_phone_number"
	// MeetingID is a numeric Zoom meeting ID used for dial-in.
	MeetingID NumberType = "meeting_id"
	// AnonymousNumber indicates the caller's number was withheld.
	AnonymousNumber NumberType = "anonymous_number"
	// ZoomRevenueAccelerator is a number used with Zoom Revenue Accelerator.
	ZoomRevenueAccelerator NumberType = "zoom_revenue_accelerator"
)

// CallType classifies a phone call as a general or emergency call.
type CallType string

const (
	// General is a standard, non-emergency phone call.
	General CallType = "general"
	// Emergency is a call placed to an emergency service (e.g. 911).
	Emergency CallType = "emergency"
)

// ExtensionType identifies the kind of Zoom Phone extension involved in a call.
type ExtensionType string

const (
	// User is an individual Zoom Phone user extension.
	User ExtensionType = "user"
	// CallQueue is a call queue extension.
	CallQueue ExtensionType = "call_queue"
	// AutoReceptionist is an auto-receptionist (IVR) extension.
	AutoReceptionist ExtensionType = "auto_receptionist"
	// CommonArea is a common area phone extension.
	CommonArea ExtensionType = "common_area"
	// ZoomRoom is a Zoom Rooms extension.
	ZoomRoom ExtensionType = "zoom_room"
	// CiscoRoom is a Cisco room system extension.
	CiscoRoom ExtensionType = "cisco_room"
	// SharedLineGroup is a shared line group extension.
	SharedLineGroup ExtensionType = "shared_line_group"
	// GroupCallPickup is a group call pickup extension.
	GroupCallPickup ExtensionType = "group_call_pickup"
	// ExternalContact is an external contact extension.
	ExternalContact ExtensionType = "external_contact"
)

// CallResult describes the final outcome of a phone call.
type CallResult string

const (
	// Answered indicates the call was answered.
	Answered CallResult = "answered"
	// Accepted indicates the call was accepted.
	Accepted CallResult = "accepted"
	// PickedUp indicates the call was picked up.
	PickedUp CallResult = "picked_up"
	// Connected indicates the call was connected.
	Connected CallResult = "connected"
	// Succeeded indicates the call completed successfully.
	Succeeded CallResult = "succeeded"
	// Voicemail indicates the call was sent to voicemail.
	Voicemail CallResult = "voicemail"
	// HangUp indicates the caller hung up before the call was answered.
	HangUp CallResult = "hang_up"
	// Canceled indicates the call was canceled before connecting.
	Canceled CallResult = "canceled"
	// CallFailed indicates the call failed to connect.
	CallFailed CallResult = "call_failed"
	// Unconnected indicates the call did not connect for an unspecified reason.
	Unconnected CallResult = "unconnected"
	// Rejected indicates the callee rejected the call.
	Rejected CallResult = "rejected"
	// Busy indicates the callee's line was busy.
	Busy CallResult = "busy"
	// RingTimeout indicates the call rang without being answered until the timeout.
	RingTimeout CallResult = "ring_timeout"
	// Overflow indicates the call was redirected due to an overflow condition.
	Overflow CallResult = "overflow"
	// NoAnswer indicates the call was not answered.
	NoAnswer CallResult = "no_answer"
	// InvalidKey indicates an invalid key was pressed in an IVR menu.
	InvalidKey CallResult = "invalid_key"
	// InvalidOperation indicates an invalid operation occurred during the call.
	InvalidOperation CallResult = "invalid_operation"
	// Abandoned indicates the caller disconnected while waiting in a queue.
	Abandoned CallResult = "abandoned"
	// SystemBlocked indicates the call was blocked by a system rule.
	SystemBlocked CallResult = "system_blocked"
	// ServiceUnavailable indicates the telephony service was unavailable.
	ServiceUnavailable CallResult = "service_unavailable"
)

// TimeType selects which timestamp is used when filtering phone call history by
// a date range.
type TimeType string

const (
	// StartTime filters by the time the call started.
	StartTime TimeType = "start_time"
	// EndTime filters by the time the call ended.
	EndTime TimeType = "end_time"
)

// RecordingStatus indicates whether a call was recorded.
type RecordingStatus string

const (
	// Recorded indicates the call has an associated recording.
	Recorded RecordingStatus = "recorded"
	// NonRecorded indicates the call was not recorded.
	NonRecorded RecordingStatus = "non_recorded"
)

// CallForwardingType controls which destinations are permitted when a user
// forwards calls.
type CallForwardingType int

const (
	// LowRestriction prohibits forwarding to external numbers.
	LowRestriction CallForwardingType = 1 //External numbers not allowed
	// MediumRestriction prohibits forwarding to external numbers and external contacts.
	MediumRestriction CallForwardingType = 2 //External numbers and external contacts not allowed
	// HighRestriction prohibits forwarding to external numbers, external contacts,
	// and internal extensions without inbound automatic call recording.
	HighRestriction CallForwardingType = 3 //External numbers, external contacts and internal extension without inbound automatic call recording not allowed
	// NoRestriction places no restrictions on call forwarding destinations.
	NoRestriction CallForwardingType = 4 //No restrictions on call forwarding
)

// CallTransferType controls which destinations are permitted when a user
// transfers a call.
type CallTransferType int

const (
	// NoRestrictionCallTransfer places no restrictions on call transfer destinations.
	NoRestrictionCallTransfer CallTransferType = 1
	// MediumRestrictionCallTransfer applies a medium level of restriction.
	MediumRestrictionCallTransfer CallTransferType = 2
	// HighRestrictionCallTransfer applies a high level of restriction.
	HighRestrictionCallTransfer CallTransferType = 3
	// LowRestrictionCallTransfer applies a low level of restriction.
	LowRestrictionCallTransfer CallTransferType = 4
)

// CallNotPickedUpAction determines what happens when a parked call is not
// retrieved within the expiration period.
type CallNotPickedUpAction int

const (
	// ForwardToVoicemail sends the unanswered call to voicemail.
	ForwardToVoicemail CallNotPickedUpAction = 0
	// Disconnect disconnects the unanswered call.
	Disconnect CallNotPickedUpAction = 9
	// ForwardToExtension forwards the unanswered call to a specified extension.
	ForwardToExtension CallNotPickedUpAction = 50
	// RingBackToParker rings back the user who originally parked the call.
	RingBackToParker CallNotPickedUpAction = 100
)

// EventType classifies the type of event that occurred during a phone call.
type EventType string

const (
	// Incoming indicates an incoming call event.
	Incoming EventType = "incoming"
	// TransferFromZoomContactCenter indicates a call transfer from Zoom Contact Center.
	TransferFromZoomContactCenter EventType = "transfer_from_zoom_contact_center"
	// SharedLineIncoming indicates an incoming call to a shared line group extension.
	SharedLineIncoming EventType = "shared_line_incoming"
	// Outgoing indicates an outgoing call event.
	Outgoing EventType = "outgoing"
	// CallMeOn indicates a Call Me On event where the user initiates a call to a specified number.
	CallMeOn EventType = "call_me_on"
	// OutgoingToZoomContactCenter indicates an outgoing call transfer to Zoom Contact Center.
	OutgoingToZoomContactCenter EventType = "outgoing_to_zoom_contact_center"
	// WarmTransfer indicates a warm transfer event where the user transfers an active call to another destination after speaking to the recipient.
	WarmTransfer EventType = "warm_transfer"
	// Forward indicates a call forwarding event where the user forwards an active call to another destination without speaking to the recipient first.
	Forward EventType = "forward"
	// RingToMember indicates a call event where the call is ringing to a specific member (e.g. in a call queue).
	RingToMember EventType = "ring_to_member"
	// Overflow indicates a call event where the call is redirected due to an overflow condition (e.g. all agents are busy in a call queue).
	EventOverflow EventType = "overflow"
	// DirectTransfer indicates a direct transfer event where the user transfers an active call to another destination without speaking to the recipient first,
	// and the call is not classified as a "Forward" event (e.g. transferring a call from one extension to another within the same account).
	DirectTransfer EventType = "direct_transfer"
	// Barge indicates a barge event where the user barges into an active call between other parties (e.g. in a call queue or shared line group).
	Barge EventType = "barge"
	// Monitor indicates a monitor event where the user listens in on an active call between other parties without participating (e.g. in a call queue or shared line group).
	Monitor EventType = "monitor"
	// Whisper indicates a whisper event where the user speaks to one party in an active call between other parties without the other parties hearing (e.g. in a call queue or shared line group).
	Whisper EventType = "whisper"
	// Listen indicates a listen event where the user listens in on an active call between other parties and can speak to all parties (e.g. in a call queue or shared line group).
	Listen EventType = "listen"
	// Takeover indicates a takeover event where the user takes over an active call between other parties, disconnecting the original parties from the call (e.g. in a call queue or shared line group).
	Takeover EventType = "takeover"
	// ConferenceBarge indicates a conference barge event where the user barges into an active conference call between other parties (e.g. in a call queue or shared line group).
	ConferenceBarge EventType = "conference_barge"
	// Park indicates a call park event where the user parks an active call, making it available for retrieval by other users (e.g. in a call queue or shared line group
	Park EventType = "park"
	// Timeout indicates a call event where the call has reached a timeout condition (e.g. ringing timeout, hold timeout, park timeout).
	Timeout EventType = "timeout"
	// ParkPickUp indicates a call event where a parked call is picked up by a user (e.g. in a call queue or shared line group
	ParkPickUp EventType = "park_pick_up"
	// Merge indicates a call event where the user merges an active call with another call, creating a conference call between all parties (e.g. in a call queue or shared line group).
	Merge EventType = "merge"
	// Shared indicates a call event where the user shares content (e.g. screen sharing, file sharing) during an active call (e.g. in a call queue or shared line group).
	Shared EventType = "shared"
)
