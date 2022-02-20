package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://192.168.1.120:3000"}
	r.Use(cors.New(config))
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetId := "181njrl_PeCETmTjiKM9H4AJP2l_4-e-TPS7bBoWgTPQ"
	readRange := "Fur!B2:B189"
	resp, err := srv.Spreadsheets.Get(spreadsheetId).Ranges(readRange).IncludeGridData(true).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	placed := make(map[string]int)
	for _, sht := range resp.Sheets {
		for _, row := range sht.Data {
			for p, cell := range row.RowData {
				for _, val := range cell.Values {
					placed[val.FormattedValue] = p + 1
				}
			}
		}
	}

	keys := make([]string, len(placed))

	i := 0

	for k := range placed {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	r.GET("/api/furs", func(c *gin.Context) {
		c.JSON(200, placed)
	})

	r.GET("/api/furs/names", func(c *gin.Context) {
		c.JSON(200, keys)
	})

	r.GET("/api/furs/:name", func(c *gin.Context) {
		c.JSON(200, placed[c.Param("name")])
	})

	r.GET("/api/furs/compare/:first/:second", func(c *gin.Context) {
		if placed[c.Param("first")] == placed[c.Param("second")] {
			c.String(400, "Both furs are the same. Please change your selection and try again.")
		} else if placed[c.Param("first")] < placed[c.Param("second")] {
			c.String(200, fmt.Sprintf("%s is dominant to %s", c.Param("first"), c.Param("second")))
		} else {
			c.String(200, fmt.Sprintf("%s is recessive to %s", c.Param("first"), c.Param("second")))
		}
	})

	r.Run()
}
