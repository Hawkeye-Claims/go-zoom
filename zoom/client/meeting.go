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

type MeetingsServicers interface {
	Get(ctx context.Context, opts ...MeetingGetOptions) ([]*models.Meeting, *http.Response, error)
	Create(ctx context.Context, userId string, meetingAttributes MeetingAttributes) (*models.Meeting, *http.Response, error)
	Delete(ctx context.Context, meetingId int, opts ...MeetingDeleteOptions) (*http.Response, error)
	Update(ctx context.Context, meetingId int, meetingAttributes *MeetingUpdateAttributes, opts ...MeetingUpdateOptions) (*http.Response, error)
}

type MeetingsService struct {
	client *Client
}

var _ MeetingsServicers = (*MeetingsService)(nil)

type MeetingGetOptions func(*meetingsGetOptions)

type MeetingUpdateOptions func(*meetingsUpdateOptions)

type MeetingDeleteOptions func(*meetingsDeleteOptions)

type meetingsGetOptions struct {
	meetingId           string `url:"meetingId,omitempty"`
	userId              string `url:"userId,omitempty"`
	queryParameters     *MeetingQueryParameters
	listQueryParameters *MeetingListQueryParameters
}

type meetingsUpdateOptions struct {
	queryParameters *MeetingUpdateQueryParameters
}

type meetingsDeleteOptions struct {
	queryParameters *MeetingDeleteQueryParameters
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

type MeetingUpdateQueryParameters struct {
	OccurrenceId string `url:"occurrence_id,omitempty"`
}

type MeetingDeleteQueryParameters struct {
	OccurrenceId          string `url:"occurrence_id,omitempty"`
	ScheduleForReminder   bool   `url:"schedule_for_reminder,omitempty"`
	CancelMeetingReminder bool   `url:"cancel_meeting_reminder,omitempty"`
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

func WithMeetingUpdateQueryParameters(params *MeetingUpdateQueryParameters) MeetingUpdateOptions {
	return func(opts *meetingsUpdateOptions) {
		opts.queryParameters = params
	}
}

func WithMeetingDeleteQueryParameters(params *MeetingDeleteQueryParameters) MeetingDeleteOptions {
	return func(o *meetingsDeleteOptions) {
		o.queryParameters = params
	}
}
func (m *MeetingsService) Get(ctx context.Context, opts ...MeetingGetOptions) ([]*models.Meeting, *http.Response, error) {
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
		res, err := m.client.request(ctx, http.MethodGet, endpoint, query, nil, meeting)
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

	res, err := m.client.request(ctx, http.MethodGet, endpoint, query, nil, queryResponse)
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
		res, err = m.client.request(ctx, http.MethodGet, endpoint, pageQuery, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		meetings = append(meetings, queryResponse.Meetings...)
	}

	return meetings, res, nil
}

type MeetingAttributes struct {
	Agenda          string                        `json:"agenda,omitempty"`
	DefaultPassword bool                          `json:"default_password,omitempty"`
	Duration        int                           `json:"duration,omitempty"`
	Password        string                        `json:"password,omitempty"`
	PreSchedule     bool                          `json:"pre_schedule,omitempty"`
	Recurrence      models.MeetingRecurrence      `json:"recurrence"`
	ScheduleFor     string                        `json:"schedule_for,omitempty"`
	Settings        models.MeetingSettings        `json:"settings"`
	StartTime       time.Time                     `json:"start_time"`
	TemplateID      string                        `json:"template_id,omitempty"`
	Timezone        string                        `json:"timezone,omitempty"`
	Topic           string                        `json:"topic,omitempty"`
	TrackingFields  []models.MeetingTrackingField `json:"tracking_fields"`
	Type            enums.MeetingType             `json:"type,omitempty"`
}

func (m *MeetingsService) Create(ctx context.Context, userId string, meetingAttributes MeetingAttributes) (*models.Meeting, *http.Response, error) {
	var response models.Meeting

	res, err := m.client.request(ctx, http.MethodPost, fmt.Sprintf("/users/%s/meetings", url.PathEscape(userId)), nil, meetingAttributes, &response)
	if err != nil {
		return &models.Meeting{}, res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusCreated {
		return &models.Meeting{}, res, fmt.Errorf("Expected status code %d, got %d", http.StatusCreated, res.StatusCode)
	}
	return &response, res, nil
}

type MeetingUpdateAttributes struct {
	Agenda         string                        `json:"agenda,omitempty"`
	Duration       int                           `json:"duration,omitempty"`
	Password       string                        `json:"password,omitempty"`
	PreSchedule    bool                          `json:"pre_schedule,omitempty"`
	Recurrence     models.MeetingRecurrence      `json:"recurrence"`
	ScheduleFor    string                        `json:"schedule_for"`
	Settings       models.MeetingSettings        `json:"settings"`
	StartTime      time.Time                     `json:"start_time"`
	TemplateID     string                        `json:"template_id"`
	Timezone       string                        `json:"timezone,omitempty"`
	Topic          string                        `json:"topic,omitempty"`
	TrackingFields []models.MeetingTrackingField `json:"tracking_fields"`
	Type           enums.MeetingType             `json:"type"`
}

func (m *MeetingsService) Update(ctx context.Context, meetingId int, meetingAttributes *MeetingUpdateAttributes, opts ...MeetingUpdateOptions) (*http.Response, error) {
	options := meetingsUpdateOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	res, err := m.client.request(ctx, http.MethodPatch, fmt.Sprintf("/meetings/%d", meetingId), options.queryParameters, meetingAttributes, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}

func (m *MeetingsService) Delete(ctx context.Context, meetingId int, opts ...MeetingDeleteOptions) (*http.Response, error) {
	options := meetingsDeleteOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	res, err := m.client.request(ctx, http.MethodDelete, fmt.Sprintf("/meetings/%d", meetingId), options.queryParameters, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}
