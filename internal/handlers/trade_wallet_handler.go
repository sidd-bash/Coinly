package handlers

import (
    "net/http"

    "github.com/sidd-bash/coinly/internal/config"
    "github.com/sidd-bash/coinly/internal/models"
    "github.com/sidd-bash/coinly/internal/services"
    "github.com/gin-gonic/gin"
)

func BuyCrypto(c *gin.Context) {
    userID := c.GetUint("user_id")

    var input struct {
        Symbol string  `json:"symbol"` // e.g. "BTC"
        Amount float64 `json:"amount"` // crypto amount to buy
    }

    if err := c.ShouldBindJSON(&input); err != nil || input.Amount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }

    // Get latest market price
    prices, err := services.GetPrices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch market prices"})
        return
    }

    symbol := input.Symbol
    rate := prices[symbolToCoinID(symbol)]
    if rate == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported currency"})
        return
    }

    // Total cost in INR
    totalCost := rate * input.Amount

    // Fetch INR wallet
    var inrWallet models.Wallet
    if err := config.DB.Where("user_id = ? AND currency = ?", userID, "INR").First(&inrWallet).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "INR wallet not found"})
        return
    }

    // Check INR balance
    if inrWallet.Balance < totalCost {
        c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient INR balance"})
        return
    }

    // Deduct INR, add crypto
    inrWallet.Balance -= totalCost
    config.DB.Save(&inrWallet)

    var cryptoWallet models.Wallet
    result := config.DB.Where("user_id = ? AND currency = ?", userID, symbol).First(&cryptoWallet)
    if result.RowsAffected == 0 {
        cryptoWallet = models.Wallet{UserID: userID, Currency: symbol, Balance: input.Amount}
        config.DB.Create(&cryptoWallet)
    } else {
        cryptoWallet.Balance += input.Amount
        config.DB.Save(&cryptoWallet)
    }

    // Record order & transaction
    order := models.Order{
        UserID: userID, Type: "BUY", Symbol: symbol, Amount: input.Amount,
        Price: rate, Status: "SUCCESS",
    }
    config.DB.Create(&order)

    txn := models.Transaction{
        UserID: userID, FromCurrency: "INR", ToCurrency: symbol,
        Amount: input.Amount, Rate: rate,
    }
    config.DB.Create(&txn)

    c.JSON(http.StatusOK, gin.H{
        "message":  "Buy order successful",
        "order":    order,
        "balances": gin.H{"INR": inrWallet.Balance, symbol: cryptoWallet.Balance},
    })
}

func SellCrypto(c *gin.Context) {
    userID := c.GetUint("user_id")

    var input struct {
        Symbol string  `json:"symbol"`
        Amount float64 `json:"amount"`
    }

    if err := c.ShouldBindJSON(&input); err != nil || input.Amount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
        return
    }

    // Get market price
    prices, err := services.GetPrices()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch market prices"})
        return
    }

    symbol := input.Symbol
    rate := prices[symbolToCoinID(symbol)]
    if rate == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported currency"})
        return
    }

    // Fetch crypto wallet
    var cryptoWallet models.Wallet
    if err := config.DB.Where("user_id = ? AND currency = ?", userID, symbol).First(&cryptoWallet).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no wallet found for " + symbol})
        return
    }

    // Check balance
    if cryptoWallet.Balance < input.Amount {
        c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient crypto balance"})
        return
    }

    // Deduct crypto, add INR
    cryptoWallet.Balance -= input.Amount
    config.DB.Save(&cryptoWallet)

    var inrWallet models.Wallet
    config.DB.Where("user_id = ? AND currency = ?", userID, "INR").First(&inrWallet)
    inrWallet.Balance += rate * input.Amount
    config.DB.Save(&inrWallet)

    // Record order & transaction
    order := models.Order{
        UserID: userID, Type: "SELL", Symbol: symbol, Amount: input.Amount,
        Price: rate, Status: "SUCCESS",
    }
    config.DB.Create(&order)

    txn := models.Transaction{
        UserID: userID, FromCurrency: symbol, ToCurrency: "INR",
        Amount: input.Amount, Rate: rate,
    }
    config.DB.Create(&txn)

    c.JSON(http.StatusOK, gin.H{
        "message":  "Sell order successful",
        "order":    order,
        "balances": gin.H{"INR": inrWallet.Balance, symbol: cryptoWallet.Balance},
    })
}

func GetOrders(c *gin.Context) {
    userID := c.GetUint("user_id")

    var orders []models.Order
    config.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&orders)

    c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func symbolToCoinID(symbol string) string {
    switch symbol {
    case "BTC":
        return "bitcoin"
    case "ETH":
        return "ethereum"
    default:
        return ""
    }
}
