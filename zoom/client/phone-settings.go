package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TheSlowpes/go-zoom/zoom/models"
)

// PhoneSettingsServicer is the interface implemented by PhoneSettingsService.
// It declares the read operation for Zoom Phone account-level settings.
// Note: the Update method exists on PhoneSettingsService but is intentionally
// not included in this interface.
type PhoneSettingsServicer interface {
	// Get retrieves the current Zoom Phone account settings.
	Get(ctx context.Context) (*models.PhoneAccountSettings, *http.Response, error)
}

// PhoneSettingsService implements PhoneSettingsServicer and provides access to
// Zoom Phone account settings API endpoints.
type PhoneSettingsService struct {
	client *Client
}

// Compile-time assertion that PhoneSettingsService satisfies the
// PhoneSettingsServicer interface.
var _ PhoneSettingsServicer = (*PhoneSettingsService)(nil)

// Get retrieves the current Zoom Phone account-level settings.
func (s *PhoneSettingsService) Get(ctx context.Context) (*models.PhoneAccountSettings, *http.Response, error) {
	var settings *models.PhoneAccountSettings
	res, err := s.client.request(ctx, http.MethodGet, "/phone/settings", nil, nil, &settings)
	if err != nil {
		return nil, res, fmt.Errorf("Error making request: %w", err)
	}

	return settings, res, nil
}

// SettingsAttributes holds the subset of Zoom Phone account settings that can
// be updated via the Update method.
type SettingsAttributes struct {
	// BillingAccountId is the ID of the billing account to associate with the
	// phone account.
	BillingAccountId string
	// BYOC, when true, enables Bring Your Own Carrier for the account.
	BYOC bool
	// MultipleSites, when true, enables the multiple-sites feature.
	MultipleSites bool
	// SiteCode, when true, enables site codes for the account.
	SiteCode bool
	// ShortExtensionLength sets the length of short extensions when site codes
	// are enabled.
	ShortExtensionLength int
	// ShowDeviceIPForCallLog, when true, shows device IP addresses in call
	// logs.
	ShowDeviceIPForCallLog bool
}

// Update patches the Zoom Phone account settings with the values contained in
// attributes. Only the fields present in SettingsAttributes are updated; all
// other settings are left unchanged.
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
