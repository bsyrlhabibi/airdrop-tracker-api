# GUIDE.md — Cara Menjalankan Airdrop Tracker API

Panduan lengkap untuk menjalankan project ini — lokal dan deploy ke Fly.io.

---

## Daftar Isi

1. [Persiapan](#1-persiapan)
2. [Install Go](#2-install-go)
3. [Clone Project](#3-clone-project)
4. [Setup Environment](#4-setup-environment)
5. [Install Dependencies](#5-install-dependencies)
6. [Jalankan Server (Lokal)](#6-jalankan-server-lokal)
7. [Buka Swagger UI](#7-buka-swagger-ui)
8. [Test API Manual](#8-test-api-manual)
9. [Deploy ke Fly.io](#9-deploy-ke-flyio)
10. [Troubleshooting](#10-troubleshooting)

---

## 1. Persiapan

Pastikan komputer kamu sudah punya:

| Tool | Cek di Terminal | Belum punya? |
|------|----------------|--------------|
| **Go** (v1.22+) | `go version` | [Install Go](https://go.dev/dl/) |
| **Git** | `git --version` | [Install Git](https://git-scm.com/) |
| **Terminal** | Buka Terminal / CMD | — |

> **Windows:** Pakai PowerShell atau CMD
> **Mac:** Pakai Terminal (Cmd + Space → "Terminal")
> **Linux:** Pakai Terminal biasa

---

## 2. Install Go

### Windows
1. Buka https://go.dev/dl/
2. Download file `go1.22.x.windows-amd64.msi`
3. Double-click → Install (ikuti wizard)
4. Buka PowerShell baru, ketik: `go version`
5. Harus muncul: `go version go1.22.x windows/amd64`

### Mac
```bash
# Pakai Homebrew (recommended)
brew install go

# Atau download dari https://go.dev/dl/
```

### Linux (Ubuntu/Debian)
```bash
# Download
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz

# Extract
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# Tambah ke PATH (tambahkan di ~/.bashrc atau ~/.zshrc)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Cek
go version
```

---

## 3. Clone Project

```bash
# Clone repository
git clone https://github.com/bsyrlhabibi/airdrop-tracker-api.git

# Masuk ke folder project
cd airdrop-tracker-api
```

> **Belum punya git?** Download ZIP dari GitHub → Extract → Buka folder di Terminal

---

## 4. Setup Environment

Project butuh file `.env` untuk konfigurasi. Cara buatnya:

```bash
# Copy template
cp .env.example .env
```

Isi file `.env` (bisa pakai text editor biasa):

```env
APP_PORT=8080
JWT_SECRET=ganti-dengan-random-string-kamu
DB_PATH=data/airdrop.db
```

**Penjelasan:**

| Variable | Fungsi | Contoh |
|----------|--------|--------|
| `APP_PORT` | Port server berjalan | `8080` (default, boleh ganti) |
| `JWT_SECRET` | Kunci enkripsi token | Bebas, makin random makin aman |
| `DB_PATH` | Lokasi file database | `data/airdrop.db` (default) |

> **Tips:** `JWT_SECRET` bisa diisi apa saja, contoh: `rahasia-bita-2024-xyz`

---

## 5. Install Dependencies

Setelah clone dan setup `.env`, **wajib jalankan ini dulu** untuk download semua library yang dibutuhkan:

```bash
go mod tidy
```

Perintah ini akan:
- Download semua dependency (Gin, GORM, JWT, dll)
- Bersihkan dependency yang tidak dipakai
- Buat/update file `go.sum` (checksum)

> **Tanpa ini, project tidak bisa jalan!**

Kalau berhasil, outputnya kurang lebih:

```
go: finding module for package gorm.io/gorm
go: finding module for package github.com/gin-gonic/gin
go: found github.com/gin-gonic/gin in github.com/gin-gonic/gin v1.12.0
...
```

---

## 6. Jalankan Server (Lokal)

### Cara 1: Pakai Makefile (Recommended)

```bash
make run
```

### Cara 2: Tanpa Makefile

```bash
go run cmd/server/main.go
```

### Kalau Berhasil

Terminal akan menampilkan:

```
Database connected
Migration done
[GIN-debug] Listening and serving HTTP on :8080
Server running on :8080
Swagger UI: http://localhost:8080/swagger/index.html
```

> **Jangan tutup Terminal ini!** Server harus tetap jalan.

---

## 7. Buka Swagger UI

1. Buka browser (Chrome, Firefox, Safari, dll)
2. Ketik di address bar:

```
http://localhost:8080/swagger/index.html
```

3. Kamu akan melihat halaman Swagger UI dengan semua endpoint API

### Cara Pakai Swagger

#### Step 1: Register Akun
1. Klik `/api/auth/register` → **Try it out**
2. Isi body:
```json
{
  "email": "test@email.com",
  "password": "password123",
  "name": "Bita"
}
```
3. Klik **Execute**
4. Harus muncul response `201` dengan pesan "Registered"

#### Step 2: Login
1. Klik `/api/auth/login` → **Try it out**
2. Isi body:
```json
{
  "email": "test@email.com",
  "password": "password123"
}
```
3. Klik **Execute**
4. Copy `token` dari response

#### Step 3: Authorize
1. Klik tombol **Authorize** 🔓 di pojok kanan atas
2. Paste token dengan format: `Bearer <token-kamu>`
   Contoh: `Bearer eyJhbGciOiJIUzI1NiIs...`
3. Klik **Authorize** → **Close**

#### Step 4: Pakai Endpoint Lain
Sekarang semua endpoint yang terkunci 🔒 sudah bisa diakses!

Contoh buat airdrop pertama:
1. Klik `/api/airdrops` → **Try it out**
2. Isi body:
```json
{
  "name": "zkSync",
  "chain": "Ethereum",
  "category": "rumored",
  "priority": "high",
  "url": "https://zksync.io",
  "notes": "Bridge minimal $100 setiap minggu"
}
```
3. Klik **Execute**

---

## 8. Test API Manual (Terminal)

Selain Swagger, kamu juga bisa test pakai `curl`:

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@email.com","password":"password123","name":"Bita"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@email.com","password":"password123"}'
```

### Get Airdrops (pakai token dari login)
```bash
curl http://localhost:8080/api/airdrops \
  -H "Authorization: Bearer <TOKEN_KAMU>"
```

### Create Airdrop
```bash
curl -X POST http://localhost:8080/api/airdrops \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN_KAMU>" \
  -d '{"name":"zkSync","chain":"Ethereum","category":"rumored","priority":"high"}'
```

### Get Dashboard
```bash
curl http://localhost:8080/api/dashboard \
  -H "Authorization: Bearer <TOKEN_KAMU>"
```

---

## 9. Deploy ke Fly.io

### Step 1: Install Flyctl

```bash
# Mac / Linux
curl -L https://fly.io/install.sh | sh

# Windows (PowerShell)
iwr https://fly.io/install.ps1 -useb | iex
```

### Step 2: Login

```bash
flyctl auth login
```

Browser terbuka → Login / daftar akun Fly.io

### Step 3: Launch App

```bash
cd airdrop-tracker-api
flyctl launch --no-deploy
```

Jawab pertanyaan:
```
? App name → airdrop-tracker-api (atau biarkan random)
? Select region → Singapore (sin)
? Would you like to set up a Postgresql database? → No
? Would you like to deploy now? → No
```

### Step 4: Create Volume (Persistent Storage)

```bash
flyctl volume create airdrop_data --size 1 --region sin
```

> Ini supaya database SQLite tidak hilang kalau VM restart

### Step 5: Set JWT Secret

```bash
flyctl secrets set JWT_SECRET=ganti-dengan-random-string-kamu
```

> ⚠️ **WAJIB!** Tanpa ini, auth tidak akan work

### Step 6: Deploy

```bash
flyctl deploy
```

Tunggu sampai selesai. Output akhir:
```
--> v0 deployed successfully
```

### Step 7: Cek

```bash
# Status
flyctl status

# Logs
flyctl logs

# Test API
curl https://airdrop-tracker-api.fly.dev/swagger/index.html
```

### Step 8: Dapat URL

```
https://airdrop-tracker-api.fly.dev
```

Buka di browser → Swagger UI muncul → API ready! 🎉

---

### Update & Redeploy

Kalau ada perubahan code:

```bash
# Commit & push ke GitHub
git add .
git commit -m "update: description"
git push

# Deploy ke Fly.io
flyctl deploy
```

Atau pakai Makefile:
```bash
make deploy
```

---

## 10. Troubleshooting

### "command not found: go"
→ Go belum ter-install atau belum masuk PATH. Ulangi [Step 2](#2-install-go).

### "port already in use"
→ Port 8080 sudah dipakai app lain. Ganti port di `.env`:
```env
APP_PORT=3000
```

### "no such file or directory: data/airdrop.db"
→ Folder `data/` belum ada. Buat manual:
```bash
mkdir -p data
```

### Swagger UI blank / 404
→ Pastikan server jalan (Step 6). Coba regenerate docs:
```bash
swag init -g cmd/server/main.go
make run
```

### "connection refused" saat test API
→ Server belum jalan. Pastikan Terminal tempat `make run` masih terbuka dan menampilkan `Listening on :8080`.

### "missing go.sum entry" atau error dependency
→ Jalankan ulang:
```bash
go mod tidy
```

### Lupa password / token expired
→ Register ulang dengan email baru, atau hapus database dan mulai fresh:
```bash
rm data/airdrop.db
make run
```

### Fly.io: "app not found"
→ Pastikan sudah jalankan `flyctl launch` dulu

### Fly.io: data hilang setelah restart
→ Pastikan volume sudah dibuat:
```bash
flyctl volume create airdrop_data --size 1 --region sin
```

---

## Perintah Berguna

### Lokal
```bash
go mod tidy                     # Install/update dependencies
make run                        # Jalankan server
make build                      # Build binary (hasilnya di bin/server)
make clean                      # Hapus database & binary
make swag                       # Regenerate Swagger docs
make test                       # Jalankan test (jika ada)
```

### Fly.io
```bash
flyctl status                   # Cek status app
flyctl logs                     # Lihat logs
flyctl deploy                   # Deploy update
flyctl ssh console              # SSH ke VM
flyctl secrets list             # Lihat env variables
```

---

## Ringkasan Cepat (Copy-Paste)

### Lokal
```bash
git clone https://github.com/bsyrlhabibi/airdrop-tracker-api.git
cd airdrop-tracker-api
cp .env.example .env
go mod tidy
make run
```

Lalu buka browser: `http://localhost:8080/swagger/index.html`

### Deploy
```bash
flyctl launch --no-deploy
flyctl volume create airdrop_data --size 1 --region sin
flyctl secrets set JWT_SECRET=rahasia-kamu
flyctl deploy
```

Lalu buka browser: `https://airdrop-tracker-api.fly.dev/swagger/index.html`

---

Selamat mencoba! 🚀
