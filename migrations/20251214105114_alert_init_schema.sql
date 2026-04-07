-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA alert;

CREATE TABLE alert.errors (
    id VARCHAR PRIMARY KEY,
    code TEXT NOT NULL,
    message TEXT NOT NULL,
    details JSONB NOT NULL,
    service TEXT NOT NULL,
    operation TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    alerted BOOLEAN NOT NULL
);

CREATE INDEX idx_errors_service_operation_alerted
    ON alert.errors (service, operation, alerted);

CREATE INDEX idx_errors_created_at
    ON alert.errors (created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS alert.idx_errors_created_at;
DROP INDEX IF EXISTS alert.idx_errors_service_operation_alerted;
DROP TABLE IF EXISTS alert.errors;
DROP SCHEMA IF EXISTS alert;

-- +goose StatementEnd
