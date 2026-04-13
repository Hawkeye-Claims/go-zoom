package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
)

// runPhone dispatches phone subcommands.
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

// runPhoneCallHistory fetches account, user, or UUID-targeted call history.
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
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional get query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional get query parameters")
	if err := fs.Parse(args[1:]); err != nil {
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

// runPhoneRecordings fetches account, user, or call-targeted recordings.
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
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional get query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional get query parameters")
	if err := fs.Parse(args[1:]); err != nil {
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

// runPhoneSettings fetches account-level phone settings.
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

// runPhoneUsers fetches one phone user or lists phone users.
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
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional get query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional get query parameters")
	if err := fs.Parse(args[1:]); err != nil {
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
