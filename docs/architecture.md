# Arquitectura

## Estructura de capas (DDD + Clean Architecture)

```
cmd/
├── api/                        → Entry point: DI (Uber FX), HTTP server, router
└── handler/
    └── user/
        ├── createcompany/      → Handler: crear usuario company
        └── login/              → Handler: autenticar usuario

internal/
├── domain/
│   └── user.go                 → Entidad User, roles, errores de dominio, validaciones
├── usecase/
│   └── user/
│       ├── createcompany/      → Lógica: validar, asignar rol, persistir
│       └── login/              → Lógica: validar credenciales, verificar password (bcrypt)
├── repository/
│   └── eventradatabase/        → Acceso a MySQL: INSERT, SELECT by ID, SELECT by username/email
├── utils/
│   ├── security/               → HashPassword / CheckPassword (bcrypt)
│   └── config/                 → Configuración por env vars o archivo JSON
```

## Flujo de una request — CreateCompanyUser

```
POST /user/company
    │
    ▼
Handler (createcompany)
    │  decode JSON → domain.User
    ▼
UseCase (createcompany)
    │  TrimSpace(name)
    │  validar user_name, email, password
    │  setRole(company)
    ▼
Repository (eventradatabase)
    │  INSERT INTO user
    │  SELECT user by id
    ▼
MySQL
    │
    ▼
UseCase → Handler → 201 Created { id, name, email, role }
```

## Flujo de una request — Login

```
POST /user/login
    │
    ▼
Handler (login)
    │  decode JSON → domain.User
    ▼
UseCase (login)
    │  validar: user_name o email + password
    │  GetUserByUserName o GetUserByEmail
    │  CheckPassword(storedHash, plainPassword)
    │  limpiar password del resultado
    ▼
Repository (eventradatabase)
    │  SELECT user WHERE username = ? / email = ?
    ▼
MySQL
    │
    ▼
UseCase → Handler → 200 OK { id, user_name, email, role }
                  → 401 Unauthorized si credenciales inválidas
```

## Regla de dependencias

```
Handler → UseCase → Repository → Domain
   │          │                     ↑
   │          └── utils/security    │
   │                                │
   └── (nunca importa domain directamente salvo para errores)
```

- El dominio no importa ningún paquete externo.
- Las interfaces se definen donde se usan.
- El mapeo entre capas lo hace siempre la capa exterior (`toDomain()` en request, `NewResponse()` desde dominio).

## Inyección de dependencias (Uber FX)

```
*sql.DB
    └── *Repository
            ├── (como createcompany.UserRepository) → *createcompany.UseCase → *CreateHandler
            └── (como login.UserRepository)         → *login.UseCase        → *login.Handler
                                                                                      │
                                                                               http.Handler
                                                                                      │
                                                                               *http.Server
```

## Infraestructura local

```
Docker Compose
├── db   → MySQL 8.0 en puerto 3307 (host) / 3306 (interno)
└── api  → Go API en puerto 8080, restart: on-failure

Cloudflare Tunnel → expone localhost:8080 a internet (URL pública temporal)
```

## Scripts disponibles

| Script | Descripción |
|--------|-------------|
| `bash scripts/start.sh` | Levanta Docker + cloudflared |
| `bash scripts/stop.sh` | Baja todos los contenedores |
| `bash scripts/reload.sh` | Reconstruye todo desde cero (borra datos) |
| `bash scripts/logs.sh` | Logs en tiempo real de la API |
