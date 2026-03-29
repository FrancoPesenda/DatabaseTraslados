# Documentación (guía viva)

Este directorio es la **fuente de verdad** sobre qué debe hacer la aplicación y cómo evolucionarla sin romper el diseño.

## ¿Qué es esta app?

- API en Go para conectar un backend con una base MySQL (modelada en MySQL Workbench).
- Arquitectura: **DDD + Clean Architecture**.

## Cómo leer estos docs

- `architecture.md`: visión general de capas, dependencias y convenciones.
- `domain/`: definición de conceptos por dominio (lenguaje ubicuo).
- `adr/`: decisiones de arquitectura (registros cortos, con contexto).

## Próximos documentos sugeridos

- Contratos de endpoints (OpenAPI o colección Postman).
- Esquema de base de datos (tablas, constraints, índices) y cómo se versiona con migraciones.

