package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
	// getJobs(3)
	// fmt.Println("Got Jobs")

	// jobList, err := processJobs("jobs.json")
	// if err != nil {
	// 	panic(err)
	// }

	// err = writeJobsToFile(jobList, "jobList.json")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Number of jobs processed:", len(jobList))

}
