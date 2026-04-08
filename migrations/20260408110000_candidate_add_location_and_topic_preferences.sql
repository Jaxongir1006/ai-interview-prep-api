-- +goose Up
-- +goose StatementBegin
ALTER TABLE candidate.candidate_profiles
    ADD COLUMN IF NOT EXISTS location VARCHAR(255);

CREATE TABLE candidate.candidate_topic_preferences (
    id BIGSERIAL PRIMARY KEY,
    candidate_profile_id BIGINT NOT NULL,
    topic_key VARCHAR(100) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_candidate_topic_preferences_profile_id
    ON candidate.candidate_topic_preferences (candidate_profile_id);

CREATE UNIQUE INDEX uk_candidate_topic_preferences_profile_topic
    ON candidate.candidate_topic_preferences (candidate_profile_id, topic_key);

ALTER TABLE candidate.candidate_topic_preferences
    ADD CONSTRAINT fk_candidate_topic_preferences_profile
        FOREIGN KEY (candidate_profile_id)
        REFERENCES candidate.candidate_profiles(id)
        ON DELETE CASCADE;

ALTER TABLE candidate.candidate_topic_preferences
    ADD CONSTRAINT chk_candidate_topic_preferences_priority
        CHECK (priority >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS candidate.candidate_topic_preferences
    DROP CONSTRAINT IF EXISTS chk_candidate_topic_preferences_priority;

ALTER TABLE IF EXISTS candidate.candidate_topic_preferences
    DROP CONSTRAINT IF EXISTS fk_candidate_topic_preferences_profile;

DROP INDEX IF EXISTS uk_candidate_topic_preferences_profile_topic;
DROP INDEX IF EXISTS idx_candidate_topic_preferences_profile_id;

DROP TABLE IF EXISTS candidate.candidate_topic_preferences;

ALTER TABLE candidate.candidate_profiles
    DROP COLUMN IF EXISTS location;
-- +goose StatementEnd
