package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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
