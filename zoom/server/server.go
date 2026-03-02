package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// HandlerFunc is the function signature for processing a raw Zoom webhook
// event payload. Implementations should unmarshal the JSON bytes and handle
// the event accordingly. Returning a non-nil error causes the server to
// respond with HTTP 500.
type HandlerFunc func(payload []byte) error

// WebhookServer is an HTTP server that receives, verifies, and dispatches
// Zoom webhook event notifications. It validates the HMAC-SHA256 signature and
// request timestamp on every incoming request before dispatching to the
// registered HandlerFunc for the event type.
type WebhookServer struct {
	listenAddr string
	mux        *http.ServeMux
	registry   map[string]HandlerFunc
	token      string
}

// timestampTolerance is the maximum age (or future skew) allowed for the
// x-zm-request-timestamp header. Requests outside this window are rejected to
// prevent replay attacks.
const timestampTolerance = 5 * time.Minute

// makeHandler returns a HandlerFunc that unmarshals the raw payload into a
// Notification[T], extracts the typed inner payload, and sends it to ch. If
// the channel is full the event is dropped (non-blocking send).
func makeHandler[T any](ch chan<- T) HandlerFunc {
	return func(payload []byte) error {
		var envelope Notification[T]
		if err := json.Unmarshal(payload, &envelope); err != nil {
			return fmt.Errorf("failed to unmarshal event payload: %w", err)
		}
		select {
		case ch <- envelope.Payload:
		default:
		}
		return nil
	}
}

// HandlerOption is a functional option for configuring a WebhookServer at
// construction time via NewWebhookServer.
type HandlerOption func(*WebhookServer)

// WithHandler returns a HandlerOption that registers ch to receive the typed
// payload for every Zoom webhook event matching eventType. T must match the
// payload type for the given event (e.g. MeetingEvent for "meeting.started").
func WithHandler[T any](eventType string, ch chan T) HandlerOption {
	return func(s *WebhookServer) {
		s.registry[eventType] = makeHandler(ch)
	}
}

// NewWebhookServer creates a new WebhookServer that listens on listenAddr and
// handles webhook requests at webhookPath. secretToken is the Zoom webhook
// secret token used for HMAC-SHA256 signature verification. Additional event
// handlers can be registered via HandlerOption values (e.g. WithHandler).
func NewWebhookServer(listenAddr, webhookPath, secretToken string, opts ...HandlerOption) *WebhookServer {
	server := &WebhookServer{
		listenAddr: listenAddr,
		mux:        http.NewServeMux(),
		token:      secretToken,
		registry:   make(map[string]HandlerFunc),
	}
	server.mux.HandleFunc(webhookPath, server.processEvent)
	for _, opt := range opts {
		opt(server)
	}
	return server
}

// processEvent is the internal HTTP handler that validates incoming Zoom
// webhook requests and dispatches them to the appropriate registered
// HandlerFunc. It verifies the HMAC-SHA256 signature and request timestamp
// before processing, and responds to URL validation challenges automatically.
func (s *WebhookServer) processEvent(w http.ResponseWriter, r *http.Request) {
	var header webhookHeader
	defer r.Body.Close()

	requestTimestamp := r.Header.Get("x-zm-request-timestamp")
	requestSignature := r.Header.Get("x-zm-signature")

	if requestTimestamp == "" {
		http.Error(w, "missing request timestamp", http.StatusUnauthorized)
		return
	}

	tsMs, err := strconv.ParseInt(requestTimestamp, 10, 64)
	if err != nil {
		http.Error(w, "invalid request timestamp", http.StatusUnauthorized)
		return
	}

	age := time.Since(time.UnixMilli(tsMs))
	if age.Abs() > timestampTolerance {
		http.Error(w, "request timestamp outside allowed window", http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("v0:%s:%s", requestTimestamp, string(bodyBytes))
	expected := fmt.Sprintf("v0=%s", generateHMAC(message, s.token))

	if !hmac.Equal([]byte(expected), []byte(requestSignature)) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	if err = json.Unmarshal(bodyBytes, &header); err != nil {
		http.Error(w, "failed to decode event", http.StatusBadRequest)
		return
	}

	switch header.Event {
	case "endpoint.url_validation":
		if err = s.handleValidateToken(w, bodyBytes); err != nil {
			http.Error(w, "failed to handle url validation", http.StatusInternalServerError)
		}
	default:
		handler, ok := s.registry[header.Event]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = handler(bodyBytes); err != nil {
			http.Error(w, "failed to handle event", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// Start begins listening for incoming webhook requests on the configured
// address. It blocks until the server encounters a fatal error, which it
// returns.
func (s *WebhookServer) Start() error {
	return http.ListenAndServe(s.listenAddr, s.mux)
}

// handleValidateToken handles the Zoom "endpoint.url_validation" challenge by
// computing the HMAC-SHA256 of the plain token and writing the expected JSON
// response back to Zoom.
func (s *WebhookServer) handleValidateToken(w http.ResponseWriter, bodyBytes []byte) error {
	var envelope Notification[validateTokenPayload]
	if err := json.Unmarshal(bodyBytes, &envelope); err != nil {
		return fmt.Errorf("failed to unmarshal url validation payload: %w", err)
	}

	resp := validateTokenResponse{
		PlainToken:     envelope.Payload.PlainToken,
		EncryptedToken: generateHMAC(envelope.Payload.PlainToken, s.token),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

// generateHMAC computes the HMAC-SHA256 of message using secret and returns
// the result as a lowercase hexadecimal string.
func generateHMAC(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := h.Sum(nil)
	return hex.EncodeToString(signature)
}
