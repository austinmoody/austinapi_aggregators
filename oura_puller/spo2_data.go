package main

import (
	"context"
	"github.com/austinmoody/austinapi_db/austinapi_db"
	"github.com/austinmoody/go_oura"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

func processSpo2(startDate time.Time, endDate time.Time) {
	client := go_oura.NewClient(OuraPersonalToken)

	log.Printf("Pulling Oura Spo2 Data")
	results, err := client.GetSpo2Readings(startDate, endDate, nil)
	if err != nil {
		log.Fatalf("Failed to get Spo2 data: %v", err)
	}

	for _, result := range results.Items {
		insertSpo2(result)
		log.Printf("Processed Spo2 with ID '%s'", result.ID)
	}

}

func insertSpo2(spo2 go_oura.DailySpo2Reading) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, DatabaseConnectionString)
	if err != nil {
		log.Fatalf("DB Connection error: %v", err)
	}
	defer conn.Close(ctx)

	apiDb := austinapi_db.New(conn)

	params := austinapi_db.SaveSpo2Params{
		Date:        spo2.Day.Time,
		AverageSpo2: spo2.Percentage.Average,
	}

	err = apiDb.SaveSpo2(ctx, params)
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}

}
