package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Expense struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Title       string         `json:"title" gorm:"not null"`
	Amount      float64        `json:"amount" gorm:"not null"`
	CategoryID  uuid.UUID      `json:"category_id" gorm:"type:uuid;not null"`
	Date        time.Time      `json:"date" gorm:"not null;index"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

func (e *Expense) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type CreateExpenseRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=255"`
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Description string    `json:"description" binding:"max=1000"`
}

type UpdateExpenseRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=255"`
	Amount      *float64   `json:"amount" binding:"omitempty,gt=0"`
	CategoryID  *uuid.UUID `json:"category_id"`
	Date        *time.Time `json:"date"`
	Description *string    `json:"description" binding:"omitempty,max=1000"`
}

type ExpenseFilters struct {
	CategoryID *uuid.UUID `form:"category_id"`
	StartDate  *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate    *time.Time `form:"end_date" time_format:"2006-01-02"`
	MinAmount  *float64   `form:"min_amount"`
	MaxAmount  *float64   `form:"max_amount"`
	Search     string     `form:"search"`
	Page       int        `form:"page,default=1"`
	Limit      int        `form:"limit,default=20"`
	SortBy     string     `form:"sort_by,default=date"`
	SortOrder  string     `form:"sort_order,default=desc"`
}
