package main

import (
	"os"

	"go-dev/docs"
	"go-dev/internal/database"
	"go-dev/internal/handlers"
	"go-dev/internal/middleware"
	"go-dev/internal/repository"
	"go-dev/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Subscription Management API
// @version 1.0
// @description REST API для управления подписками пользователей
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Настройка логгера
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)

	dbURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.WithFields(logrus.Fields{
		"port":         port,
		"database_url": dbURL,
	}).Info("Starting subscription service")

	// Подключение к базе
	db, err := database.Connect(dbURL)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	logger.Info("Connected to database successfully")

	// Миграции
	if err := database.RunMigrations(dbURL); err != nil {
		logger.WithError(err).Fatal("Failed to run migrations")
	}

	logger.Info("Database migrations completed")

	// Репозитории и сервисы
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, logger)

	// Роутер
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.Logger(logger))

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := router.Group("/api/v1")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
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

	logger.WithField("port", port).Info("Server starting")
	logger.Info("Swagger documentation available at: http://localhost:" + port + "/swagger/index.html")

	if err := router.Run(":" + port); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
