package server

import (
	"github.com/TheSlowpes/go-zoom/zoom/enums"
	"github.com/TheSlowpes/go-zoom/zoom/models"
)

type Notification[T any] struct {
	Event   string `json:"event"`
	EventTs int64  `json:"event_ts"`
	Payload T      `json:"payload"`
}

type webhookHeader struct {
	Event string `json:"event"`
}

type validateTokenPayload struct {
	PlainToken string `json:"plain_token"`
}

type validateTokenResponse struct {
	PlainToken     string `json:"plain_token"`
	EncryptedToken string `json:"encrypted_token"`
}

type MeetingEvent struct {
	AccountID  string         `json:"account_id"`
	Object     models.Meeting `json:"object"`
	Operator   string         `json:"operator,omitempty"`
	OperatorID string         `json:"operator_id,omitempty"`
	Operation  string         `json:"operation,omitempty"`
}

type UserEvent struct {
	AccountID    string                 `json:"account_id"`
	Operator     string                 `json:"operator,omitempty"`
	OperatorID   string                 `json:"operator_id,omitempty"`
	Object       models.User            `json:"object"`
	CreationType enums.UserCreateAction `json:"creation_type,omitempty"`
}

type AICallSummaryEvent struct {
	AccountID string               `json:"account_id"`
	Object    models.AICallSummary `json:"object"`
}

type PhoneCallElementEvent struct {
	AccountID string `json:"account_id"`
	Object    struct {
		CallElements []models.CallElement `json:"call_elements"`
	} `json:"object"`
	UserID string `json:"user_id"`
}

type PhoneCallHistoryEvent struct {
	AccountID string `json:"account_id"`
	Object    struct {
		CallLogs []models.CallHistory `json:"call_logs"`
	}
	UserID string `json:"user_id"`
}
