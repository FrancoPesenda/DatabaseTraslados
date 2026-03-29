# Go API (DDD + Clean) + MySQL

## Ejecutar

Configuración (recomendado: variables de entorno con `.env`):

- Copiá `env_configuration` a `.env` y completá valores (NO se commitea).

En Git Bash:

```bash
cp env_configuration .env
```

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
- `cmd/api/dotenv.go`: carga `.env` (sin dependencias externas).
- `config/`: configuración por scope.
- `internal/domain`: entidades + contratos (dominio).
- `internal/usecase`: lógica de negocio (casos de uso).
- `internal/repository`: persistencia (adaptadores a MySQL).
- `internal/handler`: handlers HTTP (deserializa/serializa).
- `docs`: guía para futuros desarrollos.

