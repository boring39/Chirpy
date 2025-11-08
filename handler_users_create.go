package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
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
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		response := errorResponse{Error: "Something went wrong"}
		statusCode = http.StatusInternalServerError
		respondWithJSON(response, statusCode, w)
		return
	}

	response = User{
		UserId:    user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	statusCode = http.StatusCreated
	respondWithJSON(response, statusCode, w)

}
