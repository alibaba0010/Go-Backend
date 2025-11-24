package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	redisPkg "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/argon2"

	"github.com/alibaba0010/postgres-api/internal/config"
	"github.com/alibaba0010/postgres-api/internal/database"
	"github.com/alibaba0010/postgres-api/internal/dto"
	"github.com/alibaba0010/postgres-api/internal/errors"
	"github.com/alibaba0010/postgres-api/internal/logger"
	"github.com/alibaba0010/postgres-api/internal/models"
	"github.com/alibaba0010/postgres-api/internal/utils"
)

// RegisterUser handles the DB work for signing up a new user.
// It checks for an existing email, hashes the password and inserts the user.
// Returns the created user (with ID populated) or an AppError for controller to return.
func RegisterUser(ctx context.Context, input dto.SignupInput) (*models.User, *errors.AppError) {
	// Validate input using same validation rules as controllers previously used
	validate := validator.New()
	dto.RegisterValidators(validate)

	// Run validation and convert errors to friendly messages
	if err := validate.Struct(input); err != nil {
		if ves, ok := err.(validator.ValidationErrors); ok {
			logger.Log.Info("validation errors during registration", zap.Error(err))
			var messages []string
			for _, fe := range ves {
				var msg string
				field := fe.Field()
				switch fe.Tag() {
				case "required":
					msg = fmt.Sprintf("%s is required", field)
				case "min":
					msg = fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
				case "max":
					msg = fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
				case "email":
					msg = fmt.Sprintf("%s must be a valid email address", field)
				case "password_special":
					msg = "password must contain at least one uppercase letter, one lowercase letter, one digit, and one special character"
				case "eqfield":
					// fe.Param() holds the field the current field must equal (e.g., Password)
					msg = fmt.Sprintf("%s must match %s", field, fe.Param())
				default:
					msg = fmt.Sprintf("%s is invalid", field)
				}
				messages = append(messages, msg)
			}
			return nil, errors.ValidationErrors(messages)
		}
		// Non-validation error
		return nil, errors.ValidationError(err.Error())
	}

	logger.Log.Info("Signup request for email: " + input.Name)

	// Check if user already exists
	exists, err := database.DB.NewSelect().Model((*models.User)(nil)).
		Where("email = ?", input.Email).
		Exists(ctx)
	if err != nil {
		return nil, errors.InternalError(err)
	}
	if exists {
		return nil, errors.DuplicateError("email")
	}

	// We don't hash the password here; it will be hashed with argon2id at activation.
	// Store signup payload (including raw password) temporarily in Redis instead.
	user := &models.User{
		Name:  input.Name,
		Email: input.Email,
		// Password will be set on activation
	}

	// Ensure ID is set to a UUIDv7
	newUUID, err := utils.GenerateUUIDv7()
	if err != nil {
		return nil, errors.InternalError(err)
	}
	user.ID = newUUID.String()

	// Generate verification token
	token, err := utils.GenerateToken()
	if err != nil {
		return nil, errors.InternalError(err)
	}

	// Store the signup payload (including raw password) in Redis with TTL 15 minutes.
	// NOTE: storing raw password temporarily is a security risk in production —
	// consider hashing with argon2id before storing or encrypting the payload.
	payload := struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: input.Password,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.InternalError(err)
	}

	key := "verify:" + token
	ttl := 15 * time.Minute
	if err := database.RedisClient.Set(ctx, key, b, ttl).Err(); err != nil {
		return nil, errors.InternalError(err)
	}

	// Build verification URL — Auth routes are mounted under /api/v1/auth, so
	// the full verify path is /api/v1/auth/verify
	cfg := config.LoadConfig()
	verifyURL := fmt.Sprintf("http://localhost:%s/api/v1/auth/verify?token=%s", cfg.Port, token)
	html := VerifyMail(user.Name, verifyURL)
	if err := SendEmail(user.Email, "Verify your email", html); err != nil {
		logger.Log.Error("failed to send verification email", zap.Error(err))
		// Attempt to delete token to avoid leaking
		_ = database.RedisClient.Del(ctx, key).Err()
		return nil, errors.InternalError(err)
	}

	// Per new flow, registration doesn't persist the user yet — activation will.
	return nil, nil
}
func ActivateUser(ctx context.Context, token string) (*models.User, *errors.AppError) {
	key := "verify:" + token
	data, err := database.RedisClient.Get(ctx, key).Bytes()
	if err == redisPkg.Nil {
		return nil, errors.ValidationError("invalid or expired token")
	}
	if err != nil {
		return nil, errors.InternalError(err)
	}

	var payload struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		_ = database.RedisClient.Del(ctx, key).Err()
		return nil, errors.InternalError(err)
	}

	// Hash password with argon2id
	hashedPwd, err := hashPassword(payload.Password)
	if err != nil {
		return nil, errors.InternalError(err)
	}

	user := &models.User{
		ID:       payload.ID,
		Name:     payload.Name,
		Email:    payload.Email,
		Password: hashedPwd,
	}

	// Insert into DB
	_, err = database.DB.NewInsert().Model(user).
		Returning("id").
		Exec(ctx)
	if err != nil {
		return nil, errors.InternalError(err)
	}

	// Token used within TTL -> remove it
	if err := database.RedisClient.Del(ctx, key).Err(); err != nil {
		logger.Log.Error("failed to delete verification token", zap.Error(err))
	}

	return user, nil
}

// hashPassword hashes the provided password using argon2id and returns an encoded
// string containing the salt and hash.
func hashPassword(password string) (string, error) {
	// Parameters
	var (
		timeParam uint32 = 1
		memory    uint32 = 64 * 1024
		threads   uint8  = 4
		keyLen    uint32 = 32
		saltLen   uint32 = 16
	)

	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, timeParam, memory, threads, keyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memory, timeParam, threads, b64Salt, b64Hash)
	return encoded, nil
}