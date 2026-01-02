package main

import (
	"fmt"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/http"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/api/http/handlers"
)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println(handlers.TestHandler())
	fmt.Println("API is running...")
	http.HelloHandler()
}
