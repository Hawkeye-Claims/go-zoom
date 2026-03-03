package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Hawkeye-Claims/go-zoom/zoom/enums"
	"github.com/Hawkeye-Claims/go-zoom/zoom/models"
)

// PhoneService groups all Zoom Phone sub-services under a single struct
// accessible via Client.Phone. Each field corresponds to a distinct Zoom Phone
// API domain.
type PhoneService struct {
	client *Client
	// CallHistory provides access to call-log and call-history operations.
	CallHistory *PhoneCallHistoryService
	// Settings provides access to account-level phone settings operations.
	Settings *PhoneSettingsService
	// Recordings provides access to call-recording operations.
	Recordings *PhoneRecordingsService
	// Users provides access to Zoom Phone user profile operations.
	Users *PhoneUsersService
}

// NewPhoneService initialises all Zoom Phone sub-services and attaches them to
// client.Phone. It must be called before accessing any Phone service method.
func NewPhoneService(client *Client) {
	client.Phone = &PhoneService{client: client}
	client.Phone.CallHistory = &PhoneCallHistoryService{client: client}
	client.Phone.Settings = &PhoneSettingsService{client: client}
	client.Phone.Recordings = &PhoneRecordingsService{client: client}
	client.Phone.Users = &PhoneUsersService{client: client}
}

// PhoneCallHistoryServicer is the interface implemented by
// PhoneCallHistoryService. It declares the operations available for
// retrieving and managing call history records via the Zoom Phone API.
type PhoneCallHistoryServicer interface {
	// Get retrieves call history records. Provide WithPhoneCallHistoryUUID to
	// fetch a specific record, WithUserIdForPhoneCallHistory to fetch records
	// for a user, or no option to list account-level call history.
	Get(ctx context.Context, opts ...PhoneCallHistoryGetOptions) ([]*models.CallHistory, *http.Response, error)
	// AddClientCode associates a client code with a specific call log entry.
	AddClientCode(ctx context.Context, callLogId, clientCode string) (*http.Response, error)
	// DeleteUserCallHistory removes a specific call log entry from a user's
	// history.
	DeleteUserCallHistory(ctx context.Context, userId, callLogId string) (*http.Response, error)
	// GetCallElement retrieves a specific call element by its ID.
	GetCallElement(ctx context.Context, callElementId string) (*models.CallElement, *http.Response, error)
	// GetAICallSummary retrieves an AI-generated call summary for the given
	// user and summary ID.
	GetAICallSummary(ctx context.Context, userId, aiCallSummaryId string) (*models.AICallSummary, *http.Response, error)
}

// PhoneCallHistoryService implements PhoneCallHistoryServicer and provides
// access to Zoom Phone call history API endpoints.
type PhoneCallHistoryService struct {
	client *Client
}

// Compile-time assertion that PhoneCallHistoryService satisfies the
// PhoneCallHistoryServicer interface.
var _ PhoneCallHistoryServicer = (*PhoneCallHistoryService)(nil)

// PhoneCallHistoryGetOptions is a functional option for configuring a call
// history Get request.
type PhoneCallHistoryGetOptions func(*phoneCallHistoryGetOptions)

// phoneCallHistoryGetOptions holds the resolved configuration for a call
// history Get call.
type phoneCallHistoryGetOptions struct {
	phoneCallHistoryUUID            string
	userId                          string
	phoneCallHistoryQueryParameters *PhoneCallHistoryQueryParameters
}

// PhoneCallHistoryQueryParameters holds optional query parameters for
// filtering call history results.
type PhoneCallHistoryQueryParameters struct {
	// From is the start of the date range (YYYY-MM-DD) for the query.
	From string `url:"from,omitempty"`
	// To is the end of the date range (YYYY-MM-DD) for the query.
	To string `url:"to,omitempty"`
	// Keyword filters results by a search keyword (e.g. phone number or name).
	Keyword string `url:"keyword,omitempty"`
	// Directions filters results by call direction (inbound/outbound).
	Directions []enums.Direction `url:"directions,omitempty"`
	// ConnectTypes filters results by how the call was connected.
	ConnectTypes []enums.ConnectType `url:"connect_types,omitempty"`
	// NumberTypes filters results by the type of phone number involved.
	NumberTypes []enums.NumberType `url:"number_types,omitempty"`
	// CallTypes filters results by the call type.
	CallTypes []enums.CallType `url:"call_types,omitempty"`
	// ExtensionTypes filters results by the type of extension involved.
	ExtensionTypes []enums.ExtensionType `url:"extension_types,omitempty"`
	// CallResult filters results by the call outcome.
	CallResult []enums.CallResult `url:"call_result,omitempty"`
	// GroupIDs filters results to calls associated with the given group IDs.
	GroupIDs []string `url:"group_ids,omitempty"`
	// SiteIDs filters results to calls associated with the given site IDs.
	SiteIDs []string `url:"site_ids,omitempty"`
	// Department filters results by the user's department.
	Department string `url:"department,omitempty"`
	// CostCenter filters results by the user's cost center.
	CostCenter string `url:"cost_center,omitempty"`
	// TimeType specifies whether From/To apply to the call start or end time.
	TimeType enums.TimeType `url:"time_type,omitempty"`
	// RecordingStatus filters results by the recording status of the call.
	RecordingStatus enums.RecordingStatus `url:"recording_status,omitempty"`
	// WithVoicemail, when true, includes calls that have a voicemail.
	WithVoicemail bool `url:"with_voicemail,omitempty"`
}

// WithPhoneCallHistoryUUID returns a PhoneCallHistoryGetOptions that fetches a
// single call history record by its UUID.
func WithPhoneCallHistoryUUID(uuid string) PhoneCallHistoryGetOptions {
	return func(o *phoneCallHistoryGetOptions) {
		o.phoneCallHistoryUUID = uuid
	}
}

// WithUserIdForPhoneCallHistory returns a PhoneCallHistoryGetOptions that
// lists call history records for the user identified by userId.
func WithUserIdForPhoneCallHistory(userId string) PhoneCallHistoryGetOptions {
	return func(o *phoneCallHistoryGetOptions) {
		o.userId = userId
	}
}

// WithPhoneCallHistoryQueryParameters returns a PhoneCallHistoryGetOptions
// that attaches the given query parameters to the call history request.
func WithPhoneCallHistoryQueryParameters(params *PhoneCallHistoryQueryParameters) PhoneCallHistoryGetOptions {
	return func(o *phoneCallHistoryGetOptions) {
		o.phoneCallHistoryQueryParameters = params
	}
}

// Get retrieves Zoom Phone call history records. Provide
// WithPhoneCallHistoryUUID to fetch a specific record by UUID,
// WithUserIdForPhoneCallHistory to retrieve a user's call log (paginated
// automatically), or no option to list the account-level call history
// (paginated automatically). It is an error to specify both a UUID and
// query parameters, or both a UUID and a user ID.
func (p *PhoneCallHistoryService) Get(ctx context.Context, opts ...PhoneCallHistoryGetOptions) ([]*models.CallHistory, *http.Response, error) {
	options := phoneCallHistoryGetOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	if options.phoneCallHistoryUUID != "" && options.phoneCallHistoryQueryParameters != nil {
		return nil, nil, fmt.Errorf("Cannot specify both phoneCallHistoryUUID and phoneCallHistoryQueryParameters")
	}

	if options.userId != "" && options.phoneCallHistoryUUID != "" {
		return nil, nil, fmt.Errorf("Cannot specify both userId and phoneCallHistoryUUID")
	}

	if options.phoneCallHistoryUUID != "" {
		endpoint := "/phone/call_history/" + url.PathEscape(options.phoneCallHistoryUUID)
		callHistory := &models.CallHistory{}
		res, err := p.client.request(ctx, http.MethodGet, endpoint, options.phoneCallHistoryQueryParameters, nil, callHistory)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		return []*models.CallHistory{callHistory}, res, nil
	}

	type response struct {
		*PaginationResponse
		CallHistory []*models.CallHistory `json:"call_history"`
	}

	queryResponse := &response{}
	var callHistory []*models.CallHistory

	if options.userId != "" {
		type userResponse struct {
			*PaginationResponse

			CallLogs []*models.CallHistory `json:"call_logs"`
		}
		endpoint := fmt.Sprintf("/phone/users/%s/call_history", url.PathEscape(options.userId))

		queryResponse := &userResponse{}
		res, err := p.client.request(ctx, http.MethodGet, endpoint, options.phoneCallHistoryQueryParameters, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}

		callHistory = append(callHistory, queryResponse.CallLogs...)

		type userCallHistoryPageQuery struct {
			*PhoneCallHistoryQueryParameters
			*PaginationOptions
		}
		for {
			if queryResponse.NextPageToken == "" {
				break
			}
			nextPageToken := queryResponse.NextPageToken
			pageQuery := &userCallHistoryPageQuery{
				PhoneCallHistoryQueryParameters: options.phoneCallHistoryQueryParameters,
				PaginationOptions: &PaginationOptions{
					NextPageToken: &nextPageToken,
				},
			}
			res, err = p.client.request(ctx, http.MethodGet, endpoint, pageQuery, nil, queryResponse)
			if err != nil {
				return nil, res, fmt.Errorf("Error making request: %w", err)
			}
			callHistory = append(callHistory, queryResponse.CallLogs...)
		}

		return callHistory, res, nil
	}

	endpoint := "/phone/call_history"

	res, err := p.client.request(ctx, http.MethodGet, endpoint, options.phoneCallHistoryQueryParameters, nil, queryResponse)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	callHistory = append(callHistory, queryResponse.CallHistory...)

	type callHistoryPageQuery struct {
		*PhoneCallHistoryQueryParameters
		*PaginationOptions
	}
	for {
		if queryResponse.NextPageToken == "" {
			break
		}
		nextPageToken := queryResponse.NextPageToken
		pageQuery := &callHistoryPageQuery{
			PhoneCallHistoryQueryParameters: options.phoneCallHistoryQueryParameters,
			PaginationOptions: &PaginationOptions{
				NextPageToken: &nextPageToken,
			},
		}
		res, err = p.client.request(ctx, http.MethodGet, endpoint, pageQuery, nil, queryResponse)
		if err != nil {
			return nil, res, fmt.Errorf("Error making request: %w", err)
		}
		callHistory = append(callHistory, queryResponse.CallHistory...)
	}

	return callHistory, res, nil
}

// AddClientCode associates a client-defined code string with the call log
// entry identified by callLogId. This is useful for tagging calls with
// internal billing or project codes.
func (p *PhoneCallHistoryService) AddClientCode(ctx context.Context, callLogId, clientCode string) (*http.Response, error) {
	type body struct {
		ClientCode string `json:"client_code"`
	}
	requestBody := &body{ClientCode: clientCode}
	res, err := p.client.request(ctx, http.MethodPost, fmt.Sprintf("/phone/call_history/%s/client_code", url.PathEscape(callLogId)), nil, requestBody, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d but got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}

// DeleteUserCallHistory removes the call log entry identified by callLogId
// from the call history of the user identified by userId.
func (p *PhoneCallHistoryService) DeleteUserCallHistory(ctx context.Context, userId, callLogId string) (*http.Response, error) {
	res, err := p.client.request(ctx, http.MethodDelete, fmt.Sprintf("/phone/users/%s/call_history/%s", url.PathEscape(userId), url.PathEscape(callLogId)), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d but got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}

// GetCallElement retrieves the call element record identified by
// callElementId.
func (p *PhoneCallHistoryService) GetCallElement(ctx context.Context, callElementId string) (*models.CallElement, *http.Response, error) {
	var callElement models.CallElement
	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/call_elements/%s", url.PathEscape(callElementId)), nil, nil, &callElement)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("Expected status code %d but got %d", http.StatusOK, res.StatusCode)
	}
	return &callElement, res, nil
}

// GetAICallSummary retrieves the AI-generated call summary identified by
// aiCallSummaryId for the user identified by userId.
func (p *PhoneCallHistoryService) GetAICallSummary(ctx context.Context, userId, aiCallSummaryId string) (*models.AICallSummary, *http.Response, error) {
	var aiCallSummary models.AICallSummary
	res, err := p.client.request(ctx, http.MethodGet, fmt.Sprintf("/phone/user/%s/ai_call_summary/%s", url.PathEscape(userId), url.PathEscape(aiCallSummaryId)), nil, nil, &aiCallSummary)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("Expected status code %d but got %d", http.StatusOK, res.StatusCode)
	}
	return &aiCallSummary, res, nil
}
