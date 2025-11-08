package main

import (
	"net/http"
	"os"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	var statusCode int
	var response any

	if os.Getenv("PLATFORM") != "dev" {
		response := errorResponse{Error: "Something went wrong"}
		statusCode = http.StatusForbidden
		respondWithJSON(response, statusCode, w)
		return
	}

	// decoder := json.NewDecoder(r.Body)
	// params := parameters{}
	// err := decoder.Decode(&params)
	// if err != nil {
	// 	// decoder error handling
	// 	response := errorResponse{Error: "Something went wrong"}
	// 	statusCode = http.StatusBadRequest
	// 	respondWithJSON(response, statusCode, w)
	// 	return
	// }
	cfg.fileserverHits.Store(0)
	err := cfg.db.Reset(r.Context())
	if err != nil {
		response := errorResponse{Error: "Internal Database Error"}
		statusCode = http.StatusInternalServerError
		respondWithJSON(response, statusCode, w)
		// 	return
	}
	statusCode = http.StatusOK
	response = errorResponse{Error: "success"}
	respondWithJSON(response, statusCode, w)
}
