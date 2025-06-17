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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not read request body")
		return
	}
	defer r.Body.Close()

	var payload MidtransNotificationPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid notification payload")
		return
	}

	signature := generateSignatureKey(payload.OrderID, payload.StatusCode, payload.GrossAmount, apiCfg.MidtransServerKey)
	if signature != payload.SignatureKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid signature")
		return
	}

	updateReq := database.UpdateOrderRequest{}
	var paymentStatus enums.PaymentStatus
	var orderStatus enums.OrderStatus

	switch payload.TransactionStatus {
	case "capture":
		if payload.FraudStatus == "accept" {
			paymentStatus = enums.PaymentStatusSuccess
			orderStatus = enums.OrderStatusConfirmed
		} else if payload.FraudStatus == "challenge" {
			paymentStatus = enums.PaymentStatusChallenge
		}
	case "settlement":
		paymentStatus = enums.PaymentStatusSuccess
		orderStatus = enums.OrderStatusConfirmed
	case "deny":
		paymentStatus = enums.PaymentStatusDeny
		orderStatus = enums.OrderStatusCancelled
	case "cancel", "expire":
		paymentStatus = enums.PaymentStatusFailure
		orderStatus = enums.OrderStatusCancelled
	case "pending":
		paymentStatus = enums.PaymentStatusPending
		orderStatus = enums.OrderStatusPending
	default:
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Webhook received, no action taken"})
		return
	}

	updateReq.PaymentStatus = &paymentStatus
	updateReq.OrderStatus = &orderStatus

	_, err = database.UpdateOrder(r.Context(), apiCfg.Firestore, payload.OrderID, updateReq)
	if err != nil {
		fmt.Printf("WEBHOOK_ERROR: Failed to update order %s: %v\n", payload.OrderID, err)
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Webhook processed successfully"})
}

func generateSignatureKey(orderID, statusCode, grossAmount, serverKey string) string {
	str := fmt.Sprintf("%s%s%s%s", orderID, statusCode, grossAmount, serverKey)
	hasher := sha512.New()
	hasher.Write([]byte(str))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
