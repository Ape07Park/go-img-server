# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
go run ./main/main.go

# Build
go build ./...

# Run tests
go test ./...

# Run a single package's tests
go test ./internal/storage/...

# Tidy dependencies
go mod tidy
```

## Configuration

The server is configured entirely via environment variables (defaults shown):

| Variable     | Default                    | Description                     |
|--------------|----------------------------|---------------------------------|
| `PORT`       | `8080`                     | HTTP listen port                |
| `UPLOAD_DIR` | `./uploads`                | Root directory for stored files |
| `BASE_URL`   | `http://localhost:8080`    | Domain prepended to image URLs  |
| `API_KEY`    | `dev-secret-key`           | Secret for `X-API-Key` header   |

## Architecture

The entry point (`main/main.go`) manually wires all dependencies (no DI framework). The key design decision is the `Storage` interface in `internal/storage/storage.go`, which makes the storage backend swappable — the current `LocalStorage` implementation can be replaced with a MinIO implementation by changing a single line in `main.go`.

### Request flow

```
HTTP Request
  └─ Gin router (main.go)
       ├─ Public:  GET /i/:project/:filename  → ImageHandler.Serve
       └─ /api/v1 group (APIKeyAuth middleware checks X-API-Key header)
            └─ /projects/:project
                 ├─ POST   /images              → UploadHandler.Upload
                 ├─ GET    /images              → ListHandler.List
                 ├─ DELETE /images/:filename    → DeleteHandler.Delete
                 └─ GET    /images/:filename/download → ImageHandler.Download
```

### Storage layout on disk

Files are stored at `{UPLOAD_DIR}/{project}/{unix_nano_timestamp}.{ext}`. The `sanitize()` function in `local.go` strips `..`, `/`, and `\` from both project and filename path segments to prevent path traversal.

### Placeholder packages

`internal/model/dto.go`, `internal/repository/repository.go`, `internal/service/service.go`, and `pkg/response/response.go` are currently empty stubs. The actual data model (`FileInfo`) lives in `internal/storage/storage.go`.

### Upload constraints

- Max file size: 10 MB (enforced in `UploadHandler`)
- Allowed MIME types: `image/jpeg`, `image/png`, `image/gif`, `image/webp` (checked via `Content-Type` header of the multipart part)
- Served images include a 30-day `Cache-Control` header
