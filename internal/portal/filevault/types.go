package filevault

import "slices"

type ContentGroup string

const (
	ContentGroupImage       ContentGroup = "image"
	ContentGroupPDF         ContentGroup = "pdf"
	ContentGroupDocument    ContentGroup = "document"
	ContentGroupSpreadsheet ContentGroup = "spreadsheet"
	ContentGroupAll         ContentGroup = "all"
)

func contentTypes() map[ContentGroup][]string {
	return map[ContentGroup][]string{
		ContentGroupImage: {
			"image/png",
			"image/jpeg",
			"image/webp",
		},
		ContentGroupPDF: {
			"application/pdf",
		},
		ContentGroupDocument: {
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.oasis.opendocument.text",
		},
		ContentGroupSpreadsheet: {
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.oasis.opendocument.spreadsheet",
		},
		ContentGroupAll: {}, // computed dynamically
	}
}

func allContentTypes() []string {
	seen := make(map[string]bool)
	var all []string

	for group, types := range contentTypes() {
		if group == ContentGroupAll {
			continue
		}
		for _, t := range types {
			if !seen[t] {
				seen[t] = true
				all = append(all, t)
			}
		}
	}

	return all
}

// AllowedContentTypes returns the list of allowed content types for a group.
func AllowedContentTypes(contentGroup ContentGroup) []string {
	if contentGroup == ContentGroupAll {
		return allContentTypes()
	}
	return contentTypes()[contentGroup]
}

// IsAllowedContentType checks if a MIME type is allowed in the given content group.
func IsAllowedContentType(group ContentGroup, mimeType string) bool {
	types := AllowedContentTypes(group)
	return slices.Contains(types, mimeType)
}
