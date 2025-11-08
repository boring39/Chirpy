package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/boring39/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	reqParams := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqParams)
	if err != nil {
		// decoder error handling
		log.Printf("Failed to decode request body: %s\n", err.Error())
		respondWithJSON(errorResponse{Error: "Something went wrong"}, http.StatusBadRequest, w)
		return
	}

	cleaned, err := validateSanitizeChirp(reqParams.Body)
	if err != nil {
		respondWithJSON(errorResponse{Error: err.Error()}, http.StatusBadRequest, w)
		return
	}

	// prepare valid response
	chirpParams := database.CreateChirpParams{
		Body:   cleaned,
		Userid: reqParams.UserId,
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		log.Printf("Failed to create chirp in database: %s\n", err.Error())
		log.Printf("Failed to create chirp in database: %+v\n", chirpParams)
		respondWithJSON(errorResponse{Error: "Something went wrong"}, http.StatusInternalServerError, w)
		return
	}

	response := Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.Userid,
	}
	respondWithJSON(response, http.StatusCreated, w)
}

func validateSanitizeChirp(str string) (string, error) {
	const maxChirpLength = 140
	if len(str) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	return sanitizeInput(str), nil
}
func sanitizeInput(str string) string {
	badWords := [3]string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(str, " ")
	for i, word := range words {
		match := false
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				match = true
				break
			}
		}
		if match {
			words[i] = "****"
			continue
		}
	}
	return strings.Join(words, " ")
}
