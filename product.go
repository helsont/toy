package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

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

func ToProductJson(product *Product) ProductJSON {
	return ProductJSON{ID: product.ID, Code: product.Code, Price: product.Price, CreatedAt: product.CreatedAt, UpdatedAt: product.UpdatedAt, DeletedAt: product.DeletedAt}
}

func FromProductJSON(value string) (*ProductJSON, error) {
	var product *ProductJSON
	byteValue := []byte(value)
	err := json.Unmarshal(byteValue, &product)
	return product, err
}

func createProduct(code string, price uint) (*Product, error) {
	product := &Product{Code: code, Price: price}
	err := db.Create(product).Error

	if err != nil {
		fmt.Printf("[createProduct] Error creating the product: %s\n", fmt.Sprint(err))
		return nil, err
	}

	return product, nil
}

func getProducts() ([]Product, error) {
	products := make([]Product, 0)

	err := db.Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}

func getProductByID(id string) (*Product, error) {
	var product Product
	err := db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func deleteProductByID(id string) error {
	err := db.Where("id = ?", id).Delete(Product{}).Error
	if err != nil {
		return err
	}
	return nil
}
