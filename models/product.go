package models

import (
	"time"
)

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null" binding:"required"`
	Description string    `json:"description" gorm:"type:varchar(255);not null" binding:"required"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null" binding:"required"`
	Stock       int       `json:"stock" gorm:"not null" binding:"required,gte=0"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func ProductSeeder() []Product {
	return []Product{
		{
			Name:        "Product A",
			Description: "This is Product A description",
			Price:       29.99,
			Stock:       100,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Product B",
			Description: "This is Product B description",
			Price:       49.99,
			Stock:       200,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Product C",
			Description: "This is Product C description",
			Price:       99.99,
			Stock:       150,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Name:        "Product D",
			Description: "This is Product D description",
			Price:       19.99,
			Stock:       50,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}
