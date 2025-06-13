package main

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Rizz404/midtrans-handler/internal/database"
	"github.com/Rizz404/midtrans-handler/internal/enums"
)

// MidtransNotificationPayload merepresentasikan data yang dikirim oleh Midtrans
type MidtransNotificationPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	PaymentType       string `json:"payment_type"`
}

func (apiCfg *apiConfig) handlerMidtransWebhook(w http.ResponseWriter, r *http.Request) {
	// 1. Baca body dari request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not read request body")
		return
	}
	defer r.Body.Close()

	// 2. Unmarshal body ke struct payload
	var payload MidtransNotificationPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid notification payload")
		return
	}

	// 3. Verifikasi Signature Key (SANGAT PENTING UNTUK KEAMANAN)
	signature := generateSignatureKey(payload.OrderID, payload.StatusCode, payload.GrossAmount, apiCfg.MidtransServerKey)
	if signature != payload.SignatureKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid signature")
		return
	}

	// 4. Siapkan data untuk update order
	updateReq := database.UpdateOrderRequest{}
	var paymentStatus enums.PaymentStatus
	var orderStatus enums.OrderStatus

	// 5. Tentukan status baru berdasarkan status transaksi dari Midtrans
	switch payload.TransactionStatus {
	case "capture":
		if payload.FraudStatus == "accept" {
			paymentStatus = enums.PaymentStatusPaid
			orderStatus = enums.OrderStatusPreparing
		}
	case "settlement":
		paymentStatus = enums.PaymentStatusPaid
		orderStatus = enums.OrderStatusPreparing
	case "pending":
		paymentStatus = enums.PaymentStatusUnpaid
		orderStatus = enums.OrderStatusPending
	case "deny", "cancel":
		paymentStatus = enums.PaymentStatusUnpaid
		orderStatus = enums.OrderStatusCancelled
	case "expire":
		paymentStatus = enums.PaymentStatusUnpaid
		orderStatus = enums.OrderStatusCancelled
	default:
		// Jika ada status lain yang tidak kita tangani, kita tidak melakukan apa-apa
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Webhook received, no action taken"})
		return
	}

	updateReq.PaymentStatus = &paymentStatus
	updateReq.OrderStatus = &orderStatus

	// 6. Update order di database
	_, err = database.UpdateOrder(r.Context(), apiCfg.Firestore, payload.OrderID, updateReq)
	if err != nil {
		// Meskipun gagal update, kita tetap kirim 200 OK ke Midtrans
		// agar mereka tidak mengirim ulang notifikasi.
		// Kegagalan ini harus di-log secara serius untuk diperiksa manual.
		fmt.Printf("WEBHOOK_ERROR: Failed to update order %s: %v\n", payload.OrderID, err)
	}

	// 7. Kirim respons 200 OK ke Midtrans untuk mengonfirmasi penerimaan
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Webhook processed successfully"})
}

// generateSignatureKey membuat hash SHA-512 sesuai aturan Midtrans
func generateSignatureKey(orderID, statusCode, grossAmount, serverKey string) string {
	str := fmt.Sprintf("%s%s%s%s", orderID, statusCode, grossAmount, serverKey)
	hasher := sha512.New()
	hasher.Write([]byte(str))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
