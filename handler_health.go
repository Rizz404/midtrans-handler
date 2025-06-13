package main

import (
	"net/http"
)

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	respondWithJSON(w, http.StatusOK, response)
}
