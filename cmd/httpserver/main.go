package main

import (
	"httpServer/internal/handler"
	"httpServer/internal/http1"
)

const addr = ":8080"

func main() {

	mux := http1.DefaultMux
	// Just some example Handler Funcs
	mux.HandleFunc("/hello", handler.Hello)
	// mux.HandleFunc("/hello", handler.Hello2)
	mux.HandleFunc("GET /users", handler.GetUsers)
	mux.HandleFunc("POST /users", handler.CreateUser)
	// 	mux.HandleFunc("GET /tag/{name}", func(w http.ResponseWriter, r *http.Request) {
	// 		name := r.PathValue("name")
	// 		fmt.Fprint(w, "Hello, %s", name)
	//
	// 	})

	if err := http1.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
