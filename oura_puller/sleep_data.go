package main

import (
	"context"
	"github.com/austinmoody/austinapi_db/austinapi_db"
	"github.com/austinmoody/go_oura"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func processSleep(startDate time.Time, endDate time.Time) {

	// Something to think about... specifying a start & end a day apart.
	// so Start = 2024-01-29 & End = 2024-01-30
	// Sleep will bring back 1 item dated 2024-01-29
	// Daily Sleep will bring back 1 items: 2024-01-29 and 2024-01-30

	// Pull Oura token from environment
	//ouraPersonalToken := os.Getenv("OURA_PERSONAL_TOKEN")

	// Here we would use oura_go to pull sleep data and parse out what we need
	client := go_oura.NewClient(OuraPersonalToken)

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
			sleep.TotalSleep += sleepDoc.TotalSleepDuration
			sleep.LightSleep += sleepDoc.LightSleepDuration
			sleep.DeepSleep += sleepDoc.DeepSleepDuration
			sleep.RemSleep += sleepDoc.RemSleepDuration
			sleepData[sleepDoc.Day.Time] = sleep
		} else {
			// Add new
			sleepData[sleepDoc.Day.Time] = austinapi_db.Sleep{
				Date:       sleepDoc.Day.Time,
				TotalSleep: sleepDoc.TotalSleepDuration,
				LightSleep: sleepDoc.LightSleepDuration,
				DeepSleep:  sleepDoc.DeepSleepDuration,
				RemSleep:   sleepDoc.RemSleepDuration,
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
			sleep.Rating = dailySleep.Score
			sleepData[dailySleep.Day.Time] = sleep
		} else {
			// Add new
			sleepData[dailySleep.Day.Time] = austinapi_db.Sleep{
				Date:   dailySleep.Day.Time,
				Rating: dailySleep.Score,
			}
		}

		log.Printf("Processed Daily Sleep with ID: %s\n", dailySleep.ID)
	}

	insertSleepData(sleepData)
}

func insertSleepData(sleepData map[time.Time]austinapi_db.Sleep) {

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, DatabaseConnectionString)
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
