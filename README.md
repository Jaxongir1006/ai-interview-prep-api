# [Project Name]

<!-- PROJECT-SPECIFIC: Replace [Project Name] above with your project name -->

> Built on [Go Enterprise Blueprint](https://github.com/rise-and-shine/go-enterprise-blueprint)

## Overview

<!-- PROJECT-SPECIFIC: Describe your project here -->

[Brief description of what your project does, its purpose, mission, and key features.]

- [Architecture](docs/architecture/)
- [Guidelines](docs/guidelines/)
- [Specs](docs/specs/)

## Use Case Index

| Use Case (module:operation-id) | Type           | Permissions                    | Docs                                                                                |
| ------------------------------ | -------------- | ------------------------------ | ----------------------------------------------------------------------------------- |
| `auth:admin-login`             | user_action    | -                              | [spec](docs/specs/modules/auth/usecases/user/admin-login.md)                        |
| `auth:login`                   | user_action    | -                              | [spec](docs/specs/modules/auth/usecases/user/login.md)                              |
| `auth:register`                | user_action    | -                              | [spec](docs/specs/modules/auth/usecases/user/register.md)                           |
| `auth:verify-email`            | user_action    | -                              | [spec](docs/specs/modules/auth/usecases/user/verify-email.md)                       |
| `auth:resend-verification-email` | user_action   | -                              | [spec](docs/specs/modules/auth/usecases/user/resend-verification-email.md)          |
| `auth:refresh-token`           | user_action    | -                              | [spec](docs/specs/modules/auth/usecases/user/refresh-token.md)                      |
| `auth:logout`                  | user_action    | -                              | [spec](docs/specs/modules/auth/usecases/user/logout.md)                             |
| `auth:get-my-sessions`         | user_action    | -                              | [spec](docs/specs/modules/auth/usecases/user/get-my-sessions.md)                    |
| `auth:delete-my-session`       | user_action    | -                              | [spec](docs/specs/modules/auth/usecases/user/delete-my-session.md)                  |
| `auth:get-auth-stats`          | user_action    | `auth:user:read`               | [spec](docs/specs/modules/auth/usecases/user/get-auth-stats.md)                     |
| `auth:get-users`               | user_action    | `auth:user:read`               | [spec](docs/specs/modules/auth/usecases/user/get-users.md)                          |
| `auth:create-user`             | user_action    | `auth:user:manage`             | [spec](docs/specs/modules/auth/usecases/user/create-user.md)                        |
| `auth:update-user`             | user_action    | `auth:user:manage`             | [spec](docs/specs/modules/auth/usecases/user/update-user.md)                        |
| `auth:disable-user`            | user_action    | `auth:user:manage`             | [spec](docs/specs/modules/auth/usecases/user/disable-user.md)                       |
| `auth:enable-user`             | user_action    | `auth:user:manage`             | [spec](docs/specs/modules/auth/usecases/user/enable-user.md)                        |
| `auth:create-role`             | user_action    | `auth:role:manage`             | [spec](docs/specs/modules/auth/usecases/rbac/create-role.md)                        |
| `auth:update-role`             | user_action    | `auth:role:manage`             | [spec](docs/specs/modules/auth/usecases/rbac/update-role.md)                        |
| `auth:delete-role`             | user_action    | `auth:role:manage`             | [spec](docs/specs/modules/auth/usecases/rbac/delete-role.md)                        |
| `auth:get-roles`               | user_action    | `auth:role:read`               | [spec](docs/specs/modules/auth/usecases/rbac/get-roles.md)                          |
| `auth:set-role-permissions`    | user_action    | `auth:role:manage`             | [spec](docs/specs/modules/auth/usecases/rbac/set-role-permissions.md)               |
| `auth:get-role-permissions`    | user_action    | `auth:role:read`               | [spec](docs/specs/modules/auth/usecases/rbac/get-role-permissions.md)               |
| `auth:set-user-roles`          | user_action    | `auth:access:manage`           | [spec](docs/specs/modules/auth/usecases/rbac/set-user-roles.md)                     |
| `auth:get-user-roles`          | user_action    | `auth:access:read`             | [spec](docs/specs/modules/auth/usecases/rbac/get-user-roles.md)                     |
| `auth:set-user-permissions`    | user_action    | `auth:access:manage`           | [spec](docs/specs/modules/auth/usecases/rbac/set-user-permissions.md)               |
| `auth:get-user-permissions`    | user_action    | `auth:access:read`             | [spec](docs/specs/modules/auth/usecases/rbac/get-user-permissions.md)               |
| `auth:get-user-sessions`       | user_action    | `auth:session:read`            | [spec](docs/specs/modules/auth/usecases/session/get-user-sessions.md)               |
| `auth:delete-session`          | user_action    | `auth:session:manage`          | [spec](docs/specs/modules/auth/usecases/session/delete-session.md)                  |
| `auth:delete-user-sessions`    | user_action    | `auth:session:manage`          | [spec](docs/specs/modules/auth/usecases/session/delete-user-sessions.md)            |
| `auth:create-superadmin`       | manual_command | -                              | [spec](docs/specs/modules/auth/usecases/user/create-superadmin.md)                  |
| `auth:clean-expired-sessions`  | async_task     | -                              | [spec](docs/specs/modules/auth/usecases/session/clean-expired-sessions.md)          |
| `audit:get-action-logs`        | user_action    | `audit:action-log:read`        | [spec](docs/specs/modules/audit/usecases/actionlog/get-action-logs.md)              |
| `audit:get-status-change-logs` | user_action    | `audit:status-change-log:read` | [spec](docs/specs/modules/audit/usecases/statuschangelog/get-status-change-logs.md) |
| `candidate:complete-onboarding` | user_action   | -                              | [spec](docs/specs/modules/candidate/usecases/profile/complete-onboarding.md)        |
| `platform:list-queues`         | user_action    | `taskmill:view`                | [spec](docs/specs/modules/platform/usecases/taskmill/list-queues.md)                |
| `platform:get-queue-stats`     | user_action    | `taskmill:view`                | [spec](docs/specs/modules/platform/usecases/taskmill/get-queue-stats.md)            |
| `platform:list-dlq-tasks`      | user_action    | `taskmill:view`                | [spec](docs/specs/modules/platform/usecases/taskmill/list-dlq-tasks.md)             |
| `platform:list-task-results`   | user_action    | `taskmill:view`                | [spec](docs/specs/modules/platform/usecases/taskmill/list-task-results.md)          |
| `platform:list-schedules`      | user_action    | `taskmill:view`                | [spec](docs/specs/modules/platform/usecases/taskmill/list-schedules.md)             |
| `platform:requeue-from-dlq`    | user_action    | `taskmill:manage`              | [spec](docs/specs/modules/platform/usecases/taskmill/requeue-from-dlq.md)           |
| `platform:purge-queue`         | user_action    | `taskmill:manage`              | [spec](docs/specs/modules/platform/usecases/taskmill/purge-queue.md)                |
| `platform:purge-dlq`           | user_action    | `taskmill:manage`              | [spec](docs/specs/modules/platform/usecases/taskmill/purge-dlq.md)                  |
| `platform:cleanup-results`     | user_action    | `taskmill:manage`              | [spec](docs/specs/modules/platform/usecases/taskmill/cleanup-results.md)            |
| `platform:trigger-schedule`    | user_action    | `taskmill:manage`              | [spec](docs/specs/modules/platform/usecases/taskmill/trigger-schedule.md)           |
| `platform:list-errors`         | user_action    | `alert:view`                   | [spec](docs/specs/modules/platform/usecases/alerterror/list-errors.md)              |
| `platform:get-error`           | user_action    | `alert:view`                   | [spec](docs/specs/modules/platform/usecases/alerterror/get-error.md)                |
| `platform:get-error-stats`     | user_action    | `alert:view`                   | [spec](docs/specs/modules/platform/usecases/alerterror/get-error-stats.md)          |
| `platform:cleanup-errors`      | user_action    | `alert:manage`                 | [spec](docs/specs/modules/platform/usecases/alerterror/cleanup-errors.md)           |

## How to Run

### Production

<!-- PROJECT-SPECIFIC: Build, Deploy, Run and other CLI commands for devops team -->

```bash
# Build
# TODO: ...

# Run
# TODO: ...
```

### Development

**Prerequisites:**

- Docker
- Go 1.25+
- OS: macOS or Linux (some features won't work on Windows)

**Run locally:**

```bash
make run       # Start infrastructure (database, etc.) and run the application
```

**Run tests:**

```bash
make test         # Run unit tests
make test-system  # Run system tests
```
