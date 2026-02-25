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
}

type PhoneCallHistoryService struct {
	client *Client
}

var _ PhoneCallHistoryServicer = (*PhoneCallHistoryService)(nil)

type PhoneCallHistoryGetOptions func(*phoneCallHistoryGetOptions)

type phoneCallHistoryGetOptions struct {
	phoneCallHistoryUUID            string
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
