package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/slsolo/dominancecharts/internal/data"
)

type config struct {
	port int64
	env  string
}

type application struct {
	config   config
	errorLog *log.Logger
	infoLog  *log.Logger
	models   *data.TraitModels
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)
	err := godotenv.Load(".env")

	if err != nil {
		errorLog.Println(err)
		errorLog.Println("In production, fetching config from Heroku config parameters...")
	}

	// Declare an instance of the config struct.
	var cfg config
	val, err := strconv.ParseInt(os.Getenv("PORT"), 10, 32)
	cfg.port = val
	cfg.env = os.Getenv("ENVIRONMENT")

	traits, err := data.NewTraitsFromGDoc()
	if err != nil {
		errorLog.Fatalf("Error fetching data from Dominance Charts: %v\n", err)
	}
	fmt.Println(*&traits.Data)
	app := &application{
		config:   cfg,
		errorLog: errorLog,
		infoLog:  infoLog,
		models:   traits,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	infoLog.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
