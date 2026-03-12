package middlewares

import (
	"log"
	"net/http"
	"time"
)

// Logger mencatat setiap request yang masuk beserta durasi eksekusinya.
// Berguna untuk monitoring dan debugging.
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Catat waktu mulai
		start := time.Now()

		// Log request yang masuk
		log.Printf("[REQUEST] %s %s dari %s", r.Method, r.RequestURI, r.RemoteAddr)

		// Jalankan handler
		next(w, r)

		// Hitung dan log durasi setelah handler selesai
		duration := time.Since(start)
		log.Printf("[DONE]    %s %s selesai dalam %v", r.Method, r.RequestURI, duration)
	}
}