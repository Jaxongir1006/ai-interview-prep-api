package audit

import "context"

// Portal defines the audit logging interface exposed to other modules.
type Portal interface {
	// Log records an action log with optional configuration via LogOption.
	// Uses the caller's transaction context (Borrow pattern).
	Log(ctx context.Context, action Action, opts ...LogOption) error
}

// Action represents a user action to be logged.
type Action struct {
	Module      string
	OperationID string
	Payload     any
}

// StatusChange represents an entity state transition to be logged.
type StatusChange struct {
	EntityType string
	EntityID   string
	Status     string
}

// LogConfig holds resolved option values for a Log call.
type LogConfig struct {
	Tags          []string
	GroupKey      string
	StatusChanges []StatusChange
}

// LogOption configures optional parameters for Portal.Log.
type LogOption func(*LogConfig)

// WithTags attaches categorical tags to the action log.
func WithTags(tags ...string) LogOption {
	return func(c *LogConfig) { c.Tags = tags }
}

// WithGroupKey sets the business group identifier (e.g. "tender:<uuid>").
func WithGroupKey(groupKey string) LogOption {
	return func(c *LogConfig) { c.GroupKey = groupKey }
}

// WithStatusChanges attaches status change logs to the action log.
func WithStatusChanges(changes ...StatusChange) LogOption {
	return func(c *LogConfig) { c.StatusChanges = changes }
}

// BuildLogConfig applies options and returns the resolved config.
func BuildLogConfig(opts []LogOption) LogConfig {
	var c LogConfig
	for _, opt := range opts {
		opt(&c)
	}
	return c
}
