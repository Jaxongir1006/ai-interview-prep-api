# Platform Module ERD

The platform module does not own any domain entities. It operates on Taskmill infrastructure tables (`taskmill.task_queue`, `taskmill.task_results`, `taskmill.task_schedules`) via the `console` package.

The `errors` table is managed by `rise-and-shine/pkg/observability/alert`. Platform module has read-only + delete access.

```mermaid
erDiagram
    errors {
        VARCHAR id PK
        TEXT code "error code (e.g. INTERNAL_ERROR)"
        TEXT message "error message"
        JSONB details "additional context (trace_id, actor_id, etc.)"
        TEXT service "service name"
        TEXT operation "operation where error occurred"
        TIMESTAMPTZ created_at
        BOOLEAN alerted "whether notification was sent"
    }
```
