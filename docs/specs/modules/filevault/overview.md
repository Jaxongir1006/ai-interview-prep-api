# Filevault Module

## Purpose

The Filevault module provides centralized file storage and management capabilities for the application. It handles file uploads to object storage, manages file metadata in the database, and provides mechanisms for other modules to attach files to their domain entities.

## Responsibilities

- Accept and validate file uploads (size, content type)
- Store files in object storage with organized path structure
- Track file metadata and storage status in the database
- Provide file download with streaming and ETag-based caching
- Enable other modules to attach/detach files to their entities via the portal interface
- Support content group restrictions for entity-specific file type requirements
- Maintain file association ordering via sort order

## Domain Main Entities

| Entity | Description |
| ------ | ----------- |
| `File` | Represents an uploaded file with its metadata, storage location, and optional entity association |

See [ERD.md](./ERD.md) for entity relationships.

## Configuration

| Field | Type | Default | Description |
| ----- | ---- | ------- | ----------- |
| `max_file_size_mb` | int64 | 10 | Maximum allowed file size for uploads in megabytes |

## Portal Interface

The Filevault module exposes a portal interface for other modules to manage file associations:

| Method | Description |
| ------ | ----------- |
| `Attach` | Associate uploaded files with an entity |
| `Replace` | Atomically replace all files for an entity association |
| `ListByEntity` | List all files attached to an entity |
| `DeleteByEntity` | Soft-delete all files attached to an entity |

Portal methods that modify data (`Attach`, `Replace`) require the caller to lend a transaction via the context for atomicity with the caller's operations.
