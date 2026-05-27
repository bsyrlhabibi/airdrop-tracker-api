# 📖 GUIDE — Cara Menjalankan Airdrop Tracker API

Panduan **lengkap dari nol** untuk orang yang belum pernah pakai Go.
Ikuti langkah-langkah berurutan, jangan skip!

---

## 📋 Daftar Isi

1. [Persiapan](#1-persiapan)
2. [Install Go](#2-install-go)
3. [Clone Project](#3-clone-project)
4. [Setup Environment](#4-setup-environment)
5. [Install Dependencies](#5-install-dependencies)
6. [Jalankan Server](#6-jalankan-server)
7. [Buka Swagger UI](#7-buka-swagger-ui)
8. [Test API Pertama Kali](#8-test-api-pertama-kali)
9. [Generate Swagger Docs](#9-generate-swagger-docs)
10. [Docker (Opsional)](#10-docker-opsional)
11. [Troubleshooting](#11-troubleshooting)
12. [Cheat Sheet](#12-cheat-sheet)

---

## 1. Persiapan

Pastikan komputer kamu punya **2 tools** ini:

### Cek di Terminal / CMD / PowerShell:

```bash
go version
```

```bash
git --version
```

| Tool     | Output yang benar                      | Belum punya?                          |
|----------|----------------------------------------|---------------------------------------|
| **Go**   | `go version go1.22.x ...`              | [Install Go](#2-install-go)           |
| **Git**  | `git version 2.x.x`                   | [Install Git](https://git-scm.com/)   |

> 💡 **Terminal?**
> - **Windows:** Tekan `Win + R` → ketik `cmd` → Enter. Atau pakai PowerShell.
> - **Mac:** Tekan `Cmd + Space` → ketik `Terminal` → Enter.
> - **Linux:** `Ctrl + Alt + T`

---

## 2. Install Go

Pilih sesuai OS kamu:

### 🪟 Windows

1. Buka **https://go.dev/dl/**
2. Download file: `go1.22.x.windows-amd64.msi`
3. **Double-click** file `.msi` → Next → Next → Install
4. **Tutup** semua terminal, **buka baru**
5. Ketik:
   ```bash
   go version
   ```
6. Harus muncul: `go version go1.22.x windows/amd64`

> ⚠️ **Kalau belum muncul**, restart komputer dulu, lalu coba lagi.

### 🍎 Mac

```bash
# Cara paling gampang — pakai Homebrew
brew install go

# Belum punya Homebrew? Install dulu:
# /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
# Lalu: brew install go
```

Atau download manual dari https://go.dev/dl/ → pilih `go1.22.x.darwin-arm64.pkg`

### 🐧 Linux (Ubuntu / Debian)

```bash
# Download Go
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz

# Extract ke /usr/local
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# Tambahkan ke PATH (copy-paste sekaligus)
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Cek
go version
```

> Harus muncul: `go version go1.22.5 linux/amd64`

---

## 3. Clone Project

```bash
git clone https://github.com/bsyrlhabibi/airdrop-tracker-api.git
cd airdrop-tracker-api
```

**Cek isi folder:**
```bash
ls
```

Harus terlihat: `cmd/`, `internal/`, `Makefile`, `.env.example`, dll.

> 💡 **Gak punya git?** Buka https://github.com/bsyrlhabibi/airdrop-tracker-api → klik tombol hijau **Code** → **Download ZIP** → Extract → buka folder di terminal.

---

## 4. Setup Environment

Project butuh file `.env` untuk konfigurasi database dan port.

### Buat file .env:

```bash
cp .env.example .env
```

### Edit file .env:

Buka file `.env` dengan text editor apapun (Notepad, VS Code, nano, dll):

```env
APP_PORT=8080
JWT_SECRET=rahasia-kuat-123-abc
DB_PATH=data/airdrop.db
```

**Penjelasan:**

| Variable     | Apa ini?                              | Contoh                    |
|--------------|-----------------------------------------|---------------------------|
| `APP_PORT`   | Port server berjalan                    | `8080` (default, aman)    |
| `JWT_SECRET` | Kunci enkripsi token login              | Bebas, makin random makin bagus |
| `DB_PATH`    | Lokasi file database SQLite             | `data/airdrop.db` (default) |

> 💡 **Tips JWT_SECRET:** Bisa isi apa saja, contoh: `my-super-secret-key-2024`
> Atau generate random: `openssl rand -hex 32`

---

## 5. Install Dependencies

**Wajib jalankan ini** sebelum pertama kali run:

```bash
go mod tidy
```

Perintah ini akan:
- ✅ Download semua library yang dibutuhkan (Gin, GORM, JWT, Excelize, dll)
- ✅ Bersihkan dependency yang tidak dipakai
- ✅ Buat/update file `go.sum`

**Tunggu sampai selesai** (mungkin 1-2 menit pertama kali).

Output yang benar (tidak ada error):
```
go: downloading github.com/gin-gonic/gin v1.12.0
go: downloading gorm.io/gorm v1.25.x
...
```

> ⚠️ **Error?** Pastikan kamu sudah di dalam folder project (`cd airdrop-tracker-api`).

---

## 6. Jalankan Server

### Cara 1 — Pakai Makefile (Recommended ✅)

```bash
make run
```

### Cara 2 — Tanpa Makefile

```bash
go run cmd/server/main.go
```

### Output yang benar:

```
Database connected
Migration done
[GIN-debug] Listening and serving HTTP on :8080
Server running on :8080
Swagger UI: http://localhost:8080/swagger/index.html
```

> 🚨 **JANGAN TUTUP terminal ini!** Server harus tetap jalan.
> Buka terminal baru kalau mau jalanin command lain.

---

## 7. Buka Swagger UI

1. Buka **browser** (Chrome / Firefox / Safari)
2. Ketik di address bar:

```
http://localhost:8080/swagger/index.html
```

3. Harus muncul halaman **Swagger UI** dengan daftar semua API

> ✅ **Kalau halaman muncul = server berjalan dengan benar!**

---

## 8. Test API Pertama Kali

### Step 1: Register Akun

Di Swagger UI:

1. Klik **POST /api/auth/register**
2. Klik tombol **Try it out**
3. Isi **Request body**:
   ```json
   {
     "email": "test@email.com",
     "password": "password123",
     "name": "Bita"
   }
   ```
4. Klik **Execute**
5. Response harusnya: `201 Created`

### Step 2: Login

1. Klik **POST /api/auth/login**
2. Klik **Try it out**
3. Isi body:
   ```json
   {
     "email": "test@email.com",
     "password": "password123"
   }
   ```
4. Klik **Execute**
5. **Copy token** dari response:
   ```json
   {
     "token": "eyJhbGciOiJIUzI1NiIs..."
   }
   ```
   Copy nilai `token` (seluruh string panjang itu).

### Step 3: Authorize

1. Klik tombol **🔓 Authorize** di pojok kanan atas Swagger
2. Di kolom input, ketik:
   ```
   Bearer eyJhbGciOiJIUzI1NiIs...
   ```
   (Ganti dengan token yang kamu copy tadi. **Jangan lupa** kata `Bearer` + spasi di depan!)
3. Klik **Authorize** → **Close**

### Step 4: Pakai API!

Sekarang semua endpoint 🔒 sudah terbuka!

**Coba buat airdrop pertama:**

1. Klik **POST /api/airdrops**
2. Klik **Try it out**
3. Isi body:
   ```json
   {
     "name": "zkSync",
     "chain": "Ethereum",
     "priority": "high",
     "status": "active",
     "url": "https://zksync.io",
     "notes": "Bridge minimal $100 setiap minggu"
   }
   ```
4. Klik **Execute**
5. Response: `201 Created` ✅

**Coba lihat dashboard:**

1. Klik **GET /api/dashboard**
2. Klik **Try it out** → **Execute**
3. Harus muncul summary stats!

---

## 9. Generate Swagger Docs

Swagger docs sudah di-generate dan masuk ke folder `docs/`.
Kalau kamu mengubah comment `@Summary`, `@Description`, dll di handler, **wajib regenerate**:

```bash
make swag
```

Atau tanpa Makefile:

```bash
swag init -g cmd/server/main.go
```

Lalu restart server:

```bash
make run
```

### Install swag CLI (kalau belum punya):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Pastikan `$GOPATH/bin` ada di PATH:

```bash
# Tambahkan di ~/.bashrc atau ~/.zshrc
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.bashrc
```

---

## 10. Docker (Opsional)

Kalau kamu punya Docker, bisa jalankan tanpa install Go:

### Build image:

```bash
docker build -t airdrop-api .
```

### Run container:

```bash
docker run -d \
  --name airdrop-api \
  -p 8080:8080 \
  -v airdrop_data:/app/data \
  --restart unless-stopped \
  airdrop-api
```

### Cek logs:

```bash
docker logs -f airdrop-api
```

### Stop & remove:

```bash
docker stop airdrop-api
docker rm airdrop-api
```

---

## 11. Troubleshooting

### ❌ `command not found: go`

Go belum ter-install atau belum masuk PATH.

```bash
# Cek apakah Go ada di system
which go

# Kalau kosong, install Go dulu (Step 2)
# Kalau sudah install, restart terminal / komputer
```

### ❌ `port already in use: 8080`

Port 8080 sudah dipakai app lain.

```bash
# Cek siapa yang pakai
lsof -i :8080

# Ganti port di .env
echo "APP_PORT=3001" >> .env

# Atau matiin process yang pakai
kill -9 <PID>
```

### ❌ `no such file or directory: data/airdrop.db`

Folder `data/` belum ada. Buat manual:

```bash
mkdir -p data
```

Lalu jalankan lagi: `make run`

### ❌ Swagger UI blank / 404

1. Pastikan server **sudah jalan** (Step 6)
2. Coba regenerate docs:
   ```bash
   make swag
   make run
   ```
3. Buka **baru** di browser: `http://localhost:8080/swagger/index.html`

### ❌ `connection refused` saat test API

Server belum jalan. Pastikan terminal `make run` masih terbuka dan menampilkan `Server running on :8080`.

### ❌ `missing go.sum entry`

Dependency belum lengkap. Jalankan ulang:

```bash
go mod tidy
make run
```

### ❌ `CGO_ENABLED` error (SQLite compile)

SQLite butuh CGO. Kalau error:

```bash
# Linux — install gcc dulu
sudo apt install gcc

# Mac — install Xcode Command Line Tools
xcode-select --install

# Windows — install TDM-GCC atau MinGW
```

### ❌ Lupa password / mau reset database

Hapus file database dan mulai fresh:

```bash
rm data/airdrop.db
make run
```

> ⚠️ **Semua data hilang!** Pastikan backup dulu kalau perlu.

### ❌ Token expired / Unauthorized

Token JWT expired. Login ulang:

1. POST `/api/auth/login` dengan email & password
2. Copy token baru
3. Klik **Authorize** di Swagger → paste token baru

---

## 12. Cheat Sheet

### Sehari-hari

```bash
# Jalankan server
make run

# Build binary (hasil: bin/server)
make build

# Regenerate Swagger setelah edit comment handler
make swag

# Run tests
make test

# Hapus database + binary (fresh start)
make clean
```

### Curl (test dari terminal)

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@email.com","password":"password123","name":"Bita"}'

# Login (copy token dari response)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@email.com","password":"password123"}'

# List airdrops (ganti TOKEN)
curl http://localhost:8080/api/airdrops \
  -H "Authorization: Bearer TOKEN"

# Create airdrop
curl -X POST http://localhost:8080/api/airdrops \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"name":"zkSync","chain":"Ethereum","priority":"high","status":"active"}'

# Dashboard stats
curl http://localhost:8080/api/dashboard \
  -H "Authorization: Bearer TOKEN"

# Export Excel
curl http://localhost:8080/api/export/excel \
  -H "Authorization: Bearer TOKEN" \
  -o export.xlsx
```

---

## 🎯 Ringkasan Cepat (Copy-Paste)

```bash
# Clone + Setup + Run
git clone https://github.com/bsyrlhabibi/airdrop-tracker-api.git
cd airdrop-tracker-api
cp .env.example .env
go mod tidy
make run
```

Buka browser: **http://localhost:8080/swagger/index.html** 🚀

---

Selamat mencoba! 🎉

Ada masalah? Cek [Troubleshooting](#11-troubleshooting) dulu sebelum bertanya.
