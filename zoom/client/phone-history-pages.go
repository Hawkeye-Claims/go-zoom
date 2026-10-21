package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Hawkeye-Claims/go-zoom/zoom/models"
)

// PhoneUserCallHistoryPager provides read-only paging without changing the
// existing PhoneCallHistoryServicer contract or its downstream implementations.
type PhoneUserCallHistoryPager interface {
	ListUserPage(ctx context.Context, userID string, params *PhoneCallHistoryQueryParameters, pagination *PaginationOptions) (*PhoneUserCallHistoryPage, *http.Response, error)
}

var _ PhoneUserCallHistoryPager = (*PhoneCallHistoryService)(nil)

// PhoneUserCallHistoryPage preserves both arrays documented by Zoom's owner
// history endpoint. They are alternate representations and should not be added
// together. Prefer CallElements when non-nil, otherwise use CallLogs.
type PhoneUserCallHistoryPage struct {
	PaginationResponse
	// From and To are the bounds echoed by Zoom.
	From string `json:"from"`
	To   string `json:"to"`
	// CallElements contains modern flat call elements, if returned.
	CallElements []*models.CallHistory `json:"call_elements"`
	// CallLogs contains the older documented call log representation.
	CallLogs []*models.CallHistory `json:"call_logs"`
}

// ListUserPage reads one call history page for a user ID or email. S2S apps
// require phone:read:list_call_logs:admin. Reuse the filters and pass the returned
// NextPageToken to read subsequent pages. Tokens expire after 15 minutes. This
// method does not retry or follow cursors; errors are not end-of-data. The HTTP
// response body is closed before returning.
func (p *PhoneCallHistoryService) ListUserPage(ctx context.Context, userID string, params *PhoneCallHistoryQueryParameters, pagination *PaginationOptions) (*PhoneUserCallHistoryPage, *http.Response, error) {
	if userID == "" {
		return nil, nil, fmt.Errorf("user ID is required")
	}
	query := &struct {
		*PhoneCallHistoryQueryParameters
		*PaginationOptions
	}{params, pagination}
	page := &PhoneUserCallHistoryPage{}
	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/users/%s/call_history", url.PathEscape(userID)), query, nil, page)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}
	return page, res, nil
}
