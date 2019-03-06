package itinerary

import (
	"database/sql"
	"time"
)

var (
	trips []*Trip
)

type Trip struct {
	ID        int
	Location  string
	StartDate time.Time
	EndDate   time.Time
	Budget    float64
}

type Day struct {
	ID          int
	Date        time.Time
	Location    string
	Activities  string
	Restaurants string
	Hotel       string
	TripID      int
}

type ItineraryService struct {
	db *sql.DB
}

func NewService(db *sql.DB) *ItineraryService {
	return &ItineraryService{
		db: db,
	}
}

const (
	insertTripQuery = `INSERT INTO trips (trip_name, budget, start_date, end_date) VALUES (?, ?, ?, ?)`

	selectLastID = `SELECT LAST_INSERT_ID()`

	createDetailsQuery = "INSERT INTO details (date, trip_id) VALUES (?, ?)"

	selectTripsQuery = "SELECT id, trip_name, budget, start_date, end_date FROM trips"

	selectDetailsQuery = "SELECT id, date, day_location, activities, restaurants, hotel, trip_id FROM details WHERE trip_id = ?"

	updateDetailsQuery = "UPDATE details SET day_location = ?, activities = ?, restaurants = ?, hotel = ? WHERE id = ?"
)

func (i *ItineraryService) AddTrip(tripName string, budget float64, startDate time.Time, endDate time.Time) error {

	trxn, err := i.db.Begin()
	if err != nil {
		trxn.Rollback()
		return err
	}

	_, err = trxn.Exec(insertTripQuery, tripName, budget, startDate, endDate)
	if err != nil {
		trxn.Rollback()
		return err
	}

	row := trxn.QueryRow(selectLastID)

	var tripID int
	err = row.Scan(
		&tripID,
	)
	if err != nil {
		trxn.Rollback()
		return err
	}

	tripDays := countDays(tripID, startDate, endDate)

	for _, day := range tripDays {
		_, err := trxn.Exec(createDetailsQuery, day.Date, day.TripID)
		if err != nil {
			trxn.Rollback()
			return err
		}
	}

	err = trxn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (i *ItineraryService) ListTrips() ([]Trip, error) {
	rows, err := i.db.Query(selectTripsQuery)
	if err != nil {
		return nil, err
	}

	var trips []Trip
	for rows.Next() {
		var trip Trip

		err := rows.Scan(
			&trip.ID,
			&trip.Location,
			&trip.Budget,
			&trip.StartDate,
			&trip.EndDate,
		)
		if err != nil {
			return nil, err
		}

		trips = append(trips, trip)
	}

	return trips, nil
}

func (i *ItineraryService) ListDays(tripID int) ([]Day, error) {
	rows, err := i.db.Query(selectDetailsQuery, tripID)
	if err != nil {
		return nil, err
	}

	var days []Day
	for rows.Next() {
		var day Day

		err := rows.Scan(
			&day.ID,
			&day.Date,
			&day.Location,
			&day.Activities,
			&day.Restaurants,
			&day.Hotel,
			&day.TripID,
		)

		if err != nil {
			return nil, err
		}

		days = append(days, day)
	}

	return days, nil
}

func countDays(tripID int, startDate time.Time, endDate time.Time) []Day {

	var days []Day
	currentDate := startDate

	for !currentDate.After(endDate) {

		newDay := Day{
			Date:   currentDate,
			TripID: tripID,
		}

		days = append(days, newDay)

		currentDate = currentDate.Add(24 * time.Hour)
	}
	return days
}

func (i *ItineraryService) UpdateDetails(day Day) error {
	_, err := i.db.Exec(updateDetailsQuery, day.Location, day.Activities, day.Restaurants, day.Hotel, day.ID)

	if err != nil {
		return err
	}

	return nil
}
