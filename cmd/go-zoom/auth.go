package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
)

// runAuth dispatches auth subcommands.
func (c *cli) runAuth(args []string) error {
	if len(args) == 0 {
		c.printAuthUsage()
		return nil
	}

	switch args[0] {
	case "test":
		return c.runAuthTest(args[1:])
	case "authorization-code":
		return c.runAuthAuthorizationCode(args[1:])
	case "help", "--help", "-h":
		c.printAuthUsage()
		return nil
	default:
		c.printAuthUsage()
		return fmt.Errorf("unknown auth command %q", args[0])
	}
}

// runAuthTest validates auth configuration and tests account-credentials auth.
func (c *cli) runAuthTest(args []string) error {
	fs := flag.NewFlagSet("auth test", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}

	zoomClient, err := c.newClient(cfg)
	if err != nil {
		return err
	}

	if cfg.grantType == "authorization_code" {
		fmt.Fprintln(c.stdout, "Authorization Code configuration is valid.")
		fmt.Fprintln(c.stdout, "Run `auth authorization-code` to complete interactive OAuth.")
		return nil
	}

	users, _, err := zoomClient.Users.Get(context.Background(), client.WithUserId("me"))
	if err != nil {
		return fmt.Errorf("authentication check failed: %w", err)
	}

	fmt.Fprintf(c.stdout, "Authentication succeeded. Retrieved %d user record(s).\n", len(users))
	return nil
}

// runAuthAuthorizationCode starts a local server for authorization-code OAuth.
func (c *cli) runAuthAuthorizationCode(args []string) error {
	fs := flag.NewFlagSet("auth authorization-code", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, false, true)
	authTimeout := fs.Duration("auth-timeout", 5*time.Minute, "How long to wait for OAuth callback")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}

	cfg.grantType = "authorization_code"
	zoomClient, err := c.newClient(cfg)
	if err != nil {
		return err
	}

	redirectURL, err := url.Parse(cfg.redirectURI)
	if err != nil {
		return fmt.Errorf("invalid redirect URI: %w", err)
	}
	if redirectURL.Host == "" {
		return errors.New("redirect URI must include a host and port")
	}

	callbackPath := redirectURL.Path
	if callbackPath == "" {
		callbackPath = "/"
	}

	done := make(chan struct{})
	mux := http.NewServeMux()
	mux.Handle("/oauth/login", zoomClient.RequestAuthorization())

	callbackHandler := zoomClient.HandleOAuthCallback()
	mux.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		callbackHandler.ServeHTTP(w, r)
		if r.URL.Query().Get("code") != "" {
			select {
			case <-done:
			default:
				close(done)
			}
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("go-zoom CLI OAuth helper is running"))
	})

	srv := &http.Server{
		Addr:              redirectURL.Host,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if listenErr := srv.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			errCh <- listenErr
		}
	}()

	fmt.Fprintf(c.stdout, "Open this URL in your browser to authenticate:\nhttp://%s/oauth/login\n", redirectURL.Host)
	fmt.Fprintf(c.stdout, "Waiting up to %s for callback on %s ...\n", authTimeout.String(), cfg.redirectURI)

	select {
	case <-done:
		_, _, err = zoomClient.Users.Get(context.Background(), client.WithUserId("me"))
		if err != nil {
			_ = srv.Shutdown(context.Background())
			return fmt.Errorf("oauth callback received but token validation failed: %w", err)
		}
		fmt.Fprintln(c.stdout, "Authorization succeeded.")
	case listenErr := <-errCh:
		return fmt.Errorf("oauth helper server failed: %w", listenErr)
	case <-time.After(*authTimeout):
		_ = srv.Shutdown(context.Background())
		return errors.New("timed out waiting for OAuth callback")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	return nil
}

// printAuthUsage prints usage for auth commands.
func (c *cli) printAuthUsage() {
	fmt.Fprintln(c.stderr, `auth commands:
  auth test                 Validate auth config (uses SDK client)
  auth authorization-code   Start local OAuth helper flow

Examples:
  go run ./cmd/go-zoom auth test --grant-type account_credentials
  go run ./cmd/go-zoom auth authorization-code --redirect-uri http://localhost:8080/oauth/callback`)
}

// ensureNoUnexpectedArgs validates that no trailing positional arguments remain.
func ensureNoUnexpectedArgs(args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf("unexpected arguments: %s", strings.Join(args, " "))
}
