package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/maxviazov/dolina-flower-order-backend/internal/config"
	"github.com/maxviazov/dolina-flower-order-backend/internal/handlers"
	"github.com/maxviazov/dolina-flower-order-backend/internal/logger"
	"github.com/maxviazov/dolina-flower-order-backend/internal/repository/sqlite"
	"github.com/maxviazov/dolina-flower-order-backend/internal/services"
)

type App struct {
	config *config.Config
	logger *logger.Logger
	server *http.Server
	router *gin.Engine
	repo   *sqlite.Repository
}

func New() *App {
	return &App{
		logger: logger.GetLogger(),
	}
}

func (a *App) Initialize() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	a.config = cfg

	if err := a.logger.Initialize(cfg); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	a.logger.Info("Application initializing...")

	repo, err := sqlite.NewRepository("./flowers.db")
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	a.repo = repo

	if a.config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	a.router = gin.New()
	a.setupMiddleware()
	a.setupRoutes()

	a.server = &http.Server{
		Addr:    a.config.GetServerAddress(),
		Handler: a.router,
	}

	a.logger.Info("Application initialized successfully")
	return nil
}

func (a *App) setupMiddleware() {
	a.router.Use(gin.Logger())
	a.router.Use(gin.Recovery())
	a.router.Use(a.corsMiddleware())
}

func (a *App) setupRoutes() {
	a.router.GET("/health", a.healthCheck)

	orderService := services.NewOrderService(a.repo)
	orderHandler := handlers.NewOrderHandler(orderService)
	flowerHandler := handlers.NewFlowerHandler(orderService)

	api := a.router.Group("/api/v1")
	{
		api.GET("/ping", a.ping)
		api.GET("/flowers", flowerHandler.GetAvailableFlowers)

		orders := api.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
		}
	}
}

func (a *App) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// healthCheck обработчик проверки здоровья приложения
func (a *App) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "dolina-flower-order-backend",
	})
}

// ping простой обработчик для тестирования
func (a *App) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (a *App) Run(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	a.logger.WithField("address", a.server.Addr).Info("Starting server")

	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	a.logger.Info("Server started successfully")

	select {
	case <-quit:
		a.logger.Info("Shutdown signal received")
	case <-ctx.Done():
		a.logger.Info("Context cancelled")
	}

	return a.Shutdown()
}

func (a *App) Shutdown() error {
	a.logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if a.repo != nil {
		a.repo.Close()
	}

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error("Server forced to shutdown")
		return err
	}

	a.logger.Info("Server shutdown completed")
	return nil
}
