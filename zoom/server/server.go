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

type HandlerFunc func(payload []byte) error

type WebhookServer struct {
	listenAddr string
	mux        *http.ServeMux
	registry   map[string]HandlerFunc
	token      string
}

const timestampTolerance = 5 * time.Minute

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

type HandlerOption func(*WebhookServer)

func WithHandler[T any](eventType string, ch chan T) HandlerOption {
	return func(s *WebhookServer) {
		s.registry[eventType] = makeHandler(ch)
	}
}

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

func (s *WebhookServer) Start() error {
	return http.ListenAndServe(s.listenAddr, s.mux)
}

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

func generateHMAC(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := h.Sum(nil)
	return hex.EncodeToString(signature)
}
