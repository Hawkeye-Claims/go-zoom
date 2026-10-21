package client_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
	"github.com/Hawkeye-Claims/go-zoom/zoom/enums"
	"github.com/Hawkeye-Claims/go-zoom/zoom/tokenmutex"
)

type phoneTransport func(*http.Request) (*http.Response, error)

func (f phoneTransport) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func phoneResponse(req *http.Request, status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Retry-After": {"12"}}, Body: io.NopCloser(strings.NewReader(body)), Request: req}
}

func newPhoneClient(t *testing.T, transport phoneTransport, opts ...client.ClientOption) *client.Client {
	t.Helper()
	c, err := client.NewClient(&http.Client{Transport: transport}, "account", "client", "secret", opts...)
	if err != nil {
		t.Fatal(err)
	}
	client.NewPhoneService(c)
	return c
}

func tokenResponse(req *http.Request) *http.Response {
	return phoneResponse(req, 200, `{"access_token":"token","expires_in":3600}`)
}

func TestPhoneVoicemailPages(t *testing.T) {
	ctx := context.Background()
	params := &client.PhoneVoicemailQueryParameters{From: "2026-10-19T00:00:00Z", To: "2026-10-20T00:00:00Z", Status: "all"}
	size := 1
	wantQuery := url.Values{"from": {params.From}, "to": {params.To}, "status": {"all"}, "page_size": {"1"}}
	calls := 0
	c := newPhoneClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "zoom.us" {
			return tokenResponse(req), nil
		}
		calls++
		if req.Method != "GET" || req.URL.Host != "api.zoom.us" || req.URL.EscapedPath() != "/v2/phone/users/owner%2Fid/voice_mails" {
			t.Fatalf("unexpected request: %s %s", req.Method, req.URL)
		}
		if !reflect.DeepEqual(req.URL.Query(), wantQuery) {
			t.Fatalf("query = %v, want %v", req.URL.Query(), wantQuery)
		}
		if calls == 1 {
			return phoneResponse(req, 200, `{"from":"2026-10-19T00:00:00Z","to":"2026-10-20T00:00:00Z","page_size":1,"page_count":2,"total_records":2,"next_page_token":"next+/=","voice_mails":[{"id":"vm1","call_element_id":"el1","call_history_id":"history1","call_id":"call1","caller_number":"+15551234567","callee_number":"100","caller_name":"Caller","date_time":"2026-10-19T12:00:00Z","duration":18,"status":"read","download_url":"https://example.invalid/audio"}]}`), nil
		}
		return phoneResponse(req, 200, `{"page_size":1,"total_records":2,"voice_mails":[{"id":"vm2","status":"unread"}]}`), nil
	})
	page, res, err := c.Phone.Voicemails.ListUserPage(ctx, "owner/id", params, &client.PaginationOptions{PageSize: &size})
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 || page.TotalRecords != 2 || page.PageCount != 2 || page.PageSize != 1 || page.From != params.From || page.To != params.To || len(page.VoiceMails) != 1 {
		t.Fatalf("page metadata: %+v", page)
	}
	vm := page.VoiceMails[0]
	if vm.ID != "vm1" || vm.CallElementID != "el1" || vm.CallHistoryID != "history1" || vm.CallID != "call1" || vm.CallerNumber != "+15551234567" || vm.CalleeNumber != "100" || vm.Duration != 18 || vm.DateTime.IsZero() || vm.Transcription != nil {
		t.Fatalf("voicemail: %+v", vm)
	}
	wantQuery.Set("next_page_token", page.NextPageToken)
	page, _, err = c.Phone.Voicemails.ListUserPage(ctx, "owner/id", params, &client.PaginationOptions{PageSize: &size, NextPageToken: &page.NextPageToken})
	if err != nil {
		t.Fatal(err)
	}
	if page.NextPageToken != "" || len(page.VoiceMails) != 1 || page.VoiceMails[0].ID != "vm2" || calls != 2 {
		t.Fatalf("last page: %+v, requests %d", page, calls)
	}
}

func TestPhoneVoicemailTranscription(t *testing.T) {
	for _, status := range []int{0, 1, 2, 4, 999, 1234} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			c := newPhoneClient(t, func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "zoom.us" {
					return tokenResponse(req), nil
				}
				if req.URL.EscapedPath() != "/v2/phone/voice_mails/vm%2Fid" || req.Method != "GET" {
					t.Fatalf("unexpected detail request: %s %s", req.Method, req.URL)
				}
				return phoneResponse(req, 200, fmt.Sprintf(`{"id":"vm/id","duration":18,"transcription":{"status":%d,"content":"text","engine":"otter"}}`, status)), nil
			})
			vm, _, err := c.Phone.Voicemails.Get(context.Background(), "vm/id")
			if err != nil {
				t.Fatal(err)
			}
			if vm.Transcription == nil || vm.Transcription.Status != status || vm.Transcription.Content != "text" || vm.Transcription.Engine != "otter" || vm.Duration != 18 {
				t.Fatalf("detail: %+v", vm)
			}
		})
	}
}

func TestPhoneReadErrors(t *testing.T) {
	for _, status := range []int{400, 401, 403, 404, 429, 500} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			c := newPhoneClient(t, func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "zoom.us" {
					return tokenResponse(req), nil
				}
				return phoneResponse(req, status, `{"code":400,"message":"expired cursor or denied"}`), nil
			})
			next := "expired"
			page, res, err := c.Phone.Voicemails.ListUserPage(context.Background(), "owner", nil, &client.PaginationOptions{NextPageToken: &next})
			if page != nil || err == nil || res == nil || res.StatusCode != status || res.Header.Get("Retry-After") != "12" {
				t.Fatalf("page=%v response=%v err=%v", page, res, err)
			}
			vm, res, err := c.Phone.Voicemails.Get(context.Background(), "vm")
			if vm != nil || err == nil || res == nil || res.StatusCode != status {
				t.Fatalf("vm=%v response=%v err=%v", vm, res, err)
			}
			history, res, err := c.Phone.CallHistory.ListUserPage(context.Background(), "owner", nil, nil)
			if history != nil || err == nil || res == nil || res.StatusCode != status {
				t.Fatalf("history=%v response=%v err=%v", history, res, err)
			}
		})
	}
}

func TestPhoneReadsUseExistingS2STokens(t *testing.T) {
	ctx := context.Background()
	store := tokenmutex.NewDefault()
	tokens, resources := 0, 0
	deny := false
	c := newPhoneClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "zoom.us" {
			tokens++
			user, password, ok := req.BasicAuth()
			if req.Method != "POST" || req.URL.Path != "/oauth/token" || req.URL.Query().Get("grant_type") != "account_credentials" || req.URL.Query().Get("account_id") != "account" || !ok || user != "client" || password != "secret" {
				t.Fatalf("unexpected token request: %s", req.URL)
			}
			return phoneResponse(req, 200, fmt.Sprintf(`{"access_token":"token%d","expires_in":3600}`, tokens)), nil
		}
		resources++
		if req.Header.Get("Authorization") != fmt.Sprintf("Bearer token%d", tokens) {
			t.Fatalf("authorization = %q", req.Header.Get("Authorization"))
		}
		if deny {
			return phoneResponse(req, 401, `{"code":124,"message":"expired"}`), nil
		}
		return phoneResponse(req, 200, `{"id":"vm"}`), nil
	}, client.WithToken(store))
	read := func() {
		t.Helper()
		if _, _, err := c.Phone.Voicemails.Get(ctx, "vm"); err != nil {
			t.Fatal(err)
		}
	}
	read()
	read()
	if tokens != 1 {
		t.Fatalf("cached token acquisitions = %d", tokens)
	}
	if err := store.Lock(ctx); err != nil {
		t.Fatal(err)
	}
	if err := store.Set(ctx, "expired", time.Now().Add(-time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := store.Unlock(ctx); err != nil {
		t.Fatal(err)
	}
	read()
	if tokens != 2 {
		t.Fatalf("expiry acquisitions = %d", tokens)
	}
	deny = true
	_, res, err := c.Phone.Voicemails.Get(ctx, "vm")
	if err == nil || res.StatusCode != 401 || tokens != 2 || resources != 4 {
		t.Fatalf("401 must not replay: tokens=%d resources=%d err=%v", tokens, resources, err)
	}
	deny = false
	read()
	if tokens != 3 || resources != 5 {
		t.Fatalf("next request must reacquire: tokens=%d resources=%d", tokens, resources)
	}
}

func TestPhoneUserHistoryPages(t *testing.T) {
	ctx := context.Background()
	params := &client.PhoneCallHistoryQueryParameters{From: "2026-10-19T00:00:00Z", To: "2026-10-20T00:00:00Z", CallResult: []enums.CallResult{enums.CallResult("no_answer")}}
	calls := 0
	c := newPhoneClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "zoom.us" {
			return tokenResponse(req), nil
		}
		calls++
		if req.URL.EscapedPath() != "/v2/phone/users/owner%2Fid/call_history" || req.URL.Query().Get("call_results") != "no_answer" || req.URL.Query().Has("call_result") || req.URL.Query().Get("from") != params.From || req.URL.Query().Get("to") != params.To {
			t.Fatalf("unexpected history query: %s", req.URL)
		}
		if req.URL.Query().Get("next_page_token") == "" {
			return phoneResponse(req, 200, `{"page_size":1,"total_records":2,"next_page_token":"next","call_elements":[{"call_element_id":"el1","call_history_uuid":"history1","call_id":"call1","direction":"inbound","result":"no_answer","result_reason":"ring_timeout","event":"ring","talk_time":0,"wait_time":15,"voicemail_id":"vm1","caller_did_number":"+15551234567","caller_email":"caller@example.com","callee_email":"owner@example.com","start_time":"2026-10-19T12:00:00Z"}],"call_logs":[{"id":"legacy","result":"voicemail"}]}`), nil
		}
		return phoneResponse(req, 200, `{"call_elements":[{"call_element_id":"el2","result":"answered_by_other","talk_time":12}]}`), nil
	})
	page, _, err := c.Phone.CallHistory.ListUserPage(ctx, "owner/id", params, nil)
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalRecords != 2 || page.NextPageToken != "next" || len(page.CallElements) != 1 || len(page.CallLogs) != 1 || page.CallLogs[0].Result != "voicemail" {
		t.Fatalf("page: %+v", page)
	}
	e := page.CallElements[0]
	if e.CallElementID != "el1" || e.CallHistoryUUID != "history1" || e.Result != "no_answer" || e.ResultReason != "ring_timeout" || e.Event != "ring" || e.WaitTime != 15 || e.VoicemailID != "vm1" || e.CallerDIDNumber != "+15551234567" || e.CalleeEmail != "owner@example.com" || e.StartTime.IsZero() {
		t.Fatalf("element: %+v", e)
	}
	history, _, err := c.Phone.CallHistory.Get(ctx, client.WithUserIdForPhoneCallHistory("owner/id"), client.WithPhoneCallHistoryQueryParameters(params))
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 2 || history[0].CallElementID != "el1" || history[1].Result != "answered_by_other" || history[1].TalkTime != 12 || calls != 3 {
		t.Fatalf("history: %+v requests %d", history, calls)
	}
}

func TestPhoneUserHistoryLegacyAndIncompletePages(t *testing.T) {
	for _, mode := range []string{"legacy", "error", "repeated", "empty_modern"} {
		t.Run(mode, func(t *testing.T) {
			calls := 0
			c := newPhoneClient(t, func(req *http.Request) (*http.Response, error) {
				if req.URL.Host == "zoom.us" {
					return tokenResponse(req), nil
				}
				calls++
				if calls == 1 {
					if mode == "empty_modern" {
						return phoneResponse(req, 200, `{"call_elements":[],"call_logs":[{"id":"duplicate"}]}`), nil
					}
					return phoneResponse(req, 200, `{"next_page_token":"next","call_logs":[{"id":"legacy1","call_path_id":"path1","result":"no_answer"}]}`), nil
				}
				if mode == "error" {
					return phoneResponse(req, 429, `{"code":429,"message":"rate limit"}`), nil
				}
				if mode == "repeated" {
					return phoneResponse(req, 200, `{"next_page_token":"next","call_logs":[{"id":"legacy1"}]}`), nil
				}
				return phoneResponse(req, 200, `{"call_logs":[{"id":"legacy2","result":"voicemail"}]}`), nil
			})
			history, res, err := c.Phone.CallHistory.Get(context.Background(), client.WithUserIdForPhoneCallHistory("owner"))
			switch mode {
			case "legacy":
				if err != nil || len(history) != 2 || history[0].CallHistoryUUID != "path1" || history[1].Result != "voicemail" {
					t.Fatalf("history=%+v err=%v", history, err)
				}
			case "empty_modern":
				if err != nil || len(history) != 0 || calls != 1 {
					t.Fatalf("empty modern history=%+v err=%v calls=%d", history, err, calls)
				}
			default:
				if err == nil || history != nil || calls != 2 || res == nil {
					t.Fatalf("partial history=%+v response=%v err=%v calls=%d", history, res, err, calls)
				}
			}
		})
	}
}

func TestGetCallElementUsesCurrentPath(t *testing.T) {
	c := newPhoneClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "zoom.us" {
			return tokenResponse(req), nil
		}
		if req.URL.EscapedPath() != "/v2/phone/call_element/el%2Fid" {
			t.Fatalf("path = %s", req.URL.EscapedPath())
		}
		return phoneResponse(req, 200, `{"call_element_id":"el/id","call_history_uuid":"history","talk_time":12,"wait_time":3}`), nil
	})
	element, _, err := c.Phone.CallHistory.GetCallElement(context.Background(), "el/id")
	if err != nil {
		t.Fatal(err)
	}
	if element.CallElementID != "el/id" || element.CallHistoryUUID != "history" || element.TalkTime != 12 || element.WaitTime != 3 {
		t.Fatalf("element=%+v", element)
	}
}

func TestPhonePageValidationAndCancellation(t *testing.T) {
	c := newPhoneClient(t, func(req *http.Request) (*http.Response, error) { return nil, req.Context().Err() })
	if _, _, err := c.Phone.Voicemails.ListUserPage(context.Background(), "", nil, nil); err == nil {
		t.Fatal("missing owner accepted")
	}
	if _, _, err := c.Phone.Voicemails.Get(context.Background(), ""); err == nil {
		t.Fatal("missing voicemail accepted")
	}
	if _, _, err := c.Phone.CallHistory.ListUserPage(context.Background(), "", nil, nil); err == nil {
		t.Fatal("missing history owner accepted")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := c.Phone.Voicemails.ListUserPage(ctx, "owner", nil, nil); err == nil {
		t.Fatal("canceled request succeeded")
	}
}

func TestPhonePagesRejectMalformedJSON(t *testing.T) {
	c := newPhoneClient(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.Host == "zoom.us" {
			return tokenResponse(req), nil
		}
		return phoneResponse(req, 200, `{"next_page_token":"next",`), nil
	})
	page, res, err := c.Phone.Voicemails.ListUserPage(context.Background(), "owner", nil, nil)
	if page != nil || res == nil || err == nil {
		t.Fatalf("malformed voicemail page=%v response=%v err=%v", page, res, err)
	}
	history, res, err := c.Phone.CallHistory.ListUserPage(context.Background(), "owner", nil, nil)
	if history != nil || res == nil || err == nil {
		t.Fatalf("malformed history page=%v response=%v err=%v", history, res, err)
	}
}
