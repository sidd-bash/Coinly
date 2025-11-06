package models

import "gorm.io/gorm"

type Wallet struct {
    gorm.Model
    UserID   uint    `json:"user_id"`
    Currency string  `json:"currency"`
    Balance  float64 `json:"balance"`
}
