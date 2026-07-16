# Examples

## Authentication

The recommended entry point is `client.New(teamSlug, opts...)`. It manages
Auth0 authentication, automatic token refresh, and telemetry.

### Machine-to-machine (Client Credentials)

The most common mode for backend services. The SDK fetches an access token from
`https://{tenantDomain}/oauth/token` using the OAuth2 Client Credentials grant and
refreshes it automatically before expiry. The tenant domain is passed to
`option.WithClientCredentials` right after the context.

```go
package main

import (
    "context"

    client "github.com/auth0/teams-go/client"
    option "github.com/auth0/teams-go/option"
)

func main() {
    ctx := context.Background()

    c, err := client.New(
        "acme", // team slug
        option.WithClientCredentials(
            ctx,
            "acme-auth.us.auth0.com", // Auth0 tenant domain
            "<client-id>",
            "<client-secret>",
        ),
    )
    if err != nil {
        panic(err)
    }

    _ = c
}
```

By default the token audience is derived as `https://{teamSlug}.teams.auth0.com/api`.

### Machine-to-machine with an explicit audience

Use this when the audience differs from the default derived value.

```go
c, err := client.New(
    "acme",
    option.WithClientCredentialsAndAudience(
        ctx,
        "acme-auth.us.auth0.com", // Auth0 tenant domain
        "<client-id>",
        "<client-secret>",
        "https://acme.teams.auth0.com/api",
    ),
)
```

### Custom token source

Supply your own `oauth2.TokenSource` when you manage tokens elsewhere (for example, a
shared in-process token cache).

```go
import "golang.org/x/oauth2"

var src oauth2.TokenSource = myTokenSource()

c, err := client.New(
    "acme",
    option.WithTokenSource(src),
)
```

### Static bearer token

Use a pre-obtained access token directly. The SDK does not refresh static tokens.

```go
c, err := client.New(
    "acme",
    option.WithStaticToken("<token>"),
)
```

## URL derivation

| Value    | Default                                  |
| -------- | ---------------------------------------- |
| Base URL | `https://{teamSlug}.teams.auth0.com`     |
| Audience | `https://{teamSlug}.teams.auth0.com/api` |
| Token    | `https://{tenantDomain}/oauth/token`     |

The base URL deliberately omits the `/api` suffix because the generated endpoints prepend
it. The audience includes `/api` because that is the value the Teams API validates.

Both can be overridden: pass `option.WithBaseURL(...)` to set the base URL and
`option.WithClientCredentialsAndAudience(...)` to set the audience.

## Telemetry

By default the SDK sends an `Auth0-Client` header containing the SDK name, version, and Go
runtime version. You can add custom environment entries or disable the header entirely.

```go
// Add a custom entry to the Auth0-Client telemetry header.
c, err := client.New(
    "acme",
    option.WithStaticToken("<token>"),
    option.WithAuth0ClientEnvEntry("myapp", "2.0.0"),
)

// Disable the Auth0-Client header.
c, err = client.New(
    "acme",
    option.WithStaticToken("<token>"),
    option.WithNoAuth0ClientInfo(),
)
```

## Debug logging

`option.WithDebug(true)` logs each request and response, with sensitive headers (such as
`Authorization`) redacted.

```go
c, err := client.New(
    "acme",
    option.WithClientCredentials(ctx, "acme-auth.us.auth0.com", "<client-id>", "<client-secret>"),
    option.WithDebug(true),
)
```
