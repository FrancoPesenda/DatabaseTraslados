# Casos de Uso

## Crear usuario de empresa (CreateCompanyUser)

**Objetivo:** Registrar un nuevo usuario en el sistema con el rol `company`.

**Actor:** Sistema externo o administrador que consume la API.

**Flujo principal:**
1. El cliente envía `user_name`, `email` y `password`.
2. El sistema valida que los tres campos estén presentes.
3. El sistema asigna automáticamente el rol `company` al nuevo usuario.
4. El sistema persiste el usuario en la base de datos.
5. El sistema retorna el usuario creado con su `id`, `name`, `email` y `role`.

**Reglas de negocio:**
- `user_name` es obligatorio y debe ser único en el sistema.
- `email` es obligatorio y debe ser único en el sistema.
- `password` es obligatoria. No se almacena en texto plano.
- El rol se asigna siempre como `company`, el cliente no puede elegirlo.
- Los espacios al inicio y fin del `name` se eliminan automáticamente.

**Flujos alternativos:**
- Si falta `user_name`, `email` o `password` → error 400.
- Si `user_name` o `email` ya existen en la base → error 500 (constraint violation).

---

## Health Check (Ping)

**Objetivo:** Verificar que la API está operativa y puede recibir tráfico.

**Actor:** Scripts de monitoreo, scripts de inicio, balanceadores de carga.

**Flujo principal:**
1. El cliente envía `GET /ping`.
2. El sistema responde `pong` con status 200.

**Reglas de negocio:** Ninguna. No requiere autenticación ni conexión a la base de datos.
