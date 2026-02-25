package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/TheSlowpes/go-zoom/zoom/enums"
	"github.com/TheSlowpes/go-zoom/zoom/models"
)

type MeetingsSummaryServicers interface {
	GetSummary(ctx context.Context, opts ...MeetingsSummaryGetOptions) ([]*models.MeetingSummary, *http.Response, error)
	DeleteSummary(ctx context.Context, meetingId string) (*http.Response, error)
}

var _ MeetingsSummaryServicers = (*MeetingsService)(nil)

type MeetingsSummaryGetOptions func(*meetingsSummaryGetOptions)

type meetingsSummaryGetOptions struct {
	meetingId       string
	queryParameters *MeetingSummaryQueryParameters
}

type MeetingSummaryQueryParameters struct {
	From            time.Time             `url:"from,omitempty"`
	To              time.Time             `url:"to,omitempty"`
	TimeFilterField enums.TimeFilterField `url:"time_filter_field"`
}

func WithMeetingIdForSummary(meetingId string) MeetingsSummaryGetOptions {
	return func(o *meetingsSummaryGetOptions) {
		o.meetingId = meetingId
	}
}

func WithMeetingSummaryQueryParameters(params *MeetingSummaryQueryParameters) MeetingsSummaryGetOptions {
	return func(o *meetingsSummaryGetOptions) {
		o.queryParameters = params
	}
}

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
	}

	return meetingSummaries, res, nil
}

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
