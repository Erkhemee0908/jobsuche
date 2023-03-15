package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
)

// This function retrieves job data from a specific URL and writes the response body to a file called "jobs.json".
func getJobs(jobNum ...int) {

	if len(jobNum) == 0 {
		jobNum = append(jobNum, 3)
	}

	// The URL to send a GET request to
	// url := "https://rest.arbeitsagentur.de/jobboerse/jobsuche-service/pc/v4/jobs?size=1&angebotsart=4&was=Fachinformatiker%2Fin%20Anwendungsentwicklung&wo=Berlin&umkreis=10"

	numberOfJobs := "size=" + strconv.Itoa(jobNum[0])

	url := "https://rest.arbeitsagentur.de/jobboerse/jobsuche-service/pc/v4/jobs?"
	url2 := "&angebotsart=4&was=Fachinformatiker%2Fin%20Anwendungsentwicklung&wo=Berlin&umkreis=10"

	url = url + numberOfJobs + url2

	// Create a new HTTP request with the given URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Create an HTTP client to send the request
	client := &http.Client{}

	// Send the HTTP request and get the response
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body into a byte slice
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Create a new file called "jobs.json"
	file, err := os.Create("jobs.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write the response body to the file
	_, err = file.Write(body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

}

// This function takes a file path as input, reads in the contents of the file, and processes the job data contained in it.
// It returns a slice of Job structs and an error, if any occurred.
func processJobs(file string) ([]Job, error) {

	// Define a LocalJob struct with fields for the job title, reference number, and employer.
	type LocalJob struct {
		Titel       string `json:"titel"`
		Refnum      string `json:"refnr"`
		Arbeitgeber string `json:"arbeitgeber"`
	}

	// Define a jobs struct with a field for a slice of LocalJob structs.
	var jobs struct {
		Stellenangebote []LocalJob `json:"stellenangebote"`
	}

	// Read in the contents of the file at the given file path.
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into the jobs struct.
	err = json.Unmarshal(content, &jobs)
	if err != nil {
		return nil, err
	}

	// Create an empty slice to hold Job structs.
	var jobList []Job

	// Iterate over each LocalJob in the slice of LocalJob structs.
	for _, job := range jobs.Stellenangebote {
		sleepRandomDuration()

		// Call the getJobDetail function to retrieve the job detail.
		jobDetail, err := getJobDetail(job.Refnum)
		if err != nil {
			return nil, err
		}

		// Construct a new Job struct with the relevant job data and the retrieved job detail.
		Job := Job{
			Titel:       job.Titel,
			Refnum:      job.Refnum,
			Arbeitgeber: job.Arbeitgeber,
			Description: jobDetail,
		}

		// Append the new Job struct to the slice of Job structs.
		jobList = append(jobList, Job)
	}
	// Return the slice of Job structs and a nil error.

	return jobList, nil
}

// This function takes a slice of Job structs and a filename as input and writes the job data to a file in JSON format.
// It returns an error, if any occurred.

// The function marshals the slice of Job structs into JSON format with indentation.
// It then writes the JSON data to a file with the given filename and file permissions.
// If there are no errors, the function returns a nil error.
// func writeJobsToFile(jobs []Job, filename string) error {

// 	// Marshal the slice of Job structs into JSON format with indentation.
// 	result, err := json.MarshalIndent(jobs, "", "    ")
// 	if err != nil {
// 		return err
// 	}

// 	// Write the JSON data to a file with the given filename and file permissions.
// 	err = os.WriteFile(filename, result, 0644)
// 	if err != nil {
// 		return err
// 	}

// 	// If there are no errors, return a nil error.
// 	return nil
// }

func getJobDetail(refnum string) (string, error) {
	url := "https://www.arbeitsagentur.de/jobsuche/jobdetail/" + refnum

	// Create a new Chrome instance.
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout of 5 seconds for the entire process.
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Navigate to the URL you want to scrape.
	err := chromedp.Run(ctx, chromedp.Navigate(url))
	if err != nil {
		return "", err
	}

	// Wait for the element with refnum "jobdetails-beschreibung" to be visible.
	err = chromedp.Run(ctx, chromedp.WaitVisible("#jobdetails-titel"))
	if err != nil {
		return "", err
	}

	var html string

	// Wait for the element with refnum "jobdetails-beschreibung" to be visible.
	err = chromedp.Run(ctx, chromedp.WaitVisible("#jobdetails-beschreibung"))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {

			return "N/A", nil
		}
		return "", err
	}

	// Get the content of the element with refnum "jobdetails-beschreibung".
	err = chromedp.Run(ctx, chromedp.InnerHTML("#jobdetails-beschreibung", &html, chromedp.ByID))
	if err != nil {
		return "", err
	}

	// Return the scraped HTML as a string.
	return html, nil
}

func sleepRandomDuration() {
	// Generate a random duration between 1 and 10 seconds
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	duration := time.Duration(r.Intn(100)+1) * time.Millisecond

	// Sleep for the random duration
	time.Sleep(duration)
}
