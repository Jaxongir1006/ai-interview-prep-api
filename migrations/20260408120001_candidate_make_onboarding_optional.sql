-- +goose Up
-- +goose StatementBegin
ALTER TABLE candidate.candidate_profiles
    ALTER COLUMN target_role DROP NOT NULL,
    ALTER COLUMN experience_level DROP NOT NULL;

ALTER TABLE candidate.candidate_profiles
    DROP CONSTRAINT IF EXISTS chk_candidate_profiles_experience_level;

ALTER TABLE candidate.candidate_profiles
    ADD CONSTRAINT chk_candidate_profiles_experience_level
        CHECK (experience_level IS NULL OR experience_level IN ('junior', 'mid', 'senior'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE candidate.candidate_profiles
    DROP CONSTRAINT IF EXISTS chk_candidate_profiles_experience_level;

UPDATE candidate.candidate_profiles
SET
    target_role = COALESCE(target_role, 'unknown'),
    experience_level = COALESCE(experience_level, 'junior');

ALTER TABLE candidate.candidate_profiles
    ADD CONSTRAINT chk_candidate_profiles_experience_level
        CHECK (experience_level IN ('junior', 'mid', 'senior'));

ALTER TABLE candidate.candidate_profiles
    ALTER COLUMN experience_level SET NOT NULL,
    ALTER COLUMN target_role SET NOT NULL;
-- +goose StatementEnd
