package cli

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lindsshldz/itinerary-cli/itinerary"
	"github.com/manifoldco/promptui"
)

type CLI struct {
	itineraryService *itinerary.ItineraryService
}

func New(itineraryService *itinerary.ItineraryService) *CLI {
	return &CLI{
		itineraryService: itineraryService,
	}
}

const (
	addTripCmd    = "Add Trip"
	addDetailsCmd = "Add fun details to your trip"
	printTripCmd  = "Print out your trip itinerary"

	tripDayDescriptionTemplate = `
%s: %s
	You have planned to %s, eat at %s, and sleep at %s.
		`
)

func (c *CLI) MainMenu() {

	for {
		fmt.Println()

		prompt := promptui.Select{
			Label: "Select Action",
			Items: []string{
				addTripCmd,
				addDetailsCmd,
				printTripCmd,
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case addTripCmd:
			err := c.addTripPrompt()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

		case addDetailsCmd:
			err := c.addDetails()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

		case printTripCmd:
			err := c.printItinerary()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

}

func (c *CLI) addTripPrompt() error {
	locationPrompt := promptui.Prompt{
		Label: "Location",
	}
	location, err := locationPrompt.Run()
	if err != nil {
		return err
	}

	tripStartDate, err := datePromptHelper("Trip Start Date")
	if err != nil {
		return err
	}

	tripEndDate, err := datePromptHelper("Trip End Date")
	if err != nil {
		return err
	}

	tripBudget, err := numberPromptHelper("Overall Budget")
	if err != nil {
		return err
	}

	err = c.itineraryService.AddTrip(location, tripBudget, tripStartDate, tripEndDate)
	if err != nil {
		return err
	}

	fmt.Println("Added " + location + " trip!")

	return nil
}

func (c *CLI) addDetails() error {
	availableTrips, err := c.itineraryService.ListTrips()
	if err != nil {
		return err
	}

	if len(availableTrips) == 0 {
		fmt.Println("Create a trip first")
		return nil
	}
	var options []string
	for _, trip := range availableTrips {
		options = append(options, trip.Location)
	}

	selectTripPrompt := promptui.Select{
		Label: "Select Trip",
		Items: options,
	}

	chosenIndex, _, err := selectTripPrompt.Run()
	if err != nil {
		return err
	}

	chosenTrip := availableTrips[chosenIndex]

	availableDays, err := c.itineraryService.ListDays(chosenTrip.ID)
	if err != nil {
		return err
	}

	if len(availableDays) == 0 {
		fmt.Println("Must add start and end dates")
		return nil
	}
	var dayOptions []string
	for _, day := range availableDays {
		dayOptions = append(dayOptions, day.Date.Format("01-02-2006"))
	}

	selectDayPrompt := promptui.Select{
		Label: "Select Day",
		Items: dayOptions,
	}

	chosenDayIndex, _, err := selectDayPrompt.Run()
	if err != nil {
		return err
	}

	chosenDay := availableDays[chosenDayIndex]

	dayLocation, err := detailPromptHelper("Location on that day")
	if err != nil {
		return err
	}

	dayActivities, err := detailPromptHelper("Add activities for the day")
	if err != nil {
		return err
	}

	dayRestaurant, err := detailPromptHelper("Where you will eat that day")
	if err != nil {
		return err
	}

	dayHotel, err := detailPromptHelper("Hotel name for that night")
	if err != nil {
		return err
	}

	chosenDay.Location = dayLocation
	chosenDay.Activities = dayActivities
	chosenDay.Restaurants = dayRestaurant
	chosenDay.Hotel = dayHotel

	err = c.itineraryService.UpdateDetails(chosenDay)
	if err != nil {
		return err
	}

	fmt.Println("Added details to " + chosenDay.Location)

	return nil
}

func (c *CLI) printItinerary() error {

	availableTrips, err := c.itineraryService.ListTrips()
	if err != nil {
		return err
	}

	if len(availableTrips) == 0 {
		fmt.Println("Create a trip first")
		return nil
	}
	var options []string
	for _, trip := range availableTrips {
		options = append(options, trip.Location)
	}

	selectTripPrompt := promptui.Select{
		Label: "Select Trip",
		Items: options,
	}

	chosenIndex, _, err := selectTripPrompt.Run()
	if err != nil {
		return err
	}

	chosenTrip := availableTrips[chosenIndex]

	days, err := c.itineraryService.ListDays(chosenTrip.ID)
	if err != nil {
		return err
	}

	fmt.Println("*" + chosenTrip.Location + " Trip*")
	fmt.Println("Overall Budget: $" + strconv.FormatFloat(chosenTrip.Budget, 'f', 2, 64))
	for _, day := range days {
		fmt.Printf(tripDayDescriptionTemplate, day.Date.Format("01-02-2006"), day.Location, day.Activities, day.Restaurants, day.Hotel)
	}

	return nil

}

func datePromptHelper(label string) (time.Time, error) {

	validate := func(input string) error {
		_, err := time.Parse("01-02-2006", input)
		if err != nil {
			return errors.New("Needs to be in full date format: mm-dd-yyyy")
		}
		return nil
	}
	datePrompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	dateStr, err := datePrompt.Run()
	if err != nil {
		return time.Time{}, err
	}
	date, err := time.Parse("01-02-2006", dateStr)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func detailPromptHelper(label string) (string, error) {

	validate := func(input string) error {
		if len(input) > 20 {
			return errors.New("can only have 20 characters")
		}
		return nil
	}
	detailPrompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	detail, err := detailPrompt.Run()
	if err != nil {
		return detail, err
	}
	return detail, nil
}

func numberPromptHelper(label string) (float64, error) {

	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("needs to be a number")
		}
		return nil
	}
	numberPrompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	numberStr, err := numberPrompt.Run()
	if err != nil {
		return 0, err
	}
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}
