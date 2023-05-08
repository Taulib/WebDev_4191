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
	"quiz1.taulib.net/internal/data"
)

// Application version number

const version = "1.0.0"

// the config settings
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// Dependency Injection
type application struct {
	config config
	logger *log.Logger
	Models data.Models
}

func main() {
	var tum config

	// read the flags that are needed to populate config
	flag.IntVar(&tum.port, "port", 4000, "API Server Port")
	flag.StringVar(&tum.env, "env", "development", "Enviornment (development | staging | production)")
	flag.StringVar(&tum.db.dsn, "db-dsn", os.Getenv("COURSES_DB_DSN"), "PostgresSQL DSN")
	flag.IntVar(&tum.db.maxOpenConns, "db-max-open-conns", 25, "PostgresSQL max open connections")
	flag.IntVar(&tum.db.maxIdleConns, "db-max-idle-conns", 25, "PostgresSQL max idle connections")
	flag.StringVar(&tum.db.maxIdleTime, "db-max-idle-time", "15m", "PostgresSQL max connection idle time")
	flag.Parse()

	// creating a logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	//createing the connection pool
	db, err := openDB(tum)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	//log sucessful connection pool
	logger.Println("database connection pool was established")
	//creating a instance of our application structure
	app := &application{
		config: tum,
		logger: logger,
		Models: data.NewModels(db),
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
	db.SetMaxOpenConns(tum.db.maxOpenConns)
	db.SetMaxIdleConns(tum.db.maxIdleConns)
	duration, err := time.ParseDuration(tum.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	// create a context with a 5-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// we weill create a map that will store validation errors
type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// methods that operate on our Validator type
// check if the mao has any entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// add an entry to the map if the key does not already exist
func (v *Validator) AddError(key, message string) {
	// check if the key is already in the map
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// check to see if a element can be found in a list of items
func In(element string, list ...string) bool {
	for i := range list {
		if element == list[i] {
			return true
		}

	}
	return false
}

//
