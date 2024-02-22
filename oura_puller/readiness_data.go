// Oura calls it Readiness, for AustinAPI called Preparedness
package main

import (
	"context"
	"github.com/austinmoody/austinapi_db/austinapi_db"
	"github.com/austinmoody/go_oura"
	"github.com/jackc/pgx/v5"
	"log"
	"time"
)

func processReadiness(startDate time.Time, endDate time.Time) {

	client := go_oura.NewClient(OuraPersonalToken)

	log.Printf("Pulling Oura Readiness Data")
	readinessData, err := client.GetReadinesses(startDate, endDate, nil)
	if err != nil {
		log.Fatalf("Failed to get readiness data: %v", err)
	}

	for _, rd := range readinessData.Items {
		insertReadinessData(rd)
		log.Printf("Processed Readiness with ID '%s'", rd.Id)
	}

}

func insertReadinessData(readiness go_oura.DailyReadiness) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, DatabaseConnectionString)
	if err != nil {
		log.Fatalf("DB Connection error: %v", err)
	}
	defer conn.Close(ctx)

	apiDb := austinapi_db.New(conn)

	params := austinapi_db.SavePreparednessParams{
		Date:   readiness.Day.Time,
		Rating: readiness.Score,
	}

	err = apiDb.SavePreparedness(ctx, params)
	if err != nil {
		log.Fatalf("Insert error: %v", err)
	}
}
