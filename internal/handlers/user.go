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

type UserHandler struct {
	service *service.UserService
	logger  *logrus.Logger
}

func NewUserHandler(service *service.UserService, logger *logrus.Logger) *UserHandler {
	if logger == nil {
		logger = logrus.New()
	}

	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

// Create создает нового пользователя
// @Summary Создать пользователя
// @Description Создает нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "Данные пользователя"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithFields(logrus.Fields{
		"name":  req.Name,
		"email": req.Email,
	}).Info("Creating user")

	user, err := h.service.Create(&req)
	if err != nil {
		h.logger.WithError(err).Error("Failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithField("user_id", user.ID).Info("User created successfully")
	c.JSON(http.StatusCreated, user)
}

// GetByID получает пользователя по ID
// @Summary Получить пользователя по ID
// @Description Возвращает пользователя по указанному ID
// @Tags users
// @Produce json
// @Param id path string true "ID пользователя (UUID)"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	h.logger.WithField("user_id", id).Info("Getting user by ID")

	user, err := h.service.GetByID(id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to get user")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// List возвращает список пользователей
// @Summary Получить список пользователей
// @Description Возвращает список пользователей с пагинацией
// @Tags users
// @Produce json
// @Param limit query int false "Лимит записей"
// @Param offset query int false "Смещение"
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	h.logger.WithFields(logrus.Fields{
		"limit":  limit,
		"offset": offset,
	}).Info("Listing users")

	users, err := h.service.List(limit, offset)
	if err != nil {
		h.logger.WithError(err).Error("Failed to list users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Update обновляет пользователя
// @Summary Обновить пользователя
// @Description Обновляет существующего пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "ID пользователя (UUID)"
// @Param user body models.UpdateUserRequest true "Данные для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithField("user_id", id).Info("Updating user")

	err = h.service.Update(id, &req)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to update user")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.logger.WithField("user_id", id).Info("User updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Delete удаляет пользователя
// @Summary Удалить пользователя
// @Description Удаляет пользователя по ID
// @Tags users
// @Produce json
// @Param id path string true "ID пользователя (UUID)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	h.logger.WithField("user_id", id).Info("Deleting user")

	err = h.service.Delete(id)
	if err != nil {
		h.logger.WithError(err).WithField("id", id).Error("Failed to delete user")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	h.logger.WithField("user_id", id).Info("User deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
