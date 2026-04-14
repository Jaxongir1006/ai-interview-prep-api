-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS interview;

CREATE TABLE interview.interview_sessions (
    id VARCHAR PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    title VARCHAR(255) NOT NULL,
    target_role VARCHAR(100) NOT NULL,
    experience_level VARCHAR(50) NOT NULL,
    difficulty VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    question_count INTEGER NOT NULL DEFAULT 0,
    answered_count INTEGER NOT NULL DEFAULT 0,
    total_score NUMERIC(5,2),
    total_duration_seconds BIGINT NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    abandoned_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE interview.interview_questions (
    id BIGSERIAL PRIMARY KEY,
    session_id VARCHAR NOT NULL,
    topic_key VARCHAR(100) NOT NULL,
    topic_name VARCHAR(255) NOT NULL,
    difficulty VARCHAR(50) NOT NULL,
    question_type VARCHAR(50) NOT NULL,
    question_text TEXT NOT NULL,
    expected_answer TEXT,
    evaluation_rubric JSONB,
    source VARCHAR(50) NOT NULL,
    source_question_id VARCHAR(255),
    ai_provider VARCHAR(100),
    ai_model VARCHAR(100),
    prompt_version VARCHAR(100),
    position INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE interview.interview_answers (
    id BIGSERIAL PRIMARY KEY,
    session_question_id BIGINT NOT NULL,
    answer_text TEXT NOT NULL,
    time_spent_seconds BIGINT NOT NULL DEFAULT 0,
    submitted_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE interview.interview_reviews (
    id BIGSERIAL PRIMARY KEY,
    session_question_id BIGINT NOT NULL,
    answer_id BIGINT,
    reviewer_type VARCHAR(50) NOT NULL,
    score NUMERIC(5,2),
    correctness_rate NUMERIC(5,4),
    feedback TEXT,
    rubric_scores JSONB,
    strengths JSONB,
    improvements JSONB,
    ai_provider VARCHAR(100),
    ai_model VARCHAR(100),
    prompt_version VARCHAR(100),
    metadata JSONB,
    reviewed_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_interview_sessions_user_id
    ON interview.interview_sessions (user_id);

CREATE INDEX idx_interview_sessions_user_status
    ON interview.interview_sessions (user_id, status);

CREATE INDEX idx_interview_sessions_started_at
    ON interview.interview_sessions (started_at);

CREATE UNIQUE INDEX uk_interview_questions_session_position
    ON interview.interview_questions (session_id, position);

CREATE INDEX idx_interview_questions_session_id
    ON interview.interview_questions (session_id);

CREATE INDEX idx_interview_questions_topic_key
    ON interview.interview_questions (topic_key);

CREATE UNIQUE INDEX uk_interview_answers_session_question_id
    ON interview.interview_answers (session_question_id);

CREATE UNIQUE INDEX uk_interview_reviews_session_question_id
    ON interview.interview_reviews (session_question_id);

CREATE INDEX idx_interview_reviews_answer_id
    ON interview.interview_reviews (answer_id);

ALTER TABLE interview.interview_questions
    ADD CONSTRAINT fk_interview_questions_session
        FOREIGN KEY (session_id)
        REFERENCES interview.interview_sessions(id)
        ON DELETE CASCADE;

ALTER TABLE interview.interview_answers
    ADD CONSTRAINT fk_interview_answers_question
        FOREIGN KEY (session_question_id)
        REFERENCES interview.interview_questions(id)
        ON DELETE CASCADE;

ALTER TABLE interview.interview_reviews
    ADD CONSTRAINT fk_interview_reviews_question
        FOREIGN KEY (session_question_id)
        REFERENCES interview.interview_questions(id)
        ON DELETE CASCADE;

ALTER TABLE interview.interview_reviews
    ADD CONSTRAINT fk_interview_reviews_answer
        FOREIGN KEY (answer_id)
        REFERENCES interview.interview_answers(id)
        ON DELETE SET NULL;

ALTER TABLE interview.interview_sessions
    ADD CONSTRAINT chk_interview_sessions_status
        CHECK (status IN ('in_progress', 'completed', 'abandoned', 'scoring'));

ALTER TABLE interview.interview_sessions
    ADD CONSTRAINT chk_interview_sessions_difficulty
        CHECK (difficulty IN ('easy', 'medium', 'hard', 'mixed'));

ALTER TABLE interview.interview_sessions
    ADD CONSTRAINT chk_interview_sessions_experience_level
        CHECK (experience_level IN ('junior', 'mid', 'senior'));

ALTER TABLE interview.interview_sessions
    ADD CONSTRAINT chk_interview_sessions_question_count
        CHECK (question_count >= 0);

ALTER TABLE interview.interview_sessions
    ADD CONSTRAINT chk_interview_sessions_answered_count
        CHECK (answered_count >= 0);

ALTER TABLE interview.interview_sessions
    ADD CONSTRAINT chk_interview_sessions_total_score
        CHECK (total_score IS NULL OR (total_score >= 0 AND total_score <= 100));

ALTER TABLE interview.interview_sessions
    ADD CONSTRAINT chk_interview_sessions_total_duration
        CHECK (total_duration_seconds >= 0);

ALTER TABLE interview.interview_questions
    ADD CONSTRAINT chk_interview_questions_difficulty
        CHECK (difficulty IN ('easy', 'medium', 'hard'));

ALTER TABLE interview.interview_questions
    ADD CONSTRAINT chk_interview_questions_question_type
        CHECK (question_type IN ('technical', 'behavioral', 'system_design', 'coding'));

ALTER TABLE interview.interview_questions
    ADD CONSTRAINT chk_interview_questions_source
        CHECK (source IN ('ai', 'manual', 'catalog'));

ALTER TABLE interview.interview_questions
    ADD CONSTRAINT chk_interview_questions_position
        CHECK (position > 0);

ALTER TABLE interview.interview_answers
    ADD CONSTRAINT chk_interview_answers_time_spent
        CHECK (time_spent_seconds >= 0);

ALTER TABLE interview.interview_reviews
    ADD CONSTRAINT chk_interview_reviews_reviewer_type
        CHECK (reviewer_type IN ('ai', 'manual'));

ALTER TABLE interview.interview_reviews
    ADD CONSTRAINT chk_interview_reviews_score
        CHECK (score IS NULL OR (score >= 0 AND score <= 100));

ALTER TABLE interview.interview_reviews
    ADD CONSTRAINT chk_interview_reviews_correctness_rate
        CHECK (correctness_rate IS NULL OR (correctness_rate >= 0 AND correctness_rate <= 1));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS interview.interview_reviews
    DROP CONSTRAINT IF EXISTS chk_interview_reviews_correctness_rate;

ALTER TABLE IF EXISTS interview.interview_reviews
    DROP CONSTRAINT IF EXISTS chk_interview_reviews_score;

ALTER TABLE IF EXISTS interview.interview_reviews
    DROP CONSTRAINT IF EXISTS chk_interview_reviews_reviewer_type;

ALTER TABLE IF EXISTS interview.interview_answers
    DROP CONSTRAINT IF EXISTS chk_interview_answers_time_spent;

ALTER TABLE IF EXISTS interview.interview_questions
    DROP CONSTRAINT IF EXISTS chk_interview_questions_position;

ALTER TABLE IF EXISTS interview.interview_questions
    DROP CONSTRAINT IF EXISTS chk_interview_questions_source;

ALTER TABLE IF EXISTS interview.interview_questions
    DROP CONSTRAINT IF EXISTS chk_interview_questions_question_type;

ALTER TABLE IF EXISTS interview.interview_questions
    DROP CONSTRAINT IF EXISTS chk_interview_questions_difficulty;

ALTER TABLE IF EXISTS interview.interview_sessions
    DROP CONSTRAINT IF EXISTS chk_interview_sessions_total_duration;

ALTER TABLE IF EXISTS interview.interview_sessions
    DROP CONSTRAINT IF EXISTS chk_interview_sessions_total_score;

ALTER TABLE IF EXISTS interview.interview_sessions
    DROP CONSTRAINT IF EXISTS chk_interview_sessions_answered_count;

ALTER TABLE IF EXISTS interview.interview_sessions
    DROP CONSTRAINT IF EXISTS chk_interview_sessions_question_count;

ALTER TABLE IF EXISTS interview.interview_sessions
    DROP CONSTRAINT IF EXISTS chk_interview_sessions_experience_level;

ALTER TABLE IF EXISTS interview.interview_sessions
    DROP CONSTRAINT IF EXISTS chk_interview_sessions_difficulty;

ALTER TABLE IF EXISTS interview.interview_sessions
    DROP CONSTRAINT IF EXISTS chk_interview_sessions_status;

ALTER TABLE IF EXISTS interview.interview_reviews
    DROP CONSTRAINT IF EXISTS fk_interview_reviews_answer;

ALTER TABLE IF EXISTS interview.interview_reviews
    DROP CONSTRAINT IF EXISTS fk_interview_reviews_question;

ALTER TABLE IF EXISTS interview.interview_answers
    DROP CONSTRAINT IF EXISTS fk_interview_answers_question;

ALTER TABLE IF EXISTS interview.interview_questions
    DROP CONSTRAINT IF EXISTS fk_interview_questions_session;

DROP INDEX IF EXISTS interview.idx_interview_reviews_answer_id;
DROP INDEX IF EXISTS interview.uk_interview_reviews_session_question_id;
DROP INDEX IF EXISTS interview.uk_interview_answers_session_question_id;
DROP INDEX IF EXISTS interview.idx_interview_questions_topic_key;
DROP INDEX IF EXISTS interview.idx_interview_questions_session_id;
DROP INDEX IF EXISTS interview.uk_interview_questions_session_position;
DROP INDEX IF EXISTS interview.idx_interview_sessions_started_at;
DROP INDEX IF EXISTS interview.idx_interview_sessions_user_status;
DROP INDEX IF EXISTS interview.idx_interview_sessions_user_id;

DROP TABLE IF EXISTS interview.interview_reviews;
DROP TABLE IF EXISTS interview.interview_answers;
DROP TABLE IF EXISTS interview.interview_questions;
DROP TABLE IF EXISTS interview.interview_sessions;
DROP SCHEMA IF EXISTS interview;
-- +goose StatementEnd
