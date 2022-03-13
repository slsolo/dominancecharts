package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

type server struct {
	router *gin.Engine
	pool   *redis.Pool
}

func (s *server) getData(command, hash, key string) ([]string, error) {
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
		v, _ = redis.String(conn.Do("HGET", hash, key))
		values = append(values, v)
	}
	if err != nil {
		fmt.Printf("%v\n", err)
		return make([]string, 0), err
	}
	return values, nil
}

func (s *server) routes() {
	s.router.GET("/api/:trait/names", func(c *gin.Context) {
		values, err := s.getData("names", strings.Title(c.Param("trait")), "")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if len(values) == 0 {
			c.String(http.StatusNotFound, fmt.Sprintf("No %ss found.", c.Param("trait")))
			return
		}
		sort.Strings(values)
		c.JSON(http.StatusOK, values)
	})
	s.router.GET("/api/:trait", func(c *gin.Context) {
		placed := make(map[string]int)
		values, err := s.getData("values", strings.Title(c.Param("trait")), "")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if len(values) == 0 {
			c.String(http.StatusNotFound, fmt.Sprintf("no %ss found", c.Param("trait")))
			return
		}

		respLen := len(values)
		for i := 0; i < respLen; i += 2 {
			placed[values[i]], _ = strconv.Atoi(values[i+1])
		}
		c.JSON(http.StatusOK, placed)
	})

	s.router.GET("/api/:trait/:name", func(c *gin.Context) {
		values, err := s.getData("names", strings.Title(c.Param("trait")), c.Param("name"))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if len(values) == 0 {
			c.String(http.StatusNotFound, "%s %s not found\n", c.Param("Trait"), c.Param("name"))
			return
		}
		var i int
		i, err = strconv.Atoi(values[0])
		if err != nil {
			c.String(http.StatusInternalServerError, "Value is not an integer")
			return
		}
		c.JSON(200, i)
	})

	s.router.GET("/api/:trait/compare/:first/:second", func(c *gin.Context) {
		trait := strings.Title(c.Param("trait"))
		first := c.Param("first")
		second := c.Param("second")
		var fi, si int
		values, err := s.getData("single", trait, first)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if len(values) == 0 {
			c.String(http.StatusNotFound, "%s %s not found\n", trait, first)
			return
		}

		fi, err = strconv.Atoi(values[0])
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		values, err = s.getData("single", trait, second)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if len(values) == 0 {
			c.String(http.StatusNotFound, "%s %s not found\n", trait, second)
			return
		}

		si, err = strconv.Atoi(values[0])
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if first == second {
			c.String(400, "Both %ss are the same. Please change your selection and try again.", trait)
		} else if fi < si {
			c.String(200, fmt.Sprintf("%s is dominant to %s", first, second))
		} else {
			c.String(200, fmt.Sprintf("%s is recessive to %s", first, second))
		}
	})
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("In production, fetching config from Heroku config parameters...")
	}
	srv := server{
		router: gin.Default(),
		pool: &redis.Pool{
			MaxIdle:     10,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", os.Getenv("REDIS_TLS_URL"))
			},
		},
	}
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	srv.router.Use(cors.New(config))
	srv.routes()

	srv.router.Run()
}
