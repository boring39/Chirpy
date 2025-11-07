package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type validResponse struct {
		Valid bool `json:"valid"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	var statusCode int
	var response any
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// decoder error handling
		response := errorResponse{Error: "Something went wrong"}
		statusCode = http.StatusBadRequest
		respondWithJSON(response, http.StatusBadRequest, w)
		return
	}

	if len(params.Body) > 140 {
		//prepare error response
		response = errorResponse{Error: "Chirp is too long"}
		statusCode = http.StatusBadRequest
	} else {
		// prepare valid response
		response = validResponse{Valid: true}
		statusCode = http.StatusOK
	}

	respondWithJSON(response, statusCode, w)
}

func respondWithJSON(response any, statusCode int, w http.ResponseWriter) {
	responseBody, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		statusCode = http.StatusInternalServerError
		responseBody = []byte(`{"error":"Something went wrong"}`)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(responseBody)
}
