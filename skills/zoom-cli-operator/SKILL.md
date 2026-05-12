---
name: zoom-cli-operator
description: Use this skill when the user wants to operate Zoom through the go-zoom CLI. This includes checking upcoming or existing meetings, creating/updating/deleting meetings, listing AI meeting summaries, checking phone call history, downloading call recordings or transcripts, checking phone users/settings, validating credentials, or asking for the right `go run ./cmd/go-zoom ...` command. Use it whenever the user asks to look up Zoom data, create a Zoom meeting, inspect transcripts/summaries/recordings/call logs, or troubleshoot CLI usage in this repo.
---

# Using The go-zoom CLI

This repository includes a lightweight Zoom CLI at `cmd/go-zoom`. It may be used either as a built binary or directly from the source tree.

Prefer the built binary when it is available in `PATH` or the user says they have installed/built it:

```sh
go-zoom <command> [subcommand] [flags]
```

When working from the repository source tree, use:

```sh
go run ./cmd/go-zoom <command> [subcommand] [flags]
```

Prefer this CLI when the user wants to operate Zoom directly. Do not suggest Zoom's separate global `zoom` command for these workflows; this skill is about this project's `go-zoom` CLI.

## Authentication

Service commands use Server-to-Server OAuth (`account_credentials`). Credentials can come from environment variables or flags:

```sh
export ZOOM_ACCOUNT_ID=...
export ZOOM_CLIENT_ID=...
export ZOOM_CLIENT_SECRET=...
```

Validate credentials before doing real work when auth status is unknown:

```sh
go run ./cmd/go-zoom auth test --grant-type account_credentials
```

The CLI has an `auth authorization-code` helper for user OAuth. Non-auth service commands still reject `--grant-type authorization_code`; use account credentials for `users`, `meetings`, and `phone` commands unless the CLI has been extended to persist and reuse user OAuth tokens.

### User OAuth Testing

Zoom OAuth apps normally require a registered HTTPS redirect URI, so do not assume `localhost` redirect URIs will work. The CLI supports two practical test flows.

Manual paste-back flow, useful when you have a registered public HTTPS redirect URI but no callback server:

```sh
go run ./cmd/go-zoom auth authorization-code \
  --manual \
  --redirect-uri https://YOUR_REGISTERED_DOMAIN/oauth/callback
```

The command prints the SDK-built Zoom authorization URL. Open it, approve access, then paste the full redirected callback URL or the authorization code into the CLI. Prefer pasting the full callback URL because it includes the OAuth `state` value for validation.

Ngrok callback flow, useful for local testing when the Zoom app can be configured with the ngrok HTTPS URL:

```sh
ngrok http 8080
```

Copy the `https://...ngrok...` forwarding URL, add `/oauth/callback`, and register that exact URL in the Zoom app redirect URI allowlist. In another terminal, run:

```sh
go run ./cmd/go-zoom auth authorization-code \
  --redirect-uri https://YOUR_NGROK_HOST/oauth/callback \
  --listen-addr localhost:8080
```

Use `--listen-addr` whenever the redirect URI host is public HTTPS but the callback server should bind locally. Without it, the CLI tries to listen on the redirect URI host.

Never print secrets back to the user. If credentials are missing, tell the user which environment variables or flags are required without asking them to paste secrets into chat.

## Command Discovery

When unsure, ask the CLI for help rather than guessing:

```sh
go-zoom --help
go-zoom meetings --help
go-zoom phone --help
go-zoom phone recordings --help
```

If `go-zoom` is not available, run the same command through `go run ./cmd/go-zoom` from the repository root.

The CLI command tree currently includes:

- `auth test`
- `auth authorization-code`
- `users get|create|update|delete`
- `meetings get|create|update|delete`
- `meetings summary get|delete`
- `phone call-history get|add-client-code|delete|call-element get|ai-summary get`
- `phone recordings get|download-recording|download-transcript|delete|enable-auto-delete|disable-auto-delete|recover`
- `phone settings get|update`
- `phone users get|profile-settings get`

## JSON Inputs

Create/update commands require `--json` or `--json-file`. Filterable get commands use `--query-json` or `--query-json-file`.

The CLI decodes JSON strictly, so unknown fields fail. Prefer `--json-file` for anything long, sensitive, or shell-quote-heavy.

Example inline JSON:

```sh
go-zoom meetings create \
  --user-id me \
  --json '{"topic":"Project kickoff","type":2,"start_time":"2026-05-13T15:00:00Z","duration":30}'
```

Example query JSON:

```sh
go-zoom meetings get \
  --user-id me \
  --query-json '{"type":"upcoming"}'
```

## Meetings

Use `meetings get` to fetch one meeting or list meetings for a user. It requires exactly one of `--meeting-id` or `--user-id`.

```sh
# List meetings for a user, such as upcoming meetings when the API accepts the type filter.
go-zoom meetings get \
  --user-id me \
  --query-json '{"type":"upcoming"}'

# Fetch one meeting.
go-zoom meetings get --meeting-id 123456789
```

Create meetings with `meetings create`:

```sh
go-zoom meetings create \
  --user-id me \
  --json-file ./meeting.json
```

Useful meeting JSON fields include `topic`, `type`, `start_time`, `duration`, `timezone`, `agenda`, `password`, and `settings`. `start_time` should be an RFC3339 timestamp such as `2026-05-13T15:00:00Z`. Scheduled meetings commonly use `type: 2`.

Update and delete meetings:

```sh
go-zoom meetings update \
  --meeting-id 123456789 \
  --json '{"topic":"Updated topic"}'

go-zoom meetings delete --meeting-id 123456789
```

## Meeting Summaries

Use meeting summaries for Zoom AI-generated meeting summary data:

```sh
# List summaries, optionally filtered by date fields.
go-zoom meetings summary get

# Get summary for one meeting.
go-zoom meetings summary get --meeting-id 123456789

# Delete a summary.
go-zoom meetings summary delete --meeting-id 123456789
```

If the user asks for "transcripts" of meetings, check whether they mean AI meeting summaries or phone call recording transcripts. This CLI currently exposes meeting summaries and phone recording transcripts; it does not expose a separate meeting-recording transcript command.

## Phone Call History And AI Call Summaries

Use call history for Zoom Phone calls:

```sh
# Account-wide call history.
go-zoom phone call-history get

# A user's call history with date filters.
go-zoom phone call-history get \
  --user-id USER_ID \
  --query-json '{"from":"2026-05-01","to":"2026-05-12"}'

# One call-history record.
go-zoom phone call-history get --uuid CALL_HISTORY_UUID
```

Related commands:

```sh
go-zoom phone call-history call-element get --call-element-id CALL_ELEMENT_ID
go-zoom phone call-history ai-summary get --user-id USER_ID --ai-call-summary-id SUMMARY_ID
go-zoom phone call-history add-client-code --call-log-id CALL_LOG_ID --client-code CLIENT_CODE
go-zoom phone call-history delete --user-id USER_ID --call-log-id CALL_LOG_ID
```

## Phone Recordings And Transcripts

Use recordings commands for Zoom Phone call recordings:

```sh
# Account-level recordings with filters.
go-zoom phone recordings get \
  --query-json '{"from":"2026-05-01","to":"2026-05-12"}'

# Recordings for one user.
go-zoom phone recordings get --user-id USER_ID

# Recordings for one call.
go-zoom phone recordings get --call-id CALL_ID
```

Download outputs to files:

```sh
go-zoom phone recordings download-recording \
  --file-id FILE_ID \
  --output ./recording.mp3

go-zoom phone recordings download-transcript \
  --recording-id RECORDING_ID \
  --output ./transcript.json
```

Be careful with destructive or state-changing recording commands. Confirm user intent before deleting, recovering, or changing auto-delete:

```sh
go-zoom phone recordings delete --recording-id RECORDING_ID
go-zoom phone recordings enable-auto-delete --recording-id RECORDING_ID
go-zoom phone recordings disable-auto-delete --recording-id RECORDING_ID
go-zoom phone recordings recover --recording-id RECORDING_ID
```

## Phone Users And Settings

Phone users:

```sh
go-zoom phone users get
go-zoom phone users get --user-id USER_ID
go-zoom phone users get --query-json '{"status":"active"}'
go-zoom phone users profile-settings get --user-id USER_ID
```

Phone account settings:

```sh
go-zoom phone settings get
go-zoom phone settings update --json-file ./phone-settings.json
```

Treat `phone settings update` as a state-changing operation. Show the proposed JSON and ask for confirmation if the user has not clearly authorized the change.

## Users

Common user commands:

```sh
go-zoom users get --user-id me
go-zoom users get --query-json '{"status":"active"}'
go-zoom users create --action create --json-file ./user.json
go-zoom users update --user-id USER_ID --json-file ./user-update.json
go-zoom users delete --user-id USER_ID
```

Confirm before destructive user operations, especially `users delete`.

## Operating Style

When the user asks for information:

1. Check auth if needed.
2. Run the narrowest read-only CLI command that answers the question.
3. Summarize the relevant fields instead of dumping huge JSON unless the user asks for raw output.
4. If the command returns an error, report the exact command and the useful part of the error without exposing secrets.

When the user asks to create or modify data:

1. Build the smallest valid JSON payload.
2. Prefer writing JSON to a temporary or workspace file when quoting would be fragile.
3. For destructive operations, confirm first unless the user explicitly gave a direct instruction with the target ID.
4. Run the command and summarize the created/changed resource ID, URL, status, or output file path.

Do not claim the CLI supports endpoints that are not in the command tree. If the user asks for an unsupported operation, say it is not exposed by the CLI yet and offer to add SDK/CLI support.
