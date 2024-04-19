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

	"github.com/a1d1yar/assingment1_Golang/internal/data"
	_ "github.com/lib/pq"
)

type configuration struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config   configuration
	logger   *log.Logger
	database *data.DBModel
}

func main() {
	var cfg configuration

	flag.IntVar(&cfg.port, "port", 5432, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:Aldiyar2004@localhost:5432/a.maratovDB?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDatabase(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("Database connection pool established")

	dbModel := &data.DBModel{
		DB: db,
	}

	app := &application{
		config:   cfg,
		logger:   logger,
		database: dbModel,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("Starting %s server on %s", cfg.env, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)
}

func openDatabase(cfg configuration) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
