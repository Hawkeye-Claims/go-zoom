package main

import "fmt"

// run dispatches top-level commands.
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

// printRootUsage prints the root CLI usage text.
func (c *cli) printRootUsage() {
	fmt.Fprintln(c.stderr, `go-zoom CLI wrapper

Usage:
  go run ./cmd/go-zoom <command> [subcommand] [flags]

Commands:
  auth                                  Authentication helpers (SDK auth options)
  users get|create|update|delete        Manage Zoom users
  meetings get|create|update|delete     Manage Zoom meetings
  meetings summary get|delete           Get/list or delete AI meeting summaries
  phone call-history get                Get account/user/single call history
  phone call-history add-client-code    Tag a call log with a client code
  phone call-history delete             Delete a user's call log entry
  phone call-history call-element get   Fetch a call element by ID
  phone call-history ai-summary get     Fetch an AI call summary
  phone recordings get                  Get account/user/call recordings
  phone recordings download-recording   Download a recording file
  phone recordings download-transcript  Download transcript JSON
  phone recordings delete               Permanently delete a recording
  phone recordings enable-auto-delete   Enable auto-delete on a recording
  phone recordings disable-auto-delete  Disable auto-delete on a recording
  phone recordings recover              Recover a recording from the trash
  phone settings get|update             Get or patch account phone settings
  phone users get                       Get one phone user or list phone users
  phone users profile-settings get      Get a phone user's profile settings

Use "<command> --help" for command-specific flags.`)
}
