package auth0

import (
	"testing"

	core "github.com/auth0/teams-go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- SanitizeDomain ---

func TestSanitizeDomain_PlainDomain(t *testing.T) {
	assert.Equal(t, "mytenant.auth0.com", SanitizeDomain("mytenant.auth0.com"))
}

func TestSanitizeDomain_HTTPSPrefix(t *testing.T) {
	assert.Equal(t, "mytenant.auth0.com", SanitizeDomain("https://mytenant.auth0.com"))
}

func TestSanitizeDomain_HTTPPrefix(t *testing.T) {
	assert.Equal(t, "mytenant.auth0.com", SanitizeDomain("http://mytenant.auth0.com"))
}

func TestSanitizeDomain_TrailingSlash(t *testing.T) {
	assert.Equal(t, "mytenant.auth0.com", SanitizeDomain("mytenant.auth0.com/"))
}

func TestSanitizeDomain_HTTPSAndTrailingSlash(t *testing.T) {
	assert.Equal(t, "mytenant.auth0.com", SanitizeDomain("https://mytenant.auth0.com/"))
}

func TestSanitizeDomain_MultipleTrailingSlashes(t *testing.T) {
	assert.Equal(t, "mytenant.auth0.com", SanitizeDomain("https://mytenant.auth0.com///"))
}

// --- DeriveURLs ---

func TestDeriveURLs_Defaults(t *testing.T) {
	opts := &core.RequestOptions{TenantDomain: "acme-auth.us.auth0.com"}
	baseURL, tokenURL, audience := DeriveURLs("acme", opts)
	// Base URL must NOT include /api (raw_client.go prepends /api/ itself).
	assert.Equal(t, "https://acme.teams.auth0.com", baseURL)
	assert.Equal(t, "https://acme-auth.us.auth0.com/oauth/token", tokenURL)
	// Audience MUST include /api (the aud claim the server validates).
	assert.Equal(t, "https://acme.teams.auth0.com/api", audience)
}

func TestDeriveURLs_CustomBaseURL(t *testing.T) {
	opts := &core.RequestOptions{TenantDomain: "acme-auth.us.auth0.com", BaseURL: "https://custom.example.com"}
	baseURL, tokenURL, audience := DeriveURLs("acme", opts)
	assert.Equal(t, "https://custom.example.com", baseURL)
	assert.Equal(t, "https://acme-auth.us.auth0.com/oauth/token", tokenURL)
	// The default audience is derived from the slug, independent of the base URL
	// override. Supply an explicit audience when the environment's API host
	// differs from the slug-derived host.
	assert.Equal(t, "https://acme.teams.auth0.com/api", audience)
}

func TestDeriveURLs_CustomAudience(t *testing.T) {
	opts := &core.RequestOptions{TenantDomain: "acme-auth.us.auth0.com", Audience: "https://custom-api.example.com/"}
	baseURL, tokenURL, audience := DeriveURLs("acme", opts)
	assert.Equal(t, "https://acme.teams.auth0.com", baseURL)
	assert.Equal(t, "https://acme-auth.us.auth0.com/oauth/token", tokenURL)
	assert.Equal(t, "https://custom-api.example.com/", audience)
}

func TestDeriveURLs_CustomBaseURLAndAudience(t *testing.T) {
	opts := &core.RequestOptions{
		TenantDomain: "acme-auth.us.auth0.com",
		BaseURL:      "https://custom.example.com",
		Audience:     "https://custom-api.example.com/",
	}
	baseURL, tokenURL, audience := DeriveURLs("acme", opts)
	assert.Equal(t, "https://custom.example.com", baseURL)
	assert.Equal(t, "https://acme-auth.us.auth0.com/oauth/token", tokenURL)
	assert.Equal(t, "https://custom-api.example.com/", audience)
}

func TestDeriveURLs_NoTenantDomain(t *testing.T) {
	// Without a tenant domain (static token / token source), no token URL is
	// derived because those modes never contact the token endpoint.
	opts := &core.RequestOptions{}
	baseURL, tokenURL, audience := DeriveURLs("acme", opts)
	assert.Equal(t, "https://acme.teams.auth0.com", baseURL)
	assert.Equal(t, "", tokenURL)
	assert.Equal(t, "https://acme.teams.auth0.com/api", audience)
}

// --- ValidateOptions ---

func TestValidateOptions_EmptyTeamSlug(t *testing.T) {
	opts := &core.RequestOptions{Token: "tok"}
	err := ValidateOptions("", opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "team slug must not be empty")
}

func TestValidateOptions_EmptyTeamSlug_WithBaseURL(t *testing.T) {
	// The team slug is always required, even when an explicit base URL is set.
	opts := &core.RequestOptions{Token: "tok", BaseURL: "https://my-team.teams.auth0.com"}
	err := ValidateOptions("", opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "team slug must not be empty")
}

func TestValidateOptions_EmptyTenantDomain_ClientCredentials(t *testing.T) {
	// tenant_domain is required for client credentials (it builds the token URL).
	opts := &core.RequestOptions{ClientID: "id", ClientSecret: "secret"}
	err := ValidateOptions("acme", opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tenant domain must not be empty when using client credentials")
}

func TestValidateOptions_EmptyTenantDomain_StaticToken(t *testing.T) {
	// A static token never contacts the token endpoint, so tenant_domain is not
	// required; team_slug alone is enough.
	opts := &core.RequestOptions{Token: "tok"}
	err := ValidateOptions("acme", opts)
	assert.NoError(t, err)
}

func TestValidateOptions_TenantDomainNotADomain(t *testing.T) {
	// A tenant domain without a "." looks like a slug, not a domain. Only the
	// client-credentials flow validates the domain shape.
	opts := &core.RequestOptions{TenantDomain: "acme", ClientID: "id", ClientSecret: "secret"}
	err := ValidateOptions("acme", opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not look like a domain")
}

func TestValidateOptions_NoAuth(t *testing.T) {
	opts := &core.RequestOptions{}
	err := ValidateOptions("acme", opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must provide either client credentials")
}

func TestValidateOptions_BothAuthModes(t *testing.T) {
	opts := &core.RequestOptions{
		TenantDomain: "acme-auth.us.auth0.com",
		ClientID:     "id",
		ClientSecret: "secret",
		Token:        "tok",
	}
	err := ValidateOptions("acme", opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only one authentication mode")
}

func TestValidateOptions_ClientCredentials_Valid(t *testing.T) {
	opts := &core.RequestOptions{
		TenantDomain: "acme-auth.us.auth0.com",
		ClientID:     "id",
		ClientSecret: "secret",
	}
	err := ValidateOptions("acme", opts)
	assert.NoError(t, err)
}

func TestValidateOptions_ClientCredentials_MissingClientID(t *testing.T) {
	opts := &core.RequestOptions{
		TenantDomain: "acme-auth.us.auth0.com",
		ClientSecret: "secret",
	}
	err := ValidateOptions("acme", opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client ID must not be empty")
}

func TestValidateOptions_ClientCredentials_OrphanClientID(t *testing.T) {
	opts := &core.RequestOptions{
		TenantDomain: "acme-auth.us.auth0.com",
		ClientID:     "id",
	}
	err := ValidateOptions("acme", opts)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client secret must not be empty")
}

func TestValidateOptions_StaticToken_Valid(t *testing.T) {
	opts := &core.RequestOptions{Token: "my-token"}
	err := ValidateOptions("acme", opts)
	assert.NoError(t, err)
}
