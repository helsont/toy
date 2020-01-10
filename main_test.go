package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type Container struct {
	handler *echo.Echo
	server  *httptest.Server
	db      *gorm.DB
}

var container *Container

func TestMain(m *testing.M) {
	fmt.Println("Before")
	handler := GetHandler()
	server := httptest.NewServer(handler)
	db = InitDB()
	container = &Container{handler, server, db}

	fmt.Println("Starting test")
	tx := db.Begin()

	code := m.Run()

	fmt.Println("After")

	tx.Rollback()
	server.Close()
	TearDown()

	fmt.Println("Quitting")
	os.Exit(code)
}

func getTestClient(t *testing.T) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  container.server.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

func TestGetProductsEmpty(t *testing.T) {
	client := getTestClient(t)
	client.GET("/products").Expect().Status(200)
}

func TestCreateProduct(t *testing.T) {
	client := getTestClient(t)
	code := "Shoes"
	var price uint = 50

	body := client.POST("/products").
		WithQuery("code", code).
		WithQuery("price", price).
		Expect().
		Status(200).
		Body().
		Raw()

	// Read JSON from response
	var product *ProductJSON
	byteValue := []byte(body)
	err := json.Unmarshal(byteValue, &product)
	assert.NoError(t, err, "Found error parsing json")

	assert.NotNil(t, product.ID)
	assert.Equal(t, code, product.Code)
	assert.Equal(t, price, product.Price)
}

func TestGetProductByID(t *testing.T) {
	client := getTestClient(t)

	product, err := createProduct("Jack Shoes", 55)
	assert.NoError(t, err)

	body := client.GET(fmt.Sprintf("/products/%d", product.ID)).
		Expect().
		Status(200).
		Body().
		Raw()
	jsonProduct, err2 := fromProductJSON(body)
	// fmt.Printf("json product: %s\n", fmt.Sprint(jsonProduct))
	// fmt.Printf("database product: %s\n", fmt.Sprint(product))
	assert.Nil(t, err2)

	assert.Equal(t, product.ID, jsonProduct.ID)
	assert.Equal(t, product.Code, jsonProduct.Code)
	assert.Equal(t, product.Price, jsonProduct.Price)
	assert.True(t, product.CreatedAt.Equal(jsonProduct.CreatedAt))
	assert.True(t, product.UpdatedAt.Equal(jsonProduct.UpdatedAt))
	assert.Nil(t, product.DeletedAt)
	assert.Nil(t, jsonProduct.DeletedAt)
}
