package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/lindsshldz/itinerary-api/db"
	"github.com/lindsshldz/itinerary-api/itinerary"
	"github.com/lindsshldz/itinerary-api/server"
)

const port = ":8000"

func main() {

	db, err := db.ConnectDatabase("itinerary_db.config")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	itineraryService := itinerary.NewService(db)

	itineraryServer := server.NewServer(itineraryService)

	router := mux.NewRouter()
	router.HandleFunc("/trips/{tripID}/days", itineraryServer.ListDaysAtTripHandler).Methods("GET")
	router.HandleFunc("/trips/{tripID}/days/{dayID}", itineraryServer.AddDetailsToDayHandler).Methods("PUT")
	router.HandleFunc("/trips", itineraryServer.CreateTripHandler).Methods("POST")
	router.HandleFunc("/trips", itineraryServer.ListTripsHandler).Methods("GET")

	http.Handle("/", router)

	fmt.Println("Waiting for requests on port:", port)
	http.ListenAndServe(port, nil)

}
