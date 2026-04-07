# Audit Module

## Purpose

Provides a centralized abstraction for other modules to store and manage audit logs of the system. Acts as the single entry point for audit logging across the application.

## Responsibilities

- Provide an abstraction layer (via portal) for other modules to record audit logs
- Store and manage user action logs (who did what, when, and where)
- Store and manage object status change logs (what changed, from which state to which state)
- Ensure consistent audit log format across all modules

## Domain Main Entities

| Entity              | Description                                         |
| ------------------- | --------------------------------------------------- |
| `action_log`        | Records who performed what action, when, and where  |
| `status_change_log` | Records entity state transitions tied to an action  |

See ERD.md for entity relationships.
