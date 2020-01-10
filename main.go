package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	// "github.com/julienschmidt/httprouter"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4/middleware"
	"github.com/helsont/toy/models/product"
)

type ErrorJSON struct {
	Message string `json:"string"`
}

func toProductJSON(product *Product) ProductJSON {
	return ProductJSON{ID: product.ID, Code: product.Code, Price: product.Price, CreatedAt: product.CreatedAt, UpdatedAt: product.UpdatedAt, DeletedAt: product.DeletedAt}
}

func fromProductJSON(value string) (*ProductJSON, error) {
	var product *ProductJSON
	byteValue := []byte(value)
	err := json.Unmarshal(byteValue, &product)
	return product, err
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
	listener, err := net.Listen("tcp", ":3001")
	if err != nil {
		e.Logger.Fatal(listener)
	}
	e.Listener = listener

	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// err := echoServer(w, r)
			// if err != nil {
			// 	log.Printf("echo server: %v", err)
			// }
		}),
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}
	e.Server = server

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/products", getProductsAPI)
	e.GET("/products/:id", getProductByIDAPI)
	e.POST("/products", createProductAPI)
	e.DELETE("/products/:id", deleteProductByIDAPI)

	// e.GET("/ws", handleWebsocket)

	return e
}

func InitDB() *gorm.DB {
	conn, err := gorm.Open("postgres", "host=localhost port=5432 user=helson dbname=toy sslmode=disable")

	if err != nil {
		fmt.Printf("[InitDB] Error intializing database:\n %s", fmt.Sprint(err))
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

	// Start server
	e.Logger.Fatal(e.Start(""))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// func handleWebsocket(w http.ResponseWriter, r *http.Request) {

// }

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
		jsonProducts[i] = toProductJSON(&products[i])
	}
	return c.JSON(200, jsonProducts)
}

func getProducts() ([]Product, error) {
	products := make([]Product, 0)

	err := db.Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
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

	return c.JSON(200, toProductJSON(product))
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

	c.JSON(200, toProductJSON(productOrm))
	return nil
}

func getProductByID(id string) (*Product, error) {
	var product Product
	err := db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
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

func deleteProductByIDAPI(c echo.Context) error {
	id := c.Param("id")
	err := deleteProductByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error deleting the Product")
	}
	return c.NoContent(http.StatusOK)
}

func deleteProductByID(id string) error {
	err := db.Where("id = ?", id).Delete(Product{}).Error
	if err != nil {
		return err
	}
	return nil
}

func handleWebsocket(c echo.Context) error {

	// websocket.Handler(func(ws *websocket.Conn) {
	// 	defer ws.Close()
	// 	for {
	// 		// Write
	// 		err := websocket.Message.Send(ws, "Hello, Client!")
	// 		if err != nil {
	// 			c.Logger().Error(err)
	// 		}

	// 		// Read
	// 		msg := ""
	// 		err = websocket.Message.Receive(ws, &msg)
	// 		if err != nil {
	// 			c.Logger().Error(err)
	// 		}
	// 		fmt.Printf("%s\n", msg)
	// 	}
	// }).ServeHTTP(c.Response(), c.Request())
	return nil
}
