// Package teams provides a Go client for the Auth0 Teams API.
//
// This SDK handles authentication, automatic token management, retries with
// exponential backoff, and telemetry out of the box.
//
// # Client Initialization
//
// Create a client using [client.New] with your team slug and one of the
// supported authentication options. For the client-credentials flow, the
// Auth0 tenant domain is passed to the option right after the context:
//
//	c, err := client.New(
//	    "acme", // team slug
//	    option.WithClientCredentials(ctx, "acme-auth.us.auth0.com", clientID, clientSecret),
//	)
//
// Three authentication modes are supported (mutually exclusive):
//
//   - Client Credentials — [option.WithClientCredentials]
//   - Custom Token Source — [option.WithTokenSource]
//   - Static Token — [option.WithStaticToken]
//
// # Resource Hierarchy
//
// The top-level client exposes API resources through a nested structure:
//
//	c.Members              // List, invite, get, update, and remove members
//	c.Members.Tenants      // Tenant membership for a member
//	c.Tenants              // List, create, get, update, and delete tenants
//	c.Tenants.Members      // Member management scoped to a tenant
//	c.Environments         // Manage tenant environments
//	c.Activity             // Read team activity logs
//
// # Error Handling
//
// All API methods return typed errors that can be inspected with [errors.As]:
//
//	var notFound *teams.NotFoundError
//	if errors.As(err, &notFound) {
//	    fmt.Println(notFound.Body)
//	}
//
// Available error types: [BadRequestError], [UnauthorizedError],
// [ForbiddenError], [NotFoundError], [ConflictError], [TooManyRequestsError],
// [InternalServerError].
//
// # Pointer Helpers
//
// Helper functions like [String], [Bool], [Int], and [Time] create pointers
// to primitive values for use with optional request fields.
package teams
