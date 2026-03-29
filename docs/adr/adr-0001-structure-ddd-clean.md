# ADR-0001: Estructura DDD + Clean Architecture

## Contexto

Necesitamos una API en Go que se conecte a MySQL y pueda crecer con nuevos casos de uso sin volverse difícil de mantener.

## Decisión

Adoptamos una estructura por capas:

- `internal/domain`: entidades + contratos de repositorio.
- `internal/application`: casos de uso.
- `internal/infrastructure`: detalles de DB y proveedores.
- `internal/interfaces`: HTTP.

## Consecuencias

- El dominio queda aislado de frameworks/DB.
- La infraestructura puede cambiar (MySQL, tests, mocks) sin reescribir el dominio.

