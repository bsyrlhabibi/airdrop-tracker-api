# Airdrop Tracker API

Personal airdrop task management backend — built with **Go**, **Gin**, **GORM**, and **SQLite**.

Manage airdrops, track daily tasks, monitor wallets, and view farming progress — all from a clean REST API.

---

## Features

- **Auth** — Register & login with JWT authentication
- **Airdrops** — CRUD airdrop tracking with category, priority, status, deadline
- **Tasks** — Task checklist per airdrop (one-time, daily, weekly, monthly)
- **Wallets** — Multi-wallet management per chain
- **Dashboard** — Summary stats (total airdrops, tasks completed, pending)
- **Swagger UI** — Interactive API documentation built-in

---

## Tech Stack

| Component | Library |
|-----------|---------|
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io) |
| Database | [SQLite](https://www.sqlite.org/) (via `gorm.io/driver/sqlite`) |
| Auth | [golang-jwt](https://github.com/golang-jwt/jwt) + bcrypt |
| Config | [godotenv](https://github.com/joho/godotenv) |
| Docs | [swaggo/swag](https://github.com/swaggo/swag) + Swagger UI |

---

## Project Structure

```
airdrop-tracker/
├── cmd/server/main.go              # Entry point
├── internal/
│   ├── config/config.go            # Environment config
│   ├── database/database.go        # SQLite connection + migration
│   ├── model/                      # Data models (User, Airdrop, Task, Wallet)
│   ├── repository/                 # Data access layer
│   ├── service/                    # Business logic (extensible)
│   ├── handler/                    # HTTP handlers
│   ├── middleware/auth.go          # JWT middleware
│   └── router/router.go           # Route definitions
├── docs/                           # Swagger generated docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── data/                           # SQLite database file
├── .env.example                    # Environment template
├── .gitignore
├── Makefile
├── GUIDE.md                        # Step-by-step setup guide
└── README.md
```

---

## API Endpoints

### Public
| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login, get JWT token |

### Protected (Bearer Token)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/airdrops` | List all airdrops |
| POST | `/api/airdrops` | Create airdrop |
| GET | `/api/airdrops/:id` | Get airdrop detail |
| PUT | `/api/airdrops/:id` | Update airdrop |
| DELETE | `/api/airdrops/:id` | Delete airdrop |
| GET | `/api/airdrops/:id/tasks` | List tasks |
| POST | `/api/airdrops/:id/tasks` | Create task |
| PUT | `/api/tasks/:id/complete` | Mark task done |
| PUT | `/api/tasks/:id/reset` | Reset task |
| DELETE | `/api/tasks/:id` | Delete task |
| GET | `/api/wallets` | List wallets |
| POST | `/api/wallets` | Add wallet |
| DELETE | `/api/wallets/:id` | Delete wallet |
| GET | `/api/dashboard` | Stats summary |

---

## Quick Start

```bash
# 1. Clone
git clone https://github.com/bsyrlhabibi/airdrop-tracker-api.git
cd airdrop-tracker-api

# 2. Setup env
cp .env.example .env

# 3. Run
make run

# 4. Open Swagger
# Browser → http://localhost:8080/swagger/index.html
```

> For detailed setup instructions, see [GUIDE.md](./GUIDE.md)

---

## License

MIT
