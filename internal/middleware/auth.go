package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {

    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        secret := []byte(os.Getenv("JWT_SECRET"))

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return secret, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            c.Abort()
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        userIDFloat, ok := claims["user_id"].(float64)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
            c.Abort()
            return
        }

        userID := uint(userIDFloat)
        c.Set("user_id", userID)
        c.Set("username", claims["username"])

        c.Next()
    }
}
