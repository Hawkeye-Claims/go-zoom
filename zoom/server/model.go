// Package server provides an HTTP webhook server for receiving and verifying
// Zoom event notifications. It handles HMAC-SHA256 signature validation,
// URL validation challenges, and dispatches typed event payloads to
// registered channel-based handlers.
package server

import (
	"github.com/TheSlowpes/go-zoom/zoom/enums"
	"github.com/TheSlowpes/go-zoom/zoom/models"
)

// Notification is a generic wrapper for all Zoom webhook event payloads. T is
// the concrete payload type for a specific event (e.g. MeetingEvent,
// UserEvent).
type Notification[T any] struct {
	// Event is the Zoom event type string (e.g. "meeting.started").
	Event string `json:"event"`
	// EventTs is the Unix timestamp (milliseconds) at which the event occurred.
	EventTs int64 `json:"event_ts"`
	// Payload contains the event-specific data.
	Payload T `json:"payload"`
}

// webhookHeader is used internally to unmarshal just the event type field from
// an incoming webhook body before dispatching to the appropriate handler.
type webhookHeader struct {
	Event string `json:"event"`
}

// validateTokenPayload is the payload received in an
// "endpoint.url_validation" event from Zoom.
type validateTokenPayload struct {
	PlainToken string `json:"plain_token"`
}

// validateTokenResponse is the response body sent back to Zoom in reply to
// an "endpoint.url_validation" challenge.
type validateTokenResponse struct {
	PlainToken     string `json:"plain_token"`
	EncryptedToken string `json:"encrypted_token"`
}

// MeetingEvent is the payload type for meeting-related Zoom webhook events
// (e.g. "meeting.started", "meeting.ended").
type MeetingEvent struct {
	// AccountID is the Zoom account that owns the meeting.
	AccountID string `json:"account_id"`
	// Object is the meeting object associated with the event.
	Object models.Meeting `json:"object"`
	// Operator is the email or user ID of the user who triggered the event,
	// if applicable.
	Operator string `json:"operator,omitempty"`
	// OperatorID is the Zoom user ID of the operator, if applicable.
	OperatorID string `json:"operator_id,omitempty"`
	// Operation describes the action that triggered the event (e.g. "update").
	Operation string `json:"operation,omitempty"`
}

// UserEvent is the payload type for user-related Zoom webhook events
// (e.g. "user.created", "user.deleted").
type UserEvent struct {
	// AccountID is the Zoom account that owns the user.
	AccountID string `json:"account_id"`
	// Operator is the email or user ID of the admin who performed the action,
	// if applicable.
	Operator string `json:"operator,omitempty"`
	// OperatorID is the Zoom user ID of the operator, if applicable.
	OperatorID string `json:"operator_id,omitempty"`
	// Object is the user associated with the event.
	Object models.User `json:"object"`
	// CreationType indicates how the user was created (e.g. via SSO or
	// invitation), if applicable.
	CreationType enums.UserCreateAction `json:"creation_type,omitempty"`
}

// AICallSummaryEvent is the payload type for AI call summary webhook events
// (e.g. "phone.ai_call_summary_completed").
type AICallSummaryEvent struct {
	// AccountID is the Zoom account associated with the event.
	AccountID string `json:"account_id"`
	// Object contains the AI call summary data.
	Object models.AICallSummary `json:"object"`
}

// PhoneCallElementEvent is the payload type for phone call-element webhook
// events (e.g. "phone.callee_missed").
type PhoneCallElementEvent struct {
	// AccountID is the Zoom account associated with the event.
	AccountID string `json:"account_id"`
	// Object contains the list of call elements associated with the event.
	Object struct {
		CallElements []models.CallElement `json:"call_elements"`
	} `json:"object"`
	// UserID is the Zoom user ID associated with the event.
	UserID string `json:"user_id"`
}

// PhoneCallHistoryEvent is the payload type for phone call-history webhook
// events (e.g. "phone.call_history_completed").
type PhoneCallHistoryEvent struct {
	// AccountID is the Zoom account associated with the event.
	AccountID string `json:"account_id"`
	// Object contains the list of call log entries associated with the event.
	Object struct {
		CallLogs []models.CallHistory `json:"call_logs"`
	} `json:"object"`
	// UserID is the Zoom user ID associated with the event.
	UserID string `json:"user_id"`
}
