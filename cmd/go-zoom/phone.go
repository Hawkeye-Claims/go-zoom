package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
)

// runPhone dispatches phone subcommands.
func (c *cli) runPhone(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone commands: call-history, recordings, settings, users")
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
		fmt.Fprintln(c.stderr, "phone commands: call-history, recordings, settings, users")
		return nil
	default:
		return fmt.Errorf("unknown phone command %q", args[0])
	}
}

// runPhoneCallHistory dispatches phone call-history subcommands.
func (c *cli) runPhoneCallHistory(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone call-history commands: get, add-client-code, delete, call-element get, ai-summary get")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runPhoneCallHistoryGet(args[1:])
	case "add-client-code":
		return c.runPhoneCallHistoryAddClientCode(args[1:])
	case "delete":
		return c.runPhoneCallHistoryDelete(args[1:])
	case "call-element":
		return c.runPhoneCallHistoryCallElement(args[1:])
	case "ai-summary":
		return c.runPhoneCallHistoryAICallSummary(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "phone call-history commands: get, add-client-code, delete, call-element get, ai-summary get")
		return nil
	default:
		return fmt.Errorf("unknown phone call-history command %q", args[0])
	}
}

// runPhoneCallHistoryGet fetches account, user, or UUID-targeted call history.
func (c *cli) runPhoneCallHistoryGet(args []string) error {
	fs := flag.NewFlagSet("phone call-history get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	uuid := fs.String("uuid", "", "Call history UUID")
	userID := fs.String("user-id", "", "User ID")
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional get query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional get query parameters")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *uuid != "" && *userID != "" {
		return errors.New("--uuid and --user-id cannot be used together")
	}
	if *uuid != "" && (*queryJSON != "" || *queryJSONFile != "") {
		return errors.New("--uuid cannot be used with query parameter input")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	opts := make([]client.PhoneCallHistoryGetOptions, 0, 2)
	if *uuid != "" {
		opts = append(opts, client.WithPhoneCallHistoryUUID(*uuid))
	} else if *userID != "" {
		opts = append(opts, client.WithUserIdForPhoneCallHistory(*userID))
	}

	queryParameters, err := readOptionalJSONInputWithFlags[client.PhoneCallHistoryQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
	if err != nil {
		return err
	}
	if queryParameters != nil {
		opts = append(opts, client.WithPhoneCallHistoryQueryParameters(queryParameters))
	}

	history, _, err := zoomClient.Phone.CallHistory.Get(context.Background(), opts...)
	if err != nil {
		return err
	}

	return c.writeJSON(history)
}

// runPhoneCallHistoryAddClientCode tags a call log entry with a client code.
func (c *cli) runPhoneCallHistoryAddClientCode(args []string) error {
	fs := flag.NewFlagSet("phone call-history add-client-code", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	callLogID := fs.String("call-log-id", "", "Call log ID (required)")
	clientCode := fs.String("client-code", "", "Client code to associate (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *callLogID == "" {
		return errors.New("--call-log-id is required")
	}
	if *clientCode == "" {
		return errors.New("--client-code is required")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	if _, err := zoomClient.Phone.CallHistory.AddClientCode(context.Background(), *callLogID, *clientCode); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "client code added to call log %s\n", *callLogID)
	return nil
}

// runPhoneCallHistoryDelete removes a call log entry from a user's history.
func (c *cli) runPhoneCallHistoryDelete(args []string) error {
	fs := flag.NewFlagSet("phone call-history delete", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "User ID (required)")
	callLogID := fs.String("call-log-id", "", "Call log ID (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *userID == "" {
		return errors.New("--user-id is required")
	}
	if *callLogID == "" {
		return errors.New("--call-log-id is required")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	if _, err := zoomClient.Phone.CallHistory.DeleteUserCallHistory(context.Background(), *userID, *callLogID); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "call log %s deleted from user %s\n", *callLogID, *userID)
	return nil
}

// runPhoneCallHistoryCallElement dispatches call-element subcommands.
func (c *cli) runPhoneCallHistoryCallElement(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone call-history call-element commands: get")
		return nil
	}
	switch args[0] {
	case "get":
		return c.runPhoneCallHistoryCallElementGet(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "phone call-history call-element commands: get")
		return nil
	default:
		return fmt.Errorf("unknown phone call-history call-element command %q", args[0])
	}
}

// runPhoneCallHistoryCallElementGet fetches a call element by ID.
func (c *cli) runPhoneCallHistoryCallElementGet(args []string) error {
	fs := flag.NewFlagSet("phone call-history call-element get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	callElementID := fs.String("call-element-id", "", "Call element ID (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *callElementID == "" {
		return errors.New("--call-element-id is required")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	callElement, _, err := zoomClient.Phone.CallHistory.GetCallElement(context.Background(), *callElementID)
	if err != nil {
		return err
	}

	return c.writeJSON(callElement)
}

// runPhoneCallHistoryAICallSummary dispatches ai-summary subcommands.
func (c *cli) runPhoneCallHistoryAICallSummary(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone call-history ai-summary commands: get")
		return nil
	}
	switch args[0] {
	case "get":
		return c.runPhoneCallHistoryAICallSummaryGet(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "phone call-history ai-summary commands: get")
		return nil
	default:
		return fmt.Errorf("unknown phone call-history ai-summary command %q", args[0])
	}
}

// runPhoneCallHistoryAICallSummaryGet fetches an AI call summary.
func (c *cli) runPhoneCallHistoryAICallSummaryGet(args []string) error {
	fs := flag.NewFlagSet("phone call-history ai-summary get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "User ID (required)")
	summaryID := fs.String("ai-call-summary-id", "", "AI call summary ID (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *userID == "" {
		return errors.New("--user-id is required")
	}
	if *summaryID == "" {
		return errors.New("--ai-call-summary-id is required")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	summary, _, err := zoomClient.Phone.CallHistory.GetAICallSummary(context.Background(), *userID, *summaryID)
	if err != nil {
		return err
	}

	return c.writeJSON(summary)
}

// runPhoneRecordings dispatches phone recordings commands.
func (c *cli) runPhoneRecordings(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone recordings commands: get, download-recording, download-transcript, delete, enable-auto-delete, disable-auto-delete, recover")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runPhoneRecordingsGet(args[1:])
	case "download-recording":
		return c.runPhoneRecordingsDownloadRecording(args[1:])
	case "download-transcript":
		return c.runPhoneRecordingsDownloadTranscript(args[1:])
	case "delete":
		return c.runPhoneRecordingsDelete(args[1:])
	case "enable-auto-delete":
		return c.runPhoneRecordingsEnableAutoDelete(args[1:])
	case "disable-auto-delete":
		return c.runPhoneRecordingsDisableAutoDelete(args[1:])
	case "recover":
		return c.runPhoneRecordingsRecover(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "phone recordings commands: get, download-recording, download-transcript, delete, enable-auto-delete, disable-auto-delete, recover")
		return nil
	default:
		return fmt.Errorf("unknown phone recordings command %q", args[0])
	}
}

// runPhoneRecordingsGet fetches account, user, or call-targeted recordings.
func (c *cli) runPhoneRecordingsGet(args []string) error {
	fs := flag.NewFlagSet("phone recordings get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "User ID")
	callID := fs.String("call-id", "", "Call ID")
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional get query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional get query parameters")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *userID != "" && *callID != "" {
		return errors.New("--user-id and --call-id cannot be used together")
	}
	if (*userID != "" || *callID != "") && (*queryJSON != "" || *queryJSONFile != "") {
		return errors.New("query parameter input can only be used for account-level recordings requests")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	opts := make([]client.CallRecordingGetOptions, 0, 1)
	if *userID != "" {
		opts = append(opts, client.WithRecordingUserId(*userID))
	} else if *callID != "" {
		opts = append(opts, client.WithRecordingCallId(*callID))
	} else {
		queryParameters, queryErr := readOptionalJSONInputWithFlags[client.CallRecordingQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
		if queryErr != nil {
			return queryErr
		}
		if queryParameters != nil {
			opts = append(opts, client.WithCallRecordingQueryParameters(queryParameters))
		}
	}

	recordings, _, err := zoomClient.Phone.Recordings.Get(context.Background(), opts...)
	if err != nil {
		return err
	}

	return c.writeJSON(recordings)
}

// runPhoneRecordingsDownloadRecording downloads a recording file by file ID.
func (c *cli) runPhoneRecordingsDownloadRecording(args []string) error {
	fs := flag.NewFlagSet("phone recordings download-recording", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	fileID := fs.String("file-id", "", "Recording file ID")
	outputPath := fs.String("o", "", "Output file path (required)")
	fs.StringVar(outputPath, "output", "", "Output file path (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}

	if *fileID == "" {
		return errors.New("--file-id is required")
	}
	if *outputPath == "" {
		return errors.New("-o/--output is required")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	outputFile, err := os.Create(*outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	_, err = zoomClient.Phone.Recordings.DownloadCallRecording(context.Background(), *fileID, outputFile)
	closeErr := outputFile.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return fmt.Errorf("failed to close output file: %w", closeErr)
	}

	_, _ = fmt.Fprintf(c.stdout, "recording written to %s\n", *outputPath)
	return nil
}

// runPhoneRecordingsDownloadTranscript downloads transcript JSON by recording ID.
func (c *cli) runPhoneRecordingsDownloadTranscript(args []string) error {
	fs := flag.NewFlagSet("phone recordings download-transcript", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	recordingID := fs.String("recording-id", "", "Recording ID")
	outputPath := fs.String("o", "", "Output file path (required)")
	fs.StringVar(outputPath, "output", "", "Output file path (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}

	if *recordingID == "" {
		return errors.New("--recording-id is required")
	}
	if *outputPath == "" {
		return errors.New("-o/--output is required")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	transcript, _, err := zoomClient.Phone.Recordings.DownloadCallTranscript(context.Background(), *recordingID)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(transcript, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal transcript output: %w", err)
	}
	if err := os.WriteFile(*outputPath, b, 0o644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	_, _ = fmt.Fprintf(c.stdout, "transcript written to %s\n", *outputPath)
	return nil
}

// runPhoneSettings dispatches phone settings commands.
func (c *cli) runPhoneSettings(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone settings commands: get, update")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runPhoneSettingsGet(args[1:])
	case "update":
		return c.runPhoneSettingsUpdate(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "phone settings commands: get, update")
		return nil
	default:
		return fmt.Errorf("unknown phone settings command %q", args[0])
	}
}

// runPhoneSettingsGet fetches account-level phone settings.
func (c *cli) runPhoneSettingsGet(args []string) error {
	fs := flag.NewFlagSet("phone settings get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
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

// phoneSettingsUpdateInput mirrors client.SettingsAttributes but with JSON
// tags so it can be supplied via --json or --json-file. Each field is a
// pointer so that omitted fields are not sent in the PATCH request.
type phoneSettingsUpdateInput struct {
	BillingAccountId       *string `json:"billing_account_id,omitempty"`
	BYOC                   *bool   `json:"byoc,omitempty"`
	MultipleSites          *bool   `json:"multiple_sites,omitempty"`
	SiteCode               *bool   `json:"site_code,omitempty"`
	ShortExtensionLength   *int    `json:"short_extension_length,omitempty"`
	ShowDeviceIPForCallLog *bool   `json:"show_device_ip_for_call_log,omitempty"`
}

// runPhoneSettingsUpdate patches account-level phone settings.
func (c *cli) runPhoneSettingsUpdate(args []string) error {
	fs := flag.NewFlagSet("phone settings update", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	jsonInput := fs.String("json", "", "Inline JSON object for phone settings attributes")
	jsonFile := fs.String("json-file", "", "Path to JSON file for phone settings attributes")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}

	input, err := readJSONInput[phoneSettingsUpdateInput](*jsonInput, *jsonFile)
	if err != nil {
		return err
	}

	attrs := &client.SettingsAttributes{
		BillingAccountId:       input.BillingAccountId,
		BYOC:                   input.BYOC,
		MultipleSites:          input.MultipleSites,
		SiteCode:               input.SiteCode,
		ShortExtensionLength:   input.ShortExtensionLength,
		ShowDeviceIPForCallLog: input.ShowDeviceIPForCallLog,
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	if _, err := zoomClient.Phone.Settings.Update(context.Background(), attrs); err != nil {
		return err
	}

	_, _ = fmt.Fprintln(c.stdout, "phone settings updated")
	return nil
}

// runPhoneUsers dispatches phone users commands.
func (c *cli) runPhoneUsers(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone users commands: get, profile-settings get")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runPhoneUsersGet(args[1:])
	case "profile-settings":
		return c.runPhoneUsersProfileSettings(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "phone users commands: get, profile-settings get")
		return nil
	default:
		return fmt.Errorf("unknown phone users command %q", args[0])
	}
}

// runPhoneUsersGet fetches one phone user or lists phone users.
func (c *cli) runPhoneUsersGet(args []string) error {
	fs := flag.NewFlagSet("phone users get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "Phone user ID")
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional get query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional get query parameters")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}

	if *userID != "" && (*queryJSON != "" || *queryJSONFile != "") {
		return errors.New("--user-id cannot be used with query parameter input")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	opts := make([]client.PhoneUserGetOptions, 0, 1)
	if *userID != "" {
		opts = append(opts, client.WithPhoneUserID(*userID))
	} else {
		queryParameters, queryErr := readOptionalJSONInputWithFlags[client.PhoneUserQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
		if queryErr != nil {
			return queryErr
		}
		if queryParameters != nil {
			opts = append(opts, client.WithPhoneUserQueryParameters(queryParameters))
		}
	}

	users, _, err := zoomClient.Phone.Users.Get(context.Background(), opts...)
	if err != nil {
		return err
	}

	return c.writeJSON(users)
}

// runPhoneUsersProfileSettings dispatches profile-settings commands.
func (c *cli) runPhoneUsersProfileSettings(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "phone users profile-settings commands: get")
		return nil
	}
	switch args[0] {
	case "get":
		return c.runPhoneUsersProfileSettingsGet(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "phone users profile-settings commands: get")
		return nil
	default:
		return fmt.Errorf("unknown phone users profile-settings command %q", args[0])
	}
}

// runPhoneUsersProfileSettingsGet fetches phone profile settings for a user.
func (c *cli) runPhoneUsersProfileSettingsGet(args []string) error {
	fs := flag.NewFlagSet("phone users profile-settings get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "Phone user ID (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *userID == "" {
		return errors.New("--user-id is required")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	settings, _, err := zoomClient.Phone.Users.GetProfileSetting(context.Background(), *userID)
	if err != nil {
		return err
	}

	return c.writeJSON(settings)
}

// runPhoneRecordingsDelete permanently removes a recording.
func (c *cli) runPhoneRecordingsDelete(args []string) error {
	return c.runPhoneRecordingSimpleCmd(args, "phone recordings delete", "recording-id", "Recording ID (required)",
		func(zc *client.Client, recordingID string) (*http.Response, error) {
			return zc.Phone.Recordings.Delete(context.Background(), recordingID)
		}, "recording %s deleted")
}

// runPhoneRecordingsEnableAutoDelete enables auto-delete on a recording.
func (c *cli) runPhoneRecordingsEnableAutoDelete(args []string) error {
	return c.runPhoneRecordingSimpleCmd(args, "phone recordings enable-auto-delete", "recording-id", "Recording ID (required)",
		func(zc *client.Client, recordingID string) (*http.Response, error) {
			return zc.Phone.Recordings.EnableAutoDelete(context.Background(), recordingID)
		}, "auto-delete enabled on recording %s")
}

// runPhoneRecordingsDisableAutoDelete disables auto-delete on a recording.
func (c *cli) runPhoneRecordingsDisableAutoDelete(args []string) error {
	return c.runPhoneRecordingSimpleCmd(args, "phone recordings disable-auto-delete", "recording-id", "Recording ID (required)",
		func(zc *client.Client, recordingID string) (*http.Response, error) {
			return zc.Phone.Recordings.DisableAutoDelete(context.Background(), recordingID)
		}, "auto-delete disabled on recording %s")
}

// runPhoneRecordingsRecover restores a recording from the trash.
func (c *cli) runPhoneRecordingsRecover(args []string) error {
	return c.runPhoneRecordingSimpleCmd(args, "phone recordings recover", "recording-id", "Recording ID (required)",
		func(zc *client.Client, recordingID string) (*http.Response, error) {
			return zc.Phone.Recordings.Recover(context.Background(), recordingID)
		}, "recording %s recovered")
}

// runPhoneRecordingSimpleCmd is a shared helper for recording commands that
// take a single ID flag and return only a status message on success.
func (c *cli) runPhoneRecordingSimpleCmd(
	args []string,
	cmdName string,
	idFlag string,
	idFlagHelp string,
	invoke func(*client.Client, string) (*http.Response, error),
	successFormat string,
) error {
	fs := flag.NewFlagSet(cmdName, flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	id := fs.String(idFlag, "", idFlagHelp)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *id == "" {
		return fmt.Errorf("--%s is required", idFlag)
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}
	client.NewPhoneService(zoomClient)

	if _, err := invoke(zoomClient, *id); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, successFormat+"\n", *id)
	return nil
}
