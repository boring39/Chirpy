package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Failed to create chirp in database: %s\n", err.Error())
		respondWithJSON(errorResponse{Error: "Something went wrong"}, http.StatusInternalServerError, w)
		return
	}

	response := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		response[i] = Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.Userid,
		}
	}
	respondWithJSON(response, http.StatusOK, w)
}
