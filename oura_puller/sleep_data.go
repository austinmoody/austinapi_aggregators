package main

import (
	"fmt"
	"github.com/austinmoody/go_oura"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func ProcessSleep(startDate time.Time, endDate time.Time) {

	// TODO: Something to think about... specifying a start & end a day apart.
	// so Start = 2024-01-29 & End = 2024-01-30
	// Sleep will bring back 1 item dated 2024-01-29
	// Daily Sleep will bring back 1 items: 2024-01-29 and 2024-01-30

	// Pull Oura token from environment
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ouraPersonalToken := os.Getenv("OURA_PERSONAL_TOKEN")

	// Here we would use oura_go to pull sleep data and parse out what we need
	client := go_oura.NewClient(ouraPersonalToken)

	// Get Sleep from API, has information about durations & start/end time etc...
	fmt.Printf("Pulling Sleep Items\n")
	sleepDocs, err := client.GetSleeps(startDate, endDate, nil)
	if err != nil {
		fmt.Printf("Error getting Sleep Items: %v\n", err)
		return
	}

	for _, sleepDoc := range sleepDocs.Items {
		InsertSleepData(sleepDoc.Day.Format("2006-01-02"), 0, sleepDoc.TotalSleepDuration)
		fmt.Printf("\tSleep Item %s Inserted\n", sleepDoc.ID)
	}

	// Get Daily Sleep from API, has the Score
	fmt.Printf("Pulling Daily Sleep Items\n")
	dailySleeps, err := client.GetDailySleeps(startDate, endDate, nil)
	if err != nil {
		fmt.Printf("Error getting Daily Sleep items: %v\n", err)
	}

	for _, dailySleep := range dailySleeps.Items {
		InsertDailySleepData(dailySleep.Day.Format("2006-01-02"), dailySleep.Score)
		fmt.Printf("\tDaily Sleep Item %s Inserted\n", dailySleep.ID)
	}
}

func InsertSleepData(date string, rating int, totalDuration int) {
	InsertData("INSERT INTO sleep (date, rating, total_duration) VALUES ($1, $2, $3) ON CONFLICT (date) DO UPDATE SET rating = EXCLUDED.rating;", date, rating, totalDuration)
}

func InsertDailySleepData(date string, score int64) {
	InsertData("INSERT INTO sleep (date, rating) VALUES ($1, $2) ON CONFLICT (date) DO UPDATE SET rating = EXCLUDED.rating;", date, score)
}
