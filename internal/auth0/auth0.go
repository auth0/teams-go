package auth0

import (
	"errors"
	"strings"

	core "github.com/auth0/teams-go/core"
)

// SanitizeDomain strips protocol prefixes and trailing slashes from the domain.
// It loops to handle edge cases like "http://http://example.com".
func SanitizeDomain(domain string) string {
	d := domain
	for {
		trimmed := strings.TrimPrefix(d, "https://")
		trimmed = strings.TrimPrefix(trimmed, "http://")
		if trimmed == d {
			break
		}
		d = trimmed
	}
	d = strings.TrimRight(d, "/")
	return d
}

// DeriveURLs computes the base URL, token URL, and audience for a team.
//
// The base URL and audience deliberately differ. The generated raw_client.go
// files build endpoint paths as baseURL + "/api/tenants" (and similar), so the
// base URL must NOT include an "/api" suffix or every request would resolve to
// "/api/api/...". The audience, on the other hand, is the exact "aud" claim the
// Teams API validates and DOES include "/api".
//
// The default audience is derived from the team slug
// (https://{teamSlug}.teams.auth0.com/api), independent of any base URL
// override. When the base URL points at a non-production environment whose API
// host differs from the slug-derived host, supply the matching audience
// explicitly via option.WithClientCredentialsAndAudience.
//
// The token URL is derived from the tenant domain carried on the options (set
// via option.WithClientCredentials). It is only meaningful for the
// client-credentials flow, so it is empty when no tenant domain is configured.
//
// An explicit base URL override (via option.WithBaseURL) and an explicit
// audience override (via option.WithClientCredentialsAndAudience) take
// precedence over the derived values.
func DeriveURLs(teamSlug string, options *core.RequestOptions) (baseURL, tokenURL, audience string) {
	baseURL = options.BaseURL
	if baseURL == "" {
		baseURL = "https://" + teamSlug + ".teams.auth0.com"
	}
	if options.TenantDomain != "" {
		tokenURL = "https://" + SanitizeDomain(options.TenantDomain) + "/oauth/token"
	}
	audience = options.Audience
	if audience == "" {
		audience = "https://" + teamSlug + ".teams.auth0.com/api"
	}
	return
}

// ValidateOptions checks that the team slug and auth options are valid.
//
// The team slug is always required: it is the SDK's identity and derives both
// the base URL and the token audience. option.WithBaseURL overrides where
// requests are sent, but it does not replace the slug.
//
// The tenant domain lives on the options (set via option.WithClientCredentials)
// because it is only used to build the OAuth token URL, which only the
// client-credentials flow contacts.
func ValidateOptions(teamSlug string, options *core.RequestOptions) error {
	if teamSlug == "" {
		return errors.New("auth0: team slug must not be empty")
	}

	// Detect auth modes by their unique discriminating fields. ClientSecret
	// uniquely identifies client credentials; Token a static bearer; and
	// TokenSource a custom source.
	hasClientCredentials := options.ClientSecret != ""
	hasStaticToken := options.Token != ""
	hasTokenSource := options.TokenSource != nil

	// ClientID alone (without ClientSecret) looks like an incomplete client
	// credentials config.
	hasOrphanClientID := options.ClientID != "" && !hasClientCredentials

	// tenant_domain is only used to build the OAuth token URL, which only the
	// client-credentials flow needs. A static token or a custom token source
	// never contacts the token endpoint, so tenant_domain is not required for
	// those modes.
	if hasClientCredentials || hasOrphanClientID {
		tenantDomain := options.TenantDomain
		if tenantDomain == "" {
			return errors.New("auth0: tenant domain must not be empty when using client credentials (set it via option.WithClientCredentials)")
		}
		// A tenant domain without a "." is almost certainly not a domain. This
		// catches the common misuse of passing the team slug as the tenant
		// domain, e.g. WithClientCredentials(ctx, teamSlug, ...).
		if !strings.Contains(SanitizeDomain(tenantDomain), ".") {
			return errors.New("auth0: tenant domain does not look like a domain (did you pass the team slug as the tenant domain?)")
		}
	}

	// Count how many auth modes are active.
	count := 0
	if hasClientCredentials || hasOrphanClientID {
		count++
	}
	if hasStaticToken {
		count++
	}
	if hasTokenSource {
		count++
	}

	if count > 1 {
		return errors.New("auth0: only one authentication mode may be used (client credentials, token source, or static token)")
	}
	if count == 0 {
		return errors.New("auth0: must provide either client credentials (WithClientCredentials), a token source (WithTokenSource), or a static token (WithStaticToken/WithToken)")
	}
	if hasClientCredentials || hasOrphanClientID {
		if options.ClientID == "" {
			return errors.New("auth0: client credentials: client ID must not be empty")
		}
		if options.ClientSecret == "" {
			return errors.New("auth0: client credentials: client secret must not be empty")
		}
	}
	return nil
}
