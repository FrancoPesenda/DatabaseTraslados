package user

import (
	"errors"
	"time"
)

var (
	ErrEventNameRequired  = errors.New("field 'name' is required")
	ErrLocationIDRequired = errors.New("field 'location_id' is required")
	ErrStartDateRequired  = errors.New("field 'start_date' is required")
	ErrEndDateRequired    = errors.New("field 'end_date' is required")
	ErrInvalidEventDates  = errors.New("start_date must be before or equal to end_date")
)

type Event struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	LocationID int    `json:"location_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Image      string `json:"image,omitempty"`
}

func (e Event) IsNameValid() bool {
	return e.Name != ""
}

func (e Event) IsLocationIDValid() bool {
	return e.LocationID > 0
}

func (e Event) IsStartDateValid() bool {
	return e.StartDate != "" && isValidDate(e.StartDate)
}

func (e Event) IsEndDateValid() bool {
	return e.EndDate != "" && isValidDate(e.EndDate)
}

func (e Event) HasValidDateRange() bool {
	start, err := time.Parse("2006-01-02", e.StartDate)
	if err != nil {
		return false
	}

	end, err := time.Parse("2006-01-02", e.EndDate)
	if err != nil {
		return false
	}

	return !end.Before(start)
}

func isValidDate(value string) bool {
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}
