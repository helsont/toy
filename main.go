package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})
	e.GET("/products", getProductsAPI)
	e.GET("/products/:id", getProductByIDAPI)
	e.POST("/products", createProductAPI)
	e.DELETE("/products/:id", deleteProductByIDAPI)
	e.Logger.Fatal(e.Start(":3001"))

	return e
}

func GetWebsocketHandler() {
	listener, err := net.Listen("tcp", "127.0.0.1:3002")
	if err != nil {
		panic(err)
	}
	// listener.Close()
	fmt.Println("HELO")
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
			if err != nil {
				panic(err)
			}
			defer c.Close(websocket.StatusInternalError, "the sky is falling")

			ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
			defer cancel()

			var v interface{}
			err = wsjson.Read(ctx, c, &v)
			if err != nil {
				panic(err)
			}

			log.Printf("received: %v", v)

			c.Close(websocket.StatusNormalClosure, "")

		}),
		ReadTimeout:  20 * time.Minute,
		WriteTimeout: 20 * time.Minute,
	}
	// server.Close()

	go func() {
		err := server.Serve(listener)
		fmt.Println("Starting WS server")
		if err != http.ErrServerClosed {
			log.Fatalf("failed to listen and serve: %v", err)
		}
	}()
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
	GetWebsocketHandler()

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
