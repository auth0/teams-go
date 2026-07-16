package option

import (
	"context"

	core "github.com/auth0/teams-go/core"
	"golang.org/x/oauth2"
)

// WithClientCredentials configures OAuth2 client credentials authentication.
// The SDK will automatically fetch, cache, and refresh access tokens
// using the client_credentials grant type.
//
// tenantDomain is the Auth0 tenant domain (e.g. "acme.us.auth0.com") used to
// build the OAuth token endpoint. The domain is sanitized automatically
// (https:// prefix and trailing slashes are stripped).
//
// The audience defaults to https://{teamSlug}.teams.auth0.com/api.
// To specify a custom audience, use WithClientCredentialsAndAudience instead.
func WithClientCredentials(ctx context.Context, tenantDomain, clientID, clientSecret string) *core.ClientCredentialsOption {
	return &core.ClientCredentialsOption{
		Context:      ctx,
		TenantDomain: tenantDomain,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// WithClientCredentialsAndAudience configures OAuth2 client credentials authentication
// with a custom audience. The SDK will automatically fetch, cache, and refresh
// access tokens using the client_credentials grant type.
//
// tenantDomain is the Auth0 tenant domain (e.g. "acme.us.auth0.com") used to
// build the OAuth token endpoint. See WithClientCredentials for details.
func WithClientCredentialsAndAudience(ctx context.Context, tenantDomain, clientID, clientSecret, audience string) *core.ClientCredentialsOption {
	return &core.ClientCredentialsOption{
		Context:      ctx,
		TenantDomain: tenantDomain,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
	}
}

// WithTokenSource configures a custom oauth2.TokenSource for authentication.
// The SDK will use the provided token source to obtain access tokens.
// This allows full control over how tokens are fetched, cached, and refreshed.
//
// The token source is mutually exclusive with other authentication modes
// (client credentials and static token).
func WithTokenSource(tokenSource oauth2.TokenSource) *core.TokenSourceOption {
	return &core.TokenSourceOption{
		TokenSource: tokenSource,
	}
}

// WithStaticToken configures a static bearer token for authentication.
// This is an alias for WithToken — provided for symmetry with WithClientCredentials.
func WithStaticToken(token string) *core.TokenOption {
	return &core.TokenOption{
		Token: token,
	}
}

// WithNoAuth0ClientInfo disables the Auth0-Client telemetry header.
func WithNoAuth0ClientInfo() *core.NoAuth0ClientInfoOption {
	return &core.NoAuth0ClientInfoOption{}
}

// WithAuth0ClientEnvEntry adds a custom key-value entry to the Auth0-Client
// header's env map (e.g., WithAuth0ClientEnvEntry("myapp", "1.0.0")).
func WithAuth0ClientEnvEntry(key, value string) *core.Auth0ClientEnvEntryOption {
	return &core.Auth0ClientEnvEntryOption{
		Key:   key,
		Value: value,
	}
}

// WithDebug enables HTTP request/response debug logging. When enabled, the SDK
// logs every outgoing request and incoming response via log.Printf. Sensitive
// headers (Authorization, Cookie, etc.) are redacted automatically.
func WithDebug(debug bool) *core.DebugOption {
	return &core.DebugOption{
		Debug: debug,
	}
}
