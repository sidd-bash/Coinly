package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/sidd-bash/coinly/internal/handlers"
    "github.com/sidd-bash/coinly/internal/middleware"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Welcome to Coinly API 🚀"})
    })

    // Auth routes
    r.POST("/register", handlers.Register)
    r.POST("/login", handlers.Login)

    // Protected group
    auth := r.Group("/api")
    auth.Use(middleware.AuthMiddleware())

    auth.POST("/trades", handlers.CreateTrade)
    auth.GET("/trades", handlers.GetTrades)
    auth.GET("/users", handlers.GetUsers)

    return r
}