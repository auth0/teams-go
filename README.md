![Go SDK for Auth0 Teams](https://cdn.auth0.com/website/sdks/banners/teams-go-banner.png)

<div align="center">

[![GoDoc](https://pkg.go.dev/badge/github.com/auth0/teams-go.svg)](https://pkg.go.dev/github.com/auth0/teams-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/auth0/teams-go?style=flat-square)](https://goreportcard.com/report/github.com/auth0/teams-go)
[![Release](https://img.shields.io/github/v/release/auth0/teams-go?include_prereleases&style=flat-square)](https://github.com/auth0/teams-go/releases)
[![License](https://img.shields.io/github/license/auth0/teams-go.svg?style=flat-square)](https://github.com/auth0/teams-go/blob/main/LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/auth0/teams-go/ci.yml?branch=main&style=flat-square)](https://github.com/auth0/teams-go/actions?query=branch%3Amain)
[![Codecov](https://img.shields.io/codecov/c/github/auth0/teams-go/main?style=flat-square)](https://codecov.io/gh/auth0/teams-go/tree/main)

</div>

<div align="center">

📚 [Documentation](#documentation) • 🚀 [Getting Started](#getting-started) • 💬 [Feedback](#feedback)

</div>

The Auth0 Teams Go SDK provides convenient, strongly-typed access to the Auth0 Teams API from Go. It manages authentication, automatic token refresh, retries, and telemetry for you, so you can focus on calling the API.

## Documentation

- [Godoc](https://pkg.go.dev/github.com/auth0/teams-go) - the generated API documentation for this SDK.
- [Examples](./examples.md) - telemetry, debug logging, and other usage examples.
- [Runnable examples](./examples) - a `go run`-able program for each authentication mode and resource.
- [API Reference](./reference.md) - a full reference for every endpoint and type in this library.
- [Docs site](https://www.auth0.com/docs) - explore our documentation, quickstart guides, and API reference.

## Getting Started

### Requirements

This library follows the [same support policy as Go](https://go.dev/doc/devel/release#policy). The last two major Go releases are actively supported, and compatibility issues with older versions are not guaranteed.

- Go 1.25+

### Installation

```shell
go get github.com/auth0/teams-go
```

### Usage

Instantiate and use the client with the following:

```go
package example

import (
    context "context"

    teams "github.com/auth0/teams-go"
    client "github.com/auth0/teams-go/client"
    option "github.com/auth0/teams-go/option"
)

func do() {
    c, err := client.New(
        "acme", // team slug
        option.WithClientCredentials(
            context.TODO(),
            "acme-auth.us.auth0.com", // Auth0 tenant domain
            "<client-id>",
            "<client-secret>",
        ),
    )
    if err != nil {
        return
    }
    request := &teams.InviteMembersRequestContent{
        Email: "jane.smith@company.com",
        Role:  teams.TeamMemberRoleEnumTeamsContributor,
        TenantIDs: []teams.TenantID{
            "538c9e21-e3d5-4ad6-b3d0-352c62369fb0",
            "b1c9e21e-1234-4ad6-b3d0-352c6236abcd",
        },
        TenantRoles: []teams.TenantMemberRoleEnum{
            teams.TenantMemberRoleEnumOwner,
            teams.TenantMemberRoleEnumEditorUsers,
        },
        ClientIDs: []teams.ClientID{
            "client_123",
            "client_456",
        },
    }
    c.Members.Invite(
        context.TODO(),
        request,
    )
}
```

`client.New(teamSlug, opts...)` is the recommended entry point. It manages Auth0 authentication, token refresh, and telemetry for you.

## Authentication

The client supports several mutually exclusive authentication modes. The team slug is the only positional argument.

### Client Credentials (M2M)

Tokens are fetched and refreshed automatically. The Auth0 tenant domain is passed right after the context.

```go
c, err := client.New(
    "acme",
    option.WithClientCredentials(ctx, "acme-auth.us.auth0.com", "<client-id>", "<client-secret>"),
)
```

You can also specify an explicit token audience:

```go
c, err := client.New(
    "acme",
    option.WithClientCredentialsAndAudience(ctx, "acme-auth.us.auth0.com", "<client-id>", "<client-secret>", "https://acme.teams.auth0.com/api"),
)
```

### Custom Token Source

Supply your own `oauth2.TokenSource`.

```go
c, err := client.New(
    "acme",
    option.WithTokenSource(myTokenSource),
)
```

### Static Token

Use a static bearer token (no automatic refresh).

```go
c, err := client.New(
    "acme",
    option.WithStaticToken("<token>"),
)
```

By default the base URL is derived as `https://{teamSlug}.teams.auth0.com` and the token audience as `https://{teamSlug}.teams.auth0.com/api`. The token endpoint is `https://{tenantDomain}/oauth/token`. See [examples.md](./examples.md) for telemetry and debug examples.

## Base URL

You can override the derived base URL with the `option.WithBaseURL` option. You can configure any arbitrary base URL, which is particularly useful in test environments.

```go
c, err := client.New(
    "acme",
    option.WithStaticToken("<token>"),
    option.WithBaseURL("https://acme.teams.auth0.com"),
)
```

## Request Options

A variety of request options are included to adapt the behavior of the library, which includes configuring authorization tokens, or providing your own instrumented `*http.Client`.

These request options can either be specified on the client so that they're applied on every request, or for an individual request, like so:

> Providing your own `*http.Client` is recommended. Otherwise, the `http.DefaultClient` will be used, and your client will wait indefinitely for a response (unless the per-request, context-based timeout is used).

```go
// Specify default options applied on every request.
c, err := client.New(
    "acme",
    option.WithStaticToken("<token>"),
    option.WithHTTPClient(
        &http.Client{
            Timeout: 5 * time.Second,
        },
    ),
)

// Specify options for an individual request.
response, err := c.Members.Invite(
    ...,
    option.WithToken("<token>"),
)
```

## Raw Responses

You can access the raw HTTP response data by using the `WithRawResponse` field on the client. This is useful when you need to examine the response headers received from the API call. (When the endpoint is paginated, the raw HTTP response data will be included automatically in the Page response object.)

```go
response, err := client.Members.WithRawResponse.Invite(...)
if err != nil {
    return err
}
fmt.Printf("Got response headers: %v", response.Header)
fmt.Printf("Got status code: %d", response.StatusCode)
```

## Error Handling

Structured error types are returned from API calls that return non-success status codes. These errors are compatible with the `errors.Is` and `errors.As` APIs, so you can access the error like so:

```go
response, err := client.Members.Invite(...)
if err != nil {
    var apiError *core.APIError
    if errors.As(err, apiError) {
        // Do something with the API error ...
    }
    return err
}
```

## Explicit Null Values

If you want to send the explicit `null` JSON value through an optional parameter, you can use the setters\
that come with every object. Calling a setter method for a property will flip a bit in the `explicitFields` bitfield for that setter's object; during serialization, any property with a flipped bit will have its omittable status stripped, so zero or `nil` values will be sent explicitly rather than omitted altogether:

```go
type ExampleRequest struct {
    // An optional string parameter.
    Name *string `json:"name,omitempty" url:"-"`

    // Private bitmask of fields set to an explicit value and therefore not to be omitted
    explicitFields *big.Int `json:"-" url:"-"`
}

request := &ExampleRequest{}
request.SetName(nil)

response, err := client.Members.Invite(ctx, request, ...)
```

## Retries

The SDK is instrumented with automatic retries with exponential backoff. A request will be retried as long as the request is deemed retryable and the number of retry attempts has not grown larger than the configured retry limit (default: 2).

Which status codes are retried depends on the `retryStatusCodes` generator configuration:

**`legacy`** (current default): retries on
- [408](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/408) (Timeout)
- [429](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429) (Too Many Requests)
- [5XX](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#server_error_responses) (All server errors, including 500)

**`recommended`**: retries on
- [408](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/408) (Timeout)
- [429](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429) (Too Many Requests)
- [502](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/502) (Bad Gateway)
- [503](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/503) (Service Unavailable)
- [504](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/504) (Gateway Timeout)

If the `Retry-After` header is present in the response, the SDK will prioritize respecting its value exactly over the default exponential backoff.

Use the `option.WithMaxAttempts` option to configure this behavior for the entire client or an individual request:

```go
c, err := client.New(
    "acme",
    option.WithStaticToken("<token>"),
    option.WithMaxAttempts(1),
)

response, err := c.Members.Invite(
    ...,
    option.WithMaxAttempts(1),
)
```

## Timeouts

Setting a timeout for each individual request is as simple as using the standard context library. Setting a one second timeout for an individual API call looks like the following:

```go
ctx, cancel := context.WithTimeout(ctx, time.Second)
defer cancel()

response, err := client.Members.Invite(ctx, ...)
```

## Feedback

### Contributing

We appreciate feedback and contribution to this SDK! Before you get started, please see the following:

- [Auth0's general contribution guidelines](https://github.com/auth0/open-source-template/blob/master/GENERAL-CONTRIBUTING.md)
- [Auth0's code of conduct guidelines](https://github.com/auth0/open-source-template/blob/master/CODE-OF-CONDUCT.md)

While we value open-source contributions to this SDK, this library is generated programmatically. Additions made directly to this library would have to be moved over to our generation code, otherwise they would be overwritten upon the next generated release. Feel free to open a PR as a proof of concept, but know that we will not be able to merge it as-is. We suggest opening an issue first to discuss with us! On the other hand, contributions to the README are always very welcome.

### Raise an Issue

To provide feedback or report a bug, please [raise an issue on our issue tracker](https://github.com/auth0/teams-go/issues).

### Vulnerability Reporting

Please do not report security vulnerabilities on the public GitHub issue tracker. The [Responsible Disclosure Program](https://auth0.com/responsible-disclosure-policy) details the procedure for disclosing security issues.

## What is Auth0?

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://cdn.auth0.com/website/sdks/logos/auth0_dark_mode.png" width="150">
    <source media="(prefers-color-scheme: light)" srcset="https://cdn.auth0.com/website/sdks/logos/auth0_light_mode.png" width="150">
    <img alt="Auth0 Logo" src="https://cdn.auth0.com/website/sdks/logos/auth0_light_mode.png" width="150">
  </picture>
</p>
<p align="center">
  Auth0 is an easy to implement, adaptable authentication and authorization platform. To learn more checkout <a href="https://auth0.com/why-auth0">Why Auth0?</a>
</p>
<p align="center">
  This project is licensed under the Apache-2.0 license. See the <a href="./LICENSE">LICENSE</a> file for more info.
</p>
