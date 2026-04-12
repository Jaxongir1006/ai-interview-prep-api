#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

BINARY_NAME="system-test-app"
BINARY_PATH="$PROJECT_DIR/bin/$BINARY_NAME"
HEALTH_URL="http://localhost:9876/health" # host and port should be consistent with test.yaml config
HEALTH_TIMEOUT=10
TEST_DB_NAME="test_interview-prep" # should be consistent with test.yaml config

# Default superadmin credentials for E2E testing.
# Frontend E2E tests use these to bootstrap users, roles, and permissions via API.
SUPERADMIN_USERNAME="superadmin"
SUPERADMIN_PASSWORD="superadmin123"

log() {
	echo "[$(date '+%Y-%m-%d %H:%M:%S')] [prepare-system-test]: INFO $1"
}

log_error() {
	echo "[$(date '+%Y-%m-%d %H:%M:%S')] [prepare-system-test]: ERROR $1" >&2
}

# Start infrastructure
start_infra() {
	log "starting infrastructure: docker-compose -f dev-infra.yaml up -d --build"
	docker-compose -f "$PROJECT_DIR/dev-infra.yaml" up -d --build
	log "infrastructure is ready"
}

# Wait for postgres to be ready
wait_for_postgres() {
	log "waiting for postgres to be ready"
	for i in $(seq 1 30); do
		if docker-compose -f "$PROJECT_DIR/dev-infra.yaml" exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
			log "postgres is ready"
			return 0
		fi
		sleep 1
	done
	log_error "timeout waiting for postgres"
	return 1
}

# Drop and recreate test database for a clean state
create_test_db() {
	wait_for_postgres
	log "recreating test database: $TEST_DB_NAME"
	docker-compose -f "$PROJECT_DIR/dev-infra.yaml" exec -T postgres \
		psql -U postgres \
		-c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$TEST_DB_NAME' AND pid <> pg_backend_pid()" \
		-c "DROP DATABASE IF EXISTS \"$TEST_DB_NAME\"" \
		-c "CREATE DATABASE \"$TEST_DB_NAME\""
	log "test database ready"
}

# Kill existing app if running
kill_existing_app() {
	if pgrep -f "$BINARY_NAME" >/dev/null 2>&1; then
		log "killing existing $BINARY_NAME process"
		pkill -SIGINT -f "$BINARY_NAME" || true

		# Wait for graceful shutdown (max 5 seconds)
		for i in {1..10}; do
			if ! pgrep -f "$BINARY_NAME" >/dev/null 2>&1; then
				log "existing process terminated"
				return 0
			fi
			sleep 0.5
		done

		# Force kill if still running
		log "force killing existing process"
		pkill -SIGKILL -f "$BINARY_NAME" || true
		sleep 0.5
	fi
}

# Build the application
build_app() {
	log "building application: go build -race -o $BINARY_PATH ./cmd"
	cd "$PROJECT_DIR"
	go build -race -o "$BINARY_PATH" ./cmd
}

# Seed default superadmin for E2E testing
seed_superadmin() {
	log "seeding superadmin (username: $SUPERADMIN_USERNAME)"
	cd "$PROJECT_DIR"
	ENVIRONMENT="test" "$BINARY_PATH" auth create-superadmin \
		--username "$SUPERADMIN_USERNAME" \
		--password "$SUPERADMIN_PASSWORD"
	log "superadmin seeded"
}

# Seed E2E test data (platform module: queues, tasks, errors, schedules)
seed_e2e_data() {
	log "seeding E2E test data"
	docker-compose -f "$PROJECT_DIR/dev-infra.yaml" exec -T postgres \
		psql -U postgres -d "$TEST_DB_NAME" -f /dev/stdin <"$SCRIPT_DIR/seeds/e2e.sql"
	log "E2E test data seeded"
}

# Start the application
start_app() {
	log "starting application: $BINARY_PATH run"
	cd "$PROJECT_DIR"
	ENVIRONMENT="test" GORACE="halt_on_error=1" "$BINARY_PATH" run &
	APP_PID=$!
	log "application started (PID: $APP_PID)"
}

# Wait for application to be healthy
wait_until_ready() {
	log "waiting for application to be ready at $HEALTH_URL"

	for i in $(seq 1 $((HEALTH_TIMEOUT * 2))); do
		if curl -s -o /dev/null -w "%{http_code}" "$HEALTH_URL" | grep -q "200"; then
			log "application is ready"
			return 0
		fi
		sleep 0.5
	done

	log_error "timeout waiting for application to be ready"
	return 1
}

# Main
main() {
	start_infra
	kill_existing_app
	create_test_db
	build_app
	seed_superadmin
	seed_e2e_data
	start_app
	wait_until_ready
}

main "$@"
