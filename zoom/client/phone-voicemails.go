package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Hawkeye-Claims/go-zoom/zoom/models"
)

// PhoneVoicemailsServicer provides read-only owner voicemail operations.
type PhoneVoicemailsServicer interface {
	ListUserPage(context.Context, string, *PhoneVoicemailQueryParameters, *PaginationOptions) (*PhoneVoicemailPage, *http.Response, error)
	Get(context.Context, string) (*models.PhoneVoicemail, *http.Response, error)
}

// PhoneVoicemailsService reads owner voicemail metadata and transcriptions.
type PhoneVoicemailsService struct {
	client *Client
}

var _ PhoneVoicemailsServicer = (*PhoneVoicemailsService)(nil)

// PhoneVoicemailQueryParameters filters an owner's voicemail list.
type PhoneVoicemailQueryParameters struct {
	// From and To accept YYYY-MM-DD or UTC second-resolution timestamps.
	// Zoom permits one month per query within the last six months.
	From string `url:"from,omitempty"`
	// To is required when From is set and must use the same format.
	To string `url:"to,omitempty"`
	// Status accepts all, read, or unread. Zoom defaults to all.
	Status string `url:"status,omitempty"`
	// Trash queries voicemail from trash when true.
	Trash bool `url:"trash,omitempty"`
}

// PhoneVoicemailPage is one page, not a complete owner mailbox.
type PhoneVoicemailPage struct {
	PaginationResponse
	// From and To are the query bounds echoed by Zoom.
	From string `json:"from"`
	To   string `json:"to"`
	// VoiceMails contains only this page's records.
	VoiceMails []*models.PhoneVoicemail `json:"voice_mails"`
}

// ListUserPage reads one page for a user ID or email. It requires
// phone:read:list_voicemails:admin for S2S apps. Pass the same filters and the
// returned NextPageToken on subsequent calls. Tokens expire after 15 minutes;
// errors must not be treated as end-of-data. The response body is closed.
func (p *PhoneVoicemailsService) ListUserPage(ctx context.Context, userID string, params *PhoneVoicemailQueryParameters, pagination *PaginationOptions) (*PhoneVoicemailPage, *http.Response, error) {
	if userID == "" {
		return nil, nil, fmt.Errorf("user ID is required")
	}
	query := &struct {
		*PhoneVoicemailQueryParameters
		*PaginationOptions
	}{params, pagination}
	page := &PhoneVoicemailPage{}
	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/users/%s/voice_mails", url.PathEscape(userID)), query, nil, page)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}
	return page, res, nil
}

// Get reads voicemail detail, including transcription when available. S2S apps
// require phone:read:voicemail:admin. Derive the ID from the intended owner's
// list before calling this account-authorized endpoint. No audio is fetched.
func (p *PhoneVoicemailsService) Get(ctx context.Context, voicemailID string) (*models.PhoneVoicemail, *http.Response, error) {
	if voicemailID == "" {
		return nil, nil, fmt.Errorf("voicemail ID is required")
	}
	voicemail := &models.PhoneVoicemail{}
	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/voice_mails/%s", url.PathEscape(voicemailID)), nil, nil, voicemail)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("expected status code %d but got %d", http.StatusOK, res.StatusCode)
	}
	return voicemail, res, nil
}
