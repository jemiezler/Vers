# API

## `GET /healthz`

Returns API health.

```json
{ "status": "ok" }
```

## `POST /reviews`

Runs the scaffold review pipeline for a manifest.

```json
{
  "filename": "go.mod",
  "content": "require github.com/gin-gonic/gin v1.10.0"
}
```

