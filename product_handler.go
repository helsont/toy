package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// https://stackoverflow.com/questions/51643293/how-to-query-a-gorm-model
func getProductsAPI(c echo.Context) error {
	products, err := getProducts()

	if err != nil {
		fmt.Printf("[getProductByIDAPI] Error fetching products: %s\n", fmt.Sprint(err))
		return c.JSON(500, "Unable to get the Products")
	}

	length := len(products)
	var jsonProducts = make([]ProductJSON, length)
	for i := 0; i < length; i++ {
		jsonProducts[i] = ToProductJson(&products[i])
	}
	return c.JSON(200, jsonProducts)
}

func getProductByIDAPI(c echo.Context) error {
	id := c.Param("id")

	fmt.Printf("[getProductByIDAPI] ID: %s\n", id)

	product, err := getProductByID(id)

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return c.JSON(404, "No product found")
		}
		fmt.Printf("[getProductByIDAPI] Error fetching product by id: %s\n", fmt.Sprint(err))
		return c.JSON(500, "Unable to fetch the Product")
	}

	return c.JSON(200, ToProductJson(product))
}

func createProductAPI(c echo.Context) error {
	code := c.QueryParams().Get("code")
	priceStr := c.QueryParams().Get("price")

	fmt.Printf("[createProductAPI] Code: %s Price: %s\n", code, priceStr)

	price, err := strconv.Atoi(priceStr)

	// TODO: Implement proper form validation
	if err != nil {
		fmt.Printf("[createProductAPI] Bad request: %s\n", fmt.Sprint(err))
		return c.JSON(400, "Bad request")
	}

	productOrm, err2 := createProduct(code, uint(price))
	if err2 != nil {
		fmt.Printf("[createProductAPI] Error creating the product: %s\n", fmt.Sprint(err2))
		return c.JSON(500, "Unable to create the Product")
	}

	c.JSON(200, ToProductJson(productOrm))
	return nil
}

func deleteProductByIDAPI(c echo.Context) error {
	id := c.Param("id")
	err := deleteProductByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error deleting the Product")
	}
	return c.NoContent(http.StatusOK)
}
