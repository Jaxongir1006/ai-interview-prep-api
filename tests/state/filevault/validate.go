package filevault

//nolint:gochecknoglobals // static validation maps for test state
var (
	// validFileKeys defines the allowed keys for file test data.
	// Keys correspond to file.File entity fields.
	validFileKeys = map[string]bool{
		"id":               true,
		"original_name":    true,
		"stored_name":      true,
		"content_type":     true,
		"size":             true,
		"checksum":         true,
		"path":             true,
		"entity_type":      true,
		"entity_id":        true,
		"association_type": true,
		"sort_order":       true,
		"uploaded_by":      true,
		"storage_status":   true,
	}
)
