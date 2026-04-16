package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/Hawkeye-Claims/go-zoom/zoom/client"
	"github.com/Hawkeye-Claims/go-zoom/zoom/enums"
)

// runUsers dispatches users subcommands.
func (c *cli) runUsers(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "users commands: get, create, update, delete")
		return nil
	}

	switch args[0] {
	case "get":
		return c.runUsersGet(args[1:])
	case "create":
		return c.runUsersCreate(args[1:])
	case "update":
		return c.runUsersUpdate(args[1:])
	case "delete":
		return c.runUsersDelete(args[1:])
	case "help", "--help", "-h":
		fmt.Fprintln(c.stderr, "users commands: get, create, update, delete")
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

	opts := make([]client.UserGetOptions, 0, 2)
	if *userID != "" {
		opts = append(opts, client.WithUserId(*userID))
		queryParameters, queryErr := readOptionalJSONInputWithFlags[client.UserQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
		if queryErr != nil {
			return queryErr
		}
		if queryParameters != nil {
			opts = append(opts, client.WithUserQueryParameters(queryParameters))
		}
	} else {
		queryParameters, queryErr := readOptionalJSONInputWithFlags[client.ListUserQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
		if queryErr != nil {
			return queryErr
		}
		if queryParameters != nil {
			opts = append(opts, client.WithListUserQueryParameters(queryParameters))
		}
	}

	users, _, err := zoomClient.Users.Get(context.Background(), opts...)
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

// runUsersUpdate patches an existing user from JSON payload data.
func (c *cli) runUsersUpdate(args []string) error {
	fs := flag.NewFlagSet("users update", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "User ID or email (required)")
	jsonInput := fs.String("json", "", "Inline JSON object for user update attributes")
	jsonFile := fs.String("json-file", "", "Path to JSON file for user update attributes")
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional patch query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional patch query parameters")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *userID == "" {
		return errors.New("--user-id is required")
	}

	userAttributes, err := readJSONInput[client.UserUpdateAttributes](*jsonInput, *jsonFile)
	if err != nil {
		return err
	}

	queryParameters, err := readOptionalJSONInputWithFlags[client.UserPatchQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
	if err != nil {
		return err
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	opts := make([]client.UserPatchOptions, 0, 1)
	if queryParameters != nil {
		opts = append(opts, client.WithUserPatchQueryParameters(queryParameters))
	}

	if _, err := zoomClient.Users.Update(context.Background(), *userID, &userAttributes, opts...); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "user %s updated\n", *userID)
	return nil
}

// runUsersDelete deletes a user with optional transfer/delete query parameters.
func (c *cli) runUsersDelete(args []string) error {
	fs := flag.NewFlagSet("users delete", flag.ContinueOnError)
	fs.SetOutput(c.stderr)

	cfg := c.bindClientFlags(fs, true, true)
	userID := fs.String("user-id", "", "User ID or email (required)")
	queryJSON := fs.String("query-json", "", "Inline JSON object for optional delete query parameters")
	queryJSONFile := fs.String("query-json-file", "", "Path to JSON file for optional delete query parameters")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := ensureNoUnexpectedArgs(fs.Args()); err != nil {
		return err
	}
	if *userID == "" {
		return errors.New("--user-id is required")
	}

	queryParameters, err := readOptionalJSONInputWithFlags[client.UserDeleteQueryParameters](*queryJSON, *queryJSONFile, "--query-json", "--query-json-file")
	if err != nil {
		return err
	}

	zoomClient, err := c.newServiceClient(cfg)
	if err != nil {
		return err
	}

	opts := make([]client.UserDeleteOptions, 0, 1)
	if queryParameters != nil {
		opts = append(opts, client.WithUserDeleteQueryParameters(queryParameters))
	}

	if _, err := zoomClient.Users.Delete(context.Background(), *userID, opts...); err != nil {
		return err
	}

	_, _ = fmt.Fprintf(c.stdout, "user %s deleted\n", *userID)
	return nil
}
