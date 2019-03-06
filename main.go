package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lindsshldz/itinerary-cli/cli"
	"github.com/lindsshldz/itinerary-cli/db"
	"github.com/lindsshldz/itinerary-cli/itinerary"
)

func main() {

	db, err := db.ConnectDatabase("itinerary_db.config")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	itineraryService := itinerary.NewService(db)

	cliMenu := cli.New(itineraryService)

	fmt.Println()

	cliMenu.MainMenu()

}
