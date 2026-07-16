package core

import (
	"context"

	"golang.org/x/oauth2"
)

// ClientCredentialsOption configures OAuth2 client credentials authentication.
type ClientCredentialsOption struct {
	Context      context.Context
	TenantDomain string
	ClientID     string
	ClientSecret string
	Audience     string
}

func (c *ClientCredentialsOption) applyRequestOptions(opts *RequestOptions) {
	opts.ClientID = c.ClientID
	opts.ClientSecret = c.ClientSecret
	opts.TenantDomain = c.TenantDomain
	opts.OAuthContext = c.Context
	opts.Audience = c.Audience
}

// TokenSourceOption configures a custom oauth2.TokenSource for authentication.
type TokenSourceOption struct {
	TokenSource oauth2.TokenSource
}

func (t *TokenSourceOption) applyRequestOptions(opts *RequestOptions) {
	opts.TokenSource = t.TokenSource
}

// NoAuth0ClientInfoOption disables the Auth0-Client telemetry header.
type NoAuth0ClientInfoOption struct{}

func (n *NoAuth0ClientInfoOption) applyRequestOptions(opts *RequestOptions) {
	opts.NoAuth0ClientInfo = true
}

// Auth0ClientEnvEntryOption adds a custom entry to the Auth0-Client env map.
type Auth0ClientEnvEntryOption struct {
	Key   string
	Value string
}

func (a *Auth0ClientEnvEntryOption) applyRequestOptions(opts *RequestOptions) {
	if opts.Auth0ClientEnv == nil {
		opts.Auth0ClientEnv = make(map[string]string)
	}
	opts.Auth0ClientEnv[a.Key] = a.Value
}

// DebugOption enables HTTP request/response debug logging.
type DebugOption struct {
	Debug bool
}

func (d *DebugOption) applyRequestOptions(opts *RequestOptions) {
	opts.Debug = d.Debug
}
