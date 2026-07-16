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

<details><summary><code>client.Members.Invite(request) -> *teams.InviteMembersResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint invites a new member to the team and optionally grants them access to one or more tenants with specified roles
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
request := &teams.InviteMembersRequestContent{
        Email: "jane.smith@company.com",
        Role: teams.TeamMemberRoleEnumTeamsContributor,
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
client.Members.Invite(
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

**email:** `string` — The email of the team member
    
</dd>
</dl>

<dl>
<dd>

**role:** `*teams.TeamMemberRoleEnum` 
    
</dd>
</dl>

<dl>
<dd>

**tenantIDs:** `[]teams.TenantID` — Specifies an optional list of tenant IDs and their corresponding role assignments.
    
</dd>
</dl>

<dl>
<dd>

**tenantRoles:** `[]*teams.TenantMemberRoleEnum` — The list of roles for each tenant.
    
</dd>
</dl>

<dl>
<dd>

**clientIDs:** `[]teams.ClientID` — The list of apps the user has access to for the 'editor-specific-apps' role. This property is only present for this role.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Members.Get(ID) -> *teams.GetMembersResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the team member's identity and team-level role. Tenant access is retrieved via /api/members/{id}/tenants/{environment}.
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
        "id",
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

**id:** `string` — The unique ID of the team member to retrieve.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Members.Delete(ID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint deletes a team member from the team.
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
client.Members.Delete(
        context.TODO(),
        "id",
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

**id:** `string` — The unique ID of the team member to delete.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Members.Update(ID, request) -> *teams.UpdateMembersResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint can be used to change the team-level role for the team member
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
request := &teams.UpdateMembersRequestContent{
        Role: teams.TeamMemberRoleEnumTeamsOwner,
    }
client.Members.Update(
        context.TODO(),
        "id",
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

**id:** `string` — The unique ID of the team member to update
    
</dd>
</dl>

<dl>
<dd>

**role:** `*teams.TeamMemberRoleEnum` — The new role to assign to the team member.
    
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

Delete a tenant from the current team
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

## Activity
<details><summary><code>client.Activity.List() -> *teams.ListActivityResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint returns a paginated list of activity log entries for the team.
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
request := &teams.ListActivityRequestParameters{
        Take: teams.Int(
            1,
        ),
        From: teams.String(
            "from",
        ),
        Since: teams.Time(
            teams.MustParseDateTime(
                "2024-01-15T09:30:00Z",
            ),
        ),
        Until: teams.Time(
            teams.MustParseDateTime(
                "2024-01-15T09:30:00Z",
            ),
        ),
        Type: teams.ActivityLogEventTypeEnumTeamMember.Ptr(),
        Status: teams.ActivityLogStatusEnumSuccess.Ptr(),
    }
client.Activity.List(
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

**since:** `*teams.DateTimeIso` — Return logs created at or after this ISO 8601 date format
    
</dd>
</dl>

<dl>
<dd>

**until:** `*teams.DateTimeIso` — Return logs created before this ISO 8601 timestamp
    
</dd>
</dl>

<dl>
<dd>

**type_:** `*teams.ActivityLogEventTypeEnum` — Exact match: filter by event type
    
</dd>
</dl>

<dl>
<dd>

**status:** `*teams.ActivityLogStatusEnum` — Exact match: filter by status
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Members Tenants
<details><summary><code>client.Members.Tenants.Create(ID, request) -> *teams.CreateMembersTenantsResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint can be used to add new tenants access
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
request := &teams.CreateMembersTenantsRequestContent{
        Tenants: []teams.TenantID{
            "538c9e21-e3d5-4ad6-b3d0-352c62369fb0",
            "b1c9e21e-1234-4ad6-b3d0-352c6236abcd",
        },
        Roles: []teams.TenantMemberRoleEnum{
            teams.TenantMemberRoleEnumOwner,
            teams.TenantMemberRoleEnumEditorSpecificApps,
        },
        ClientIDs: []teams.ClientID{
            "client_123",
            "client_456",
        },
    }
client.Members.Tenants.Create(
        context.TODO(),
        "id",
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

**id:** `string` — The unique ID of the team member to update
    
</dd>
</dl>

<dl>
<dd>

**tenants:** `[]teams.TenantID` — A list of tenant ids
    
</dd>
</dl>

<dl>
<dd>

**roles:** `[]*teams.TenantMemberRoleEnum` — List of roles to be assigned to the tenant member
    
</dd>
</dl>

<dl>
<dd>

**clientIDs:** `[]teams.ClientID` — List of applications. Only if nested role is editor-specific-app
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Members.Tenants.Remove(ID, request) -> *teams.DeleteMembersTenantsResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint can be used to remove tenant access
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
request := &teams.DeleteMembersTenantsRequestContent{
        Tenants: []teams.TenantID{
            "tenants",
        },
    }
client.Members.Tenants.Remove(
        context.TODO(),
        "id",
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

**id:** `string` — The unique ID of the team member.
    
</dd>
</dl>

<dl>
<dd>

**tenants:** `[]teams.TenantID` — Specifies list of tenant IDs to be removed
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Members.Tenants.Update(ID, request) -> *teams.UpdateMembersTenantsResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint can be used to update existing tenant roles.
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
request := &teams.UpdateMembersTenantsRequestContent{
        Tenants: []teams.TenantID{
            "538c9e21-e3d5-4ad6-b3d0-352c62369fb0",
        },
        Roles: []teams.TenantMemberRoleEnum{
            teams.TenantMemberRoleEnumOwner,
            teams.TenantMemberRoleEnumEditorSpecificApps,
        },
        ClientIDs: []teams.ClientID{
            "client_123",
            "client_456",
        },
    }
client.Members.Tenants.Update(
        context.TODO(),
        "id",
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

**id:** `string` — The unique ID of the team member to update
    
</dd>
</dl>

<dl>
<dd>

**tenants:** `[]teams.TenantID` — A list of tenant ids
    
</dd>
</dl>

<dl>
<dd>

**roles:** `[]*teams.TenantMemberRoleEnum` — List of roles to be assigned to the tenant member
    
</dd>
</dl>

<dl>
<dd>

**clientIDs:** `[]teams.ClientID` — List of applications. Only if nested role is editor-specific-app
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Members.Tenants.List(ID, Environment) -> *teams.ListMembersTenantsResponseContent</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

This endpoint retrieves a list of all tenants within a specific environment that the team member has access to. Environment identifiers required for this request can be obtained using the /api/environments endpoint, which provides details for all environments where tenants associated with your team exists.
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
        "id",
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

**id:** `string` — Unique ID of the team member.
    
</dd>
</dl>

<dl>
<dd>

**environment:** `string` — The cloud environment identifier (e.g. 'US-3', 'EU-1', 'acme-dev').
    
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

