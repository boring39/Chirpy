package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	pathValue := r.PathValue("chirpId")
	id, err := uuid.Parse(pathValue)
	if err != nil {
		log.Printf("Failed to parse UUID (%s): %s\n", pathValue, err.Error())
		respondWithJSON(errorResponse{Error: "Something went wrong"}, http.StatusBadRequest, w)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		log.Printf("Chirp UUID (%s) not found: %s", id, err.Error())
		respondWithJSON(errorResponse{Error: "Chirp Id not found"}, http.StatusNotFound, w)
		return
	}

	response := Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.Userid,
	}
	respondWithJSON(response, http.StatusOK, w)
}

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
