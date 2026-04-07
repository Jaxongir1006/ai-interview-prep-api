-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS filevault;

CREATE TABLE filevault.files (
    id VARCHAR PRIMARY KEY,
    original_name VARCHAR(500) NOT NULL,
    stored_name VARCHAR(500) NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    checksum VARCHAR(64),
    path VARCHAR(1000) NOT NULL,
    entity_type VARCHAR(100),
    entity_id BIGINT,
    association_type VARCHAR(100),
    sort_order INT NOT NULL DEFAULT 0,
    uploaded_by VARCHAR(255) NOT NULL,
    storage_status VARCHAR(100) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now (),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_files_entity ON filevault.files (entity_type, entity_id)
WHERE
    deleted_at IS NULL;

CREATE INDEX idx_files_uploaded_by ON filevault.files (uploaded_by)
WHERE
    deleted_at IS NULL;

CREATE INDEX idx_files_orphaned ON filevault.files (created_at)
WHERE
    entity_id IS NULL
    AND deleted_at IS NULL;

CREATE UNIQUE INDEX idx_files_storage ON filevault.files (stored_name);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS filevault.files;
DROP SCHEMA IF EXISTS filevault CASCADE;
-- +goose StatementEnd
