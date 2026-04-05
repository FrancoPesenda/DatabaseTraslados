# Arquitectura

## Estructura de capas (DDD + Clean Architecture)

```
cmd/
├── api/                        → Entry point: DI (Uber FX), HTTP server, router
└── handler/
    └── user/
        └── createcompany/      → HTTP handler: decode request, call use case, encode response

internal/
├── domain/
│   └── user.go                 → Entidad User, roles, errores de dominio, validaciones
├── usecase/
│   └── user/
│       └── createcompany/      → Lógica de negocio: validar, asignar rol, llamar al repositorio
├── repository/
│   └── eventradatabase/        → Acceso a MySQL: INSERT y SELECT
└── config/                     → Configuración por env vars o archivo JSON
```

## Flujo de una request

```
HTTP Request
    │
    ▼
Handler (cmd/handler)
    │  decode JSON → domain.User
    ▼
Use Case (internal/usecase)
    │  validar campos
    │  asignar rol company
    ▼
Repository (internal/repository)
    │  INSERT INTO user
    │  SELECT user by id
    ▼
MySQL (Docker)
    │
    ▼
Use Case → Handler → HTTP Response (JSON)
```

## Regla de dependencias

Las dependencias apuntan siempre hacia adentro:

```
Handler → UseCase → Repository → Domain
                              ↑
                         (solo stdlib)
```

- El dominio no importa ningún paquete externo.
- Las interfaces se definen donde se usan (handler define su propia interfaz de UseCase, usecase define su propia interfaz de Repository).
- El mapeo entre capas lo hace siempre la capa exterior (`toDomain()` en request, `NewResponse()` desde dominio).

## Inyección de dependencias

Se usa [Uber FX](https://github.com/uber-go/fx). Cada constructor recibe sus dependencias como parámetros. El grafo de dependencias se resuelve en `cmd/api/fxapp.go`.

```
*sql.DB
    └── *Repository
            └── (como createcompany.UserRepository)
                    └── *UseCase
                            └── (como user.UseCase)
                                    └── *CreateHandler
                                                └── http.Handler
                                                        └── *http.Server
```

## Infraestructura local

```
Docker Compose
├── db   → MySQL 8.0 en puerto 3307 (host) / 3306 (interno)
└── api  → Go API en puerto 8080

Cloudflare Tunnel → expone localhost:8080 a internet (URL pública temporal)
```
