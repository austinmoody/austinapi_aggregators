package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {

	// Pull Oura token from environmentd
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ouraPersonalToken := os.Getenv("OURA_PERSONAL_TOKEN")

	// Setup Command Line Arguments
	var customFlag TypeChoices
	customFlag.Options = []string{"sleep", "readiness"}
	flag.Var(&customFlag, "type", "Oura Ring data type to pull")

	startDateInput := flag.String("start-date", time.Now().Add(-24*time.Hour).Format("2006-01-02"), "Start date to search for Oura Ring data (defaults to yesterday")
	endDateInput := flag.String("end-date", time.Now().Format("2006-01-02"), "End date to search for Oura Ring data (defaults to today)")

	flag.Parse()

	startDate, err := time.Parse("2006-01-02", *startDateInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse the start date: %v", err)
		os.Exit(1)
	}

	endDate, err := time.Parse("2006-01-02", *endDateInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse the end date: %v", err)
		os.Exit(1)
	}

	fmt.Println("Type Specified: " + customFlag.Value)
	fmt.Println("Start Date: " + startDate.Format("2006-01-02"))
	fmt.Println("End Date: " + endDate.Format("2006-01-02"))

	if customFlag.Value == "sleep" {
		ProcessSleep(ouraPersonalToken, startDate, endDate)
	}
}
