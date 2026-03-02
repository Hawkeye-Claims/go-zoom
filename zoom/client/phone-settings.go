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
// be updated via the Update method. All fields are optional pointers; only
// non-nil fields are included in the PATCH request. Fields left as nil are
// not sent and their corresponding settings on the server remain unchanged.
type SettingsAttributes struct {
	// BillingAccountId is the ID of the billing account to associate with the
	// phone account.
	BillingAccountId *string
	// BYOC, when true, enables Bring Your Own Carrier for the account.
	BYOC *bool
	// MultipleSites, when true, enables the multiple-sites feature.
	MultipleSites *bool
	// SiteCode, when true, enables site codes for the account.
	SiteCode *bool
	// ShortExtensionLength sets the length of short extensions when site codes
	// are enabled.
	ShortExtensionLength *int
	// ShowDeviceIPForCallLog, when true, shows device IP addresses in call
	// logs.
	ShowDeviceIPForCallLog *bool
}

// Update patches the Zoom Phone account settings with the values of any
// non-nil fields in attributes. Only fields explicitly set on
// SettingsAttributes are sent in the PATCH request; nil fields are omitted
// and their corresponding settings on the server are left unchanged.
func (s *PhoneSettingsService) Update(ctx context.Context, attributes *SettingsAttributes) (*http.Response, error) {
	if attributes == nil {
		return nil, fmt.Errorf("attributes cannot be nil")
	}

	type siteCodeBody struct {
		Enable               *bool `json:"enable,omitempty"`
		ShortExtensionLength *int  `json:"short_extension_length,omitempty"`
	}
	type multipleSitesBody struct {
		Enabled  *bool         `json:"enabled,omitempty"`
		SiteCode *siteCodeBody `json:"site_code,omitempty"`
	}
	type body struct {
		BillingAccount *struct {
			ID string `json:"id,omitempty"`
		} `json:"billing_account,omitempty"`
		BYOC *struct {
			Enable bool `json:"enable"`
		} `json:"byoc,omitempty"`
		MultipleSites          *multipleSitesBody `json:"multiple_sites,omitempty"`
		ShowDeviceIPForCallLog *struct {
			Enable bool `json:"enable"`
		} `json:"show_device_ip_for_call_log,omitempty"`
	}

	requestBody := &body{}

	if attributes.BillingAccountId != nil {
		requestBody.BillingAccount = &struct {
			ID string `json:"id,omitempty"`
		}{ID: *attributes.BillingAccountId}
	}
	if attributes.BYOC != nil {
		requestBody.BYOC = &struct {
			Enable bool `json:"enable"`
		}{Enable: *attributes.BYOC}
	}
	if attributes.MultipleSites != nil || attributes.SiteCode != nil || attributes.ShortExtensionLength != nil {
		ms := &multipleSitesBody{Enabled: attributes.MultipleSites}
		if attributes.SiteCode != nil || attributes.ShortExtensionLength != nil {
			sc := &siteCodeBody{}
			if attributes.SiteCode != nil {
				sc.Enable = attributes.SiteCode
			}
			if attributes.ShortExtensionLength != nil {
				sc.ShortExtensionLength = attributes.ShortExtensionLength
			}
			ms.SiteCode = sc
		}
		requestBody.MultipleSites = ms
	}
	if attributes.ShowDeviceIPForCallLog != nil {
		requestBody.ShowDeviceIPForCallLog = &struct {
			Enable bool `json:"enable"`
		}{Enable: *attributes.ShowDeviceIPForCallLog}
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
