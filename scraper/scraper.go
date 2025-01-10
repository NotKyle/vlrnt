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
	"notkyle.org/vlrnt/utils"
)

type Match struct {
	Url, Team1, Team2 string
}

func Scrape(url string) error {
	if url == "" {
		return errors.New("url cannot be empty")
	}

	var matches []Match

	fmt.Println("Scraping", url)

	// Remove path
	path, _ := utils.GetPath(url)
	urlnopath := strings.ReplaceAll(url, path, "")

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

	fmt.Println("Allowed domains:", urlnopath)

	fmt.Println("Final URL:", urlnopath)

	fmt.Println("c", c)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// triggered when the scraper encounters an error
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// fired when the server responds
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	// triggered when a CSS selector matches an element
	c.OnHTML("a.match-item", func(e *colly.HTMLElement) {
		match := Match{}

		match.Url = e.Attr("href")

		goquerySelection := e.DOM
		teams := goquerySelection.Find("div.match-item-vs>.match-item-vs-team")

		// Example Goquery usage

		teams.Each(func(i int, s *goquery.Selection) {
			teamName := s.Find(".match-item-vs-team-name>.text-of").Text()

			fmt.Println("Team name:", teamName)
			fmt.Println("Match URL:", match.Url)
			fmt.Println("Match number:", i)

			if i == 0 {
				match.Team1 = teamName
			} else {
				match.Team2 = teamName
			}
		})

		// for _, team := range teams {
		// 	teamName := goquerySelection.Find(".match-item-vs-team-name>.text-of").Text()
		// 	fmt.Println("Team name:", teamName)
		// }

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
			"Url",
			"Team1",
			"Team2",
		}

		writer.Write(headers)

		for _, match := range matches {
			record := []string{
				match.Url,
				match.Team1,
				match.Team2,
			}

			writer.Write(record)
		}

		defer writer.Flush()
	})

	fmt.Println("Visiting", url)

	c.Visit(url)

	fmt.Println("Scraping complete")

	return nil
}
