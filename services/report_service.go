package services

import (
	"kasir-api/repositories"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// GetHariIni — laporan dari jam 00:00 sampai 23:59 hari ini
func (s *ReportService) GetHariIni() (*repositories.DailySummary, error) {
	now := time.Now()

	// Awal hari: hari ini jam 00:00:00
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Akhir hari: besok jam 00:00:00 (exclusive)
	end := start.Add(24 * time.Hour)

	return s.repo.GetDailySummary(start, end)
}

// GetByDateRange — laporan berdasarkan rentang tanggal (optional challenge)
func (s *ReportService) GetByDateRange(startDate, endDate time.Time) (*repositories.DailySummary, error) {
	return s.repo.GetDailySummary(startDate, endDate)
}