package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

type config struct {
	port int
	env  string
}

type TraitAttributes struct {
	FirstRelease bool `redis:"first_release"`
	Retired      bool `redis:"retired"`
}

type TraitData struct {
	Name       string          `redis:"name"`
	Position   int             `redis:position`
	Attributes TraitAttributes `redis:"attributes"`
}

type application struct {
	pool   *redis.Pool
	config config
	logger *log.Logger
}

func (s *application) getData(command, hash, key string) ([]string, error) {
	conn := s.pool.Get()
	defer conn.Close()
	var values []string
	var err error
	switch command {
	case "names":
		values, err = redis.Strings(conn.Do("HKEYS", hash))
		break
	case "values":
		values, err = redis.Strings(conn.Do("HGETALL", hash))
		break
	case "single":
		var v string
		v, err = redis.String(conn.Do("HGET", hash, key))
		if err != nil {
			values = append(values, err.Error())
		}
		values = append(values, v)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
		return make([]string, 0), err
	}
	return values, nil
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
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

	app := &application{
		config: cfg,
		logger: logger,
		pool: &redis.Pool{
			MaxIdle:     10,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.DialURL(os.Getenv("REDIS_TLS_URL"), redis.DialTLSSkipVerify(true))
			},
		},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
