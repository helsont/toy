package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	// "github.com/julienschmidt/httprouter"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4/middleware"
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

var db *gorm.DB

// GetDB : returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}

// GetHandler : Server handler
func GetHandler() *echo.Echo {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/products", getProducts)
	e.POST("/products", createProductAPI)

	return e
}

func InitDB() *gorm.DB {
	conn, err := gorm.Open("postgres", "host=localhost port=5432 user=helson dbname=toy sslmode=disable")

	if err != nil {
		fmt.Printf(err.Error())
		panic("failed to connect database")
	}

	// defer conn.Close()

	// Migrate the schema
	conn.AutoMigrate(&Product{})

	db = conn

	return conn
}

func TearDown() {
	db.Close()
}

func main() {
	// Initialize web server
	e := GetHandler()

	// Initialize Postgres connection
	InitDB()

	// Teardown dependencies
	defer TearDown()

	// // Create
	// db.Create(&Product{Code: "L1212", Price: 1000})

	// // Read
	// var product Product
	// db.First(&product, 1)                   // find product with id 1
	// db.First(&product, "code = ?", "L1212") // find product with code l1212

	// // Update - update product's price to 2000
	// db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	// db.Delete(&product)

	// Start server
	e.Logger.Fatal(e.Start(":3001"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// https://stackoverflow.com/questions/51643293/how-to-query-a-gorm-model
func getProducts(c echo.Context) error {
	conn := GetDB()
	products := make([]Product, 0)

	err := conn.Find(&products).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return c.JSON(404, struct {
				Message string `json:"message"`
			}{"record not found"})
		}
		c.Logger().Error(err)
	}

	return c.JSON(200, products)
}

func createProductAPI(c echo.Context) error {
	code := c.QueryParams().Get("code")
	priceStr := c.QueryParams().Get("price")
	price, err := strconv.Atoi(priceStr)

	// TODO: Implement proper form validation
	if err != nil {
		return c.JSON(400, "Bad request")
	}

	_, err2 := createProduct(code, uint(price))
	if err2 != nil {
		return c.JSON(500, "Unable to create the Product")
	}

	return nil
}

func createProduct(code string, price uint) (*Product, error) {
	product := &Product{Code: code, Price: price}
	err := db.Create(product).Error

	if err != nil {
		return nil, err
	}

	return product, nil
}
