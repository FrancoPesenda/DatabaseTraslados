# Go API (DDD + Clean) + MySQL

## Ejecutar

Variables de entorno:

- `HTTP_ADDR` (default `:8080`)
- `MYSQL_DSN` (default `root:password@tcp(127.0.0.1:3306)/app?parseTime=true`)

Run:

```bash
go run ./cmd/api
```

Healthcheck:

- `GET /health`

## Estructura

- `cmd/api`: entrypoint.
- `internal/domain`: entidades + contratos.
- `internal/application`: casos de uso (por dominio).
- `internal/infrastructure`: MySQL y persistencia.
- `internal/interfaces/http`: rutas y handlers.
- `docs`: guía para futuros desarrollos.

