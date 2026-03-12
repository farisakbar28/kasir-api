package models

import "time"

// Transaction adalah header transaksi — menyimpan total dan waktu
type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

// TransactionDetail adalah detail per item dalam satu transaksi
type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"` // omitempty = tidak muncul di JSON kalau kosong
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

// CheckoutItem adalah satu item yang dikirim dari kasir
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// CheckoutRequest adalah body request dari kasir — berisi banyak item
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}