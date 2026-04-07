package anymap_test

import (
	"go-enterprise-blueprint/pkg/anymap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		name       string
		data       map[string]any
		key        string
		defaultVal string
		expected   string
	}{
		{
			name:       "returns value when key exists with string type",
			data:       map[string]any{"name": "alice"},
			key:        "name",
			defaultVal: "default",
			expected:   "alice",
		},
		{
			name:       "converts non-string value using fmt.Sprintf",
			data:       map[string]any{"count": 42},
			key:        "count",
			defaultVal: "default",
			expected:   "42",
		},
		{
			name:       "returns default when key is missing",
			data:       map[string]any{},
			key:        "name",
			defaultVal: "default",
			expected:   "default",
		},
		{
			name:       "returns empty string default",
			data:       map[string]any{},
			key:        "name",
			defaultVal: "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := anymap.String(tt.data, tt.key, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name       string
		data       map[string]any
		key        string
		defaultVal bool
		expected   bool
	}{
		{
			name:       "returns value when key exists with anymap.bool type",
			data:       map[string]any{"active": true},
			key:        "active",
			defaultVal: false,
			expected:   true,
		},
		{
			name:       "returns false value",
			data:       map[string]any{"active": false},
			key:        "active",
			defaultVal: true,
			expected:   false,
		},
		{
			name:       "returns default when key is missing",
			data:       map[string]any{},
			key:        "active",
			defaultVal: true,
			expected:   true,
		},
		{
			name:       "returns default when value is not bool",
			data:       map[string]any{"active": "yes"},
			key:        "active",
			defaultVal: false,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := anymap.Bool(tt.data, tt.key, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTime(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	defaultTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		data       map[string]any
		key        string
		defaultVal time.Time
		expected   time.Time
	}{
		{
			name:       "returns value when key exists with time.Time type",
			data:       map[string]any{"created_at": now},
			key:        "created_at",
			defaultVal: defaultTime,
			expected:   now,
		},
		{
			name:       "returns default when key is missing",
			data:       map[string]any{},
			key:        "created_at",
			defaultVal: defaultTime,
			expected:   defaultTime,
		},
		{
			name:       "returns default when value is not time.Time",
			data:       map[string]any{"created_at": "2025-01-15"},
			key:        "created_at",
			defaultVal: defaultTime,
			expected:   defaultTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := anymap.Time(tt.data, tt.key, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTimePtr(t *testing.T) {
	now := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	defaultTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		data       map[string]any
		key        string
		defaultVal *time.Time
		expected   *time.Time
	}{
		{
			name:       "returns pointer to value when key exists",
			data:       map[string]any{"expires_at": now},
			key:        "expires_at",
			defaultVal: nil,
			expected:   &now,
		},
		{
			name:       "returns nil default when key is missing",
			data:       map[string]any{},
			key:        "expires_at",
			defaultVal: nil,
			expected:   nil,
		},
		{
			name:       "returns non-nil default when key is missing",
			data:       map[string]any{},
			key:        "expires_at",
			defaultVal: &defaultTime,
			expected:   &defaultTime,
		},
		{
			name:       "returns default when value is not time.Time",
			data:       map[string]any{"expires_at": "2025-01-15"},
			key:        "expires_at",
			defaultVal: nil,
			expected:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := anymap.TimePtr(tt.data, tt.key, tt.defaultVal)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateKeys(t *testing.T) {
	validKeys := map[string]bool{
		"name":   true,
		"age":    true,
		"active": true,
	}

	t.Run("passes when all keys are valid", func(t *testing.T) {
		data := map[string]any{"name": "alice", "age": 30}
		anymap.ValidateKeys(t, "TestFunc", validKeys, data)
	})

	t.Run("passes with empty data", func(t *testing.T) {
		data := map[string]any{}
		anymap.ValidateKeys(t, "TestFunc", validKeys, data)
	})
}
