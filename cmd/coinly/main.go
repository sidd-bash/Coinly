package main

import (
    "fmt"
    "os"

    "github.com/joho/godotenv"
    "github.com/sidd-bash/coinly/internal/config"
    "github.com/sidd-bash/coinly/internal/models"
    "github.com/sidd-bash/coinly/internal/routes"
)

func main() {
    // 🔹 Load environment variables once globally
    if err := godotenv.Load(); err != nil {
        fmt.Println("⚠️ No .env file found — using system environment variables")
    } else {
        fmt.Println("✅ Loaded .env successfully")
    }

    config.Init()

    config.DB.AutoMigrate(
        &models.User{},
        &models.Trade{},
        &models.Wallet{},
        &models.Order{},
        &models.Transaction{},
    )


    r := routes.SetupRouter()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Println("🚀 Starting Coinly backend on port", port)
    r.Run(":" + port)
}
