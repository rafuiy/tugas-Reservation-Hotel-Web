# Sistem Reservasi Hotel

Aplikasi reservasi hotel sederhana berbasis Go (Gin) dengan UI server-rendered (html/template).

## Prasyarat
- Go >= 1.22
- Postgres 18

## Panduan Instalasi (copy-paste)

1) Clone repo
```bash
git clone <URL_REPO_KAMU>
cd <NAMA FOLDER KAMU>
```

2) Set env
```bash
copy .env.example .env
```
Edit `.env` sesuai kebutuhan. Contoh:
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/reservasi_hotel?sslmode=disable
PORT=8080
JWT_SECRET=isi_secret_yang_kuat
ADMIN_SEED_EMAIL=admin@local.test
ADMIN_SEED_PASSWORD=admin123
```

3) Buat database (pilih salah satu)

Via terminal (psql):
```bash
psql -U postgres -d postgres -c "CREATE DATABASE reservasi_hotel;"
```

Atau via pgAdmin:
- Create Database → nama: `reservasi_hotel`

4) Jalankan seeder
```bash
go run ./cmd/seed/main.go
```

5) Jalankan aplikasi
```bash
go run ./cmd/server/main.go
```

Server berjalan di `http://localhost:8080` (atau sesuai `PORT`).

## Akun Admin Default (sesuai .env)
- Email: `admin@local.test`
- Password: `admin123`
