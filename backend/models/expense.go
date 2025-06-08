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
	Category  	string     	   `json:"category" gorm:"not null"`
	Date        time.Time      `json:"date" gorm:"not null;index"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
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
	Category  	string 	  `json:"category" binding:"required"`
	Date        time.Time `json:"date" binding:"required"`
	Description string    `json:"description" binding:"max=1000"`
}

type UpdateExpenseRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=255"`
	Amount      *float64   `json:"amount" binding:"omitempty,gt=0"`
	Category  	*string    `json:"category"`
	Date        *time.Time `json:"date"`
	Description *string    `json:"description" binding:"omitempty,max=1000"`
}

type ExpenseFilters struct {
	Category *string `form:"category"`
}
