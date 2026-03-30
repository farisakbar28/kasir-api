# Swagger API Documentation

Dokumentasi API Swagger telah berhasil diintegrasikan ke project kasir-api!

## 📌 Akses Dokumentasi

Setelah menjalankan server, akses dokumentasi Swagger di:

```
http://localhost:8080/swagger/index.html
```

## 🚀 Cara Menjalankan

1. **Build project:**

   ```bash
   go build -o kasir-api.exe
   ```

2. **Jalankan server:**

   ```bash
   .\kasir-api.exe
   ```

3. **Buka browser:**
   - Swagger UI: `http://localhost:8080/swagger/index.html`
   - API Health: `http://localhost:8080/health`

## 📝 File-file yang Sudah Diupdate

### 1. **main.go**

- Menambahkan imports untuk `httpSwagger` dan `docs`
- Menambahkan swagger comments untuk API info
- Menambahkan route untuk Swagger UI: `/swagger/`

### 2. **handlers/product_handler.go**

Swagger comments untuk:

- `GetAll()` - GET /api/produk (list semua produk)
- `Create()` - POST /api/produk (buat produk baru)
- `GetByID()` - GET /api/produk/{id} (ambil produk tertentu)
- `Update()` - PUT /api/produk/{id} (update produk)
- `Delete()` - DELETE /api/produk/{id} (hapus produk)

### 3. **handlers/transaction_handler.go**

Swagger comments untuk:

- `Checkout()` - POST /api/checkout (proses checkout)

### 4. **handlers/report_handler.go**

Swagger comments untuk:

- `HandleHariIni()` - GET /api/report/hari-ini (laporan penjualan)

### 5. **docs/ (Auto-generated)**

- `docs.go` - Swagger spec dalam bentuk Go code
- `swagger.json` - Swagger spec dalam format JSON
- `swagger.yaml` - Swagger spec dalam format YAML

## 🔑 Autentikasi API Key

Beberapa endpoint memerlukan API Key. Untuk testing di Swagger:

1. Klik tombol **"Authorize"** di atas dokumentasi (icon gembok)
2. Masukkan API Key Anda di field `X-API-Key`
3. Klik **"Authorize"**

## 📚 API Endpoints yang Terdokumentasi

### Products (Produk)

- **GET /api/produk** - List semua produk (optional filter: `name`)
- **POST /api/produk** - Buat produk baru
- **GET /api/produk/{id}** - Ambil produk tertentu
- **PUT /api/produk/{id}** - Update produk (perlu API Key)
- **DELETE /api/produk/{id}** - Hapus produk (perlu API Key)

### Transactions (Penjualan)

- **POST /api/checkout** - Proses checkout (perlu API Key)

### Reports (Laporan)

- **GET /api/report/hari-ini** - Laporan penjualan (optional: start_date & end_date)

### Health

- **GET /health** - Status API

## 🔄 Update Dokumentasi

Setiap kali Anda menambah endpoint atau mengubah handler, lakukan:

```bash
swag init -g main.go
go build -o kasir-api.exe
```

Swagger documentation akan otomatis terupdate.

## 📖 Format Swagger Comments

Contoh format yang digunakan:

```go
// GetAll godoc
// @Summary Get all products
// @Description Retrieve all products with optional name filter
// @Tags Products
// @Produce json
// @Param name query string false "Product name filter"
// @Success 200 {array} models.Product "List of products"
// @Failure 500 {string} string "Internal server error"
// @Router /produk [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
    // ...
}
```

## 🛠️ Customization

Untuk mengubah informasi API (judul, versi, deskripsi), edit bagian atas file `main.go`:

```go
// @title Kasir API Documentation
// @version 1.0
// @description API untuk sistem manajemen kasir/penjualan
// @host localhost:8080
// @basePath /api
// @schemes http https
```

## 📦 Dependencies Terinstall

- `github.com/swaggo/http-swagger` - Handler untuk Swagger UI
- `github.com/swaggo/swag` - Tool untuk generate Swagger docs
- `github.com/swaggo/files` - File assets untuk Swagger UI

---

**Selamat! Swagger API Documentation sudah siap digunakan!** 🎉
