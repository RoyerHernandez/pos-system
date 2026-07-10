# POS System — Go REST API

Go REST API backend for the **Callejon Bar** POS system. Designed to coexist with the existing PHP application on the same MySQL database, and eventually replace it as the sole backend.

## Architecture

```
Clients (Web, Android, iOS)
            |
            v
  Go REST API (port 8080)
  ├── Chi Router + CORS + JWT Middleware
  ├── Handlers   (HTTP request/response)
  ├── Services   (business logic, transactions)
  ├── Repositories (data access via sqlx)
  └── Models     (structs matching DB schema)
            |
            v
     MySQL (shared database)
            ^
            |
  PHP POS App (port 8000)
  (coexistence period — eventually deprecated)
```

### Project Structure

```
api/
├── cmd/server/main.go              # Entry point — wires dependencies, starts server
├── internal/
│   ├── auth/jwt.go                 # JWT generation and validation
│   ├── config/
│   │   ├── config.go               # Environment config loader
│   │   └── database.go             # MySQL connection pool (sqlx)
│   ├── handlers/                   # HTTP handlers (one per entity)
│   │   ├── auth.go
│   │   ├── cashregister.go
│   │   ├── category.go
│   │   ├── client.go
│   │   ├── dashboard.go
│   │   ├── health.go
│   │   ├── inventory.go
│   │   ├── product.go
│   │   ├── reports.go
│   │   ├── response.go            # Shared JSON/error response helpers
│   │   ├── sale.go
│   │   └── user.go
│   ├── middleware/
│   │   ├── auth.go                 # JWT extraction + context injection
│   │   ├── ratelimit.go            # Per-IP rate limiting (in-memory)
│   │   ├── requestid.go            # X-Request-ID generation/validation
│   │   └── role.go                 # Role-based access (RequireRole)
│   ├── models/                     # Structs matching DB tables
│   ├── repositories/               # Data access layer (SQL queries)
│   ├── router/router.go            # Chi router + CORS config
│   ├── services/                   # Business logic layer
│   └── utils/image.go              # Image upload with MIME validation
├── go.mod
├── go.sum
└── .env.example
```

## Tech Stack

| Component | Library | Why |
|-----------|---------|-----|
| Router | [go-chi/chi](https://github.com/go-chi/chi) v5 | Lightweight, stdlib-compatible, middleware-friendly |
| Database | [jmoiron/sqlx](https://github.com/jmoiron/sqlx) | Thin wrapper over database/sql, struct scanning |
| Auth | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) v5 | JWT generation + validation |
| Passwords | golang.org/x/crypto/bcrypt | Compatible with PHP's `$2y$` bcrypt hashes |
| Config | [joho/godotenv](https://github.com/joho/godotenv) | `.env` file loading |
| No ORM | Direct SQL | Full control, performance, no magic |

## Requirements

- Go 1.21+ (developed with 1.26)
- MySQL 8.0+
- Existing POS database (`pos`) with tables already created by the PHP app

## Setup

### 1. Environment Configuration

```bash
cd api/
cp .env.example .env
```

Edit `.env` with your values:

```env
DB_HOST=localhost
DB_PORT=3306
DB_NAME=pos
DB_USER=root
DB_PASS=

SERVER_PORT=8080

# REQUIRED — server will not start without this
# Generate with: openssl rand -base64 48
JWT_SECRET=your-strong-random-secret-here

# Comma-separated allowed origins for CORS
CORS_ORIGINS=http://localhost:8000,http://localhost:8080
```

> **Important:** `JWT_SECRET` is required. The server will refuse to start if it is empty or set to the default placeholder.

### 2. Start the Server

```bash
# From the api/ directory
go run ./cmd/server

# Or from the repo root using dev scripts
./start-dev.sh    # Starts both PHP + Go servers
./stop-dev.sh     # Stops both
```

The API starts on the port configured in `SERVER_PORT` (default: 8080).

### 3. Verify

```bash
curl http://localhost:8080/api/v1/health
# {"status":"ok","database":"connected"}
```

## Authentication

The API uses JWT Bearer tokens with access/refresh token flow.

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"usuario": "admin", "password": "admin123"}'
```

Response:
```json
{
  "access_token": "eyJhbG...",
  "refresh_token": "eyJhbG...",
  "user": { "id": 1, "usuario": "admin", "perfil": "Administrador", ... }
}
```

### Using Tokens

```bash
curl http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <access_token>"
```

### Refresh

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "<refresh_token>"}'
```

- Access tokens expire in **24 hours**
- Refresh tokens expire in **7 days**
- Access tokens cannot be used to refresh (distinct token types enforced)
- Deactivated users cannot refresh tokens

## API Endpoints

### Public (no auth required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/login` | Authenticate user |
| POST | `/api/v1/auth/refresh` | Refresh JWT tokens |
| GET | `/api/v1/health` | Health check + DB status |

### Authenticated (all roles)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/sales` | Create a sale |
| POST | `/api/v1/cashregister/open` | Open cash register shift |
| PUT | `/api/v1/cashregister/{id}/close` | Close cash register shift |
| GET | `/api/v1/cashregister/current` | Get current open register |

### Admin + Especial

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/products` | List products (paginated) |
| GET | `/api/v1/products/{id}` | Get product by ID |
| POST | `/api/v1/products` | Create product |
| PUT | `/api/v1/products/{id}` | Update product (metadata only, not stock) |
| DELETE | `/api/v1/products/{id}` | Soft-delete product |
| GET | `/api/v1/clients` | List clients |
| GET | `/api/v1/clients/{id}` | Get client by ID |
| POST | `/api/v1/clients` | Create client |
| PUT | `/api/v1/clients/{id}` | Update client |
| DELETE | `/api/v1/clients/{id}` | Soft-delete client |
| GET | `/api/v1/sales` | List sales (paginated, filterable) |
| GET | `/api/v1/sales/{id}` | Sale detail with line items |
| GET | `/api/v1/inventory` | List inventory movements |
| POST | `/api/v1/inventory` | Create manual movement (entrada/salida) |
| GET | `/api/v1/cashregister` | List all register shifts |
| GET | `/api/v1/dashboard/kpis` | Dashboard KPIs |
| GET | `/api/v1/reports/sales` | Sales report by date range |
| GET | `/api/v1/reports/products` | Product performance report |
| GET | `/api/v1/reports/clients` | Client purchase rankings |

### Admin Only

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users` | List users |
| GET | `/api/v1/users/{id}` | Get user by ID |
| POST | `/api/v1/users` | Create user |
| PUT | `/api/v1/users/{id}` | Update user |
| DELETE | `/api/v1/users/{id}` | Soft-delete user |
| GET | `/api/v1/categories` | List categories |
| GET | `/api/v1/categories/{id}` | Get category by ID |
| POST | `/api/v1/categories` | Create category |
| PUT | `/api/v1/categories/{id}` | Update category |
| DELETE | `/api/v1/categories/{id}` | Soft-delete category |
| PUT | `/api/v1/sales/{id}/cancel` | Cancel sale (reverses all) |

### Query Parameters

**List endpoints** support pagination:
- `page` — page number (default: 1)
- `per_page` — items per page (default: 25, max: 200)

**Sales list** additional filters:
- `start_date` — filter from date (YYYY-MM-DD)
- `end_date` — filter to date (YYYY-MM-DD)
- `status` — filter by status (`completada`, `cancelada`)
- `user_id` — filter by user

**Reports** support:
- `start_date`, `end_date` — date range filters
- `Accept: text/csv` header — returns CSV instead of JSON

## Sales Pipeline

Sale creation is the most critical flow. It executes **8 steps in a single database transaction**:

1. Validate open cash register (locked with `FOR UPDATE`)
2. Lock and validate stock for each product (`FOR UPDATE`)
3. Generate sequential sale code (`VTA-NNNNNN`, serialized)
4. Calculate totals (subtotal, tax, discount)
5. Insert sale record
6. For each product: insert detail line, decrement stock, record inventory movement
7. Update client purchase totals
8. Update cash register totals (by payment method)

**Cancellation** reverses all changes atomically:
- Restores product stock (with row locks)
- Creates reversal inventory movements
- Reverses client purchase totals
- Updates sale status to `cancelada`

### Payment Methods

Accepted values: `Efectivo`, `Tarjeta`, `Transferencia`

### Inventory Movements

- **Sale-driven**: automatic `salida` on sale, `entrada` on cancellation
- **Manual**: `entrada` (purchase, adjustment) or `salida` (waste, adjustment)
- All movements linked to products and users, with optional reference IDs

## Coexistence with PHP App

The Go API shares the same MySQL database with the PHP POS application. Key design decisions for safe coexistence:

- **Same DB schema**: Go models match existing table structures exactly
- **Bcrypt compatibility**: Go reads PHP-generated `$2y$` bcrypt hashes
- **Soft-delete convention**: Go uses `estado = 0` for deletes, matching PHP
- **Image paths**: uploads go to `../pos/views/img/` (same directory PHP uses)
- **Row-level locking**: `SELECT ... FOR UPDATE` prevents stock inconsistencies when both apps write concurrently
- **Sale code serialization**: prevents duplicate `VTA-NNNNNN` codes across both apps

## Security

- JWT with distinct access/refresh token types
- Bcrypt password hashing (cost 10)
- Role-based access control (Administrador, Especial, Vendedor)
- Per-IP rate limiting (100 req/min)
- CORS with configurable allowed origins
- Image upload MIME validation (magic bytes, not just extension)
- Path traversal protection on file uploads
- Sanitized error messages (no internal details leaked)
- Request ID validation (alphanumeric, max 64 chars)
- Graceful shutdown with signal handling

## Error Response Format

All errors follow a consistent structure:

```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "human-readable description"
  }
}
```

HTTP status codes: 400 (bad request), 401 (unauthorized), 403 (forbidden), 404 (not found), 500 (internal error).
