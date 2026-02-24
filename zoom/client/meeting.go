package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/TheSlowpes/go-zoom/zoom/models"
)

type MeetingsServicers interface {
	Get(ctx context.Context, opts ...MeetingGetOptions) ([]*models.Meeting, *http.Response, error)
}

type MeetingsService struct {
	client *Client
}

var _ MeetingsServicers = (*MeetingsService)(nil)

type MeetingGetOptions func(*meetingsGetOptions)

type meetingsGetOptions struct {
	meetingId           string `url:"meetingId,omitempty"`
	userId              string `url:"userId,omitempty"`
	queryParameters     *MeetingQueryParameters
	listQueryParameters *MeetingListQueryParameters
}

type MeetingQueryParameters struct {
	OccurrenceId            string `url:"occurrence_id,omitempty"`
	ShowPreviousOccurrences bool   `url:"show_previous_occurrences,omitempty"`
}

type MeetingListQueryParameters struct {
	Type     string `url:"type,omitempty"`
	From     string `url:"from,omitempty"`
	To       string `url:"to,omitempty"`
	Timezone string `url:"timezone,omitempty"`
}

func WithMeetingId(meetingId string) MeetingGetOptions {
	return func(opts *meetingsGetOptions) {
		opts.meetingId = meetingId
	}
}

func WithMeetingUserId(userId string) MeetingGetOptions {
	return func(opts *meetingsGetOptions) {
		opts.userId = userId
	}
}

func WithMeetingQueryParameters(params *MeetingQueryParameters) MeetingGetOptions {
	return func(opts *meetingsGetOptions) {
		opts.queryParameters = params
	}
}

func WithMeetingListQueryParameters(params *MeetingListQueryParameters) MeetingGetOptions {
	return func(opts *meetingsGetOptions) {
		opts.listQueryParameters = params
	}
}
func (s *MeetingsService) Get(ctx context.Context, opts ...MeetingGetOptions) ([]*models.Meeting, *http.Response, error) {
	options := meetingsGetOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.meetingId != "" && options.userId != "" {
		return nil, nil, fmt.Errorf("Cannot specify both meetingId and userId")
	}

	if options.queryParameters != nil && options.listQueryParameters != nil {
		return nil, nil, fmt.Errorf("Cannot specify both queryParameters and listQueryParameters")
	}

	query := any(nil)
	if options.queryParameters != nil {
		query = options.queryParameters
	} else if options.listQueryParameters != nil {
		query = options.listQueryParameters
	}

	if options.meetingId != "" {
		endpoint := fmt.Sprintf("/meetings/%s", url.PathEscape(options.meetingId))
		meeting := &models.Meeting{}
		res, err := s.client.request(ctx, http.MethodGet, endpoint, query, nil, meeting)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		return []*models.Meeting{meeting}, res, nil
	}

	var endpoint string
	if options.userId != "" {
		endpoint = fmt.Sprintf("/users/%s/meetings", url.PathEscape(options.userId))
	} else {
		return nil, nil, fmt.Errorf("Must specify either meetingId or userId")
	}

	type response struct {
		*PaginationResponse
		Meetings []*models.Meeting `json:"meetings"`
	}

	queryResponse := &response{}
	var meetings []*models.Meeting

	res, err := s.client.request(ctx, http.MethodGet, endpoint, query, nil, queryResponse)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	meetings = append(meetings, queryResponse.Meetings...)

	type meetingListPageQuery struct {
		*MeetingListQueryParameters
		*PaginationOptions
	}

	for {
		if queryResponse.NextPageToken == "" {
			break
		}

		nextPageToken := queryResponse.NextPageToken
		pageQuery := &meetingListPageQuery{
			MeetingListQueryParameters: options.listQueryParameters,
			PaginationOptions: &PaginationOptions{
				NextPageToken: &nextPageToken,
			},
		}
		res, err = s.client.request(ctx, http.MethodGet, endpoint, pageQuery, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		meetings = append(meetings, queryResponse.Meetings...)
	}

	return meetings, res, nil
}
