package main

import (
	"context"
	"github.com/austinmoody/austinapi_db/austinapi_db"
	"github.com/austinmoody/go_oura"
	"github.com/jackc/pgx/v5"
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
	log.Printf("Pulling Sleep Items")
	sleepDocs, err := client.GetSleeps(startDate, endDate, nil)
	if err != nil {
		log.Printf("Error getting Sleep Items: %v\n", err)
		return
	}

	sleepData := make(map[time.Time]austinapi_db.Sleep)

	for _, sleepDoc := range sleepDocs.Items {

		sleep, ok := sleepData[sleepDoc.Day.Time]
		if ok {
			// This item already exists in map, add values
			sleep.TotalSleep += int32(sleepDoc.TotalSleepDuration)
			sleep.LightSleep += int32(sleepDoc.LightSleepDuration)
			sleep.DeepSleep += int32(sleepDoc.DeepSleepDuration)
			sleep.RemSleep += int32(sleepDoc.RemSleepDuration)
			sleepData[sleepDoc.Day.Time] = sleep
		} else {
			// Add new
			sleepData[sleepDoc.Day.Time] = austinapi_db.Sleep{
				Date:       sleepDoc.Day.Time,
				TotalSleep: int32(sleepDoc.TotalSleepDuration), // TODO change go_oura or db to match
				LightSleep: int32(sleepDoc.LightSleepDuration),
				DeepSleep:  int32(sleepDoc.DeepSleepDuration),
				RemSleep:   int32(sleepDoc.RemSleepDuration),
			}
		}

		log.Printf("Processed Sleep with ID: %s\n", sleepDoc.ID)
	}

	// Get Daily Sleep from API, has the Score
	log.Printf("Pulling Daily Sleep Items\n")
	dailySleeps, err := client.GetDailySleeps(startDate, endDate, nil)
	if err != nil {
		log.Printf("Error getting Daily Sleep items: %v\n", err)
	}

	for _, dailySleep := range dailySleeps.Items {

		sleep, ok := sleepData[dailySleep.Day.Time]
		if ok {
			// This item already exists in map, add values
			sleep.Rating = int32(dailySleep.Score) // TODO change go_oura or db to match
			sleepData[dailySleep.Day.Time] = sleep
		} else {
			// Add new
			sleepData[dailySleep.Day.Time] = austinapi_db.Sleep{
				Date:   dailySleep.Day.Time,
				Rating: int32(dailySleep.Score),
			}
		}

		log.Printf("Processed Daily Sleep with ID: %s\n", dailySleep.ID)
	}

	InsertSleepData(sleepData)
}

func InsertSleepData(sleepData map[time.Time]austinapi_db.Sleep) {

	connStr := GetDatabaseConnectionString()
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("DB Connection error: %v", err)
	}
	defer conn.Close(ctx)

	apiDb := austinapi_db.New(conn)

	for _, sleep := range sleepData {

		// Don't save if one of these is missing... shouldn't happen
		if sleep.Rating == 0 || sleep.TotalSleep == 0 {
			continue
		}

		params := austinapi_db.SaveSleepParams{
			Date:       sleep.Date,
			Rating:     sleep.Rating,
			TotalSleep: sleep.TotalSleep,
			LightSleep: sleep.LightSleep,
			DeepSleep:  sleep.DeepSleep,
			RemSleep:   sleep.RemSleep,
		}

		err = apiDb.SaveSleep(ctx, params)
		if err != nil {
			log.Fatalf("Insert error: %v", err)
		}
	}
}
