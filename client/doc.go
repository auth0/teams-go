// Package client provides the top-level [Teams] client for the Auth0 Teams API.
//
// # Creating a Client
//
// Use [New] to create a client with a team slug, automatic token management,
// and telemetry. For the client-credentials flow, the Auth0 tenant domain is
// passed to the option right after the context:
//
//	c, err := client.New(
//	    "acme", // team slug
//	    option.WithClientCredentials(ctx, "acme-auth.us.auth0.com", clientID, clientSecret),
//	)
//
// # Resource Access
//
// The client exposes API resources through nested sub-clients:
//
//	c.Members              // Member management
//	c.Members.Tenants      // Tenant membership for a member
//	c.Tenants              // Tenant management
//	c.Tenants.Members      // Member management scoped to a tenant
//	c.Environments         // Tenant environment management
//	c.Activity             // Team activity logs
//
// Each sub-client also exposes a WithRawResponse field for accessing full
// HTTP response details (status code, headers, and body).
package client
