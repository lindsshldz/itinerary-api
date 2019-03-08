package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

// func (s *ItineraryServer) GetTripHandler(rw http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	tripIDStr := vars["tripID"]

// 	tripID, err := strconv.Atoi(tripIDStr)
// 	if err != nil {
// 		fmt.Println("Invalid tripID:", err)
// 		rw.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	trip, err := s.itineraryService.GetTrip(tripID)
// 	if err != nil {
// 		fmt.Println("Error getting trip:", err)
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	tripJSON, err := json.Marshal(trip)
// 	if err != nil {
// 		fmt.Println("Error marshaling trip:", err)
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	rw.Header().Add("Content-Type", "application/json")
// 	rw.Write(arcadeJSON)
// }

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
