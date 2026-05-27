# 🪂 Airdrop Tracker API

Backend REST API untuk mengelola airdrop crypto — tracking wallet, task harian, dan progress farming multi-account.

Built with **Go + Gin + GORM + SQLite**.

---

## ✨ Features

- **Auth** — Register & login dengan JWT
- **Accounts** — Multi-account (sybil) management dengan color coding
- **Airdrops** — Global airdrop catalog (chain, category, priority, status)
- **Airdrop Tasks** — Template task per airdrop (start/end date, auto-expand harian)
- **Account Tasks** — Daily tracking per account (pending → ongoing → finish/missed)
- **Categories** — Kustomisasi kategori task
- **Wallets** — Multi-wallet per account & chain
- **Dashboard** — Stats summary + per-account comparison
- **Excel Export** — 5-sheet styled export (Overview, Tasks, Wallets, Quick Reference, Airdrop Tasks)
- **Swagger UI** — Interactive API documentation

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────┐
│                  API Layer (Gin)                 │
│  handler/ → router/ → middleware/ (JWT auth)     │
├─────────────────────────────────────────────────┤
│              Data Layer (GORM + SQLite)          │
│  model/ → repository/ → database/               │
├─────────────────────────────────────────────────┤
│                 SQLite File                      │
│               data/airdrop.db                    │
└─────────────────────────────────────────────────┘
```

**Flow:**
```
Airdrop Catalog (global) ──→ Account assigns airdrop
                                   │
                                   ▼
                          Auto-sync tasks from template
                                   │
                                   ▼
                          Daily tracking per account
                          (pending/ongoing/finish/missed)
```

---

## 📁 Project Structure

```
airdrop-tracker-api/
├── cmd/server/main.go                  # Entry point
├── internal/
│   ├── config/config.go                # Env config (APP_PORT, JWT_SECRET, DB_PATH)
│   ├── database/database.go            # SQLite connection + auto-migration
│   ├── model/                          # Data models
│   │   ├── user.go                     #   User (auth)
│   │   ├── account.go                  #   Account (sybil identity)
│   │   ├── account_airdrop.go          #   Junction: account ↔ airdrop
│   │   ├── airdrop.go                  #   Airdrop (global catalog)
│   │   ├── airdrop_task.go             #   AirdropTask (template)
│   │   ├── task.go                     #   Task (daily tracking)
│   │   ├── category.go                 #   Category (task type)
│   │   └── wallet.go                   #   Wallet (address per chain)
│   ├── repository/                     # Database queries (GORM)
│   ├── handler/                        # HTTP handlers
│   │   ├── auth.go                     #   Register + Login
│   │   ├── account.go                  #   Account CRUD + assign/clone
│   │   ├── airdrop.go                  #   Airdrop CRUD
│   │   ├── airdrop_task.go             #   AirdropTask CRUD
│   │   ├── task.go                     #   Task CRUD + today/by-date
│   │   ├── category.go                 #   Category CRUD
│   │   ├── wallet.go                   #   Wallet CRUD
│   │   ├── dashboard.go                #   Stats summary
│   │   └── export.go                   #   Excel export (5 sheets)
│   ├── middleware/auth.go              # JWT middleware
│   └── router/router.go               # Route definitions
├── docs/                               # Swagger generated docs
├── data/                               # SQLite database (gitignored)
├── .env.example                        # Environment template
├── .gitignore
├── Dockerfile
├── Makefile
├── GUIDE.md                            # Setup guide (beginner-friendly)
└── README.md                           # ← Kamu di sini
```

---

## 🛠️ Tech Stack

| Component      | Library                                                    |
|----------------|------------------------------------------------------------|
| HTTP Framework | [Gin](https://github.com/gin-gonic/gin)                   |
| ORM            | [GORM](https://gorm.io)                                   |
| Database       | [SQLite](https://www.sqlite.org/) (via gorm driver)       |
| Auth           | [golang-jwt](https://github.com/golang-jwt/jwt) + bcrypt |
| Config         | [godotenv](https://github.com/joho/godotenv)              |
| Excel Export   | [excelize](https://github.com/xuri/excelize)              |
| Docs           | [swaggo/swag](https://github.com/swaggo/swag) + Swagger UI|

---

## 🚀 Quick Start

```bash
# 1. Clone
git clone https://github.com/bsyrlhabibi/airdrop-tracker-api.git
cd airdrop-tracker-api

# 2. Setup env
cp .env.example .env

# 3. Install dependencies
go mod tidy

# 4. Run
make run

# 5. Buka browser
# → http://localhost:8080/swagger/index.html
```

> 📖 **Belum pernah pakai Go?** Baca [GUIDE.md](./GUIDE.md) — panduan lengkap dari nol!

---

## 📡 API Endpoints

### Public
| Method | Path                 | Description        |
|--------|----------------------|--------------------|
| POST   | `/api/auth/register` | Register new user  |
| POST   | `/api/auth/login`    | Login, get JWT     |

### Accounts
| Method | Path                              | Description              |
|--------|-----------------------------------|--------------------------|
| GET    | `/api/accounts`                   | List all accounts        |
| POST   | `/api/accounts`                   | Create account           |
| GET    | `/api/accounts/:id`               | Get account detail       |
| PUT    | `/api/accounts/:id`               | Update account           |
| DELETE | `/api/accounts/:id`               | Delete account           |
| POST   | `/api/accounts/:id/airdrops`      | Assign airdrop to account|
| GET    | `/api/accounts/:id/airdrops`      | List account's airdrops  |
| DELETE | `/api/accounts/:id/airdrops/:aid` | Remove airdrop           |
| POST   | `/api/accounts/:id/clone`         | Clone account + wallets  |
| GET    | `/api/accounts/:id/tasks/today`   | Today's tasks            |
| GET    | `/api/accounts/:id/tasks/by-date` | Tasks by date            |

### Airdrops (Global Catalog)
| Method | Path                  | Description     |
|--------|-----------------------|-----------------|
| GET    | `/api/airdrops`       | List airdrops   |
| POST   | `/api/airdrops`       | Create airdrop  |
| GET    | `/api/airdrops/:id`   | Get airdrop     |
| PUT    | `/api/airdrops/:id`   | Update airdrop  |
| DELETE | `/api/airdrops/:id`   | Delete airdrop  |

### Airdrop Tasks (Templates)
| Method | Path                          | Description          |
|--------|-------------------------------|----------------------|
| GET    | `/api/airdrops/:id/tasks`     | List tasks           |
| POST   | `/api/airdrops/:id/tasks`     | Create task          |
| POST   | `/api/airdrops/:id/tasks/bulk`| Bulk create tasks    |
| PUT    | `/api/airdrops/:id/tasks/reorder` | Reorder tasks  |
| PUT    | `/api/airdrop-tasks/:id`      | Update task          |
| DELETE | `/api/airdrop-tasks/:id`      | Delete task          |

### Account Tasks (Daily Tracking)
| Method | Path                            | Description        |
|--------|---------------------------------|--------------------|
| GET    | `/api/account-airdrops/:id/tasks` | List tasks       |
| POST   | `/api/account-airdrops/:id/tasks` | Create task      |
| PUT    | `/api/tasks/:id`                | Update task        |
| DELETE | `/api/tasks/:id`                | Delete task        |

### Other
| Method | Path                       | Description          |
|--------|----------------------------|----------------------|
| GET    | `/api/categories`          | List categories      |
| POST   | `/api/categories`          | Create category      |
| PUT    | `/api/categories/:id`      | Update category      |
| DELETE | `/api/categories/:id`      | Delete category      |
| GET    | `/api/wallets`             | List wallets         |
| POST   | `/api/wallets`             | Add wallet           |
| DELETE | `/api/wallets/:id`         | Delete wallet        |
| GET    | `/api/dashboard`           | Stats summary        |
| GET    | `/api/dashboard/comparison`| Account comparison   |
| GET    | `/api/export/excel`        | Export to Excel      |

---

## ⚙️ Environment Variables

| Variable     | Required | Default              | Description                    |
|--------------|----------|----------------------|--------------------------------|
| `APP_PORT`   | No       | `8080`               | Server port                    |
| `JWT_SECRET` | **Yes**  | `default-secret`     | JWT signing key (change this!) |
| `DB_PATH`    | No       | `data/airdrop.db`    | SQLite database path           |

---

## 📦 Makefile Commands

```bash
make run        # Jalankan server
make build      # Build binary → bin/server
make swag       # Regenerate Swagger docs
make test       # Run tests
make clean      # Hapus database & binary
make deploy     # Deploy ke Fly.io
```

---

## 📄 License

MIT

---

> 📖 **Butuh panduan lengkap?** Lihat [GUIDE.md](./GUIDE.md)
