package models

import "github.com/jinzhu/gorm"

import "time"

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

type ProductJSON struct {
	ID        uint       `json:"id"`
	Code      string     `json:"code"`
	Price     uint       `json:"price"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
