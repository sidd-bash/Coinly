package models

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Username string  `json:"username" gorm:"unique"`
    Password string  `json:"-"` // don't expose hashed password in responses
    Balance  float64 `json:"balance"`
}

type Trade struct {
    gorm.Model
    UserID    uint    `json:"user_id"`
    Symbol    string  `json:"symbol"`
    Quantity  float64 `json:"quantity"`
    Price     float64 `json:"price"`
    TradeType string  `json:"trade_type"`
}
