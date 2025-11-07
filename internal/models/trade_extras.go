package models

import "gorm.io/gorm"

type Order struct {
    gorm.Model
    UserID uint    `json:"user_id"`
    Type   string  `json:"type"`   // BUY or SELL
    Symbol string  `json:"symbol"` // BTC, ETH, etc.
    Amount float64 `json:"amount"`
    Price  float64 `json:"price"`
    Status string  `json:"status"` // SUCCESS, FAILED
}

type Transaction struct {
    gorm.Model
    UserID       uint    `json:"user_id"`
    FromCurrency string  `json:"from_currency"`
    ToCurrency   string  `json:"to_currency"`
    Amount       float64 `json:"amount"`
    Rate         float64 `json:"rate"`
}
