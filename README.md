# vulon_backend

API REST en Go para la gestión administrativa de clubes deportivos. Permite administrar jugadores, eventos, asistencia, pagos, gastos y usuarios, con autenticación por JWT y reportes en Excel.

## Stack

- **Lenguaje:** Go 1.23
- **Router:** [gorilla/mux](https://github.com/gorilla/mux)
- **Base de datos:** PostgreSQL (driver [`lib/pq`](https://github.com/lib/pq))
- **Autenticación:** JWT (`golang-jwt/jwt/v5`)
- **Hot reload:** [Air](https://github.com/air-verse/air) (config en `.air.toml`)
- **Reportes:** [xuri/excelize](https://github.com/xuri/excelize)
- **Variables de entorno:** [joho/godotenv](https://github.com/joho/godotenv)

## Estructura del proyecto

```
backend/
├── main.go                  # Punto de entrada y registro de rutas
├── calendar/                # Helpers de fechas
├── db/                      # Conexión a PostgreSQL
├── helpers/                 # Generadores (Excel, token, jugadores)
├── middleware/              # Middleware JWT
├── models/                  # Estructuras de datos (DTOs)
├── routes/                  # Handlers HTTP por dominio
├── sql/                     # Esquema de base de datos (schema_postgres.sql)
├── validations/             # Validaciones de entrada
├── tmp/                     # Binarios de Air
└── .air.toml                # Configuración de hot reload
```

## Variables de entorno

Copiar `.env.example` como `.env` en la raíz del proyecto y completar la configuración local. El archivo `.env` está excluido de Git; no subir credenciales ni secretos al repositorio. En Railway u otro proveedor, definir estas variables en la configuración del servicio.

| Variable | Descripción |
| --- | --- |
| `DB_USER`, `DB_PASSWORD` | Usuario y contraseña de PostgreSQL. |
| `DB_SERVER`, `DB_PORT`, `DB_NAME` | Host, puerto y base de datos PostgreSQL. |
| `DB_SSLMODE` | Modo SSL de la conexión. |
| `HASH` | Secreto utilizado para firmar y validar tokens JWT. |
| `RAILWAY_ENVIRONMENT` | Cuando está definido, la aplicación no carga `.env` y usa el entorno del proveedor. |
| `ENVIROMENT` | Selecciona el modo de credenciales de Google Calendar: archivo local o JSON en una variable. El nombre está escrito así en el código. |
| `FILE_LOCATIONS` | Ruta local al JSON de credenciales de Google Calendar cuando se usa el modo de archivo. |
| `GOOGLE_CREDENTIALS_JSON` | JSON de credenciales de Google Calendar cuando se usa el modo de variable. |
| `CALENDAR_ID` | Identificador del calendario de Google Calendar. |

## Configuración de la base de datos

### PostgreSQL local

Instalar PostgreSQL 14 o posterior (localmente o mediante Docker) y crear una base de datos y un usuario. Configurar sus datos de conexión en `.env`.

### Aplicar el esquema

Una vez la BD sea accesible, ejecutar el script de creación de tablas:

```bash
psql -f sql/schema_postgres.sql
```

Ejecuta el comando con `psql` configurado para conectarse a la base de datos del proyecto.

### Railway / proveedor cloud

1. Crear un servicio PostgreSQL y configurar en el backend las variables de conexión indicadas arriba, usando los datos proporcionados por el proveedor.
2. Aplicar `sql/schema_postgres.sql` desde una consola `psql` conectada a esa base de datos.

## Cómo correrlo

### Local con hot reload (Air)

```bash
air
```

El binario se recompila en cada cambio y el servidor queda escuchando en `http://localhost:3050`.

### Local sin hot reload

```bash
go build -o tmp/main.exe .
./tmp/main.exe
```

### Compilación cruzada (Linux para deploy)

```bash
GOOS=linux GOARCH=amd64 go build -o tmp/pitz-backend .
```

## Endpoints

Base URL: `/api`. Los endpoints protegidos validan un token JWT; `loginSession` y `passwordRestoration` son públicos.

### Autenticación
| Método | Endpoint                        | Descripción                       |
| ------ | ------------------------------- | --------------------------------- |
| POST   | `/api/loginSession`             | Iniciar sesión                    |
| POST   | `/api/passwordRestoration`      | Restaurar contraseña              |

### Home
| Método | Endpoint       | Descripción           |
| ------ | -------------- | --------------------- |
| GET    | `/api/home`    | Datos del dashboard   |

### Jugadores
| Método | Endpoint                                       | Descripción                |
| ------ | ---------------------------------------------- | -------------------------- |
| GET    | `/api/players`                                 | Listar jugadores           |
| GET    | `/api/players/{id}`                            | Obtener jugador por ID     |
| POST   | `/api/players/newPlayer`                       | Crear jugador              |
| PUT    | `/api/players/editPlayer/{id}`                 | Editar jugador             |
| DELETE | `/api/players`                                 | Eliminar jugador           |

### Eventos
| Método | Endpoint                                | Descripción                 |
| ------ | --------------------------------------- | --------------------------- |
| GET    | `/api/events`                           | Listar eventos              |
| GET    | `/api/events/eventById/{id}`            | Evento por ID               |
| POST   | `/api/events/newEvent`                  | Crear evento                |
| PUT    | `/api/events/editEvent/{id}`            | Editar evento               |
| GET    | `/api/events/getEventsByMonth`          | Eventos por mes             |

### Tipos y asistencia
| Método | Endpoint                                | Descripción                     |
| ------ | --------------------------------------- | ------------------------------- |
| GET    | `/api/eventsTypes`                      | Tipos de evento                 |
| GET    | `/api/asistanceTypes`                   | Tipos de asistencia             |
| GET    | `/api/asistanceByPlayerId/{id}`         | Asistencia por jugador          |

### Pagos
| Método | Endpoint                                | Descripción                       |
| ------ | --------------------------------------- | --------------------------------- |
| GET    | `/api/payments`                         | Pagos del mes                     |
| POST   | `/api/payments/new-payment`             | Registrar pago                    |
| GET    | `/api/paymentsTypes`                    | Tipos de pago                     |
| GET    | `/api/payments/paymentById/{id}`        | Pago por ID                       |
| DELETE | `/api/payments/deletePayment/{id}`      | Eliminar pago                     |
| POST   | `/api/payments/getPaymentsReport`       | Reporte en Excel                  |

### Gastos
| Método | Endpoint                                | Descripción                       |
| ------ | --------------------------------------- | --------------------------------- |
| GET    | `/api/expenses`                         | Gastos del mes                    |
| POST   | `/api/expenses/new-expense`             | Registrar gasto                   |
| GET    | `/api/expenses/{id}`                    | Gasto por ID                      |
| DELETE | `/api/expenses/deleteExpense/{id}`      | Eliminar gasto                    |
| POST   | `/api/expenses/getExpensesReport`       | Reporte en Excel                  |

### Usuarios
| Método | Endpoint                           | Descripción                |
| ------ | ---------------------------------- | -------------------------- |
| GET    | `/api/users/`                      | Listar usuarios            |
| GET    | `/api/users/basics/{id}`           | Datos básicos de usuario   |
| GET    | `/api/users/{id}`                  | Usuario por ID             |
| PUT    | `/api/updateUserPassword/{id}`     | Actualizar contraseña      |

## CORS

El servidor responde con los headers `Access-Control-Allow-Origin: *` para los métodos `GET, POST, PUT, DELETE, OPTIONS`. En producción conviene restringir el origen.

## Notas

- El servidor escucha en el puerto `3050` (ver `main.go`).
- La conexión a la base de datos se cierra al finalizar (`defer db.CerrarConexion()`).
- El driver `lib/pq` devuelve los campos de tipo fecha/timestamp nativamente como `time.Time` (no requiere `parseTime`).
- Los placeholders SQL usan `$1`, `$2`, … (no `?`); los identificadores con uso de mayúsculas están entre comillas dobles (`"table"`).

## Diferencias clave Postgres vs MySQL

| Concepto | MySQL | PostgreSQL (este proyecto) |
| --- | --- | --- |
| Driver Go | `github.com/go-sql-driver/mysql` | `github.com/lib/pq` |
| DSN | `user:pass@tcp(host:port)/db?parseTime=true` | `host=… port=… user=… password=… dbname=… sslmode=…` |
| Placeholders | `?` | `$1`, `$2`, … |
| Identificadores | `` `tabla` `` | `"tabla"` (sin comillas si son solo minúsculas) |
| `AUTO_INCREMENT` | Sí | Usar `uuid_generate_v4()` o serial del proyecto |
| Timestamp literal | `current_timestamp()` | `CURRENT_TIMESTAMP` (sin paréntesis) |
| Formato fecha | `DATE_FORMAT(d, '%m-%Y')` | `TO_CHAR(d, 'MM-YYYY')` |
| Soft deletes | `delete_flag` TINYINT | `delete_flag` SMALLINT (idéntico en queries) |
| Booleanos | `TINYINT(1)` | `BOOLEAN` nativo |

## Smoke test

Con PostgreSQL configurado, el esquema aplicado y un usuario de prueba creado, iniciar el servidor y probar el login. Sustituir los marcadores por las credenciales del usuario de prueba:

```bash
curl -X POST http://localhost:3050/api/loginSession \
  -H "Content-Type: application/json" \
  -d '{"user":"<usuario>","password":"<contraseña>"}'
```

Usar el token devuelto para probar un endpoint protegido, por ejemplo `GET /api/players`.

## Troubleshooting

- **Error de SSL** — verificar que `DB_SSLMODE` coincida con la configuración del servidor PostgreSQL.
- **Error de autenticación** — verificar que los datos de conexión de `.env` coincidan con los de PostgreSQL.
- **Error de conexión** — verificar que PostgreSQL esté iniciado y accesible desde el backend.
- **`pq: relation "users" does not exist`** — no se ha aplicado el esquema. Ejecutar `psql … -f sql/schema_postgres.sql`.
- **Tokens JWT no válidos / expirados** — verificar la variable `HASH` y reiniciar el servidor tras cambiarla.
- **Google Calendar falla** — revisa `CALENDAR_ID`, `FILE_LOCATIONS` (o `GOOGLE_CREDENTIALS_JSON`) y permisos de la cuenta de servicio.
