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

// Load the webpage in HTML format
func fetchVersionPage(webPage string) *goquery.Document {
	// Request the webpage and get the response.
	resp, err := http.Get(webPage)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Check if the response had a 200 status code or not.
	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}

	// Get the HTML data from the response body.
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc

}

func getValidData(allVersionData []versionData) []versionData {
	numberOfDays := 10
	validData := []versionData{}
	now := time.Now()

	for _, data := range allVersionData {
		dayDifference := now.Sub(data.releaseDate).Hours() / 24

		if numberOfDays > int(dayDifference) {
			validData = append(validData, data)
		}
	}

	return validData
}

// Get all version numbers from the first page
// func getVersionNumbers () []string {

// }

// Get first page to get all the version numbers
// getFirstPage() {

// }

func main() {
	// URL to get data from.
	webPage := "https://www.microsoft.com/en-us/wdsi/definitions/antimalware-definition-release-notes"

	// Request the webpage and get the response.
	resp, err := http.Get(webPage)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Check if the response had a 200 status code or not.
	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}

	// Get the HTML data from the response body.
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Set a counter.
	count := 0
	versions := []string{}
	// Loop over the versions.
	for {
		id := "#dropDownOption_" + strconv.Itoa(count)
		titleVersionText := doc.Find(id).Text()

		if len(titleVersionText) == 0 {
			break
		}

		versions = append(versions, titleVersionText)
		// fmt.Println(count, titleVersionText)

		count++
	}

	allVersionData := []versionData{}
	for _, version := range versions {
		webPage := "https://www.microsoft.com/en-us/wdsi/definitions/antimalware-definition-release-notes?requestVersion=" + version
		versionDoc := fetchVersionPage(webPage)

		id := "#releaseDate_0"
		releaseDate := versionDoc.Find(id).Text()

		data := versionData{}

		data.version = version
		data.releaseDate, _ = time.Parse("1/2/2006 3:04:05 PM", releaseDate)

		allVersionData = append(allVersionData, data)
	}

	validData := getValidData(allVersionData)

	for _, data := range validData {
		fmt.Println(data.version, data.releaseDate)
	}

}
