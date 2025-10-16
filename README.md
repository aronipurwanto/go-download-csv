# 📦 Transaction API — Go + Fiber + PostgreSQL + CSV Export

Aplikasi ini adalah REST API modular berbasis **Golang** dan **Fiber v2**, yang mengimplementasikan:
- CRUD Transaction (MVC + Clean Architecture + SOLID)
- Middleware request/response validation
- PostgreSQL via GORM
- CSV Export endpoint dengan **auto split file >10KB**
- Configurable `.env` system via Viper
- Structured error handling (error mapper)
- Optional Excel-friendly BOM & CSV chunk download links

---

## 🧱 Struktur Direktori

```
internal/
├── app/                     # Bootstrap & Fiber setup
│   └── app.go
├── config/                  # Config loader (Viper + env)
│   └── config.go
├── deliveries/
│   └── http/
│       ├── transaction_controller.go  # Controller + export CSV
│       ├── error_map.go               # Error mapper (HTTP ↔ domain)
│       └── router.go                  # Route registration
├── domain/
│   └── transaction/
│       ├── entity.go
│       ├── repository.go
│       ├── service.go
│       └── dto.go
├── middleware/
│   ├── validate_body.go
│   └── enforce_response_envelope.go
└── pkg/
    ├── response/
    │   └── response.go
    └── errorsx/
        └── types.go
```

---

## ⚙️ Konfigurasi Environment

Gunakan `.env` atau environment variable langsung.

### Contoh `.env`
```env
APP_NAME=transaction-api
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=transactiondb
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Jakarta

DB_MAX_OPEN_CONNS=20
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=1h

DB_PORT_PUBLIC=5432
PGADMIN_EMAIL=admin@local
PGADMIN_PASSWORD=admin
```

---

## 🚀 Cara Menjalankan

### 1. Jalankan PostgreSQL lokal
```bash
docker run -d --name pg-tx   -e POSTGRES_PASSWORD=postgres   -e POSTGRES_DB=transactiondb   -p 5432:5432 postgres:16
```

### 2. Build & Run
```bash
go mod tidy
go run ./cmd/server
```

Atau jalankan semua stack dengan Docker Compose:

```bash
docker compose up -d --build
```

Aplikasi berjalan di:  
📍 `http://localhost:8080/v1/transactions`

---

## 🧩 Arsitektur & Prinsip

- **Clean Architecture** (Delivery → Service → Repository)
- **SOLID** (Single Responsibility, Dependency Inversion, dll)
- **Middleware:** validasi body & enforcement response envelope
- **Error Mapper:** konsisten 404/409/422/500
- **Configurable:** lewat `.env` via `internal/config/config.go`

---

## 🧰 Endpoint

### CRUD
| Method | Endpoint | Deskripsi |
|---------|-----------|-----------|
| POST | `/v1/transactions` | Buat transaksi baru |
| GET | `/v1/transactions/:id` | Ambil transaksi by ID |
| GET | `/v1/transactions?page=1&size=10` | Daftar transaksi |
| PUT | `/v1/transactions/:id` | Update transaksi |
| DELETE | `/v1/transactions/:id` | Hapus transaksi |

### Export CSV
`GET /v1/transactions/export.csv`

| Param | Keterangan |
|--------|-------------|
| `from` | Filter tanggal awal (`YYYY-MM-DD` / RFC3339) |
| `to` | Filter tanggal akhir |
| `part` | Unduh bagian tertentu (jika split) |
| `excel=true` | Tambahkan BOM UTF-8 agar mudah dibuka di Excel |

#### Mode Auto Split
- ≤10KB ⇒ 1 file CSV langsung diunduh  
- >10KB ⇒ server membalas JSON daftar link (part 1..N)

```json
{
  "success": true,
  "data": {
    "links": [
      "http://localhost:8080/v1/transactions/export.csv?part=1",
      "http://localhost:8080/v1/transactions/export.csv?part=2"
    ]
  },
  "meta": {
    "total_bytes_estimate": 102400,
    "chunk_limit_bytes": 10240,
    "num_parts": 10
  }
}
```

---

## 🧱 Docker Compose Setup

Jalankan PostgreSQL + API + pgAdmin:

```bash
docker compose up -d --build
```

Layanan:
- API: `http://localhost:8080`
- pgAdmin: `http://localhost:5050` (login pakai `admin@local / admin`)

---

## 🧠 Fitur Tambahan

- **Split CSV by size (10KB each part)**
- **Auto-manifest download links**
- **Excel UTF-8 BOM support (?excel=true)**
- **Response envelope & validation middleware**
- **Clean Config Loader via Viper**
- **Docker-ready architecture**

---

## 🧪 Testing

Gunakan `curl` atau `httpie`:

```bash
# Create
curl -X POST http://localhost:8080/v1/transactions   -H 'Content-Type: application/json'   -d '{"transaction_id":"TX001","amount":100000,"status":"SUCCESS"}'

# Export CSV
curl -L 'http://localhost:8080/v1/transactions/export.csv' -o tx.csv
```

---

## 🧾 Lisensi
MIT © 2025 — Ahmad Roni Purwanto
