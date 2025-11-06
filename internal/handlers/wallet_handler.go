package handlers

import (
    "net/http"

    "github.com/sidd-bash/coinly/internal/config"
    "github.com/sidd-bash/coinly/internal/models"
    "github.com/gin-gonic/gin"
)

func GetWallets(c *gin.Context) {
    userID := c.GetUint("user_id")

    var wallets []models.Wallet
    config.DB.Where("user_id = ?", userID).Find(&wallets)

    c.JSON(http.StatusOK, gin.H{"wallets": wallets})
}
