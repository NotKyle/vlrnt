package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"notkyle.org/vlrnt/db"
	"notkyle.org/vlrnt/scraper"
	"notkyle.org/vlrnt/utils"
)

func main() {
	envFile, err := godotenv.Read(".env")

	if err != nil {
		panic(err)
	}

	url := envFile["SCRAPE_URL"]

	if url == "" {
		panic(errors.New("No URL provided"))
	}

	_, err = os.Stat("db.sqlite")

	if err != nil {
		db.CreateDB()
	}

	scraper.Scrape(url)

	DoMatches()
}

func DoMatches() {
	database, err := db.Open()

	if err != nil {
		panic(err)
	}

	// fmt.Println("Getting matches")
	matches, err := db.GetMatches(database)

	if matches == nil || len(matches) == 0 {
		panic("No matches found")
	}

	if err != nil {
		panic(err)
	}

	// Convert matches to JSON array
	for _, match := range matches {
		var asJson string = ""

		// asJson = fmt.Sprintf(
		// 	"{\"id\": %d, \"url\": \"%s\", \"team1\": \"%s\", \"team2\": \"%s\", \"startTime\": \"%s\", \"endTime\": \"%s\"}",
		// 	match.ID, match.URL, match.Team1.Name, match.Team2.Name, match.StartTime, match.EndTime,
		// )

		asJson, err = utils.AsJson(match)

		fmt.Println(asJson)
	}
}
