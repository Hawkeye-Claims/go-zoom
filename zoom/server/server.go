package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HandlerFunc func(payload []byte) error

type WebhookServer struct {
	listenAddr string
	mux        *http.ServeMux
	registry   map[string]HandlerFunc
	token      string
}

func makeHandler[T any](processFn func(T) error) HandlerFunc {
	return func(payload []byte) error {
		var envelope Notification[T]
		if err := json.Unmarshal(payload, &envelope); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %w", err)
		}
		return processFn(envelope.Payload)
	}
}

type HandlerOption func(*WebhookServer)

func WithHandler[T any](eventType string, fn func(T) error) HandlerOption {
	return func(s *WebhookServer) {
		s.registry[eventType] = makeHandler(fn)
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

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	message := fmt.Sprintf("v0:%s:%s", requestTimestamp, string(bodyBytes))

	expected := fmt.Sprintf("v0=%s", generateHMAC(message, s.token))

	if !hmac.Equal([]byte(expected), []byte(requestSignature)) {
		return
	}

	if err = json.Unmarshal(bodyBytes, &header); err != nil {
		return
	}

	switch header.Event {
	case "endpoint.url_validation":

	}
}

func (s *WebhookServer) Start() error {
	return http.ListenAndServe(s.listenAddr, s.mux)
}

func (s *WebhookServer) handleValidateToken(bodyBytes []byte) error {

}

func generateHMAC(message, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := h.Sum(nil)
	return hex.EncodeToString(signature)
}
