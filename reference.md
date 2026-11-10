# Reference
## Members
<details><summary><code>client.Members.List() -> *teams.ListMembersResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint retrieves a paginated list of all members under the current team. It supports checkpoint-based pagination for efficient retrieval of large datasets
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &teams.ListMembersRequestParameters{
        Take: teams.Int(
            1,
        ),
        From: teams.String(
            "from",
        ),
    }
client.Members.List(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**take:** `*int` — Maximum number of items to return per page. Defaults to 50. Maximum 50
    
</dd>
</dl>

<dl>
<dd>

**from:** `*string` — Checkpoint cursor for pagination; use the 'next' value from a previous response to fetch the next page
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Members.Get(MemberID) -> *teams.GetMembersResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the team member's identity and team-level role. Tenant access is retrieved via /api/v1/members/{member_id}/tenants/{environment}.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Members.Get(
        context.TODO(),
        "member_id",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**memberID:** `string` — The unique ID of the team member to retrieve.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Tenants
<details><summary><code>client.Tenants.List() -> *teams.ListTenantsResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint retrieves a paginated list of tenants for a specified team. It supports cursor-based pagination to ensure efficient data retrieval across large datasets. If the private cloud support beta is enabled for the team, the response includes both public and private tenants; otherwise, only public tenants are returned. To navigate results, the next cursor from the response should be used to fetch subsequent pages.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &teams.ListTenantsRequestParameters{
        Take: teams.Int(
            1,
        ),
        From: teams.String(
            "from",
        ),
        Sort: teams.String(
            "sort",
        ),
        Environment: teams.String(
            "US-3",
        ),
        EnvironmentTag: teams.EnvironmentTypeEnumDevelopment.Ptr(),
        Locality: teams.String(
            "us",
        ),
        TenantName: teams.String(
            "acme-development-tenant",
        ),
    }
client.Tenants.List(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**take:** `*int` — Maximum number of items to return per page. Defaults to 50. Maximum 50
    
</dd>
</dl>

<dl>
<dd>

**from:** `*string` — Checkpoint cursor for pagination; use the 'next' value from a previous response to fetch the next page
    
</dd>
</dl>

<dl>
<dd>

**sort:** `*string` — Sort order for results using MongoDB-style syntax. Defaults to 'created_at:-1' (newest first, descending). Use 'created_at:1' for oldest first (ascending). Format: field:direction where direction is 1 (ascending) or -1 (descending).
    
</dd>
</dl>

<dl>
<dd>

**environment:** `*string` — Filter tenants by environment. For public cloud tenants, use the short name (e.g., 'US-3', 'EU-1'). For private cloud tenants, use the space name (e.g., 'acme-dev').
    
</dd>
</dl>

<dl>
<dd>

**environmentTag:** `*teams.EnvironmentTypeEnum` — Filter tenants by environment tag
    
</dd>
</dl>

<dl>
<dd>

**locality:** `*string` — Filter tenants by locality. For public cloud tenants, matches locality (e.g., 'us', 'eu'). For private cloud tenants, matches the space's primary locality (e.g., 'virginia', 'frankfurt').
    
</dd>
</dl>

<dl>
<dd>

**tenantName:** `*teams.TenantName` — Filter tenants by name prefix (case-sensitive).
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Tenants.Create(request) -> *teams.CreateTenantsResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint creates a new tenant under your team. Upon successful creation, the response includes the tenant's details and management client credentials. The tenant type (public or private) is determined by the attributes in the request payload. Refer to the payload schema for required attributes.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &teams.CreateTenantsRequestContent{
        PublicCloudPayload: &teams.PublicCloudPayload{
            AdminEmail: "admin_email",
            Locality: "us",
        },
    }
client.Tenants.Create(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**request:** `*teams.CreateTenantsRequestContent` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Tenants.Get(TenantID) -> *teams.GetTenantsResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint retrieves all metadata for a specific tenant, including the tenant's unique ID, name, domain, environment, locality, environment type, and creation date. It is useful for obtaining comprehensive information about a tenant for display in dashboards, auditing, or automated management tasks. The response includes all relevant fields required to identify and manage the tenant programmatically
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Tenants.Get(
        context.TODO(),
        "tenantId",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**tenantID:** `teams.TenantID` — The unique ID of the tenant to retrieve (UUID format)
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Tenants.Delete(TenantID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Delete a tenant from the current team. Tenants tagged as Production cannot be deleted through this API. To protect a tenant from accidental deletion, tag it as Production via Tenant Settings. [Link](https://auth0.com/docs/get-started/tenant-settings#environment-tag).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Tenants.Delete(
        context.TODO(),
        "tenantId",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**tenantID:** `teams.TenantID` — The id of the tenant to delete
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Environments
<details><summary><code>client.Environments.List() -> *teams.ListEnvironmentsResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint provide a list of environments associated with the Team where tenants exist.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Environments.List(
        context.TODO(),
    )
}
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Members Tenants
<details><summary><code>client.Members.Tenants.List(MemberID, Environment) -> *teams.ListMembersTenantsResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint retrieves a list of all tenants within a specific environment that the team member has access to. Environment identifiers required for this request can be obtained using the /api/v1/environments endpoint, which provides details for all environments where tenants associated with your team exists.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Members.Tenants.List(
        context.TODO(),
        "member_id",
        "environment",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**memberID:** `string` — Unique ID of the team member.
    
</dd>
</dl>

<dl>
<dd>

**environment:** `string` — The cloud environment identifier (e.g. 'prod-us-3', 'prod-eu-1', 'acme-dev').
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Tenants Members
<details><summary><code>client.Tenants.Members.List(TenantID) -> *teams.ListTenantsMembersResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint retrieves a paginated list of all members for a specific tenant.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &teams.ListTenantsMembersRequestParameters{
        Take: teams.Int(
            1,
        ),
        From: teams.String(
            "from",
        ),
    }
client.Tenants.Members.List(
        context.TODO(),
        "tenantId",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**tenantID:** `teams.TenantID` — The unique ID of the tenant.
    
</dd>
</dl>

<dl>
<dd>

**take:** `*int` — Maximum number of items to return per page. Defaults to 50. Maximum 50
    
</dd>
</dl>

<dl>
<dd>

**from:** `*string` — Checkpoint cursor for pagination; use the 'next' value from a previous response to fetch the next page
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

