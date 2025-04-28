// Package main provides the entry point for the API server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"xm-exercise/internal/api"
	"xm-exercise/internal/config"
	"xm-exercise/internal/db"
	"xm-exercise/internal/events"
	"xm-exercise/internal/logger"
	"xm-exercise/internal/utils"
)

// @title           			Company Management API
// @version         			1.0
// @description     			A company management service API
// @termsOfService  			http://swagger.io/terms/
// @contact.name    			Rezoan Tamal
// @contact.email 				my.name.in.lower.case@gmail.com
// @license.name  				MIT
// @license.url   				https://opensource.org/licenses/MIT
// @host      					localhost:8080
// @BasePath  					/api/v1
// @securityDefinitions.apikey  Bearer
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and the JWT token.

const (
	EnvLogLevel = "LOG_LEVEL"
	EnvAppEnv   = "APP_ENV"
	DevEnv      = "dev"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	isDev := strings.ToLower(utils.GetEnv(EnvAppEnv, DevEnv)) == DevEnv
	logLevel := os.Getenv(EnvLogLevel)
	if logLevel == "" {
		logLevel = "info"
	}

	if err := logger.Init(logLevel, isDev); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	database, err := db.NewDatabase(cfg.DatabaseDialect, cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer database.Close()
	logger.Info("Database connected", zap.String("dialect", cfg.DatabaseDialect))

	producer := events.NewKafkaProducer(cfg.KafkaBrokers)
	defer producer.Close()
	logger.Info("Kafka producer initialized", zap.Strings("brokers", cfg.KafkaBrokers))

	router := api.NewRouter(database, producer, cfg)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
