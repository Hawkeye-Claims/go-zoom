# go-zoom

A lightweight [Zoom API](https://marketplace.zoom.us/docs/api-reference/introduction/) client for Go. Endpoint coverage is minimal, but the existing services are designed to be easy to extend.

## Installation

```sh
go get github.com/TheSlowpes/go-zoom
```

## Client Setup

### Server-to-Server OAuth (default)

Use this for backend services with a [Server-to-Server OAuth app](https://marketplace.zoom.us/docs/guides/build/server-to-server-oauth-app/). This is the default grant type — no extra options are needed.

```go
import (
    "net/http"
    "os"

    "github.com/TheSlowpes/go-zoom/zoom/client"
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
import (
    "net/http"
    "os"

    "github.com/TheSlowpes/go-zoom/zoom/client"
)

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

## Webhook Listener

Use `server.NewWebhookServer` to receive and process Zoom webhook events. All incoming requests are verified with HMAC-SHA256 against your webhook secret token before any handler is invoked. The `endpoint.url_validation` handshake required by Zoom is handled automatically.

```go
import (
    "fmt"
    "os"

    "github.com/TheSlowpes/go-zoom/zoom/server"
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

Two event payload structs are provided out of the box:

| Type | Backing model | Fields |
|------|--------------|--------|
| `server.MeetingEvent` | `models.Meeting` | `AccountId`, `Object`, `Operator`, `OperatorId`, `Operation` |
| `server.UserEvent` | `models.User` | `AccountId`, `Object`, `Operator`, `OperatorId`, `CreationType` |

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

## Making API Requests

```go
import (
    "context"
    "fmt"
    "log"

    "github.com/TheSlowpes/go-zoom/zoom/client"
)

ctx := context.Background()

users, _, err := c.Users.Get(ctx, client.WithUserId("me"))
if err != nil {
    log.Fatal(err)
}

for _, u := range users {
    fmt.Printf("ID: %s  Name: %s  Email: %s\n", u.ID, u.DisplayName, u.Email)
}
```

Paginated endpoints are handled transparently — the client follows `next_page_token` automatically and returns the complete result set.
