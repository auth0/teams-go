package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"strings"
	"testing"

	teams "github.com/auth0/teams-go"
	core "github.com/auth0/teams-go/core"
	"github.com/auth0/teams-go/internal/auth0"
	"github.com/auth0/teams-go/internal/telemetry"
	option "github.com/auth0/teams-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// --- Domain Sanitization ---

func TestSanitizeDomain(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"mytenant.auth0.com", "mytenant.auth0.com"},
		{"https://mytenant.auth0.com", "mytenant.auth0.com"},
		{"http://mytenant.auth0.com", "mytenant.auth0.com"},
		{"https://mytenant.auth0.com/", "mytenant.auth0.com"},
		{"https://mytenant.auth0.com///", "mytenant.auth0.com"},
		{"mytenant.auth0.com/", "mytenant.auth0.com"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, auth0.SanitizeDomain(tt.input))
		})
	}
}

// TestDerivedURLs verifies the deliberate base-URL/audience split: the base URL
// must NOT carry an "/api" suffix (the generated endpoints prepend it), while the
// audience MUST carry "/api".
func TestDerivedURLs(t *testing.T) {
	baseURL, tokenURL, audience := auth0.DeriveURLs("acme", core.NewRequestOptions(option.WithClientCredentials(context.Background(), "acme-auth.us.auth0.com", "id", "secret")))
	assert.Equal(t, "https://acme.teams.auth0.com", baseURL)
	assert.Equal(t, "https://acme-auth.us.auth0.com/oauth/token", tokenURL)
	assert.Equal(t, "https://acme.teams.auth0.com/api", audience)
}

// --- Telemetry ---

func TestBuildAuth0ClientInfo(t *testing.T) {
	info := telemetry.BuildAuth0ClientInfo(nil)
	assert.Equal(t, teams.SDKName, info.Name)
	assert.Equal(t, teams.SDKVersion, info.Version)
	assert.Equal(t, runtime.Version(), info.Env["go"])
}

func TestBuildAuth0ClientInfoWithCustomEnv(t *testing.T) {
	custom := map[string]string{"myapp": "1.0.0"}
	info := telemetry.BuildAuth0ClientInfo(custom)
	assert.Equal(t, "1.0.0", info.Env["myapp"])
	assert.Equal(t, runtime.Version(), info.Env["go"])
}

func TestEncodeAuth0ClientInfo(t *testing.T) {
	info := telemetry.Auth0ClientInfo{
		Name:    "teams-go",
		Version: "0.0.1",
		Env:     map[string]string{"go": "go1.21"},
	}
	encoded, err := telemetry.EncodeAuth0ClientInfo(info)
	require.NoError(t, err)

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)

	var roundTrip telemetry.Auth0ClientInfo
	err = json.Unmarshal(decoded, &roundTrip)
	require.NoError(t, err)

	assert.Equal(t, info.Name, roundTrip.Name)
	assert.Equal(t, info.Version, roundTrip.Version)
	assert.Equal(t, info.Env, roundTrip.Env)
}

// --- Option Validation ---

func TestNew_EmptyTeamSlug(t *testing.T) {
	_, err := New("", option.WithStaticToken("tok"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "team slug must not be empty")
}

func TestNew_EmptyTenantDomain(t *testing.T) {
	// Client credentials without a tenant domain must fail.
	_, err := New("acme", option.WithClientCredentials(context.Background(), "", "id", "secret"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tenant domain must not be empty")
}

// TestNew_TenantDomainNotADomain guards against passing the team slug as the
// tenant domain. A tenant domain without a "." almost certainly means the slug
// was passed by mistake.
func TestNew_TenantDomainNotADomain(t *testing.T) {
	_, err := New("acme", option.WithClientCredentials(context.Background(), "acme", "id", "secret"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not look like a domain")
}

func TestNew_NoAuth(t *testing.T) {
	_, err := New("acme")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must provide either client credentials")
}

func TestNew_BothAuthModes(t *testing.T) {
	_, err := New("acme",
		option.WithClientCredentials(context.Background(), "mytenant.auth0.com", "id", "secret"),
		option.WithStaticToken("tok"),
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only one authentication mode")
}

func TestNew_ClientCredentials_MissingClientID(t *testing.T) {
	_, err := New("acme",
		option.WithClientCredentials(context.Background(), "mytenant.auth0.com", "", "secret"),
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client ID must not be empty")
}

func TestNew_ClientCredentials_MissingClientSecret(t *testing.T) {
	_, err := New("acme",
		option.WithClientCredentials(context.Background(), "mytenant.auth0.com", "id", ""),
	)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "client secret must not be empty")
}

func TestNew_StaticToken(t *testing.T) {
	c, err := New("acme", option.WithStaticToken("my-token"))
	require.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c.Tenants)
	assert.NotNil(t, c.Members)
}

func TestNew_WithToken(t *testing.T) {
	c, err := New("acme", option.WithToken("my-token"))
	require.NoError(t, err)
	assert.NotNil(t, c)
}

// --- Transport Chain Headers ---

func TestTransportHeaders_StaticToken(t *testing.T) {
	var capturedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New("acme",
		option.WithStaticToken("my-token"),
		option.WithBaseURL(server.URL),
		option.WithHTTPClient(&http.Client{}),
	)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	require.NoError(t, err)

	resp, err := c.options.HTTPClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Verify transport-set headers.
	assert.Equal(t, teams.SDKName+"/"+teams.SDKVersion, capturedHeaders.Get("User-Agent"))

	auth0Client := capturedHeaders.Get("Auth0-Client")
	assert.NotEmpty(t, auth0Client)

	decoded, err := base64.StdEncoding.DecodeString(auth0Client)
	require.NoError(t, err)

	var info telemetry.Auth0ClientInfo
	err = json.Unmarshal(decoded, &info)
	require.NoError(t, err)
	assert.Equal(t, teams.SDKName, info.Name)
	assert.Equal(t, teams.SDKVersion, info.Version)
	assert.Equal(t, runtime.Version(), info.Env["go"])

	// Verify X-Fern-* headers are NOT present.
	assert.Empty(t, capturedHeaders.Get("X-Fern-Language"))
	assert.Empty(t, capturedHeaders.Get("X-Fern-SDK-Name"))
	assert.Empty(t, capturedHeaders.Get("X-Fern-SDK-Version"))
}

func TestTransportHeaders_NoAuth0ClientInfo(t *testing.T) {
	var capturedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New("acme",
		option.WithStaticToken("my-token"),
		option.WithBaseURL(server.URL),
		option.WithNoAuth0ClientInfo(),
	)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	require.NoError(t, err)

	resp, err := c.options.HTTPClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, teams.SDKName+"/"+teams.SDKVersion, capturedHeaders.Get("User-Agent"))
	assert.Empty(t, capturedHeaders.Get("Auth0-Client"))
}

func TestTransportHeaders_CustomEnvEntry(t *testing.T) {
	var capturedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c, err := New("acme",
		option.WithStaticToken("my-token"),
		option.WithBaseURL(server.URL),
		option.WithAuth0ClientEnvEntry("myapp", "2.0.0"),
	)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	require.NoError(t, err)

	resp, err := c.options.HTTPClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	auth0Client := capturedHeaders.Get("Auth0-Client")
	decoded, err := base64.StdEncoding.DecodeString(auth0Client)
	require.NoError(t, err)

	var info telemetry.Auth0ClientInfo
	err = json.Unmarshal(decoded, &info)
	require.NoError(t, err)
	assert.Equal(t, "2.0.0", info.Env["myapp"])
	assert.Equal(t, runtime.Version(), info.Env["go"])
}

// --- M2M Client Credentials Flow ---

func TestNew_ClientCredentials_M2MFlow(t *testing.T) {
	tokenServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)
		assert.Contains(t, bodyStr, "grant_type=client_credentials")
		assert.Contains(t, bodyStr, "audience=")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"test-access-token","token_type":"Bearer","expires_in":86400}`))
	}))
	defer tokenServer.Close()

	var capturedHeaders http.Header
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	defer apiServer.Close()

	tokenHost := strings.TrimPrefix(tokenServer.URL, "https://")
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tokenServer.Client())

	c, err := New("acme",
		option.WithClientCredentialsAndAudience(ctx, tokenHost, "test-client-id", "test-client-secret", "https://test-audience/"),
		option.WithBaseURL(apiServer.URL),
	)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", apiServer.URL+"/test", nil)
	require.NoError(t, err)

	resp, err := c.options.HTTPClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// Verify oauth2 transport set Authorization.
	assert.Equal(t, "Bearer test-access-token", capturedHeaders.Get("Authorization"))
	// Verify transport chain headers.
	assert.Equal(t, teams.SDKName+"/"+teams.SDKVersion, capturedHeaders.Get("User-Agent"))
	assert.NotEmpty(t, capturedHeaders.Get("Auth0-Client"))
}

// TestNew_ClientCredentials_DefaultAudience verifies that, with no explicit
// audience, the slug-derived audience is sent to the token endpoint.
func TestNew_ClientCredentials_DefaultAudience(t *testing.T) {
	tokenServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)
		// Default audience is https://{teamSlug}.teams.auth0.com/api.
		assert.Contains(t, bodyStr, "audience=https")
		assert.Contains(t, bodyStr, "acme.teams.auth0.com")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":3600}`))
	}))
	defer tokenServer.Close()

	tokenHost := strings.TrimPrefix(tokenServer.URL, "https://")
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tokenServer.Client())

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer apiServer.Close()

	// No audience override → uses default derived from the team slug.
	c, err := New("acme",
		option.WithClientCredentials(ctx, tokenHost, "id", "secret"),
		option.WithBaseURL(apiServer.URL),
	)
	require.NoError(t, err)

	req, _ := http.NewRequest("GET", apiServer.URL+"/test", nil)
	resp, err := c.options.HTTPClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
}

// --- Custom Base URL ---

func TestNew_CustomBaseURL(t *testing.T) {
	c, err := New("acme",
		option.WithStaticToken("tok"),
		option.WithBaseURL("https://custom.example.com/api"),
	)
	require.NoError(t, err)
	assert.Equal(t, "https://custom.example.com/api", c.baseURL)
}

// TestNew_DefaultBaseURL verifies the base URL has no "/api" suffix (the
// generated endpoints prepend it).
func TestNew_DefaultBaseURL(t *testing.T) {
	c, err := New("acme",
		option.WithStaticToken("tok"),
	)
	require.NoError(t, err)
	assert.Equal(t, "https://acme.teams.auth0.com", c.baseURL)
}

// --- Custom Audience ---

func TestNew_CustomAudience(t *testing.T) {
	tokenServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)
		assert.Contains(t, bodyStr, "audience=https%3A%2F%2Fcustom-api.example.com%2F")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"tok","token_type":"Bearer","expires_in":3600}`))
	}))
	defer tokenServer.Close()

	tokenHost := strings.TrimPrefix(tokenServer.URL, "https://")

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer apiServer.Close()

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tokenServer.Client())

	c, err := New("acme",
		option.WithClientCredentialsAndAudience(ctx, tokenHost, "id", "secret", "https://custom-api.example.com/"),
		option.WithBaseURL(apiServer.URL),
	)
	require.NoError(t, err)

	req, _ := http.NewRequest("GET", apiServer.URL+"/test", nil)
	resp, err := c.options.HTTPClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
}

// --- Max Attempts ---

func TestNew_WithMaxAttempts(t *testing.T) {
	c, err := New("acme",
		option.WithStaticToken("tok"),
		option.WithMaxAttempts(5),
	)
	require.NoError(t, err)
	assert.NotNil(t, c)
}

// --- Tenant Domain Variants Are Accepted ---

// TestNew_DomainVariants verifies that a tenant domain passed to the
// client-credentials option is accepted regardless of protocol prefix or
// trailing slash. The base URL always derives from the team slug and is
// unaffected by the tenant domain.
func TestNew_DomainVariants(t *testing.T) {
	variants := []string{
		"mytenant.auth0.com",
		"https://mytenant.auth0.com",
		"https://mytenant.auth0.com/",
	}
	for _, domain := range variants {
		t.Run(domain, func(t *testing.T) {
			c, err := New("acme", option.WithClientCredentials(context.Background(), domain, "id", "secret"))
			require.NoError(t, err)
			assert.Equal(t, "https://acme.teams.auth0.com", c.baseURL)
		})
	}
}

// --- Debug Mode ---

func TestNew_WithDebug(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.LstdFlags)

	c, err := New("acme",
		option.WithStaticToken("my-secret-token"),
		option.WithBaseURL(server.URL),
		option.WithDebug(true),
	)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	require.NoError(t, err)
	// Simulate the Authorization header that the generated Fern caller adds.
	req.Header.Set("Authorization", "Bearer my-secret-token")

	resp, err := c.options.HTTPClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	output := buf.String()

	// Debug output should contain request and response sections.
	assert.Contains(t, output, "---[ REQUEST ]---")
	assert.Contains(t, output, "---[ RESPONSE ]---")
	assert.Contains(t, output, "GET /test")

	// The token set via WithStaticToken must be redacted, not logged in plain text.
	assert.NotContains(t, output, "my-secret-token")
	assert.Contains(t, output, "[REDACTED]")

	// The transport chain adds User-Agent and Auth0-Client, which should be visible.
	assert.Contains(t, output, "User-Agent: "+teams.SDKName+"/"+teams.SDKVersion)
}

func TestNew_WithDebugDisabled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	defer log.SetFlags(log.LstdFlags)

	c, err := New("acme",
		option.WithStaticToken("my-token"),
		option.WithBaseURL(server.URL),
		option.WithDebug(false),
	)
	require.NoError(t, err)

	req, err := http.NewRequest("GET", server.URL+"/test", nil)
	require.NoError(t, err)

	resp, err := c.options.HTTPClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	// No debug output should be produced.
	assert.Empty(t, buf.String())
}

// --- End-to-end: single /api prefix through the generated client ---

// TestE2E_SingleAPIPrefix is the critical regression guard for the base-URL/
// audience split. An end-to-end call through the generated Tenants client must
// land on .../api/tenants, NOT .../api/api/tenants. It also verifies the
// transport-chain and application-layer headers flow through correctly.
func TestE2E_SingleAPIPrefix(t *testing.T) {
	var capturedPath string
	var capturedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		capturedHeaders = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"tenants":[]}`))
	}))
	defer server.Close()

	c, err := New("acme",
		option.WithStaticToken("my-secret-token"),
		option.WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	// Call through the actual generated client.
	_, _ = c.Tenants.List(context.Background(), &teams.ListTenantsRequestParameters{})

	// The endpoint must resolve to a single "/api" prefix.
	assert.Equal(t, "/api/tenants", capturedPath)
	assert.NotContains(t, capturedPath, "/api/api")

	// Transport-layer headers (set by RoundTripper chain).
	assert.Equal(t, teams.SDKName+"/"+teams.SDKVersion, capturedHeaders.Get("User-Agent"))
	assert.NotEmpty(t, capturedHeaders.Get("Auth0-Client"))

	// Application-layer headers (set by generated code via ToHeader → Caller).
	assert.Equal(t, "Bearer my-secret-token", capturedHeaders.Get("Authorization"))

	// X-Fern headers should NOT be present.
	assert.Empty(t, capturedHeaders.Get("X-Fern-Language"))
	assert.Empty(t, capturedHeaders.Get("X-Fern-SDK-Name"))
	assert.Empty(t, capturedHeaders.Get("X-Fern-SDK-Version"))
}
