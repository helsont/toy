package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4/middleware"
)

type ErrorJSON struct {
	Message string `json:"string"`
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
