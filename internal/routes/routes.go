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

        
    auth.GET("/users", handlers.GetUsers)

    // Trading
    auth.POST("/trades", handlers.CreateTrade)
    auth.GET("/trades", handlers.GetTrades)


    // Wallets
    auth.GET("/wallets", handlers.GetWallets)

    // Market
    auth.GET("/market/prices", handlers.GetMarketPrices)

    // Trade
    auth.POST("/trade/buy", handlers.BuyCrypto)
    auth.POST("/trade/sell", handlers.SellCrypto)
    auth.GET("/orders", handlers.GetOrders)

    return r
}