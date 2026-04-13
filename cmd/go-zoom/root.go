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
  auth                     Authentication helpers (SDK auth options)
  users get|create         Get/list users or create a user
  meetings get|create      Get/list meetings or create a meeting
  meetings summary get     Get one meeting summary or list summaries
  phone call-history get   Get account/user/single call history
  phone recordings get     Get account/user/call recordings
  phone recordings download-recording   Download a recording file
  phone recordings download-transcript  Download transcript JSON
  phone settings get       Get account phone settings
  phone users get          Get one phone user or list phone users

Use "<command> --help" for command-specific flags.`)
}
