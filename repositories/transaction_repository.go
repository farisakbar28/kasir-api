package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	// Mulai database transaction
	// Kalau ada error di tengah jalan, semua perubahan dibatalkan (rollback)
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	// defer rollback — kalau fungsi return sebelum commit, semua dibatalkan
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// Step 1: Loop setiap item, cek produk ada atau tidak, hitung subtotal
	for _, item := range items {
		var productPrice, stock int
		var productName string

		// Ambil data produk dari database
		err := tx.QueryRow(
			"SELECT name, price, stock FROM products WHERE id = $1",
			item.ProductID,
		).Scan(&productName, &productPrice, &stock)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("produk dengan id %d tidak ditemukan", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Cek stok cukup atau tidak
		if stock < item.Quantity {
			return nil, fmt.Errorf("stok produk '%s' tidak cukup (stok: %d, diminta: %d)", productName, stock, item.Quantity)
		}

		// Hitung subtotal untuk item ini
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// Step 2: Kurangi stok produk
		_, err = tx.Exec(
			"UPDATE products SET stock = stock - $1 WHERE id = $2",
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Step 3: Simpan header transaksi, dapat ID transaksi
	var transactionID int
	err = tx.QueryRow(
		"INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id",
		totalAmount,
	).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// Step 4: Simpan detail transaksi untuk setiap item
	for i := range details {
		details[i].TransactionID = transactionID

		var detailID int
		err = tx.QueryRow(
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transactionID,
			details[i].ProductID,
			details[i].Quantity,
			details[i].Subtotal,
		).Scan(&detailID)
		if err != nil {
			return nil, err
		}

		details[i].ID = detailID
	}

	// Step 5: Commit — semua perubahan disimpan permanen ke database
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}