package middleware

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"xm-exercise/internal/auth"
	"xm-exercise/internal/logger"
)

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	jwtService *auth.JWTService
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService}
}

// Authenticate middleware validates JWT tokens and adds user ID to context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.WithContext(r.Context())

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Warn("Missing authorization header")
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Warn("Invalid authorization format")
			http.Error(w, "Invalid authorization format, expected 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Warn("Invalid JWT token", zap.Error(err))
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		log.Info("User authenticated", zap.String("user_id", claims.UserID))
		ctx := context.WithValue(r.Context(), logger.UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) (string, bool) {
	// Fixme: Cast the context.Value result to interface{} before type assertion
	value := ctx.Value(logger.UserIDKey)
	if value == nil {
		return "", false
	}
	userID, ok := value.(string)
	return userID, ok
}
