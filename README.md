# Vulon Backend

API REST en Go para la gestión administrativa de clubes deportivos. Incluye gestión de entidades, jugadores, equipos, eventos, asistencia, pagos, gastos y usuarios, además de autenticación JWT, integración con Google Calendar y reportes en Excel.

## Tecnologías

- Go 1.23
- PostgreSQL con [`lib/pq`](https://github.com/lib/pq)
- [`gorilla/mux`](https://github.com/gorilla/mux)
- JWT con `golang-jwt/jwt/v5`
- Google Calendar API
- Reportes Excel con `excelize`
- Recarga en caliente opcional con [Air](https://github.com/air-verse/air)

## Estructura

```text
.
├── calendar/       # Integración con Google Calendar
├── db/             # Conexión a PostgreSQL
├── helpers/        # Respuestas, tokens y generación de Excel
├── middleware/     # Middleware JWT
├── models/         # Modelos y DTOs
├── routes/         # Handlers HTTP
├── sql/            # Esquema de PostgreSQL
├── validations/    # Validaciones de entrada
├── Postamn/        # Colección de solicitudes Postman
├── main.go         # Registro de rutas y arranque del servidor
└── .air.toml        # Configuración de Air
```

## Requisitos

- Go 1.23 o posterior
- PostgreSQL 14 o posterior
- `psql` para aplicar el esquema
- Air, solo si se desea usar recarga en caliente (`go install github.com/air-verse/air@latest`)

## Configuración

1. Copia `.env.example` a `.env` en la raíz del repositorio.
2. Completa las variables de PostgreSQL y define un secreto privado para `HASH`.
3. Configura Google Calendar si se van a usar las funciones que consultan el calendario.

| Variable | Uso |
| --- | --- |
| `DB_SERVER`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Conexión a PostgreSQL. |
| `DB_SSLMODE` | Modo SSL de PostgreSQL; si se omite, el código usa `disable`. |
| `HASH` | Secreto para firmar y validar tokens JWT. Debe mantenerse privado. |
| `RAILWAY_ENVIRONMENT` | Si está definido, la aplicación no carga `.env` y usa las variables del entorno. |
| `ENVIROMENT` | Modo de credenciales de Google Calendar: `DEV` o `PROD` (el nombre está escrito así en el código). |
| `FILE_LOCATIONS` | Ruta al JSON de credenciales cuando `ENVIROMENT=DEV`. |
| `GOOGLE_CREDENTIALS_JSON` | JSON de credenciales cuando `ENVIROMENT=PROD`. |
| `CALENDAR_ID` | ID del calendario de Google Calendar. |

No subas `.env`, credenciales ni secretos al repositorio.

### Base de datos

Crea una base de datos PostgreSQL y configura sus credenciales en `.env`. Aplica el esquema desde la raíz del proyecto:

```bash
psql -h "$DB_SERVER" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f sql/schema_postgres.sql
```

En Windows PowerShell puedes pasar los valores directamente o configurar `PGHOST`, `PGPORT`, `PGUSER` y `PGDATABASE` antes de ejecutar `psql -f sql/schema_postgres.sql`.

En Railway u otro proveedor, configura las variables en el servicio y ejecuta el mismo esquema contra la base de datos remota. Cuando `RAILWAY_ENVIRONMENT` está definido, el proceso no lee el archivo `.env`.

## Ejecutar

Desde la raíz del repositorio:

```bash
go run .
```

El servidor escucha en el puerto `3050`, por lo que la URL local es `http://localhost:3050`.

Con Air:

```bash
air
```

Compilar para Linux (por ejemplo, para despliegue):

```bash
GOOS=linux GOARCH=amd64 go build -o tmp/pitz-backend .
```

## API

La base de las rutas es `/api`. Las rutas y métodos siguientes reflejan el registro actual de `main.go`; los parámetros entre llaves son variables de ruta. Para ejemplos de cuerpos JSON y solicitudes listas para ejecutar, importa [`Postamn/Pitz Backend.postman_collection.json`](Postamn/Pitz%20Backend.postman_collection.json) en Postman.

### Autenticación y dashboard

| Método | Ruta | Descripción |
| --- | --- | --- |
| POST | `/api/loginSession` | Iniciar sesión y obtener un token. |
| POST* | `/api/passwordRestoration` | Restaurar contraseña. |
| GET* | `/api/home?date=AAAA-MM-DD` | Datos del dashboard para la fecha indicada; requiere `date`. |

### Entidades

| Método | Ruta | Descripción |
| --- | --- | --- |
| GET | `/api/entities` | Listar entidades. |
| GET | `/api/entitiesByUser` | Listar entidades asignadas al usuario autenticado. |
| POST | `/api/newEntity` | Crear entidad. |
| POST | `/api/assignEntity` | Asignar una entidad a un usuario. |

### Jugadores

| Método | Ruta | Descripción |
| --- | --- | --- |
| GET | `/api/players/{entity_uid}` | Listar jugadores de una entidad. |
| GET | `/api/playersStats/{entity_uid}` | Obtener estadísticas de jugadores de una entidad. |
| GET | `/api/playerById/{id}` | Obtener jugador por ID. |
| POST | `/api/players/newPlayer` | Crear jugador. |
| PUT | `/api/players/editPlayer/{id}` | Editar jugador. |
| DELETE | `/api/players` | Eliminar jugador. |

### Equipos y categorías

| Método | Ruta | Descripción |
| --- | --- | --- |
| GET | `/api/teams?entity_uid={entity_uid}` | Listar equipos de una entidad. |
| GET | `/api/teams/{id}` | Obtener equipo por ID. |
| POST | `/api/teams/new` | Crear equipo. |
| POST | `/api/teams/delete/{id}` | Eliminar equipo. |
| GET | `/api/team/assigned/{id}` | Listar jugadores asignados a un equipo. |
| POST | `/api/team/assign` | Asignar jugador a equipo. |
| DELETE | `/api/team/unassign?assignation_uid={id}` | Quitar una asignación. |
| POST | `/api/team/assignBulk` | Asignar varios jugadores. |
| DELETE | `/api/team/unassignBulk` | Quitar varias asignaciones. |
| GET | `/api/teams-categories?entity_uid={entity_uid}` | Listar categorías de una entidad. |
| POST | `/api/teams-categories/new` | Crear categoría. |
| POST | `/api/teams-categories/assign` | Asignar categorías a un equipo. |
| POST | `/api/teams-categories/unassign` | Quitar categorías de un equipo. |

### Eventos y asistencia

| Método | Ruta | Descripción |
| --- | --- | --- |
| GET | `/api/events` | Listar eventos. |
| GET | `/api/events/eventById/{id}` | Obtener evento por ID. |
| POST | `/api/events/newEvent` | Crear evento. |
| PUT | `/api/events/editEvent/{id}` | Editar evento. |
| GET | `/api/events/getEventsByMonth` | Consultar eventos del mes. |
| GET | `/api/eventsTypes` | Listar tipos de evento. |
| GET | `/api/eventsStates` | Listar estados de evento. |
| GET | `/api/asistanceTypes` | Listar tipos de asistencia. |
| GET | `/api/asistanceByPlayerId/{id}` | Consultar asistencia por jugador. |

### Pagos y gastos

| Método | Ruta | Descripción |
| --- | --- | --- |
| GET | `/api/payments` | Consultar pagos del mes. |
| POST | `/api/payments/new-payment` | Registrar pago. |
| GET | `/api/paymentsTypes/{entity_uid}` | Listar tipos de pago de una entidad. |
| GET | `/api/paymentsCategory/{entity_uid}` | Listar categorías de pago de una entidad. |
| POST | `/api/paymentsCategory` | Crear categoría de pago. |
| GET | `/api/payments/paymentById/{id}` | Obtener pago por ID. |
| DELETE | `/api/payments/deletePayment/{id}` | Eliminar pago. |
| POST | `/api/payments/getPaymentsReport` | Generar reporte de pagos en Excel. |
| GET | `/api/expenses` | Consultar gastos del mes. |
| POST | `/api/expenses/new-expense` | Registrar gasto. |
| `Get` | `/api/expenses/{id}` | Obtener gasto por ID. **La ruta está registrada con mayúscula/minúscula incorrecta** y puede no aceptar solicitudes HTTP `GET`; cambiar el método en `main.go` a `GET`. |
| DELETE | `/api/expenses/deleteExpense/{id}` | Eliminar gasto. |
| POST | `/api/expenses/getExpensesReport` | Generar reporte de gastos en Excel. |

### Usuarios

| Método | Ruta | Descripción |
| --- | --- | --- |
| GET | `/api/users/` | Listar usuarios. |
| GET | `/api/users/basics/{id}` | Obtener datos básicos de usuario. |
| GET | `/api/users/{id}` | Obtener usuario por ID. |
| PUT | `/api/updateUserPassword/{id}` | Actualizar contraseña. |

> **Autenticación:** los handlers que requieren sesión esperan `Authorization: Bearer <token>`. La aplicación no aplica el middleware JWT globalmente a todas las rutas; la validación se configura por handler/ruta. Comprueba la colección Postman o el código antes de asumir que un endpoint está protegido. `*` indica método esperado; esas dos rutas no restringen el método HTTP en el registro actual.

## CORS

El servidor permite cualquier origen (`Access-Control-Allow-Origin: *`) y los métodos `GET`, `POST`, `PUT`, `DELETE` y `OPTIONS`. En producción conviene restringir los orígenes permitidos.

## Prueba rápida

Con PostgreSQL configurado, el esquema aplicado y el servidor en ejecución, prueba el inicio de sesión:

```bash
curl -X POST http://localhost:3050/api/loginSession \
  -H "Content-Type: application/json" \
  -d '{"user":"<usuario>","password":"<contraseña>"}'
```

Usa el token recibido como `Authorization: Bearer <token>` en los endpoints que lo requieran.

## Solución de problemas

- **Error de conexión a PostgreSQL:** verifica que el servidor esté disponible y que las variables `DB_*` sean correctas.
- **Error de SSL:** verifica que `DB_SSLMODE` coincida con la configuración del servidor.
- **`pq: relation ... does not exist`:** aplica `sql/schema_postgres.sql` en la base de datos utilizada por la aplicación.
- **JWT inválido o expirado:** revisa `HASH` y reinicia la aplicación si cambias el secreto.
- **Google Calendar no responde:** revisa `ENVIROMENT`, `CALENDAR_ID` y el origen de credenciales configurado (`FILE_LOCATIONS` o `GOOGLE_CREDENTIALS_JSON`), además de los permisos de la cuenta de servicio.
