package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {

	printBanner()

	// Setup Command Line Arguments
	var customFlag TypeChoices
	customFlag.Options = []string{"sleep", "readiness"}
	flag.Var(&customFlag, "type", "Oura Ring data type to pull")

	startDateInput := flag.String("start-date", time.Now().Add(-24*time.Hour).Format("2006-01-02"), "Start date to search for Oura Ring data (defaults to yesterday")
	endDateInput := flag.String("end-date", time.Now().Format("2006-01-02"), "End date to search for Oura Ring data (defaults to today)")

	flag.Parse()

	if customFlag.Value == "" {
		log.Println("You must specify Oura Ring data type")
		os.Exit(1)
	}

	startDate, err := time.Parse("2006-01-02", *startDateInput)
	if err != nil {
		log.Fatalf("Failed to parse the start date: %v", err)
	}

	endDate, err := time.Parse("2006-01-02", *endDateInput)
	if err != nil {
		log.Fatalf("Failed to parse the end date: %v", err)
	}

	log.Println("Processing Oura Ring Type: " + customFlag.Value)
	log.Println("Start Date: " + startDate.Format("2006-01-02"))
	log.Println("End Date: " + endDate.Format("2006-01-02"))
	log.Println("--------------------------------------------")

	if customFlag.Value == "sleep" {
		ProcessSleep(startDate, endDate)
	}
}

func printBanner() {
	banner := `
 ___            _    _        ___  ___  ___
/   \ _  _  ___| |_ (_) _ _  /   \| _ \|_ _|
| - || || |(_-/|  _|| || ' \ | - ||  _/ | |
|_|_| \_._|/__/ \__||_||_||_||_|_||_|  |___|

AustinAPI Oura Ring Data Puller
============================================
`
	log.Print(banner)
}
