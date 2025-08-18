package main

import (
	"go-dev/internal/config"
	"go-dev/internal/database"
	"go-dev/internal/handlers"
	"go-dev/internal/repository"
	"go-dev/internal/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Настройка логгера
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	if cfg.Environment == "development" {
		logger.SetLevel(logrus.DebugLevel)
	}

	// Подключение к базе данных
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Запуск миграций
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		logger.Fatal("Failed to run migrations:", err)
	}

	// Инициализация слоев приложения
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, logger)

	// Настройка роутера
	router := gin.Default()

	// Простейшие проверочные endpoints
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// API маршруты
	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "subscription-service",
			})
		})

		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("", subscriptionHandler.Create)
			subscriptions.GET("", subscriptionHandler.List)
			subscriptions.GET("/:id", subscriptionHandler.GetByID)
			subscriptions.PUT("/:id", subscriptionHandler.Update)
			subscriptions.DELETE("/:id", subscriptionHandler.Delete)
			subscriptions.GET("/total-cost", subscriptionHandler.GetTotalCost)
		}
	}

	logger.Infof("Server starting on port %s", cfg.Port)
	logger.Info("Available endpoints:")
	logger.Info("  GET    /ping")
	logger.Info("  GET    /api/v1/health")
	logger.Info("  POST   /api/v1/subscriptions")
	logger.Info("  GET    /api/v1/subscriptions")
	logger.Info("  GET    /api/v1/subscriptions/:id")
	logger.Info("  PUT    /api/v1/subscriptions/:id")
	logger.Info("  DELETE /api/v1/subscriptions/:id")
	logger.Info("  GET    /api/v1/subscriptions/total-cost")

	if err := router.Run(":" + cfg.Port); err != nil {
		logger.Fatal("Failed to start server:", err)
	}
}
