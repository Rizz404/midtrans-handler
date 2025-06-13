package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Rizz404/midtrans-handler/internal/database"
	"github.com/Rizz404/midtrans-handler/internal/enums"
)

func (apiCfg *apiConfig) handlerCreateOrder(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserID              string               `json:"userId"`
		PaymentMethodID     string               `json:"paymentMethodId"`
		OrderType           enums.OrderType      `json:"orderType"`
		EstimatedReadyTime  *time.Time           `json:"estimatedReadyTime,omitempty"`
		SpecialInstructions *string              `json:"specialInstructions,omitempty"`
		OrderItems          []database.OrderItem `json:"orderItems"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error passing json: %v", err))
		return
	}

	finalOrder, err := database.CreateOrderWithPayment(
		r.Context(),
		apiCfg.Firestore,
		apiCfg.MidtransCore,
		database.CreateOrderWithPaymentRequest{
			UserID:              params.UserID,
			PaymentMethodID:     params.PaymentMethodID,
			OrderType:           params.OrderType,
			EstimatedReadyTime:  params.EstimatedReadyTime,
			SpecialInstructions: params.SpecialInstructions,
			OrderItems:          params.OrderItems,
		},
	)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, dbOrderToOrder(*finalOrder))
}

func (apiCfg *apiConfig) handlerGetOrderByID(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("orderID")

	order, err := database.GetOrderByID(r.Context(), apiCfg.Firestore, orderID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get order with ID %s: %v", orderID, err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbOrderToOrder(*order))
}

func (apiCfg *apiConfig) handlerUpdateOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("orderID")

	decoder := json.NewDecoder(r.Body)
	params := database.UpdateOrderRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	updatedOrder, err := database.UpdateOrder(r.Context(), apiCfg.Firestore, orderID, params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't update order %s: %v", orderID, err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbOrderToOrder(*updatedOrder))
}
