package main

import (
	"httpServer/internal/handler"
	"httpServer/internal/http1"
)

const addr = ":8080"

func main() {

	hdl := http1.DefaultMux
	// Just some example Handler Funcs
	hdl.HandleFunc("/hello", handler.Hello)
	// hdl.HandleFunc("GET /users", handler.GetUsers)
	// hdl.HandleFunc("POST /users", handler.CreateUser)
	if err := http1.ListenAndServe(addr, hdl); err != nil {
		panic(err)
	}
}
