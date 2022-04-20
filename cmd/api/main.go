package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/slsolo/dominancecharts/internal/data"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
	models *data.TraitModels
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
		log.Println("In production, fetching config from Heroku config parameters...")
	}

	// Declare an instance of the config struct.
	var cfg config

	// Read the value of the port and env command-line flags into the config struct. We
	// default to using the port number 4000 and the environment "development" if no
	// corresponding flags are provided.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	traits, err := data.NewTraitsFromGDoc()
	if err != nil {
		log.Fatalf("Error fetching data from Dominance Charts: %v\n", err)
	}
	fmt.Println(*&traits.Data)
	app := &application{
		config: cfg,
		logger: logger,
		models: traits,
	}

	handler := cors.Default().Handler(app.Routes())

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
