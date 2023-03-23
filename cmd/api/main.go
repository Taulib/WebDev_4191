//FileName: Quiz1 cmd/api/main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Application version number

const version = "1.0.0"

// the config settings
type config struct {
	port int
	env  string
}

// Dependency Injection
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var tum config

	// read the flags that are needed to populate config
	flag.IntVar(&tum.port, "port", 4000, "API Server Port")
	flag.StringVar(&tum.env, "env", "development", "Enviornment (development | staging | production)")
	flag.Parse()

	// creating a logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	//creating a instance of our application structure
	app := &application{
		config: tum,
		logger: logger,
	}

	//creating our ServerMux
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	//creating our http Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", tum.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Starting our server
	logger.Printf("Starting %s server is on %s", tum.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
