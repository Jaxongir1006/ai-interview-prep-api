# Analytics Module ERD

```mermaid
erDiagram
    achievement_definitions ||--o{ candidate_achievements : "awarded as"

    candidate_progress_summaries {
        BIGSERIAL id PK
        VARCHAR user_id FK "references auth.users UUID-formatted string identifier"
        INT current_streak "default 0"
        INT longest_streak "default 0"
        INT total_interviews_taken "default 0"
        BIGINT total_time_spent_seconds "default 0"
        NUMERIC_5_2 average_score "0-100 aggregate score"
        TIMESTAMPTZ last_interview_at "nullable"
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    candidate_topic_stats {
        BIGSERIAL id PK
        VARCHAR user_id FK "references auth.users UUID-formatted string identifier"
        VARCHAR topic_key "stable topic identifier, e.g. golang-concurrency"
        INT attempts "default 0"
        BIGINT total_time_spent_seconds "default 0"
        NUMERIC_5_2 average_score "0-100 aggregate score"
        NUMERIC_5_2 best_score "0-100 aggregate score"
        TIMESTAMPTZ last_practiced_at "nullable"
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    achievement_definitions {
        BIGSERIAL id PK
        VARCHAR code UK "stable identifier, e.g. first-interview"
        VARCHAR name
        TEXT description
        VARCHAR category
        INT sort_order "default 0"
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    candidate_achievements {
        BIGSERIAL id PK
        VARCHAR user_id FK "references auth.users UUID-formatted string identifier"
        BIGINT achievement_definition_id FK
        TIMESTAMPTZ achieved_at
        JSONB metadata "nullable"
        TIMESTAMPTZ created_at
    }
```

## Schema

The analytics tables reside in the `analytics` schema.

## Notes

- `candidate_progress_summaries.user_id` is unique, enforcing one summary row per candidate
- `candidate_topic_stats` should be unique per `user_id + topic_key`
- `achievement_definitions.code` is unique and stable for business logic and API references
- `candidate_achievements` should be unique per `user_id + achievement_definition_id`
- `current_streak`, `longest_streak`, `total_interviews_taken`, and `total_time_spent_seconds` are constrained to be non-negative
- `average_score` and `best_score` are constrained to the range `0..100`
- Strengths and weaknesses should be derived from `candidate_topic_stats` at read time or through dedicated analytics queries, not stored as standalone mutable fields
- When implementing the migration, follow [Migration Guideline](../../../guidelines/13_db_migrations.md): create tables first, create indexes second, then add foreign keys and check constraints with `ALTER TABLE`
