package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
)

const (
	envZoomAccountID    = "ZOOM_ACCOUNT_ID"
	envZoomClientID     = "ZOOM_CLIENT_ID"
	envZoomClientSecret = "ZOOM_CLIENT_SECRET"
)

type cli struct {
	stdout io.Writer
	stderr io.Writer
}

type clientFlags struct {
	accountID    string
	clientID     string
	clientSecret string
	grantType    string
	redirectURI  string
	timeout      time.Duration
}

func main() {
	c := &cli{stdout: os.Stdout, stderr: os.Stderr}
	if err := c.run(os.Args[1:]); err != nil {
		fmt.Fprintf(c.stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func (c *cli) run(args []string) error {
	if len(args) == 0 {
		c.printRootUsage()
		return nil
	}

	switch args[0] {
	case "auth":
		return c.runAuth(args[1:])
	case "users":
		return c.runUsers(args[1:])
	case "meetings":
		return c.runMeetings(args[1:])
	case "phone":
		return c.runPhone(args[1:])
	case "help", "--help", "-h":
		c.printRootUsage()
		return nil
	default:
		c.printRootUsage()
		return fmt.Errorf("unknown command %q", args[0])
	}
}

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

func (c *cli) runAuthTest(args []string) error {
	fs := flag.NewFlagSet("auth test", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
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

func (c *cli) runAuthAuthorizationCode(args []string) error {
	fs := flag.NewFlagSet("auth authorization-code", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, false, true)
	authTimeout := fs.Duration("auth-timeout", 5*time.Minute, "How long to wait for OAuth callback")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
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

func (c *cli) runUsers(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "users commands: get")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runUsersGet(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "users commands: get")
		return nil
	default:
		return fmt.Errorf("unknown users command %q", args[0])
	}
}

func (c *cli) runUsersGet(args []string) error {
	fs := flag.NewFlagSet("users get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "User ID or email (optional; empty lists users)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	var users any
	if *userID != "" {
		users, _, err = zoomClient.Users.Get(context.Background(), client.WithUserId(*userID))
	} else {
		users, _, err = zoomClient.Users.Get(context.Background())
	}
	if err != nil {
		return err
	}

	return c.writeJSON(users)
}

func (c *cli) runMeetings(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "meetings commands: get, summary get")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runMeetingsGet(args[1:])
	case "summary":
		return c.runMeetingsSummary(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "meetings commands: get, summary get")
		return nil
	default:
		return fmt.Errorf("unknown meetings command %q", args[0])
	}
}

func (c *cli) runMeetingsGet(args []string) error {
	fs := flag.NewFlagSet("meetings get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	meetingID := fs.String("meeting-id", "", "Meeting ID")
	userID := fs.String("user-id", "", "User ID (required when listing meetings)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	if *meetingID == "" && *userID == "" {
		return errors.New("either --meeting-id or --user-id must be provided")
	}
	if *meetingID != "" && *userID != "" {
		return errors.New("--meeting-id and --user-id cannot be used together")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	var meetings any
	if *meetingID != "" {
		meetings, _, err = zoomClient.Meetings.Get(context.Background(), client.WithMeetingId(*meetingID))
	} else {
		meetings, _, err = zoomClient.Meetings.Get(context.Background(), client.WithMeetingUserId(*userID))
	}
	if err != nil {
		return err
	}

	return c.writeJSON(meetings)
}

func (c *cli) runMeetingsSummary(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "meetings summary commands: get")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runMeetingsSummaryGet(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "meetings summary commands: get")
		return nil
	default:
		return fmt.Errorf("unknown meetings summary command %q", args[0])
	}
}

func (c *cli) runMeetingsSummaryGet(args []string) error {
	fs := flag.NewFlagSet("meetings summary get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	meetingID := fs.String("meeting-id", "", "Meeting ID (optional)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	var summaries any
	if *meetingID != "" {
		summaries, _, err = zoomClient.Meetings.GetSummary(context.Background(), client.WithMeetingIdForSummary(*meetingID))
	} else {
		summaries, _, err = zoomClient.Meetings.GetSummary(context.Background())
	}
	if err != nil {
		return err
	}

	return c.writeJSON(summaries)
}

func (c *cli) runPhone(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone commands: call-history get, recordings get, settings get, users get")
		return nil
	}

	switch args[0] {
	case "call-history":
		return c.runPhoneCallHistory(args[1:])
	case "recordings":
		return c.runPhoneRecordings(args[1:])
	case "settings":
		return c.runPhoneSettings(args[1:])
	case "users":
		return c.runPhoneUsers(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "phone commands: call-history get, recordings get, settings get, users get")
		return nil
	default:
		return fmt.Errorf("unknown phone command %q", args[0])
	}
}

func (c *cli) runPhoneCallHistory(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone call-history commands: get")
		return nil
	}
	if args[0] != "get" {
		return fmt.Errorf("unknown phone call-history command %q", args[0])
	}

	fs := flag.NewFlagSet("phone call-history get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	uuid := fs.String("uuid", "", "Call history UUID")
	userID := fs.String("user-id", "", "User ID")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
	if *uuid != "" && *userID != "" {
		return errors.New("--uuid and --user-id cannot be used together")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	var history any
	if *uuid != "" {
		history, _, err = zoomClient.Phone.CallHistory.Get(context.Background(), client.WithPhoneCallHistoryUUID(*uuid))
	} else if *userID != "" {
		history, _, err = zoomClient.Phone.CallHistory.Get(context.Background(), client.WithUserIdForPhoneCallHistory(*userID))
	} else {
		history, _, err = zoomClient.Phone.CallHistory.Get(context.Background())
	}
	if err != nil {
		return err
	}

	return c.writeJSON(history)
}

func (c *cli) runPhoneRecordings(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone recordings commands: get")
		return nil
	}
	if args[0] != "get" {
		return fmt.Errorf("unknown phone recordings command %q", args[0])
	}

	fs := flag.NewFlagSet("phone recordings get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "User ID")
	callID := fs.String("call-id", "", "Call ID")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
	if *userID != "" && *callID != "" {
		return errors.New("--user-id and --call-id cannot be used together")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	var recordings any
	if *userID != "" {
		recordings, _, err = zoomClient.Phone.Recordings.Get(context.Background(), client.WithRecordingUserId(*userID))
	} else if *callID != "" {
		recordings, _, err = zoomClient.Phone.Recordings.Get(context.Background(), client.WithRecordingCallId(*callID))
	} else {
		recordings, _, err = zoomClient.Phone.Recordings.Get(context.Background())
	}
	if err != nil {
		return err
	}

	return c.writeJSON(recordings)
}

func (c *cli) runPhoneSettings(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone settings commands: get")
		return nil
	}
	if args[0] != "get" {
		return fmt.Errorf("unknown phone settings command %q", args[0])
	}

	fs := flag.NewFlagSet("phone settings get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	settings, _, err := zoomClient.Phone.Settings.Get(context.Background())
	if err != nil {
		return err
	}

	return c.writeJSON(settings)
}

func (c *cli) runPhoneUsers(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone users commands: get")
		return nil
	}
	if args[0] != "get" {
		return fmt.Errorf("unknown phone users command %q", args[0])
	}

	fs := flag.NewFlagSet("phone users get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "Phone user ID")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	var users any
	if *userID != "" {
		users, _, err = zoomClient.Phone.Users.Get(context.Background(), client.WithPhoneUserID(*userID))
	} else {
		users, _, err = zoomClient.Phone.Users.Get(context.Background())
	}
	if err != nil {
		return err
	}

	return c.writeJSON(users)
}

func (c *cli) bindClientFlags(fs *flag.FlagSet, includeGrantType bool, includeRedirectURI bool) *clientFlags {
	cfg := &clientFlags{}

	fs.StringVar(&cfg.accountID, "account-id", os.Getenv(envZoomAccountID), "Zoom account ID (or set ZOOM_ACCOUNT_ID)")
	fs.StringVar(&cfg.clientID, "client-id", os.Getenv(envZoomClientID), "Zoom client ID (or set ZOOM_CLIENT_ID)")
	fs.StringVar(&cfg.clientSecret, "client-secret", os.Getenv(envZoomClientSecret), "Zoom client secret (or set ZOOM_CLIENT_SECRET)")
	fs.DurationVar(&cfg.timeout, "http-timeout", 30*time.Second, "HTTP client timeout")

	if includeGrantType {
		fs.StringVar(&cfg.grantType, "grant-type", "account_credentials", "OAuth grant type: account_credentials or authorization_code")
	}
	if includeRedirectURI {
		fs.StringVar(&cfg.redirectURI, "redirect-uri", "", "OAuth redirect URI (required for authorization_code)")
	}

	return cfg
}

func (c *cli) newServiceClient(cfg *clientFlags) (*client.Client, error) {
	zoomClient, err := c.newClient(cfg)
	if err != nil {
		return nil, err
	}
	if cfg.grantType == "authorization_code" {
		return nil, errors.New("service commands support grant-type=account_credentials; use auth authorization-code for interactive OAuth")
	}
	return zoomClient, nil
}

func (c *cli) newClient(cfg *clientFlags) (*client.Client, error) {
	if strings.TrimSpace(cfg.accountID) == "" {
		return nil, fmt.Errorf("missing required account ID (use --account-id or %s)", envZoomAccountID)
	}
	if strings.TrimSpace(cfg.clientID) == "" {
		return nil, fmt.Errorf("missing required client ID (use --client-id or %s)", envZoomClientID)
	}
	if strings.TrimSpace(cfg.clientSecret) == "" {
		return nil, fmt.Errorf("missing required client secret (use --client-secret or %s)", envZoomClientSecret)
	}

	grantType := strings.TrimSpace(cfg.grantType)
	if grantType == "" {
		grantType = "account_credentials"
	}
	if grantType != "account_credentials" && grantType != "authorization_code" {
		return nil, fmt.Errorf("unsupported grant type %q", grantType)
	}

	opts := []client.ClientOption{client.WithGrantType(grantType)}
	if grantType == "authorization_code" {
		if strings.TrimSpace(cfg.redirectURI) == "" {
			return nil, errors.New("redirect URI is required when using authorization_code grant type")
		}
		opts = append(opts, client.WithRedirectURI(cfg.redirectURI))
	}

	zoomClient, err := client.NewClient(
		&http.Client{Timeout: cfg.timeout},
		cfg.accountID,
		cfg.clientID,
		cfg.clientSecret,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return zoomClient, nil
}

func (c *cli) writeJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}
	_, _ = fmt.Fprintln(c.stdout, string(b))
	return nil
}

func (c *cli) printRootUsage() {
	fmt.Fprintln(c.stderr, `go-zoom CLI wrapper

Usage:
  go run ./cmd/go-zoom <command> [subcommand] [flags]

Commands:
  auth                     Authentication helpers (SDK auth options)
  users get                Get one user or list users
  meetings get             Get one meeting or list a user's meetings
  meetings summary get     Get one meeting summary or list summaries
  phone call-history get   Get account/user/single call history
  phone recordings get     Get account/user/call recordings
  phone settings get       Get account phone settings
  phone users get          Get one phone user or list phone users

Use "<command> --help" for command-specific flags.`)
}

func (c *cli) printAuthUsage() {
	fmt.Fprintln(c.stderr, `auth commands:
  auth test                 Validate auth config (uses SDK client)
  auth authorization-code   Start local OAuth helper flow

Examples:
  go run ./cmd/go-zoom auth test --grant-type account_credentials
  go run ./cmd/go-zoom auth authorization-code --redirect-uri http://localhost:8080/oauth/callback`)
}
