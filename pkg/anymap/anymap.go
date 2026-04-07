// Package anymap provides utilities for extracting typed values from map[string]any.
package anymap

import (
	"fmt"
	"testing"
	"time"
)

// String extracts a string value from data, returns defaultVal if not present.
// If the value exists but is not a string, it converts using fmt.Sprintf.
func String(data map[string]any, key string, defaultVal string) string {
	if v, hasKey := data[key]; hasKey {
		if s, isString := v.(string); isString {
			return s
		}
		return fmt.Sprintf("%v", v)
	}
	return defaultVal
}

// StringPtr extracts a *string value from data, returns defaultVal if not present.
// If the value exists and is nil, returns nil.
func StringPtr(data map[string]any, key string, defaultVal *string) *string {
	if v, hasKey := data[key]; hasKey {
		if v == nil {
			return nil
		}
		if s, isString := v.(string); isString {
			return &s
		}
		s := fmt.Sprintf("%v", v)
		return &s
	}
	return defaultVal
}

// Bool extracts a bool value from data, returns defaultVal if not present.
func Bool(data map[string]any, key string, defaultVal bool) bool {
	if v, hasKey := data[key]; hasKey {
		if b, isBool := v.(bool); isBool {
			return b
		}
	}
	return defaultVal
}

// Time extracts a time.Time value from data, returns defaultVal if not present.
func Time(data map[string]any, key string, defaultVal time.Time) time.Time {
	if v, hasKey := data[key]; hasKey {
		if timeVal, isTime := v.(time.Time); isTime {
			return timeVal
		}
	}
	return defaultVal
}

// TimePtr extracts a *time.Time value from data, returns defaultVal if not present.
func TimePtr(data map[string]any, key string, defaultVal *time.Time) *time.Time {
	if v, hasKey := data[key]; hasKey {
		if timeVal, isTime := v.(time.Time); isTime {
			return &timeVal
		}
	}
	return defaultVal
}

// StringSlice extracts a []string value from data, returns defaultVal if not present.
func StringSlice(data map[string]any, key string, defaultVal []string) []string {
	if v, hasKey := data[key]; hasKey {
		if ss, isSlice := v.([]string); isSlice {
			return ss
		}
	}
	return defaultVal
}

// ValidateKeys checks that all keys in data exist in validKeys.
// Fails the test immediately if an unknown key is found.
func ValidateKeys(t *testing.T, funcName string, validKeys map[string]bool, data map[string]any) {
	t.Helper()

	for key := range data {
		if !validKeys[key] {
			var valid []string
			for k := range validKeys {
				valid = append(valid, k)
			}
			t.Fatalf("%s: unknown field %q (valid fields: %v)", funcName, key, valid)
		}
	}
}
