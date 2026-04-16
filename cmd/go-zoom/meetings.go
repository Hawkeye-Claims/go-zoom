package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
)

// runMeetings dispatches meetings subcommands.
func (c *cli) runMeetings(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "meetings commands: get, create, update, delete, summary get|delete")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runMeetingsGet(args[1:])
	case "create":
		return c.runMeetingsCreate(args[1:])
	case "update":
		return c.runMeetingsUpdate(args[1:])
	case "delete":
		return c.runMeetingsDelete(args[1:])
	case "summary":
		return c.runMeetingsSummary(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "meetings commands: get, create, update, delete, summary get|delete")
		return nil
	default:
		return fmt.Errorf("unknown meetings command %q", args[0])
	}
}

// runMeetingsGet fetches one meeting or lists meetings for a user.
func (c *cli) runMeetingsGet(args []string) error {
	fs := flag.NewFlagSet("meetings get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	meetingID := fs.String("meeting-id", "", "Meeting ID")
	userID := fs.String("user-id", "", "User ID (required when listing meetings)")
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional get query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional get query parameters")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
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

	opts := make([]client.MeetingGetOptions, 0, 2)
	if *meetingID != "" {
		opts = append(opts, client.WithMeetingId(*meetingID))
		queryParameters, queryErr := readOptionalJSONInputWithFlags[client.MeetingQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
		if queryErr != nil {
			return queryErr
		}
		if queryParameters != nil {
			opts = append(opts, client.WithMeetingQueryParameters(queryParameters))
		}
	} else {
		opts = append(opts, client.WithMeetingUserId(*userID))
		queryParameters, queryErr := readOptionalJSONInputWithFlags[client.MeetingListQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
		if queryErr != nil {
			return queryErr
		}
		if queryParameters != nil {
			opts = append(opts, client.WithMeetingListQueryParameters(queryParameters))
		}
	}

	meetings, _, err := zoomClient.Meetings.Get(context.Background(), opts...)
	if err != nil {
		return err
	}

	return c.writeJSON(meetings)
}

// runMeetingsCreate creates a meeting for a user from JSON payload data.
func (c *cli) runMeetingsCreate(args []string) error {
	fs := flag.NewFlagSet("meetings create", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "User ID for whom to create the meeting")
	jsonInput := fs.String("json", "", "Inline JSON object for meeting attributes")
	jsonFile := fs.String("json-file", "", "Path to JSON file for meeting attributes")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *userID == "" {
		return errors.New("--user-id is required")
	}

	meetingAttributes, err := readJSONInput[client.MeetingAttributes](*jsonInput, *jsonFile)
	if err != nil {
		return err
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	meeting, _, err := zoomClient.Meetings.Create(context.Background(), *userID, meetingAttributes)
	if err != nil {
		return err
	}

	return c.writeJSON(meeting)
}

// runMeetingsSummary dispatches meetings summary commands.
func (c *cli) runMeetingsSummary(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "meetings summary commands: get, delete")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runMeetingsSummaryGet(args[1:])
	case "delete":
		return c.runMeetingsSummaryDelete(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "meetings summary commands: get, delete")
		return nil
	default:
		return fmt.Errorf("unknown meetings summary command %q", args[0])
	}
}

// runMeetingsSummaryGet fetches a specific summary or lists all summaries.
func (c *cli) runMeetingsSummaryGet(args []string) error {
	fs := flag.NewFlagSet("meetings summary get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	meetingID := fs.String("meeting-id", "", "Meeting ID (optional)")
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional get query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional get query parameters")
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

	opts := make([]client.MeetingsSummaryGetOptions, 0, 2)
	if *meetingID != "" {
		opts = append(opts, client.WithMeetingIdForSummary(*meetingID))
	}

	queryParameters, err := readOptionalJSONInputWithFlags[client.MeetingSummaryQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
	if err != nil {
		return err
	}
	if queryParameters != nil {
		opts = append(opts, client.WithMeetingSummaryQueryParameters(queryParameters))
	}

	summaries, _, err := zoomClient.Meetings.GetSummary(context.Background(), opts...)
	if err != nil {
		return err
	}

	return c.writeJSON(summaries)
}

// runMeetingsUpdate patches an existing meeting from JSON payload data.
func (c *cli) runMeetingsUpdate(args []string) error {
	fs := flag.NewFlagSet("meetings update", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	meetingID := fs.String("meeting-id", "", "Meeting ID (numeric, required)")
	jsonInput := fs.String("json", "", "Inline JSON object for meeting update attributes")
	jsonFile := fs.String("json-file", "", "Path to JSON file for meeting update attributes")
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional update query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional update query parameters")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *meetingID == "" {
		return errors.New("--meeting-id is required")
	}

	meetingIDInt, err := strconv.Atoi(*meetingID)
	if err != nil {
		return fmt.Errorf("--meeting-id must be a numeric meeting ID: %w", err)
	}

	meetingAttributes, err := readJSONInput[client.MeetingUpdateAttributes](*jsonInput, *jsonFile)
	if err != nil {
		return err
	}

	queryParameters, err := readOptionalJSONInputWithFlags[client.MeetingUpdateQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
	if err != nil {
		return err
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	opts := make([]client.MeetingUpdateOptions, 0, 1)
	if queryParameters != nil {
		opts = append(opts, client.WithMeetingUpdateQueryParameters(queryParameters))
	}

	if _, err := zoomClient.Meetings.Update(context.Background(), meetingIDInt, &meetingAttributes, opts...); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "meeting %d updated\n", meetingIDInt)
	return nil
}

// runMeetingsDelete cancels/removes a meeting with optional delete query parameters.
func (c *cli) runMeetingsDelete(args []string) error {
	fs := flag.NewFlagSet("meetings delete", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	meetingID := fs.String("meeting-id", "", "Meeting ID (numeric, required)")
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional delete query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional delete query parameters")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *meetingID == "" {
		return errors.New("--meeting-id is required")
	}

	meetingIDInt, err := strconv.Atoi(*meetingID)
	if err != nil {
		return fmt.Errorf("--meeting-id must be a numeric meeting ID: %w", err)
	}

	queryParameters, err := readOptionalJSONInputWithFlags[client.MeetingDeleteQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
	if err != nil {
		return err
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	opts := make([]client.MeetingDeleteOptions, 0, 1)
	if queryParameters != nil {
		opts = append(opts, client.WithMeetingDeleteQueryParameters(queryParameters))
	}

	if _, err := zoomClient.Meetings.Delete(context.Background(), meetingIDInt, opts...); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "meeting %d deleted\n", meetingIDInt)
	return nil
}

// runMeetingsSummaryDelete removes the AI summary for a meeting.
func (c *cli) runMeetingsSummaryDelete(args []string) error {
	fs := flag.NewFlagSet("meetings summary delete", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	meetingID := fs.String("meeting-id", "", "Meeting ID (required)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *meetingID == "" {
		return errors.New("--meeting-id is required")
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	if _, err := zoomClient.Meetings.DeleteSummary(context.Background(), *meetingID); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "meeting %s summary deleted\n", *meetingID)
	return nil
}
