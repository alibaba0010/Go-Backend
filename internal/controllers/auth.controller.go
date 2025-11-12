package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/alibaba0010/postgres-api/internal/database"
	"github.com/alibaba0010/postgres-api/internal/dto"
	"github.com/alibaba0010/postgres-api/internal/errors"
	"github.com/alibaba0010/postgres-api/internal/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

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

	// Validate input
	if err := validate.Struct(input); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			msg := e.Field() + " failed validation: " + e.Tag()
			errors.ErrorResponse(writer, request, errors.ValidationError(msg))
			return
		}
	}

	// Check if user already exists
	exists, err := database.DB.NewSelect().Model((*models.User)(nil)).
		Where("email = ?", input.Email).
		Exists(request.Context())
	if err != nil {
		errors.ErrorResponse(writer, request, errors.InternalError(err))
		return
	}
	if exists {
		errors.ErrorResponse(writer, request, errors.DuplicateError("email"))
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		errors.ErrorResponse(writer, request, errors.InternalError(err))
		return
	}

	// Save new user
	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
	// Insert only specific columns so the DB can generate defaults (id, created_at, etc.)
	// Request the DB to RETURNING the generated id so bun will populate user.ID.
	_, err = database.DB.NewInsert().Model(user).
		Column("name", "email", "password").
		Returning("id").
		Exec(request.Context())
	if err != nil {
		errors.ErrorResponse(writer, request, errors.InternalError(err))
		return
	}

	// Return success response (without password)
	response := map[string]interface{}{
		
		"title": "User created successfully",
		"data": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(response)
}
