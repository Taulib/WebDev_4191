//FileName: Quiz1 cmd/api/main.go

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Application version number

const version = "1.0.0"

// the config settings
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
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
	flag.StringVar(&tum.db.dsn, "db-dsn", os.Getenv("COURSES_DB_DSN"), "PostgresSQL DSN")
	flag.Parse()

	// creating a logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	//createing the connection pool
	db, err := openDB(tum)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
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
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// openDB() function retruns a *sql.DB connection pool
func openDB(tum config) (*sql.DB, error) {
	db, err := sql.Open("postgres", tum.db.dsn)
	if err != nil {
		return nil, err
	}

	// create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
