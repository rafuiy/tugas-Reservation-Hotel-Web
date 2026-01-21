# Sistem Reservasi Hotel

Project ini adalah aplikasi reservasi hotel sederhana berbasis Go (Gin) dengan UI server-rendered (html/template) dan API JSON. Tidak ada frontend terpisah.

## Prasyarat
- Go >= 1.22
- Postgres lokal (tanpa Docker) atau Docker + Docker Compose

## Konfigurasi
Salin `.env.example` menjadi `.env`, lalu sesuaikan nilainya.

```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/reservasi_hotel?sslmode=disable
PORT=8080
JWT_SECRET=isi_secret_yang_kuat
ADMIN_SEED_EMAIL=admin@local.test
ADMIN_SEED_PASSWORD=admin123
```

## Menjalankan Database (Tanpa Docker)

Pastikan Postgres lokal berjalan dan `psql` tersedia di PATH.

Buat database:

```
psql -U postgres -d postgres -c "CREATE DATABASE reservasi_hotel;"
```

Jalankan migrasi:

```
psql -U postgres -d reservasi_hotel -f migrations/001_init.sql
```

## Menjalankan Database (Docker)

```
docker compose up -d
```

## Migrasi
Jika menggunakan Docker, migrasi dijalankan via `psql` di container database:

```
docker compose exec db psql -U postgres -d reservasi_hotel -f /migrations/001_init.sql
```

## Seed Admin

```
go run cmd/seed/main.go
```

Admin default (sesuai `.env`):
- Email: `admin@local.test`
- Password: `admin123`

## Menjalankan Server

```
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8080` (atau sesuai `PORT`).

## Flow Minimal (Ringkas)
1) Register user (UI `/register` atau API `/api/auth/register`).
2) Seed admin lalu login admin.
3) Admin membuat room + availability.
4) User membuat booking (status PENDING).
5) Admin approve booking.
6) User membuat payment dan menandai PAID.
7) Invoice dapat diakses via `/invoice/:booking_id` atau API.

## Catatan Ambiguitas
- Amount payment di-set dari `price_per_slot` room pada saat create payment (server-side), jadi input amount dari user diabaikan demi konsistensi.
- Status/enum memakai `CHECK constraint` di Postgres (tanpa tipe ENUM custom).

## Menjalankan Test

```
go test ./...
```

## Contoh Curl (Opsional)

```
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"User A","email":"usera@test.local","password":"secret123","phone":"081234"}'

curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"usera@test.local","password":"secret123"}'
```

