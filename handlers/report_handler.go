package handlers

import (
	"encoding/json"
	"kasir-api/services"
	"net/http"
	"time"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleHariIni(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Cek apakah ada query parameter start_date dan end_date
	// Kalau ada, pakai date range. Kalau tidak, pakai hari ini.
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	var summary interface{}
	var err error

	if startStr != "" && endStr != "" {
		// Optional challenge: filter by date range
		// Format tanggal yang diterima: 2026-01-01
		start, errParse := time.Parse("2006-01-02", startStr)
		if errParse != nil {
			http.Error(w, "Format start_date tidak valid. Gunakan format: 2026-01-01", http.StatusBadRequest)
			return
		}

		end, errParse := time.Parse("2006-01-02", endStr)
		if errParse != nil {
			http.Error(w, "Format end_date tidak valid. Gunakan format: 2026-01-01", http.StatusBadRequest)
			return
		}

		// end_date ditambah 1 hari agar tanggal end_date ikut terhitung
		end = end.Add(24 * time.Hour)

		summary, err = h.service.GetByDateRange(start, end)
	} else {
		// Default: laporan hari ini
		summary, err = h.service.GetHariIni()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}