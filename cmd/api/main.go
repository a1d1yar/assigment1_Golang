package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/a1d1yar/assingment1_Golang/internal/data"
	"github.com/a1d1yar/assingment1_Golang/internal/jsonlog"
	"github.com/a1d1yar/assingment1_Golang/internal/mailer"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	db     *data.DBModel
	mailer mailer.Mailer
	models data.Models

	wg sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:Aldiyar2004@localhost:5432/a.maratovDB?sslmode=disable", "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "5246435a643a45", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "cb6b75fc700c77", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", " <221257@astanait.edu.kz>", "SMTP sender")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	// dbase, err := openDB(cfg)
	// if err != nil {
	// 	logger.PrintFatal(err, nil)
	// }
	// defer dbase.Close()

	logger.PrintInfo("database connection pool established", nil)

	dbModel := &data.DBModel{
		DB: db,
	}

	models := data.NewModels(db)
	// models := &data.Models{
	// 	Users: dbase,
	// }

	app := &application{
		config: cfg,
		logger: logger,
		db:     dbModel,
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		models: models,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})
	err = srv.ListenAndServe()
	logger.PrintFatal(err, nil)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	// db.SetMaxOpenConns(cfg.db.maxOpenConns)
	// // Set the maximum number of idle connections in the pool. Again, passing a value
	// // less than or equal to 0 will mean there is no limit.
	// db.SetMaxIdleConns(cfg.db.maxIdleConns)
	// // Use the time.ParseDuration() function to convert the idle timeout duration string
	// // to a time.Duration type.
	// duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	// if err != nil {
	// 	return nil, err
	// }
	// // Set the maximum idle timeout.
	// db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
