---
name: zoom-endpoint-sdk
description: Use this skill when adding, changing, or reviewing Zoom REST API endpoints in the go-zoom SDK. This includes requests to add new Zoom client methods, model Zoom API request or response bodies, wire a new service into Client, update the CLI for an endpoint, handle pagination, add query options, or follow this repository's endpoint conventions. Use it whenever the user mentions Zoom endpoints, SDK coverage, API resources, service methods, models, or CLI commands in this repo, even if they do not explicitly ask for a skill.
---

# Adding Zoom Endpoints To go-zoom

This repository wraps Zoom REST API endpoints in a small Go SDK. Preserve the existing service style so new endpoints feel native to the project instead of like a generic generated client.

## First Pass

1. Identify the Zoom API resource, HTTP method, path, path parameters, query parameters, request body, response body, success status code, and whether the endpoint is paginated.
2. Inspect the closest existing implementation before editing. Good starting points are `zoom/client/user.go`, `zoom/client/meeting.go`, `zoom/client/meeting-summaries.go`, `zoom/client/phone.go`, `zoom/client/phone-recordings.go`, `zoom/client/phone-users.go`, and `zoom/client/phone-settings.go`.
3. Decide whether the endpoint belongs to an existing service (`Users`, `Meetings`, `Phone.CallHistory`, `Phone.Recordings`, `Phone.Settings`, `Phone.Users`) or needs a new service/sub-service.
4. Add only the endpoint surface the user requested. Do not generate broad API coverage unless asked.

## Client Service Pattern

Use the existing internal request path:

```go
res, err := service.client.request(ctx, http.MethodGet, endpoint, query, body, out)
```

Keep these conventions:

- Service methods accept `context.Context` first.
- Methods return the decoded resource, the raw `*http.Response`, and `error` for response-producing calls.
- Methods that do not decode a body return `(*http.Response, error)`.
- Wrap request errors with `fmt.Errorf("Error making request: %w", err)`.
- Check expected success status codes after `request` when the endpoint has a specific success code. Existing create methods expect `http.StatusCreated`; update/delete/action methods usually expect `http.StatusNoContent`; ordinary reads usually expect `http.StatusOK` when they are single-resource helpers.
- Escape every path parameter with `url.PathEscape`, including IDs, emails, UUIDs, file IDs, and meeting IDs represented as strings.
- Use `fmt.Sprintf` for paths with parameters and literal strings for fixed paths.

## Service Interfaces

Each service exposes a small interface near the service type:

```go
type ExampleServicer interface {
    Get(ctx context.Context, opts ...ExampleGetOptions) ([]*models.Example, *http.Response, error)
}

type ExampleService struct {
    client *Client
}

var _ ExampleServicer = (*ExampleService)(nil)
```

When adding a method to an existing service, update its interface unless the file already documents that a method is intentionally excluded. `PhoneSettingsService.Update` is currently an explicit exception.

When adding a new `Phone` sub-service, update `PhoneService` in `zoom/client/phone.go`, initialize it in `NewPhoneService`, and add a field comment matching the existing style.

## Options And Query Parameters

Use functional options when a method has alternate modes or optional query parameters:

```go
type ExampleGetOptions func(*exampleGetOptions)

type exampleGetOptions struct {
    exampleID       string
    queryParameters *ExampleQueryParameters
}

func WithExampleID(exampleID string) ExampleGetOptions {
    return func(o *exampleGetOptions) {
        o.exampleID = exampleID
    }
}
```

Use `url` tags for query structs because `Client.request` encodes query values with `go-querystring`:

```go
type ExampleQueryParameters struct {
    From string `url:"from,omitempty"`
    To   string `url:"to,omitempty"`
}
```

Validate mutually exclusive options before making the request. Return `nil, nil, error` for methods that normally return a decoded value plus response when no request has been made.

## Pagination

List endpoints should usually auto-paginate when Zoom returns `next_page_token`. Embed `*PaginationResponse` in the local response wrapper and append all pages before returning.

Keep any original filter query parameters on subsequent pages by combining them with `*PaginationOptions`:

```go
type examplePageQuery struct {
    *ExampleQueryParameters
    *PaginationOptions
}

nextPageToken := queryResponse.NextPageToken
pageQuery := &examplePageQuery{
    ExampleQueryParameters: options.queryParameters,
    PaginationOptions: &PaginationOptions{NextPageToken: &nextPageToken},
}
```

Decode list response arrays using the exact Zoom response key, such as `users`, `meetings`, `summaries`, `call_history`, `call_logs`, or `call_recordings`.

## Models And Attributes

Put reusable Zoom response models in `zoom/models`. Put endpoint-specific request attribute structs in the client file unless they are reused across services. Existing examples include `UserAttributes`, `MeetingAttributes`, `MeetingUpdateAttributes`, and `SettingsAttributes`.

Use these struct tag conventions:

- JSON response models usually omit `omitempty` unless existing neighboring fields use it for optional values.
- Request/update attributes should use `json:"field,omitempty"` for optional value fields.
- Use pointer fields for PATCH attributes when false, zero, or empty string must be distinguishable from omitted. `SettingsAttributes` is the clearest local example.
- Reuse existing enum types from `zoom/enums` when the value set is already represented. Add small enum types only when doing so improves type safety and matches nearby code.
- Add concise field comments for exported structs and fields. This repo documents exported API heavily, and missing comments make new endpoint code stand out.

## CLI Pattern

The CLI is optional unless the user asks for endpoint access through `cmd/go-zoom` or the existing feature clearly has CLI coverage nearby.

When adding CLI coverage:

- Follow the nested dispatcher style in `cmd/go-zoom/phone.go`, `users.go`, and `meetings.go`.
- Use `flag.NewFlagSet` with `flag.ContinueOnError` and `fs.SetOutput(c.stderr)`.
- Bind Zoom client flags through `c.bindClientFlags(fs, true, true)`.
- Validate required flags and mutually exclusive inputs before constructing the service client.
- For JSON request bodies or query parameters, reuse `readJSONInput`, `readOptionalJSONInputWithFlags`, and local input structs with JSON tags.
- Output decoded resources with `c.writeJSON`. For successful no-body operations, write a short human-readable success message.
- Initialize `Phone` endpoints with `client.NewPhoneService(zoomClient)` before accessing `zoomClient.Phone`.

## Common Traps

- `NewClient` currently initializes `Users` and `Meetings`, but not `Phone`; phone callers must call `NewPhoneService`.
- `Client.request` already handles authentication, token refresh, query encoding, JSON encoding, JSON decoding, and `io.Writer` downloads. Do not duplicate that logic.
- Avoid using raw string concatenation for unescaped path parameters.
- Do not silently drop pagination filters on later pages.
- Do not add compatibility shims or alternate names unless the user asks for a public API migration path.
- Do not invent tests that require real Zoom credentials. Prefer compile-time checks and request-shape tests only if the repo already has a test harness or the user asks for one.

## Verification

Before finishing endpoint work:

1. Run `gofmt` on changed Go files.
2. Run `go test ./...` from the repository root.
3. If CLI code changed, also run the relevant command with `--help` when feasible to catch dispatch/flag errors.
4. Summarize the service methods added, model/enum files touched, CLI commands added, and any Zoom spec assumptions that need user confirmation.
