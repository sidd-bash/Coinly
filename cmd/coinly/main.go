package main

import (
    "fmt"
    "os"
    "github.com/sidd-bash/coinly/internal/config"
    "github.com/sidd-bash/coinly/internal/models"
    "github.com/sidd-bash/coinly/internal/routes"
)

func main() {
    config.Init()

    // Auto migrate models
    config.DB.AutoMigrate(&models.User{}, &models.Trade{})

    r := routes.SetupRouter()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Println("🚀 Starting Coinly backend on port", port)
    r.Run(":" + port)
}
