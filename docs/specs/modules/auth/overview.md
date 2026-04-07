# Auth Module

## Purpose

The Auth module handles authentication, authorization, and user management for the application.

## Responsibilities

- User management (create, update, disable, enable)
- Authentication (login, token refresh, logout)
- Session management (create, revoke, cleanup)
- Role-Based Access Control (RBAC)
- Permission management (role permissions, direct user permissions)

## Domain Main Entities

| Entity           | Description                                            |
| ---------------- | ------------------------------------------------------ |
| `User`           | User account (admins are users with admin permissions) |
| `Session`        | Active authentication session with tokens              |
| `Role`           | Named collection of permissions                        |
| `RolePermission` | Permission assigned to a role                          |
| `UserRole`       | Role assigned to a user                                |
| `UserPermission` | Direct permission assigned to a user                   |

See ERD.md for entity relationships.
