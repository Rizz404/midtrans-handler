package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func paymentMethodRoutes(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()

	r.Post("/", apiCfg.handlerCreatePaymentMethod)
	r.Get("/", apiCfg.handlerGetAllPaymentMethods)
	r.Post("/bulk", apiCfg.handlerBulkCreatePaymentMethods)
	r.Get("/{paymentMethodID}", apiCfg.handlerGetPaymentMethodByID)
	r.Put("/{paymentMethodID}", apiCfg.handlerUpdatePaymentMethod)
	r.Delete("/{paymentMethodID}", apiCfg.handlerDeletePaymentMethod)

	return r
}
