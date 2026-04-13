-- +goose Up
-- +goose StatementBegin
ALTER TABLE candidate.candidate_profiles
    ADD COLUMN IF NOT EXISTS onboarding_completed BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS onboarding_completed_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS candidate.candidate_profiles
    DROP COLUMN IF EXISTS onboarding_completed_at,
    DROP COLUMN IF EXISTS onboarding_completed;
-- +goose StatementEnd
