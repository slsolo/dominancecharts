package data

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Trait struct {
	Name       string `json:"name"`
	Position   int    `json:"position"`
	Initial    bool   `json:"initial_release"`
	Retired    bool   `json:"retired"`
	Limited    bool   `json:"limited_edition"`
	Unplaced   bool   `json:"unplaced"`
	RangeStart int    `json:"range_start"`
	RangeEnd   int    `json:"range_end"`
}
type DocData map[string]map[string]Trait

type TraitModels struct {
	Data DocData
	Lock sync.RWMutex
}

func NewTraitsFromGDoc() (*TraitModels, error) {
	sheetData := make(DocData)
	var tm TraitModels
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		return &TraitModels{}, err
	}
	spreadsheetID := os.Getenv("SPREADSHEET")
	ranges := []string{"Fur!B2:B", "Eyes!B2:B", "Tails!A2:A", "Ears!A2:A", "Whisker Colour!A2:A", "Whisker Shape!A2:A", "Other!F3:F"}
	resp, err := srv.Spreadsheets.Get(spreadsheetID).Ranges(ranges...).IncludeGridData(true).Do()
	if err != nil {
		return &TraitModels{}, err
	}
	count_furs := 0
	count_eyes := 0
	count_tails := 0
	count_ears := 0
	count_whisker_colours := 0
	count_whisker_shapes := 0
	count_shades := 0
	for _, sht := range resp.Sheets {
		fmt.Printf("\n=======PROCESSING %s==============\n", sht.Properties.Title)
		sheetProcessed := false
		vals := make(map[string]Trait)
		for _, row := range sht.Data {
			for p, cell := range row.RowData {
				for _, val := range cell.Values {
					key := val.FormattedValue
					if sheetProcessed {
						continue
					}
					if key == "" {
						sheetProcessed = true
						continue
					}
					if sht.Properties.Title == "Fur" {
						count_furs++
					}
					if sht.Properties.Title == "Eyes" {
						count_eyes++
					}
					if sht.Properties.Title == "Tails" {
						count_tails++
					}
					if sht.Properties.Title == "Ears" {
						count_ears++
					}
					if sht.Properties.Title == "Whisker Colour" {
						count_whisker_colours++
					}
					if sht.Properties.Title == "Whisker Shape" {
						count_whisker_shapes++
					}
					if sht.Properties.Title == "Other" {
						count_shades++
					}
					var isRetired bool
					var isInitialRelease bool
					var isLimitedEdition bool
					var isUnplaced bool
					var trait Trait
					if strings.Contains(key, "*") {
						isInitialRelease = true
						key = strings.ReplaceAll(key, "*", "")
					}
					if strings.Contains(key, "(Retired)") {
						isRetired = true
						key = strings.ReplaceAll(key, "(Retired)", "")
					}
					if strings.Contains(key, "(Ltd. Release)") {
						isLimitedEdition = true
						key = strings.ReplaceAll(key, "(Ltd. Release)", "")
					}

					key = strings.TrimSpace(key)
					trait.Name = key
					trait.Position = p + 1
					trait.Initial = isInitialRelease
					trait.Limited = isLimitedEdition
					trait.Retired = isRetired

					fmt.Printf("%s: %v", key, trait)
					vals[key] = trait
				}
			}
		}
		sheetData[sht.Properties.Title] = vals
	}
	tm.Data = sheetData
	return &tm, nil
}
