package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alibaba0010/postgres-api/internal/dto"
	"github.com/alibaba0010/postgres-api/internal/errors"
	"github.com/alibaba0010/postgres-api/internal/services"
)

// SignupHandler godoc
//
//	@Summary		User Signup
//	@Description	Creates a new user account
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dto.SignupInput	true	"Signup request"
//	@Success		201	{object}	map[string]interface{} "User created successfully"
//	@Failure		400	{object}	map[string]string		"Validation error"
//	@Failure		409	{object}	map[string]string		"Duplicate email"
//	@Failure		500	{object}	map[string]string		"Internal server error"
//	@Router			/auth/signup [post]
func SignupHandler(writer http.ResponseWriter, request *http.Request) {
	var input dto.SignupInput

	// Decode JSON
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		errors.ErrorResponse(writer, request, errors.ValidationError("Invalid JSON body"))
		return
	}
	
	_, appErr := services.RegisterUser(request.Context(), input)
	if appErr != nil {
		errors.ErrorResponse(writer, request, appErr)
		return
	}
	// Per new flow we don't persist the user at signup; activation will.
	resp := map[string]string{
		"title":   "Successfully signed up",
		"message": "Please check your email for a verification link",
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(writer).Encode(resp)
}
func ActivateUserHandler(writer http.ResponseWriter, request *http.Request) {
	token := request.URL.Query().Get("token")
	if token == "" {
		errors.ErrorResponse(writer, request, errors.ValidationError("token is required"))
		return
	}

	user, appErr := services.ActivateUser(request.Context(), token)
	if appErr != nil {
		errors.ErrorResponse(writer, request, appErr)
		return
	}

	// Return created user (omit password)
	resp := dto.SignUpResponse{
		Title: "User activated successfully",
		Data: dto.SignUpData{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(resp)
}