package repositories

import (
	"database/sql"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

type DailySummary struct {
	TotalRevenue    int          `json:"total_revenue"`
	TotalTransaksi  int          `json:"total_transaksi"`
	ProdukTerlaris  *TopProduct  `json:"produk_terlaris"`
}

type TopProduct struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

func (repo *ReportRepository) GetDailySummary(start, end time.Time) (*DailySummary, error) {
	summary := &DailySummary{}

	// Query 1: Total revenue dan jumlah transaksi dalam rentang waktu
	err := repo.db.QueryRow(`
		SELECT 
			COALESCE(SUM(total_amount), 0),
			COUNT(id)
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`, start, end).Scan(&summary.TotalRevenue, &summary.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// Query 2: Produk terlaris (yang paling banyak terjual dalam rentang waktu)
	var topProduct TopProduct
	err = repo.db.QueryRow(`
		SELECT 
			p.name,
			SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, start, end).Scan(&topProduct.Nama, &topProduct.QtyTerjual)

	if err == sql.ErrNoRows {
		// Tidak ada transaksi hari ini — produk terlaris null
		summary.ProdukTerlaris = nil
	} else if err != nil {
		return nil, err
	} else {
		summary.ProdukTerlaris = &topProduct
	}

	return summary, nil
}