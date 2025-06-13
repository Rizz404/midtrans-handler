package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func OrderRoutes(apiCfg *apiConfig) http.Handler {
	r := chi.NewRouter()

	r.Post("/", apiCfg.handlerCreateOrder)
	r.Get("/{orderID}", apiCfg.handlerGetOrderByID)
	r.Patch("/{orderID}", apiCfg.handlerUpdateOrder)

	return r
}
