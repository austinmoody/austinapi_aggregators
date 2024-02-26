package main

import (
	"context"
	"github.com/austinmoody/austinapi_db/austinapi_db"
	"github.com/austinmoody/go_oura"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

func processHeartRate(startDate time.Time, endDate time.Time) {
	for d := startDate; d.After(endDate) == false; d = d.AddDate(0, 0, 1) {
		processHeartRateSingleDate(d)
	}
}

func processHeartRateSingleDate(startDate time.Time) {
	// only to be used for 1 day at a time
	client := go_oura.NewClient(OuraPersonalToken)
	endDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 23, 23, 59, 0, startDate.Location())
	low := 0
	high := 0
	count := 0
	total := 0

	// Get HeartRate from API, has information about durations & start/end time etc...
	log.Printf("Pulling Heart Rate items for %s", startDate.Format("2006-01-02"))
	heartrates, err := client.GetHeartRates(startDate, endDate, nil)
	if err != nil {
		log.Printf("Error getting Heart Rate Items: %v\n", err)
		return
	}

	count = len(heartrates.Items)

	for _, hr := range heartrates.Items {
		total += hr.Bpm

		if hr.Bpm > high {
			high = hr.Bpm
		}

		if low == 0 || hr.Bpm < low {
			low = hr.Bpm
		}
	}
	log.Printf("Lowest Heart Rate: %d", low)
	log.Printf("Highest Heart Rate: %d", high)
	log.Printf("Total Heart Rate Items: %d", count)
	log.Printf("Average Heart Rate: %d", total/count)

	params := austinapi_db.SaveHeartRateParams{
		Date:    startDate,
		Low:     low,
		High:    high,
		Average: total / count,
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, DatabaseConnectionString)
	if err != nil {
		log.Fatalf("DB Connection error: %v", err)
	}
	defer conn.Close(ctx)

	apiDb := austinapi_db.New(conn)

	err = apiDb.SaveHeartRate(ctx, params)
	if err != nil {
		log.Fatalf("Error saving Heart Rate to DB: %v", err)
	}

}
