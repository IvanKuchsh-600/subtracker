package main

import (
	"log"

	"github.com/IvanKuchsh-600/subtracker/internal/app"
	"github.com/IvanKuchsh-600/subtracker/internal/config"
)

func main() {
	cfg := config.MustLoad()

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	app.Run()
}

// package main

// import (
// 	"context"
// 	"log/slog"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/jackc/pgx/v5/stdlib"
// 	ginSwagger "github.com/swaggo/gin-swagger"

// 	"subtracker/internal/config"
// 	"subtracker/internal/handlers"
// 	"subtracker/internal/repository"
// 	"subtracker/internal/service"
// )

// // @title SubTracker API
// // @version 1.0
// // @description Subscription tracking service API
// // @host localhost:8080
// // @BasePath /api/v1
// func main() {
// 	// Load configuration
// 	cfg, err := config.Load()
// 	if err != nil {
// 		slog.Error("Failed to load config", "error", err)
// 		os.Exit(1)
// 	}

// 	// Setup logger
// 	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
// 		Level: slog.LevelDebug,
// 	}))

// 	// Connect to database
// 	poolConfig, err := pgxpool.ParseConfig(cfg.Database.ConnString())
// 	if err != nil {
// 		logger.Error("Failed to parse database config", "error", err)
// 		os.Exit(1)
// 	}

// 	poolConfig.MaxConns = cfg.Database.MaxConnections
// 	poolConfig.MinConns = cfg.Database.MinConnections

// 	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
// 	if err != nil {
// 		logger.Error("Failed to connect to database", "error", err)
// 		os.Exit(1)
// 	}
// 	defer pool.Close()

// 	// Test database connection
// 	if err := pool.Ping(context.Background()); err != nil {
// 		logger.Error("Database ping failed", "error", err)
// 		os.Exit(1)
// 	}

// 	logger.Info("Connected to database", "host", cfg.Database.Host, "dbname", cfg.Database.DBName)

// 	// Run migrations
// 	db := stdlib.OpenDB(*pool.Config().ConnConfig)
// 	defer db.Close()

// 	migrationSQL, err := os.ReadFile("migrations/001_create_subscriptions_table.sql")
// 	if err != nil {
// 		logger.Warn("Migration file not found, skipping", "error", err)
// 	} else {
// 		if _, err := db.Exec(string(migrationSQL)); err != nil {
// 			logger.Warn("Migration failed", "error", err)
// 		} else {
// 			logger.Info("Migrations applied successfully")
// 		}
// 	}

// 	// Setup dependencies
// 	subscriptionRepo := repository.NewSubscriptionRepository(pool)
// 	subscriptionService := service.NewSubscriptionService(subscriptionRepo, logger)
// 	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)

// 	// Setup Gin
// 	gin.SetMode(gin.ReleaseMode)
// 	router := gin.New()
// 	router.Use(gin.Recovery())
// 	router.Use(gin.Logger())

// 	// API routes
// 	api := router.Group("/api/v1")
// 	{
// 		subscriptions := api.Group("/subscriptions")
// 		{
// 			subscriptions.POST("", subscriptionHandler.Create)
// 			subscriptions.GET("", subscriptionHandler.List)
// 			subscriptions.GET("/total-cost", subscriptionHandler.GetTotalCost)
// 			subscriptions.GET("/:id", subscriptionHandler.GetByID)
// 			subscriptions.PUT("/:id", subscriptionHandler.Update)
// 			subscriptions.DELETE("/:id", subscriptionHandler.Delete)
// 		}
// 	}

// 	// Swagger
// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	// Health check
// 	router.GET("/health", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	})

// 	// Start server
// 	srv := &http.Server{
// 		Addr:         ":" + cfg.Server.Port,
// 		Handler:      router,
// 		ReadTimeout:  15 * time.Second,
// 		WriteTimeout: 15 * time.Second,
// 		IdleTimeout:  60 * time.Second,
// 	}

// 	go func() {
// 		logger.Info("Starting server", "port", cfg.Server.Port)
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.Error("Failed to start server", "error", err)
// 			os.Exit(1)
// 		}
// 	}()

// 	// Graceful shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	logger.Info("Shutting down server...")

// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	if err := srv.Shutdown(ctx); err != nil {
// 		logger.Error("Server forced to shutdown", "error", err)
// 	}

// 	logger.Info("Server exited")
// }
