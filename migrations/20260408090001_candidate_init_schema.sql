-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS candidate;

CREATE TABLE candidate.candidate_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    full_name VARCHAR(255),
    bio TEXT,
    target_role VARCHAR(100) NOT NULL,
    experience_level VARCHAR(20) NOT NULL,
    interview_goal_per_week INTEGER NOT NULL DEFAULT 3,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uk_candidate_profiles_user_id
    ON candidate.candidate_profiles (user_id);

CREATE INDEX idx_candidate_profiles_target_role
    ON candidate.candidate_profiles (target_role);

ALTER TABLE candidate.candidate_profiles
    ADD CONSTRAINT fk_candidate_profiles_user
        FOREIGN KEY (user_id)
        REFERENCES auth.users(id)
        ON DELETE CASCADE;

ALTER TABLE candidate.candidate_profiles
    ADD CONSTRAINT chk_candidate_profiles_experience_level
        CHECK (experience_level IN ('junior', 'mid', 'senior'));

ALTER TABLE candidate.candidate_profiles
    ADD CONSTRAINT chk_candidate_profiles_interview_goal_per_week
        CHECK (interview_goal_per_week >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS candidate.candidate_profiles
    DROP CONSTRAINT IF EXISTS chk_candidate_profiles_interview_goal_per_week;

ALTER TABLE IF EXISTS candidate.candidate_profiles
    DROP CONSTRAINT IF EXISTS chk_candidate_profiles_experience_level;

ALTER TABLE IF EXISTS candidate.candidate_profiles
    DROP CONSTRAINT IF EXISTS fk_candidate_profiles_user;

DROP INDEX IF EXISTS idx_candidate_profiles_target_role;
DROP INDEX IF EXISTS uk_candidate_profiles_user_id;

DROP TABLE IF EXISTS candidate.candidate_profiles;
DROP SCHEMA IF EXISTS candidate;
-- +goose StatementEnd
