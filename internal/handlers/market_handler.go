package handlers

import (
    "net/http"

    "github.com/sidd-bash/coinly/internal/services"
    "github.com/gin-gonic/gin"
)

func GetMarketPrices(c *gin.Context) {
    prices, err := services.GetPrices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"prices_in_inr": prices})
}
