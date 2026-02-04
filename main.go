package main

import (
	"fmt"
	"log"
	"os"
	datastorage "walletGolang/dataStorage"
	"walletGolang/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config.env")

	if err != nil {
		fmt.Println(err)
		return
	}

	host := "localhost"
	port := "5432"
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	db, err := datastorage.NewPostgres(host, port, user, password, dbname)

	if err != nil {
		log.Fatal(err)
		return
	}

	server := server.Server{}

	server.Start(db)

}
