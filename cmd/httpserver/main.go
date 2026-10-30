package main

import (
	"httpServer/internal/handler"
	"httpServer/internal/http1"
	"log"
	"os"
	"runtime/trace"
)

const addr = ":8080"

func main() {

	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer f.Close()

	if err := trace.Start(f); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer trace.Stop()

	mux := http1.DefaultMux
	// Just some example Handler Funcs
	mux.HandleFunc("/hello", handler.Hello)
	mux.HandleFunc("GET /users", handler.GetUsers)
	mux.HandleFunc("POST /users", handler.CreateUser)
	mux.HandleFunc("/metrics", handler.Telemetry)
	// 	mux.HandleFunc("GET /tag/{name}", func(w http.ResponseWriter, r *http.Request) {
	// 		name := r.PathValue("name")
	// 		fmt.Fprint(w, "Hello, %s", name)
	//
	// 	})

	if err := http1.ListenAndServe(addr, mux); err != nil {
		panic(err)
	}
}
