package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"

	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

var pool *redis.Pool

func fetch(ctx context.Context, done chan bool, srv *sheets.Service, spreadsheet, sheet, column string, row int) {
	conn := pool.Get()
	defer conn.Close()

	readRange := fmt.Sprintf("%s!%s%d:%s", sheet, column, row, column)
	resp, err := srv.Spreadsheets.Get(spreadsheet).Ranges(readRange).IncludeGridData(true).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	endofrange := 0
	for _, sht := range resp.Sheets {
		for _, row := range sht.Data {
			for p, cell := range row.RowData {
				for _, val := range cell.Values {
					if val.FormattedValue == "" {
						if endofrange == 0 {
							endofrange = p + 1
							break
						}
					}
				}
			}
		}
	}

	placedRange := fmt.Sprintf("%s!%s%d:%s%d", sheet, column, row, column, endofrange)
	log.Printf("Placed range %s\n", placedRange)
	resp, _ = srv.Spreadsheets.Get(spreadsheet).Ranges(placedRange).IncludeGridData(true).Do()

	for _, sht := range resp.Sheets {
		for _, row := range sht.Data {
			for p, cell := range row.RowData {
				for _, val := range cell.Values {
					fmt.Printf("%s: %d", val.FormattedValue, p+1)
					_, err = conn.Do("HSET", sheet, val.FormattedValue, p+1)
					if err != nil {
						log.Fatalf("Error setting redis value %v\n", err)
					}
				}
			}
		}
	}
	done <- true
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("In production, fetching config from Heroku config parameters...")
	}

	ctx := context.Background()
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS_TLS_URL"))
		},
	}

	srv, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	furs := make(chan bool, 1)
	eyes := make(chan bool, 1)
	tails := make(chan bool, 1)
	ears := make(chan bool, 1)
	whiskerColour := make(chan bool, 1)
	whiskerShape := make(chan bool, 1)
	shades := make(chan bool, 1)
	spreadsheetId := os.Getenv("SPREADSHEET")
	fetch(ctx, furs, srv, spreadsheetId, "Fur", "B", 2)
	fetch(ctx, eyes, srv, spreadsheetId, "Eyes", "B", 2)
	fetch(ctx, tails, srv, spreadsheetId, "Tails", "A", 2)
	fetch(ctx, ears, srv, spreadsheetId, "Ears", "A", 2)
	fetch(ctx, whiskerColour, srv, spreadsheetId, "Whisker Colour", "A", 2)
	fetch(ctx, whiskerShape, srv, spreadsheetId, "Whisker Shape", "A", 2)
	fetch(ctx, shades, srv, spreadsheetId, "Other", "F", 3)
	<-furs
	<-eyes
	<-tails
	<-ears
	<-whiskerColour
	<-whiskerShape
	<-shades
}
