package client

import (
	core "github.com/auth0/teams-go/core"
	"github.com/auth0/teams-go/internal/auth0"
	"github.com/auth0/teams-go/internal/transport"
	option "github.com/auth0/teams-go/option"
)

// New creates a new Teams client with Auth0 authentication and telemetry.
// It accepts the team slug and options built from the option package. The team
// slug is required.
//
// The base URL is derived as https://{teamSlug}.teams.auth0.com (the generated
// endpoints append /api themselves), while the token audience is derived as
// https://{teamSlug}.teams.auth0.com/api. option.WithBaseURL overrides where
// requests are sent without changing the derived audience; supply
// option.WithClientCredentialsAndAudience to override the audience as well.
//
// Authentication modes (mutually exclusive):
//   - option.WithClientCredentials(ctx, tenantDomain, clientID, clientSecret) — M2M OAuth2
//   - option.WithTokenSource(tokenSource) — custom oauth2.TokenSource
//   - option.WithStaticToken(token) or option.WithToken(token) — static bearer token
func New(teamSlug string, opts ...option.RequestOption) (*Teams, error) {
	options := core.NewRequestOptions(opts...)

	if err := auth0.ValidateOptions(teamSlug, options); err != nil {
		return nil, err
	}

	baseURL, tokenURL, audience := auth0.DeriveURLs(teamSlug, options)

	httpClient := transport.CloneHTTPClient(options.HTTPClient)

	chain, err := transport.BuildChain(transport.ChainConfig{
		Base:              httpClient.Transport,
		ClientID:          options.ClientID,
		ClientSecret:      options.ClientSecret,
		TokenURL:          tokenURL,
		Audience:          audience,
		OAuthContext:      options.OAuthContext,
		TokenSource:       options.TokenSource,
		NoAuth0ClientInfo: options.NoAuth0ClientInfo,
		Auth0ClientEnv:    options.Auth0ClientEnv,
		Debug:             options.Debug,
	})
	if err != nil {
		return nil, err
	}
	httpClient.Transport = chain

	fernOpts := []option.RequestOption{
		option.WithHTTPClient(httpClient),
		option.WithBaseURL(baseURL),
	}
	if options.MaxAttempts > 0 {
		fernOpts = append(fernOpts, option.WithMaxAttempts(options.MaxAttempts))
	}
	if options.Token != "" {
		fernOpts = append(fernOpts, option.WithToken(options.Token))
	}
	if len(options.HTTPHeader) > 0 {
		fernOpts = append(fernOpts, option.WithHTTPHeader(options.HTTPHeader))
	}

	return NewWithOptions(fernOpts...), nil
}
