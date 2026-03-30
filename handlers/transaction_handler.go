package handlers

import (
	"encoding/json"
	"kasir-api/models"
	"kasir-api/services"
	"net/http"
)

type TransactionHandler struct {
	service *services.TransactionService
}

func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Checkout godoc
// @Summary Process checkout
// @Description Create a new transaction from checkout items (requires API key)
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body models.CheckoutRequest true "Checkout items"
// @Success 201 {object} models.Transaction "Transaction created"
// @Failure 400 {string} string "Bad request or empty items"
// @Failure 500 {string} string "Internal server error"
// @Router /checkout [post]
// @Security ApiKeyAuth
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	// Baca list item dari request body
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Request body tidak valid", http.StatusBadRequest)
		return
	}

	// Validasi minimal ada satu item
	if len(req.Items) == 0 {
		http.Error(w, "Items tidak boleh kosong", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
