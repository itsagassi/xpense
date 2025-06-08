package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"xpense/backend/config"
	"xpense/backend/database"
	"xpense/backend/handlers"
	"xpense/backend/middleware"
)

func main() {
        cfg := config.Load()

        db, err := database.Initialize(cfg.DatabaseURL)
        if err != nil {
                log.Fatal("Failed to connect to database:", err)
        }

        router := gin.Default()
        router.Use(func(c *gin.Context) {
                c.Header("Access-Control-Allow-Origin", "*")
                c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
                
                if c.Request.Method == "OPTIONS" {
                        c.AbortWithStatus(204)
                        return
                }
                
                c.Next()
        })

        router.GET("/health", func(c *gin.Context) {
                c.JSON(200, gin.H{"status": "healthy"})
        })

        expenseHandler := handlers.NewExpenseHandler(db)

        protected := router.Group("/api/v1")
        protected.Use(middleware.AuthMiddleware(cfg.SupabaseJWTSecret))
        {
                protected.POST("/expenses", expenseHandler.CreateExpense)
                protected.GET("/expenses", expenseHandler.GetExpenses)
                protected.GET("/expenses/:id", expenseHandler.GetExpense)
                protected.PUT("/expenses/:id", expenseHandler.UpdateExpense)
                protected.DELETE("/expenses/:id", expenseHandler.DeleteExpense)
        }

        port := os.Getenv("PORT")
        if port == "" {
                port = "5000"
        }

        log.Printf("Server starting on port %s", port)
        if err := router.Run("0.0.0.0:" + port); err != nil {
                log.Fatal("Failed to start server:", err)
        }
}
