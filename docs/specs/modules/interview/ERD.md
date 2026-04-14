# Interview Module ERD

```mermaid
erDiagram
    interview_sessions ||--o{ interview_questions : "contains"
    interview_questions ||--o| interview_answers : "answered by"
    interview_questions ||--o| interview_reviews : "reviewed by"
    interview_answers ||--o| interview_reviews : "review source"

    interview_sessions {
        VARCHAR id PK "UUID-formatted session identifier"
        VARCHAR user_id FK "references auth.users UUID-formatted string identifier"
        VARCHAR title "display title"
        VARCHAR target_role "role selected for this session"
        VARCHAR experience_level "junior, mid, senior"
        VARCHAR difficulty "easy, medium, hard, mixed"
        VARCHAR status "in_progress, completed, abandoned, scoring"
        INT question_count "number of questions planned for the session"
        INT answered_count "number of submitted answers"
        NUMERIC_5_2 total_score "nullable, final 0-100 aggregate score"
        BIGINT total_duration_seconds "default 0"
        TIMESTAMPTZ started_at
        TIMESTAMPTZ completed_at "nullable"
        TIMESTAMPTZ abandoned_at "nullable"
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    interview_questions {
        BIGSERIAL id PK
        VARCHAR session_id FK
        VARCHAR topic_key "stable topic identifier"
        VARCHAR topic_name "display topic name shown at the time"
        VARCHAR difficulty "easy, medium, hard"
        VARCHAR question_type "technical, behavioral, system_design, coding"
        TEXT question_text "exact question shown to the candidate"
        TEXT expected_answer "nullable, internal reference answer"
        JSONB evaluation_rubric "nullable, rubric used for review"
        VARCHAR source "ai, manual, catalog"
        VARCHAR source_question_id "nullable, external/catalog source identifier"
        VARCHAR ai_provider "nullable"
        VARCHAR ai_model "nullable"
        VARCHAR prompt_version "nullable"
        INT position "1-based order inside the session"
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    interview_answers {
        BIGSERIAL id PK
        BIGINT session_question_id FK
        TEXT answer_text "candidate answer exactly as submitted"
        BIGINT time_spent_seconds "default 0"
        TIMESTAMPTZ submitted_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    interview_reviews {
        BIGSERIAL id PK
        BIGINT session_question_id FK
        BIGINT answer_id FK "nullable"
        VARCHAR reviewer_type "ai or manual"
        NUMERIC_5_2 score "nullable, 0-100"
        NUMERIC_5_2 correctness_rate "nullable, 0-1"
        TEXT feedback "nullable"
        JSONB rubric_scores "nullable"
        JSONB strengths "nullable"
        JSONB improvements "nullable"
        VARCHAR ai_provider "nullable"
        VARCHAR ai_model "nullable"
        VARCHAR prompt_version "nullable"
        JSONB metadata "nullable, internal review metadata"
        TIMESTAMPTZ reviewed_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }
```

## Schema

The interview tables reside in the `interview` schema.

## Notes

- `interview_sessions.user_id` references `auth.users(id)` conceptually; module code should still communicate across modules through portals.
- `interview_sessions.status` is constrained to `in_progress`, `completed`, `abandoned`, or `scoring`.
- `interview_sessions.difficulty` is constrained to `easy`, `medium`, `hard`, or `mixed`.
- `interview_sessions.experience_level` is constrained to `junior`, `mid`, or `senior`.
- `interview_questions.difficulty` is constrained to `easy`, `medium`, or `hard`.
- `interview_questions.source` is constrained to `ai`, `manual`, or `catalog`.
- `interview_reviews.reviewer_type` is constrained to `ai` or `manual`.
- `interview_questions` should be unique by `session_id + position`.
- `interview_answers` should allow one current answer per session question.
- `interview_reviews` should allow one current review per session question.
- Scores are constrained to `0..100`; correctness rate is constrained to `0..1`.
- Durations and counts are constrained to be non-negative.
- When implementing the migration, follow [Migration Guideline](../../../guidelines/13_db_migrations.md): create tables first, create indexes second, then add foreign keys and check constraints with `ALTER TABLE`.
