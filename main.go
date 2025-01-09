package main

import (
	"errors"

	"notkyle.org/vlrnt/scraper"
)

func main() {
	url := ""

	if url == "" {
		panic(errors.New("No URL provided"))
	}

	scraper.Scrape(url)
}
