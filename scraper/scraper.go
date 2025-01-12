package scraper

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"notkyle.org/vlrnt/db"

	"notkyle.org/vlrnt/structs"
)

func Scrape(url string) error {
	if url == "" {
		return errors.New("url cannot be empty")
	}

	var matches []structs.Match

	// Remove path
	// path, _ := utils.GetPath(url)
	// urlnopath := strings.ReplaceAll(url, path, "")

	// Add URL to allowed domains
	c := colly.NewCollector(
		colly.AllowedDomains(
			"https://www.vlr.gg",
			"www.vlr.gg",
			"vlr.gg",
			"https://www.vlr.gg/",
			"www.vlr.gg/",
			"vlr.gg/",
		),
	)

	c.OnRequest(func(r *colly.Request) {
	})

	// triggered when the scraper encounters an error
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// fired when the server responds
	c.OnResponse(func(r *colly.Response) {
	})

	// triggered when a CSS selector matches an element
	c.OnHTML("a.match-item", func(e *colly.HTMLElement) {
		match := structs.Match{}

		match.URL = e.Attr("href")

		goquerySelection := e.DOM
		teams := goquerySelection.Find("div.match-item-vs>.match-item-vs-team")

		// Example Goquery usage

		teams.Each(func(i int, s *goquery.Selection) {
			teamName := s.Find(".match-item-vs-team-name>.text-of").Text()

			// fmt.Println("Team name:", teamName)
			// fmt.Println("Match URL:", match.URL)
			// fmt.Println("Match number:", i)

			teamName = strings.TrimSpace(teamName)

			team := structs.Team{
				ID:      i,
				Name:    teamName,
				Players: []structs.Player{},
			}

			if i == 0 {
				match.Team1 = team
			} else {
				match.Team2 = team
			}
		})

		matches = append(matches, match)
	})

	// triggered once scraping is done (e.g., write the data to a CSV file)
	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("matches.csv")

		if err != nil {
			log.Fatal("Cannot create file", err)
		}

		defer file.Close()

		writer := csv.NewWriter(file)

		headers := []string{
			"URL",
			"Team1",
			"Team2",
		}

		writer.Write(headers)

		databse, err := db.Open()

		if err != nil {
			log.Fatal(err)
		}

		for _, match := range matches {
			record := []string{
				match.URL,
				match.Team1.Name,
				match.Team2.Name,
			}

			// Add to database
			// fmt.Println("Adding match to database", match)
			db.AddMatch(databse, match)

			writer.Write(record)
		}

		defer writer.Flush()
	})

	// fmt.Println("Visiting", url)

	c.Visit(url)

	// fmt.Println("Scraping complete")

	return nil
}
