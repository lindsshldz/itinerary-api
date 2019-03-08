package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lindsshldz/itinerary-api/itinerary"
)

type ItineraryServer struct {
	itineraryService *itinerary.ItineraryService
}

func NewServer(itineraryService *itinerary.ItineraryService) *ItineraryServer {
	return &ItineraryServer{
		itineraryService: itineraryService,
	}
}

func (s *ItineraryServer) ListTripsHandler(rw http.ResponseWriter, r *http.Request) {
	Trips, err := s.itineraryService.ListTrips()
	if err != nil {
		fmt.Println("Error listing Trips:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	TripsJSON, err := json.Marshal(Trips)
	if err != nil {
		fmt.Println("Error marshaling Trips:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.Write(TripsJSON)
}

type CreateTripRequest struct {
	Location  string
	Budget    float64
	StartDate time.Time
	EndDate   time.Time
}

func (s *ItineraryServer) CreateTripHandler(rw http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var newTrip CreateTripRequest
	err = json.Unmarshal(requestBody, &newTrip)
	if err != nil {
		fmt.Println("Error unmarshaling new trip details:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.itineraryService.AddTrip(newTrip.Location, newTrip.Budget, newTrip.StartDate, newTrip.EndDate)
	if err != nil {
		fmt.Println("Error creating new trip:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (s *ItineraryServer) ListDaysAtTripHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripIDStr := vars["tripID"]

	tripID, err := strconv.Atoi(tripIDStr)
	if err != nil {
		fmt.Println("Invalid arcadeID:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	daysAtTrip, err := s.itineraryService.ListDays(tripID)
	if err != nil {
		fmt.Println("Error listing days:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	daysAtTripJSON, err := json.Marshal(daysAtTrip)
	if err != nil {
		fmt.Println("Error marshaling daysAtTrip:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.Write(daysAtTripJSON)
}

type AddDayDetailsRequest struct {
	Location    string
	Activities  string
	Restaurants string
	Hotel       string
}

func (s *ItineraryServer) AddDetailsToDayHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripIDStr := vars["tripID"]

	tripID, err := strconv.Atoi(tripIDStr)
	if err != nil {
		fmt.Println("Invalid tripID:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	dayIDStr := vars["dayID"]

	dayID, err := strconv.Atoi(dayIDStr)
	if err != nil {
		fmt.Println("Invalid dayID:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	var newDayDetails AddDayDetailsRequest
	err = json.Unmarshal(requestBody, &newDayDetails)
	if err != nil {
		fmt.Println("Error unmarshaling day details:", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	day := itinerary.Day{
		ID:          dayID,
		Location:    newDayDetails.Location,
		Activities:  newDayDetails.Activities,
		Restaurants: newDayDetails.Restaurants,
		Hotel:       newDayDetails.Hotel,
		TripID:      tripID,
	}

	err = s.itineraryService.UpdateDetails(day)
	if err != nil {
		fmt.Println("Error adding details to day:", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
