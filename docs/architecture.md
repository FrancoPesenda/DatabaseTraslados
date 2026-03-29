# Arquitectura (DDD + Clean)

## Capas

- **Domain (`internal/domain/`)**: entidades y reglas del negocio. Solo depende de Go (stdlib).
- **Application (`internal/application/`)**: casos de uso (orquesta el dominio). Depende de `domain`.
- **Infrastructure (`internal/infrastructure/`)**: DB MySQL, repositorios concretos, proveedores externos.
- **Interfaces (`internal/interfaces/`)**: HTTP handlers, rutas, middlewares (entrada/salida).

## Regla de dependencias

Las dependencias deben apuntar hacia adentro:

`interfaces` → `application` → `domain`

`infrastructure` implementa contratos definidos hacia adentro (normalmente interfaces en `domain` o `application`).

## Convenciones de paquetes

- 1 carpeta por dominio: `user`, `service`, `service_type`, `event`, `payment_method`.
- Los IDs se tipan (no `string` suelto) para evitar mezclar identificadores.

