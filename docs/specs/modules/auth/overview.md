# Auth Module

## Purpose

The Auth module handles authentication, authorization, and user identity management for the application.

## Responsibilities

- User identity management for both admin and public users
- Authentication (admin login, public user login, Google OAuth login, GitHub OAuth login, token refresh, logout)
- Session management (create, revoke, cleanup)
- Password reset for public email/password users
- Role-Based Access Control (RBAC)
- Permission management (role permissions, direct user permissions)
- Verification and auth contact data management (`username`, `email`, `phone_number`, `is_verified`)
- Email verification for password-based public registration
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
- Public users created through email/password registration start with `is_verified = false`
- Public users created or linked through OAuth use the provider's verified email signal and do not receive an application email-verification message
- Interview-preparation profile data belongs in a separate business module, not in `auth`

## Email Delivery

- Auth sends email-verification messages for password-based public registration
- Local development should use SMTP capture tooling such as Mailpit
- Staging and production should use a configured transactional email provider through SMTP or provider-specific infrastructure
- Verification links should point to the frontend verification page and include the raw one-time token as a query parameter
- Password reset links should point to the frontend password reset page and include the raw one-time token as a query parameter

## Configuration

- `auth.frontend_verify_email_url` defines the frontend email verification URL
- `auth.frontend_password_reset_url` defines the frontend password reset URL
- `auth.email_verification_token_ttl` defines email verification token lifetime
- `auth.password_reset_token_ttl` defines password reset token lifetime
- Auth requires Redis for ephemeral email verification and password reset token state

See ERD.md for entity relationships.
