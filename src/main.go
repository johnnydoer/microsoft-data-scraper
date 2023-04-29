package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type versionData struct {
	version     string
	releaseDate time.Time
}

func getReleaseDates(versions []string) []versionData {
	fmt.Println("Ok")

	allData := []versionData{}

	// Put data in the struct

	return allData
}

// load the webpage in HTML format
// func fetchVersionPage() {

// }

// get all version numbers from the first page
// func getVersionNumbers () []string {

// }

// get first page to get all the version numbers
// getFirstPage() {

// }

func main() {
	// Replace with the URL of the website you want to scrape
	webPage := "https://www.microsoft.com/en-us/wdsi/definitions/antimalware-definition-release-notes"

	resp, err := http.Get(webPage)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}


	count := 0

	for {
		id := "#dropDownOption_" + strconv.Itoa(count)
		titleVersionText := doc.Find(id).Text()

		if len(titleVersionText) == 0 {
			break
		}

		fmt.Println(count, titleVersionText)
		count++

	}
}
