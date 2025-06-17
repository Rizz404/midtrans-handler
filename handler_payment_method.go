package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rizz404/midtrans-handler/internal/database"
	"github.com/Rizz404/midtrans-handler/internal/enums"
)

func (apiCfg *apiConfig) handlerCreatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name                      string                  `json:"name"`
		Description               string                  `json:"description"`
		Logo                      *string                 `json:"logo,omitempty"`
		PaymentMethodType         enums.PaymentMethodType `json:"paymentMethodType"`
		MidtransIdentifier        *string                 `json:"midtransIdentifier"`
		MinimumAmount             float64                 `json:"minimumAmount"`
		MaximumAmount             float64                 `json:"maximumAmount"`
		AdminPaymentCode          *string                 `json:"adminPaymentCode,omitempty"`
		AdminPaymentQrCodePicture *string                 `json:"adminPaymentQrCodePicture,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error passing json: %v", err))
		return
	}

	paymentMethod, err := database.CreatePaymentMethod(r.Context(), apiCfg.Firestore, database.CreatePaymentMethodRequest{
		Name:                      params.Name,
		Description:               params.Description,
		Logo:                      params.Logo,
		PaymentMethodType:         params.PaymentMethodType,
		MidtransIdentifier:        params.MidtransIdentifier,
		MinimumAmount:             params.MinimumAmount,
		MaximumAmount:             params.MaximumAmount,
		AdminPaymentCode:          params.AdminPaymentCode,
		AdminPaymentQrCodePicture: params.AdminPaymentQrCodePicture,
	})

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn't create paymentMethod: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, dbPaymentMethodToPaymentMethod(*paymentMethod))
}

func (apiCfg *apiConfig) handlerBulkCreatePaymentMethods(w http.ResponseWriter, r *http.Request) {
	var params []database.CreatePaymentMethodRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON array: %v", err))
		return
	}

	if len(params) == 0 {
		respondWithError(w, http.StatusBadRequest, "Request body must contain at least one payment method.")
		return
	}

	paymentMethods, err := database.BulkCreatePaymentMethods(r.Context(), apiCfg.Firestore, params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't bulk create payment methods: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, dbPaymentMethodsToPaymentMethods(paymentMethods))
}

func (apiCfg *apiConfig) handlerGetAllPaymentMethods(w http.ResponseWriter, r *http.Request) {
	paymentMethods, err := database.GetAllPaymentMethods(r.Context(), apiCfg.Firestore)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get payment methods: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbPaymentMethodsToPaymentMethods(paymentMethods))
}

func (apiCfg *apiConfig) handlerGetPaymentMethodByID(w http.ResponseWriter, r *http.Request) {
	paymentMethodID := r.PathValue("paymentMethodID")

	paymentMethod, err := database.GetPaymentMethodByID(r.Context(), apiCfg.Firestore, paymentMethodID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't get payment method with ID %s: %v", paymentMethodID, err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbPaymentMethodToPaymentMethod(*paymentMethod))
}

func (apiCfg *apiConfig) handlerUpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	paymentMethodID := r.PathValue("paymentMethodID")

	decoder := json.NewDecoder(r.Body)
	params := database.UpdatePaymentMethodRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	updatedPaymentMethod, err := database.UpdatePaymentMethod(r.Context(), apiCfg.Firestore, paymentMethodID, params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't update payment method %s: %v", paymentMethodID, err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbPaymentMethodToPaymentMethod(*updatedPaymentMethod))
}

func (apiCfg *apiConfig) handlerDeletePaymentMethod(w http.ResponseWriter, r *http.Request) {
	paymentMethodID := r.PathValue("paymentMethodID")

	err := database.DeletePaymentMethod(r.Context(), apiCfg.Firestore, paymentMethodID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't delete payment method %s: %v", paymentMethodID, err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
