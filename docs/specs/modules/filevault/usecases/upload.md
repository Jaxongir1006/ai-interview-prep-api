# File Upload

Uploads a file to object storage and saves its metadata to the database.

> **type**: user_action

> **operation-id**: `upload`

> **access**: POST api/v1/files/upload

> **actor**: user

> **permissions**: -

> **implementation**: [usecase.go](../../../../../internal/modules/filevault/usecase/upload/usecase.go)

## Input

Multipart form data:

- `file`: binary, required — the file to upload. Name and size are extracted from the file header; MIME type is detected from content.

## Output

```json
{
  "id": "string",
  "original_name": "string",
  "content_type": "string",
  "size": 0
}
```

## Execute

- Resolve actor ID from context

- Validate file size against configured maximum

- Detect and validate MIME type from file content

- Generate file ID and stored name

- Insert file record to database with pending status

- Upload file to object storage

- Update file record status to stored

- Return file info
