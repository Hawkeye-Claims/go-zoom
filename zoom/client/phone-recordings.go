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

type PhoneRecordingsServicer interface {
	Get(ctx context.Context, opts ...CallRecorddingGetOptions) ([]*models.CallRecording, *http.Response, error)
	DownloadCallRecording(ctx context.Context, fileId string, w io.Writer) (*http.Response, error)
	DownloadCallTranscript(ctx context.Context, recordingId string) (*models.RecordingTranscript, *http.Response, error)
}

type PhoneRecordingsService struct {
	client *Client
}

var _ PhoneRecordingsServicer = (*PhoneRecordingsService)(nil)

type CallRecordingQueryParameters struct {
	From          string `url:"from,omitempty"`
	To            string `url:"to,omitempty"`
	OwnerType     string `url:"owner_type,omitempty"`
	RecordingType string `url:"recording_type,omitempty"`
	SiteID        string `url:"site_id,omitempty"`
	QueryDateType string `url:"query_date_type,omitempty"`
	GroupID       string `url:"group_id,omitempty"`
}

type CallRecorddingGetOptions func(*callRecordingGetOptions)

type callRecordingGetOptions struct {
	callId          string
	userId          string
	queryParameters *CallRecordingQueryParameters
}

func WithRecordingUserId(userId string) CallRecorddingGetOptions {
	return func(opts *callRecordingGetOptions) {
		opts.userId = userId
	}
}

func WithRecordingCallId(callId string) CallRecorddingGetOptions {
	return func(opts *callRecordingGetOptions) {
		opts.callId = callId
	}
}

func WithCallRecordingQueryParameters(queryParameters *CallRecordingQueryParameters) CallRecorddingGetOptions {
	return func(opts *callRecordingGetOptions) {
		opts.queryParameters = queryParameters
	}
}

func (r *PhoneRecordingsService) Get(ctx context.Context, opts ...CallRecorddingGetOptions) ([]*models.CallRecording, *http.Response, error) {
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
