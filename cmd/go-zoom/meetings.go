package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
)

// runMeetings dispatches meetings subcommands.
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

// runMeetingsGet fetches one meeting or lists meetings for a user.
func (c *cli) runMeetingsGet(args []string) error {
	fs := flag.NewFlagSet("meetings get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	meetingID := fs.String("meeting-id", "", "Meeting ID")
	userID := fs.String("user-id", "", "User ID (required when listing meetings)")
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

// runMeetingsSummary dispatches meetings summary commands.
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

// runMeetingsSummaryGet fetches a specific summary or lists all summaries.
func (c *cli) runMeetingsSummaryGet(args []string) error {
	fs := flag.NewFlagSet("meetings summary get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	meetingID := fs.String("meeting-id", "", "Meeting ID (optional)")
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
