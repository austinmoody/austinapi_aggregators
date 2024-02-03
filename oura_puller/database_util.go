package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

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

func InsertData(query string, args ...interface{}) {
	connStr := GetDatabaseConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to close the connection!
	defer db.Close()

	// Prepares statement for inserting data
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	// Inserts our data
	_, err = stmt.Exec(args...)
	if err != nil {
		log.Fatal(err)
	}

}
