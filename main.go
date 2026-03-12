package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/middlewares"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
	APIKey string `mapstructure:"API_KEY"`
}

func main() {
	// ── 1. Load konfigurasi ───────────────────────────────
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
		APIKey: viper.GetString("API_KEY"),
	}

	// ── 2. Koneksi database ───────────────────────────────
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	defer db.Close()

	// ── 3. Dependency Injection ───────────────────────────

	// Product
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Report
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// ── 4. Setup middleware ───────────────────────────────
	// apiKeyMiddleware dibuat sekali, dipakai di banyak route
	apiKeyMiddleware := middlewares.APIKey(config.APIKey)

	// ── 5. Routes ─────────────────────────────────────────
	//
	// Pola pembungkusan middleware (dibaca dari dalam ke luar):
	// CORS → Logger → APIKey → Handler
	//
	// Artinya: setiap request melewati CORS dulu, lalu Logger mencatat,
	// lalu APIKey mengecek, baru sampai ke Handler.

	// GET /api/produk      → public (tanpa API key)
	// POST /api/produk     → public (tanpa API key)
	http.HandleFunc("/api/produk",
		middlewares.CORS(
			middlewares.Logger(
				productHandler.HandleProducts,
			),
		),
	)

	// GET /api/produk/{id}    → public
	// PUT /api/produk/{id}    → butuh API key
	// DELETE /api/produk/{id} → butuh API key
	http.HandleFunc("/api/produk/",
		middlewares.CORS(
			middlewares.Logger(
				apiKeyMiddleware(
					productHandler.HandleProductByID,
				),
			),
		),
	)

	// POST /api/checkout → butuh API key
	http.HandleFunc("/api/checkout",
		middlewares.CORS(
			middlewares.Logger(
				apiKeyMiddleware(
					transactionHandler.HandleCheckout,
				),
			),
		),
	)

	// GET /api/report/hari-ini → public
	http.HandleFunc("/api/report/hari-ini",
		middlewares.CORS(
			middlewares.Logger(
				reportHandler.HandleHariIni,
			),
		),
	)

	// Health check → public
	http.HandleFunc("/health",
		middlewares.CORS(
			middlewares.Logger(
				func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(`{"status":"OK","message":"API Running"}`))
				},
			),
		),
	)

	// ── 6. Start server ───────────────────────────────────
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server jalan di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}