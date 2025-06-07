package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"xpense/backend/middleware"
	"xpense/backend/models"
	"xpense/backend/utils"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

type CreateCategoryRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Color string `json:"color" binding:"required,hexcolor"`
}

type UpdateCategoryRequest struct {
	Name  *string `json:"name" binding:"omitempty,min=1,max=100"`
	Color *string `json:"color" binding:"omitempty,hexcolor"`
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category

	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
		return
	}

	if err := h.db.Where("user_id = ? or user_id is null", userID).Order("name ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve categories", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Categories retrieved successfully", categories))
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", utils.FormatValidationErrors(err)))
		return
	}

	var existing models.Category
	if err := h.db.Where("name = ? AND user_id = ?", req.Name, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, utils.ErrorResponse("Category already exists", "A category with this name already exists"))
		return
	}

	category := models.Category{
		Name:   req.Name,
		Color:  req.Color,
		UserID: &userID,
	}

	if err := h.db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create category", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Category created successfully", category))
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid category ID", "Must be a valid UUID"))
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", utils.FormatValidationErrors(err)))
		return
	}

	var category models.Category
	if err := h.db.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Category not found", "You can only update your own categories"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to find category", err.Error()))
		return
	}

	if req.Name != nil && *req.Name != category.Name {
		var existing models.Category
		if err := h.db.Where("name = ? AND user_id = ? AND id != ?", *req.Name, userID, categoryID).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, utils.ErrorResponse("Category name already exists", "A category with this name already exists"))
			return
		}
	}

	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("No updates provided", "At least one field must be updated"))
		return
	}

	if err := h.db.Model(&category).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update category", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Category updated successfully", category))
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid category ID", "Must be a valid UUID"))
		return
	}

	var category models.Category
	if err := h.db.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Category not found", "You can only delete your own categories"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to find category", err.Error()))
		return
	}

	var count int64
	if err := h.db.Model(&models.Expense{}).Where("category_id = ? AND user_id = ?", categoryID, userID).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to check category usage", err.Error()))
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Cannot delete category", "Category is used in existing expenses"))
		return
	}

	if err := h.db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete category", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Category deleted successfully", nil))
}
