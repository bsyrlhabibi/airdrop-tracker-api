# Airdrop Tracker API

A RESTful backend API for managing cryptocurrency airdrop opportunities across multiple sybil accounts. Track airdrops, tasks, wallets, gas spending, and export everything to Excel — all from one place.

## Tech Stack

- **Language:** Go 1.25+
- **HTTP Framework:** Gin
- **ORM:** GORM
- **Database:** SQLite (file-based, zero-config)
- **Auth:** JWT (golang-jwt) + bcrypt password hashing
- **Export:** Excelize (styled Excel `.xlsx` export)
- **Docs:** Swagger UI via swaggo/gin-swagger

## Project Structure

```
airdrop-tracker-api/
├── cmd/server/main.go          # App entrypoint — loads config, connects DB, starts server
├── internal/
│   ├── config/config.go        # Loads .env into Config struct
│   ├── database/database.go    # GORM connection, auto-migration, legacy schema fixes
│   ├── handler/                # HTTP handlers (one file per resource)
│   │   ├── auth.go             # Register & Login
│   │   ├── airdrop.go          # CRUD for global airdrops
│   │   ├── account.go          # CRUD for sybil accounts + assign/remove airdrops
│   │   ├── category.go         # CRUD for task categories
│   │   ├── airdrop_task.go     # CRUD for global airdrop tasks (templates)
│   │   ├── task.go             # CRUD for per-account tasks + today/date queries
│   │   ├── wallet.go           # Create/list/delete wallets
│   │   ├── dashboard.go        # Aggregated stats summary
│   │   ├── export.go           # Multi-sheet styled Excel export
│   │   └── utils.go            # Shared helpers (parseDate)
│   ├── middleware/auth.go      # JWT Bearer auth middleware
│   ├── model/                  # GORM model definitions
│   │   ├── user.go
│   │   ├── account.go
│   │   ├── airdrop.go
│   │   ├── airdrop_task.go
│   │   ├── account_airdrop.go
│   │   ├── task.go
│   │   ├── category.go
│   │   └── wallet.go
│   ├── repository/             # Data-access layer (DB queries)
│   │   ├── user_repo.go
│   │   ├── account_repo.go
│   │   ├── airdrop_repo.go
│   │   ├── airdrop_task_repo.go
│   │   ├── account_airdrop_repo.go
│   │   ├── task_repo.go
│   │   ├── category_repo.go
│   │   └── wallet_repo.go
│   └── router/router.go        # All route registration, CORS, Swagger setup
├── docs/                        # Auto-generated Swagger docs
├── .env.example                 # Environment variable template
├── go.mod
└── go.sum
```

## Environment Variables

| Variable     | Default              | Description                       |
|-------------|----------------------|-----------------------------------|
| `APP_PORT`  | `8080`               | Server listen port                |
| `DB_PATH`   | `data/airdrop.db`    | SQLite database file path         |
| `JWT_SECRET`| `default-secret`     | Secret key for signing JWT tokens |

## Installation & Running

### Prerequisites

- Go 1.25 or later
- A C compiler (required by `mattn/go-sqlite3`)

### Steps

```bash
# 1. Clone the repository
git clone <repo-url>
cd airdrop-tracker-api

# 2. Copy and edit environment config
cp .env.example .env
# Edit .env — at minimum change JWT_SECRET to a random string

# 3. Install dependencies
go mod download

# 4. Create the data directory
mkdir -p data

# 5. Run the server
go run ./cmd/server/main.go
```

The server starts on the port specified in `APP_PORT` (default `8080`).  
Swagger UI is available at: `http://localhost:8080/swagger/index.html`

### Build a binary

```bash
go build -o airdrop-server ./cmd/server
./airdrop-server
```

## Authentication

All endpoints except **Register** and **Login** require a JWT token in the `Authorization` header.

```
Authorization: Bearer <token>
```

Tokens expire after **72 hours**. The JWT payload contains `user_id` (uint).

---

## API Endpoints

### Auth (Public)

#### Register

```
POST /api/auth/register
```

**Request body:**
```json
{
  "email": "user@email.com",
  "password": "secret123",
  "name": "Bita"
}
```

**Response `201 Created`:**
```json
{
  "message": "Registered",
  "user": { "id": 1, "email": "user@email.com", "name": "Bita" }
}
```

**Error responses:**
- `400` — Validation error
- `409` — Email already exists

---

#### Login

```
POST /api/auth/login
```

**Request body:**
```json
{
  "email": "user@email.com",
  "password": "secret123"
}
```

**Response `200 OK`:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { "id": 1, "email": "user@email.com", "name": "Bita" }
}
```

**Error responses:**
- `400` — Validation error
- `401` — Invalid credentials

---

### Airdrops (Global Catalog)

#### List Airdrops

```
GET /api/airdrops
```

**Response `200 OK`:** Array of `Airdrop` objects.

---

#### Create Airdrop

```
POST /api/airdrops
```

**Request body:**
```json
{
  "name": "zkSync",
  "chain": "Ethereum",
  "category": "rumored",
  "priority": "high",
  "status": "active",
  "url": "https://zksync.io",
  "date_start": "2025-01-01",
  "date_end": "2025-12-31",
  "notes": "Bridge weekly"
}
```

| Field        | Required | Default    | Values                               |
|-------------|----------|------------|--------------------------------------|
| `name`      | Yes      | —          | Any string                           |
| `chain`     | Yes      | —          | e.g. `Ethereum`, `Arbitrum`, `Base`  |
| `category`  | No       | `"rumored"`| `rumored`, `confirmed`, `ended`      |
| `priority`  | No       | `"medium"` | `low`, `medium`, `high`              |
| `status`    | No       | `"active"` | `active`, `upcoming`, `end`, `missed`|
| `url`       | No       | `""`       | URL string                           |
| `date_start`| No       | `null`     | `YYYY-MM-DD`                         |
| `date_end`  | No       | `null`     | `YYYY-MM-DD`                         |
| `notes`     | No       | `""`       | Free text                            |

**Response `201 Created`:** The created `Airdrop` object.

---

#### Get Airdrop

```
GET /api/airdrops/:id
```

**Response `200 OK`:** Single `Airdrop` object.  
**Error `404`:** Airdrop not found.

---

#### Update Airdrop

```
PUT /api/airdrops/:id
```

**Request body:** Partial — only include fields you want to change. Empty strings are ignored.

**Response `200 OK`:** Updated `Airdrop` object.  
**Error `404`:** Airdrop not found.

---

#### Delete Airdrop

```
DELETE /api/airdrops/:id
```

**Response `200 OK`:**
```json
{ "message": "Deleted" }
```

---

### Airdrop Tasks (Global Templates)

These are checklist items that belong to a global Airdrop. When an airdrop is assigned to an account, these tasks are automatically synced.

#### List Airdrop Tasks

```
GET /api/airdrops/:id/tasks
```

**Response `200 OK`:** Array of `AirdropTask` objects.

---

#### Create Airdrop Task

```
POST /api/airdrops/:id/tasks
```

**Request body:**
```json
{
  "name": "Bridge 0.1 ETH",
  "category_id": 1,
  "status": "pending",
  "start_date": "2025-01-15",
  "end_date": "2025-06-01"
}
```

| Field        | Required | Default     | Values                                     |
|-------------|----------|-------------|--------------------------------------------|
| `name`      | Yes      | —           | Task description                           |
| `category_id`| No      | `null`      | FK to Category                             |
| `status`    | No       | `"pending"` | `pending`, `ongoing`, `finish`, `missed`   |
| `start_date`| No       | `null`      | `YYYY-MM-DD`                               |
| `end_date`  | No       | `null`      | `YYYY-MM-DD`                               |

**Response `201 Created`:** The created `AirdropTask` object.

---

#### Update Airdrop Task

```
PUT /api/airdrop-tasks/:id
```

**Request body:** Partial update — only non-empty fields are applied.

**Response `200 OK`:** Updated `AirdropTask` object.  
**Error `404`:** Task not found.

---

#### Delete Airdrop Task

```
DELETE /api/airdrop-tasks/:id
```

**Response `200 OK`:**
```json
{ "message": "Deleted" }
```

---

### Accounts (Sybil Accounts)

#### List Accounts

```
GET /api/accounts
```

**Response `200 OK`:** Array of `Account` objects (with `wallets` and `account_airdrops` if loaded).

---

#### Create Account

```
POST /api/accounts
```

**Request body:**
```json
{
  "name": "Akun 1",
  "color": "#3B82F6",
  "notes": "Main sybil account"
}
```

| Field   | Required | Default     | Description           |
|---------|----------|-------------|-----------------------|
| `name`  | Yes      | —           | Account name          |
| `color` | No       | `"#3B82F6"` | Hex color code        |
| `notes` | No       | `""`        | Free text notes       |

**Response `201 Created`:** The created `Account` object.

---

#### Get Account

```
GET /api/accounts/:id
```

Includes `wallets` and `account_airdrops` relations.  
**Response `200 OK`:** Single `Account` object.  
**Error `404`:** Account not found.

---

#### Update Account

```
PUT /api/accounts/:id
```

**Request body:** Partial update.

**Response `200 OK`:** Updated `Account` object.  
**Error `404`:** Account not found.

---

#### Delete Account

```
DELETE /api/accounts/:id
```

| Query Param | Type | Description                                |
|-------------|------|--------------------------------------------|
| `force`     | bool | `?force=true` deletes account + all data   |

Without `force`, returns `409` if the account has wallets or airdrops.

**Response `200 OK`:**
```json
{ "message": "Account deleted" }
```
```json
{ "message": "Account and all related data deleted" }
```

**Error `409`:** Account has linked data (without `?force=true`).

---

### Account → Airdrop Assignment

#### Assign Airdrop to Account

```
POST /api/accounts/:id/airdrops
```

**Request body:**
```json
{
  "airdrop_id": 1,
  "notes": "Focus on bridging"
}
```

When assigned, all global `AirdropTask` templates for that airdrop are automatically synced as per-account `Task` records.

**Response `201 Created`:** The created `AccountAirdrop` object.  
**Error `409`:** Airdrop already assigned to this account.

---

#### Remove Airdrop from Account

```
DELETE /api/accounts/:id/airdrops/:airdrop_id
```

**Response `200 OK`:**
```json
{ "message": "Airdrop removed from account" }
```

---

### Tasks (Per Account-Airdrop)

Tasks belong to an `AccountAirdrop` and track per-account progress.

#### Create Task

```
POST /api/account-airdrops/:id/tasks
```

**Request body:**
```json
{
  "name": "Bridge 0.1 ETH",
  "category_id": 1,
  "status": "pending",
  "frequency": "daily",
  "date": "2025-01-15"
}
```

| Field        | Required | Default     | Values                                    |
|-------------|----------|-------------|-------------------------------------------|
| `name`      | Yes      | —           | Task description                          |
| `category_id`| No      | `null`      | FK to Category                            |
| `status`    | No       | `"pending"` | `pending`, `ongoing`, `finish`, `missed`  |
| `frequency` | No       | `"once"`    | `once`, `daily`, `weekly`, `monthly`      |
| `date`      | No       | today       | `YYYY-MM-DD`                              |

**Response `201 Created`:** The created `Task` object.

---

#### Update Task

```
PUT /api/tasks/:id
```

**Request body:** Partial update. Additional fields:

| Field       | Type    | Description          |
|------------|---------|----------------------|
| `gas_spent`| float64 | Gas cost in USD      |
| `tx_hash`  | string  | Transaction hash     |

**Response `200 OK`:** Updated `Task` object.  
**Error `404`:** Task not found.

---

#### Delete Task

```
DELETE /api/tasks/:id
```

**Response `200 OK`:**
```json
{ "message": "Deleted" }
```

---

#### Get Today's Tasks

```
GET /api/accounts/:id/tasks/today
```

Returns all tasks scheduled for today across all airdrops in the specified account.

**Response `200 OK`:** Array of `Task` objects.

---

#### Get Tasks by Date

```
GET /api/accounts/:id/tasks/by-date?date=2025-03-15
```

| Query Param | Required | Default | Description         |
|-------------|----------|---------|---------------------|
| `date`      | No       | today   | `YYYY-MM-DD` format |

**Response `200 OK`:** Array of `Task` objects.

---

### Categories

#### List Categories

```
GET /api/categories
```

**Response `200 OK`:** Array of `Category` objects.

---

#### Create Category

```
POST /api/categories
```

**Request body:**
```json
{
  "name": "Bridge",
  "color": "#3B82F6"
}
```

| Field   | Required | Default     |
|---------|----------|-------------|
| `name`  | Yes      | —           |
| `color` | No       | `"#6B7280"` |

**Response `201 Created`:** The created `Category` object.

---

#### Update Category

```
PUT /api/categories/:id
```

**Request body:** Partial update (`name`, `color`).

**Response `200 OK`:** Updated `Category` object.

---

#### Delete Category

```
DELETE /api/categories/:id
```

**Response `200 OK`:**
```json
{ "message": "Deleted" }
```

---

### Wallets

#### List Wallets

```
GET /api/wallets
GET /api/wallets?account_id=1
```

| Query Param  | Required | Description             |
|-------------|----------|-------------------------|
| `account_id`| No       | Filter by Account ID    |

**Response `200 OK`:** Array of `Wallet` objects.

---

#### Create Wallet

```
POST /api/wallets
```

**Request body:**
```json
{
  "account_id": 1,
  "label": "Main Wallet",
  "address": "0x7245...139b",
  "chain": "Ethereum"
}
```

**Response `201 Created`:** The created `Wallet` object.

---

#### Delete Wallet

```
DELETE /api/wallets/:id
```

**Response `200 OK`:**
```json
{ "message": "Deleted" }
```

---

### Dashboard

#### Get Dashboard Summary

```
GET /api/dashboard
```

Returns aggregated statistics with a per-account breakdown.

**Response `200 OK`:**
```json
{
  "total_airdrops": 12,
  "active_airdrops": 8,
  "upcoming_airdrops": 2,
  "ended_airdrops": 1,
  "missed_airdrops": 1,
  "total_tasks": 45,
  "completed_tasks": 20,
  "pending_tasks": 18,
  "ongoing_tasks": 5,
  "missed_tasks": 2,
  "total_wallets": 6,
  "total_accounts": 3,
  "accounts": [
    {
      "id": 1,
      "name": "Akun 1",
      "color": "#3B82F6",
      "total_airdrops": 5,
      "active_airdrops": 4,
      "total_tasks": 15,
      "completed_tasks": 8,
      "pending_tasks": 5,
      "ongoing_tasks": 1,
      "missed_tasks": 1,
      "total_wallets": 2
    }
  ]
}
```

---

### Export

#### Export to Excel

```
GET /api/export/excel
```

Downloads a styled `.xlsx` file with 4 sheets:

| Sheet              | Contents                                      |
|-------------------|-----------------------------------------------|
| **Overview**       | Per-account summary (completion %, stats)     |
| **Tasks**          | All tasks with account, airdrop, category, gas|
| **Wallets**        | All wallets grouped by account                |
| **Quick Reference**| All airdrops with account assignment counts   |
| **Airdrop Tasks**  | All global airdrop task templates             |

**Response:** Binary `.xlsx` file download (`Content-Disposition: attachment`).

---

### Swagger UI

Interactive API documentation is available at:

```
GET /swagger/index.html
```

---

## curl Examples

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass123","name":"Test"}'

# Login (save token)
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"pass123"}' | jq -r '.token')

# Create an airdrop
curl -X POST http://localhost:8080/api/airdrops \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"zkSync","chain":"Ethereum","priority":"high"}'

# List all airdrops
curl http://localhost:8080/api/airdrops \
  -H "Authorization: Bearer $TOKEN"

# Create an account
curl -X POST http://localhost:8080/api/accounts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Account 1","color":"#EF4444"}'

# Assign airdrop to account
curl -X POST http://localhost:8080/api/accounts/1/airdrops \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"airdrop_id":1}'

# Get today's tasks
curl http://localhost:8080/api/accounts/1/tasks/today \
  -H "Authorization: Bearer $TOKEN"

# Get dashboard stats
curl http://localhost:8080/api/dashboard \
  -H "Authorization: Bearer $TOKEN"

# Export Excel
curl -o export.xlsx http://localhost:8080/api/export/excel \
  -H "Authorization: Bearer $TOKEN"
```

## License

MIT
