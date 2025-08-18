package handlers

import (
	"go-dev/internal/models"
	"go-dev/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
	logger  *logrus.Logger
}

func NewSubscriptionHandler(service *service.SubscriptionService, logger *logrus.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

// Create создает новую подписку
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req models.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"service_name": req.ServiceName,
		"user_id":      req.UserID,
		"price":        req.Price,
	}).Info("Creating subscription")

	subscription, err := h.service.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	h.logger.WithField("subscription_id", subscription.ID).Info("Subscription created successfully")
	c.JSON(http.StatusCreated, subscription)
}

// GetByID получает подписку по ID
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	subscription, err := h.service.GetByID(id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to get subscription")
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// List возвращает список подписок
func (h *SubscriptionHandler) List(c *gin.Context) {
	var userID *uuid.UUID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		parsedUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
			return
		}
		userID = &parsedUUID
	}

	var serviceName *string
	if serviceNameStr := c.Query("service_name"); serviceNameStr != "" {
		serviceName = &serviceNameStr
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	subscriptions, err := h.service.List(userID, serviceName, limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list subscriptions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscriptions"})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// Update обновляет подписку
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to update subscription")
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	h.logger.WithField("subscription_id", id).Info("Subscription updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Subscription updated successfully"})
}

// Delete удаляет подписку
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to delete subscription")
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	h.logger.WithField("subscription_id", id).Info("Subscription deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted successfully"})
}

// GetTotalCost подсчитывает общую стоимость подписок
func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
	startPeriod := c.Query("start_period")
	endPeriod := c.Query("end_period")

	if startPeriod == "" || endPeriod == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_period and end_period are required"})
		return
	}

	var userID *uuid.UUID
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		parsedUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
			return
		}
		userID = &parsedUUID
	}

	var serviceName *string
	if serviceNameStr := c.Query("service_name"); serviceNameStr != "" {
		serviceName = &serviceNameStr
	}

	h.logger.WithFields(logrus.Fields{
		"start_period": startPeriod,
		"end_period":   endPeriod,
		"user_id":      userID,
		"service_name": serviceName,
	}).Info("Calculating total cost")

	result, err := h.service.GetTotalCost(userID, serviceName, startPeriod, endPeriod)
	if err != nil {
		h.logger.WithError(err).Error("Failed to calculate total cost")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total cost"})
		return
	}

	c.JSON(http.StatusOK, result)
}
