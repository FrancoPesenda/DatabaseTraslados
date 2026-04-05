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
- `password` es obligatoria. No se almacena en texto plano (bcrypt).
- El rol se asigna siempre como `company`, el cliente no puede elegirlo.
- Los espacios al inicio y fin del `name` se eliminan automáticamente.

**Flujos alternativos:**
- Si falta `user_name`, `email` o `password` → 400 Bad Request.
- Si `user_name` o `email` ya existen en la base → 500 (constraint violation).

---

## Login (Login)

**Objetivo:** Autenticar un usuario existente verificando sus credenciales.

**Actor:** Usuario registrado que quiere acceder al sistema.

**Flujo principal:**
1. El cliente envía `password` y al menos uno de `user_name` o `email`.
2. El sistema busca al usuario por `user_name` (si está presente) o por `email`.
3. El sistema verifica que la password coincida con el hash almacenado (bcrypt).
4. El sistema retorna los datos del usuario sin la password.

**Reglas de negocio:**
- Se debe proveer al menos `user_name` o `email`. Si no se provee ninguno → 400.
- `password` es obligatoria → 400 si falta.
- Si el usuario no existe o la password no coincide → 401 Unauthorized (no se diferencia entre ambos casos por seguridad).
- La password nunca se retorna en la respuesta.
- Si se proveen ambos (`user_name` y `email`), se usa `user_name` con prioridad.

**Flujos alternativos:**
- Falta `user_name` y `email` → 400 Bad Request.
- Falta `password` → 400 Bad Request.
- Usuario no encontrado o password incorrecta → 401 Unauthorized.

---

## Health Check (Ping)

**Objetivo:** Verificar que la API está operativa y puede recibir tráfico.

**Actor:** Scripts de monitoreo, scripts de inicio, balanceadores de carga.

**Flujo principal:**
1. El cliente envía `GET /ping`.
2. El sistema responde `pong` con status 200.

**Reglas de negocio:** Ninguna. No requiere autenticación ni conexión a la base de datos.
