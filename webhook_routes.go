package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func webhookRoutes(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()

	r.Post("/midtrans", apiCfg.handlerMidtransWebhook)

	return r
}
