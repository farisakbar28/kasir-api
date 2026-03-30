# Kasir API 🛒

REST API sistem kasir sederhana yang dibangun dengan **Go** (tanpa framework), menggunakan **PostgreSQL** (Supabase) sebagai database, dan di-deploy di **Railway**.

---

## 📋 Daftar Isi

- [Tech Stack](#tech-stack)
- [Arsitektur](#arsitektur)
- [Database Schema](#database-schema)
- [API Documentation](#-api-documentation)
- [Instalasi & Setup](#instalasi--menjalankan-lokal)
- [Deployment](#deployment)

---

## Tech Stack

| Teknologi   | Keterangan                                   |
| ----------- | -------------------------------------------- |
| Go (Golang) | Bahasa pemrograman utama                     |
| net/http    | HTTP server bawaan Go, tanpa framework       |
| PostgreSQL  | Database via Supabase (cloud)                |
| Viper       | Manajemen konfigurasi & environment variable |
| lib/pq      | Driver PostgreSQL untuk Go                   |
| Railway     | Platform deployment cloud                    |

---

## Arsitektur

Project ini menggunakan **Layered Architecture** — setiap folder punya tanggung jawab yang terpisah dan jelas.

```
kasir-api/
├── database/              # Inisialisasi koneksi database
├── handlers/              # Menerima HTTP request & mengirim response
├── middlewares/           # API Key, CORS, Logger
├── models/                # Definisi struct/tipe data
├── repositories/          # Query SQL ke database
├── services/              # Logika bisnis
├── docs/                  # Swagger documentation
├── .env                   # Konfigurasi lokal (tidak di-commit)
├── go.mod
└── main.go                # Entry point & dependency injection
```

### Alur Request

```
HTTP Request
    │
    ▼
Middleware (CORS → Logger → API Key)
    │
    ▼
Handler  →  Service  →  Repository  →  Database
    │
    ▼
HTTP Response (JSON)
```

### Penjelasan Setiap Layer

- **Handler** — menerima request, validasi input, kirim response
- **Service** — logika bisnis (validasi stok, hitung total, dll)
- **Repository** — semua query SQL ke database
- **Model** — definisi bentuk data (struct)
- **Middleware** — dijalankan sebelum handler (auth, logging, cors)

---

## Database Schema

```sql
-- Tabel produk
CREATE TABLE products (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR NOT NULL,
    price      INT4 NOT NULL,
    stock      INT4 NOT NULL,
    deleted_at TIMESTAMP DEFAULT NULL  -- soft delete
);

-- Tabel transaksi (header)
CREATE TABLE transactions (
    id           SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabel detail transaksi (per item)
CREATE TABLE transaction_details (
    id             SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id     INT REFERENCES products(id),
    quantity       INT NOT NULL,
    subtotal       INT NOT NULL
);
```

---

## 📚 API Documentation

### Base URL

**Development:** `http://localhost:8080/api`

**Production:** Update sesuai dengan domain Anda di Railway

### Authentication

API menggunakan **API Key** untuk autentikasi endpoint tertentu.

- Header: `X-API-Key`
- Tipe: String
- Contoh: `X-API-Key: your-secret-api-key-here`

**Endpoint yang memerlukan API Key:**

- `POST /checkout` (Checkout/Transaksi)

**Endpoint yang PUBLIC (tanpa API Key):**

- `GET /produk` (List Produk)
- `POST /produk` (Buat Produk)
- `GET /produk/{id}` (Detail Produk)
- `PUT /produk/{id}` (Update Produk) - **PERLU API KEY**
- `DELETE /produk/{id}` (Hapus Produk) - **PERLU API KEY**
- `GET /report/hari-ini` (Laporan Penjualan)

### Response Format

Semua response menggunakan format **JSON**.

**Sukses (2xx):**

```json
{
  "id": 1,
  "name": "Indomie",
  "price": 5000,
  "stock": 100,
  "created_at": "2026-03-31T10:30:00Z"
}
```

**Error (4xx, 5xx):**

```json
{
  "error": "Deskripsi error",
  "status": 400
}
```

### Endpoints

#### 1️⃣ **PRODUCTS (Produk)**

##### 📌 Lihat Semua Produk

```http
GET /api/produk
```

**Query Parameters:**

- `name` (optional) - Filter produk berdasarkan nama

**Contoh Request:**

```bash
curl http://localhost:8080/api/produk
curl "http://localhost:8080/api/produk?name=indomie"
```

**Response (200 OK):**

```json
[
  {
    "id": 1,
    "name": "Indomie Goreng",
    "price": 4500,
    "stock": 50,
    "created_at": "2026-03-01T08:00:00Z"
  },
  {
    "id": 2,
    "name": "Indomie Kuah",
    "price": 4500,
    "stock": 75,
    "created_at": "2026-03-01T08:30:00Z"
  }
]
```

---

##### 📌 Buat Produk Baru

```http
POST /api/produk
Content-Type: application/json
```

**Request Body:**

```json
{
  "name": "Indomie Goreng",
  "price": 4500,
  "stock": 100
}
```

**Contoh Request:**

```bash
curl -X POST http://localhost:8080/api/produk \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Indomie Goreng",
    "price": 4500,
    "stock": 100
  }'
```

**Response (201 Created):**

```json
{
  "id": 1,
  "name": "Indomie Goreng",
  "price": 4500,
  "stock": 100,
  "created_at": "2026-03-31T10:30:00Z"
}
```

**Error (400 Bad Request):**

```
Name tidak boleh kosong
```

---

##### 📌 Lihat Detail Produk

```http
GET /api/produk/{id}
```

**Contoh Request:**

```bash
curl http://localhost:8080/api/produk/1
```

**Response (200 OK):**

```json
{
  "id": 1,
  "name": "Indomie Goreng",
  "price": 4500,
  "stock": 50,
  "created_at": "2026-03-01T08:00:00Z"
}
```

**Error (404 Not Found):**

```
Produk tidak ditemukan
```

---

##### 📌 Update Produk

```http
PUT /api/produk/{id}
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body:**

```json
{
  "name": "Indomie Goreng Pedas",
  "price": 5000,
  "stock": 75
}
```

**Contoh Request:**

```bash
curl -X PUT http://localhost:8080/api/produk/1 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "name": "Indomie Goreng Pedas",
    "price": 5000,
    "stock": 75
  }'
```

**Response (200 OK):**

```json
{
  "id": 1,
  "name": "Indomie Goreng Pedas",
  "price": 5000,
  "stock": 75,
  "created_at": "2026-03-01T08:00:00Z"
}
```

---

##### 📌 Hapus Produk

```http
DELETE /api/produk/{id}
X-API-Key: your-secret-api-key
```

**Contoh Request:**

```bash
curl -X DELETE http://localhost:8080/api/produk/1 \
  -H "X-API-Key: your-secret-api-key"
```

**Response (200 OK):**

```
Produk berhasil dihapus
```

---

#### 2️⃣ **TRANSACTIONS (Transaksi/Checkout)**

##### 📌 Proses Checkout (Buat Transaksi)

```http
POST /api/checkout
Content-Type: application/json
X-API-Key: your-secret-api-key
```

**Request Body:**

```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 2,
      "quantity": 1
    }
  ]
}
```

**Contoh Request:**

```bash
curl -X POST http://localhost:8080/api/checkout \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 2, "quantity": 1}
    ]
  }'
```

**Response (201 Created):**

```json
{
  "id": 1,
  "items": [
    {
      "product_id": 1,
      "product_name": "Indomie Goreng",
      "quantity": 2,
      "price": 4500,
      "subtotal": 9000
    },
    {
      "product_id": 2,
      "product_name": "Indomie Kuah",
      "quantity": 1,
      "price": 4500,
      "subtotal": 4500
    }
  ],
  "total_amount": 13500,
  "created_at": "2026-03-31T10:35:00Z"
}
```

**Error (400 Bad Request):**

```
Items tidak boleh kosong
```

**Error (500 Internal Server Error):**

```
Stok produk tidak cukup
```

---

#### 3️⃣ **REPORTS (Laporan)**

##### 📌 Laporan Penjualan Hari Ini

```http
GET /api/report/hari-ini
```

**Query Parameters:**

- `start_date` (optional) - Format: `YYYY-MM-DD` (contoh: `2026-03-31`)
- `end_date` (optional) - Format: `YYYY-MM-DD` (contoh: `2026-03-31`)

> Jika tidak ada parameter, dashboard menampilkan data hari ini saja.

**Contoh Request:**

```bash
# Laporan hari ini
curl http://localhost:8080/api/report/hari-ini

# Laporan range tanggal
curl "http://localhost:8080/api/report/hari-ini?start_date=2026-03-01&end_date=2026-03-31"
```

**Response (200 OK):**

```json
{
  "total_sales": 10,
  "total_amount": 45000,
  "top_products": [
    {
      "product_id": 1,
      "product_name": "Indomie Goreng",
      "quantity_sold": 5,
      "total_revenue": 22500
    },
    {
      "product_id": 2,
      "product_name": "Indomie Kuah",
      "quantity_sold": 3,
      "total_revenue": 13500
    }
  ]
}
```

---

#### 4️⃣ **HEALTH CHECK**

##### 📌 Status API

```http
GET /health
```

**Response (200 OK):**

```json
{
  "status": "OK",
  "message": "API Running"
}
```

---

### Error Codes

| Code | Pesan                 | Arti                                    |
| ---- | --------------------- | --------------------------------------- |
| 200  | OK                    | Request berhasil                        |
| 201  | Created               | Resource berhasil dibuat                |
| 400  | Bad Request           | Request invalid atau data tidak lengkap |
| 401  | Unauthorized          | API Key salah atau tidak ada            |
| 404  | Not Found             | Resource tidak ditemukan                |
| 405  | Method Not Allowed    | HTTP method tidak sesuai                |
| 500  | Internal Server Error | Error di server                         |

---

### Swagger UI Documentation

Dokumentasi API interaktif tersedia di Swagger UI:

```
http://localhost:8080/swagger/index.html
```

Gunakan interface Swagger untuk test endpoint secara langsung.

---

## Instalasi & Menjalankan Lokal

### Requirements

- Go versi 1.19 atau lebih baru
- Akun Supabase (gratis)

### Langkah-langkah

**1. Clone repository**

```bash
git clone https://github.com/username/kasir-api.git
cd kasir-api
```

**2. Install dependencies**

```bash
go mod tidy
```

**3. Buat file .env**

```bash
cp .env.example .env
```

Isi file `.env` dengan nilai yang sesuai:

```
PORT=8080
DB_CONN=postgresql://postgres.xxxx:password@aws-x-ap-southeast-1.pooler.supabase.com:6543/postgres
API_KEY=your-secret-api-key-here
```

**4. Jalankan server**

```bash
go run main.go
```

Output:

```
Database connected successfully
Server jalan di 0.0.0.0:8080
```

**5. Test health check**

```bash
curl http://localhost:8080/health
```

---

## Environment Variables

| Variable  | Keterangan                                 | Contoh                  |
| --------- | ------------------------------------------ | ----------------------- |
| `PORT`    | Port server                                | `8080`                  |
| `DB_CONN` | PostgreSQL connection string dari Supabase | `postgresql://...`      |
| `API_KEY` | Secret key untuk endpoint yang diproteksi  | `kasir-api-secret-2026` |

---

## Middleware

### 1. CORS

Mengizinkan akses dari domain yang berbeda (dibutuhkan untuk frontend).
Diterapkan ke semua endpoint.

### 2. Logger

Mencatat setiap request yang masuk beserta durasi eksekusi.

Contoh output:

```
[REQUEST] POST /api/checkout dari 127.0.0.1:54321
[DONE]    POST /api/checkout selesai dalam 45.231ms
```

### 3. API Key

Memproteksi endpoint sensitif. Request harus menyertakan header:

```
X-API-Key: your-secret-api-key-here
```

Request tanpa API key atau dengan API key salah akan mendapat response:

```
401 Unauthorized
```

---

## API Reference

### Base URL

**Lokal:**

```
http://localhost:8080
```

**Production:**

```
https://kasir-api-production-d862.up.railway.app
```

---

### Health Check

#### `GET /health`

Cek apakah server berjalan.

**Auth:** Tidak diperlukan

**Request:**

```bash
curl https://kasir-api-production-d862.up.railway.app/health
```

**Response `200 OK`:**

```json
{
  "status": "OK",
  "message": "API Running"
}
```

---

### Produk

#### `GET /api/produk`

Ambil semua produk. Bisa difilter berdasarkan nama.

**Auth:** Tidak diperlukan

**Query Parameters:**

| Parameter | Tipe   | Wajib | Keterangan                               |
| --------- | ------ | ----- | ---------------------------------------- |
| `name`    | string | Tidak | Filter produk by nama (case-insensitive) |

**Request — ambil semua:**

```bash
curl https://kasir-api-production-d862.up.railway.app/api/produk
```

**Request — search by nama:**

```bash
curl "https://kasir-api-production-d862.up.railway.app/api/produk?name=indom"
```

**Response `200 OK`:**

```json
[
  {
    "id": 1,
    "name": "Indomie Godog",
    "price": 3500,
    "stock": 10
  },
  {
    "id": 2,
    "name": "Indomie Goreng",
    "price": 3000,
    "stock": 20
  }
]
```

**Response jika tidak ada produk:**

```json
[]
```

---

#### `GET /api/produk/{id}`

Ambil satu produk berdasarkan ID.

**Auth:** Diperlukan (`X-API-Key`)

**Path Parameters:**

| Parameter | Tipe    | Keterangan |
| --------- | ------- | ---------- |
| `id`      | integer | ID produk  |

**Request:**

```bash
curl https://kasir-api-production-d862.up.railway.app/api/produk/1 \
  -H "X-API-Key: your-secret-api-key-here"
```

**Response `200 OK`:**

```json
{
  "id": 1,
  "name": "Indomie Godog",
  "price": 3500,
  "stock": 10
}
```

**Response `404 Not Found`:**

```
produk tidak ditemukan
```

**Response `400 Bad Request`:**

```
ID tidak valid
```

---

#### `POST /api/produk`

Tambah produk baru.

**Auth:** Tidak diperlukan

**Request Body:**

| Field   | Tipe    | Wajib | Keterangan         |
| ------- | ------- | ----- | ------------------ |
| `name`  | string  | Ya    | Nama produk        |
| `price` | integer | Ya    | Harga dalam rupiah |
| `stock` | integer | Ya    | Jumlah stok awal   |

**Request:**

```bash
curl -X POST https://kasir-api-production-d862.up.railway.app/api/produk \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kopi Kapal Api",
    "price": 2500,
    "stock": 200
  }'
```

**Response `201 Created`:**

```json
{
  "id": 4,
  "name": "Kopi Kapal Api",
  "price": 2500,
  "stock": 200
}
```

---

#### `PUT /api/produk/{id}`

Update data produk berdasarkan ID.

**Auth:** Diperlukan (`X-API-Key`)

**Path Parameters:**

| Parameter | Tipe    | Keterangan                   |
| --------- | ------- | ---------------------------- |
| `id`      | integer | ID produk yang akan diupdate |

**Request Body:**

| Field   | Tipe    | Wajib | Keterangan              |
| ------- | ------- | ----- | ----------------------- |
| `name`  | string  | Ya    | Nama produk baru        |
| `price` | integer | Ya    | Harga baru dalam rupiah |
| `stock` | integer | Ya    | Stok baru               |

**Request:**

```bash
curl -X PUT https://kasir-api-production-d862.up.railway.app/api/produk/1 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key-here" \
  -d '{
    "name": "Indomie Goreng Jumbo",
    "price": 4000,
    "stock": 150
  }'
```

**Response `200 OK`:**

```json
{
  "id": 1,
  "name": "Indomie Goreng Jumbo",
  "price": 4000,
  "stock": 150
}
```

**Response `404 Not Found`:**

```
produk tidak ditemukan
```

---

#### `DELETE /api/produk/{id}`

Hapus produk berdasarkan ID (soft delete — data tidak benar-benar dihapus dari database).

**Auth:** Diperlukan (`X-API-Key`)

**Path Parameters:**

| Parameter | Tipe    | Keterangan                  |
| --------- | ------- | --------------------------- |
| `id`      | integer | ID produk yang akan dihapus |

**Request:**

```bash
curl -X DELETE https://kasir-api-production-d862.up.railway.app/api/produk/1 \
  -H "X-API-Key: your-secret-api-key-here"
```

**Response `200 OK`:**

```json
{
  "message": "Product deleted successfully"
}
```

**Response `404 Not Found`:**

```
produk tidak ditemukan
```

> **Catatan:** Delete menggunakan soft delete. Produk ditandai `deleted_at` dengan timestamp dan tidak akan muncul di endpoint GET, namun data tetap tersimpan di database untuk menjaga integritas riwayat transaksi.

---

### Transaksi

#### `POST /api/checkout`

Buat transaksi baru. Sistem akan menghitung total harga dan mengurangi stok produk secara otomatis.

**Auth:** Diperlukan (`X-API-Key`)

**Request Body:**

| Field                | Tipe    | Wajib | Keterangan              |
| -------------------- | ------- | ----- | ----------------------- |
| `items`              | array   | Ya    | Daftar item yang dibeli |
| `items[].product_id` | integer | Ya    | ID produk               |
| `items[].quantity`   | integer | Ya    | Jumlah yang dibeli      |

**Request:**

```bash
curl -X POST https://kasir-api-production-d862.up.railway.app/api/checkout \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key-here" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 2, "quantity": 1}
    ]
  }'
```

**Response `201 Created`:**

```json
{
  "id": 1,
  "total_amount": 10000,
  "created_at": "2026-03-05T13:00:00Z",
  "details": [
    {
      "id": 1,
      "transaction_id": 1,
      "product_id": 1,
      "product_name": "Indomie Godog",
      "quantity": 2,
      "subtotal": 7000
    },
    {
      "id": 2,
      "transaction_id": 1,
      "product_id": 2,
      "product_name": "Vit 1000ml",
      "quantity": 1,
      "subtotal": 3000
    }
  ]
}
```

**Response `500` jika stok tidak cukup:**

```
stok produk 'Indomie Godog' tidak cukup (stok: 1, diminta: 5)
```

**Response `500` jika produk tidak ditemukan:**

```
produk dengan id 99 tidak ditemukan
```

> **Catatan:** Checkout menggunakan **database transaction**. Jika salah satu item gagal (produk tidak ada, stok kurang), seluruh transaksi dibatalkan dan stok tidak berubah.

---

### Laporan

#### `GET /api/report/hari-ini`

Ambil laporan penjualan. Default menampilkan data hari ini. Bisa difilter dengan query parameter untuk rentang tanggal tertentu.

**Auth:** Tidak diperlukan

**Query Parameters:**

| Parameter    | Tipe   | Wajib | Keterangan                          |
| ------------ | ------ | ----- | ----------------------------------- |
| `start_date` | string | Tidak | Tanggal mulai, format: `2026-01-01` |
| `end_date`   | string | Tidak | Tanggal akhir, format: `2026-03-05` |

**Request — laporan hari ini:**

```bash
curl https://kasir-api-production-d862.up.railway.app/api/report/hari-ini
```

**Request — laporan by rentang tanggal:**

```bash
curl "https://kasir-api-production-d862.up.railway.app/api/report/hari-ini?start_date=2026-01-01&end_date=2026-03-05"
```

**Response `200 OK`:**

```json
{
  "total_revenue": 45000,
  "total_transaksi": 5,
  "produk_terlaris": {
    "nama": "Indomie Godog",
    "qty_terjual": 12
  }
}
```

**Response jika belum ada transaksi hari ini:**

```json
{
  "total_revenue": 0,
  "total_transaksi": 0,
  "produk_terlaris": null
}
```

---

## HTTP Status Codes

| Code  | Nama                  | Kapan digunakan                                      |
| ----- | --------------------- | ---------------------------------------------------- |
| `200` | OK                    | Request berhasil                                     |
| `201` | Created               | Resource baru berhasil dibuat                        |
| `400` | Bad Request           | Input tidak valid atau JSON rusak                    |
| `401` | Unauthorized          | API key tidak ada atau salah                         |
| `404` | Not Found             | Resource tidak ditemukan                             |
| `405` | Method Not Allowed    | HTTP method tidak didukung                           |
| `500` | Internal Server Error | Error di server (stok kurang, produk tidak ada, dll) |

---

## Build & Deployment

### Build Binary Lokal

```bash
# Build standar
go build -o kasir-api

# Build optimized (ukuran lebih kecil)
go build -ldflags="-s -w" -o kasir-api

# Jalankan binary
./kasir-api
```

### Cross Compilation

```bash
# Build untuk Linux (dari Windows/Mac)
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o kasir-api

# Build untuk Windows (dari Mac/Linux)
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o kasir-api.exe

# Build untuk macOS (dari Windows/Linux)
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o kasir-api
```

### Deploy ke Railway

1. Push kode ke GitHub
2. Buka [railway.app](https://railway.app) dan connect repository
3. Set environment variables di Railway dashboard:
   ```
   PORT=8080
   DB_CONN=postgresql://...
   API_KEY=your-secret-api-key-here
   ```
4. Railway otomatis detect Go dan deploy

---

## Catatan Pengembangan

### Soft Delete

Produk menggunakan soft delete. Saat DELETE dipanggil, kolom `deleted_at` diisi timestamp, bukan dihapus dari database. Ini untuk menjaga integritas data riwayat transaksi — jika produk dihapus sungguhan, data transaksi lama yang mereferensikan produk tersebut akan rusak.

### Database Transaction pada Checkout

Proses checkout menggunakan database transaction (`BEGIN` / `COMMIT` / `ROLLBACK`). Jika ada satu item yang gagal diproses (produk tidak ada, stok kurang), seluruh operasi dibatalkan dan tidak ada perubahan yang tersimpan ke database.

### In-Memory vs Database

Sesi 1 project ini menggunakan in-memory storage (data hilang saat restart). Sejak sesi 2, semua data disimpan di PostgreSQL via Supabase sehingga data tetap ada meski server di-restart atau di-redeploy.
