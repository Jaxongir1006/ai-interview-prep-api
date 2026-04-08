# Auth Module

## Purpose

The Auth module handles authentication, authorization, and user identity management for the application.

## Responsibilities

- User identity management for both admin and public users
- Authentication (admin login, public user login, Google OAuth login, GitHub OAuth login, token refresh, logout)
- Session management (create, revoke, cleanup)
- Role-Based Access Control (RBAC)
- Permission management (role permissions, direct user permissions)
- Verification and auth contact data management (`username`, `email`, `phone_number`, `is_verified`)
- External identity linking for OAuth providers

## Domain Main Entities

| Entity           | Description                                                                 |
| ---------------- | --------------------------------------------------------------------------- |
| `User`           | Shared auth identity for admins and public users                            |
| `OAuthAccount`   | External OAuth identity linked to a user account                            |
| `Session`        | Active authentication session with tokens                                   |
| `Role`           | Named collection of permissions                                             |
| `RolePermission` | Permission assigned to a role                                               |
| `UserRole`       | Role assigned to a user                                                     |
| `UserPermission` | Direct permission assigned to a user                                        |

## Identity Model

- Admin users authenticate with `username + password`
- Public users authenticate with `email + password`
- Public users may also authenticate with Google OAuth or GitHub OAuth
- Both account types live in the same `auth.users` table
- Interview-preparation profile data belongs in a separate business module, not in `auth`

See ERD.md for entity relationships.
