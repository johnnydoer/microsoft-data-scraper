package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type versionData struct {
	version     string
	releaseDate time.Time
}

var wg = sync.WaitGroup{}

// Load the webpage in HTML format.
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
func getVersionNumbers(doc *goquery.Document) []string {
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

	return versions
}

// Get first page to get all the version numbers
func getFirstPage() *goquery.Document {
	// URL to get data from.
	webPage := "https://www.microsoft.com/en-us/wdsi/definitions/antimalware-definition-release-notes"

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

func extractReleaseDate(versionDoc *goquery.Document) time.Time {
	id := "#releaseDate_0"
	releaseDateStr := versionDoc.Find(id).Text()

	releaseDate, _ := time.Parse("1/2/2006 3:04:05 PM", releaseDateStr)

	return releaseDate
}

func getReleaseDates(version string, c chan versionData) {
	webPage := "https://www.microsoft.com/en-us/wdsi/definitions/antimalware-definition-release-notes?requestVersion=" + version

	// Function to fetch the data from URL.
	doc := fetchVersionPage(webPage)

	// Extract the release date in time format from the document.
	releaseDate := extractReleaseDate(doc)

	var data versionData

	data.releaseDate = releaseDate
	data.version = version

	// fmt.Println("Before") // DEBUG
	c <- data
	// fmt.Println("After") // DEBUG

	wg.Done()

	// fmt.Println("GoRoutine for version:", version) // DEBUG
}

func main() {
	// Request the webpage and get the response.
	doc := getFirstPage()

	// Extract all possible version numbers from the first page.
	versions := getVersionNumbers(doc)
	fmt.Println("First Page done.")

	//
	allVersionData := []versionData{}

	// Channel to communicate all version data in.
	versionDatach := make(chan versionData, len(versions))

	// Loop over the versions from 1st page to get the release dates and save all the data.
	for _, version := range versions {
		wg.Add(1)
		go getReleaseDates(version, versionDatach)
	}
	// fmt.Println("Versions loop done.") // DEBUG

	wg.Wait()
	close(versionDatach)

	// fmt.Println("Waiting done.") // DEBUG

	for vData := range versionDatach {
		allVersionData = append(allVersionData, vData)
	}

	validData := getValidData(allVersionData)

	for _, data := range validData {
		fmt.Println(data.version, data.releaseDate)
	}
}

