package main

import (
	"errors"

	"github.com/joho/godotenv"

	"notkyle.org/vlrnt/scraper"
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

	scraper.Scrape(url)
}
