package middlewares

import "net/http"

// CORS menambahkan header yang diperlukan agar browser bisa mengakses API
// dari domain yang berbeda (cross-origin request).
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Izinkan semua origin mengakses API ini
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Izinkan method-method HTTP yang kita pakai
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Izinkan header yang dibutuhkan — termasuk X-API-Key untuk auth
		w.Header().Set("Access-Control-Allow-Headers", "X-API-Key, Content-Type")

		// Handle preflight request dari browser
		// Browser mengirim OPTIONS dulu sebelum request asli untuk cek izin CORS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Lanjutkan ke handler berikutnya
		next(w, r)
	}
}