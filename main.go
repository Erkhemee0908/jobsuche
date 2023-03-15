package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
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

	http.ListenAndServe(":8080", nil)
}
