-- =============================================================================
-- E2E Test Seeds: Platform Module
--
-- Seeds test data for frontend E2E tests. Run after backend starts and
-- superadmin is seeded. Auth/audit data is bootstrapped via API calls
-- in the frontend test setup; platform data has no create API so we
-- insert directly.
--
-- Target database: test_interview-prep
-- =============================================================================

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- 1. TaskMill: Queues with tasks (generates queue stats via the view)
-- ─────────────────────────────────────────────────────────────────────────────

-- Queue "auth": 3 available, 1 in-flight, 1 in DLQ
INSERT INTO taskmill.task_queue
    (queue_name, task_group_id, operation_id, meta, payload, scheduled_at, visible_at, priority, attempts, max_attempts, idempotency_key, created_at, updated_at)
VALUES
    ('auth', NULL, 'clean-expired-sessions', '{}', '{"batch_size": 100}',
     NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour', 0, 0, 3,
     'seed-auth-001', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour'),
    ('auth', NULL, 'clean-expired-sessions', '{}', '{"batch_size": 100}',
     NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes', 0, 0, 3,
     'seed-auth-002', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes'),
    ('auth', NULL, 'clean-expired-sessions', '{}', '{"batch_size": 100}',
     NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '15 minutes', 0, 0, 3,
     'seed-auth-003', NOW() - INTERVAL '15 minutes', NOW() - INTERVAL '15 minutes'),
    -- In-flight (visible_at in the future)
    ('auth', NULL, 'clean-expired-sessions', '{}', '{"batch_size": 50}',
     NOW() - INTERVAL '5 minutes', NOW() + INTERVAL '5 minutes', 0, 1, 3,
     'seed-auth-inflight-001', NOW() - INTERVAL '5 minutes', NOW());

-- DLQ entry for auth queue
INSERT INTO taskmill.task_queue
    (queue_name, task_group_id, operation_id, meta, payload, scheduled_at, visible_at, priority, attempts, max_attempts, idempotency_key, created_at, updated_at, dlq_at, dlq_reason)
VALUES
    ('auth', NULL, 'clean-expired-sessions', '{}', '{"batch_size": 200}',
     NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', 0, 3, 3,
     'seed-auth-dlq-001', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '1 hour',
     NOW() - INTERVAL '1 hour', '{"error": "connection timeout"}');

-- Queue "platform": 2 available, 1 in DLQ
INSERT INTO taskmill.task_queue
    (queue_name, task_group_id, operation_id, meta, payload, scheduled_at, visible_at, priority, attempts, max_attempts, idempotency_key, created_at, updated_at)
VALUES
    ('platform', 'report-batch-001', 'generate-error-report', '{"source": "alert"}', '{"date": "2026-02-13"}',
     NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours', 1, 0, 3,
     'seed-platform-001', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours'),
    ('platform', 'report-batch-001', 'generate-error-report', '{"source": "alert"}', '{"date": "2026-02-14"}',
     NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour', 1, 0, 3,
     'seed-platform-002', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour');

INSERT INTO taskmill.task_queue
    (queue_name, task_group_id, operation_id, meta, payload, scheduled_at, visible_at, priority, attempts, max_attempts, idempotency_key, created_at, updated_at, dlq_at, dlq_reason)
VALUES
    ('platform', NULL, 'cleanup-old-results', '{}', '{"days": 30}',
     NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours', 0, 3, 3,
     'seed-platform-dlq-001', NOW() - INTERVAL '3 hours', NOW() - INTERVAL '2 hours',
     NOW() - INTERVAL '2 hours', '{"error": "deadlock detected"}');


-- ─────────────────────────────────────────────────────────────────────────────
-- 2. TaskMill: Task Results (completed tasks)
-- ─────────────────────────────────────────────────────────────────────────────

INSERT INTO taskmill.task_results
    (id, queue_name, task_group_id, operation_id, meta, payload, priority, attempts, max_attempts, idempotency_key, scheduled_at, created_at, completed_at)
VALUES
    (1001, 'auth', NULL, 'clean-expired-sessions', '{}', '{"batch_size": 100}',
     0, 1, 3, 'seed-result-001',
     NOW() - INTERVAL '6 hours', NOW() - INTERVAL '6 hours', NOW() - INTERVAL '5 hours'),
    (1002, 'auth', NULL, 'clean-expired-sessions', '{}', '{"batch_size": 100}',
     0, 1, 3, 'seed-result-002',
     NOW() - INTERVAL '4 hours', NOW() - INTERVAL '4 hours', NOW() - INTERVAL '3 hours'),
    (1003, 'platform', 'report-batch-000', 'generate-error-report', '{"source": "alert"}', '{"date": "2026-02-12"}',
     1, 1, 3, 'seed-result-003',
     NOW() - INTERVAL '24 hours', NOW() - INTERVAL '24 hours', NOW() - INTERVAL '23 hours'),
    (1004, 'platform', 'report-batch-000', 'generate-error-report', '{"source": "alert"}', '{"date": "2026-02-11"}',
     1, 2, 3, 'seed-result-004',
     NOW() - INTERVAL '48 hours', NOW() - INTERVAL '48 hours', NOW() - INTERVAL '47 hours');


-- ─────────────────────────────────────────────────────────────────────────────
-- 3. TaskMill: Schedules
-- ─────────────────────────────────────────────────────────────────────────────
-- Auth module auto-registers "clean-expired-sessions" on startup.
-- Add one more so the UI shows multiple rows.

INSERT INTO taskmill.task_schedules
    (operation_id, queue_name, cron_pattern, next_run_at, last_run_at, last_run_status, last_error, run_count, created_at, updated_at)
VALUES
    ('cleanup-old-results', 'platform', '0 3 * * *',
     DATE_TRUNC('day', NOW()) + INTERVAL '1 day 3 hours',
     NOW() - INTERVAL '21 hours', 'success', NULL, 14,
     NOW() - INTERVAL '14 days', NOW() - INTERVAL '21 hours')
ON CONFLICT (operation_id) DO NOTHING;


-- ─────────────────────────────────────────────────────────────────────────────
-- 4. Alert: Errors (different codes, services, operations for filter testing)
-- ─────────────────────────────────────────────────────────────────────────────

INSERT INTO alert.errors
    (id, code, message, details, service, operation, created_at, alerted)
VALUES
    ('a1b2c3d4-0001-4000-8000-000000000001', 'TIMEOUT', 'Request timed out after 30s',
     '{"endpoint": "/api/v1/auth/get-users", "duration_ms": "30000"}',
     'blueprint', 'get-users',
     NOW() - INTERVAL '2 hours', false),

    ('a1b2c3d4-0002-4000-8000-000000000002', 'DB_CONNECTION', 'Failed to acquire connection from pool',
     '{"pool_size": "10", "active": "10", "waiting": "5"}',
     'blueprint', 'create-user',
     NOW() - INTERVAL '90 minutes', false),

    ('a1b2c3d4-0003-4000-8000-000000000003', 'VALIDATION', 'Invalid input: username already exists',
     '{"field": "username", "value": "duplicate_user"}',
     'blueprint', 'create-user',
     NOW() - INTERVAL '1 hour', true),

    ('a1b2c3d4-0004-4000-8000-000000000004', 'KAFKA_PUBLISH', 'Failed to publish event to Kafka',
     '{"topic": "auth.user.created", "error": "broker not available"}',
     'blueprint', 'create-user',
     NOW() - INTERVAL '45 minutes', false),

    ('a1b2c3d4-0005-4000-8000-000000000005', 'TIMEOUT', 'Request timed out after 30s',
     '{"endpoint": "/api/v1/platform/list-queues", "duration_ms": "30000"}',
     'blueprint', 'list-queues',
     NOW() - INTERVAL '30 minutes', true),

    ('a1b2c3d4-0006-4000-8000-000000000006', 'INTERNAL', 'Unexpected nil pointer in session handler',
     '{"stack_trace": "auth/session.go:142"}',
     'blueprint', 'get-my-sessions',
     NOW() - INTERVAL '15 minutes', false);

COMMIT;
