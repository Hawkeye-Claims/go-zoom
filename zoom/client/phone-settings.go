package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TheSlowpes/go-zoom/zoom/models"
)

type PhoneSettingsServicer interface {
	Get(ctx context.Context) (*models.PhoneAccountSettings, *http.Response, error)
}

type PhoneSettingsService struct {
	client *Client
}

var _ PhoneSettingsServicer = (*PhoneSettingsService)(nil)

func (s *PhoneSettingsService) Get(ctx context.Context) (*models.PhoneAccountSettings, *http.Response, error) {
	var settings *models.PhoneAccountSettings
	res, err := s.client.request(ctx, http.MethodGet, "/phone/settings", nil, nil, &settings)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return settings, res, nil
}

type SettingsAttributes struct {
	BillingAccountId       string
	BYOC                   bool
	MultipleSites          bool
	SiteCode               bool
	ShortExtensionLength   int
	ShowDeviceIPForCallLog bool
}

func (s *PhoneSettingsService) Update(ctx context.Context, attributes *SettingsAttributes) (*http.Response, error) {
	type body struct {
		BillingAccount models.BillingAccount `json:"billing_account"`
		BYOC           models.BYOC           `json:"byoc"`
		MultipleSites  struct {
			Enabled  bool `json:"enabled"`
			SiteCode struct {
				Enable               bool `json:"enable"`
				ShortExtensionLength int  `json:"short_extension_length,omitempty"`
			} `json:"site_code"`
		} `json:"multiple_sites"`
		ShowDeviceIPForCallLog models.ShowDeviceIPForCallLog `json:"show_device_ip_for_call_log"`
	}

	requestBody := &body{
		BillingAccount: models.BillingAccount{
			ID: attributes.BillingAccountId,
		},
		BYOC: models.BYOC{
			Enable: attributes.BYOC,
		},
		MultipleSites: struct {
			Enabled  bool `json:"enabled"`
			SiteCode struct {
				Enable               bool `json:"enable"`
				ShortExtensionLength int  `json:"short_extension_length,omitempty"`
			} `json:"site_code"`
		}{
			Enabled: attributes.MultipleSites,
			SiteCode: struct {
				Enable               bool `json:"enable"`
				ShortExtensionLength int  `json:"short_extension_length,omitempty"`
			}{
				Enable:               attributes.SiteCode,
				ShortExtensionLength: attributes.ShortExtensionLength,
			},
		},
	}

	res, err := s.client.request(ctx, http.MethodPatch, "/phone/settings", nil, requestBody, nil)
	if err != nil {
		return nil, fmt.Errorf("Error making request: %w", err)
	}
	if res.StatusCode != http.StatusNoContent {
		return res, fmt.Errorf("Expected status code %d, got %d", http.StatusNoContent, res.StatusCode)
	}
	return res, nil
}
