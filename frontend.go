package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	getJobs(3)
	fmt.Println("Got Jobs")

	jobList, err := processJobs("jobs.json")
	if err != nil {
		panic(err)
	}

	err = writeJobsToFile(jobList, "jobList.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("Number of jobs processed:", len(jobList))

	fmt.Println("Starting application...")
	myApp := app.New()
	myWindow := myApp.NewWindow("Hello, World!")
	myLabel := widget.NewLabel("Number of jobs processed: " + strconv.Itoa(len(jobList)))
	myContainer := container.New(layout.NewVBoxLayout(), myLabel)
	myWindow.SetContent(myContainer)
	myWindow.ShowAndRun()
}
