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

var validData = []versionData{}

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

// Get all version numbers from the first page
func getVersionNumbers(doc *goquery.Document) []string {
	// Set a counter.
	count := 0
	versions := []string{}

	// Loop over the versions.
	for {
		id := "#dropDownOption_" + strconv.Itoa(count)
		titleVersionText := doc.Find(id).Text()

		// titleVersionText is empty then page was not found.
		if len(titleVersionText) == 0 {
			break
		}

		versions = append(versions, titleVersionText)
		// fmt.Println(count, titleVersionText) //DEBUG

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

	// Parse string to time format.
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

func getValidData(c chan versionData) {
	numberOfDays := 10

	now := time.Now()

	data := <-c
	dayDifference := now.Sub(data.releaseDate).Hours() / 24

	if numberOfDays > int(dayDifference) {
		validData = append(validData, data)
	}

	wg.Done()
}

func main() {
	// Request the webpage and get the response.
	doc := getFirstPage()

	// Extract all possible version numbers from the first page.
	versions := getVersionNumbers(doc)
	fmt.Println("First Page done.")

	// Channel to communicate all version data in.
	// Note unbuffered channel works here because we are also reading concurrently.
	// If the had first putting all data into a channel and then reading it then we require buffered channel to hold all the inputs.
	versionDatach := make(chan versionData)

	// Loop over the versions from 1st page to get the release dates and save all the data.
	for _, version := range versions {
		wg.Add(2)
		go getReleaseDates(version, versionDatach)
		go getValidData(versionDatach)

	}
	// fmt.Println("Versions loop done.") // DEBUG

	wg.Wait()
	close(versionDatach)

	for _, data := range validData {
		fmt.Println(data.version, data.releaseDate)
	}

	// TODO: Save data to file or send to other service or push to PubSub queue.
}
