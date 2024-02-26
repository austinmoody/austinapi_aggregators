package main

import (
	"context"
	"github.com/austinmoody/austinapi_db/austinapi_db"
	"github.com/austinmoody/go_oura"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

func processStress(startDate time.Time, endDate time.Time) {
	client := go_oura.NewClient(OuraPersonalToken)

	log.Printf("Pulling Oura Stress Data")
	results, err := client.GetStresses(startDate, endDate, nil)
	if err != nil {
		log.Fatalf("Failed to get Stress data: %v", err)
	}

	for _, result := range results.Items {
		insertStress(result)
		log.Printf("Processed Stress with ID '%s'", result.ID)
	}

}

func insertStress(stress go_oura.DailyStress) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, DatabaseConnectionString)
	if err != nil {
		log.Fatalf("DB Connection error: %v", err)
	}
	defer conn.Close(ctx)

	apiDb := austinapi_db.New(conn)

	params := austinapi_db.SaveStressParams{
		Date:               stress.Day.Time,
		HighStressDuration: stress.StressHigh,
	}

	err = apiDb.SaveStress(ctx, params)
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}

}
