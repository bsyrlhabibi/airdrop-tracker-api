# GUIDE.md — Cara Menjalankan Airdrop Tracker API

Panduan lengkap untuk menjalankan project ini di komputer kamu, dari nol sampai bisa pakai Swagger.

---

## Daftar Isi

1. [Persiapan](#1-persiapan)
2. [Install Go](#2-install-go)
3. [Clone Project](#3-clone-project)
4. [Setup Environment](#4-setup-environment)
5. [Jalankan Server](#5-jalankan-server)
6. [Buka Swagger UI](#6-buka-swagger-ui)
7. [Test API Manual](#7-test-api-manual)
8. [Troubleshooting](#8-troubleshooting)

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

## 5. Jalankan Server

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

## 6. Buka Swagger UI

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

## 7. Test API Manual (Terminal)

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

## 8. Troubleshooting

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
→ Pastikan server jalan (Step 5). Coba regenerate docs:
```bash
swag init -g cmd/server/main.go
make run
```

### "connection refused" saat test API
→ Server belum jalan. Pastikan Terminal tempat `make run` masih terbuka dan menampilkan `Listening on :8080`.

### Lupa password / token expired
→ Register ulang dengan email baru, atau hapus database dan mulai fresh:
```bash
rm data/airdrop.db
make run
```

---

## Perintah Berguna

```bash
make run          # Jalankan server
make build        # Build binary (hasilnya di bin/server)
make clean        # Hapus database & binary
make test         # Jalankan test (jika ada)

swag init -g cmd/server/main.go   # Regenerate Swagger docs
```

---

## Butuh Bantuan?

- Buka issue di GitHub
- Atau tanya di grup project

---

Selamat mencoba! 🚀
