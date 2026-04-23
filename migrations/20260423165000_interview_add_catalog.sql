-- +goose Up
-- +goose StatementBegin
CREATE TABLE interview.interview_target_roles (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE interview.interview_experience_levels (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE interview.interview_topics (
    id BIGSERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE interview.interview_target_role_topics (
    id BIGSERIAL PRIMARY KEY,
    target_role_id BIGINT NOT NULL,
    topic_id BIGINT NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX uk_interview_target_roles_key
    ON interview.interview_target_roles (key);

CREATE UNIQUE INDEX uk_interview_experience_levels_key
    ON interview.interview_experience_levels (key);

CREATE UNIQUE INDEX uk_interview_topics_key
    ON interview.interview_topics (key);

CREATE UNIQUE INDEX uk_interview_target_role_topics_role_topic
    ON interview.interview_target_role_topics (target_role_id, topic_id);

CREATE INDEX idx_interview_target_role_topics_topic_id
    ON interview.interview_target_role_topics (topic_id);

ALTER TABLE interview.interview_target_role_topics
    ADD CONSTRAINT fk_interview_target_role_topics_role
        FOREIGN KEY (target_role_id)
        REFERENCES interview.interview_target_roles(id)
        ON DELETE CASCADE;

ALTER TABLE interview.interview_target_role_topics
    ADD CONSTRAINT fk_interview_target_role_topics_topic
        FOREIGN KEY (topic_id)
        REFERENCES interview.interview_topics(id)
        ON DELETE CASCADE;

ALTER TABLE interview.interview_target_roles
    ADD CONSTRAINT chk_interview_target_roles_display_order
        CHECK (display_order >= 0);

ALTER TABLE interview.interview_experience_levels
    ADD CONSTRAINT chk_interview_experience_levels_display_order
        CHECK (display_order >= 0);

ALTER TABLE interview.interview_topics
    ADD CONSTRAINT chk_interview_topics_display_order
        CHECK (display_order >= 0);

ALTER TABLE interview.interview_target_role_topics
    ADD CONSTRAINT chk_interview_target_role_topics_display_order
        CHECK (display_order >= 0);

INSERT INTO interview.interview_target_roles (key, name, description, display_order)
VALUES
    ('python', 'Python Developer', 'Backend and application interviews focused on Python', 10),
    ('golang', 'Go Developer', 'Backend interviews focused on Go services', 20),
    ('javascript', 'JavaScript Developer', 'Frontend and full-stack interviews focused on JavaScript', 30),
    ('java', 'Java Developer', 'Backend interviews focused on Java services', 40);

INSERT INTO interview.interview_experience_levels (key, name, description, display_order)
VALUES
    ('junior', 'Junior', 'Early-career interview expectations', 10),
    ('mid', 'Mid-level', 'Intermediate interview expectations', 20),
    ('senior', 'Senior', 'Senior interview expectations', 30);

INSERT INTO interview.interview_topics (key, name, description, category, display_order)
VALUES
    ('algorithms', 'Algorithms', 'Data structures, complexity, and problem-solving', 'core', 10),
    ('system-design', 'System Design', 'Scalable architecture and distributed systems', 'backend', 20),
    ('database-design', 'Database Design', 'Relational modeling, indexing, and query design', 'backend', 30);

INSERT INTO interview.interview_target_role_topics (target_role_id, topic_id, display_order)
SELECT r.id, t.id, t.display_order
FROM interview.interview_target_roles r
CROSS JOIN interview.interview_topics t
WHERE r.key IN ('python', 'golang', 'javascript', 'java')
  AND t.key IN ('algorithms', 'system-design', 'database-design');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS interview.interview_target_role_topics
    DROP CONSTRAINT IF EXISTS chk_interview_target_role_topics_display_order;

ALTER TABLE IF EXISTS interview.interview_topics
    DROP CONSTRAINT IF EXISTS chk_interview_topics_display_order;

ALTER TABLE IF EXISTS interview.interview_experience_levels
    DROP CONSTRAINT IF EXISTS chk_interview_experience_levels_display_order;

ALTER TABLE IF EXISTS interview.interview_target_roles
    DROP CONSTRAINT IF EXISTS chk_interview_target_roles_display_order;

ALTER TABLE IF EXISTS interview.interview_target_role_topics
    DROP CONSTRAINT IF EXISTS fk_interview_target_role_topics_topic;

ALTER TABLE IF EXISTS interview.interview_target_role_topics
    DROP CONSTRAINT IF EXISTS fk_interview_target_role_topics_role;

DROP INDEX IF EXISTS interview.idx_interview_target_role_topics_topic_id;
DROP INDEX IF EXISTS interview.uk_interview_target_role_topics_role_topic;
DROP INDEX IF EXISTS interview.uk_interview_topics_key;
DROP INDEX IF EXISTS interview.uk_interview_experience_levels_key;
DROP INDEX IF EXISTS interview.uk_interview_target_roles_key;

DROP TABLE IF EXISTS interview.interview_target_role_topics;
DROP TABLE IF EXISTS interview.interview_topics;
DROP TABLE IF EXISTS interview.interview_experience_levels;
DROP TABLE IF EXISTS interview.interview_target_roles;
-- +goose StatementEnd
