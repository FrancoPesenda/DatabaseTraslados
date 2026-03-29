# Go API (DDD + Clean) + MySQL

## Ejecutar

Configuración por archivo JSON:

- `config/<scope>/infrastructure_config.json`
- Scope por defecto: `local`
- Podés cambiarlo con `CONFIG_SCOPE` (ej: `dev`, `prod`)

Run:

```bash
go run ./cmd/api
```

Healthcheck:

- `GET /health`

## Estructura

- `cmd/api`: entrypoint.
- `cmd/api/factory.go`: inyección de dependencias y construcción de la app.
- `cmd/api/config.go`: lectura de `config/<scope>/infrastructure_config.json`.
- `config/`: configuración por scope.
- `internal/domain`: entidades + contratos (dominio).
- `internal/usecase`: lógica de negocio (casos de uso).
- `internal/repository`: persistencia (adaptadores a MySQL).
- `internal/handler`: handlers HTTP (deserializa/serializa).
- `docs`: guía para futuros desarrollos.

