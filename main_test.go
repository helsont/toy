package main

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/labstack/echo/v4"
)

type SomeData struct {
	handler *echo.Echo
	server  *httptest.Server
}

var container *SomeData

func TestMain(m *testing.M) {
	fmt.Println("Before")
	handler := GetHandler()
	server := httptest.NewServer(handler)
	container = &SomeData{handler, server}
	_ = InitDB()

	fmt.Println("Starting test")
	code := m.Run()
	fmt.Println("After")

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
