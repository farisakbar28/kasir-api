package middlewares

import "net/http"

// APIKey adalah middleware yang memvalidasi API key dari header request.
// Hanya request yang membawa header "X-API-Key" dengan nilai yang benar
// yang diizinkan melanjutkan ke handler.
func APIKey(validAPIKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Ambil API key dari header request
			apiKey := r.Header.Get("X-API-Key")

			// Cek apakah header ada
			if apiKey == "" {
				http.Error(w, "API key required", http.StatusUnauthorized)
				return
			}

			// Cek apakah API key cocok
			if apiKey != validAPIKey {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// API key valid — lanjutkan ke handler
			next(w, r)
		}
	}
}