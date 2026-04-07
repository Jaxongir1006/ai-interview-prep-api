package audit

//nolint:gochecknoglobals // static validation maps for test state
var (
	validActionLogKeys = map[string]bool{
		"user_id": true, "module": true, "operation_id": true,
		"request_payload": true, "ip_address": true, "user_agent": true,
		"tags": true, "group_key": true,
		"trace_id": true, "created_at": true,
	}

	validStatusChangeLogKeys = map[string]bool{
		"action_log_id": true, "entity_type": true, "entity_id": true,
		"status": true, "trace_id": true, "created_at": true,
	}
)
