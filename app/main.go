package main

import (
	"github.com/TP2-Voice-Agora/backend/internal/lib/logger/prettyslog"
	"github.com/TP2-Voice-Agora/backend/internal/repository/postgres"
	"github.com/TP2-Voice-Agora/backend/internal/services/auth"
	"github.com/TP2-Voice-Agora/backend/internal/services/http-server"
	"github.com/TP2-Voice-Agora/backend/internal/services/ideas"
	"github.com/TP2-Voice-Agora/backend/internal/services/users"
	_ "github.com/joho/godotenv"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	// Load environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Logger
	logger := prettyslog.New()

	// Repository
	repo := &postgres.PostgresRepository{}
	err := repo.ConnectDB(dbURL, *logger)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer repo.CloseConnectDB()

	// Services
	ideaService := ideas.New(*logger, repo)
	authService := auth.New(*logger, repo, 2*time.Hour, jwtSecret)
	userService := users.New(*logger, repo)

	// HTTP Server
	server := http_server.NewHTTPServer(ideaService, authService, userService, logger)
	handler := server.SetupRoutes()

	logger.Info("Server starting...", slog.String("port", port))
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
