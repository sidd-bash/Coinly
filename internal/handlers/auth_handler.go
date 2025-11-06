package handlers

import (
    "net/http"
    "os"
    "time"

    "github.com/sidd-bash/coinly/internal/config"
    "github.com/sidd-bash/coinly/internal/models"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

// --- REGISTER ---
func Register(c *gin.Context) {
    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
        return
    }

    user := models.User{Username: input.Username, Password: string(hashed), Balance: 10000}
    if err := config.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists"})
        return
    }

    // ✅ Auto-create INR wallet with ₹100,000
    wallet := models.Wallet{UserID: user.ID, Currency: "INR", Balance: 100000}
    config.DB.Create(&wallet)

    c.JSON(http.StatusCreated, gin.H{
        "message": "user registered successfully",
        "wallet":  wallet,
    })

}

// --- LOGIN ---
func Login(c *gin.Context) {
    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    result := config.DB.Where("username = ?", input.Username).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // --- dynamically fetch secret here ---
    jwtSecret := []byte(os.Getenv("JWT_SECRET"))

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id":  user.ID,
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
