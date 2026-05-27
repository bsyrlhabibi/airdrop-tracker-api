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
- **Dashboard** — Stats summary per account
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

### 🔓 Public

| Method | Path                 | Description        |
|--------|----------------------|--------------------|
| POST   | `/api/auth/register` | Register new user  |
| POST   | `/api/auth/login`    | Login, get JWT     |

### 🔒 Protected (Bearer Token)

#### Accounts

| Method | Path                                       | Description               |
|--------|--------------------------------------------|---------------------------|
| GET    | `/api/accounts`                            | List all accounts         |
| POST   | `/api/accounts`                            | Create account            |
| GET    | `/api/accounts/:id`                        | Get account detail        |
| PUT    | `/api/accounts/:id`                        | Update account            |
| DELETE | `/api/accounts/:id`                        | Delete account            |
| POST   | `/api/accounts/:id/airdrops`               | Assign airdrop to account |
| DELETE | `/api/accounts/:id/airdrops/:airdrop_id`   | Remove airdrop            |

#### Airdrops (Global Catalog)

| Method | Path                  | Description        |
|--------|-----------------------|--------------------|
| GET    | `/api/airdrops`       | List airdrops      |
| POST   | `/api/airdrops`       | Create airdrop     |
| GET    | `/api/airdrops/:id`   | Get airdrop detail |
| PUT    | `/api/airdrops/:id`   | Update airdrop     |
| DELETE | `/api/airdrops/:id`   | Delete airdrop     |

#### Airdrop Tasks (Templates)

| Method | Path                      | Description                     |
|--------|---------------------------|---------------------------------|
| GET    | `/api/airdrops/:id/tasks` | List template tasks             |
| POST   | `/api/airdrops/:id/tasks` | Create template task            |
| PUT    | `/api/airdrop-tasks/:id`  | Update template task            |
| DELETE | `/api/airdrop-tasks/:id`  | Delete template + cascade daily |

#### Account Tasks (Daily Tracking)

| Method | Path                              | Description              |
|--------|-----------------------------------|--------------------------|
| POST   | `/api/account-airdrops/:id/tasks` | Create daily task        |
| PUT    | `/api/tasks/:id`                  | Update task status       |
| DELETE | `/api/tasks/:id`                  | Delete task              |
| GET    | `/api/accounts/:id/tasks/today`   | Get today's tasks        |
| GET    | `/api/accounts/:id/tasks/by-date` | Get tasks by date        |

#### Categories

| Method | Path                  | Description      |
|--------|-----------------------|------------------|
| GET    | `/api/categories`     | List categories  |
| POST   | `/api/categories`     | Create category  |
| PUT    | `/api/categories/:id` | Update category  |
| DELETE | `/api/categories/:id` | Delete category  |

#### Wallets

| Method | Path               | Description    |
|--------|--------------------|----------------|
| GET    | `/api/wallets`     | List wallets   |
| POST   | `/api/wallets`     | Add wallet     |
| DELETE | `/api/wallets/:id` | Delete wallet  |

#### Dashboard & Export

| Method | Path                | Description                |
|--------|---------------------|----------------------------|
| GET    | `/api/dashboard`    | Stats summary              |
| GET    | `/api/export/excel` | Export to Excel (5 sheets) |

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
make run            # Jalankan server
make build          # Build binary → bin/server
make swag           # Regenerate Swagger docs
make swag-install   # Install swag CLI
make test           # Run tests
make clean          # Hapus database & binary
make deploy         # Deploy ke Fly.io
```

---

## 📄 License

MIT

---

> 📖 **Butuh panduan lengkap?** Lihat [GUIDE.md](./GUIDE.md)
