package main

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestEchoHandler(t *testing.T) {
	handler := GetHandler()
	_ = InitDB()
	defer TearDown()

	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	testEcho(e)
}

func testEcho(e *httpexpect.Expect) {
	e.GET("/products").Expect().Status(200)
}
