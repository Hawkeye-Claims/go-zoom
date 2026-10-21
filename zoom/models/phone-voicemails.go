package models

import "time"

// PhoneVoicemail contains the metadata returned by owner voicemail lists and
// voicemail detail. Transcription is absent when Zoom does not return it.
type PhoneVoicemail struct {
	// ID identifies the voicemail independently of its associated call.
	ID string `json:"id,omitempty"`
	// CallElementID identifies the call element that produced the voicemail.
	CallElementID string `json:"call_element_id,omitempty"`
	// CallHistoryID identifies the associated call history.
	CallHistoryID string `json:"call_history_id,omitempty"`
	// CallID identifies the associated phone call.
	CallID string `json:"call_id,omitempty"`
	// CallLogID is the legacy call log identifier, when returned.
	CallLogID string `json:"call_log_id,omitempty"`
	// CallerName is the caller's display name.
	CallerName string `json:"caller_name,omitempty"`
	// CallerNumber is the caller's phone number or extension.
	CallerNumber string `json:"caller_number,omitempty"`
	// CallerNumberType is 1 for internal or 2 for external numbers.
	CallerNumberType int `json:"caller_number_type,omitempty"`
	// CallerAccountCode identifies the caller's account for inter-account dialing.
	CallerAccountCode string `json:"caller_account_code,omitempty"`
	// CalleeName is the recipient's display name.
	CalleeName string `json:"callee_name,omitempty"`
	// CalleeNumber is the recipient's phone number or extension.
	CalleeNumber string `json:"callee_number,omitempty"`
	// CalleeNumberType is 1 internal, 2 external, or 3 customized emergency.
	CalleeNumberType int `json:"callee_number_type,omitempty"`
	// CalleeAccountCode identifies the recipient's account.
	CalleeAccountCode string `json:"callee_account_code,omitempty"`
	// DateTime is the voicemail creation instant.
	DateTime time.Time `json:"date_time"`
	// Duration is the voicemail length in seconds.
	Duration int `json:"duration"`
	// Status is read or unread, distinct from transcription processing status.
	Status string `json:"status,omitempty"`
	// DownloadURL requires bearer authentication. Reading metadata or a
	// transcription does not download audio or change the voicemail's status.
	DownloadURL string `json:"download_url,omitempty"`
	// Transcription is included by detail reads when Zoom returns it.
	Transcription *PhoneVoicemailTranscription `json:"transcription,omitempty"`
}

// PhoneVoicemailTranscription preserves Zoom's raw processing status, including
// unknown future values. Status 0 is unavailable, 1 processing, 2 successful,
// and 4 disabled. A missing transcription is distinct from status 0.
type PhoneVoicemailTranscription struct {
	// Content is the transcribed text, if available.
	Content string `json:"content"`
	// Status is Zoom's numeric transcription state.
	Status int `json:"status"`
	// Engine identifies the transcription provider.
	Engine string `json:"engine,omitempty"`
}
