package scraper

import (
	"errors"
	"fmt"
	"notkyle.org/vlrnt/utils"
)

func Scrape(url string) error {
	if url == "" {
		return errors.New("url cannot be empty")
	}

	fmt.Println("Scraping", url)

	// urlParts := strings.Split(url, "/")
	method, err := utils.GetMethod(url)

	if err != nil {
		return err
	}

	domain, err := utils.GetDomain(url)
	if err != nil {
		return err
	}

	tld, err := utils.GetTLD(url)
	if err != nil {
		return err
	}

	path, err := utils.GetPath(url)
	if err != nil {
		return err
	}

	fmt.Println("Domain:", domain)
	fmt.Println("TLD:", tld)
	fmt.Println("Path:", path)
	fmt.Println("Method:", method)

	subdomain, err := utils.GetSubDomain(url)
	if err != nil {
		return err
	}

	fmt.Println("Subodmain:", subdomain)

	return nil
}
