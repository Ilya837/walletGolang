package main

import (
	"fmt"
	"log"
	"os"
	datastorage "walletGolang/dataStorage"
	"walletGolang/server"

	"github.com/joho/godotenv"
)

func startServer() {
	err := godotenv.Load("config.env")

	if err != nil {
		fmt.Println(err)
		return
	}

	host := os.Getenv("POSTGRES_HOST")
	dbPort := "5432"
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	db, err := datastorage.NewPostgres(host, dbPort, user, password, dbName)

	if err != nil {
		log.Fatal(err)
		return
	}

	server := server.Server{}

	servePort := os.Getenv("SERVER_PORT")

	server.Start(db, servePort)
}

func main() {
	startServer()
}
