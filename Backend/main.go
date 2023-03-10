package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	getJobs()

	jobList, err := processJobs("jobs.json")
	if err != nil {
		panic(err)
	}

	err = writeJobsToFile(jobList, "jobList.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of jobs processed:", len(jobList))

}

type Job struct {
	Titel       string `json:"titel"`
	Refnum      string `json:"refnr"`
	Arbeitgeber string `json:"arbeitgeber"`
	Description string `json:"description"`
}

// TODO: Allow getJobs function to recieve a parameter for the number of jobs to retrieve, make default 2
func getJobs() {

	url := "https://rest.arbeitsagentur.de/jobboerse/jobsuche-service/pc/v4/jobs?size=2&angebotsart=4&was=Fachinformatiker%2Fin%20Anwendungsentwicklung&wo=Berlin&umkreis=10"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Write the response body to a file
	file, err := os.Create("jobs.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

}

func processJobs(file string) ([]Job, error) {
	// Define a LocalJob struct
	type LocalJob struct {
		Titel       string `json:"titel"`
		Refnum      string `json:"refnr"`
		Arbeitgeber string `json:"arbeitgeber"`
	}

	var jobs struct {
		Stellenangebote []LocalJob `json:"stellenangebote"`
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &jobs)
	if err != nil {
		return nil, err
	}

	var jobList []Job

	for _, job := range jobs.Stellenangebote {
		start := time.Now()
		sleepRandomDuration()
		jobDetail, err := getJobDetail(job.Refnum)
		if err != nil {
			return nil, err
		}
		duration := time.Since(start)

		fmt.Println(duration)

		Job := Job{
			Titel:       job.Titel,
			Refnum:      job.Refnum,
			Arbeitgeber: job.Arbeitgeber,
			Description: jobDetail,
		}
		jobList = append(jobList, Job)
	}

	return jobList, nil
}

func writeJobsToFile(jobs []Job, filename string) error {
	result, err := json.MarshalIndent(jobs, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, result, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getJobDetail(refnum string) (string, error) {
	url := "https://www.arbeitsagentur.de/jobsuche/jobdetail/" + refnum

	// Create a new Chrome instance.
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Navigate to the URL you want to scrape.
	err := chromedp.Run(ctx, chromedp.Navigate(url))
	if err != nil {
		return "", err
	}

	// Wait for the element with refnum "jobdetails-beschreibung" to be visible.
	err = chromedp.Run(ctx, chromedp.WaitVisible("#jobdetails-beschreibung"))
	if err != nil {
		return "", err
	}

	// Get the content of the element with refnum "jobdetails-beschreibung".
	var html string
	err = chromedp.Run(ctx, chromedp.InnerHTML("#jobdetails-beschreibung", &html, chromedp.ByID))
	if err != nil {
		return "", err
	}

	// Return the scraped HTML as a string.
	return html, nil
}

func sleepRandomDuration() {
	// Generate a random duration between 1 and 10 seconds
	rand.Seed(time.Now().UnixNano())
	duration := time.Duration(rand.Intn(3)+1) * time.Second

	// Sleep for the random duration
	time.Sleep(duration)
}
