package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"dolina-flower-order-backend/internal/config"
	"dolina-flower-order-backend/internal/logger"
)

// App представляет основное приложение
type App struct {
	config *config.Config
	logger *logger.Logger
	server *http.Server
	router *gin.Engine
}

// New создает новый экземпляр приложения
func New() *App {
	return &App{
		logger: logger.GetLogger(),
	}
}

// Initialize инициализирует приложение
func (a *App) Initialize() error {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	a.config = cfg

	// Инициализируем логгер
	if err := a.logger.Initialize(cfg); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	a.logger.Info("Application initializing...")

	// Настраиваем Gin
	if a.config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Создаем роутер
	a.router = gin.New()
	a.setupMiddleware()
	a.setupRoutes()

	// Создаем HTTP сервер
	a.server = &http.Server{
		Addr:         a.config.GetServerAddress(),
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
		IdleTimeout:  a.config.Server.IdleTimeout,
	}

	a.logger.Infof("Application initialized successfully on %s", a.config.GetServerAddress())
	return nil
}

// setupMiddleware настраивает middleware
func (a *App) setupMiddleware() {
	// Логирование запросов
	a.router.Use(a.loggingMiddleware())

	// Recovery middleware
	a.router.Use(gin.Recovery())

	// CORS middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = a.config.Security.CORSOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	a.router.Use(cors.New(corsConfig))
}

// setupRoutes настраивает маршруты
func (a *App) setupRoutes() {
	// Health check
	a.router.GET("/health", a.healthCheck)

	// API группа
	api := a.router.Group("/api/v1")
	{
		// Пока базовые маршруты для MVP
		api.GET("/ping", a.ping)
	}
}

// loggingMiddleware создает middleware для логирования запросов
func (a *App) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Обрабатываем запрос
		c.Next()

		// Логируем после обработки
		duration := time.Since(start)
		a.logger.LogRequest(
			c.Request.Context(),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
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

// Run запускает приложение
func (a *App) Run(ctx context.Context) error {
	// Канал для получения сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		a.logger.Infof("Starting server on %s", a.config.GetServerAddress())
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	a.logger.Info("Server started successfully")

	// Ждем сигнал завершения
	select {
	case <-quit:
		a.logger.Info("Shutdown signal received")
	case <-ctx.Done():
		a.logger.Info("Context cancelled")
	}

	return a.Shutdown()
}

// Shutdown корректно завершает работу приложения
func (a *App) Shutdown() error {
	a.logger.Info("Shutting down server...")

	// Создаем контекст с таймаутом для завершения
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Завершаем HTTP сервер
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	a.logger.Info("Server shutdown completed")
	return nil
}
