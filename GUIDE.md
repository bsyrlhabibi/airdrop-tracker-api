# Airdrop Tracker API — Beginner's Guide

Welcome! This guide walks you through setting up and using the Airdrop Tracker API from scratch. No prior experience with Go is needed — just follow the steps.

---

## What Is This?

This is a **backend API** (a server that responds to HTTP requests) for tracking cryptocurrency airdrops. It lets you:

- **Create accounts** (sybil wallets) to participate in airdrops
- **Track airdrops** you want to farm (e.g., zkSync, Arbitrum, Base)
- **Manage tasks** for each airdrop per account (e.g., "Bridge 0.1 ETH")
- **Record gas spent** and transaction hashes
- **See a dashboard** with progress stats
- **Export everything to Excel** with styled sheets

The API does **not** have a UI — it returns JSON data. You connect a mobile or web frontend to it (or just test with curl/Postman).

---

## Prerequisites

Before you start, make sure you have:

1. **Go** (version 1.25 or later) — download from https://go.dev/dl/
2. **A C compiler** — needed for the SQLite driver
   - **macOS:** `xcode-select --install`
   - **Ubuntu/Debian:** `sudo apt install gcc`
   - **Windows:** Comes with Go installer or use TDM-GCC
3. **Git** — to clone the repository

Verify Go is installed:

```bash
go version
# Should print: go version go1.25.x ...
```

---

## Step 1: Get the Code

```bash
git clone <repo-url>
cd airdrop-tracker-api
```

---

## Step 2: Configure Environment

The app reads settings from a `.env` file. There's a template ready for you:

```bash
cp .env.example .env
```

Open `.env` in any text editor. You'll see:

```env
APP_PORT=8080
JWT_SECRET=change-me-to-random-string
DB_PATH=data/airdrop.db
```

**What each setting does:**

- `APP_PORT` — The port the server listens on. Default `8080` is fine.
- `JWT_SECRET` — A secret password used to sign login tokens. **Change this** to something random. You can generate one with: `openssl rand -hex 32`
- `DB_PATH` — Where the SQLite database file is saved. Default is fine.

---

## Step 3: Install Dependencies

```bash
go mod download
```

This downloads all the Go libraries the project needs (Gin, GORM, JWT, etc.).

---

## Step 4: Create the Data Folder

```bash
mkdir -p data
```

This is where the SQLite database file (`airdrop.db`) will be created automatically on first run.

---

## Step 5: Start the Server

```bash
go run ./cmd/server/main.go
```

You should see output like:

```
Database connected
Migration done
Server running on :8080
Swagger UI: http://localhost:8080/swagger/index.html
```

The server is now running! Leave this terminal open.

---

## Step 6: Test It

Open a **new terminal** and try these commands:

### 6.1 Register a user

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"me@test.com","password":"mypassword","name":"Me"}'
```

Expected response:

```json
{
  "message": "Registered",
  "user": { "id": 1, "email": "me@test.com", "name": "Me" }
}
```

### 6.2 Login and get a token

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"me@test.com","password":"mypassword"}'
```

Expected response:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": { "id": 1, "email": "me@test.com", "name": "Me" }
}
```

**Copy the token value.** You'll use it for every other request.

### 6.3 Save the token for convenience

```bash
export TOKEN="eyJhbGciOiJIUzI1NiIs..."
```

Now you can use `$TOKEN` in all subsequent curl commands.

### 6.4 Create your first airdrop

```bash
curl -X POST http://localhost:8080/api/airdrops \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"zkSync","chain":"Ethereum","priority":"high","notes":"Bridge weekly"}'
```

### 6.5 Create a sybil account

```bash
curl -X POST http://localhost:8080/api/accounts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Account 1","color":"#EF4444","notes":"Main farming account"}'
```

### 6.6 Assign the airdrop to the account

```bash
curl -X POST http://localhost:8080/api/accounts/1/airdrops \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"airdrop_id":1}'
```

### 6.7 Check the dashboard

```bash
curl http://localhost:8080/api/dashboard \
  -H "Authorization: Bearer $TOKEN"
```

### 6.8 Export to Excel

```bash
curl -o my-airdrops.xlsx http://localhost:8080/api/export/excel \
  -H "Authorization: Bearer $TOKEN"
```

Open `my-airdrops.xlsx` in Excel or Google Sheets — it has styled sheets with your data.

---

## Using Postman (Alternative to curl)

If you prefer a GUI:

1. Download [Postman](https://www.postman.com/downloads/)
2. Create a new request
3. Set the method (GET, POST, PUT, DELETE)
4. Enter the URL (e.g., `http://localhost:8080/api/airdrops`)
5. Go to the **Headers** tab and add:
   - Key: `Authorization`
   - Value: `Bearer <your-token>`
6. For POST/PUT, go to the **Body** tab → select **raw** → **JSON** and paste the request body
7. Click **Send**

---

## Folder Structure Explained

Here's what each folder does:

| Folder/File | What It Does |
|---|---|
| `cmd/server/main.go` | **Start here.** The entrypoint that boots everything up. |
| `internal/config/` | Reads `.env` file and creates a config struct. |
| `internal/database/` | Connects to SQLite and runs automatic migrations (creates tables). |
| `internal/model/` | Defines your data structures (User, Airdrop, Account, Task, etc.). Think of these as database table schemas. |
| `internal/repository/` | Database query functions. Each file handles DB operations for one model. |
| `internal/handler/` | HTTP request handlers. Each file handles one resource's API endpoints. |
| `internal/middleware/` | The JWT authentication middleware that checks tokens on protected routes. |
| `internal/router/` | Maps URLs to handlers and configures CORS + Swagger. |
| `docs/` | Auto-generated Swagger documentation files. |

---

## Data Model Overview

Here's how the data connects:

```
User
├── Accounts (sybil accounts)
│   ├── Wallets (blockchain addresses)
│   └── AccountAirdrops (assigned airdrops)
│       └── Tasks (per-account task progress)
├── Airdrops (global airdrop catalog)
│   └── AirdropTasks (global task templates)
└── Categories (task labels like "Bridge", "Swap", "Daily")
```

**Key relationships:**

- An **Airdrop** is a global entry (e.g., "zkSync"). It belongs to one user.
- An **Account** is a sybil identity (e.g., "Account 1"). It also belongs to one user.
- **AccountAirdrop** links an Account to an Airdrop (many-to-many).
- **AirdropTask** is a global task template (e.g., "Bridge 0.1 ETH") for an Airdrop.
- **Task** is a per-account task created when you assign an airdrop to an account.
- **Category** is a label you can attach to tasks (e.g., "Bridge", "Staking").
- **Wallet** is a blockchain address tied to an Account.

---

## Status Values Reference

### Airdrop Status
- `active` — Currently farming
- `upcoming` — Not started yet
- `end` — Airdrop has ended
- `missed` — Missed the window

### Task Status
- `pending` — Not started
- `ongoing` — In progress
- `finish` — Completed
- `missed` — Missed deadline

### Task Frequency
- `once` — One-time task
- `daily` — Repeat every day
- `weekly` — Repeat every week
- `monthly` — Repeat every month

---

## Common Troubleshooting

### "gcc: command not found" when running the server

The SQLite driver needs a C compiler. Install it:

```bash
# Ubuntu/Debian
sudo apt install gcc

# macOS
xcode-select --install
```

### "port already in use" error

Another process is using port 8080. Either:

- Kill it: `lsof -i :8080` then `kill <PID>`
- Change the port in `.env`: `APP_PORT=3000`

### "Authorization header required" error

You forgot to include the JWT token. Add this header:

```
-H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### "Invalid token" error

Your token has expired (tokens last 72 hours) or is malformed. Log in again to get a new token.

### "Email already exists" when registering

That email is already taken. Use a different email or log in with the existing one.

### Database is corrupted or behaving strangely

Delete the database file and restart — the server will recreate it:

```bash
rm data/airdrop.db
go run ./cmd/server/main.go
```

**Warning:** This deletes all your data. Back up `data/airdrop.db` first if needed.

### "Account has wallets or airdrops" when deleting

The account has linked data. Use `?force=true` to delete everything:

```bash
curl -X DELETE "http://localhost:8080/api/accounts/1?force=true" \
  -H "Authorization: Bearer $TOKEN"
```

---

## Quick Reference: All API Endpoints

| Method | URL | Auth | Description |
|--------|-----|------|-------------|
| POST | `/api/auth/register` | No | Register new user |
| POST | `/api/auth/login` | No | Login, get JWT token |
| GET | `/api/airdrops` | Yes | List all airdrops |
| POST | `/api/airdrops` | Yes | Create airdrop |
| GET | `/api/airdrops/:id` | Yes | Get airdrop details |
| PUT | `/api/airdrops/:id` | Yes | Update airdrop |
| DELETE | `/api/airdrops/:id` | Yes | Delete airdrop |
| GET | `/api/airdrops/:id/tasks` | Yes | List airdrop task templates |
| POST | `/api/airdrops/:id/tasks` | Yes | Create airdrop task template |
| PUT | `/api/airdrop-tasks/:id` | Yes | Update airdrop task template |
| DELETE | `/api/airdrop-tasks/:id` | Yes | Delete airdrop task template |
| GET | `/api/accounts` | Yes | List all accounts |
| POST | `/api/accounts` | Yes | Create account |
| GET | `/api/accounts/:id` | Yes | Get account details |
| PUT | `/api/accounts/:id` | Yes | Update account |
| DELETE | `/api/accounts/:id` | Yes | Delete account (?force=true) |
| POST | `/api/accounts/:id/airdrops` | Yes | Assign airdrop to account |
| DELETE | `/api/accounts/:id/airdrops/:aid` | Yes | Remove airdrop from account |
| POST | `/api/account-airdrops/:id/tasks` | Yes | Create per-account task |
| PUT | `/api/tasks/:id` | Yes | Update task |
| DELETE | `/api/tasks/:id` | Yes | Delete task |
| GET | `/api/accounts/:id/tasks/today` | Yes | Get today's tasks |
| GET | `/api/accounts/:id/tasks/by-date` | Yes | Get tasks for a date |
| GET | `/api/categories` | Yes | List categories |
| POST | `/api/categories` | Yes | Create category |
| PUT | `/api/categories/:id` | Yes | Update category |
| DELETE | `/api/categories/:id` | Yes | Delete category |
| GET | `/api/wallets` | Yes | List wallets (?account_id=N) |
| POST | `/api/wallets` | Yes | Create wallet |
| DELETE | `/api/wallets/:id` | Yes | Delete wallet |
| GET | `/api/dashboard` | Yes | Dashboard statistics |
| GET | `/api/export/excel` | Yes | Export Excel file |
| GET | `/swagger/index.html` | No | Swagger UI docs |

---

## Building for Production

To build a standalone binary:

```bash
go build -o airdrop-server ./cmd/server
```

Then deploy and run:

```bash
JWT_SECRET=your-production-secret DB_PATH=/var/data/airdrop.db APP_PORT=8080 ./airdrop-server
```

You can also set these as system environment variables instead of prefixing the command.
