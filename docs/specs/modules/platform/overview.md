# Platform Module

## Purpose

The platform module serves as the operational backbone of the system. It does not implement business features — instead, it exposes internal infrastructure APIs and background jobs that support platform administration, monitoring, and maintenance.

## Responsibilities

- Taskmill Console APIs — view queue stats, DLQ tasks, task results, schedules; manage retries, purges, and manual triggers
- Platform housekeeping — background cleanup tasks for old task results and queue entries
- Error monitoring — browse application errors, view details, get statistics, cleanup old records
- Error catalog — expose registered error codes across all modules (planned)
- Health and diagnostics — readiness/liveness probes (planned)

## Domain Main Entities

| Entity | Description |
| ------ | ----------- |

This module has no domain entities of its own. It operates on infrastructure-level data (Taskmill tables) via the `console` package.

See ERD.md for entity relationships.
