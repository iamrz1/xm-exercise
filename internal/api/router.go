package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	"xm-exercise/internal/api/handlers"
	appMiddleware "xm-exercise/internal/api/middleware"
	"xm-exercise/internal/auth"
	"xm-exercise/internal/config"
	"xm-exercise/internal/db"
	_ "xm-exercise/internal/docs"
	"xm-exercise/internal/events"
	"xm-exercise/internal/logger"
)

// NewRouter creates a new router with all application routes
func NewRouter(database *db.Database, producer *events.KafkaProducer, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpiration)

	companyRepo := db.NewCompanyRepository(database)
	userRepo := db.NewUserRepository(database)

	authHandler := handlers.NewAuthHandler(userRepo, jwtService)
	companyHandler := handlers.NewCompanyHandler(companyRepo, producer)

	authMiddleware := appMiddleware.NewAuthMiddleware(jwtService)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(appMiddleware.LoggerMiddleware)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)

		cr := chi.NewRouter()
		cr.Get("/{id}", companyHandler.Get)
		cr.With(authMiddleware.Authenticate).Post("/", companyHandler.Create)
		cr.With(authMiddleware.Authenticate).Patch("/{id}", companyHandler.Patch)
		cr.With(authMiddleware.Authenticate).Delete("/{id}", companyHandler.Delete)
		r.Mount("/companies", cr)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
	))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// nolint:errcheck
		// reason: The response is simple and a failure is unlikely to be recoverable.
		w.Write([]byte("OK"))
	})

	logger.Info("API routes initialized")
	return r
}
