package routes

import (
    "github.com/gin-gonic/gin"
    "coinly/internal/handlers"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Welcome to Coinly API 🚀"})
    })

    // User routes
    r.POST("/users", handlers.CreateUser)
    r.GET("/users", handlers.GetUsers)

    // Trade routes
    r.POST("/trades", handlers.CreateTrade)
    r.GET("/trades", handlers.GetTrades)

    return r
}
