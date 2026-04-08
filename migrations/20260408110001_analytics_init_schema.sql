-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS analytics;

CREATE TABLE analytics.candidate_progress_summaries (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    total_interviews_taken INTEGER NOT NULL DEFAULT 0,
    total_time_spent_seconds BIGINT NOT NULL DEFAULT 0,
    average_score NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    last_interview_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE analytics.candidate_topic_stats (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    topic_key VARCHAR(100) NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    total_time_spent_seconds BIGINT NOT NULL DEFAULT 0,
    average_score NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    best_score NUMERIC(5,2) NOT NULL DEFAULT 0.00,
    last_practiced_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE analytics.achievement_definitions (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE analytics.candidate_achievements (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    achievement_definition_id BIGINT NOT NULL,
    achieved_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uk_candidate_progress_summaries_user_id
    ON analytics.candidate_progress_summaries (user_id);

CREATE UNIQUE INDEX uk_candidate_topic_stats_user_topic
    ON analytics.candidate_topic_stats (user_id, topic_key);

CREATE UNIQUE INDEX uk_achievement_definitions_code
    ON analytics.achievement_definitions (code);

CREATE UNIQUE INDEX uk_candidate_achievements_user_definition
    ON analytics.candidate_achievements (user_id, achievement_definition_id);

CREATE INDEX idx_candidate_topic_stats_user_id
    ON analytics.candidate_topic_stats (user_id);

CREATE INDEX idx_candidate_achievements_user_id
    ON analytics.candidate_achievements (user_id);

ALTER TABLE analytics.candidate_achievements
    ADD CONSTRAINT fk_candidate_achievements_definition
        FOREIGN KEY (achievement_definition_id)
        REFERENCES analytics.achievement_definitions(id)
        ON DELETE CASCADE;

ALTER TABLE analytics.candidate_progress_summaries
    ADD CONSTRAINT chk_candidate_progress_summaries_current_streak
        CHECK (current_streak >= 0);

ALTER TABLE analytics.candidate_progress_summaries
    ADD CONSTRAINT chk_candidate_progress_summaries_longest_streak
        CHECK (longest_streak >= 0);

ALTER TABLE analytics.candidate_progress_summaries
    ADD CONSTRAINT chk_candidate_progress_summaries_total_interviews
        CHECK (total_interviews_taken >= 0);

ALTER TABLE analytics.candidate_progress_summaries
    ADD CONSTRAINT chk_candidate_progress_summaries_total_time
        CHECK (total_time_spent_seconds >= 0);

ALTER TABLE analytics.candidate_progress_summaries
    ADD CONSTRAINT chk_candidate_progress_summaries_average_score
        CHECK (average_score >= 0 AND average_score <= 100);

ALTER TABLE analytics.candidate_topic_stats
    ADD CONSTRAINT chk_candidate_topic_stats_attempts
        CHECK (attempts >= 0);

ALTER TABLE analytics.candidate_topic_stats
    ADD CONSTRAINT chk_candidate_topic_stats_total_time
        CHECK (total_time_spent_seconds >= 0);

ALTER TABLE analytics.candidate_topic_stats
    ADD CONSTRAINT chk_candidate_topic_stats_average_score
        CHECK (average_score >= 0 AND average_score <= 100);

ALTER TABLE analytics.candidate_topic_stats
    ADD CONSTRAINT chk_candidate_topic_stats_best_score
        CHECK (best_score >= 0 AND best_score <= 100);

ALTER TABLE analytics.achievement_definitions
    ADD CONSTRAINT chk_achievement_definitions_sort_order
        CHECK (sort_order >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS analytics.achievement_definitions
    DROP CONSTRAINT IF EXISTS chk_achievement_definitions_sort_order;

ALTER TABLE IF EXISTS analytics.candidate_topic_stats
    DROP CONSTRAINT IF EXISTS chk_candidate_topic_stats_best_score;

ALTER TABLE IF EXISTS analytics.candidate_topic_stats
    DROP CONSTRAINT IF EXISTS chk_candidate_topic_stats_average_score;

ALTER TABLE IF EXISTS analytics.candidate_topic_stats
    DROP CONSTRAINT IF EXISTS chk_candidate_topic_stats_total_time;

ALTER TABLE IF EXISTS analytics.candidate_topic_stats
    DROP CONSTRAINT IF EXISTS chk_candidate_topic_stats_attempts;

ALTER TABLE IF EXISTS analytics.candidate_progress_summaries
    DROP CONSTRAINT IF EXISTS chk_candidate_progress_summaries_average_score;

ALTER TABLE IF EXISTS analytics.candidate_progress_summaries
    DROP CONSTRAINT IF EXISTS chk_candidate_progress_summaries_total_time;

ALTER TABLE IF EXISTS analytics.candidate_progress_summaries
    DROP CONSTRAINT IF EXISTS chk_candidate_progress_summaries_total_interviews;

ALTER TABLE IF EXISTS analytics.candidate_progress_summaries
    DROP CONSTRAINT IF EXISTS chk_candidate_progress_summaries_longest_streak;

ALTER TABLE IF EXISTS analytics.candidate_progress_summaries
    DROP CONSTRAINT IF EXISTS chk_candidate_progress_summaries_current_streak;

ALTER TABLE IF EXISTS analytics.candidate_achievements
    DROP CONSTRAINT IF EXISTS fk_candidate_achievements_definition;

DROP INDEX IF EXISTS idx_candidate_achievements_user_id;
DROP INDEX IF EXISTS idx_candidate_topic_stats_user_id;
DROP INDEX IF EXISTS uk_candidate_achievements_user_definition;
DROP INDEX IF EXISTS uk_achievement_definitions_code;
DROP INDEX IF EXISTS uk_candidate_topic_stats_user_topic;
DROP INDEX IF EXISTS uk_candidate_progress_summaries_user_id;

DROP TABLE IF EXISTS analytics.candidate_achievements;
DROP TABLE IF EXISTS analytics.achievement_definitions;
DROP TABLE IF EXISTS analytics.candidate_topic_stats;
DROP TABLE IF EXISTS analytics.candidate_progress_summaries;
DROP SCHEMA IF EXISTS analytics;
-- +goose StatementEnd
