package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
	"github.com/Hawkeye-Claims/go-zoom/zoom/enums"
)

// runUsers dispatches users subcommands.
func (c *cli) runUsers(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "users commands: get, create")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runUsersGet(args[1:])
	case "create":
		return c.runUsersCreate(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "users commands: get, create")
		return nil
	default:
		return fmt.Errorf("unknown users command %q", args[0])
	}
}

// runUsersGet fetches one user or lists users when --user-id is omitted.
func (c *cli) runUsersGet(args []string) error {
	fs := flag.NewFlagSet("users get", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "User ID or email (optional; empty lists users)")
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

// runUsersCreate creates a user from JSON payload data.
func (c *cli) runUsersCreate(args []string) error {
	fs := flag.NewFlagSet("users create", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	action := fs.String("action", string(enums.Create), "User create action: create, auto_create, cust_create, sso_create")
	jsonInput := fs.String("json", "", "Inline JSON object for user attributes")
	jsonFile := fs.String("json-file", "", "Path to JSON file for user attributes")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}

	userAction := enums.UserCreateAction(*action)
	switch userAction {
	case enums.Create, enums.AutoCreate, enums.CustCreate, enums.SSOCreate:
	default:
		return fmt.Errorf("unsupported action %q", *action)
	}

	userAttributes, err := readJSONInput[client.UserAttributes](*jsonInput, *jsonFile)
	if err != nil {
		return err
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	user, _, err := zoomClient.Users.Create(context.Background(), userAction, userAttributes)
	if err != nil {
		return err
	}

	return c.writeJSON(user)
}
