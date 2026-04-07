# File Download

Download a file that has been attached to an entity.

> **type**: user_action

> **operation-id**: `download`

> **access**: GET `/api/v1/files/download?id={id}`

> **actor**: user

> **permissions**: -

> **implementation**: [usecase.go](../../../../../internal/modules/filevault/usecase/download/usecase.go)

## Input

- `id`: string (UUID) (required)

## Output

- Binary file stream

## Execute

- Look up file record by ID
- Verify file is attached to an entity (`entity_id` is not null)
- Verify file storage status is `stored`
- Check `If-None-Match` header against file checksum; return 304 if match
- Retrieve file from object storage
- Stream file content to response with appropriate headers
