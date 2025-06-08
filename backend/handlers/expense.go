package handlers

import (
	"net/http"
	"time"
        "fmt"

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

	expense := models.Expense{
		UserID:      userID,
		Title:       req.Title,
		Amount:      req.Amount,
		Category:    req.Category,
		Date:        req.Date,
		Description: req.Description,
	}

	if err := h.db.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create expense", err.Error()))
		return
	}

	if err := h.db.First(&expense, expense.ID).Error; err != nil {
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

	query := h.db.Where("user_id = ?", userID)

	if filters.Category != nil {
		query = query.Where("category = ?", *filters.Category)
	}

	var total int64
	if err := query.Model(&models.Expense{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to count expenses", err.Error()))
		return
	}

	var expenses []models.Expense
	if err := query.Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve expenses", err.Error()))
		return
	}

	response := gin.H{
		"data": expenses,
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
	if err := h.db.Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, utils.ErrorResponse("Expense not found", "The requested expense does not exist"))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve expense", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Expense retrieved successfully", expense))
}

func (h *ExpenseHandler) GetExpenseTotalPerCategories(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
		return
	}

	query := h.db.Table("expenses").
		Select("category as name, SUM(amount) as value").
		Where("user_id = ?", userID).
		Group("category")

	var results []struct {
		Name  string  `json:"name"` 
		Value float64 `json:"value"`
	}

	if err := query.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve expense totals", err.Error()))
		return
	}
	if len(results) == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("No expenses found", "You have no expenses recorded"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Expense totals per category retrieved successfully", results))
}


func (h *ExpenseHandler) GetExpenseTotalPerMonth(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
		return
	}

	query := h.db.Table("expenses").
		Select("TO_CHAR(DATE_TRUNC('month', date), 'Mon') AS name, SUM(amount) AS total").
		Where("user_id = ?", userID).
		Group("name").
		Order("MIN(date)")

	var results []struct {
		Name  string  `json:"name"` 
		Total float64 `json:"total"` 
	}

	if err := query.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve monthly totals", err.Error()))
		return
	}
	if len(results) == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("No expenses found", "You have no expenses recorded"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Expense totals per month retrieved successfully", results))
}

func (h *ExpenseHandler) GetExpenseTotalPerWeek(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Unauthorized", err.Error()))
		return
	}

	query := h.db.Table("expenses").
		Select("EXTRACT(WEEK FROM date) AS week_num, SUM(amount) AS total").
		Where("user_id = ?", userID).
		Group("week_num").
		Order("week_num")

	var rawResults []struct {
		WeekNum int     `json:"week_num"`
		Total   float64 `json:"total"`
	}

	if err := query.Scan(&rawResults).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to retrieve weekly totals", err.Error()))
		return
	}
	if len(rawResults) == 0 {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("No expenses found", "You have no expenses recorded"))
		return
	}

	formatted := make([]map[string]interface{}, len(rawResults))
	for i, r := range rawResults {
		formatted[i] = map[string]interface{}{
			"name":  fmt.Sprintf("Week %d", r.WeekNum),
			"total": r.Total,
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Expense totals per week retrieved successfully", formatted))
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

	updates := make(map[string]any)
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.Category != nil {
		updates["category"] = *req.Category
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

	if err := h.db.First(&expense, expense.ID).Error; err != nil {
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
