package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/TheSlowpes/go-zoom/zoom/models"
)

// PhoneRecordingsServicer is the interface implemented by
// PhoneRecordingsService. It declares the operations available for managing
// Zoom Phone call recordings via the Zoom Phone API.
type PhoneRecordingsServicer interface {
	// Get retrieves call recordings. Provide WithRecordingCallId to fetch
	// recordings for a specific call, WithRecordingUserId to list recordings
	// for a user, or no option to list all account-level recordings.
	Get(ctx context.Context, opts ...CallRecordingGetOptions) ([]*models.CallRecording, *http.Response, error)
	// Delete permanently removes the recording identified by recordingId.
	Delete(ctx context.Context, recordingId string) (*http.Response, error)
	// DownloadCallRecording streams the audio file identified by fileId into
	// the supplied io.Writer.
	DownloadCallRecording(ctx context.Context, fileId string, w io.Writer) (*http.Response, error)
	// DownloadCallTranscript retrieves the transcript for the recording
	// identified by recordingId.
	DownloadCallTranscript(ctx context.Context, recordingId string) (*models.RecordingTranscript, *http.Response, error)
	// EnableAutoDelete enables automatic deletion for the recording identified
	// by recordingId.
	EnableAutoDelete(ctx context.Context, recordingId string) (*http.Response, error)
	// DisableAutoDelete disables automatic deletion for the recording
	// identified by recordingId.
	DisableAutoDelete(ctx context.Context, recordingId string) (*http.Response, error)
	// Recover restores the recording identified by recordingId from the trash.
	Recover(ctx context.Context, recordingId string) (*http.Response, error)
}

// PhoneRecordingsService implements PhoneRecordingsServicer and provides
// access to Zoom Phone recording API endpoints.
type PhoneRecordingsService struct {
	client *Client
}

// Compile-time assertion that PhoneRecordingsService satisfies the
// PhoneRecordingsServicer interface.
var _ PhoneRecordingsServicer = (*PhoneRecordingsService)(nil)

// CallRecordingQueryParameters holds optional query parameters for filtering
// call recording results.
type CallRecordingQueryParameters struct {
	// From is the start of the date range (YYYY-MM-DD) for the query.
	From string `url:"from,omitempty"`
	// To is the end of the date range (YYYY-MM-DD) for the query.
	To string `url:"to,omitempty"`
	// OwnerType filters results by the type of recording owner
	// (e.g. "user" or "call_queue").
	OwnerType string `url:"owner_type,omitempty"`
	// RecordingType filters results by recording type
	// (e.g. "OnDemand" or "Automatic").
	RecordingType string `url:"recording_type,omitempty"`
	// SiteID filters results to recordings associated with a specific site.
	SiteID string `url:"site_id,omitempty"`
	// QueryDateType specifies whether the date range applies to the recording
	// start or end time.
	QueryDateType string `url:"query_date_type,omitempty"`
	// GroupID filters results to recordings associated with a specific group.
	GroupID string `url:"group_id,omitempty"`
}

// CallRecordingGetOptions is a functional option for configuring a call
// recording Get request.
type CallRecordingGetOptions func(*callRecordingGetOptions)

// callRecordingGetOptions holds the resolved configuration for a recording Get
// call.
type callRecordingGetOptions struct {
	callId          string
	userId          string
	queryParameters *CallRecordingQueryParameters
}

// WithRecordingUserId returns a CallRecordingGetOptions that lists call
// recordings for the user identified by userId.
func WithRecordingUserId(userId string) CallRecordingGetOptions {
	return func(opts *callRecordingGetOptions) {
		opts.userId = userId
	}
}

// WithRecordingCallId returns a CallRecordingGetOptions that fetches the
// recording(s) associated with the call identified by callId.
func WithRecordingCallId(callId string) CallRecordingGetOptions {
	return func(opts *callRecordingGetOptions) {
		opts.callId = callId
	}
}

// WithCallRecordingQueryParameters returns a CallRecordingGetOptions that
// attaches the given query parameters to the recording list request.
func WithCallRecordingQueryParameters(queryParameters *CallRecordingQueryParameters) CallRecordingGetOptions {
	return func(opts *callRecordingGetOptions) {
		opts.queryParameters = queryParameters
	}
}

// Get retrieves Zoom Phone call recordings. At most one option may be
// supplied: WithRecordingCallId to fetch recordings for a specific call,
// WithRecordingUserId to list recordings for a user (paginated
// automatically), or no option to list all account recordings (paginated
// automatically). Supplying more than one option returns an error.
func (r *PhoneRecordingsService) Get(ctx context.Context, opts ...CallRecordingGetOptions) ([]*models.CallRecording, *http.Response, error) {
	options := &callRecordingGetOptions{}
	counter := 0
	for _, opt := range opts {
		opt(options)
		counter++
	}
	if counter > 1 {
		return nil, nil, errors.New("Only one of WithRecordingUserId, WithRecordingCallId, or WithCallRecordingQueryParameters can be used")
	}

	if options.callId != "" {
		endpoint := fmt.Sprintf("/phone/call_logs/%s/recordings", url.PathEscape(options.callId))
		callRecording := &models.CallRecording{}
		res, err := r.client.request(ctx, http.MethodGet, endpoint, nil, nil, callRecording)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		return []*models.CallRecording{callRecording}, res, nil
	}

	type response struct {
		*PaginationResponse
		CallRecordings []*models.CallRecording `json:"call_recordings"`
	}

	queryResponse := &response{}
	var callRecordings []*models.CallRecording

	if options.userId != "" {
		endpoint := fmt.Sprintf("/phone/users/%s/recordings", url.PathEscape(options.userId))

		res, err := r.client.request(ctx, http.MethodGet, endpoint, options.queryParameters, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}

		callRecordings = append(callRecordings, queryResponse.CallRecordings...)

		type callRecordingPageQuery struct {
			*PaginationOptions
		}
		for {
			if queryResponse.NextPageToken == "" {
				break
			}
			nextPageToken := queryResponse.NextPageToken
			pageQuery := &callRecordingPageQuery{
				PaginationOptions: &PaginationOptions{
					NextPageToken: &nextPageToken,
				},
			}
			res, err := r.client.request(ctx, http.MethodGet, endpoint, pageQuery, nil, queryResponse)
			if err != nil {
				return nil, res, fmt.Errorf("Error making request: %w", err)
			}
			callRecordings = append(callRecordings, queryResponse.CallRecordings...)
		}

		return callRecordings, res, nil
	}

	endpoint := "/phone/recordings"

	res, err := r.client.request(ctx, http.MethodGet, endpoint, options.queryParameters, nil, queryResponse)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	callRecordings = append(callRecordings, queryResponse.CallRecordings...)

	type callRecordingPageQuery struct {
		*CallRecordingQueryParameters
		*PaginationOptions
	}
	for {
		if queryResponse.NextPageToken == "" {
			break
		}
		nextPageToken := queryResponse.NextPageToken
		pageQuery := &callRecordingPageQuery{
			CallRecordingQueryParameters: options.queryParameters,
			PaginationOptions: &PaginationOptions{
				NextPageToken: &nextPageToken,
			},
		}
		res, err := r.client.request(ctx, http.MethodGet, endpoint, pageQuery, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		callRecordings = append(callRecordings, queryResponse.CallRecordings...)
	}

	return callRecordings, res, nil
}

// DownloadCallRecording streams the audio file for the recording identified
// by fileId into w. The raw HTTP response is returned alongside any error.
func (r *PhoneRecordingsService) DownloadCallRecording(ctx context.Context, fileId string, w io.Writer) (*http.Response, error) {
	res, err := r.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/recordings/download/%s", url.PathEscape(fileId)), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
	_, err = io.Copy(w, res.Body)
	return res, err
}

// DownloadCallTranscript retrieves the transcript for the recording identified
// by recordingId.
func (r *PhoneRecordingsService) DownloadCallTranscript(ctx context.Context, recordingId string) (*models.RecordingTranscript, *http.Response, error) {
	var transcript models.RecordingTranscript
	res, err := r.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/recording_transcript/download/%s", url.PathEscape(recordingId)), nil, nil, &transcript)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
	return &transcript, res, nil
}

// Delete permanently removes the call recording identified by recordingId.
func (r *PhoneRecordingsService) Delete(ctx context.Context, recordingId string) (*http.Response, error) {
	res, err := r.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/recordings/%s", url.PathEscape(recordingId)), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}

// EnableAutoDelete sets the auto-delete flag to true on the recording
// identified by recordingId, scheduling it for automatic removal.
func (r *PhoneRecordingsService) EnableAutoDelete(ctx context.Context, recordingId string) (*http.Response, error) {
	type body struct {
		AutoDeleteEnable bool `json:"auto_delete_enable"`
	}
	requestBody := &body{AutoDeleteEnable: true}
	res, err := r.client.request(ctx, http.MethodPatch, fmt.Sprintf("/phone/recordings/%s", url.PathEscape(recordingId)), nil, requestBody, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}

// DisableAutoDelete sets the auto-delete flag to false on the recording
// identified by recordingId, preventing it from being automatically removed.
func (r *PhoneRecordingsService) DisableAutoDelete(ctx context.Context, recordingId string) (*http.Response, error) {
	type body struct {
		AutoDeleteEnable bool `json:"auto_delete_enable"`
	}
	requestBody := &body{AutoDeleteEnable: false}
	res, err := r.client.request(ctx, http.MethodPatch, fmt.Sprintf("/phone/recordings/%s", url.PathEscape(recordingId)), nil, requestBody, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}

// Recover restores the recording identified by recordingId from the trash.
func (r *PhoneRecordingsService) Recover(ctx context.Context, recordingId string) (*http.Response, error) {
	type body struct {
		Action string `json:"action"`
	}
	requestBody := &body{Action: "recover"}
	res, err := r.client.request(ctx, http.MethodPost, fmt.Sprintf("/phone/recordings/%s/trash", url.PathEscape(recordingId)), nil, requestBody, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}
