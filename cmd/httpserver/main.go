package main

import (
	"bytes"
	"context"
	"httpServer/internal/handler"
	"httpServer/internal/http1"
	"httpServer/internal/telemetry"
	"log"
	"os"
)

func main() {

	// Initalizing environment varaibles
	getEnvFile()
	addr, there := os.LookupEnv("ADDRESS")
	if !there {
		log.Fatal("Address was not set in environment variable")
		os.Exit(3)
	}

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
		log.Fatal("Address was not set in environment variable")
		os.Exit(3)
	}
}

func getEnvFile() {
	content, err := os.ReadFile(".env")
	if err != nil {
		log.Print(".env file is not in current directory")
		return
	}
	lines := bytes.Split(content, []byte("\n"))

	for _, line := range lines {
		key, val, found := bytes.Cut(line, []byte("="))
		if !found {
			continue
		}
		// fmt.Printf("key: %s, value: %s\n", key, val)
		err := os.Setenv(string(key), string(val))
		if err != nil {
			panic("env key & values are wrong")
		}
	}
}
