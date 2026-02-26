package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/TheSlowpes/go-zoom/zoom/enums"
	"github.com/TheSlowpes/go-zoom/zoom/models"
)

type PhoneService struct {
	client      *Client
	CallHistory *PhoneCallHistoryService
}

func NewPhoneService(client *Client) {
	client.Phone = &PhoneService{client: client}
	client.Phone.CallHistory = &PhoneCallHistoryService{client: client}
}

type PhoneCallHistoryServicer interface {
	Get(ctx context.Context, opts ...PhoneCallHistoryGetOptions) ([]*models.CallHistory, *http.Response, error)
	AddClientCode(callLogId, clientCode string) (*http.Response, error)
	DeleteUserCallHistory(userId, callLogId string) (*http.Response, error)
}

type PhoneCallHistoryService struct {
	client *Client
}

var _ PhoneCallHistoryServicer = (*PhoneCallHistoryService)(nil)

type PhoneCallHistoryGetOptions func(*phoneCallHistoryGetOptions)

type phoneCallHistoryGetOptions struct {
	phoneCallHistoryUUID            string
	userId                          string
	phoneCallHistoryQueryParameters *PhoneCallHistoryQueryParameters
}

type PhoneCallHistoryQueryParameters struct {
	From            string                `url:"from,omitempty"`
	To              string                `url:"to,omitempty"`
	Keyword         string                `url:"keyword,omitempty"`
	Directions      []enums.Direction     `url:"directions,omitempty"`
	ConnectTypes    []enums.ConnectType   `url:"connect_types,omitempty"`
	NumberTypes     []enums.NumberType    `url:"number_types,omitempty"`
	CallTypes       []enums.CallType      `url:"call_types,omitempty"`
	ExtensionTypes  []enums.ExtensionType `url:"extension_types,omitempty"`
	CallResult      []enums.CallResult    `url:"call_result,omitempty"`
	GroupIDs        []string              `url:"group_ids,omitempty"`
	SiteIDs         []string              `url:"site_ids,omitempty"`
	Department      string                `url:"department,omitempty"`
	CostCenter      string                `url:"cost_center,omitempty"`
	TimeType        enums.TimeType        `url:"time_type,omitempty"`
	RecordingStatus enums.RecordingStatus `url:"recording_status,omitempty"`
	WithVoicemail   bool                  `url:"with_voicemail,omitempty"`
}

func WithPhoneCallHistoryUUID(uuid string) PhoneCallHistoryGetOptions {
	return func(o *phoneCallHistoryGetOptions) {
		o.phoneCallHistoryUUID = uuid
	}
}

func WithUserIdForPhoneCallHistory(userId string) PhoneCallHistoryGetOptions {
	return func(o *phoneCallHistoryGetOptions) {
		o.userId = userId
	}
}

func WithPhoneCallHistoryQueryParameters(params *PhoneCallHistoryQueryParameters) PhoneCallHistoryGetOptions {
	return func(o *phoneCallHistoryGetOptions) {
		o.phoneCallHistoryQueryParameters = params
	}
}

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

func (p *PhoneCallHistoryService) AddClientCode(callLogId, clientCode string) (*http.Response, error) {
	type body struct {
		ClientCode string `json:"client_code"`
	}
	requestBody := &body{ClientCode: clientCode}
	res, err := p.client.request(context.Background(), http.MethodPost, fmt.Sprintf("/phone/call_history/%s/client_code", url.PathEscape(callLogId)), nil, requestBody, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d but got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}

func (p *PhoneCallHistoryService) DeleteUserCallHistory(userId, callLogId string) (*http.Response, error) {
	res, err := p.client.request(context.Background(), http.MethodDelete, fmt.Sprintf("/phone/users/%s/call_history/%s", url.PathEscape(userId), url.PathEscape(callLogId)), nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d but got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}

func (p *PhoneCallHistoryService) GetCallElement(callElementId string) (*models.CallElement, *http.Response, error) {
	var callElement models.CallElement
	res, err := p.client.request(context.Background(), http.MethodGet, fmt.Sprintf("/phone/call_elements/%s", url.PathEscape(callElementId)), nil, nil, &callElement)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("Expected status code %d but got %d", http.StatusOK, res.StatusCode)
	}
	return &callElement, res, nil
}

func (p *PhoneCallHistoryService) GetAICallSummary(userId, aiCallSummaryId string) (*models.AICallSummary, *http.Response, error) {
	var aiCallSummary models.AICallSummary
	res, err := p.client.request(context.Background(), http.MethodGet, fmt.Sprintf("/phone/user/%s/ai_call_summary/%s", url.PathEscape(userId), url.PathEscape(aiCallSummaryId)), nil, nil, &aiCallSummary)
}
