package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"xm-exercise/internal/auth"
	"xm-exercise/internal/db"
	"xm-exercise/internal/logger"
	"xm-exercise/pkg/models"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	userRepo   *db.UserRepository
	jwtService *auth.JWTService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo *db.UserRepository, jwtService *auth.JWTService) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserRegistration true "User registration data"
// @Success 200 {object} TokenResponse "User registered successfully"
// @Failure 400 {string} string "Invalid request body or validation error"
// @Failure 409 {string} string "Name or email already taken"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.WithContext(ctx)
	var creds models.UserRegistration
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := creds.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := h.userRepo.ExistsByEmail(creds.Email)
	if err != nil {
		http.Error(w, "Error checking email", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	user := models.User{
		ID:           uuid.New().String(),
		Name:         creds.Name,
		Email:        creds.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.userRepo.Create(user); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Generate token
	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(models.TokenResponse{
		Token: token,
	}); err != nil {
		log.Error("Failed to encode response data",
			zap.Error(err),
		)
	}
}

// Login godoc
// @Summary Login a user
// @Description Login with username and password and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserLogin true "Login credentials"
// @Success 200 {object} TokenResponse "User logged in successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.WithContext(ctx)
	var creds models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByEmail(creds.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.jwtService.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(models.TokenResponse{
		Token: token,
	}); err != nil {
		log.Error("Failed to encode response data",
			zap.Error(err),
		)
	}
}
