package handlers

import (
        "net/http"
        "time"

        "github.com/gin-gonic/gin"
        "github.com/google/uuid"
        "gorm.io/gorm"

        "xpense/backend/middleware"
        "xpense/backend/models"
        "xpense/backend/utils"
)

type ExpenseHandler struct {
        db *gorm.DB
}

func NewExpenseHandler(db *gorm.DB) *ExpenseHandler {
        return &ExpenseHandler{db: db}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
        userID, err := middleware.GetUserID(c)
        if err != nil {
                c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
                return
        }

        var req models.CreateExpenseRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", utils.FormatValidationErrors(err)))
                return
        }

        var category models.Category
        if err := h.db.Where("id = ? AND (user_id = ?)", req.CategoryID, userID).First(&category).Error; err != nil {
                if err == gorm.ErrRecordNotFound {
                        c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid category", "Category not found or not accessible"))
                        return
                }
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Database error", err.Error()))
                return
        }

        expense := models.Expense{
                UserID:      userID,
                Title:       req.Title,
                Amount:      req.Amount,
                CategoryID:  req.CategoryID,
                Date:        req.Date,
                Description: req.Description,
        }

        if err := h.db.Create(&expense).Error; err != nil {
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create expense", err.Error()))
                return
        }

        if err := h.db.Preload("Category").First(&expense, expense.ID).Error; err != nil {
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to load expense", err.Error()))
                return
        }

        c.JSON(http.StatusCreated, utils.SuccessResponse("Expense created successfully", expense))
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
        userID, err := middleware.GetUserID(c)
        if err != nil {
                c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
                return
        }

        var filters models.ExpenseFilters
        if err := c.ShouldBindQuery(&filters); err != nil {
                c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid query parameters", err.Error()))
                return
        }

        if filters.Page < 1 {
                filters.Page = 1
        }
        if filters.Limit < 1 || filters.Limit > 100 {
                filters.Limit = 20
        }
        if filters.SortBy == "" {
                filters.SortBy = "date"
        }
        if filters.SortOrder != "asc" && filters.SortOrder != "desc" {
                filters.SortOrder = "desc"
        }

        query := h.db.Where("user_id = ?", userID)

        if filters.CategoryID != nil {
                query = query.Where("category_id = ?", *filters.CategoryID)
        }
        if filters.StartDate != nil {
                query = query.Where("date >= ?", *filters.StartDate)
        }
        if filters.EndDate != nil {
                query = query.Where("date <= ?", *filters.EndDate)
        }
        if filters.MinAmount != nil {
                query = query.Where("amount >= ?", *filters.MinAmount)
        }
        if filters.MaxAmount != nil {
                query = query.Where("amount <= ?", *filters.MaxAmount)
        }
        if filters.Search != "" {
                query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
        }

        var total int64
        if err := query.Model(&models.Expense{}).Count(&total).Error; err != nil {
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to count expenses", err.Error()))
                return
        }

        offset := (filters.Page - 1) * filters.Limit
        var expenses []models.Expense
        if err := query.Preload("Category").
                Order(filters.SortBy + " " + filters.SortOrder).
                Limit(filters.Limit).
                Offset(offset).
                Find(&expenses).Error; err != nil {
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve expenses", err.Error()))
                return
        }

        totalPages := (int(total) + filters.Limit - 1) / filters.Limit

        response := gin.H{
                "data": expenses,
                "pagination": gin.H{
                        "page":        filters.Page,
                        "limit":       filters.Limit,
                        "total":       total,
                        "total_pages": totalPages,
                        "has_next":    filters.Page < totalPages,
                        "has_prev":    filters.Page > 1,
                },
        }

        c.JSON(http.StatusOK, utils.SuccessResponse("Expenses retrieved successfully", response))
}

func (h *ExpenseHandler) GetExpense(c *gin.Context) {
        userID, err := middleware.GetUserID(c)
        if err != nil {
                c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
                return
        }

        expenseID, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid expense ID", "Expense ID must be a valid UUID"))
                return
        }

        var expense models.Expense
        if err := h.db.Preload("Category").Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
                if err == gorm.ErrRecordNotFound {
                        c.JSON(http.StatusNotFound, utils.ErrorResponse("Expense not found", "The requested expense does not exist"))
                        return
                }
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve expense", err.Error()))
                return
        }

        c.JSON(http.StatusOK, utils.SuccessResponse("Expense retrieved successfully", expense))
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
        userID, err := middleware.GetUserID(c)
        if err != nil {
                c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
                return
        }

        expenseID, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid expense ID", "Expense ID must be a valid UUID"))
                return
        }

        var req models.UpdateExpenseRequest
        if err := c.ShouldBindJSON(&req); err != nil {
                c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", utils.FormatValidationErrors(err)))
                return
        }

        var expense models.Expense
        if err := h.db.Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
                if err == gorm.ErrRecordNotFound {
                        c.JSON(http.StatusNotFound, utils.ErrorResponse("Expense not found", "The requested expense does not exist"))
                        return
                }
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to find expense", err.Error()))
                return
        }

        if req.CategoryID != nil {
                var category models.Category
                if err := h.db.Where("id = ? AND (user_id = ? OR user_id is null)", *req.CategoryID, userID, true).First(&category).Error; err != nil {
                        if err == gorm.ErrRecordNotFound {
                                c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid category", "Category not found or not accessible"))
                                return
                        }
                        c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Database error", err.Error()))
                        return
                }
        }

        updates := make(map[string]any)
        if req.Title != nil {
                updates["title"] = *req.Title
        }
        if req.Amount != nil {
                updates["amount"] = *req.Amount
        }
        if req.CategoryID != nil {
                updates["category_id"] = *req.CategoryID
        }
        if req.Date != nil {
                updates["date"] = *req.Date
        }
        if req.Description != nil {
                updates["description"] = *req.Description
        }

        if len(updates) == 0 {
                c.JSON(http.StatusBadRequest, utils.ErrorResponse("No updates provided", "At least one field must be updated"))
                return
        }

        updates["updated_at"] = time.Now()

        if err := h.db.Model(&expense).Updates(updates).Error; err != nil {
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update expense", err.Error()))
                return
        }

        if err := h.db.Preload("Category").First(&expense, expense.ID).Error; err != nil {
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to load updated expense", err.Error()))
                return
        }

        c.JSON(http.StatusOK, utils.SuccessResponse("Expense updated successfully", expense))
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
        userID, err := middleware.GetUserID(c)
        if err != nil {
                c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
                return
        }

        expenseID, err := uuid.Parse(c.Param("id"))
        if err != nil {
                c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid expense ID", "Expense ID must be a valid UUID"))
                return
        }

        // Find existing expense
        var expense models.Expense
        if err := h.db.Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
                if err == gorm.ErrRecordNotFound {
                        c.JSON(http.StatusNotFound, utils.ErrorResponse("Expense not found", "The requested expense does not exist"))
                        return
                }
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to find expense", err.Error()))
                return
        }

        if err := h.db.Delete(&expense).Error; err != nil {
                c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete expense", err.Error()))
                return
        }

        c.JSON(http.StatusOK, utils.SuccessResponse("Expense deleted successfully", nil))
}
