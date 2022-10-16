package models

import "time"

type Product struct {
	Name        string        `json:"name" validate:"required"`
	Description string        `json:"description"`
	IsActive    bool          `json:"is_active"`
	CreatedAt   time.Duration `json:"created_at"`
	ProductType string        `json:"product_type"`
	Price       float64       `json:"price"`
	Currency    string        `json:"currency"`
}
