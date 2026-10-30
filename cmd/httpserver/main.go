package main

import (
	"context"
	"httpServer/internal/handler"
	"httpServer/internal/http1"
	"httpServer/internal/telemetry"
	"log"
)

const addr = ":8080"

func main() {

	// for debugging purposes
	// 	f, err := os.Create("trace.out")
	// 	if err != nil {
	// 		log.Fatalf("failed to create trace output file: %v", err)
	// 	}
	// 	defer f.Close()
	//
	// 	if err := trace.Start(f); err != nil {
	// 		log.Fatalf("failed to start trace: %v", err)
	// 	}
	// 	defer trace.Stop()

	// init tracing
	tp, err := telemetry.InitTracer()
	if err != nil {
		log.Fatalf("failed to initialize tracer: %v", err)
	}
	defer tp.Shutdown(context.Background())

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

	// adding telemetry middleware
	wrappedMux := telemetry.TracingMiddleWare(mux)

	if err := http1.ListenAndServe(addr, wrappedMux); err != nil {
		panic(err)
	}
}
