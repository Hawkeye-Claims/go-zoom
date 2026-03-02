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

// MeetingsServicers is the interface implemented by MeetingsService. It
// declares the full set of meeting management operations available through the
// Zoom Meetings API.
type MeetingsServicers interface {
	// Get retrieves one or more meetings. Provide WithMeetingId to fetch a
	// single meeting by ID, or WithMeetingUserId to list all meetings for a
	// user (pagination is handled automatically).
	Get(ctx context.Context, opts ...MeetingGetOptions) ([]*models.Meeting, *http.Response, error)
	// Create schedules a new meeting for the specified user.
	Create(ctx context.Context, userId string, meetingAttributes MeetingAttributes) (*models.Meeting, *http.Response, error)
	// Delete cancels and removes the meeting identified by meetingId.
	Delete(ctx context.Context, meetingId int, opts ...MeetingDeleteOptions) (*http.Response, error)
	// Update patches an existing meeting identified by meetingId with the
	// supplied attributes.
	Update(ctx context.Context, meetingId int, meetingAttributes *MeetingUpdateAttributes, opts ...MeetingUpdateOptions) (*http.Response, error)
}

// MeetingsService implements MeetingsServicers and provides access to the Zoom
// Meetings API endpoints.
type MeetingsService struct {
	client *Client
}

// Compile-time assertion that MeetingsService satisfies the MeetingsServicers
// interface.
var _ MeetingsServicers = (*MeetingsService)(nil)

// MeetingGetOptions is a functional option for configuring a meeting Get
// request.
type MeetingGetOptions func(*meetingsGetOptions)

// MeetingUpdateOptions is a functional option for configuring a meeting Update
// (PATCH) request.
type MeetingUpdateOptions func(*meetingsUpdateOptions)

// MeetingDeleteOptions is a functional option for configuring a meeting Delete
// request.
type MeetingDeleteOptions func(*meetingsDeleteOptions)

// meetingsGetOptions holds the resolved configuration for a meeting Get call.
type meetingsGetOptions struct {
	meetingId           string `url:"meetingId,omitempty"`
	userId              string `url:"userId,omitempty"`
	queryParameters     *MeetingQueryParameters
	listQueryParameters *MeetingListQueryParameters
}

// meetingsUpdateOptions holds the resolved configuration for a meeting Update
// call.
type meetingsUpdateOptions struct {
	queryParameters *MeetingUpdateQueryParameters
}

// meetingsDeleteOptions holds the resolved configuration for a meeting Delete
// call.
type meetingsDeleteOptions struct {
	queryParameters *MeetingDeleteQueryParameters
}

// MeetingQueryParameters holds optional query parameters for a single-meeting
// Get request.
type MeetingQueryParameters struct {
	// OccurrenceId filters the response to a specific occurrence of a
	// recurring meeting.
	OccurrenceId string `url:"occurrence_id,omitempty"`
	// ShowPreviousOccurrences, when true, includes past occurrences in the
	// response.
	ShowPreviousOccurrences bool `url:"show_previous_occurrences,omitempty"`
}

// MeetingListQueryParameters holds optional query parameters for listing a
// user's meetings.
type MeetingListQueryParameters struct {
	// Type filters the list to a specific meeting type (e.g. "scheduled").
	Type string `url:"type,omitempty"`
	// From is the start of the date range for the query (RFC 3339 date string).
	From string `url:"from,omitempty"`
	// To is the end of the date range for the query (RFC 3339 date string).
	To string `url:"to,omitempty"`
	// Timezone specifies the time zone used to interpret From and To.
	Timezone string `url:"timezone,omitempty"`
}

// MeetingUpdateQueryParameters holds optional query parameters for a meeting
// Update request.
type MeetingUpdateQueryParameters struct {
	// OccurrenceId targets a specific occurrence of a recurring meeting for
	// the update.
	OccurrenceId string `url:"occurrence_id,omitempty"`
}

// MeetingDeleteQueryParameters holds optional query parameters for a meeting
// Delete request.
type MeetingDeleteQueryParameters struct {
	// OccurrenceId targets a specific occurrence of a recurring meeting for
	// deletion.
	OccurrenceId string `url:"occurrence_id,omitempty"`
	// ScheduleForReminder, when true, sends a reminder to the host before
	// deletion.
	ScheduleForReminder bool `url:"schedule_for_reminder,omitempty"`
	// CancelMeetingReminder, when true, sends cancellation notifications to
	// meeting participants.
	CancelMeetingReminder bool `url:"cancel_meeting_reminder,omitempty"`
}

// WithMeetingId returns a MeetingGetOptions that fetches a single meeting by
// its Zoom meeting ID.
func WithMeetingId(meetingId string) MeetingGetOptions {
	return func(opts *meetingsGetOptions) {
		opts.meetingId = meetingId
	}
}

// WithMeetingUserId returns a MeetingGetOptions that lists all meetings
// belonging to the user identified by userId.
func WithMeetingUserId(userId string) MeetingGetOptions {
	return func(opts *meetingsGetOptions) {
		opts.userId = userId
	}
}

// WithMeetingQueryParameters returns a MeetingGetOptions that attaches the
// given query parameters to a single-meeting Get request.
func WithMeetingQueryParameters(params *MeetingQueryParameters) MeetingGetOptions {
	return func(opts *meetingsGetOptions) {
		opts.queryParameters = params
	}
}

// WithMeetingListQueryParameters returns a MeetingGetOptions that attaches the
// given query parameters to a list-meetings request.
func WithMeetingListQueryParameters(params *MeetingListQueryParameters) MeetingGetOptions {
	return func(opts *meetingsGetOptions) {
		opts.listQueryParameters = params
	}
}

// WithMeetingUpdateQueryParameters returns a MeetingUpdateOptions that attaches
// the given query parameters to a meeting Update request.
func WithMeetingUpdateQueryParameters(params *MeetingUpdateQueryParameters) MeetingUpdateOptions {
	return func(opts *meetingsUpdateOptions) {
		opts.queryParameters = params
	}
}

// WithMeetingDeleteQueryParameters returns a MeetingDeleteOptions that attaches
// the given query parameters to a meeting Delete request.
func WithMeetingDeleteQueryParameters(params *MeetingDeleteQueryParameters) MeetingDeleteOptions {
	return func(o *meetingsDeleteOptions) {
		o.queryParameters = params
	}
}

// Get retrieves Zoom meetings. Provide WithMeetingId to fetch a single meeting
// by its ID, or WithMeetingUserId to list all meetings for a user (pagination
// is followed automatically). It is an error to specify both a meeting ID and
// a user ID, or to combine MeetingQueryParameters with
// MeetingListQueryParameters.
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

// MeetingAttributes holds the fields for creating a new Zoom meeting.
type MeetingAttributes struct {
	// Agenda is an optional description or agenda for the meeting.
	Agenda string `json:"agenda,omitempty"`
	// DefaultPassword, when true, generates a random password for the meeting.
	DefaultPassword bool `json:"default_password,omitempty"`
	// Duration is the scheduled meeting duration in minutes.
	Duration int `json:"duration,omitempty"`
	// Password is an optional passcode for joining the meeting.
	Password string `json:"password,omitempty"`
	// PreSchedule, when true, creates a pre-scheduled meeting.
	PreSchedule bool `json:"pre_schedule,omitempty"`
	// Recurrence holds the recurrence settings for a recurring meeting.
	Recurrence models.MeetingRecurrence `json:"recurrence"`
	// ScheduleFor is the email or user ID of the user the meeting is scheduled
	// on behalf of.
	ScheduleFor string `json:"schedule_for,omitempty"`
	// Settings holds the meeting configuration options.
	Settings models.MeetingSettings `json:"settings"`
	// StartTime is the scheduled start time of the meeting.
	StartTime time.Time `json:"start_time"`
	// TemplateID is the ID of a meeting template to apply.
	TemplateID string `json:"template_id,omitempty"`
	// Timezone specifies the time zone for the meeting start time.
	Timezone string `json:"timezone,omitempty"`
	// Topic is the meeting title.
	Topic string `json:"topic,omitempty"`
	// TrackingFields holds custom tracking field values for reporting.
	TrackingFields []models.MeetingTrackingField `json:"tracking_fields"`
	// Type specifies the meeting type (e.g. scheduled, recurring).
	Type enums.MeetingType `json:"type,omitempty"`
}

// Create schedules a new Zoom meeting for the user identified by userId using
// the provided meeting attributes.
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

// MeetingUpdateAttributes holds the fields that can be patched on an existing
// Zoom meeting. All fields are optional; only non-zero values are sent.
type MeetingUpdateAttributes struct {
	// Agenda is the updated meeting description or agenda.
	Agenda string `json:"agenda,omitempty"`
	// Duration is the updated meeting duration in minutes.
	Duration int `json:"duration,omitempty"`
	// Password is the updated meeting passcode.
	Password string `json:"password,omitempty"`
	// PreSchedule, when true, marks the meeting as pre-scheduled.
	PreSchedule bool `json:"pre_schedule,omitempty"`
	// Recurrence holds updated recurrence settings.
	Recurrence models.MeetingRecurrence `json:"recurrence"`
	// ScheduleFor is the email or user ID of the host the meeting is scheduled
	// for.
	ScheduleFor string `json:"schedule_for"`
	// Settings holds the updated meeting configuration options.
	Settings models.MeetingSettings `json:"settings"`
	// StartTime is the updated meeting start time.
	StartTime time.Time `json:"start_time"`
	// TemplateID is the updated meeting template ID.
	TemplateID string `json:"template_id"`
	// Timezone specifies the time zone for the updated start time.
	Timezone string `json:"timezone,omitempty"`
	// Topic is the updated meeting title.
	Topic string `json:"topic,omitempty"`
	// TrackingFields holds updated custom tracking field values.
	TrackingFields []models.MeetingTrackingField `json:"tracking_fields"`
	// Type is the updated meeting type.
	Type enums.MeetingType `json:"type"`
}

// Update patches an existing Zoom meeting identified by meetingId with the
// supplied attributes. Optional query parameters can be supplied via
// MeetingUpdateOptions.
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

// Delete cancels and removes the Zoom meeting identified by meetingId.
// Optional query parameters (e.g. cancellation notifications) can be supplied
// via MeetingDeleteOptions.
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
