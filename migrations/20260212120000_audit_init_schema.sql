-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA audit;

-- action_logs: range-partitioned by created_at (monthly)
CREATE TABLE audit.action_logs (
    id BIGSERIAL,
    user_id VARCHAR,
    module VARCHAR NOT NULL,
    operation_id VARCHAR NOT NULL,
    request_payload JSONB NOT NULL DEFAULT '{}',
    ip_address VARCHAR NOT NULL DEFAULT '',
    user_agent VARCHAR NOT NULL DEFAULT '',
    tags VARCHAR[] NOT NULL DEFAULT '{}',
    group_key VARCHAR,
    trace_id VARCHAR NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- status_change_logs: range-partitioned by created_at (monthly)
-- NOTE: No FK from action_log_id to action_logs.id because PostgreSQL requires
-- unique constraints on partitioned tables to include the partition key (created_at),
-- making a simple FK reference to (id) alone impossible.
CREATE TABLE audit.status_change_logs (
    id BIGSERIAL,
    action_log_id BIGINT NOT NULL,
    entity_type VARCHAR NOT NULL,
    entity_id VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    trace_id VARCHAR NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Indexes for action_logs
CREATE INDEX idx_action_logs_user_id ON audit.action_logs (user_id);
CREATE INDEX idx_action_logs_module_operation ON audit.action_logs (module, operation_id);
CREATE INDEX idx_action_logs_trace_id ON audit.action_logs (trace_id);
CREATE INDEX idx_action_logs_tags ON audit.action_logs USING GIN (tags);
CREATE INDEX idx_action_logs_group_key ON audit.action_logs (group_key);

-- Indexes for status_change_logs
CREATE INDEX idx_status_change_logs_action_log_id ON audit.status_change_logs (action_log_id);
CREATE INDEX idx_status_change_logs_entity ON audit.status_change_logs (entity_type, entity_id);
CREATE INDEX idx_status_change_logs_trace_id ON audit.status_change_logs (trace_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS audit CASCADE;
-- +goose StatementEnd
