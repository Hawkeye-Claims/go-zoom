package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Hawkeye-Claims/go-zoom/zoom/enums"
	"github.com/Hawkeye-Claims/go-zoom/zoom/models"
)

// MeetingsSummaryServicers is the interface implemented by MeetingsService for
// AI meeting-summary operations. It declares retrieval and deletion of
// meeting summaries via the Zoom Meetings Summaries API.
type MeetingsSummaryServicers interface {
	// GetSummary retrieves one or more AI meeting summaries. When
	// WithMeetingIdForSummary is supplied it returns the summary for that
	// specific meeting; otherwise it lists all summaries, following pagination
	// automatically.
	GetSummary(ctx context.Context, opts ...MeetingsSummaryGetOptions) ([]*models.MeetingSummary, *http.Response, error)
	// DeleteSummary removes the AI meeting summary for the meeting identified
	// by meetingId.
	DeleteSummary(ctx context.Context, meetingId string) (*http.Response, error)
}

// Compile-time assertion that MeetingsService satisfies the
// MeetingsSummaryServicers interface.
var _ MeetingsSummaryServicers = (*MeetingsService)(nil)

// MeetingsSummaryGetOptions is a functional option for configuring a meeting
// summary Get request.
type MeetingsSummaryGetOptions func(*meetingsSummaryGetOptions)

// meetingsSummaryGetOptions holds the resolved configuration for a summary Get
// call.
type meetingsSummaryGetOptions struct {
	meetingId       string
	queryParameters *MeetingSummaryQueryParameters
}

// MeetingSummaryQueryParameters holds optional query parameters for listing
// meeting summaries.
type MeetingSummaryQueryParameters struct {
	// From is the start of the date range to filter summaries.
	From time.Time `url:"from,omitempty"`
	// To is the end of the date range to filter summaries.
	To time.Time `url:"to,omitempty"`
	// TimeFilterField specifies which timestamp field (start time or end time)
	// From and To are applied to.
	TimeFilterField enums.TimeFilterField `url:"time_filter_field"`
}

// WithMeetingIdForSummary returns a MeetingsSummaryGetOptions that fetches the
// AI summary for the meeting identified by meetingId.
func WithMeetingIdForSummary(meetingId string) MeetingsSummaryGetOptions {
	return func(o *meetingsSummaryGetOptions) {
		o.meetingId = meetingId
	}
}

// WithMeetingSummaryQueryParameters returns a MeetingsSummaryGetOptions that
// attaches the given query parameters to a list-summaries request.
func WithMeetingSummaryQueryParameters(params *MeetingSummaryQueryParameters) MeetingsSummaryGetOptions {
	return func(o *meetingsSummaryGetOptions) {
		o.queryParameters = params
	}
}

// GetSummary retrieves AI meeting summaries. When WithMeetingIdForSummary is
// provided it returns a slice containing only that meeting's summary.
// Otherwise it lists all summaries, following pagination automatically so the
// returned slice contains every summary.
func (m *MeetingsService) GetSummary(ctx context.Context, opts ...MeetingsSummaryGetOptions) ([]*models.MeetingSummary, *http.Response, error) {
	options := meetingsSummaryGetOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.meetingId != "" {
		endpoint := fmt.Sprintf("/meetings/%s/meeting_summary", url.PathEscape(options.meetingId))
		meetingSummary := &models.MeetingSummary{}
		res, err := m.client.request(ctx, http.MethodGet, endpoint, options.queryParameters, nil, meetingSummary)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		return []*models.MeetingSummary{meetingSummary}, res, nil
	}

	type response struct {
		*PaginationResponse
		Summaries []*models.MeetingSummary `json:"summaries"`
	}

	queryResponse := &response{}
	var meetingSummaries []*models.MeetingSummary

	endpoint := "/meetings/meeting_summaries"

	res, err := m.client.request(ctx, http.MethodGet, endpoint, options.queryParameters, nil, queryResponse)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	meetingSummaries = append(meetingSummaries, queryResponse.Summaries...)

	type meetingSummaryPageQuery struct {
		*MeetingSummaryQueryParameters
		*PaginationOptions
	}
	for {
		if queryResponse.NextPageToken == "" {
			break
		}
		nextPageToken := queryResponse.NextPageToken
		pageQuery := &meetingSummaryPageQuery{
			MeetingSummaryQueryParameters: options.queryParameters,
			PaginationOptions: &PaginationOptions{
				NextPageToken: &nextPageToken,
			},
		}
		res, err = m.client.request(ctx, http.MethodGet, endpoint, pageQuery, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		meetingSummaries = append(meetingSummaries, queryResponse.Summaries...)
	}

	return meetingSummaries, res, nil
}

// DeleteSummary removes the AI meeting summary for the meeting identified by
// meetingId.
func (m *MeetingsService) DeleteSummary(ctx context.Context, meetingId string) (*http.Response, error) {
	res, err := m.client.request(ctx, http.MethodDelete, fmt.Sprintf("/meetings/%s/meeting_summary", url.PathEscape(meetingId)), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}
