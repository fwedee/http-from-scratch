package main

import (
	"httpServer/internal/handler"
	"httpServer/internal/http1"
	"net/http"
)

const addr = ":8080"

func main() {

	http.HandleFunc("GET /hello", handler.Hello)
	http.HandleFunc("GET /users", handler.GetUsers)
	http.HandleFunc("POST /users", handler.CreateUser)
	if err := http1.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
