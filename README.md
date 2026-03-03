# go-zoom

A lightweight [Zoom API](https://marketplace.zoom.us/docs/api-reference/introduction/) client for Go. Endpoint coverage is minimal, but the existing services are designed to be easy to extend.

## Installation

```sh
go get github.com/Hawkeye-Claims/go-zoom
```

## Client Setup

### Server-to-Server OAuth (default)

Use this for backend services with a [Server-to-Server OAuth app](https://marketplace.zoom.us/docs/guides/build/server-to-server-oauth-app/). This is the default grant type — no extra options are needed.

```go
import (
    "net/http"
    "os"

    "github.com/Hawkeye-Claims/go-zoom/zoom/client"
)

c, err := client.NewClient(
    &http.Client{},
    os.Getenv("ZOOM_ACCOUNT_ID"),
    os.Getenv("ZOOM_CLIENT_ID"),
    os.Getenv("ZOOM_CLIENT_SECRET"),
)
```

Tokens are fetched and cached automatically. They are refreshed on expiry or when a `401` response is received.

### Authorization Code OAuth

Use this for user-facing apps where end users authorize your application via the Zoom OAuth flow. You must supply a redirect URI and register the two HTTP middleware handlers on your server.

```go
c, err := client.NewClient(
    &http.Client{},
    os.Getenv("ZOOM_ACCOUNT_ID"),
    os.Getenv("ZOOM_CLIENT_ID"),
    os.Getenv("ZOOM_CLIENT_SECRET"),
    client.WithGrantType("authorization_code"),
    client.WithRedirectURI("https://yourapp.example.com/oauth/callback"),
)
if err != nil {
    log.Fatal(err)
}

// Redirect users to Zoom's authorization page
http.Handle("/oauth/login", c.RequestAuthorization())

// Handle the callback from Zoom, exchange the code for a token
http.Handle("/oauth/callback", c.HandleOAuthCallback())
```

`RequestAuthorization` redirects the incoming request to Zoom's authorization URL. `HandleOAuthCallback` exchanges the authorization code for an access token and stores it. After the callback completes, the client is ready to make API requests on behalf of the user. Refresh tokens are rotated automatically.

### Custom Token Store

By default, tokens are stored in memory. For horizontally scaled deployments, implement the `TokenMutex` interface and inject it via `WithToken`:

```go
c, err := client.NewClient(
    &http.Client{},
    os.Getenv("ZOOM_ACCOUNT_ID"),
    os.Getenv("ZOOM_CLIENT_ID"),
    os.Getenv("ZOOM_CLIENT_SECRET"),
    client.WithToken(myRedisTokenStore),
)
```

The `TokenMutex` interface:

```go
type TokenMutex interface {
    Lock(context.Context) error
    Unlock(context.Context) error
    Get(context.Context) (string, error)
    GetRefreshToken(context.Context) (string, error)
    Set(context.Context, string, time.Time) error
    SetRefreshToken(context.Context, string) error
    Clear(context.Context) error
}
```

## Services

`NewClient` initializes `c.Users` and `c.Meetings` automatically. The `Phone` service tree must be initialized separately by calling `NewPhoneService`:

```go
c, err := client.NewClient(...)

// Initialize Phone sub-services (CallHistory, Recordings, Settings, Users)
client.NewPhoneService(c)
```

After calling `NewPhoneService`, the following are available: `c.Phone.CallHistory`, `c.Phone.Recordings`, `c.Phone.Settings`, and `c.Phone.Users`.

### Users — `c.Users`

```go
// Get a single user
users, _, err := c.Users.Get(ctx, client.WithUserId("me"))

// List all users (auto-paginated)
users, _, err := c.Users.Get(ctx, client.WithListUserQueryParameters(&client.ListUserQueryParameters{
    Status: enums.ActiveUser,
}))

// Create a user
user, _, err := c.Users.Create(ctx, enums.Create, client.UserAttributes{...})

// Update a user
_, err := c.Users.Update(ctx, "userId", &client.UserUpdateAttributes{...})

// Delete a user
_, err := c.Users.Delete(ctx, "userId")
```

### Meetings — `c.Meetings`

```go
// Get a single meeting
meetings, _, err := c.Meetings.Get(ctx, client.WithMeetingId("meetingId"))

// List meetings for a user (auto-paginated)
meetings, _, err := c.Meetings.Get(ctx, client.WithMeetingUserId("userId"))

// Create a meeting
meeting, _, err := c.Meetings.Create(ctx, "userId", client.MeetingAttributes{...})

// Update a meeting
_, err := c.Meetings.Update(ctx, meetingId, &client.MeetingUpdateAttributes{...})

// Delete a meeting
_, err := c.Meetings.Delete(ctx, meetingId)
```

#### Meeting Summaries

```go
// Get AI-generated summaries for a meeting
summaries, _, err := c.Meetings.GetSummary(ctx, client.WithMeetingIdForSummary("meetingId"))

// Delete a meeting summary
_, err := c.Meetings.DeleteSummary(ctx, "meetingId")
```

### Phone — `c.Phone`

#### Call History — `c.Phone.CallHistory`

```go
// Get account-wide call history (auto-paginated)
history, _, err := c.Phone.CallHistory.Get(ctx)

// Get call history for a specific user (auto-paginated)
history, _, err := c.Phone.CallHistory.Get(ctx,
    client.WithUserIdForPhoneCallHistory("userId"),
    client.WithPhoneCallHistoryQueryParameters(&client.PhoneCallHistoryQueryParameters{
        From: "2024-01-01",
        To:   "2024-01-31",
    }),
)

// Get a single call history record by UUID
history, _, err := c.Phone.CallHistory.Get(ctx, client.WithPhoneCallHistoryUUID("uuid"))

// Get a single call element
element, _, err := c.Phone.CallHistory.GetCallElement("callElementId")

// Get an AI call summary
summary, _, err := c.Phone.CallHistory.GetAICallSummary("userId", "aiCallSummaryId")

// Add a client code to a call log entry
_, err := c.Phone.CallHistory.AddClientCode("callLogId", "clientCode")

// Delete a user's call history entry
_, err := c.Phone.CallHistory.DeleteUserCallHistory("userId", "callLogId")
```

#### Recordings — `c.Phone.Recordings`

```go
// List call recordings for a user (auto-paginated)
recordings, _, err := c.Phone.Recordings.Get(ctx, client.WithRecordingUserId("userId"))

// Get recordings for a specific call
recordings, _, err := c.Phone.Recordings.Get(ctx, client.WithRecordingCallId("callId"))

// Download a recording to an io.Writer
_, err := c.Phone.Recordings.DownloadCallRecording(ctx, "fileId", w)

// Download a transcript
transcript, _, err := c.Phone.Recordings.DownloadCallTranscript(ctx, "recordingId")

// Enable or disable auto-delete for a recording
_, err := c.Phone.Recordings.EnableAutoDelete("recordingId")
_, err := c.Phone.Recordings.DisableAutoDelete("recordingId")

// Recover a deleted recording
_, err := c.Phone.Recordings.Recover(ctx, "recordingId")

// Delete a recording
_, err := c.Phone.Recordings.Delete(ctx, "recordingId")
```

#### Settings — `c.Phone.Settings`

```go
// Get account-level phone settings
settings, _, err := c.Phone.Settings.Get(ctx)

// Update account-level phone settings
_, err := c.Phone.Settings.Update(ctx, &client.SettingsAttributes{...})
```

#### Phone Users — `c.Phone.Users`

```go
// List all phone users (auto-paginated)
users, _, err := c.Phone.Users.Get(ctx)

// Filter phone users
users, _, err := c.Phone.Users.Get(ctx,
    client.WithPhoneUserQueryParameters(&client.PhoneUserQueryParameters{
        Status: enums.ActiveUser,
    }),
)

// Get a single phone user by ID
users, _, err := c.Phone.Users.Get(ctx, client.WithPhoneUserID("userId"))

// Get a phone user's profile settings
settings, _, err := c.Phone.Users.GetProfileSetting(ctx, "userId")
```

## Webhook Listener

Use `server.NewWebhookServer` to receive and process Zoom webhook events. All incoming requests are verified with HMAC-SHA256 against your webhook secret token before any handler is invoked. The `endpoint.url_validation` handshake required by Zoom is handled automatically.

```go
import (
    "fmt"
    "os"

    "github.com/Hawkeye-Claims/go-zoom/zoom/server"
)

meetingCh := make(chan server.MeetingEvent, 10)
userCh    := make(chan server.UserEvent, 10)

ws := server.NewWebhookServer(
    ":8080",
    "/zoom/webhook",
    os.Getenv("ZOOM_WEBHOOK_SECRET"),
    server.WithHandler("meeting.created", meetingCh),
    server.WithHandler("user.created",    userCh),
)

go ws.Start()

for evt := range meetingCh {
    fmt.Println("New meeting:", evt.Object.Topic)
}
```

`WithHandler` uses Go generics to route each event type to a typed channel. The payload is unmarshalled directly into the channel's element type — no type assertions needed.

### Built-in Payload Types

Several event payload structs are provided out of the box:

| Type | Backing model | Fields |
|------|--------------|--------|
| `server.MeetingEvent` | `models.Meeting` | `AccountId`, `Object`, `Operator`, `OperatorId`, `Operation` |
| `server.UserEvent` | `models.User` | `AccountId`, `Object`, `Operator`, `OperatorId`, `CreationType` |
| `server.AICallSummaryEvent` | `models.AICallSummary` | `AccountId`, `Object` |
| `server.PhoneCallElementEvent` | `[]models.CallElement` | `AccountId`, `Object.CallElements`, `UserID` |
| `server.PhoneCallHistoryEvent` | `[]models.CallHistory` | `AccountId`, `Object.CallLogs`, `UserID` |

### Custom Payload Types

You are not limited to the built-in types. Any struct can be used as the payload generic for `WithHandler`. Define your own struct matching the `payload` field of the Zoom event you want to handle and pass a channel of that type:

```go
type MyCustomPayload struct {
    AccountId string `json:"account_id"`
    // ... fields matching the Zoom event payload
}

ch := make(chan MyCustomPayload, 10)
server.WithHandler("recording.completed", ch)
```

The full event envelope is defined as `server.Notification[T]` and is available if you need access to the top-level `Event` string or `EventTs` timestamp:

```go
type Notification[T any] struct {
    Event   string `json:"event"`
    EventTs int64  `json:"event_ts"`
    Payload T      `json:"payload"`
}
```

## Pagination

Paginated endpoints follow `next_page_token` automatically and return the complete result set. No additional handling is required by the caller.
