package handlers

import (
    "net/http"
    "github.com/sidd-bash/coinly/internal/config"
    "github.com/sidd-bash/coinly/internal/models"
    "github.com/gin-gonic/gin"
)

func CreateTrade(c *gin.Context) {
    var trade models.Trade
    if err := c.ShouldBindJSON(&trade); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    config.DB.First(&user, trade.UserID)

    if user.ID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    if trade.TradeType == "BUY" {
        user.Balance -= trade.Price * trade.Quantity
    } else if trade.TradeType == "SELL" {
        user.Balance += trade.Price * trade.Quantity
    }

    config.DB.Save(&user)
    config.DB.Create(&trade)

    c.JSON(http.StatusCreated, gin.H{"trade": trade, "updated_balance": user.Balance})
}

func GetTrades(c *gin.Context) {
    var trades []models.Trade
    config.DB.Find(&trades)
    c.JSON(http.StatusOK, trades)
}
