package main

import (
	"database/sql"
	"fmt"
	"github.com/austinmoody/go_oura"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

func ProcessSleep(personalToken string, startDate time.Time, endDate time.Time) {

	// TODO: Something to think about... specifying a start & end a day apart.
	// so Start = 2024-01-29 & End = 2024-01-30
	// Sleep will bring back 1 item dated 2024-01-29
	// Daily Sleep will bring back 1 items: 2024-01-29 and 2024-01-30

	// Here we would use oura_go to pull sleep data and parse out what we need
	client := go_oura.NewClient(personalToken)

	// Get Sleep from API, has information about durations & start/end time etc...
	sleepDocs, err := client.GetSleeps(startDate, endDate, nil)
	if err != nil {
		fmt.Printf("Error getting Sleep Items: %v\n", err)
		return
	}

	for _, sleepDoc := range sleepDocs.Items {
		InsertSleepData(sleepDoc.Day.Format("2006-01-02"), 0, sleepDoc.TotalSleepDuration)
	}

	// Get Daily Sleep from API, has the Score
	dailySleeps, err := client.GetDailySleeps(startDate, endDate, nil)
	if err != nil {
		fmt.Printf("Error getting Daily Sleep items: %v\n", err)
	}

	for _, dailySleep := range dailySleeps.Items {
		InsertDailySleepData(dailySleep.Day.Format("2006-01-02"), dailySleep.Score)
	}
}

func GetDatabaseConnectionString() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseUser := os.Getenv("DATABASE_USER")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	sslMode := os.Getenv("DATABASE_SSLMODE")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		databaseHost,
		databasePort,
		databaseUser,
		databasePassword,
		databaseName,
		sslMode,
	)
}

func InsertSleepData(date string, rating int, totalDuration int) {
	connStr := GetDatabaseConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to close the connection!
	defer db.Close()

	// Prepares statement for inserting data
	stmt, err := db.Prepare("INSERT INTO sleep (date, rating, total_duration) VALUES ($1, $2, $3) ON CONFLICT (date) DO UPDATE SET rating = EXCLUDED.rating;")
	if err != nil {
		log.Fatal(err)
	}

	// Inserts our data
	_, err = stmt.Exec(date, rating, totalDuration)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sleep data successfully inserted!")
}

func InsertDailySleepData(date string, score int64) {
	connStr := GetDatabaseConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to close the connection!
	defer db.Close()

	// Prepares statement for inserting data
	stmt, err := db.Prepare("INSERT INTO sleep (date, rating) VALUES ($1, $2) ON CONFLICT (date) DO UPDATE SET rating = EXCLUDED.rating;")
	if err != nil {
		log.Fatal(err)
	}

	// Inserts our data
	_, err = stmt.Exec(date, score)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Daily Sleep data successfully inserted!")

}
