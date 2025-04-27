package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamrz1/xm-exercise/internal/api/handlers"
	"github.com/iamrz1/xm-exercise/internal/api/middleware"
	"github.com/iamrz1/xm-exercise/internal/auth"
	"github.com/iamrz1/xm-exercise/internal/config"
	"github.com/iamrz1/xm-exercise/internal/db"
	"github.com/iamrz1/xm-exercise/internal/events"
)

// NewRouter creates a new router with all application routes
func NewRouter(database *db.PostgresDB, producer *events.KafkaProducer, cfg *config.Config) http.Handler {
	r := mux.NewRouter()

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiration)

	// Initialize repositories
	companyRepo := db.NewCompanyRepository(database)
	userRepo := db.NewUserRepository(database)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	companyHandler := handlers.NewCompanyHandler(companyRepo, producer)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Auth routes
	r.HandleFunc("/api/v1/auth/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods(http.MethodPost)

	// Company routes
	companies := r.PathPrefix("/api/v1/companies").Subrouter()
	companies.Use(authMiddleware.Authenticate)

	companies.HandleFunc("", companyHandler.Create).Methods(http.MethodPost)
	companies.HandleFunc("/{id}", companyHandler.Get).Methods(http.MethodGet)
	companies.HandleFunc("/{id}", companyHandler.Patch).Methods(http.MethodPatch)
	companies.HandleFunc("/{id}", companyHandler.Delete).Methods(http.MethodDelete)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	return r
}
