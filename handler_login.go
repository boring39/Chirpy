package main

import (
	"encoding/json"
	"net/http"

	"github.com/boring39/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// decoder error handling
		respondWithJSON(errorResponse{Error: "Something went wrong"}, http.StatusBadRequest, w)
		return
	}
	email := params.Email

	user, err := cfg.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		respondWithJSON(errorResponse{Error: "Something went wrong"}, http.StatusInternalServerError, w)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if !match || err != nil {
		respondWithJSON(errorResponse{Error: "Incorrect email or password"}, http.StatusUnauthorized, w)
		return
	}

	response := User{
		UserId:    user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(response, http.StatusOK, w)

}
